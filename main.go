package main

import (
	"Session/session"
	signinView "Session/view"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

var store = session.StoreSession{Sessions: make(map[string]*session.Session)}

type user struct {
	Username string
	Age      int
}

func main() {
	serverAddr := "localhost:8000"
	server := http.Server{
		Addr:    serverAddr,
		Handler: http.DefaultServeMux,
	}

	http.HandleFunc("/signin", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			err := signinView.Signin().Render(request.Context(), writer)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		case http.MethodPost:
			err := request.ParseForm()
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			username := request.Form.Get("username")
			password := request.Form.Get("password")
			if username == "user" && password == "123" {
				newSession := store.Create()
				newSession.Values["user"] = user{
					Username: username,
					Age:      rand.Intn(18),
				}
				http.SetCookie(writer, &http.Cookie{
					Name:  "Session",
					Value: newSession.ID,
				})

				fmt.Fprintf(writer, fmt.Sprintf("username : %s", username))
				return
			}
			fmt.Fprintf(writer, fmt.Sprintf("Incorrect username or password"))
			return
		}
	})

	http.HandleFunc("/logout", func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("Session")
		if err != nil {
			return
		}

		store.Delete(cookie.Value)

		http.SetCookie(writer, &http.Cookie{
			Name:   "Session",
			Value:  "",
			MaxAge: -1,
		})

		fmt.Fprint(writer, "Logged out successfully!")
	})

	http.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("Session")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		storeSession, err := store.Get(cookie.Value)
		if err != nil {
			if errors.Is(err, &session.SessionNotFound{}) {
				http.SetCookie(writer, &http.Cookie{
					Name:   "Session",
					Value:  "",
					MaxAge: -1,
				})
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if sessionUser, ok := storeSession.Values["user"].(user); ok {
			fmt.Fprintf(writer, "Username: %s, Age: %d", sessionUser.Username, sessionUser.Age)
		} else {
			http.Error(writer, "User not found in session", http.StatusInternalServerError)
			return
		}
	})

	fmt.Printf("Listening on http://%s\n", serverAddr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
		return
	}
}
