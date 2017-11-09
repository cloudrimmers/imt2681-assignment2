package httperror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrors(t *testing.T) {
	var r *httptest.ResponseRecorder

	// Test1
	r = httptest.NewRecorder()
	ServiceUnavailable(r, errors.New("HTTP Service unavailable"))

	if r.Result().StatusCode != http.StatusServiceUnavailable {
		t.Error("ERROR Wrong status code: ", r.Result().StatusCode)
	}

	// Test 2
	r = httptest.NewRecorder()
	InternalServer(r, errors.New("HTTP Internal Server"))

	if r.Result().StatusCode != http.StatusInternalServerError {
		t.Error("ERROR Wrong status code: ", r.Result().StatusCode)
	}
	// Test 3
	r = httptest.NewRecorder()
	NotFound(r, errors.New("HTTP Not found"))

	if r.Result().StatusCode != http.StatusNotFound {
		t.Error("ERROR Wrong status code: ", r.Result().StatusCode)
	}
	// Test 4
	r = httptest.NewRecorder()
	BadRequest(r, errors.New("HTTP Bad Request"))

	if r.Result().StatusCode != http.StatusBadRequest {
		t.Error("ERROR Wrong status code: ", r.Result().StatusCode)
	}

}
