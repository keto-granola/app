package admin

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"github.com/keto-granola/keto-granola/internal/apperr"
	"github.com/keto-granola/keto-granola/internal/config"
	"github.com/keto-granola/keto-granola/internal/webassets"
)

type Handler struct {
	templates    *template.Template
	assetsLoader *webassets.Loader
	clientURL    string
	devEnv       bool
}

type ShellData struct {
	ClientURL string
	DevEnv    bool
	AssetSrc  string
	AssetCSS  []string
}

func NewHandler(
	assetsLoader *webassets.Loader,
	tmpl *template.Template,
	clientURL string,
	env config.Environment,
) *Handler {
	return &Handler{
		templates:    tmpl,
		assetsLoader: assetsLoader,
		clientURL:    clientURL,
		devEnv:       env == config.EnvironmentDevelopment,
	}
}

func (h *Handler) ServeShell(e echo.Context) error {
	const adminEntryPath = "src/admin/entries/admin.tsx"
	asset, err := webassets.ResolveEntry(h.assetsLoader, adminEntryPath, h.clientURL, h.devEnv)
	if err != nil {
		return apperr.ToHTTPError(err)
	}

	shellData := &ShellData{
		ClientURL: h.clientURL,
		DevEnv:    h.devEnv,
		AssetSrc:  asset.Src,
		AssetCSS:  asset.CSS,
	}

	return h.templates.ExecuteTemplate(e.Response(), "admin.html", shellData)
}
