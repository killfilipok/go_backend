package mySqlFuncs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/killfilipok/backend_stuff/03_project/structs"
)

func SendUserObjBack(user structs.User, w http.ResponseWriter, db *sql.DB) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// w.WriteHeader(http.StatusOK)

	w.Write(userJSON)
}

func RowExists(query string, db *sql.DB, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		// glog.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}
