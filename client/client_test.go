package main

import (
	"testing"

	"github.com/ap4y/gpgdb/lib"
)

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
	entity, err := lib.NewEntity("../fixtures/sec/secring.gpg", func(entity *lib.Entity) ([]byte, error) {
		return []byte("12345678"), nil
	})
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
