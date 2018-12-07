package main

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/slack"
)

const htmlIndex = `<html><body>
<a href="/oauth2Login">Log in with My oauth2 server</a>
</body></html>
`

var (
	oauth2Config = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/oauth2Callback",
		ClientID:     "495996213110.494794973285",
		ClientSecret: "8bd2a8764b6772546d8ed31008b6b503",
		Scopes:       []string{"identity.basic", "identity.email"},
		Endpoint:     slack.Endpoint,
	}
	// Some random string, random for each request
	oauthStateString = "random"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/oauth2Login", handleOauth2Login)
	http.HandleFunc("/oauth2Callback", handleOauth2Callback)
	fmt.Println("app run http://localhost:3000")
	fmt.Println(http.ListenAndServe(":3000", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)
}

func handleOauth2Login(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleOauth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "token: %s\n", token.AccessToken)
}
