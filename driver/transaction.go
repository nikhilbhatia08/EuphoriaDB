package driver

import "github.com/nikhilbhatia08/EuphoriaDB/transactions"

type Tx struct {
	conn *EuphoriaDBConnection
	tx   *transactions.Transaction
}

func (t *Tx) Commit() error {
	err := t.tx.Commit()
	t.conn.transaction = nil
	return err
}

func (t *Tx) Rollback() error {
	err := t.tx.Rollback()
	t.conn.transaction = nil
	return err
}
