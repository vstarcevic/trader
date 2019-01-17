package api

import "github.com/vstarcevic/trader/trader"

// Service for api
type Service struct {
	ContactRepo      trader.ContactRepository
	SubscriptionRepo trader.SubscriptionRepository
}
