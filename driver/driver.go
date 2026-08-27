package driver

import (
	"database/sql"
	"database/sql/driver"
	"sync"
)

var _ driver.Driver = (*EuphoriaDBDriver)(nil)

type EuphoriaDBDriver struct{}

var (
	instancesMu sync.Mutex
	instances   = map[string]*EuphoriaDB{}
)

func init() {
	sql.Register("euphoriadb", &EuphoriaDBDriver{})
}

func getOrCreateDB(directory string) (*EuphoriaDB, error) {
	instancesMu.Lock()
	defer instancesMu.Unlock()

	if db, ok := instances[directory]; ok {
		return db, nil
	}
	db, err := NewEuphoriaDB(directory)
	if err != nil {
		return nil, err
	}
	instances[directory] = db
	return db, nil
}

func (d *EuphoriaDBDriver) Open(directory string) (driver.Conn, error) {
	db, err := getOrCreateDB(directory)
	if err != nil {
		return nil, err
	}
	return &EuphoriaDBConnection{db: db}, nil
}
