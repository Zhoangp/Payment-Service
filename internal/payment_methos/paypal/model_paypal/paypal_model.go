package model_paypal

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
type CaptureOrderResponse struct {
	Status string `json:"status"`
}
