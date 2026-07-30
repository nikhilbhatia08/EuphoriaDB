package buffer

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
	"github.com/nikhilbhatia08/EuphoriaDB/log"
)

type Buffer struct {
	FileManager *filemgr.FileManager
	LogManager  *log.LogManager
	Contents    *filemgr.Page
	Block       *filemgr.BlockId
	Pins        int32
	TxNum       int32
	Lsn         int32
}

func NewBuffer(fileManager *filemgr.FileManager, logManager *log.LogManager) *Buffer {
	contents := filemgr.NewPage(fileManager.BlockSize())
	return &Buffer{
		FileManager: fileManager,
		LogManager:  logManager,
		Contents:    contents,
		Pins:        0,
		TxNum:       -1,
		Lsn:         -1,
	}
}

func (b *Buffer) SetModified(txNum int32, lsn int32) {
	b.TxNum = txNum
	b.Lsn = lsn
}

func (b *Buffer) IsPinned() bool {
	return b.Pins > 0
}

func (b *Buffer) AssignToBlock(BlockId *filemgr.BlockId) error {
	if err := b.flush(); err != nil {
		return fmt.Errorf("unable to assign block to buffer. err: %w", err)
	}

	b.Block = BlockId
	if err := b.FileManager.Read(BlockId, b.Contents); err != nil {
		return fmt.Errorf("unable to read contents of block to page. err: %w", err)
	}
	b.Pins = 0

	return nil
}

func (b *Buffer) flush() error {
	if b.TxNum == -1 {
		return nil
	}

	if err := b.LogManager.Flush(b.Lsn); err != nil {
		return fmt.Errorf("unable to flush. err: %w", err)
	}
	if err := b.FileManager.Write(b.Block, b.Contents); err != nil {
		return fmt.Errorf("unable to write buffer contents to disk. err: %w", err)
	}
	b.TxNum = -1

	return nil
}

func (b *Buffer) Pin() {
	b.Pins++
}

func (b *Buffer) Unpin() {
	b.Pins--
}
