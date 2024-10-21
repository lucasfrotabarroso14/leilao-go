package auction

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"leilao-go/configuration/logger"
	"leilao-go/internal/entity/auction_entity"
	"leilao-go/internal/internal_error"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	TimeStamp   int64                           `bson:"time_stamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auction_entity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := AuctionEntityMongo{
		Id:          auction_entity.Id,
		ProductName: auction_entity.ProductName,
		Category:    auction_entity.Category,
		Description: auction_entity.Description,
		Condition:   auction_entity.Condition,
		Status:      auction_entity.Status,
		TimeStamp:   auction_entity.TimeStamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")

	}
	return nil

}
