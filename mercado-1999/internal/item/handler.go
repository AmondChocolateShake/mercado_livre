package item

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/AmondChocolateShake/mercado_livre/mercado-1999/internal/apierror"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/items", h.create)
	r.GET("/items/:id", h.get)
	r.PUT("/items/:id", h.update)
	r.DELETE("/items/:id", h.delete)
}

type createRequest struct {
	SellerID          int64           `json:"seller_id"`
	Title             string          `json:"title"`
	Price             decimal.Decimal `json:"price"`
	AvailableQuantity int             `json:"available_quantity"`
	Condition         string          `json:"condition"`
	BuyingMode        string          `json:"buying_mode"`
}

type updateRequest struct {
	Title             *string          `json:"title"`
	Price             *decimal.Decimal `json:"price"`
	AvailableQuantity *int             `json:"available_quantity"`
	Status            *string          `json:"status"`
}

func (h *Handler) create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, http.StatusBadRequest, "bad_request", err)
		return
	}

	it, err := h.svc.Create(c.Request.Context(), CreateInput{
		SellerID:          req.SellerID,
		Title:             req.Title,
		Price:             req.Price,
		AvailableQuantity: req.AvailableQuantity,
		Condition:         req.Condition,
		BuyingMode:        req.BuyingMode,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, it)
}

func (h *Handler) get(c *gin.Context) {
	it, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (h *Handler) update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, http.StatusBadRequest, "bad_request", err)
		return
	}

	it, err := h.svc.Update(c.Request.Context(), c.Param("id"), UpdateInput{
		Title:             req.Title,
		Price:             req.Price,
		AvailableQuantity: req.AvailableQuantity,
		Status:            req.Status,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, it)
}

func (h *Handler) delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		apierror.Respond(c, http.StatusNotFound, "not_found", err)
	case errors.Is(err, ErrItemClosed):
		apierror.Respond(c, http.StatusConflict, "conflict", err)
	case errors.Is(err, ErrInvalidID),
		errors.Is(err, ErrSellerNotFound),
		errors.Is(err, ErrInvalidTitle),
		errors.Is(err, ErrInvalidPrice),
		errors.Is(err, ErrInvalidQuantity),
		errors.Is(err, ErrInvalidCondition),
		errors.Is(err, ErrInvalidBuyingMode),
		errors.Is(err, ErrAuctionQuantity),
		errors.Is(err, ErrInvalidStatus),
		errors.Is(err, ErrNoChanges):
		apierror.Respond(c, http.StatusBadRequest, "bad_request", err)
	default:
		log.Printf("item: %v", err)
		apierror.Respond(c, http.StatusInternalServerError, "internal_error", errors.New("internal server error"))
	}
}
