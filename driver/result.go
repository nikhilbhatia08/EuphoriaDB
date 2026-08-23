package driver

import "errors"

type Result struct {
	rowsAffected int64
}

func (r *Result) LastInsertId() (int64, error) {
	return 0, errors.New("LastInsertId not supported")
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}
