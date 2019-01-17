package api

// ContactCreateResponse is api response model
type ContactCreateResponse struct {
	Contactid int `json:"contactid"`
}

// CheckSubscriptionResponse is api response model
type CheckSubscriptionResponse struct {
	SubscriptionExists bool `json:"subscriptionexists"`
}

// ContactSubscriptionResponse is api response
type ContactSubscriptionResponse struct {
	ContactResponse      ContactResponse      `json:"contact"`
	SubscriptionResponse SubscriptionResponse `json:"subscriptionresponse"`
}

// ContactResponse is api response model
type ContactResponse struct {
	Clientid   string `json:"clientid"`
	Broker     string `json:"broker"`
	Country    string `json:"country"`
	Language   string `json:"language"`
	Identifier string `json:"identifier"`
}

// SubscriptionResponse is api response model
type SubscriptionResponse struct {
	SubscriptionType string `json:"subscriptiontype"`
	AccountCode      string `json:"accountcode"`
}
