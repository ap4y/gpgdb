package lib

import "testing"

func TestNewEntityPublic(t *testing.T) {
	entity, err := NewEntity("../fixtures/pub/some_user.gpg", nil)
	if err != nil {
		t.Fatalf("Unable to read entity: %s", err)
	}

	if entity.User != "some_user" {
		t.Errorf("Invalid entity user: %s", entity.User)
	}

	if len(entity.EntityList) != 1 {
		t.Errorf("Invalid entityList: %#v", entity.EntityList)
	}
}

func TestNewEntitySecret(t *testing.T) {
	entity, err := NewEntity("../fixtures/sec/secring.gpg", func(entity *Entity) ([]byte, error) {
		return []byte("12345678"), nil
	})
	if err != nil {
		t.Fatalf("Unable to read entity: %s", err)
	}

	if len(entity.EntityList) != 1 {
		t.Fatalf("Invalid entityList: %#v", entity.EntityList)
	}

	if entity.EntityList[0].PrivateKey.Encrypted {
		t.Errorf("Returned Encrypred Private Key")
	}

	if entity.EntityList[0].PrimaryKey == nil {
		t.Errorf("Returned Encrypred Public Key")
	}

	if entity.EntityList[0].Subkeys[0].PrivateKey.Encrypted {
		t.Errorf("Returned Encrypred SubKey")
	}
}

func TestGetIdentity(t *testing.T) {
	entity, err := NewEntity("../fixtures/pub/some_user.gpg", nil)
	if err != nil {
		t.Fatalf("Unable to read entity: %s", err)
	}

	if entity.GetIdentity() != "Some User <foo@bar.com>" {
		t.Errorf("Invalid entity user: %s", entity.GetIdentity())
	}
}

func TestPGPEntity(t *testing.T) {
	entity, err := NewEntity("../fixtures/pub/some_user.gpg", nil)
	if err != nil {
		t.Fatalf("Unable to read entity: %s", err)
	}

	if entity.PGPEntity() == nil {
		t.Errorf("PGPEntity is nil")
	}
}
