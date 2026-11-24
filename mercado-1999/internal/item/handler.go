package item

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
}

type createRequest struct {
	SellerID          int64  `json:"seller_id"`
	Title             string `json:"title"`
	Price             int64  `json:"price"`
	AvailableQuantity int    `json:"available_quantity"`
	Condition         string `json:"condition"`
	BuyingMode        string `json:"buying_mode"`
}

// TODO
// 1. c.ShouldBindJSON(&req) 실패 → 400
// 2. createRequest → CreateInput으로 옮겨 h.svc.Create 호출
// 3. 에러면 에러 종류에 맞는 상태 코드로 응답 (apierror.Respond 사용)
//    - 없는 상품: 404 / 잘못된 입력: 400 / 그 외: 500 (내부 에러 메시지는 숨기고 로그로)
// 4. 성공이면 201 + 상품 JSON
func (h *Handler) create(c *gin.Context) {
	apierror.Respond(c, http.StatusNotImplemented, "not_implemented", errNotImplemented)
}

// TODO
// c.Param("id")로 받은 ID로 h.svc.Get 호출 → 200 또는 에러 응답
func (h *Handler) get(c *gin.Context) {
	apierror.Respond(c, http.StatusNotImplemented, "not_implemented", errNotImplemented)
}
