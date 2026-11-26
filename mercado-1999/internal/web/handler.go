package web

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/item"
)

//go:embed templates/*.html static/*
var files embed.FS

const listLimit = 50

type Handler struct {
	items *item.Service
	pages map[string]*template.Template
}

func NewHandler(items *item.Service) (*Handler, error) {
	pages := map[string]*template.Template{}
	for _, name := range []string{"index", "item", "not_found"} {
		t, err := template.New("layout.html").Funcs(funcs).ParseFS(files,
			"templates/layout.html",
			"templates/"+name+".html",
		)
		if err != nil {
			return nil, err
		}
		pages[name] = t
	}
	return &Handler{items: items, pages: pages}, nil
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	static, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	r.StaticFS("/static", http.FS(static))
	r.GET("/", h.index)
	r.GET("/articulo/:id", h.item)
}

type indexPage struct {
	Query string
	Items []item.Item
}

type itemPage struct {
	Query string
	Item  *item.Item
}

func (h *Handler) index(c *gin.Context) {
	query := c.Query("q")
	items, err := h.items.Search(c.Request.Context(), query, listLimit)
	if err != nil {
		h.serverError(c, err)
		return
	}
	h.render(c, http.StatusOK, "index", indexPage{Query: query, Items: items})
}

func (h *Handler) item(c *gin.Context) {
	it, err := h.items.Get(c.Request.Context(), c.Param("id"))
	if errors.Is(err, item.ErrNotFound) || errors.Is(err, item.ErrInvalidID) {
		h.render(c, http.StatusNotFound, "not_found", indexPage{})
		return
	}
	if err != nil {
		h.serverError(c, err)
		return
	}
	h.render(c, http.StatusOK, "item", itemPage{Item: it})
}

func (h *Handler) render(c *gin.Context, status int, page string, data any) {
	var buf bytes.Buffer
	if err := h.pages[page].ExecuteTemplate(&buf, "layout.html", data); err != nil {
		h.serverError(c, err)
		return
	}
	c.Data(status, "text/html; charset=utf-8", buf.Bytes())
}

func (h *Handler) serverError(c *gin.Context, err error) {
	log.Printf("web: %v", err)
	c.String(http.StatusInternalServerError, "Error interno del servidor")
}
