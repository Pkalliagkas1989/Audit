package handlers

// import (
// 	"net/http" // Package for handling HTTP requests and responses
// )

// // RenderErrorTemplate renders the "errors.html" template with a given status code and error message
// func RenderErrorTemplate(w http.ResponseWriter, statusCode int, message string) {
// 	// Set the status code of the response
// 	w.WriteHeader(statusCode)

// 	// Define the data structure to pass to the template
// 	data := struct {
// 		StatusCode int
// 		Message    string
// 	}{
// 		StatusCode: statusCode,
// 		Message:    message,
// 	}

// 	// Attempt to render the "errors.html" template with the data provided
// 	if err := Templates.ExecuteTemplate(w, "errors.html", data); err != nil {
// 		// If template execution fails, set status code and write error message
// 		w.WriteHeader(statusCode)
// 		w.Write([]byte("Error rendering error page"))
// 		return
// 	}
// }
