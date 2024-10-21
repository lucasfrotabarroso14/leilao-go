package bid

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"leilao-go/configuration/logger"
	"leilao-go/internal/entity/auction_entity"
	"leilao-go/internal/entity/bid_entity"
	"leilao-go/internal/infra/database/auction"
	"leilao-go/internal/internal_error"
	"sync"
)

type BidEntityMongo struct {
	Id        string  `bson:"_id"`
	UserId    string  `bson:"user_id"`
	AuctionId string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	TimeStamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository *auction.AuctionRepository
}

func NewBidRepository(database *mongo.Database, auctionRepository *auction.AuctionRepository) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepository,
	}
}

func (bd *BidRepository) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup
	for _, bid := range bidEntities {
		wg.Add(1)
		go func(bidValue bid_entity.Bid) {
			defer wg.Done()
			auctionEntity, err := bd.AuctionRepository.FindAuctionById(ctx, bidValue.AuctionId)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}
			if auctionEntity.Status != auction_entity.Active {
				return
			}
			var bidEntityMongo = &BidEntityMongo{
				Id:        bidValue.AuctionId,
				UserId:    bidValue.UserId,
				AuctionId: bidValue.AuctionId,
				Amount:    bidValue.Amount,
				TimeStamp: bidValue.TimesTamp.Unix(),
			}
			if _, err := bd.Collection.InsertOne(ctx, bidEntityMongo); err != nil {
				logger.Error("Error trying to insert bid", err)
				return
			}

		}(bid)
	}
	wg.Wait()
	return nil
}
