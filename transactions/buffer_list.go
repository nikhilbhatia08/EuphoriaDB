package transactions

import (
	"github.com/nikhilbhatia08/EuphoriaDB/buffer"
	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

type BufferList struct {
	Buffers       map[filemgr.BlockId]*buffer.Buffer
	Pins          map[buffer.Buffer]int
	BufferManager *buffer.BufferManager
}

func NewBufferList(bufferManager *buffer.BufferManager) *BufferList {
	return &BufferList{
		Buffers:       map[filemgr.BlockId]*buffer.Buffer{},
		Pins:          map[buffer.Buffer]int{},
		BufferManager: bufferManager,
	}
}

func (bl *BufferList) Pin(block *filemgr.BlockId) error {
	if buffer, ok := bl.Buffers[*block]; ok {
		bl.Pins[*buffer]++
		return nil
	}

	buffer, err := bl.BufferManager.Pin(block)
	if err != nil {
		return err
	}
	bl.Buffers[*block] = buffer
	bl.Pins[*buffer]++

	return nil
}

func (bl *BufferList) Unpin(block *filemgr.BlockId) {
	buffer, ok := bl.Buffers[*block]
	if !ok {
		return
	}

	bl.Pins[*buffer]--
	if bl.Pins[*buffer] <= 0 {
		bl.BufferManager.Unpin(buffer)
		delete(bl.Pins, *buffer)
		delete(bl.Buffers, *block)
	}
}

func (bl *BufferList) UnpinAll() {
	for _, buffer := range bl.Buffers {
		bl.BufferManager.Unpin(buffer)
	}
	bl.Buffers = map[filemgr.BlockId]*buffer.Buffer{}
	bl.Pins = map[buffer.Buffer]int{}
}

func (bl *BufferList) GetBuffer(block *filemgr.BlockId) *buffer.Buffer {
	buffer, ok := bl.Buffers[*block]
	if !ok {
		return nil
	}

	return buffer
}
