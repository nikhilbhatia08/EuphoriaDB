package buffer

import (
	"fmt"
	"sync"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type BufferManager struct {
	BufferPool []*Buffer
	AvailableBuffers int32

	mu sync.Mutex
}

func NewBufferManager(fileManager *filemgr.FileManager, logManager *log.LogManager, numberOfBuffers int32) *BufferManager {
	buffers := make([]*Buffer, numberOfBuffers)
	for i := 0; i < int(numberOfBuffers); i++ {
		buffers[i] = NewBuffer(fileManager, logManager)
	}

	return &BufferManager{
		BufferPool: buffers,
		AvailableBuffers: numberOfBuffers,
	}
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
