package structs

import jwt "github.com/dgrijalva/jwt-go"

type Token struct {
	UserId string
	jwt.StandardClaims
}

type MyImage struct {
	Owner string `json:"owner", db:"owner"`
	Image string `json:"image", db:"image"`
}

type RepoCred struct {
	Owner string `json:"owner", db:"-"`
	Repo  string `json:"repo", db:"-"`
}

type Note struct {
	Owner     string `json:"owner", db:"owner"`
	Title     string `json:"title", db:"title"`
	Text      string `json:"text", db:"text"`
	CreatedAt int32  `json:"createdAt", db:"createdAt"`
	Uid       string `json:"uid", db:"uid"`
}

type User struct {
	Username    string `json:"username", db:"username"`
	Password    string `json:"password", db:"password"`
	Uid         string `json:"uid", db:"uid"`
	Token       string `json:"token", db:"-"`
	ImageUrl    string `json:"imageurl", db:"imageurl"`
	GithubToken string `json:"githubtoken", db:"githubtoken"`
}

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
	Psecret string `json:"psecret`
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}
