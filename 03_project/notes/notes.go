package notes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/killfilipok/backend_stuff/03_project/database"
	"github.com/killfilipok/backend_stuff/03_project/mySqlFuncs"
	"github.com/killfilipok/backend_stuff/03_project/structs"
	"github.com/rs/xid"
)

func PostNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	userUid := r.Context().Value("user").(string)

	note := &structs.Note{}
	err := json.NewDecoder(r.Body).Decode(note)

	note.Owner = userUid
	note.Uid = xid.New().String()
	note.CreatedAt = int32(time.Now().Unix())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = database.DBCon.Query("insert into notes(owner,title,text,uid,created) values ($1, $2, $3, $4, $5)",
		note.Owner, note.Title, note.Text, note.Uid, note.CreatedAt); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mySqlFuncs.SendObjBack(note, w)
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	userUid := r.Context().Value("user").(string)

	rows, err := database.DBCon.Query("SELECT * FROM notes WHERE owner=$1", userUid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	notes := make([]structs.Note, 0)
	for rows.Next() {
		note := structs.Note{}

		err := rows.Scan(&note.Owner, &note.Title, &note.Text, &note.Uid, &note.CreatedAt)

		if err != nil {
			fmt.Println(err, 1)
			http.Error(w, http.StatusText(500), 500)
			return
		} else {
			notes = append(notes, note)
		}

	}
	if err = rows.Err(); err != nil {
		fmt.Println(err, 2)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	mySqlFuncs.SendObjBack(notes, w)
}

func parseNotes(jsonBuffer []byte) ([]structs.Note, error) {

	notes := []structs.Note{}

	err := json.Unmarshal(jsonBuffer, &notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func DeleteNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	userUid := r.Context().Value("user").(string)

	var list []structs.Note
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	list, err = parseNotes(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	for _, note := range list {
		_, err := database.DBCon.Query("DELETE FROM notes WHERE owner=$1 AND uid=$2", userUid, note.Uid)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}

func DeleteAllNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	userUid := r.Context().Value("user").(string)

	_, err := database.DBCon.Query("DELETE FROM notes WHERE owner=$1", userUid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	userUid := r.Context().Value("user").(string)

	note := &structs.Note{}
	err := json.NewDecoder(r.Body).Decode(note)

	note.Owner = userUid

	_, err = database.DBCon.Exec("UPDATE notes SET owner = $1, title=$2, text=$3, uid=$4, created=$5 WHERE owner=$1 AND uid=$4",
		userUid, note.Title, note.Text, note.Uid, note.CreatedAt)
	if err != nil || note.Uid == "" {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	mySqlFuncs.SendObjBack(note, w)
}
