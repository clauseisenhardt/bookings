package handlers

import (
	"fmt"
	"net/http"

	"github.com/clauseisenhardt/bookings/pkg/config"
	"github.com/clauseisenhardt/bookings/pkg/models"
	"github.com/clauseisenhardt/bookings/pkg/render"
)

// Repo the repositiry used by handles
var Repo *Repository

// Repository is the repositiory type
type Repository struct {
	App *config.AppConfig
}

// NewRepo cretes a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Send the data to template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (m *Repository) Calc(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "This is the Calc page\n")
	if err != nil {
		fmt.Println(err)
	}
	sum := addValues(2, 2)
	_, _ = fmt.Fprintf(w, fmt.Sprintf("2 + 2 is : %d", sum))
}
func addValues(x, y int) int {
	return x + y
}
