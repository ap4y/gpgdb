package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ap4y/gpgdb/lib"
)

func TestPutKey(t *testing.T) {
	r, _ := http.NewRequest("PUT", "/keys/foo", nil)
	request := &lib.Request{"foo", "", nil, []byte("bar"), r}
	response := httptest.NewRecorder()

	PutKey(response, request, &lib.MockDB{})

	if response.Code != http.StatusOK {
		t.Fatal("Invalid response code")
	}
}

func TestListKeys(t *testing.T) {
	r, _ := http.NewRequest("GET", "/keys", nil)
	request := &lib.Request{"foo", "", nil, nil, r}
	response := httptest.NewRecorder()

	ListKeys(response, request, &lib.MockDB{})
	if response.Code != http.StatusOK {
		t.Fatal("Invalid response code")
	}

	if body := response.Body.String(); body != "{\"keys\":[\"foo\",\"bar\"]}" {
		t.Errorf("Invalid response: %s", body)
	}
}

func TestGetKey(t *testing.T) {
	r, _ := http.NewRequest("GET", "/keys/foo", nil)
	request := &lib.Request{"foo", "", nil, nil, r}
	response := httptest.NewRecorder()

	GetKey(response, request, &lib.MockDB{})
	if response.Code != http.StatusOK {
		t.Fatal("Invalid response code")
	}

	if body := response.Body.String(); body != "{\"key\":\"\",\"value\":\"bar\"}" {
		t.Errorf("Invalid response: %s", body)
	}
}

func TestDeleteKey(t *testing.T) {
	r, _ := http.NewRequest("DELET", "/keys/foo", nil)
	request := &lib.Request{"foo", "", nil, nil, r}
	response := httptest.NewRecorder()

	DeleteKey(response, request, &lib.MockDB{})

	if response.Code != http.StatusOK {
		t.Fatal("Invalid response code")
	}
}
