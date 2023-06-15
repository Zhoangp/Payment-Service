package usecase

import (
	"fmt"
	"github.com/Zhoangp/Payment-Service/pb"
)

func (uc *paymentUseCase) GetPayment(req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {

	userIdDecoded, err := uc.h.Decode(req.UserId)
	if err != nil {
		return nil, err
	}
	payments, err := uc.repo.FindData(map[string]any{"user_id": userIdDecoded})
	if err != nil {
		return nil, err
	}
	var res pb.GetPaymentResponse
	for _, payment := range payments {
		var item pb.Payment
		item.PaymentId = uc.h.Encode(payment.Id)
		item.Total = payment.Total
		item.PaymentDate = payment.CreatedAt.Format("2006-01-02 15:04:05")
		for _, paymentCourse := range payment.PaymentCourses {
			item.ListItem = append(item.ListItem, &pb.ListItemPayment{
				Title:    paymentCourse.Title,
				Price:    paymentCourse.Price,
				Discount: fmt.Sprintf("%f", paymentCourse.Discount),
				Amount:   paymentCourse.Amount,
			})
		}
		item.Currency = "USD"
		res.Payment = append(res.Payment, &item)

	}

	return &res, nil
}
