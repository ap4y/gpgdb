package lib

type MockDB struct{}

func (db *MockDB) Keys(user string) ([]string, error) {
	return []string{"foo", "bar"}, nil
}

func (db *MockDB) Put(user, key string, value []byte) error {
	return nil
}

func (db *MockDB) Get(user, key string) (string, error) {
	return "bar", nil
}

func (db *MockDB) Delete(user, key string) error {
	return nil
}

func (db *MockDB) Close() error {
	return nil
}
