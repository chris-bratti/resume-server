package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"net/smtp"
)

type ContactDetails struct {
    Email  		string
    FirstName 	string
    LastName 	string
	Message		string
}

type PersonalInfo struct {
	Email	string
	Phone	string
	Success bool
}

var personalInfo PersonalInfo

func homePage(page string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        /*
		cookie, _ := r.Cookie("authelia_session")
		if cookie != nil {
			http.Error(w, "You are not authorized to view this site", http.StatusForbidden)
		}else {
			var tpl = template.Must(template.ParseFiles("./html/home.html"))
			tpl.Execute(w, nil)
		}
		*/

		var tpl = template.Must(template.ParseFiles("./html/" + page + ".html"))
		err := tpl.Execute(w, personalInfo)

		if err != nil {
			http.Error(w, "Whoops! That page cannot be located!\nError code: 404", http.StatusNotFound)
		}
    }
}

func formHandler(w http.ResponseWriter, r *http.Request){
	/*
		cookie, _ := r.Cookie("authelia_session")
		if cookie != nil {
			http.Error(w, "You are not authorized to view this site", http.StatusForbidden)
		}else {
			var tpl = template.Must(template.ParseFiles("./html/home.html"))
			tpl.Execute(w, nil)
		}
		*/

	tmpl := template.Must(template.ParseFiles("./html/contact.html"))
	
	if r.Method != http.MethodPost {
		tmpl.Execute(w, personalInfo)
		return
	}

	details := ContactDetails{
		Email:   	r.FormValue("email"),
		FirstName: 	r.FormValue("firstname"),
		LastName: 	r.FormValue("lastname"),
		Message: 	r.FormValue("message"),
	}

	err := sendEmail(details)

	if err != nil {
		log.Fatal(err)
	}

	personalInfo.Success = true

	tmpl.Execute(w, personalInfo)

	personalInfo.Success = false
	
}

type justFilesFilesystem struct {
	fs http.FileSystem
}
  
func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}
  
type neuteredReaddirFile struct {
	http.File
}
  
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func sendEmail(contactDetails ContactDetails) error{
	sendAddress := os.Getenv("EMAIL")
	auth := smtp.PlainAuth("", sendAddress, os.Getenv("KEY"), "smtp.gmail.com")

	to := []string{sendAddress}

	msg := "From: " + contactDetails.FirstName + " " + contactDetails.LastName + "\r\n" +
	"Subject: New message from Resume-Server!\r\n" +
	"Message contents: \r\n" +
	contactDetails.Message + "\r\n" +
	"Reply email: " + contactDetails.Email
	
	body := []byte(msg)
	err := smtp.SendMail("smtp.gmail.com:587", auth, sendAddress, to, body)

	return err
}


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	personalInfo = PersonalInfo{os.Getenv("EMAIL"), os.Getenv("PHONE"), false}

	fs := justFilesFilesystem{http.Dir("html")}

	assets := justFilesFilesystem{http.Dir("assets")}

	assetsFileServer := http.FileServer(assets)

	fileServer := http.FileServer(fs)

	mux := http.NewServeMux()
	mux.Handle("/html/", http.StripPrefix("/html/", fileServer))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsFileServer))
	mux.HandleFunc("/", homePage("home"))
	mux.HandleFunc("/resume", homePage("resume"))
	mux.HandleFunc("/home", homePage("home"))
	mux.HandleFunc("/about-me", homePage("about-me"))
	mux.HandleFunc("/contact", formHandler)
	mux.HandleFunc("/about-site", homePage("about-site"))
	http.ListenAndServe(":"+port, mux)
}

