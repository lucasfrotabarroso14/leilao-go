package bid

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"leilao-go/configuration/logger"
	"leilao-go/internal/entity/bid_entity"
	"leilao-go/internal/internal_error"
	"time"
)

func (bd *BidRepository) FindBidByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}
	cursor, err := bd.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId), err)
		return nil, internal_error.NewInternalServerError("Error Trying to find bids by auctionId")
	}

	var bidEntityMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId), err)
		return nil, internal_error.NewInternalServerError("Error Trying to find bids by auctionId")
	}
	var bidEntities []bid_entity.Bid
	for _, bidEntityMongo := range bidEntityMongo {
		bidEntities = append(bidEntities, bid_entity.Bid{
			AuctionId: bidEntityMongo.AuctionId,
			Id:        bidEntityMongo.Id,
			UserId:    bidEntityMongo.UserId,
			Amount:    bidEntityMongo.Amount,
			TimesTamp: time.Unix(bidEntityMongo.TimeStamp, 0),
		})
	}
	return bidEntities, nil

}

func (bd *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}
	opts := options.FindOne().SetSort(bson.D{{"amount", -1}})
	var bidEntityMongo BidEntityMongo
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionId %s", auctionId), err)
		return nil, internal_error.NewInternalServerError("Error Trying to find bids by auctionId")
	}
	return &bid_entity.Bid{
		AuctionId: bidEntityMongo.AuctionId,
		Id:        bidEntityMongo.Id,
		UserId:    bidEntityMongo.UserId,
		Amount:    bidEntityMongo.Amount,
		TimesTamp: time.Unix(bidEntityMongo.TimeStamp, 0),
	}, nil
}
