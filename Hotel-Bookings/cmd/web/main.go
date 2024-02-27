package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/config"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/handlers"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/render"
)

var portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false
	_ = SetUpAppConfig()

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// http.HandleFunc("/", repo.Home)
	// http.HandleFunc("/about", repo.About)
	//http.HandleFunc("/favicon.ico", repo.EmptyFunc)

	fmt.Println("Run Web Application at localhost" + portNumber)
	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}

func SetUpAppConfig() *handlers.Repository {
	// Add template to AppConfig then send to render
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err)
	}
	app.TemplateCache = templateCache
	// if you want to change tmpl file and check easier, set app.UserCache = false
	app.UseCache = app.InProduction
	render.NewTemplates(&app)
	// Send AppConfig to handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	return repo
}
