package payment

import "time"

type PayResponse struct {
	ID           string     `json:"id"`
	Status       string     `json:"status"`
	Paid         bool       `json:"paid"`
	Amount       Amount     `json:"amount"`
	Confirmation ConfirmURL `json:"confirmation"`
	CreatedAt    time.Time  `json:"created_at"`
	Description  string     `json:"description"`
	Metadata     Metadata   `json:"metadata"`
	Recipient    Recipient  `json:"recipient"`
	Refundable   bool       `json:"refundable"`
	Test         bool       `json:"test"`
}

type PayResponseToClient struct {
	ID           string     `json:"id"`
	Status       string     `json:"status"`
	Paid         bool       `json:"paid"`
	Confirmation ConfirmURL `json:"confirmation"`
}

type PayResponseToClientData struct {
	Data PayResponseToClient `json:"data"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ConfirmURL struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Metadata struct {
	UserLogin string `json:"user_login"`
	OrderID   string `json:"order_id"`
	Type      string `json:"type"`
}

type Recipient struct {
	AccountID string `json:"account_id"`
	GatewayID string `json:"gateway_id"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type ItemCustomer struct {
	Description    string `json:"description"`
	Amount         Amount `json:"amount"`
	VatCode        int    `json:"vat_code"`
	Quantity       int    `json:"quantity"`
	PaymentSubject string `json:"payment_subject"`
}

type Customer struct {
	Email string `json:"email"`
}

type Receipt struct {
	Customer Customer       `json:"customer"`
	Items    []ItemCustomer `json:"items"`
}

type Payment struct {
	Amount       Amount       `json:"amount"`
	Capture      bool         `json:"capture"`
	Confirmation Confirmation `json:"confirmation"`
	Description  string       `json:"description"`
	Receipt      Receipt      `json:"receipt"`
	Metadata     Metadata     `json:"metadata"`
}

type Notification struct {
	Type   string          `json:"type"`
	Event  string          `json:"event"`
	Object PayNotification `json:"object"`
}

type PayNotification struct {
	ID                   string               `json:"id"`
	Status               string               `json:"status"`
	Paid                 bool                 `json:"paid"`
	Amount               Amount               `json:"amount"`
	AuthorizationDetails AuthorizationDetails `json:"authorization_details"`
	CreatedAt            time.Time            `json:"created_at"`
	Description          string               `json:"description"`
	ExpiresAt            time.Time            `json:"expires_at"`
	Metadata             Metadata             `json:"metadata"`
	PaymentMethod        PayMethod            `json:"payment_method"`
	Refundable           bool                 `json:"refundable"`
	Test                 bool                 `json:"test"`
}

type AuthorizationDetails struct {
	RRN          string       `json:"rrn"`
	AuthCode     string       `json:"auth_code"`
	ThreeDSecure ThreeDSecure `json:"three_d_secure"`
}

type ThreeDSecure struct {
	Applied bool `json:"applied"`
}

type PayMethod struct {
	Type  string      `json:"type"`
	ID    string      `json:"id"`
	Saved bool        `json:"saved"`
	Card  CardDetails `json:"card"`
	Title string      `json:"title"`
}

type CardDetails struct {
	First6        string `json:"first6"`
	Last4         string `json:"last4"`
	ExpiryMonth   string `json:"expiry_month"`
	ExpiryYear    string `json:"expiry_year"`
	CardType      string `json:"card_type"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}
