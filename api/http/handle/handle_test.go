package handle

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/ntdat104/go-clean-architecture/api/error_code"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	response := NewResponse(c)

	assert.NotNil(t, response)
	assert.Equal(t, c, response.Ctx)
}

func TestResponse_ToResponse(t *testing.T) {
	testCases := []struct {
		name           string
		data           interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "test1",
			data:           map[string]interface{}{"key": "value"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{"key":"value"}}`,
		},
		{
			name:           "test nil",
			data:           nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			response := NewResponse(c)

			response.ToResponse(tc.data)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestResponse_ToResponseList(t *testing.T) {
	testCases := []struct {
		name           string
		list           interface{}
		totalRows      int
		page           int
		pageSize       int
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "返回列表数据",
			list:           []map[string]interface{}{{"id": 1}, {"id": 2}},
			totalRows:      10,
			page:           1,
			pageSize:       2,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"code":0,"message":"success","data":{"list":[{"id":1},{"id":2}],"pager":{"page":1,"page_size":2,"total_rows":10}}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?page="+
				""+strconv.Itoa(tc.page)+"&page_size="+strconv.Itoa(tc.pageSize), nil)

			c.Set("page", tc.page)
			c.Set("page_size", tc.pageSize)

			c.JSON(http.StatusOK, StandardResponse{
				Code:    0,
				Message: "success",
				Data: gin.H{
					"list": tc.list,
					"pager": gin.H{
						"page":       tc.page,
						"page_size":  tc.pageSize,
						"total_rows": tc.totalRows,
					},
				},
			})

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestResponse_ToErrorResponse(t *testing.T) {
	testCases := []struct {
		name           string
		err            *error_code.Error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "test 1",
			err:            error_code.ServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":10000,"message":"server internal error"}`,
		},
		{
			name:           "test 2",
			err:            error_code.InvalidParams.WithDetails("字段不能为空"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":10001,"message":"invalid params","data":{"details":["字段不能为空"]}}`,
		},
		{
			name:           "test 3",
			err:            &error_code.Error{Code: 10002, Msg: "msg", HTTP: http.StatusBadRequest, DocRef: "https://example.com/docs"},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"code":10002,"message":"msg","doc_ref":"https://example.com/docs"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			response := NewResponse(c)

			response.ToErrorResponse(tc.err)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]interface{}{"test": "value"}
	Success(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"code":0,"message":"success","data":{"test":"value"}}`, w.Body.String())
}

func TestError(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "API notfound",
			err:            error_code.NotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"code":10002,"message":"record not found","data":{"details":null}}`,
		},
		{
			name:           "Internal server error",
			err:            assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":10000,"message":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tc.err)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestGetQueryInt(t *testing.T) {
	testCases := []struct {
		name         string
		queryParam   string
		defaultValue int
		expectedInt  int
	}{
		{
			name:         "有效整数参数",
			queryParam:   "?id=123",
			defaultValue: 0,
			expectedInt:  123,
		},
		{
			name:         "无效整数参数",
			queryParam:   "?id=abc",
			defaultValue: 0,
			expectedInt:  0,
		},
		{
			name:         "无参数",
			queryParam:   "",
			defaultValue: 10,
			expectedInt:  10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test"+tc.queryParam, nil)

			result := GetQueryInt(c, "id", tc.defaultValue)

			assert.Equal(t, tc.expectedInt, result)
		})
	}
}

func TestGetQueryString(t *testing.T) {
	testCases := []struct {
		name          string
		queryParam    string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "有参数",
			queryParam:    "?name=test",
			defaultValue:  "default",
			expectedValue: "test",
		},
		{
			name:          "无参数",
			queryParam:    "",
			defaultValue:  "default",
			expectedValue: "default",
		},
		{
			name:          "空参数",
			queryParam:    "?name=",
			defaultValue:  "default",
			expectedValue: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test"+tc.queryParam, nil)

			result := GetQueryString(c, "name", tc.defaultValue)

			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

func TestIsNil(t *testing.T) {
	var nilPtr *string
	var nilSlice []string
	var nilMap map[string]string
	var nilChan chan string
	var nilInterface interface{} = nil

	nonNilPtr := new(string)
	nonNilSlice := make([]string, 0)
	nonNilMap := make(map[string]string)
	nonNilChan := make(chan string)

	testCases := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "nil值",
			value:    nil,
			expected: true,
		},
		{
			name:     "nil指针",
			value:    nilPtr,
			expected: true,
		},
		{
			name:     "nil切片",
			value:    nilSlice,
			expected: true,
		},
		{
			name:     "nil映射",
			value:    nilMap,
			expected: true,
		},
		{
			name:     "nil通道",
			value:    nilChan,
			expected: true,
		},
		{
			name:     "nil接口",
			value:    nilInterface,
			expected: true,
		},
		{
			name:     "非nil指针",
			value:    nonNilPtr,
			expected: false,
		},
		{
			name:     "非nil切片",
			value:    nonNilSlice,
			expected: false,
		},
		{
			name:     "非nil映射",
			value:    nonNilMap,
			expected: false,
		},
		{
			name:     "非nil通道",
			value:    nonNilChan,
			expected: false,
		},
		{
			name:     "整数值",
			value:    1,
			expected: false,
		},
		{
			name:     "字符串值",
			value:    "test",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsNil(tc.value)

			assert.Equal(t, tc.expected, result)
		})
	}
}
