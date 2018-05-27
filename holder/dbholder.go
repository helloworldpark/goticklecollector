package holder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	"github.com/helloworldpark/goticklecollector/collector"
)

type coinBuffer []collector.Coin

// DBCredential contains info needed to connect to DB
type DBCredential struct {
	instanceConnectionName string
	databaseUser           string
	password               string
	dbName                 string
	tableName              string
}

// DBHolder provides an interface to writing to DB
type DBHolder struct {
	isInitialized bool
	buffers       []*coinBuffer
	bufferIndex   int
	capacity      int
	credential    DBCredential

	db        *sql.DB
	dbName    string
	tableName string
}

// LoadCredential load DB credential from json file
func LoadCredential(filePath string) DBCredential {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var cred DBCredential
	if err := json.Unmarshal(raw, &cred); err != nil {
		panic(err)
	}
	return cred
}

// NewDBHolder generates DBHolder from credential.
// Call this only once.
func NewDBHolder(credential DBCredential, capacity int) DBHolder {
	buffers := make([]*coinBuffer, 2)
	b1 := make(coinBuffer, 0, capacity)
	b2 := make(coinBuffer, 0, capacity)
	buffers[0] = &b1
	buffers[1] = &b2

	dbHolder := DBHolder{
		buffers:     buffers,
		bufferIndex: 0,
		capacity:    capacity,
		credential:  credential}
	return dbHolder
}

// Init initializes DB
func (h *DBHolder) Init() {
	if h.isInitialized {
		return
	}
	h.isInitialized = true
	config := mysql.Cfg(
		h.credential.instanceConnectionName,
		h.credential.databaseUser,
		h.credential.password)
	config.DBName = h.credential.dbName

	h.dbName = h.credential.dbName
	h.tableName = h.credential.tableName
	db, err := mysql.DialCfg(config)
	if err != nil {
		panic(err)
	}
	h.db = db
}

// ConnectDBChannel connects input to DB writing buffer.
func (h DBHolder) ConnectDBChannel(chans []<-chan collector.Coin) {
	fanIn := merge(chans...)
	go func() {
		for coin := range fanIn {
			buffer := h.buffers[h.bufferIndex]
			*buffer = append(*buffer, coin)
			if len(*buffer) >= h.capacity {
				// Flush to DB
				h.flush()
				// Switch active buffer
				h.switchBuffer()
			}
		}
	}()
}

func (h *DBHolder) switchBuffer() {
	newBuffer := make(coinBuffer, 0)
	h.buffers[h.bufferIndex] = &newBuffer
	h.bufferIndex++
	h.bufferIndex %= 2
}

func (h DBHolder) flush() {
	tx, err := h.db.Begin()
	if err != nil {
		tx.Rollback()
		return
	}
	// vendor, currency, price, qty, timestamp
	qstring := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?)", h.tableName)
	buffer := h.buffers[h.bufferIndex]
	for _, coin := range *buffer {
		_, err := tx.Exec(qstring, coin.Vendor, coin.Currency, coin.Price, coin.Qty, coin.Timestamp)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

func merge(cs ...<-chan collector.Coin) <-chan collector.Coin {
	var wg sync.WaitGroup
	out := make(chan collector.Coin)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan collector.Coin) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
