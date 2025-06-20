package handlers

// import (
// 	"net/http"
// 	"text/template"
// )

// var Templates = template.Must(template.ParseGlob("templates/*.html"))

// // HomeHandler handles requests to the home page ("/")
// func HomeHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if the requested URL path is "/" (the home page)
// 	if r.URL.Path != "/" {
// 		// If not, return a 404 Not Found response with an error message
// 		RenderErrorTemplate(w, http.StatusNotFound, "Page Not Found")
// 		return
// 	}

// 	// Render the "index.html" template with the prepared data
// 	err := Templates.ExecuteTemplate(w, "index.html", nil)
// 	// If there's an error rendering the template, return a 500 Internal Server Error response
// 	if err != nil {
// 		RenderErrorTemplate(w, http.StatusInternalServerError, "could not execute template")
// 		return
// 	}
// }
