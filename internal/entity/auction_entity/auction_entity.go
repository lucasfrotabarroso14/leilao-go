package auction_entity

import (
	"context"
	"github.com/google/uuid"
	"leilao-go/internal/internal_error"
	"time"
)

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	TimeStamp   time.Time
}

type ProductCondition int
type AuctionStatus int

func CreateAuctions(productName string, category, description string, condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		TimeStamp:   time.Now(),
	}
	if err := auction.Validate(); err != nil {
		return nil, err
	}
	return auction, nil
}

func (a *Auction) Validate() *internal_error.InternalError {
	if len(a.ProductName) <= 0 ||
		len(a.Category) <= 2 ||
		len(a.Description) <= 10 && a.Condition != New {
		return internal_error.NewBadRequestError("invalid auction object")
	}
	return nil
}

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auction_entity *Auction) *internal_error.InternalError
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) ([]Auction, *internal_error.InternalError)
	FindAuctionById(ctx context.Context, id string) (*Auction, *internal_error.InternalError)
}
