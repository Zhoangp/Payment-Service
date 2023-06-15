package http

import (
	"context"
	"github.com/Zhoangp/Payment-Service/pb"
)

func (hdl *paymentHandler) GetPaypal(ctx context.Context, req *pb.GetPayalRequest) (*pb.GetPayalResponse, error) {
	res, err := hdl.uc.GetPaypal(req)
	if err != nil {
		return &pb.GetPayalResponse{Error: HandleError(err)}, nil
	}
	return res, nil
}
