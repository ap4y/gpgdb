package lib

import (
	"net/http"
	"testing"
)

var es, _ = NewEntityStorage("../fixtures/pub")
var entity = es.EntityForIdentity("Some User <foo@bar.com>")
var signEntity, _ = NewEntity("../fixtures/sec/secring.gpg", func(entity *Entity) ([]byte, error) {
	return []byte("12345678"), nil
})

func TestSignatureMessageGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo", nil)
	request := &Request{"foo", "bar", nil, nil, r}

	expected := "GET /foo X-SIGNER=bar"
	result := request.SignatureMessage()
	if result != expected {
		t.Errorf("Invalid SignatureMessage, expected: %s, got: %s", expected, result)
	}
}

func TestSignatureMessagePost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/foo", nil)
	request := &Request{"foo", "bar", nil, []byte("{\"foo\":\"bar\"}"), r}

	expected := "POST /foo X-SIGNER=bar {\"foo\":\"bar\"}"
	result := request.SignatureMessage()
	if result != expected {
		t.Errorf("Invalid SignatureMessage, expected: %s, got: %s", expected, result)
	}
}

func TestSignature(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo", nil)
	request := &Request{"some_user", "Some User <foo@bar.com>", nil, nil, r}

	result, err := request.Signature(signEntity)
	if err != nil {
		t.Fatalf("Enable to generate signature: %s", err)
	}

	if result == "" {
		t.Errorf("Signature is empty")
	}
}

func TestSign(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo", nil)
	request := &Request{"some_user", "Some User <foo@bar.com>", nil, nil, r}

	err := request.Sign(signEntity)
	if err != nil {
		t.Fatalf("Enable to generate signature: %s", err)
	}

	if request.Header.Get(pgpNameHeader) != "Some User <foo@bar.com>" {
		t.Errorf("Invalid Signer Header")
	}

	if request.Header.Get(signatureHeader) == "" {
		t.Errorf("Empty signature header")
	}
}

func TestCheckSignature(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo", nil)
	request := &Request{"some_user", "Some User <foo@bar.com>", entity.EntityList, nil, r}

	request.Sign(signEntity)

	err := request.CheckSignature()
	if err != nil {
		t.Errorf("Invalid Signature: %s", err)
	}
}

func TestAuthenticatedRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/foo", nil)
	request := &Request{"some_user", "Some User <foo@bar.com>", entity.EntityList, nil, r}

	signature, _ := request.Signature(signEntity)

	r.Header.Add(pgpNameHeader, "Some User <foo@bar.com>")
	r.Header.Add(signatureHeader, signature)

	_, err := AuthenticatedRequest(r, es)
	if err != nil {
		t.Fatalf("Enable to create request: %s", err)
	}
}
