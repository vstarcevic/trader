package postgres

import (
	"database/sql"
	"log"

	"github.com/vstarcevic/trader/trader"
)

// SubscriptionRepository represent repository
type SubscriptionRepository struct {
	db *sql.DB
	tr *sql.Tx
}

// NewSubscriptionRepository returns repository object
func NewSubscriptionRepository(dbconn *sql.DB) trader.SubscriptionRepository {
	return &SubscriptionRepository{db: dbconn}
}

// CheckSubscription checks if subscription exists for given clientid and subscription type
func (sup SubscriptionRepository) CheckSubscription(clientid string, subscriptionType string) (bool, error) {
	q := ` 
	SELECT count(1) FROM contactSubscription cs 
	INNER JOIN contact c ON cs.contactid = c.contactid
	INNER JOIN subscription s ON s.subscriptionid = cs.subscriptionid
	INNER JOIN subscriptiontype st ON st.subscriptiontypeid = s.subscriptiontypeid
	WHERE c.clientid = $1 
	AND st.name = $2`

	var result int

	err := sup.db.QueryRow(q, clientid, subscriptionType).Scan(&result)

	if err != nil {
		return false, err
	}

	if result < 1 {
		return false, nil
	}
	return true, nil
}

// CheckSubscriptionByIDs checks if subscription exists for given ids in table
func (sup SubscriptionRepository) CheckSubscriptionByIDs(contactid int, subscriptionid int) (bool, error) {
	q := `SELECT COUNT(1) FROM contactSubscription  
	WHERE contactid = $1 AND subscriptionid = $2`

	var result int

	err := sup.db.QueryRow(q, contactid, subscriptionid).Scan(&result)

	if err != nil {
		return false, err
	}

	if result < 1 {
		return false, nil
	}
	return true, nil
}

// GetSubscriptionTypeIDFromName returns id of SubscriptionType, or 0 if not found
func (sup SubscriptionRepository) GetSubscriptionTypeIDFromName(name string) int {
	query := `
	SELECT subscriptionTypeid FROM subscriptiontype WHERE name = $1`

	var subscriptionTypeid int

	err := sup.db.QueryRow(query, name).Scan(&subscriptionTypeid)

	if err != nil {
		return 0
	}

	return subscriptionTypeid
}

// GetSubscriptionIDByTypeAndAccount returns subscriptionId for given subscriptiontypeid and account, 0 if not found
func (sup SubscriptionRepository) GetSubscriptionIDByTypeAndAccount(subscriptionTypeid int, account string) int {
	query := `SELECT subscriptionid FROM subscription WHERE subscriptiontypeid = $1 AND account = $2`

	var subscriptionid int

	err := sup.db.QueryRow(query, subscriptionTypeid, account).Scan(&subscriptionid)

	if err != nil {
		return 0
	}

	return subscriptionid

}

// CreateType creates SubscriptionType and returns subscriptionTypeId
func (sup SubscriptionRepository) CreateType(name string) (int, error) {

	query := `INSERT INTO subscriptionType (name)
	VALUES ($1) RETURNING subscriptionTypeid`

	var subscriptionTypeid int

	err := sup.db.QueryRow(query, name).Scan(&subscriptionTypeid)

	if err != nil {
		return 0, err
	}

	return subscriptionTypeid, nil

}

// Create creates subscription with given model
func (sup SubscriptionRepository) Create(subscription *trader.Subscription) (int, error) {

	query := `
	INSERT INTO Subscription (subscriptionTypeid, account)
	VALUES ($1, $2) RETURNING subscriptionid`

	var subsid int

	err := sup.db.QueryRow(query, subscription.Subscriptiontypeid, subscription.Account).Scan(&subsid)

	if err != nil {
		return 0, err
	}

	return subsid, nil
}

// CreateContactSubscription creates subscription
func (sup SubscriptionRepository) CreateContactSubscription(contactid int, subscriptionid int) error {

	q := `INSERT INTO contactSubscription (contactid,subscriptionid)
	VALUES ($1, $2)`

	_, err := sup.db.Exec(q, contactid, subscriptionid)

	return err

}

// FindByAny searches for given columns
func (sup SubscriptionRepository) FindByAny(search string) ([]trader.ContactSubscription, error) {

	q := `
	SELECT c.contactid,
	c.clientid, 
	c.broker,
	c.country,
	c.language,
	c.identifier,
	st.name as subscriptionType,
	s.account	
	FROM contactSubscription cs
	FULL JOIN contact c ON cs.contactid = c.contactid
	FULL JOIN subscription s ON s.subscriptionid = cs.subscriptionid
	FULL JOIN subscriptiontype st ON st.subscriptiontypeid = s.subscriptiontypeid
	WHERE
	cast(c.clientid as varchar) = $1
	OR c.broker = $1
	OR c.country = $1
	OR c.language = $1
	OR c.identifier = $1
	UNION  
	SELECT c.contactid,
	c.clientid, 
	c.broker,
	c.country,
	c.language,
	c.identifier,
	st.name as subscriptionType,
	s.account	
	FROM contactSubscription cs
	FULL JOIN contact c ON cs.contactid = c.contactid
	FULL JOIN subscription s ON s.subscriptionid = cs.subscriptionid
	FULL JOIN subscriptiontype st ON st.subscriptiontypeid = s.subscriptiontypeid
	WHERE account = $1
	UNION  
	SELECT c.contactid,
	c.clientid, 
	c.broker,
	c.country,
	c.language,
	c.identifier,
	st.name as subscriptionType,
	s.account	
	FROM contactSubscription cs
	FULL JOIN contact c ON cs.contactid = c.contactid
	FULL JOIN subscription s ON s.subscriptionid = cs.subscriptionid
	FULL JOIN subscriptiontype st ON st.subscriptiontypeid = s.subscriptiontypeid
	WHERE st.name = $1
	`

	rows, err := sup.db.Query(q, search)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var res []trader.ContactSubscription
	var dataRow trader.ContactSubscription
	for rows.Next() {
		err := rows.Scan(&dataRow.Contact.Contactid,
			&dataRow.Contact.Clientid,
			&dataRow.Contact.Broker,
			&dataRow.Contact.Country,
			&dataRow.Contact.Language,
			&dataRow.Contact.Identifier,
			&dataRow.SubscriptionType.Name,
			&dataRow.Subscription.Account)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		res = append(res, dataRow)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil

}

// GetFirstSubscriptionForType this method returns first (min) subscriptionid from Subscription table, for given subscriptionTypeid
// this should not be like that, but it was not so clear in task so it is implemented this way
// e.g. "Subscribe client to new subscription type" could be to create new subscription type or not? AccountCode is also
// unclear what is it's connection to Subscription, but I think that for this test purpose will do good
func (sup SubscriptionRepository) GetFirstSubscriptionForType(subscriptionTypeid string) int {
	query := `SELECT min(subscriptionid) FROM subscription WHERE subscriptiontypeid = $1`

	var subscriptionid int

	err := sup.db.QueryRow(query, subscriptionTypeid).Scan(&subscriptionid)

	if err != nil {
		return 0
	}

	return subscriptionid
}
