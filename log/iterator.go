package log

import (
	"fmt"
	"runtime"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

type Iterator struct {
	fileManager     *filemgr.FileManager
	block           *filemgr.BlockId
	page            *filemgr.Page
	currentPosition int
	boundary        int
}

func NewIterator(fileManager *filemgr.FileManager, block *filemgr.BlockId) (*Iterator, error) {
	page := filemgr.NewPage(fileManager.BlockSize())

	iterator := &Iterator{
		fileManager: fileManager,
		block:       block,
		page:        page,
	}
	if err := iterator.move(block); err != nil {
		return nil, fmt.Errorf("error moving to block. err: %w", err)
	}

	return iterator, nil
}

func (it *Iterator) HasNext() bool {
	return it.currentPosition < it.fileManager.BlockSize() || it.block.BlockNum() > 0
}

func (it *Iterator) Next() ([]byte, error) {
	if it.currentPosition == it.fileManager.BlockSize() {
		if it.block.BlockNum() == 0 {
			return nil, fmt.Errorf("no more records to read")
		}

		it.block = &filemgr.BlockId{File: it.block.Filename(), Blknum: it.block.BlockNum() - 1}
		if err := it.move(it.block); err != nil {
			return nil, fmt.Errorf("error moving to block. err: %w", err)
		}
	}

	record := it.page.GetBytes(it.currentPosition)
	it.currentPosition += getIntSize() + len(record)
	return record, nil
}

func (it *Iterator) move(block *filemgr.BlockId) error {
	if err := it.fileManager.Read(block, it.page); err != nil {
		return fmt.Errorf("error reading block into page. err: %w", err)
	}

	it.boundary = it.page.GetInt(0)
	it.currentPosition = it.boundary
	return nil
}

func getIntSize() int {
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		return 4
	}

	return 8
}
