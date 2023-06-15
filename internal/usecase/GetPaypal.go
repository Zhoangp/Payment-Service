package usecase

import (
	"errors"
	"github.com/Zhoangp/Payment-Service/pb"
	"github.com/Zhoangp/Payment-Service/pkg/common"
)

func (uc *paymentUseCase) GetPaypal(req *pb.GetPayalRequest) (*pb.GetPayalResponse, error) {
	userIdDecoded, err := uc.h.Decode(req.UserId)
	if err != nil {
		return nil, err
	}
	paypal, err := uc.repo.GetInforPaypal(userIdDecoded)
	if paypal == nil {
		return nil, common.NewCustomError(errors.New("Paypal account not found"), "Paypal account not found")
	}
	return &pb.GetPayalResponse{
		Email: paypal.Email,
	}, nil
}
