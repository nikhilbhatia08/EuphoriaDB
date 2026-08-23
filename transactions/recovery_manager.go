package transactions

import (
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/buffer"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type RecoveryMgr struct {
	LogMgr      *log.LogManager
	BufferMgr   *buffer.BufferManager
	Transaction *Transaction
	TxNum       int
}

func NewRecoveryMgr(logMgr *log.LogManager, bufferMgr *buffer.BufferManager, transaction *Transaction, txNum int) *RecoveryMgr {
	return &RecoveryMgr{
		LogMgr:      logMgr,
		BufferMgr:   bufferMgr,
		Transaction: transaction,
		TxNum:       txNum,
	}
}

func (rm *RecoveryMgr) Commit() error {
	lsn, err := WriteCommitRecordToLog(rm.LogMgr, rm.TxNum)
	if err != nil {
		return err
	}

	if err := rm.LogMgr.Flush(lsn); err != nil {
		return err
	}

	if err := rm.BufferMgr.FlushAll(rm.TxNum); err != nil {
		return err
	}

	return nil
}

func (rm *RecoveryMgr) Rollback() error {
	if err := rm.doRollback(); err != nil {
		return err
	}

	lsn, err := WriteRollbackRecordToLog(rm.LogMgr, rm.TxNum)
	if err != nil {
		return err
	}

	if err := rm.LogMgr.Flush(lsn); err != nil {
		return err
	}

	if err := rm.BufferMgr.FlushAll(rm.TxNum); err != nil {
		return err
	}

	return nil
}

func (rm *RecoveryMgr) Recover() error {
	if err := rm.doRecover(); err != nil {
		return err
	}

	if err := rm.BufferMgr.FlushAll(rm.TxNum); err != nil {
		return err
	}

	lsn, err := WriteCheckpointRecordToLog(rm.LogMgr)
	if err != nil {
		return err
	}

	return rm.LogMgr.Flush(lsn)
}

func (rm *RecoveryMgr) doRollback() error {
	iterator, err := rm.LogMgr.Iterator()
	if err != nil {
		return err
	}

	for iterator.HasNext() {
		bytes, err := iterator.Next()
		if err != nil {
			return err
		}

		record, err := CreateLogRecord(bytes)
		if err != nil {
			return err
		}

		if record.TxNumber() == rm.TxNum {
			if record.Op() == Start {
				break
			}
			if err := record.Undo(rm.Transaction); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rm *RecoveryMgr) doRecover() error {
	transactions := make([]int, 0, 10)
	iterator, err := rm.LogMgr.Iterator()
	if err != nil {
		return err
	}

	for iterator.HasNext() {
		bytes, err := iterator.Next()
		if err != nil {
			return err
		}

		record, err := CreateLogRecord(bytes)
		if err != nil {
			return err
		}

		if record.Op() == Checkpoint {
			return nil
		}

		if record.Op() == Commit || record.Op() == Rollback {
			transactions = append(transactions, record.TxNumber())
		} else if !contains(transactions, record.TxNumber()) {
			if err := record.Undo(rm.Transaction); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rm *RecoveryMgr) SetInt(buffer *buffer.Buffer, offset int, newValue int) (int, error) {
	oldval := buffer.Contents.GetInt(offset)
	block := buffer.Block
	return WriteSetIntToLog(rm.LogMgr, rm.TxNum, block, offset, oldval)
}

func (rm *RecoveryMgr) SetString(buffer *buffer.Buffer, offset int, newValue string) (int, error) {
	oldvalue, err := buffer.Contents.GetString(offset)
	if err != nil {
		return -1, err
	}

	block := buffer.Block
	return WriteSetStringToLog(rm.LogMgr, rm.TxNum, block, offset, oldvalue)
}

func (rm *RecoveryMgr) SetBool(buffer *buffer.Buffer, offset int, newValue bool) (int, error) {
	oldVal := buffer.Contents.GetBool(offset)
	block := buffer.Block
	return WriteSetBoolToLog(rm.LogMgr, rm.TxNum, block, offset, oldVal)
}

func (rm *RecoveryMgr) SetDate(buffer *buffer.Buffer, offset int, newVal time.Time) (int, error) {
	oldVal := buffer.Contents.GetDate(offset)
	block := buffer.Block
	return WriteSetDateToLog(rm.LogMgr, rm.TxNum, block, offset, oldVal)
}

func contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
