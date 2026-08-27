package driver

import (
	"database/sql/driver"
	"errors"

	// "github.com/nikhilbhatia08/EuphoriaDB/server"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type EuphoriaDBConnection struct {
	db          *EuphoriaDB
	transaction *transactions.Transaction
}

func (c *EuphoriaDBConnection) Close() error {
	return nil
}

func (c *EuphoriaDBConnection) Prepare(query string) (driver.Stmt, error) {
	return &Statement{
		conn:  c,
		query: query,
	}, nil
}

func (c *EuphoriaDBConnection) Begin() (driver.Tx, error) {
	if c.transaction != nil {
		return nil, errors.New("already in a transaction")
	}
	newTx := c.db.NewTransaction()
	c.transaction = newTx
	return &Tx{
		conn: c,
		tx:   newTx,
	}, nil
}
