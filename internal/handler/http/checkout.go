package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zhoangp/Payment-Service/pb"
	"github.com/Zhoangp/Payment-Service/pkg/client"
	"github.com/Zhoangp/Payment-Service/pkg/common"
	"strconv"
)

func (hdl paymentHandler) CheckOutWithPaypal(ctx context.Context, req *pb.CheckoutRequest) (*pb.CheckoutResponse, error) {
	cartClient, err := client.InitCartServiceClient(hdl.cf)
	if err != nil {
		return nil, err
	}
	res, err := cartClient.GetCart(ctx, &pb.GetCartRequest{
		Id: req.UserId,
	})
	total, err := strconv.Atoi(res.TotalCourse)
	if err != nil {
		return &pb.CheckoutResponse{
			Error: HandleError(err),
		}, nil
	}
	if total == 0 {
		return &pb.CheckoutResponse{
			Error: HandleError(common.NewCustomError(errors.New("\"Cart is empty\""), "Cart is empty")),
		}, nil
	}
	if err != nil {
		fmt.Println(err)
	}
	response, err := hdl.uc.CheckOutWithPaypal(res.Cart, req.UserId)
	if err != nil {
		return nil, err
	}
	return response, nil
}
