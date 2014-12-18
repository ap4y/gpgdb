package lib

import "testing"

func TestEntityForIdentity(t *testing.T) {
	es, err := NewEntityStorage("../fixtures/pub")

	if err != nil {
		t.Fatalf("Enable to create IdentityStorage: %s", err)
	}

	entity := es.EntityForIdentity("Some User <foo@bar.com>")
	if entity == nil {
		t.Fatal("Returned Entity is nil")
	}

	if entity.User != "some_user" {
		t.Errorf("Invalid user %s", entity.User)
	}
}
