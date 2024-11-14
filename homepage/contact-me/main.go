package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

// ContactForm holds the form data
type ContactForm struct {
	Name    string
	Email   string
	Message string
}

// Enable CORS for handling cross-origin requests
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Render the HTML template
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parsing the template, assuming index.html exists
	tmplParsed, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	err = tmplParsed.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		log.Println(err)
	}
}

// Handle form submissions
func handleContact(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Allow cross-origin requests

	log.Println("Received request for /contact")

	if r.Method == "POST" {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Populate the form data
		form := ContactForm{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Message: r.FormValue("message"),
		}

		// Compose email
		subject := "New Contact Message"
		body := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", form.Name, form.Email, form.Message)

		// SMTP server configuration
		from := "your-email@gmail.com"
		password := os.Getenv("EMAIL_PASSWORD")
		if password == "" {
			log.Fatal("EMAIL_PASSWORD environment variable not set")
		}
		to := []string{"anxielworld@gmail.com"}

		msg := []byte("To: anxielworld@gmail.com\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n")

		// Sending the email
		err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, to, msg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
			log.Printf("Failed to send email: %v", err)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent successfully"))
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Serve the contact page with form when visiting the root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index.html", nil)
	})

	// Handle form submissions at /contact
	http.HandleFunc("/contact", handleContact)

	// Start the server
	log.Println("Server started at http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
