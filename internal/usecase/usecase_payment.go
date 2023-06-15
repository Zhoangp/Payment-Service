package usecase

import (
	"fmt"
	"github.com/Zhoangp/Payment-Service/config"
	"github.com/Zhoangp/Payment-Service/internal/model"
	"github.com/Zhoangp/Payment-Service/internal/payment_methos/paypal"
	"github.com/Zhoangp/Payment-Service/pb"
	"github.com/Zhoangp/Payment-Service/pkg/utils"
	"strconv"
)

type Repository interface {
	NewPaypalMethod(paypal model.Paypal) error
	FindData(condition map[string]interface{}) ([]model.Payment, error)
	GetInformationPaypal(userId int) (*model.Paypal, error)
	NewPayment(payment *model.Payment) (*model.Payment, error)
	GetInforPaypal(userId int) (*model.Paypal, error)
}
type paymentUseCase struct {
	repo   Repository
	cf     *config.Config
	paypal *paypal.Paypal
	h      *utils.Hasher
}

func NewPaymentUseCase(repo Repository, cf *config.Config, paypal *paypal.Paypal, h *utils.Hasher) *paymentUseCase {
	return &paymentUseCase{repo: repo, cf: cf, paypal: paypal, h: h}
}
func (c *paymentUseCase) CheckOutWithPaypal(cart *pb.Cart, userId string) (*pb.CheckoutResponse, error) {
	linkApprove, err := c.paypal.CreateOrder(cart)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(utils.TokenPayload{cart, userId}, c.cf.Service.AccessTokenExpiredIn, c.cf.Service.Secret)

	return &pb.CheckoutResponse{Links: &pb.Links{
		Href:   linkApprove,
		Method: "GET",
		Rel:    "approve",
	}, OrderToken: &pb.Token{
		Token:  token.AccessToken,
		Expire: strconv.FormatInt(token.ExpiresAt, 10),
	}}, nil
}
func (c *paymentUseCase) CaptureWithPaypal(token string, orderId string) (*utils.TokenPayload, error) {
	err := c.paypal.CaptureOrder(orderId)
	if err != nil {
		return nil, err
	}
	res, _ := utils.ValidateJWT(token, c.cf)
	userIdDecoded, err := c.h.Decode(res.UserId)
	if err != nil {
		return nil, err
	}
	payment := model.Payment{
		UserId:       userIdDecoded,
		Total:        res.Cart.TotalPrice,
		RefundStatus: false,
	}
	for _, item := range res.Cart.Courses {
		courseIdDecoded, err := c.h.Decode(item.Id)
		if err != nil {
			return nil, err
		}

		payment.PaymentCourses = append(payment.PaymentCourses, model.PaymentCourse{
			CourseId: courseIdDecoded,
			Price:    item.Price,
			Currency: item.Currency,
			Amount:   item.Price,
			Title:    item.Title,
		})
	}
	fmt.Println(res.UserId)
	if _, err := c.repo.NewPayment(&payment); err != nil {
		return nil, err
	}

	return res, nil
}
func WidthDrawMoney() {
	//itemsPayout := []model_paypal.ItemsPayout{}
	//
	//for index, course := range res.Cart.Courses {
	//	userId, err := c.h.Decode(course.InstructorId)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	pp, err := c.repo.GetInformationPaypal(userId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	cost, err := strconv.ParseFloat(course.Price, 64)
	//	itemsPayout = append(itemsPayout, model_paypal.ItemsPayout{
	//		RecipientType: "PAYPAL_ID",
	//		Amount: model_paypal.Amount{
	//			Value:    fmt.Sprintf("%f", cost*0.6),
	//			Currency: course.Currency,
	//		},
	//		Note:                 "Thanks for using service",
	//		SenderItemID:         strconv.Itoa(28012002 + index),
	//		Receiver:             pp.PaypalId,
	//		NotificationLanguage: "en-US",
	//	})
	//}
	//if err := c.paypal.Payout(itemsPayout); err != nil {
	//	return nil, err
	//
	//}
}
func (c *paymentUseCase) ConnectPaypal(token string, userId string) (string, error) {

	t, err := c.paypal.GetAccessToken("authorization_code", token)
	if err != nil {
		return "", err
	}
	fmt.Println(t)
	user, err := c.paypal.GetUserInfor(t)
	if err != nil {
		return "", err

	}
	baseUserId, err := c.h.Decode(userId)
	if err != nil {
		return "", err

	}
	a := model.Paypal{
		Email:    user.Email,
		PaypalId: user.PayerID,
		UserId:   baseUserId,
	}
	if err := c.repo.NewPaypalMethod(a); err != nil {
		return "", err

	}
	return user.Email, nil

}
