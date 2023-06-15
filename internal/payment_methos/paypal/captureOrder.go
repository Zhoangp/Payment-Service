package paypal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zhoangp/Payment-Service/internal/payment_methos/paypal/model_paypal"
	"github.com/Zhoangp/Payment-Service/pkg/common"
	"io/ioutil"
	"net/http"
)

func (p *Paypal) CaptureOrder(orderID string) error {
	p.GetAccessToken("", "")

	orderCaptureURL := p.cf.Paypal.BaseUrl + p.cf.Paypal.CaptureOrderApi
	captureURL := fmt.Sprintf(orderCaptureURL, orderID)

	req, err := http.NewRequest("POST", captureURL, nil)
	if err != nil {
		fmt.Println(err)
		return common.NewCustomError(errors.New("Failed to capture"), "Failed to capture order")

	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.token.TokenType+" "+p.token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

		return common.NewCustomError(errors.New("Failed to capture"), "Failed to capture order")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return common.NewCustomError(errors.New("Failed to capture"), "Failed to capture order")

	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Println(err)
		return common.NewCustomError(errors.New("Failed to capture"), "Failed to capture order")
	}

	var captureResponse model_paypal.CaptureOrderResponse
	err = json.Unmarshal(body, &captureResponse)
	if err != nil {
		fmt.Println(err)

		return common.NewCustomError(errors.New("Failed to capture"), "Failed to capture order")
	}

	return nil
}
