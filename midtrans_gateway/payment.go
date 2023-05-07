package midtrans_gateway

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Config struct {
	ServerKey string
	Env       midtrans.EnvironmentType
}

type Payment struct {
	MidtransConfig *Config
	SnapClient     *snap.Client
}

func NewPayment(mc *Config) *Payment {
	var snapClient = new(snap.Client)
	snapClient.New(mc.ServerKey, mc.Env)

	return &Payment{SnapClient: snapClient, MidtransConfig: mc}
}

type CustomerDetailsRequest struct {
	FName    string `json:"first_name,omitempty"`
	LName    string `json:"last_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	GrossAmt int64  `json:"gross_amount"`
}

type CreateTransactionResponse struct {
	OrderID     string `json:"order_id"`
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

var ErrCannotCreateTransaction = fmt.Errorf("cannot create transaction")

// CreateTransaction https://docs.midtrans.com/reference/backend-integration
func (p *Payment) CreateTransaction(request *CustomerDetailsRequest) (*CreateTransactionResponse, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  uuid.NewString(),
			GrossAmt: request.GrossAmt,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.FName,
			LName: request.LName,
			Email: request.Email,
			Phone: request.Phone,
		},
	}

	snapResp, err := p.SnapClient.CreateTransaction(req)
	if err != nil {
		return nil, ErrCannotCreateTransaction
	}

	return &CreateTransactionResponse{
		OrderID:     req.TransactionDetails.OrderID,
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

type TransactionStatus struct {
	StatusCode        string `json:"status_code,omitempty"`
	TransactionID     string `json:"transaction_id,omitempty"`
	GrossAmount       string `json:"gross_amount,omitempty"`
	Currency          string `json:"currency,omitempty"`
	OrderID           string `json:"order_id,omitempty"`
	PaymentType       string `json:"payment_type,omitempty"`
	SignatureKey      string `json:"signature_key,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	StatusMessage     string `json:"status_message,omitempty"`
	MerchantID        string `json:"merchant_id,omitempty"`
	BillKey           string `json:"bill_key,omitempty"`
	BillerCode        string `json:"biller_code,omitempty"`
	TransactionTime   string `json:"transaction_time,omitempty"`
	ExpiryTime        string `json:"expiry_time,omitempty"`
}

// TransactionStatus https://docs.midtrans.com/reference/get-transaction-status
func (p *Payment) TransactionStatus(orderId string) (*TransactionStatus, error) {
	url := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/status", orderId)
	if p.MidtransConfig.Env == midtrans.Production {
		url = fmt.Sprintf("https://api.midtrans.com/v2/%s/status", orderId)
	}

	apiKey := fmt.Sprintf(p.MidtransConfig.ServerKey + ":")

	var transactionStatus TransactionStatus
	err := p.SnapClient.HttpClient.Call("GET", url, &apiKey, nil, nil, &transactionStatus)
	if err != nil {
		return nil, err
	}

	return &transactionStatus, nil
}
