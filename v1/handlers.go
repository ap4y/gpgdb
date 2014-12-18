package v1

import (
	"net/http"

	"github.com/ap4y/gpgdb/lib"
)

func PutKey(w http.ResponseWriter, req *lib.Request, db lib.DBService) {
	key, value := req.URL.Query().Get(":key"), req.EncryptedBody
	if err := db.Put(req.UserName, key, value); err != nil {
		lib.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"key": key, "value": string(value)}
	lib.WriteJSON(w, response)
}

func ListKeys(w http.ResponseWriter, req *lib.Request, db lib.DBService) {
	keys, err := db.Keys(req.UserName)
	if err != nil {
		lib.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string][]string{"keys": keys}
	lib.WriteJSON(w, response)
}

func GetKey(w http.ResponseWriter, req *lib.Request, db lib.DBService) {
	key := req.URL.Query().Get(":key")
	value, err := db.Get(req.UserName, key)
	if err != nil {
		lib.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{"key": key, "value": string(value)}
	lib.WriteJSON(w, response)
}

func DeleteKey(w http.ResponseWriter, req *lib.Request, db lib.DBService) {
	key := req.URL.Query().Get(":key")
	if err := db.Delete(req.UserName, key); err != nil {
		lib.ErrorJSON(w, err.Error(), http.StatusBadRequest)
	}
}
