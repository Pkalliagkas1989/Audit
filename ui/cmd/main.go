package main

import (
	"log"
	"net/http"
)

func checkSession(r *http.Request) bool {
	// Get session_id cookie from the request
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	// Create a request to backend to verify session
	req, err := http.NewRequest("GET", "http://localhost:8080/forum/api/session/verify", nil)
	if err != nil {
		return false
	}
	req.AddCookie(cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}


func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/index":
		http.ServeFile(w, r, "./static/templates/index.html")
	case "/":
		http.ServeFile(w, r, "./static/templates/index.html")
	case "/login":
	if checkSession(r) {
		http.Redirect(w, r, "/user", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./static/templates/login.html")
	case "/register":
	if checkSession(r) {
		// If session exists, redirect to user feed
		http.Redirect(w, r, "/user", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./static/templates/register.html")
	case "/guest":
		http.ServeFile(w, r, "./static/templates/guest/guest_mainpage.html")
	case "/guest/feed":
		http.ServeFile(w, r, "./static/templates/guest/guest_feed.html")
	case "/guest/category":
		http.ServeFile(w, r, "./static/templates/guest/guest_category.html")
	case "/guest/post":
		http.ServeFile(w, r, "./static/templates/guest/guest_post.html")
	case "/user":
	if !checkSession(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./static/templates/user.html")
	case "/guest/error":
		http.ServeFile(w, r, "./static/templates/error.html")

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
