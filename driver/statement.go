package driver

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type Statement struct {
	conn  *EuphoriaDBConnection
	query string
}

func (st *Statement) Close() error {
	return nil
}

func (st *Statement) NumInput() int {
	return -1
}

func (st *Statement) Exec(args []driver.Value) (driver.Result, error) {
	var tx *transactions.Transaction
	if st.conn.transaction == nil {
		tx = st.conn.db.NewTransaction()
	} else {
		tx = st.conn.transaction
	}

	planner := st.conn.db.Planner()

	lower := strings.ToLower(strings.TrimSpace(st.query))
	if strings.HasPrefix(lower, "select") {
		return nil, fmt.Errorf("select statement in exec is not allowed")
	}

	rowsAffected, err := planner.ExecuteUpdate(st.query, tx)
	if err != nil {
		if st.conn.transaction == nil {
			_ = tx.Rollback()
		}
		return nil, err
	}

	if st.conn.transaction == nil {
		if err := tx.Commit(); err != nil {
			return nil, err
		}
	}

	return &Result{
		rowsAffected: int64(rowsAffected),
	}, nil
}

func (st *Statement) Query(args []driver.Value) (driver.Rows, error) {
	var tx *transactions.Transaction
	if st.conn.transaction == nil {
		tx = st.conn.db.NewTransaction()
	} else {
		tx = st.conn.transaction
	}

	lower := strings.ToLower(strings.TrimSpace(st.query))
	if !strings.HasPrefix(lower, "select") {
		return nil, fmt.Errorf("query not allowed with non-select statement")
	}

	planner := st.conn.db.Planner()

	plan, err := planner.CreateQueryPlanner(st.query, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	sc, err := plan.Open()
	if err != nil {
		if st.conn.transaction == nil {
			_ = tx.Rollback()
		}
		return nil, err
	}

	return &Rows{
		stmt: st,
		tx:   tx,
		scan: sc,
		plan: plan,
	}, nil
}
