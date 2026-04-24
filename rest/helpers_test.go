package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseID(t *testing.T) {
	t.Run("ValidID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/member/123", nil)
		req.SetPathValue("id", "123")
		id, err := parseID(req, "id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if id != 123 {
			t.Errorf("Expected 123, got %d", id)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/member/abc", nil)
		req.SetPathValue("id", "abc")
		_, err := parseID(req, "id")
		if err == nil {
			t.Error("Expected error for non-integer ID, got nil")
		}
	})

	t.Run("MissingKey", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/member/123", nil)
		_, err := parseID(req, "missing")
		if err == nil {
			t.Error("Expected error for missing key, got nil")
		}
	})
}

func TestParseUintID(t *testing.T) {
	t.Run("ValidUint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/member/456", nil)
		req.SetPathValue("id", "456")
		id, err := parseUintID(req, "id")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if id != 456 {
			t.Errorf("Expected 456, got %d", id)
		}
	})

	t.Run("NegativeUint", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/member/-1", nil)
		req.SetPathValue("id", "-1")
		_, err := parseUintID(req, "id")
		if err == nil {
			t.Error("Expected error for negative uint, got nil")
		}
	})
}

func TestParseIDFromString(t *testing.T) {
	t.Run("ValidString", func(t *testing.T) {
		id, err := parseIDFromString("789")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if id != 789 {
			t.Errorf("Expected 789, got %d", id)
		}
	})

	t.Run("InvalidString", func(t *testing.T) {
		_, err := parseIDFromString("xyz")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestJsonError(t *testing.T) {
	err := errors.New("test error")
	expected := `{"status":"error","error":"test error"}`
	result := jsonError(err)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestSendError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("test error")
	sendError(w, err, http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	expectedBody := `{"status":"error","error":"test error"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestSendJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	sendJSON(w, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := `{"foo":"bar"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestLogin_EmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/getJWT", nil)
	login(w, req)

	if w.Code != http.StatusNotAcceptable {
		t.Errorf("Expected status 406, got %d", w.Code)
	}
}
