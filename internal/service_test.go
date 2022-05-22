package internal

import (
	"net/http"
	"net/http/httptest"
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

	assert.Equal(t, 200, recorder.Code)
}
