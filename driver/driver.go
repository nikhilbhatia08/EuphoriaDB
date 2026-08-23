package driver

import (
	"database/sql"
	"database/sql/driver"

	"github.com/nikhilbhatia08/EuphoriaDB/server"
)

const databaseName = "euphoriadb"

func init() {
	sql.Register(databaseName, &EuphoriaDBDriver{})
}

var _ driver.Driver = (*EuphoriaDBDriver)(nil)

type EuphoriaDBDriver struct{}

func (d *EuphoriaDBDriver) Open(directory string) (driver.Conn, error) {
	db, err := server.NewEuphoriaDB(directory)
	if err != nil {
		return nil, err
	}
	return &EuphoriaDBConnection{
		db: db,
	}, nil
}
