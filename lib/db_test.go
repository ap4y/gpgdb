package lib

import (
	"os"
	"testing"
)

func TestKeys(t *testing.T) {
	db := createDB(t)
	defer cleanDB(db)

	db.Put("ap4y", "foo", []byte("bar"))
	keys, err := db.Keys("ap4y")
	if err != nil {
		t.Fatalf("Unable to set value %s", err)
	}

	if len(keys) != 1 {
		t.Fatalf("Invalid keys returned: %v", keys)
	}

	if keys[0] != "foo" {
		t.Errorf("Invalid keys returned: %v", keys)
	}
}

func TestPut(t *testing.T) {
	db := createDB(t)
	defer cleanDB(db)

	err := db.Put("ap4y", "foo", []byte("bar"))
	if err != nil {
		t.Fatalf("Unable to set value %s", err)
	}
}

func TestGet(t *testing.T) {
	db := createDB(t)
	defer cleanDB(db)

	db.Put("ap4y", "foo", []byte("bar"))
	val, err := db.Get("ap4y", "foo")
	if err != nil {
		t.Fatalf("Unable to get value %s", err)
	}

	if val != "bar" {
		t.Errorf("Invalid value in Read, expected bar, got %s", val)
	}
}

func TestDelete(t *testing.T) {
	db := createDB(t)
	defer cleanDB(db)

	db.Put("ap4y", "foo", []byte("bar"))
	err := db.Delete("ap4y", "foo")
	if err != nil {
		t.Fatalf("Unable to delete value %s", err)
	}
}

func createDB(t *testing.T) *DB {
	db, err := NewDB("./test.db")
	if err != nil {
		t.Fatalf("Unable to create DB %s", err)
	}

	return db
}

func cleanDB(db *DB) {
	db.Close()
	os.Remove("./test.db")
}
