package addrbook

import (
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

type AddrBook struct {
	IsMaster  bool
	Addrs     map[string]string
	datastore *bolt.DB
	lock      sync.Mutex
	readIdx   int
}

func New(isMaster bool, masterDBPath string) (*AddrBook, error) {
	var db *bolt.DB
	if isMaster {
		db, err := bolt.Open(masterDBPath, 0600, nil)
		if err != nil {
			return nil, err
		}
	}
	return &AddrBook{
		IsMaster:  isMaster,
		Addrs:     make(map[string]string),
		datastore: db,
	}, nil
}

func getKey() string {
	today := time.Now()
	return fmt.Sprintf("%d-%d-%d", today.Month(), today.Day(), today.Year())
}

func (a *AddrBook) save(tx *bolt.Tx) error {
	key := getKey()
	bucket, err := tx.CreateBucketIfNotExists([]byte(key))
	if err != nil {
		return err
	}
	inner, err := bucket.CreateBucketIfNotExists([]byte("addresses"))
	if err != nil {
		return err
	}
	for k, v := range a.Addrs {
		err = inner.Put([]byte(k), []byte(v))
		if err != nil {
			return err
		}
	}
}

func (a *AddrBook) find(tx *bolt.Tx) error {
	key := getKey()
	bucket := tx.Bucket(key)
	inner := bucket.Bucket([]byte("addresses"))
	cursor := inner.Cursor()
	k, v := cursor.First()
	for len(v) > 0 {
		a.Addrs[k] = v
		k, v = cursor.Next()
	}
}

func (a *AddrBook) read(p []byte) (n int, err error) {
	payload, err := json.Marshal(a.Addrs)
	if err != nil {
		return 0, err
	}
	for i := range payload[a.readIdx:] {
		if i >= len(p) {
			a.readIdx = i
			return i, nil
		}
		p[i] = payload[i]
	}
	a.readIdx = 0
	return len(p), nil
}

func (a *AddrBook) readMaster(p []byte) (n int, err error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	err := a.datastore.View(a.find)
	if err != nil {
		return 0, err
	}
	return a.read(p)
}

func (a *AddrBook) write(p []byte) (n int, err error) {
	additional := make(map[string]string)
	err := json.Unmarshal(p, addtional)
	if err != nil {
		return 0, err
	}
	for k, v := range additional {
		a.Addrs[k] = v
	}
	return len(p), nil
}

func (a *AddrBook) writeMaster(p []byte) (n int, err error) {
	length, err := a.write(p)
	if err != nil {
		return length, err
	}
	defer a.lock.Unlock()
	a.lock.Lock()
	err = a.datastore.Update(a.save)
	return length, err
}

func (a *AddrBook) Read(p []byte) (n int, err error) {
	if a.IsMaster {
		return a.readMaster(p)
	}
	return a.read(p)
}

func (a *AddrBook) Write(p []byte) (n int, err error) {
	if a.IsMaster {
		return a.writeMaster(p)
	}
	return a.write(p)
}

func (a *AddrBook) Close() error {
	if a.IsMaster {
		return a.datastore.Close()
	}
	return nil
}
