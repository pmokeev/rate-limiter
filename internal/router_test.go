package internal

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/magiconair/properties/assert"
)

func CreateMockAddr() string {
	server, err := miniredis.Run()
	if err != nil {
		return ""
	}
	return server.Addr()
}

func TestSingleRate(t *testing.T) {
	mockAddr := CreateMockAddr()
	if mockAddr == "" {
		t.Fatal("Error while creating mock redis address")
	}

	service := NewService(mockAddr)
	router := NewRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
	request.Header.Add("X-Forwarded-For", "10.10.10.10")

	router.InitRouter().ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), "99")
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
}

func TestOneHundredRequests(t *testing.T) {
	mockAddr := CreateMockAddr()
	if mockAddr == "" {
		t.Fatal("Error while creating mock redis address")
	}

	service := NewService(mockAddr)
	router := NewRouter(service)

	for i := 1; i < 100; i++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
		request.Header.Add("X-Forwarded-For", "10.10.10.10")

		router.InitRouter().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), strconv.Itoa(100-i))
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
	}
}

func TestOneHundredRequestsAndOneMore(t *testing.T) {
	mockAddr := CreateMockAddr()
	if mockAddr == "" {
		t.Fatal("Error while creating mock redis address")
	}

	service := NewService(mockAddr)
	router := NewRouter(service)

	for i := 1; i < 100; i++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
		request.Header.Add("X-Forwarded-For", "10.10.10.10")

		router.InitRouter().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), strconv.Itoa(100-i))
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
	request.Header.Add("X-Forwarded-For", "10.10.10.10")

	router.InitRouter().ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), "0")
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
}

func TestCheckStatusCodeAfterClearRate(t *testing.T) {
	mockAddr := CreateMockAddr()
	if mockAddr == "" {
		t.Fatal("Error while creating mock redis address")
	}

	service := NewService(mockAddr)
	router := NewRouter(service)

	for i := 1; i < 100; i++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
		request.Header.Add("X-Forwarded-For", "10.10.10.10")

		router.InitRouter().ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), strconv.Itoa(100-i))
		assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
	}

	recorderClear := httptest.NewRecorder()
	requestClear := httptest.NewRequest(http.MethodPost, "/api/clear", nil)
	requestClear.Header.Add("X-Forwarded-For", "10.10.10.10")

	router.InitRouter().ServeHTTP(recorderClear, requestClear)

	assert.Equal(t, http.StatusOK, recorderClear.Code)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/send", nil)
	request.Header.Add("X-Forwarded-For", "10.10.10.10")

	router.InitRouter().ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Remaining"), strconv.Itoa(99))
	assert.Equal(t, recorder.Header().Get("X-Ratelimit-Limit"), "100")
}
