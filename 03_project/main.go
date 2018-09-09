package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/killfilipok/backend_stuff/03_project/database"
	"github.com/killfilipok/backend_stuff/03_project/github"
	googleAuth "github.com/killfilipok/backend_stuff/03_project/googleAuth"
	"github.com/killfilipok/backend_stuff/03_project/imageservice"
	"github.com/killfilipok/backend_stuff/03_project/mySqlFuncs"
	"github.com/killfilipok/backend_stuff/03_project/notes"

	"github.com/killfilipok/backend_stuff/03_project/JwtAuth"
	"github.com/killfilipok/backend_stuff/03_project/structs"

	_ "github.com/lib/pq"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

func init() {
	var err error
	database.DBCon, err = sql.Open("postgres", "postgres://postgres:whiteCup3721@127.0.0.1/backend_db?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = database.DBCon.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	fs := http.FileServer(http.Dir(imageservice.UploadPath))
	http.Handle("/images/", http.StripPrefix("/images", fs))
	http.HandleFunc("/upload", JwtAuth.JwtAuthentication(imageservice.UploadFileHandler()))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/notes/update", JwtAuth.JwtAuthentication(notes.UpdateNote))
	http.HandleFunc("/notes/delete", JwtAuth.JwtAuthentication(notes.DeleteNotes))
	http.HandleFunc("/notes/delete/all", JwtAuth.JwtAuthentication(notes.DeleteAllNotes))
	http.HandleFunc("/notes/get", JwtAuth.JwtAuthentication(notes.GetNotes))
	http.HandleFunc("/notes/create", JwtAuth.JwtAuthentication(notes.PostNote))
	http.HandleFunc("/delete/all", deleteAll)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signin/google", googleAuth.LoginHandler)
	http.HandleFunc("/signin/google/callback", googleAuth.CallbackHandler)
	http.HandleFunc("/signin/github", JwtAuth.JwtAuthentication(github.LoginHandler))
	http.HandleFunc("/signin/github/callback", github.CallbackHandler)
	http.HandleFunc("/github/repo", JwtAuth.JwtAuthentication(github.GetRepos))

	http.ListenAndServe(":3000", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	rows, err := database.DBCon.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	users := make([]structs.User, 0)
	for rows.Next() {
		user := structs.User{}

		err := rows.Scan(&user.Username, &user.Password, &user.Uid, &user.ImageUrl, &user.GithubToken)

		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		users = append(users, user)

	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "home.gohtml", users)
}

func signup(w http.ResponseWriter, r *http.Request) {
	creds := &structs.User{}
	err := json.NewDecoder(r.Body).Decode(creds)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if mySqlFuncs.RowExists("select uid from users where username=$1", creds.Username) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	id := xid.New()
	creds.Uid = id.String()

	if _, err = database.DBCon.Query("insert into users values ($1, $2, $3,$4,$5)", creds.Username, string(hashedPassword), creds.Uid, "", ""); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tk := &structs.Token{UserId: id.String()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	mySqlFuncs.SendObjBack(structs.User{creds.Username, string(hashedPassword), id.String(), tokenString, "", ""}, w)
}

func signin(w http.ResponseWriter, r *http.Request) {
	creds := &structs.User{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedCreds := &structs.User{}

	result := database.DBCon.QueryRow("select * from users where username=$1", creds.Username)

	err = result.Scan(&storedCreds.Username, &storedCreds.Password, &storedCreds.Uid, &storedCreds.ImageUrl, &storedCreds.GithubToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tk := &structs.Token{UserId: storedCreds.Uid}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	storedCreds.Token = tokenString

	mySqlFuncs.SendObjBack(*storedCreds, w)
}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	database.DBCon.Query("DELETE FROM users")
	database.DBCon.Query("DELETE FROM notes")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
