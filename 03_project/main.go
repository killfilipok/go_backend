package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/killfilipok/backend_stuff/03_project/googleAuth"
	"github.com/killfilipok/backend_stuff/03_project/mySqlFuncs"
	"github.com/killfilipok/backend_stuff/03_project/structs"

	_ "github.com/lib/pq"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:whiteCup3721@127.0.0.1/backend_db?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	googleAuth.Init(db)
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/", homePage)
	http.HandleFunc("/signin/google", googleAuth.LoginHandler)
	http.HandleFunc("/signin/google/callback", googleAuth.CallbackHandler)

	http.ListenAndServe(":3000", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// db.Query("DELETE FROM users") //WARNING WARNING WARNING
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	users := make([]structs.User, 0)
	for rows.Next() {
		user := structs.User{}
		// uid:= nil
		err := rows.Scan(&user.Username, &user.Password, &user.Uid)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		} else {
			users = append(users, user)
		}

	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "home.gohtml", users)
}

func signup(w http.ResponseWriter, r *http.Request) {
	creds := &structs.User{}
	err := json.NewDecoder(r.Body).Decode(creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mySqlFuncs.RowExists("select uid from users where username=$1", db, creds.Username) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	id := xid.New()

	if _, err = db.Query("insert into users values ($1, $2, $3)", creds.Username, string(hashedPassword), id.String()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mySqlFuncs.SendUserObjBack(structs.User{creds.Username, string(hashedPassword), id.String()}, w, db)
}

func signin(w http.ResponseWriter, r *http.Request) {
	creds := &structs.User{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedCreds := &structs.User{}

	result := db.QueryRow("select * from users where username=$1", creds.Username)

	err = result.Scan(&storedCreds.Username, &storedCreds.Password, &storedCreds.Uid)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mySqlFuncs.SendUserObjBack(*storedCreds, w, db)
}
