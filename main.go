package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id        int     `json:"id"`
	Ident     string  `json:"ident"`
	Email     *string `json:"email"`
	Title     *string `json:"title"`
	Pw        *string `json:"pw"`
	ApiKey    *string `json:"api_key"`
	IsActive  bool    `json:"true"`
	IsService bool    `json:"false"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal((err))
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/user/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/user", createUser(db)).Methods("POST")
	router.HandleFunc("/user/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/user/{id}", deleteUser(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":7777", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := "SELECT id, ident, email, title, pw, api_key, is_active, is_service FROM public.usr"
		log.Println(q)
		rows, err := db.Query(q)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer rows.Close()

		users := []User{}

		for rows.Next() {
			var u User
			if err := rows.Scan(&u.Id, &u.Ident, &u.Email, &u.Title, &u.Pw, &u.ApiKey, &u.IsActive, &u.IsService); err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT id, ident, email, title, pw, api_key, is_active, is_service FROM usr WHERE id = $1", id).Scan(&u.Id, &u.Ident, &u.Email, &u.Title, &u.Pw, &u.ApiKey, &u.IsActive, &u.IsService)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(u)
	}
}

func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		err := db.QueryRow("INSERT INTO usr(ident, title) VALUES($1, $1) RETURNING id", u.Ident, u.Title)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(u)
	}
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		vars := mux.Vars(r)
		id := vars["id"]
		_, err := db.Exec("UPDATE usr SET ident = $1, title = $2 WHERE id = $3", u.Ident, u.Title, id)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(u)
	}
}

func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT id, ident, email, title, pw, api_key, is_active, is_service FROM usr WHERE id = $1", id).Scan(&u.Id, &u.Ident, &u.Email, &u.Title, &u.Pw, &u.ApiKey, &u.IsActive, &u.IsService)
		if err != nil {
			log.Fatal(err)
		} else {
			_, err := db.Exec("DELETE usr WHERE id = $3", id)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode("User deleted")
		}
	}
}
