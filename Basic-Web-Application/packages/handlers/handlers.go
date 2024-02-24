package handlers

import (
	"net/http"

	"github.com/doantaydo/Learning-GO_Web-Application/Basic-Web-Application/packages/config"
	"github.com/doantaydo/Learning-GO_Web-Application/Basic-Web-Application/packages/models"
	"github.com/doantaydo/Learning-GO_Web-Application/Basic-Web-Application/packages/render"
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
	remoteID := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteID)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

func (m *Repository) EmptyFunc(w http.ResponseWriter, r *http.Request) {
	return
}
