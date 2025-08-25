package repo

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMySQLContainer(t *testing.T) {
	ctx := context.Background()

	// Use a dedicated network for isolation and better connectivity
	networkName := "test-network"
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name: networkName,
		},
	})
	if err != nil {
		t.Fatalf("failed to create network: %v", err)
	}
	defer network.Remove(ctx)

	// Start MySQL container
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		Networks:     []string{networkName},
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "user",
			"MYSQL_PASSWORD":      "pass",
		},
		// Use a more robust health check wait strategy
		WaitingFor: wait.ForListeningPort("3306").
			WithStartupTimeout(1 * time.Minute),
	}
	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer mysqlC.Terminate(ctx)

	// Get container host/port
	host, _ := mysqlC.Host(ctx)
	port, _ := mysqlC.MappedPort(ctx, "3306")
	dsn := fmt.Sprintf("user:pass@tcp(%s:%s)/testdb?parseTime=true", host, port.Port())

	// Open the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Wait for the database to be reachable using Ping
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	// Read and execute the schema file
	schema, err := ioutil.ReadFile("../schema/mysql_schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		t.Fatalf("failed to apply schema: %v", err)
	}

	// Insert test data
	_, err = db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Verify
	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		t.Fatalf("failed to query users: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}

	log.Println("✅ MySQL container test passed")
}
