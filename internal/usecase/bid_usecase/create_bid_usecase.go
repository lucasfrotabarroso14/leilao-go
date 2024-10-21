package bid_usecase

import (
	"context"
	"leilao-go/configuration/logger"
	"leilao-go/internal/entity/bid_entity"
	"leilao-go/internal/internal_error"
	"os"
	"strconv"
	"time"
)

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}
type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	AuctionId string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	TimesTamp time.Time `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type BidUseCase struct {
	bidRepository       bid_entity.BidEntityRepository
	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

func NewBidUseCase(bidRepository bid_entity.BidEntityRepository) NewBidUseCaseInterface {
	maxSizeInterval := getMaxBatchSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		bidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}
	bidUseCase.TriggerCreateRoutine(context.Background())
	return bidUseCase
}

var bidBatch []bid_entity.Bid

type NewBidUseCaseInterface interface {
	FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError)
	FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError)
	CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError
}

func (bu *BidUseCase) TriggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)
		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						if err := bu.bidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("Error trying to process bid batch list", err)
						}
					}
					return
				}
				bidBatch = append(bidBatch, bidEntity)
				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.bidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("Error trying to process bid batch list", err)
					}
					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				if err := bu.bidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("Error trying to process bid batch list", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)

			}

		}
	}()
}
func (bu *BidUseCase) CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError {

	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.AuctionId, bidInputDTO.Amount)
	if err != nil {
		return err

	}
	bu.bidChannel <- *bidEntity
	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}
	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}
	return value
}
