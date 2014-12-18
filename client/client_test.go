package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ap4y/gpgdb/lib"
)

var entity, _ = lib.NewEntity("../fixtures/sec/secring.gpg", func(entity *lib.Entity) ([]byte, error) {
	return []byte("12345678"), nil
})

func TestEncrypt(t *testing.T) {
	entity, _ := lib.NewEntity("../fixtures/pub/some_user.gpg", nil)
	client := NewClient("http://example.com", entity)
	value, err := client.Encrypt("foo")
	if err != nil {
		t.Fatalf("Failed encrypting value: %s", err)
	}

	if value == nil {
		t.Errorf("Invalid encrypted value")
	}
}

func TestDecrypt(t *testing.T) {
	client := NewClient("http://example.com", entity)
	encrypted, _ := client.Encrypt("foo")
	value, err := client.Decrypt(encrypted)

	if err != nil {
		t.Fatalf("Failed decrypting value: %s", err)
	}

	if string(value) != "foo" {
		t.Errorf("Invalid decrypted value: %s", value)
	}
}

func TestPut(t *testing.T) {
	client := &Client{"http://example.com", entity, &MockClient{200, ""}}
	if err := client.Put("foo", "bar"); err != nil {
		t.Fatalf("Unable to put value: %s", err)
	}
}

func TestKeys(t *testing.T) {
	client := &Client{"http://example.com", entity, &MockClient{200, "{\"keys\":[\"foo\"]}"}}
	value, err := client.Keys()
	if err != nil {
		t.Fatalf("Unable to get keys: %s", err)
	}

	if len(value) != 1 {
		t.Fatal("Invalid Keys returned")
	}

	if value[0] != "foo" {
		t.Errorf("Invalid Keys value: %s", value)
	}
}

func TestGet(t *testing.T) {
	client := &Client{"http://example.com", entity, nil}
	enc, _ := client.Encrypt("bar")
	json := "{\"key\":\"foo\",\"value\":\"" + string(enc) + "\"}"
	client.HTTPClient = &MockClient{200, json}

	value, err := client.Get("foo")
	if err != nil {
		t.Fatalf("Unable to get value: %s", err)
	}

	if string(value) != "bar" {
		t.Errorf("Invalid Get value: %s", value)
	}
}

func TestDelete(t *testing.T) {
	client := &Client{"http://example.com", entity, &MockClient{200, ""}}
	if err := client.Delete("foo"); err != nil {
		t.Fatalf("Unable to delete value: %s", err)
	}
}

type MockClient struct {
	statusCode int
	body       string
}

func (c *MockClient) Do(request *http.Request) (*http.Response, error) {
	response := &http.Response{
		Status:     string(c.statusCode),
		StatusCode: c.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(c.body))}

	return response, nil
}
