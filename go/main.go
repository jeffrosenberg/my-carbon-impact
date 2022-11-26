package main

import (
	"net/http"

	"github.com/jeffrosenberg/my-carbon-impact/internal/html"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
)

func init() {
	logger = logging.GetLogger().With().Logger()
}

func main() {
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/profile", profile)
	logger.Info().Msg("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	p := html.DashboardParams{
		Title:   "Dashboard",
		Message: "Hello from dashboard",
	}
	html.Dashboard(w, p)
}

func profile(w http.ResponseWriter, r *http.Request) {
	p := html.ProfileParams{
		Title:   "Profile Show",
		Message: "Hello from profile show",
	}
	html.CreateProfile(w, p)
}
