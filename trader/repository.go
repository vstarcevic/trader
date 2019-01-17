package trader

// Contact entity
type Contact struct {
	Contactid  int
	Clientid   int
	Broker     string
	Country    string
	Language   string
	Identifier string
}

// ContactRepository specify Contact persistance API
type ContactRepository interface {

	// GetByID returns Contact for given id
	GetByID(string) (Contact, error)
	Create(contact *Contact) (int, error)
	GetIDByData(clientid string, language string, identifier string) int
}

// Subscription entity
type Subscription struct {
	Subscriptionid     int
	Subscriptiontypeid int
	Account            string
}

// SubscriptionRepository specify Subscription persistance API
type SubscriptionRepository interface {
	FindByAny(search string) ([]ContactSubscription, error)

	Create(subscription *Subscription) (int, error)

	CreateType(typeName string) (int, error)
	GetSubscriptionTypeIDFromName(name string) int

	GetSubscriptionIDByTypeAndAccount(subscriptionTypeid int, account string) int
	GetFirstSubscriptionForType(subscriptionTypeid string) int
	CheckSubscription(clientid string, subscriptionType string) (bool, error)
	CheckSubscriptionByIDs(contactid int, subscriptionid int) (bool, error)
	CreateContactSubscription(contactid int, subscriptionid int) error
}

// ContactSubscription entity
type ContactSubscription struct {
	Contact          Contact
	Subscription     Subscription
	SubscriptionType SubscriptionType
}

// SubscriptionType entity
type SubscriptionType struct {
	SubscriptionTypeid int
	Name               string
}
