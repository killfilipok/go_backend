package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/killfilipok/backend_stuff/03_project/database"
	"github.com/killfilipok/backend_stuff/03_project/structs"
	"golang.org/x/oauth2"
)

func GetRepos(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(string)

	user := &structs.User{}

	user.Uid = userId

	result := database.DBCon.QueryRow("SELECT * from users WHERE uid=$1", userId)

	err := result.Scan(&user.Username, &user.Password, &user.Uid, &user.ImageUrl, &user.GithubToken)

	checkErr := func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	fmt.Println(user)

	checkErr()

	token, err := TokenFromJSON(user.GithubToken)

	checkErr()

	repoCreds := &structs.RepoCred{}
	err = json.NewDecoder(r.Body).Decode(repoCreds)

	checkErr()

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)

	fmt.Println(token)
	fmt.Println(repoCreds.Owner)
	fmt.Println(repoCreds.Repo)

	commitInfo, _, err := client.Repositories.ListCommits(oauth2.NoContext, repoCreds.Owner, repoCreds.Repo, nil)

	if err != nil {
		fmt.Printf("Problem in commit information %v\n", err)
		checkErr()
	}

	urlsSlice := []string{}
	fmt.Printf("%+v\n", commitInfo) // Last commit information
	for _, val := range commitInfo {
		urlsSlice = append(urlsSlice, val.GetURL())
	}
	arrayInJSON, err := json.Marshal(urlsSlice)
	checkErr()
	w.Write(arrayInJSON)
}
