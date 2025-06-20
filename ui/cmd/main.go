package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var (
	APIBaseURL    = "http://localhost:8080/forum/api"
	AuthURI       = APIBaseURL + "/session/verify"
	DataURI       = APIBaseURL + "/allData"
	LoginURI      = APIBaseURL + "/session/login"
	LogoutURI     = APIBaseURL + "/session/logout"
	RegisterURI   = APIBaseURL + "/register"
	CategoriesURI = APIBaseURL + "/categories"
	ReactionsURI  = APIBaseURL + "/react"
	CommentsURI   = APIBaseURL + "/comments"
	CreatePostURI = APIBaseURL + "/posts/create"
	MyPostsURI    = APIBaseURL + "/user/posts"
	LikedPostsURI = APIBaseURL + "/user/liked"
)

func router(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/index":
        http.ServeFile(w, r, "./static/templates/index.html")
    case "/login":
        http.ServeFile(w, r, "./static/templates/login.html")
    case "/register":
        http.ServeFile(w, r, "./static/templates/register.html")
    case "/guest":
        http.ServeFile(w, r, "./static/templates/guest.html")
    case "/user":
        http.ServeFile(w, r, "./static/templates/user.html")
    case "/error":
        http.ServeFile(w, r, "./static/templates/error.html")
    case "/config":
        w.Header().Set("Content-Type", "application/json")
        config := map[string]string{
            "APIBaseURL":    APIBaseURL,
            "AuthURI":       AuthURI,
            "DataURI":       DataURI,
            "LoginURI":      LoginURI,
            "LogoutURI":     LogoutURI,
            "RegisterURI":   RegisterURI,
            "CategoriesURI": CategoriesURI,
            "ReactionsURI":  ReactionsURI,
            "CommentsURI":   CommentsURI,
            "CreatePostURI": CreatePostURI,
            "MyPostsURI":    MyPostsURI,
            "LikedPostsURI": LikedPostsURI,
        }
        json.NewEncoder(w).Encode(config)
    default:
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, "./static/templates/error.html")    
	}
}

func main() {
    // Static files
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Use the custom router for all other paths
    http.HandleFunc("/", router)

    log.Println("Serving on http://localhost:8081/")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}

