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
	Success bool
}

var personalInfo PersonalInfo

func homePage(page string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

		var tpl = template.Must(template.ParseFiles("./html/" + page + ".html"))
		err := tpl.Execute(w, personalInfo)

		if err != nil {
			http.Error(w, "Whoops! That page cannot be located!\nError code: 404", http.StatusNotFound)
		}
    }
}

func formHandler(w http.ResponseWriter, r *http.Request){
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
	sendAddress := os.Getenv("FROM_EMAIL")
	toAddress := os.Getenv("EMAIL")
	auth := smtp.PlainAuth("", sendAddress, os.Getenv("KEY"), "smtp.gmail.com")

	to := []string{toAddress}

	msg := "From: " + contactDetails.FirstName + " " + contactDetails.LastName + "\r\n" +
	"Subject: New message from RhysBratti.com!\r\n" +
	"Message contents: \r\n" +
	contactDetails.Message + "\r\n" +
	"Reply email: " + contactDetails.Email
	
	body := []byte(msg)
	err := smtp.SendMail("smtp.gmail.com:587", auth, sendAddress, to, body)

	return err
}


func main() {

	err := godotenv.Load("./env/.env")
	if err != nil {
		log.Println("Warning: Error loading .env file")
		log.Println(".env file should be located in ./env/.env")
		log.Println("If running locally, run the make-configs.sh script and then fill in key values")
		log.Println("If running in Docker, mount /env volume to container using -v /path/to/env:/app/env")
	}

	_, err = os.Stat("media")
    if os.IsNotExist(err) {
        log.Println("Warning: No /media folder found. This will cause images to not load")
		log.Println("If running locally, run the make-configs.sh script to generate /media folder, then place photos in there")
		log.Println("If running in Docker, mount /media volume to container using -v /path/to/media:/app/media")
    }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	personalInfo = PersonalInfo{os.Getenv("EMAIL"), false}

	fs := justFilesFilesystem{http.Dir("html")}
	assets := justFilesFilesystem{http.Dir("assets")}
	media := justFilesFilesystem{http.Dir("media")}

	assetsFileServer := http.FileServer(assets)
	mediaFileServer := http.FileServer(media)
	fileServer := http.FileServer(fs)

	mux := http.NewServeMux()
	mux.Handle("/html/", http.StripPrefix("/html/", fileServer))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsFileServer))
	mux.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))
	mux.HandleFunc("/", homePage("home"))
	mux.HandleFunc("/resume", homePage("resume"))
	mux.HandleFunc("/home", homePage("home"))
	mux.HandleFunc("/about-me", homePage("about-me"))
	mux.HandleFunc("/contact", formHandler)
	mux.HandleFunc("/about-site", homePage("about-site"))
	http.ListenAndServe(":"+port, mux)
}

