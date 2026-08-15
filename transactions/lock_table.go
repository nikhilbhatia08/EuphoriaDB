package transactions

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

const maxWaitTime = 10 * time.Second

type LockTable struct {
	Locks map[filemgr.BlockId]int

	mu   sync.Mutex
	cond *sync.Cond
}

func NewLockTable() *LockTable {
	lockTable := &LockTable{
		Locks: map[filemgr.BlockId]int{},
	}
	lockTable.cond = sync.NewCond(&lockTable.mu)

	return lockTable
}

func (lt *LockTable) SLock(block *filemgr.BlockId) error {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), maxWaitTime)
	defer cancel()

	stopFunc := context.AfterFunc(contextWithTimeout, func() {
		lt.cond.L.Lock()
		lt.cond.Broadcast()
		lt.cond.L.Unlock()
	})
	defer stopFunc()

	for {
		if !lt.hasXLock(block) {
			val := lt.getLockVal(block)
			lt.Locks[*block] = val + 1

			return nil
		}

		lt.cond.Wait()

		if contextWithTimeout.Err() != nil {
			if errors.Is(contextWithTimeout.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("could not acquire shared lock on block %v. err: %v", block, contextWithTimeout.Err())
			}
		}
	}
}

func (lt *LockTable) XLock(block *filemgr.BlockId) error {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), maxWaitTime)
	defer cancel()

	stopFunc := context.AfterFunc(contextWithTimeout, func() {
		lt.cond.L.Lock()
		lt.cond.Broadcast()
		lt.cond.L.Unlock()
	})
	defer stopFunc()

	for {
		if !lt.hasOtherSLocks(block) {
			lt.Locks[*block] = -1
			return nil
		}

		lt.cond.Wait()

		if contextWithTimeout.Err() != nil {
			if errors.Is(contextWithTimeout.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("could not acquire shared lock on block %v. err: %v", block, contextWithTimeout.Err())
			}
		}
	}
}

func (lt *LockTable) Unlock(block *filemgr.BlockId) {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	val := lt.getLockVal(block)
	if val > 1 {
		lt.Locks[*block] = val - 1
	} else {
		delete(lt.Locks, *block)
		lt.cond.Broadcast()
	}
}

func (lt *LockTable) hasXLock(block *filemgr.BlockId) bool {
	return lt.getLockVal(block) < 0
}

func (lt *LockTable) hasOtherSLocks(block *filemgr.BlockId) bool {
	return lt.getLockVal(block) > 1
}

func (lt *LockTable) getLockVal(block *filemgr.BlockId) int {
	val, ok := lt.Locks[*block]
	if !ok {
		return 0
	}

	return val
}
