package auction_controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"leilao-go/configuration/rest_err"
	"leilao-go/internal/entity/auction_entity"
	"leilao-go/internal/usecase/auction_usecase"
	"net/http"
	"strconv"
)

func (u *AuctionController) FindAuctionById(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid Fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})
		c.JSON(errRest.Code, errRest)
		return
	}
	auctionData, err := u.auctionUseCase.FindAuctionById(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return

	}
	c.JSON(http.StatusOK, auctionData)

}
func (u *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	category := c.Query("category")
	productName := c.Query("productName")

	statusNumber, errConv := strconv.Atoi(status)
	if errConv != nil {
		errRest := rest_err.NewBadRequestError("Error trying to validate auction status param")
		c.JSON(errRest.Code, errRest)
		return
	}

	auctions, err := u.auctionUseCase.FindAuctions(context.Background(),
		auction_entity.AuctionStatus(auction_usecase.AuctionStatus(statusNumber)), category, productName)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctions)
}
func (u *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	userId := c.Param("auctionId")

	if err := uuid.Validate(userId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid Fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})
		c.JSON(errRest.Code, errRest)
		return
	}
	auctionData, err := u.auctionUseCase.FindWinningBidByAuctionId(context.Background(), userId)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return

	}
	c.JSON(http.StatusOK, auctionData)

}
