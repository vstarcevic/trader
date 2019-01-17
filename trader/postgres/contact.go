package postgres

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/vstarcevic/trader/trader"
)

// ContactRepository represents repository
type ContactRepository struct {
	db *sql.DB
}

// NewContactRepository returns repository object
func NewContactRepository(dbconn *sql.DB) trader.ContactRepository {
	return &ContactRepository{db: dbconn}
}

// GetByID gets Contact for given id, or error if not exists
func (co ContactRepository) GetByID(id string) (trader.Contact, error) {
	query := `SELECT contactid, clientid, broker, country, language, identifier FROM contact WHERE contactid = $1`
	intid, _ := strconv.Atoi(id)
	contact := trader.Contact{Contactid: intid}
	err := co.db.QueryRow(query, id).Scan(&contact.Contactid, &contact.Clientid, &contact.Broker, &contact.Country, &contact.Language, &contact.Identifier)

	if err != nil {
		empty := trader.Contact{}
		if err == sql.ErrNoRows {
			return empty, errors.New("Contact does not exists")
		}
		return empty, err
	}

	return contact, nil
}

// Create contact in db
func (co ContactRepository) Create(contact *trader.Contact) (int, error) {

	query := `
	INSERT INTO contact (clientid, broker, country, language, identifier)
	VALUES ($1, $2, $3, $4, $5) RETURNING contactid`

	var contactid int

	err := co.db.QueryRow(query, contact.Clientid, contact.Broker, contact.Country, contact.Language, contact.Identifier).Scan(&contactid)

	if err != nil {
		return 0, err
	}

	return contactid, nil
}

// GetIDByData returns contactid if found, 0 if not
func (co ContactRepository) GetIDByData(clientid string, language string, identifier string) int {

	query := `
	SELECT contactid FROM contact 
	WHERE clientid = $1 AND language = $2 AND identifier = $3`

	var contactid int

	err := co.db.QueryRow(query, clientid, language, identifier).Scan(&contactid)

	if err != nil {
		return 0
	}

	return contactid
}
