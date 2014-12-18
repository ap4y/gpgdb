package lib

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
)

type Entity struct {
	User       string
	EntityList openpgp.EntityList
}

type PromptFunction func(entity *Entity) ([]byte, error)

func NewEntity(path string, callback PromptFunction) (*Entity, error) {
	entity, err := entityForPath(path)
	if err != nil {
		return nil, err
	}

	if err := entity.decrypt(callback); err != nil {
		return nil, err
	}

	return entity, nil
}

func (e *Entity) GetIdentity() string {
	pgpEntity := e.PGPEntity()
	if pgpEntity == nil {
		return ""
	}

	for key := range pgpEntity.Identities {
		return key
	}

	return ""
}

func (e *Entity) PGPEntity() *openpgp.Entity {
	if len(e.EntityList) == 0 {
		return nil
	}

	return e.EntityList[0]
}

func entityForPath(path string) (*Entity, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, errors.New("Unable to load identity for non-file path")
	}

	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	entityList, err := openpgp.ReadKeyRing(reader)
	if err != nil {
		return nil, err
	}

	return &Entity{strings.Replace(fi.Name(), ".gpg", "", -1), entityList}, nil
}

func (e *Entity) decrypt(callback PromptFunction) error {
	if callback == nil {
		return nil
	}

	for _, entity := range e.EntityList {
		if entity.PrivateKey == nil || !entity.PrivateKey.Encrypted {
			continue
		}

		passphrase, err := callback(e)
		if err != nil {
			return err
		}

		if err := entity.PrivateKey.Decrypt(passphrase); err != nil {
			return err
		}

		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
				if err := subkey.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
