package transactions

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/buffer"
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type Transaction struct {
	recoveryMgr *RecoveryMgr
	concurrencyMgr *ConcurrencyMgr
	bufferManager  *buffer.BufferManager
	fileMgr        *filemgr.FileManager
	buffers        *BufferList
	txNum          int
}

const END_OF_FILE = -1

var (
	nextTxNum   = 0
	nextTxNumMu sync.Mutex
)

func nextTxNumber() int {
	nextTxNumMu.Lock()
	defer nextTxNumMu.Unlock()
	nextTxNum++
	return nextTxNum
}

func NewTransaction(bufferManager *buffer.BufferManager, logMgr *log.LogManager, filemgr *filemgr.FileManager, lockTable *LockTable) *Transaction {
	tx := &Transaction{
		concurrencyMgr: NewConcurrencyMgr(lockTable),
		bufferManager:  bufferManager,
		fileMgr:        filemgr,
		txNum:          nextTxNumber(),
		buffers:        NewBufferList(bufferManager),
	}
	tx.recoveryMgr = NewRecoveryMgr(logMgr, bufferManager, tx, tx.txNum)

	return tx
}

func (tx *Transaction) Commit() error {
	if err := tx.recoveryMgr.Commit(); err != nil {
		return err
	}

	tx.concurrencyMgr.Release()
	tx.buffers.UnpinAll()
	fmt.Printf("Transaction %d committed", tx.txNum)

	return nil
}

func (tx *Transaction) Rollback() error {
	if err := tx.recoveryMgr.Rollback(); err != nil {
		return err
	}

	tx.concurrencyMgr.Release()
	tx.buffers.UnpinAll()

	return nil
}

func (tx *Transaction) Recover() error {
	tx.bufferManager.FlushAll(tx.txNum)
	if err := tx.recoveryMgr.Recover(); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) Pin(block *filemgr.BlockId) error {
	if err := tx.buffers.Pin(block); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) Unpin(block *filemgr.BlockId) {
	tx.buffers.Unpin(block)
}

func (tx *Transaction) GetInt(block *filemgr.BlockId, offset int) (int, error) {
	if err := tx.concurrencyMgr.SLock(block); err != nil {
		return math.MinInt, err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return math.MinInt, fmt.Errorf("buffer for block %s not found", block)
	}
	value := buffer.Contents.GetInt(offset)

	return value, nil
}

func (tx *Transaction) GetString(block *filemgr.BlockId, offset int) (string, error) {
	if err := tx.concurrencyMgr.SLock(block); err != nil {
		return "", err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return "", fmt.Errorf("buffer for block %s not found", block)
	}

	return buffer.Contents.GetString(offset)
}

func (tx *Transaction) GetBool(block *filemgr.BlockId, offset int) (bool, error) {
	if err := tx.concurrencyMgr.SLock(block); err != nil {
		return false, err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return false, fmt.Errorf("buffer for block %s not found", block)
	}

	value := buffer.Contents.GetBool(offset)
	return value, nil
}

func (tx *Transaction) GetDate(block *filemgr.BlockId, offset int) (time.Time, error) {
	if err := tx.concurrencyMgr.SLock(block); err != nil {
		return time.Time{}, err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return time.Time{}, fmt.Errorf("buffer for block %s not found", block)
	}

	value := buffer.Contents.GetDate(offset)
	return value, nil
}

func (tx *Transaction) SetInt(block *filemgr.BlockId, offset int, value int, okToLog bool) error {
	var err error
	if err := tx.concurrencyMgr.XLock(block); err != nil {
		return err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return fmt.Errorf("buffer for block %s not found", block)
	}

	lsn := -1
	if okToLog {
		if lsn, err = tx.recoveryMgr.SetInt(buffer, offset, value); err != nil {
			return err
		}
	}

	page := buffer.Contents
	page.SetInt(offset, value)
	buffer.SetModified(tx.txNum, lsn)

	return nil
}

func (tx *Transaction) SetBool(block *filemgr.BlockId, offset int, value bool, okToLog bool) error {
	var err error
	if err := tx.concurrencyMgr.XLock(block); err != nil {
		return err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return fmt.Errorf("buffer for block %s not found", block)
	}

	lsn := -1
	if okToLog {
		if lsn, err = tx.recoveryMgr.SetBool(buffer, offset, value); err != nil {
			return err
		}
	}

	page := buffer.Contents
	page.SetBool(offset, value)
	buffer.SetModified(tx.txNum, lsn)

	return nil
}

func (tx *Transaction) SetDate(block *filemgr.BlockId, offset int, value time.Time, okToLog bool) error {
	var err error
	if err := tx.concurrencyMgr.XLock(block); err != nil {
		return err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return fmt.Errorf("buffer for block %s not found", block)
	}

	lsn := -1
	if okToLog {
		if lsn, err = tx.recoveryMgr.SetDate(buffer, offset, value); err != nil {
			return err
		}
	}

	page := buffer.Contents
	page.SetDate(offset, value)
	buffer.SetModified(tx.txNum, lsn)

	return nil
}

func (tx *Transaction) SetString(block *filemgr.BlockId, offset int, value string, okToLog bool) error {
	var err error
	if err := tx.concurrencyMgr.XLock(block); err != nil {
		return err
	}

	buffer := tx.buffers.GetBuffer(block)
	if buffer == nil {
		return fmt.Errorf("buffer for block %s not found", block)
	}

	lsn := -1
	if okToLog {
		if lsn, err = tx.recoveryMgr.SetString(buffer, offset, value); err != nil {
			return err
		}
	}

	page := buffer.Contents
	page.SetString(offset, value)
	buffer.SetModified(tx.txNum, lsn)

	return nil
}

func (tx *Transaction) Size(filename string) (int, error) {
	dummyBlock := filemgr.NewBlockId(filename, END_OF_FILE)
	if err := tx.concurrencyMgr.SLock(dummyBlock); err != nil {
		return -1, err
	}

	return tx.fileMgr.Length(filename)
}

func (tx *Transaction) Append(filename string) (*filemgr.BlockId, error) {
	dummyBlock := filemgr.NewBlockId(filename, END_OF_FILE)
	if err := tx.concurrencyMgr.XLock(dummyBlock); err != nil {
		return nil, err
	}

	return tx.fileMgr.Append(filename)
}

func (tx *Transaction) BlockSize() int {
	return tx.fileMgr.BlockSize()
}

func (tx *Transaction) AvailableBuffers() int {
	return tx.bufferManager.AvailableBuffers
}
