package bid_usecase

import (
	"context"
	"leilao-go/internal/internal_error"
)

func (bu *BidUseCase) FindBidByAuctionId(ctx context.Context, auctionId string) ([]BidOutputDTO, *internal_error.InternalError) {
	bidList, err := bu.bidRepository.FindBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}
	var bidOutputList []BidOutputDTO
	for _, bid := range bidList {
		bidOutputList = append(bidOutputList, BidOutputDTO{
			Id:        bid.Id,
			UserId:    bid.UserId,
			Amount:    bid.Amount,
			AuctionId: bid.AuctionId,
			TimesTamp: bid.TimesTamp,
		})
	}
	return bidOutputList, nil

}

func (bu *BidUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*BidOutputDTO, *internal_error.InternalError) {
	bidEntity, err := bu.bidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	bidOutput := &BidOutputDTO{
		Id:        bidEntity.Id,
		UserId:    bidEntity.UserId,
		Amount:    bidEntity.Amount,
		AuctionId: bidEntity.AuctionId,
		TimesTamp: bidEntity.TimesTamp,
	}
	return bidOutput, nil
}
