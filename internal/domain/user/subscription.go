package user

type Subscription string

const (
	SubscriptionFree    Subscription = "FREE"
	SubscriptionPremium Subscription = "PREMIUM"
)

func (s Subscription) IsValid() bool {
	return s == SubscriptionFree || s == SubscriptionPremium
}
