package main

import (
	"chatServer/dbaccess"
	"chatServer/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

const (
	signupResponseCodeUnknownError       = -100
	signupResponseCodeSuccess            = 0
	signupResponseCodeBadRequest         = -1
	signupResponseCodeInvalidUsername    = -2
	signupResponseCodeUsernameExists     = -3
	signupResponseCodeInvalidPassword    = -4
	signupResponseCodeInvalidEmail       = -5
	signupResponseCodeInvalidDisplayName = -6
)

type signupResponse struct {
	Result    bool `json:"result"`
	ErrorCode int  `json:"errorCode"`
}

// RunWebServer opens a HTTP server port for serving HTTPS/API requests.
func RunWebServer() {
	log.Println("Starting web server on :80 and :443...")

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/dosignup", signupRequest)
	if err := http.ListenAndServeTLS(":443", "domain.crt", "domain.key", nil); err != nil {
		log.Fatalf("ListenAndServeTLS error: %v", err)
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func hashRequest(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to serve hash password request:", err)
		}
	}()

	password := strings.TrimSpace(req.URL.Query()["password"][0])
	hashed := utils.HashPassword(password)
	fmt.Fprintf(w, "%v", hashed)
}

func signupRequest(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Failed to process signup request:", err)
			log.Println(string(debug.Stack()))
		}
	}()

	log.Printf("Processing sign up request...")

	req.ParseForm()
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	email := req.Form.Get("email")
	displayName := req.Form.Get("displayname")

	if username == "" || password == "" || email == "" || displayName == "" {
		signupFailed(w, signupResponseCodeBadRequest)
		return
	}

	// Validate username.
	var rxUsername = regexp.MustCompile("^[A-Za-z0-9_]{5,15}$")
	if !rxUsername.MatchString(username) {
		signupFailed(w, signupResponseCodeInvalidUsername)
		return
	}

	u, err := dbaccess.DBAccess.GetUserByUsername(username)
	if err != nil {
		signupFailed(w, signupResponseCodeUnknownError)
		return
	}

	if u != nil {
		signupFailed(w, signupResponseCodeUsernameExists)
		return
	}

	// Validate password.
	var rxPassword = regexp.MustCompile("^[a-zA-Z[:graph:]0-9£]{8,20}$")
	if len(password) < 8 || len(password) > 20 || !rxPassword.MatchString(password) {
		signupFailed(w, signupResponseCodeInvalidPassword)
		return
	}

	// Validate email.
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 1 || len(email) > 320 || !rxEmail.MatchString(email) {
		signupFailed(w, signupResponseCodeInvalidEmail)
		return
	}

	// Validate display name.
	if len(displayName) < 1 || len(displayName) > 20 {
		signupFailed(w, signupResponseCodeBadRequest)
		return
	}

	// Now we can save, but the username could still clash after this point.
	err = dbaccess.DBAccess.CreateAccount(username, utils.HashPassword(password), email, displayName, utils.GenerateGUID())
	if err != nil {
		signupFailed(w, signupResponseCodeUnknownError)
		return
	}

	signupSuccess(w)
}

func signupFailed(w http.ResponseWriter, errorCode int) {
	response := signupResponse{Result: false, ErrorCode: errorCode}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func signupSuccess(w http.ResponseWriter) {
	response := signupResponse{Result: true, ErrorCode: 0}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
