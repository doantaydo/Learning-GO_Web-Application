package render

import (
	"net/http"
	"testing"

	"github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	session.Put(r.Context(), "error", "456")
	session.Put(r.Context(), "warning", "789")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Flash value of 123 not found in session!")
	}
	if result.Error != "456" {
		t.Error("Flash value of 456 not found in session!")
	}
	if result.Warning != "789" {
		t.Error("Flash value of 789 not found in session!")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	templateCache, err := CreateTemplateCache()

	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templateCache

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})

	if err != nil {
		t.Error(err)
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})

	if err == nil {
		t.Error("Render template that does not exist")
	}
}

func TestNewTemplate(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
