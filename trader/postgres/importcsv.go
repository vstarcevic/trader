package postgres

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/vstarcevic/trader/trader"
)

// CsvContactData represent one row in CSV file
type CsvContactData struct {
	Clientid         string
	Broker           string
	Country          string
	Language         string
	Identifier       string
	SubscriptionType string
	AccountCode      string
}

// ImportCsv struct represent import for postgres
type ImportCsv struct {
	db               *sql.DB
	filename         string
	contactRepo      trader.ContactRepository
	subscriptionRepo trader.SubscriptionRepository
}

// NewImport represent instance of new import with given parameters
func NewImport(dbconn *sql.DB, filename string) ImportCsv {
	return ImportCsv{
		db:       dbconn,
		filename: filename,
	}
}

// ImportContactData imports rows into corresponding tables in db
func (i ImportCsv) ImportContactData() {
	go i.readContactRow()
}

func (i ImportCsv) readContactRow() {

	csvFile, _ := os.Open(i.filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	isFirstLine := true
	for {
		line, error := reader.Read()
		if isFirstLine {
			isFirstLine = false
			continue
		}
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		row := CsvContactData{
			Clientid:         line[0],
			Broker:           line[1],
			Country:          line[2],
			Language:         line[3],
			Identifier:       line[4],
			SubscriptionType: line[5],
			AccountCode:      line[6]}

		i.insertContactSubscriptionRow(&row)

	}

}

func (i ImportCsv) insertContactSubscriptionRow(row *CsvContactData) {

	tx, err := i.db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction ", err)
	}
	defer tx.Rollback()

	i.contactRepo = NewContactRepository(tx)
	i.subscriptionRepo = NewSubscriptionRepository(tx)

	contid := i.contactRepo.GetIDByData(row.Clientid, row.Language, row.Identifier)
	if contid == 0 {
		clientid, _ := strconv.Atoi(row.Clientid)

		contid, _ = i.contactRepo.Create(&trader.Contact{
			Clientid:   clientid,
			Broker:     row.Broker,
			Country:    row.Country,
			Language:   row.Language,
			Identifier: row.Identifier,
		})
	}

	subsTypeid := i.subscriptionRepo.GetSubscriptionTypeIDFromName(row.SubscriptionType)
	if subsTypeid == 0 {
		subsTypeid, _ = i.subscriptionRepo.CreateType(row.SubscriptionType)
	}

	subscriptionid := i.subscriptionRepo.GetSubscriptionIDByTypeAndAccount(subsTypeid, row.AccountCode)
	if subscriptionid == 0 {
		subscriptionid, _ = i.subscriptionRepo.Create(&trader.Subscription{Subscriptiontypeid: subsTypeid, Account: row.AccountCode})
	}

	subscriptionExists, _ := i.subscriptionRepo.CheckSubscriptionByIDs(contid, subscriptionid)
	if !subscriptionExists {
		i.subscriptionRepo.CreateContactSubscription(contid, subscriptionid)
	}

	tx.Commit()

}
