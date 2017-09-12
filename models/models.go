package models

type ErrorResponse struct {
	ID          string `json:"error"`             // eg: unauthorized
	Description string `json:"error_description"` // eg: "full auth required to access this resource"
}
type OfferResponse struct {
	CardNumber string   `json:"cardNumber"`
	Offers     []Coupon `json:"offers"`
	//Facets     map[string]struct{} `json:"facets"`
}

type Coupon struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"expirationDate"`
	URL            string  `json:"url"`
	Loaded         bool    `json:"loaded"`
	LegalText      string  `json:"legalText"`
	Title          string  `json:"title"`
	Price          float32 `json:"price"`
	CouponSource   string  `json:"couponSource"`
	CouponCategory string  `json:"couponCategory"`
	PriceQualifier string  `json:"priceQualifier"`
	Source         string  `json:"source"`
}

// for uploads
type CouponPayload struct {
	CouponID string `json:"offerNumber"`
}
type Profile struct {
	// there's also a couple duplicate fields like cardNumber and firstName...
	CardNumber     string `json:"card_number"`
	FirstName      string `json:"first_name"`
	ID             string `json:"id"`
	Login          string `json:"login"`
	PreferredStore string `json:"preferred_store"`
	StoreNumber    string `json:"storeNumber"`
}
type Token struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
	ExpiresIn    *int    `json:"expires_in"`
	TokenType    string  `json:"token_type"`
	Scope        string  `json:"scope"`
}
