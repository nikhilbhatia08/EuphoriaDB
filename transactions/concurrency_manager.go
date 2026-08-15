package transactions

import (
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

type ConcurrencyMgr struct {
	lockTable *LockTable
	locks     map[filemgr.BlockId]string
}

func NewConcurrencyMgr(lockTable *LockTable) *ConcurrencyMgr {
	return &ConcurrencyMgr{
		lockTable: lockTable,
		locks:     map[filemgr.BlockId]string{},
	}
}

func (cm *ConcurrencyMgr) SLock(block *filemgr.BlockId) error {
	if _, ok := cm.locks[*block]; !ok {
		if err := cm.lockTable.SLock(block); err != nil {
			return err
		}

		cm.locks[*block] = "s"
	}

	return nil
}

func (cm *ConcurrencyMgr) XLock(block *filemgr.BlockId) error {
	if !cm.hasXLock(block) {
		if err := cm.SLock(block); err != nil {
			return err
		}
		if err := cm.lockTable.XLock(block); err != nil {
			return err
		}

		cm.locks[*block] = "x"
	}

	return nil
}

func (cm *ConcurrencyMgr) Release() {
	for block := range cm.locks {
		cm.lockTable.Unlock(&block)
	}

	cm.locks = map[filemgr.BlockId]string{}
}

func (cm *ConcurrencyMgr) hasXLock(block *filemgr.BlockId) bool {
	lock, ok := cm.locks[*block]
	if ok && lock == "x" {
		return true
	}

	return false
}
