package handlers

import (
	"fmt"
	"net/http"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/packages/config"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/packages/models"
	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/packages/render"
)

type Repository struct {
	App *config.AppConfig
}

// the repository used by the  handler
var Repo *Repository

// creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello."

	// remoteID := r.RemoteAddr
	// m.App.Session.Put(r.Context(), "remote_ip", remoteID)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Println(r)
	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

func (m *Repository) EmptyFunc(w http.ResponseWriter, r *http.Request) {
	return
}
