package web

import (
	"html/template"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/keto-granola/keto-granola/internal/apperr"
	"github.com/keto-granola/keto-granola/internal/config"
	"github.com/keto-granola/keto-granola/internal/product"
	"github.com/keto-granola/keto-granola/internal/webassets"
)

type Handler struct {
	service      *product.Service
	templates    *template.Template
	assetsLoader *webassets.Loader
	clientURL    string
	devEnv       bool
}

type ProductData struct {
	Product   *product.GetProductResponse
	ClientURL string
	DevEnv    bool
	AssetSrc  string
	AssetCSS  []string
}

func NewHandler(svc *product.Service,
	assetsLoader *webassets.Loader,
	tmpl *template.Template,
	clientURL string,
	env config.Environment,
) *Handler {
	return &Handler{
		service:      svc,
		templates:    tmpl,
		assetsLoader: assetsLoader,
		clientURL:    clientURL,
		devEnv:       env == config.EnvironmentDevelopment,
	}
}

func (h *Handler) GetProductPage(e echo.Context) error {
	ID, err := uuid.Parse(e.Param("id"))
	if err != nil {
		return apperr.ToHTTPError(apperr.Validation("request.Validate", "VALIDATION_ERROR", apperr.ErrMsgValidation))
	}

	prod, err := h.service.GetProduct(e.Request().Context(), ID)
	if err != nil {
		return apperr.ToHTTPError(err)
	}

	islandEntryPath := webassets.IslandEntryPath("add-to-cart")

	var assetSrc string
	var assetCSS []string

	if h.devEnv {
		assetSrc = h.clientURL + "/" + islandEntryPath
	} else {
		if assetSrc, err = h.assetsLoader.Asset(islandEntryPath); err != nil {
			return apperr.ToHTTPError(err)
		}

		if assetCSS, err = h.assetsLoader.AssetCSS(islandEntryPath); err != nil {
			return apperr.ToHTTPError(err)
		}
	}

	productData := &ProductData{
		Product:   prod,
		ClientURL: h.clientURL,
		DevEnv:    h.devEnv,
		AssetSrc:  assetSrc,
		AssetCSS:  assetCSS,
	}

	return h.templates.ExecuteTemplate(e.Response(), "product.html", productData)
}
