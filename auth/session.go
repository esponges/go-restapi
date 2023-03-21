package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func Secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Not Auth!!", http.StatusForbidden)
		return
	}

	fmt.Println("No secrets for you")
}

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// do auth stuff

	session.Values["authenticated"] = true
	session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordMatch(password string) {
	hashed, err := HashPassword(password)
	if err != nil {
		fmt.Println("you did something wrong dude")
	}

	match := CheckPasswordHash(password, hashed)

	if !match {
		fmt.Println("Something wrong hashing")
	}
	fmt.Printf("Success hashing pw " + password + "to " + hashed)
}

func HashPasswordHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pw := params["pw"]
	fmt.Printf("handler")

	CheckPasswordMatch(pw)
}
