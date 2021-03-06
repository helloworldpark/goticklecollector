package holder

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/helloworldpark/goticklecollector/logger"

	_ "github.com/go-sql-driver/mysql" // SQL Connection
	"github.com/helloworldpark/goticklecollector/collector"
)

type coinBuffer []collector.Coin

// DBCredential contains info needed to connect to DB
type DBCredential struct {
	InstanceConnectionName string `json:"instanceConnectionName"`
	DatabaseUser           string `json:"databaseUser"`
	Password               string `json:"password"`
	DBName                 string `json:"dbName"`
	TableName              string `json:"tableName"`
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

// NewDBHolder generates DBHolder from credential.
// Call this only once.
func NewDBHolder(credPath string, capacity int) DBHolder {
	credential := loadCredential(credPath)

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

// loadCredential load DB credential from json file
func loadCredential(filePath string) DBCredential {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Panic("%v", err)
	}

	var cred DBCredential
	if err := json.Unmarshal(raw, &cred); err != nil {
		logger.Panic("%v", err)
	}
	return cred
}

// Init initializes DB
func (h *DBHolder) Init() {
	if h.isInitialized {
		return
	}
	h.isInitialized = true

	h.dbName = h.credential.DBName
	h.tableName = h.credential.TableName

	openingQ := openingQuery(h.credential)
	db, err := sql.Open("mysql", openingQ)
	if err != nil {
		logger.Panic("%v", err)
	}
	logger.Info("[DB] Connected")
	h.db = db
}

func openingQuery(credential DBCredential) string {
	var buf bytes.Buffer
	buf.WriteString(credential.DatabaseUser)
	buf.WriteByte(':')
	buf.WriteString(credential.Password)
	buf.WriteByte('@')
	buf.WriteByte('(')
	buf.WriteString("127.0.0.1:3306")
	buf.WriteByte(')')
	buf.WriteByte('/')
	buf.WriteString(credential.DBName)
	return buf.String()
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
				logger.Info("[DB] Flushed")
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
		logger.Error("[DB] %v", err)
		return
	}
	// vendor, currency, price, qty, timestamp
	qstring := fmt.Sprintf("INSERT INTO %s (vendor, currency, price, qty, timestamp) VALUES (?, ?, ?, ?, ?)", h.tableName)

	buffer := h.buffers[h.bufferIndex]
	for _, coin := range *buffer {
		_, err := tx.Exec(qstring, coin.Vendor, coin.Currency, coin.Price, coin.Qty, coin.Timestamp)
		if err != nil {
			logger.Error("[DB] %v", err)
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
