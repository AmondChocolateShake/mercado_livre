package order

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
	r.POST("/orders", h.create)
	r.GET("/orders/:id", h.get)
}

type createRequest struct {
	ItemID   string `json:"item_id"`
	BuyerID  int64  `json:"buyer_id"`
	Quantity int    `json:"quantity"`
}

// TODO
// 1. JSON 바인딩 실패 → 400
// 2. h.svc.Create 호출
// 3. 에러 종류별 상태 코드
//    - 없는 주문/상품: 404
//    - 판매 종료, 재고 부족처럼 "지금 상태와 충돌": 409
//    - 잘못된 입력: 400
//    - 그 외: 500
// 4. 성공이면 201 + 주문 JSON
func (h *Handler) create(c *gin.Context) {
	apierror.Respond(c, http.StatusNotImplemented, "not_implemented", errNotImplemented)
}

// TODO
// c.Param("id")로 h.svc.Get 호출 → 200 또는 에러 응답
func (h *Handler) get(c *gin.Context) {
	apierror.Respond(c, http.StatusNotImplemented, "not_implemented", errNotImplemented)
}
