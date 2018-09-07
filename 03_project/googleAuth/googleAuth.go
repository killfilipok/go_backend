package googleAuth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/killfilipok/backend_stuff/03_project/mySqlFuncs"
	"github.com/killfilipok/backend_stuff/03_project/structs"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var db *sql.DB
var cred structs.Credentials
var conf *oauth2.Config
var state string
var store = sessions.NewCookieStore([]byte("secret"))

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func Init(DB *sql.DB) {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	db = DB
	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://127.0.0.1:3000/signin/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, _ := conf.Exchange(oauth2.NoContext, code)
	fmt.Fprintf(w, token.AccessToken)

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	contents, _ := ioutil.ReadAll(response.Body)
	var user *structs.GoogleUser
	_ = json.Unmarshal(contents, &user)

	// fmt.Println(user)

	var hashedPassword []byte
	if mySqlFuncs.RowExists("select uid from users where uid=$1", db, user.ID) {
		result := db.QueryRow("select * from users where uid=$1", user.ID)

		userObj := structs.User{}

		err = result.Scan(&userObj.Username, &userObj.Password, &userObj.Uid)
		fmt.Println("SignIn with google")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("SignIn with google")
			return
		}

	} else {
		res, err := bcrypt.GenerateFromPassword([]byte(user.ID+cred.Psecret), 8)

		hashedPassword = res
		fmt.Println("SignUp with google")
		if _, err = db.Query("insert into users values ($1, $2, $3)", user.Email, string(hashedPassword), user.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("SignUp with google")
			return
		}

	}

	mySqlFuncs.SendUserObjBack(structs.User{user.Email, string(hashedPassword), user.ID}, w, db)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := uniuri.New()
	url := conf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
