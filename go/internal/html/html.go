package html

import (
	"embed"
	"html/template"
	"io"

	"github.com/jeffrosenberg/my-carbon-impact/internal/profile"
	"github.com/jeffrosenberg/my-carbon-impact/pkg/logging"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog"
)

//go:embed *
var files embed.FS

var (
	dashboard     = parse("dashboard.html")
	getProfile    = parse("profile/get_profile.html")
	createProfile = parse("profile/create_profile.html")
	logger        zerolog.Logger
)

func init() {
	logger = logging.GetLogger().
		With().
		Str("action", "serve_html").
		Logger()
}

type DashboardParams struct {
	Title   string
	Message string
}

func Dashboard(w io.Writer, p DashboardParams) error {
	logger.Info().Str("page", "dashboard").Msg("Serving dashboard page")
	return dashboard.Execute(w, p)
}

type ProfileParams struct {
	Title   string
	Message string
	Profile profile.Profile
}

func GetProfile(w io.Writer, p ProfileParams) error {
	logger.Info().
		Str("entity", "profile").
		Str("operation", "get").
		Msg("Serving get profile page")
	return getProfile.Execute(w, p)
}

func CreateProfile(w io.Writer, p ProfileParams) error {
	logger.Info().
		Str("entity", "profile").
		Str("operation", "create").
		Msg("Serving create profile page")
	return createProfile.Execute(w, p)
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").
			Funcs(sprig.FuncMap()).
			ParseFS(files, "layout.html", file))
}
