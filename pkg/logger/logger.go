package logger

import (
	"errors"
	"os"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// DefaultLogFileSizeInMB Default log file size with Megabyte unit.
	DefaultLogFileSizeInMB = 512
)

var globalLogger atomic.Value

// Level map between string level with zapcore level.
type Level string

const (
	// DEBUG logs.
	DEBUG Level = "DEBUG"
	// INFO level is the default logging.
	INFO Level = "INFO"
	// WARN level logs are more important than Info.
	WARN Level = "WARN"
	// ERROR logs are high-priority.
	ERROR Level = "ERROR"
	// FATAL log message and then calls os.Exit(1).
	FATAL Level = "FATAL"
)

var levelMap = map[Level]zapcore.LevelEnabler{
	DEBUG: zapcore.DebugLevel,
	INFO:  zapcore.InfoLevel,
	WARN:  zapcore.WarnLevel,
	ERROR: zapcore.ErrorLevel,
	FATAL: zapcore.FatalLevel,
}

// Config allows users to configure log level and log file.
type Config struct {
	Level         Level
	FileLogConfig FileLogConfig
}

// FileLogConfig allows users to configure detail log file such as file path, max size of file, max file to backup,....
type FileLogConfig struct {
	IsUseFile  bool
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// NewDefaultConfig returns the default config with INFO level and log to console.
func NewDefaultConfig() Config {
	return Config{
		Level: INFO,
		FileLogConfig: FileLogConfig{
			IsUseFile: false,
		},
	}
}

// NewProductionConfig returns the production config with INFO level.
func NewProductionConfig(isUseFile bool, filePath string) Config {
	// Check if a filePath is provided and if we should use a file for logging.
	if isUseFile && filePath != "" {
		// Get the current time.
		currentTime := time.Now()
		// Format the time as yyyy-mm-dd.
		formattedDate := currentTime.Format("2006-01-02")
		// Construct the final file path with the date.
		// For example, if filePath is "/var/logs/", the new path will be "/var/logs/2025-08-25.log".
		filePath = filePath + formattedDate + ".log"
	}

	return Config{
		Level: INFO,
		FileLogConfig: FileLogConfig{
			IsUseFile:  isUseFile,
			FilePath:   filePath,
			MaxSize:    DefaultLogFileSizeInMB,
			MaxBackups: 0,
			MaxAge:     0,
			Compress:   true,
		},
	}
}

// Init creates the global logger that is used everywhere in project with your config.
func Init(cfg Config) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}

	// Use the new functions to build the components.
	cores, err := buildCores(cfg)
	if err != nil {
		return err
	}
	zapOptions := buildZapOptions()

	// Combine cores and create the logger with the options.
	logger := zap.New(zapcore.NewTee(cores...), zapOptions...)
	globalLogger.Store(logger)

	return nil
}

// buildCores constructs and returns a slice of zapcore.Core instances based on the provided configuration.
// It includes a console core and a file core if file logging is enabled.
func buildCores(cfg Config) ([]zapcore.Core, error) {
	var cores []zapcore.Core

	// Create a core for console output.
	// This core uses a console encoder and always logs at the DebugLevel.
	consoleEncoder := getConsoleEncoder()
	consoleSyncer, err := getConsolLogSyncer()
	if err != nil {
		return nil, err
	}
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel))

	// If file logging is enabled in the configuration, create and add a file core.
	if cfg.FileLogConfig.IsUseFile {
		fileSyncer, err := getFileLogSyncer(cfg.FileLogConfig)
		if err != nil {
			return nil, err
		}
		fileEncoder := getFileEncoder()
		// The file core uses the log level specified in the configuration.
		cores = append(cores, zapcore.NewCore(fileEncoder, fileSyncer, getLevel(cfg.Level)))
	}

	return cores, nil
}

// buildZapOptions returns a slice of zap.Option instances for configuring the logger.
// It adds options for caller information, stack traces, and skipping a caller frame.
func buildZapOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
}

// InitProduction creates the global logger that is used everywhere in project with production config logger.
func InitProduction(filePath string) error {
	isUseFile := len(filePath) > 0
	return Init(NewProductionConfig(isUseFile, filePath))
}

// Sync flushs any buffered log entries. It should be call before program exit.
func Sync() error {
	return getGlobalLog().Sync()
}

// Debug logs a message at Debug level.
func Debug(msg string, fields ...zap.Field) {
	getGlobalLog().Debug(msg, fields...)
}

// Info logs a message at Info level.
func Info(msg string, fields ...zap.Field) {
	getGlobalLog().Info(msg, fields...)
}

// Error logs a message at Error level.
func Error(msg string, fields ...zap.Field) {
	getGlobalLog().Error(msg, fields...)
}

// Warn logs a message at Warn level.
func Warn(msg string, fields ...zap.Field) {
	getGlobalLog().Warn(msg, fields...)
}

// Fatal logs a message at Fatal level.
func Fatal(msg string, fields ...zap.Field) {
	getGlobalLog().Fatal(msg, fields...)
}

func getGlobalLog() *zap.Logger {
	return globalLogger.Load().(*zap.Logger)
}

func setGlobalLog(logger *zap.Logger) {
	globalLogger.Store(logger)
}

func validateConfig(cfg Config) error {
	if len(cfg.Level) == 0 {
		return errors.New("missing level logger")
	}

	fileLogCfg := cfg.FileLogConfig
	if fileLogCfg.IsUseFile && len(fileLogCfg.FilePath) == 0 {
		return errors.New("file path must be not empty")
	}

	if fileLogCfg.MaxAge < 0 {
		return errors.New("MaxAge must be greater than or equal to 0")
	}

	if fileLogCfg.MaxBackups < 0 {
		return errors.New("MaxBackups must be greater than or equal to 0")
	}

	if fileLogCfg.MaxSize < 0 {
		return errors.New("MaxSize must be greater than or equal to 0")
	}

	return nil
}

func getLevel(level Level) zapcore.LevelEnabler {
	zapLevel, ok := levelMap[level]
	if !ok {
		return zapcore.InfoLevel
	}

	return zapLevel
}

// getFileEncoder returns a JSON encoder for file logging.
func getFileEncoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	})
	return zapcore.NewJSONEncoder(cfg)
}

// getConsoleEncoder returns a console encoder for colored console output.
func getConsoleEncoder() zapcore.Encoder {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	})
	return zapcore.NewConsoleEncoder(cfg)
}

func getFileLogSyncer(cfg FileLogConfig) (zapcore.WriteSyncer, error) {
	if st, err := os.Stat(cfg.FilePath); err == nil {
		if st.IsDir() {
			return nil, errors.New("file path is invalid")
		}
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	if cfg.MaxSize == 0 {
		lumberJackLogger.MaxSize = DefaultLogFileSizeInMB
	}

	return zapcore.AddSync(lumberJackLogger), nil
}

func getConsolLogSyncer() (zapcore.WriteSyncer, error) {
	writer, _, err := zap.Open([]string{"stdout"}...)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
