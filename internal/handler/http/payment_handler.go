package http

import (
	"context"
	"github.com/Zhoangp/Payment-Service/config"
	"github.com/Zhoangp/Payment-Service/pb"
	"github.com/Zhoangp/Payment-Service/pb/course"
	"github.com/Zhoangp/Payment-Service/pkg/client"
	"github.com/Zhoangp/Payment-Service/pkg/common"
	"github.com/Zhoangp/Payment-Service/pkg/utils"
)

type PaymentUseCase interface {
	CheckOutWithPaypal(cart *pb.Cart, userId string) (*pb.CheckoutResponse, error)
	CaptureWithPaypal(token string, orderId string) (*utils.TokenPayload, error)
	ConnectPaypal(token string, userId string) (string, error)
	GetPayment(req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error)
	GetPaypal(req *pb.GetPayalRequest) (*pb.GetPayalResponse, error)
}
type paymentHandler struct {
	uc PaymentUseCase
	pb.UnimplementedPaymentServiceServer
	cf *config.Config
}

func NewPaymentHandler(uc PaymentUseCase, cf *config.Config) *paymentHandler {
	return &paymentHandler{uc: uc, cf: cf}
}
func HandleError(err error) *pb.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &pb.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &pb.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
func (hdl paymentHandler) CaptureWithPaypal(ctx context.Context, request *pb.CaptureRequest) (*pb.CaptureResponse, error) {
	payload, err := hdl.uc.CaptureWithPaypal(request.Token, request.OrderId)
	if err != nil {
		return &pb.CaptureResponse{
			Error: HandleError(err),
		}, nil
	}
	cartClient, err := client.InitCartServiceClient(hdl.cf)
	if err != nil {
		return &pb.CaptureResponse{
			Error: HandleError(err),
		}, nil
	}

	_, err = cartClient.ResetCart(ctx, &pb.ResetCartRequest{
		CartId: payload.Cart.Id,
	})

	if err != nil {
		return &pb.CaptureResponse{
			Error: HandleError(err),
		}, nil
	}

	courseClient, err := client.InitCourseServiceClient(hdl.cf)
	if err != nil {
		return &pb.CaptureResponse{
			Error: HandleError(err),
		}, nil
	}
	for _, i := range payload.Cart.Courses {
		resCourse, err := courseClient.Enrollment(ctx, &course.EnrollmentRequest{
			UserId:   payload.UserId,
			CourseId: i.Id,
		})
		if err != nil {
			return &pb.CaptureResponse{
				Error: HandleError(err),
			}, nil
		}
		if resCourse.Error != nil {
			return &pb.CaptureResponse{
				Error: resCourse.Error,
			}, nil
		}
	}
	return &pb.CaptureResponse{}, nil
}
func (hdl paymentHandler) ConnectPaypalAccount(ctx context.Context, req *pb.ConnectPaypalRequest) (*pb.ConnectPaypalResponse, error) {
	email, err := hdl.uc.ConnectPaypal(req.IdentifyToken.Token, req.UserId)
	if err != nil {
		return &pb.ConnectPaypalResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.ConnectPaypalResponse{
		Email: email,
	}, nil
}
