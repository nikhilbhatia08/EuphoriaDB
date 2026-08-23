package server

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/buffer"
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
	"github.com/nikhilbhatia08/EuphoriaDB/metadata"
	"github.com/nikhilbhatia08/EuphoriaDB/plan_impl"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

const (
	blockSize  = 2048
	bufferSize = 128
	logFile    = "euphoriadb.log"
)

type EuphoriaDB struct {
	fileManager     *filemgr.FileManager
	bufferManager   *buffer.BufferManager
	logManager      *log.LogManager
	metadataManager *metadata.MetadataManager
	lockTable       *transactions.LockTable
	queryPlanner    plan_impl.QueryPlanner
	updatePlanner   plan_impl.UpdatePlanner
	planner         *plan_impl.Planner
}

func NewEuphoriaDB(directoryName string) (*EuphoriaDB, error) {
	db := &EuphoriaDB{}

	var err error
	if db.fileManager, err = filemgr.NewFileManager(directoryName, blockSize); err != nil {
		return nil, err
	}
	if db.logManager, err = log.NewLogManager(db.fileManager, logFile); err != nil {
		return nil, err
	}

	db.bufferManager = buffer.NewBufferManager(db.fileManager, db.logManager, bufferSize)
	db.lockTable = transactions.NewLockTable()

	isNew := db.fileManager.IsNew()
	transaction := db.NewTransaction()

	if isNew {
		fmt.Println("Creating new database")
	} else {
		fmt.Println("recovering database")
		if err := transaction.Recover(); err != nil {
			return nil, err
		}
	}

	if db.metadataManager, err = metadata.NewMetadataManager(isNew, transaction); err != nil {
		return nil, err
	}

	if err = transaction.Commit(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *EuphoriaDB) NewTransaction() *transactions.Transaction {
	return transactions.NewTransaction(db.bufferManager, db.logManager, db.fileManager, db.lockTable)
}

func (db *EuphoriaDB) MetadataManager() *metadata.MetadataManager {
	return db.metadataManager
}

func (db *EuphoriaDB) Planner() *plan_impl.Planner {
	return db.planner
}

func (db *EuphoriaDB) FileManager() *filemgr.FileManager {
	return db.fileManager
}

func (db *EuphoriaDB) LogManager() *log.LogManager {
	return db.logManager
}

func (db *EuphoriaDB) BufferManager() *buffer.BufferManager {
	return db.bufferManager
}
