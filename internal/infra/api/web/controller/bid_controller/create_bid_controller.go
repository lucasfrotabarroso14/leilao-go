package bid_controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"leilao-go/configuration/rest_err"
	"leilao-go/internal/infra/api/web/validation"
	"leilao-go/internal/usecase/bid_usecase"
	"net/http"
)

type BidController struct {
	bidUseCase bid_usecase.NewBidUseCaseInterface
}

func NewBidController(bidUseCase bid_usecase.NewBidUseCaseInterface) *BidController {
	return &BidController{
		bidUseCase: bidUseCase,
	}
}

func (u *BidController) CreateBid(c *gin.Context) {
	var bidInputDTO bid_usecase.BidInputDTO
	if err := c.ShouldBindJSON(&bidInputDTO); err != nil {
		restErr := validation.ValidateErr(err)
		c.JSON(restErr.Code, restErr)
		return

	}
	err := u.bidUseCase.CreateBid(context.Background(), bidInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		c.JSON(restErr.Code, restErr)
		return
	}
	c.Status(http.StatusCreated)
}
