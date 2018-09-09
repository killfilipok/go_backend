package mySqlFuncs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/killfilipok/backend_stuff/03_project/database"
)

func SendObjBack(obj interface{}, w http.ResponseWriter) {
	objJSON, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// w.WriteHeader(http.StatusOK)

	w.Write(objJSON)
}

func RowExists(query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := database.DBCon.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		// glog.Fatalf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}
