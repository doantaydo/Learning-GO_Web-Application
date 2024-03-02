package models

import "github.com/doantaydo/Learning-GO_Web-Application/Hotel-Bookings/internal/forms"

// Hold all data types
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
