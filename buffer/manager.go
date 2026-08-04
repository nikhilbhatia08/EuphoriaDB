package buffer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

const maxWaitTime = 10 * time.Second

type BufferManager struct {
	BufferPool []*Buffer
	AvailableBuffers int32

	mu sync.Mutex
	cond *sync.Cond
}

func NewBufferManager(fileManager *filemgr.FileManager, logManager *log.LogManager, numberOfBuffers int32) *BufferManager {
	buffers := make([]*Buffer, numberOfBuffers)
	for i := 0; i < int(numberOfBuffers); i++ {
		buffers[i] = NewBuffer(fileManager, logManager)
	}

	bufferManager := &BufferManager{
		BufferPool: buffers,
		AvailableBuffers: numberOfBuffers,
	}
	bufferManager.cond = sync.NewCond(&bufferManager.mu)
	return bufferManager
}

func (bm *BufferManager) FlushAll(txNum int32) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for _, buffer := range bm.BufferPool {
		if buffer.TxNum != txNum {
			continue
		}

		if err := buffer.flush(); err != nil {
			return fmt.Errorf("error flushing buffer. err: %w", err)
		}
	}

	return nil
}

func (bm *BufferManager) Unpin(buffer *Buffer) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	buffer.Unpin()
	if !buffer.IsPinned() {
		bm.AvailableBuffers++
		bm.cond.Broadcast()
	}
}

func (bm *BufferManager) Pin(block *filemgr.BlockId) (*Buffer, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitTime)
	defer cancel()

	stop := context.AfterFunc(ctx, func() {
		bm.cond.L.Lock()
		bm.cond.Broadcast()
		bm.cond.L.Unlock()
	})
	defer stop()

	for {
		buffer, err := bm.tryToPin(block)
		if err != nil {
			return nil, fmt.Errorf("error pinning block. err: %w", err)
		} else if buffer != nil {
			return buffer, nil
		}

		bm.cond.Wait()

		if ctx.Err() != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, fmt.Errorf("buffer aborted while waiting for block %s. err: %w", block.ToString(), ctx.Err())
			}

			return nil, ctx.Err()
		}
	}
}

func (bm *BufferManager) tryToPin(blockId *filemgr.BlockId) (*Buffer, error) {
	buff := bm.findExistingBuffer(blockId)
	if buff == nil {
		buff = bm.chooseUnpinnedBuffer()
		if buff == nil {
			return nil, nil
		}

		if err := buff.AssignToBlock(blockId); err != nil {
			return nil, fmt.Errorf("unable to pin block. err: %w", err)
		}
	}

	if !buff.IsPinned() {
		bm.AvailableBuffers--
	}
	buff.Pin()
	return buff, nil
}

func (bm *BufferManager) findExistingBuffer(BlockId *filemgr.BlockId) *Buffer {
	for _, buffer := range bm.BufferPool {
		if buffer.Block != nil && buffer.Block.Equals(BlockId) {
			return buffer
		}
	}

	return nil
}

func (bm *BufferManager) chooseUnpinnedBuffer() *Buffer {
	for _, buffer := range bm.BufferPool {
		if !buffer.IsPinned() {
			return buffer
		}
	}

	return nil
}
