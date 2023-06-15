package http

import (
	"context"
	"github.com/Zhoangp/Payment-Service/pb"
)

func (hdl paymentHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	res, err := hdl.uc.GetPayment(req)
	if err != nil {
		return &pb.GetPaymentResponse{Error: HandleError(err)}, nil
	}
	return res, nil
}
