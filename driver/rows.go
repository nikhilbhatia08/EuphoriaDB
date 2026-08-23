package driver

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

type Rows struct {
	stmt    *Statement
	tx      *transactions.Transaction
	scan    scan.TableScan
	plan    plan.Plan
	done    bool
	columns []string
}

func (r *Rows) Columns() []string {
	if r.columns == nil {
		sch := r.plan.Schema()
		fields := sch.Fields()
		r.columns = make([]string, len(fields))
		copy(r.columns, fields)
	}
	return r.columns
}

func (r *Rows) Close() error {
	if r.done {
		return nil
	}
	r.done = true
	r.scan.Close()

	return r.tx.Commit()
}

func (r *Rows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}

	next, err := r.scan.Next()
	if err != nil {
		_ = r.tx.Rollback()
		r.done = true
		return err
	}
	if !next {
		r.done = true
		if err := r.tx.Commit(); err != nil {
			return err
		}
		return io.EOF
	}

	columns := r.Columns()
	for i, col := range columns {
		columnType := r.plan.Schema().FieldType(col)

		var v interface{}
		switch columnType {
		case types.Integer:
			v, err = r.scan.GetInt(col)
			if err != nil {
				return err
			}
		case types.Varchar:
			v, err = r.scan.GetString(col)
			if err != nil {
				return err
			}
		case types.Boolean:
			v, err = r.scan.GetBool(col)
			if err != nil {
				return err
			}
		case types.Date:
			v, err = r.scan.GetDate(col)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported field type: %v", columnType)
		}
		dest[i] = v
	}
	return nil
}
