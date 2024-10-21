package auction

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"leilao-go/configuration/logger"
	"leilao-go/internal/entity/auction_entity"
	"leilao-go/internal/internal_error"
	"time"
)

func (ar *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}
	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auction by id =%s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		TimeStamp:   time.Unix(auctionEntityMongo.TimeStamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {

	filter := bson.M{}

	if status != 0 {
		filter["status"] = status

	}
	if category != "" {
		filter["category"] = category

	}

	if productName != "" {
		filter["productName"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}

	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auctions "), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction ")
	}
	defer cursor.Close(ctx)

	var auctionEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionEntityMongo); err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internal_error.NewInternalServerError("Error trying to find auctions")
	}
	var auctionEntity []auction_entity.Auction
	for _, auctionMongo := range auctionEntityMongo {
		auctionEntity = append(auctionEntity, auction_entity.Auction{
			Id:          auctionMongo.Id,
			ProductName: auctionMongo.ProductName,
			Category:    auctionMongo.Category,
			Description: auctionMongo.Description,
			Condition:   auctionMongo.Condition,
			Status:      auctionMongo.Status,
			TimeStamp:   time.Unix(auctionMongo.TimeStamp, 0),
		})

	}
	return auctionEntity, nil

}
