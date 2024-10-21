package auction_usecase

import (
	"context"
	"leilao-go/internal/entity/auction_entity"
	"leilao-go/internal/internal_error"
	"leilao-go/internal/usecase/bid_usecase"
)

func (au *AuctionUseCase) FindAuctionById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {

	auctionEntity, err := au.auctionRepositoryInterface.FindAuctionById(ctx, id)
	if err != nil {
		return nil, err

	}
	return &AuctionOutputDTO{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		TimeStamp:   auctionEntity.TimeStamp,
	}, nil

}

func (au *AuctionUseCase) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntities, err := au.auctionRepositoryInterface.FindAuctions(ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err

	}

	var AuctionsOutPuts []AuctionOutputDTO
	for _, value := range auctionEntities {
		AuctionsOutPuts = append(AuctionsOutPuts, AuctionOutputDTO{
			Id:          value.Id,
			ProductName: value.ProductName,
			Category:    value.Category,
			Description: value.Description,
			Condition:   ProductCondition(value.Condition),
			Status:      AuctionStatus(value.Status),
			TimeStamp:   value.TimeStamp,
		})

	}
	return AuctionsOutPuts, nil

}

func (au *AuctionUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auction, err := au.auctionRepositoryInterface.FindAuctionById(ctx, auctionId)
	if err != nil {
		return nil, err
	}
	auctionOutputDTO := AuctionOutputDTO{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   ProductCondition(auction.Condition),
		Status:      AuctionStatus(auction.Status),
		TimeStamp:   auction.TimeStamp,
	}
	bidWinning, err := au.bidRepositoryInterface.FindWinningBidByAuctionId(ctx, auction.Id)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutputDTO,
			Bid:     nil,
		}, nil
	}

	bidOutputDTO := &bid_usecase.BidOutputDTO{
		Id:        bidWinning.Id,
		UserId:    bidWinning.UserId,
		AuctionId: bidWinning.AuctionId,
		Amount:    bidWinning.Amount,
		TimesTamp: bidWinning.TimesTamp,
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		Bid:     bidOutputDTO,
	}, nil

}
