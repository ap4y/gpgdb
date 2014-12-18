package lib

import "path/filepath"

type EntityStorage struct {
	entities []*Entity
}

func NewEntityStorage(path string) (*EntityStorage, error) {
	files, err := filepath.Glob(path + "/*.gpg")
	if err != nil {
		return nil, err
	}

	var entities []*Entity
	for _, file := range files {
		entity, err := NewEntity(file, nil)
		if err != nil {
			continue
		}

		entities = append(entities, entity)
	}

	return &EntityStorage{entities}, nil
}

func (es *EntityStorage) EntityForIdentity(identity string) *Entity {
	for _, entity := range es.entities {
		if entity.GetIdentity() == identity {
			return entity
		}
	}

	return nil
}
