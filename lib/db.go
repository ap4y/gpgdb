package lib

import (
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const keySeparator = "~"

type DBService interface {
	Keys(user string) ([]string, error)
	Put(user, key string, value []byte) error
	Get(user, key string) (string, error)
	Delete(user, key string) error
	Close() error
}

type DB struct {
	levelDB *leveldb.DB
}

func NewDB(path string) (*DB, error) {
	levelDB, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &DB{levelDB}, nil
}

func (db *DB) Keys(user string) ([]string, error) {
	prefix := user + keySeparator
	iter := db.levelDB.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	var keys []string
	for iter.Next() {
		key := strings.Replace(string(iter.Key()), prefix, "", -1)
		keys = append(keys, key)
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (db *DB) Put(user, key string, value []byte) error {
	userKey := db.userKey(user, key)
	return db.levelDB.Put(userKey, value, nil)
}

func (db *DB) Get(user, key string) (string, error) {
	userKey := db.userKey(user, key)
	value, err := db.levelDB.Get(userKey, nil)
	if err != nil {
		return "", err
	}

	return string(value), err
}

func (db *DB) Delete(user, key string) error {
	userKey := db.userKey(user, key)
	return db.levelDB.Delete(userKey, nil)
}

func (db *DB) Close() error {
	return db.levelDB.Close()
}

func (db *DB) userKey(user, key string) []byte {
	return []byte(user + keySeparator + key)
}
