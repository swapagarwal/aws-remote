package main

import (
	"fmt"
	auth "github.com/swapagarwal/aws-remote/Godeps/_workspace/src/github.com/abbot/go-http-auth"
	"net/http"
	"os"
)

func Secret(user, realm string) string {
	if user == os.Getenv("login") {
		return os.Getenv("password")
	}
	return ""
}

func handle(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "Welcome, %s!", r.Username)
}

func main() {
	authenticator := auth.NewBasicAuthenticator("aws-remote.herokuapp.com", Secret)
	http.HandleFunc("/", authenticator.Wrap(handle))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
