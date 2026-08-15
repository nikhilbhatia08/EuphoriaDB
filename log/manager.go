package log

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/nikhilbhatia08/EuphoriaDB/filemgr"
)

type LogManager struct {
	FileManager    *filemgr.FileManager
	Logfile        string
	LogPage        *filemgr.Page
	CurrentBlock   *filemgr.BlockId
	latestLSN      int
	latestSavedLSN int
	mu             sync.Mutex
}

func NewLogManager(fileManager *filemgr.FileManager, logfile string) (*LogManager, error) {
	if fileManager == nil {
		return nil, fmt.Errorf("error initializing log Manager. FileManager not initialized")
	}

	blockSize := fileManager.BlockSize()
	page := filemgr.NewPage(blockSize)

	logSize, err := fileManager.Length(logfile)
	if err != nil {
		return nil, fmt.Errorf("error fetching length of file: %s", logfile)
	}

	var currentBlock *filemgr.BlockId
	if logSize == 0 {
		currentBlock, err = appendNewBlock(fileManager, logfile, page)
		if err != nil {
			return nil, fmt.Errorf("error appending a new block. err: %w", err)
		}
	} else {
		currentBlock = filemgr.NewBlockId(logfile, logSize-1)
		if err := fileManager.Read(currentBlock, page); err != nil {
			return nil, fmt.Errorf("error reading last block of log file: %s. err: %w", logfile, err)
		}
	}

	return &LogManager{
		FileManager:  fileManager,
		Logfile:      logfile,
		LogPage:      page,
		CurrentBlock: currentBlock,
		latestLSN:    0,
	}, nil
}

func (lm *LogManager) Append(logRecord []byte) (int, error) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	boundary := lm.LogPage.GetInt(0)
	recordSize := len(logRecord)
	var bytesNeeded, intSize int
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		bytesNeeded = recordSize + 4
		intSize = 4
	} else {
		bytesNeeded = recordSize + 8
		intSize = 8
	}

	if boundary-int(bytesNeeded) < intSize {
		if err := lm.flush(); err != nil {
			return 0, fmt.Errorf("error flushing page contents to disk, file: %s, err: %w", lm.Logfile, err)
		}

		var err error
		lm.CurrentBlock, err = appendNewBlock(lm.FileManager, lm.Logfile, lm.LogPage)
		if err != nil {
			return 0, fmt.Errorf("error appending new block. err: %w", err)
		}

		boundary = lm.LogPage.GetInt(0)
	}

	recordPosition := boundary - bytesNeeded
	lm.LogPage.Setbytes(recordPosition, logRecord)
	lm.LogPage.SetInt(0, recordPosition)

	lm.latestLSN += 1
	return lm.latestLSN, nil
}

func (lm *LogManager) Iterator() (*Iterator, error) {
	if err := lm.flush(); err != nil {
		return nil, fmt.Errorf("error flushing log manager before creating iterator. err: %w", err)
	}

	return NewIterator(lm.FileManager, lm.CurrentBlock)
}

func (lm *LogManager) flush() error {
	if err := lm.FileManager.Write(lm.CurrentBlock, lm.LogPage); err != nil {
		return fmt.Errorf("error writing page contents to disk. err: %w", err)
	}

	lm.latestSavedLSN = lm.latestLSN
	return nil
}

func (lm *LogManager) Flush(lsn int) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	if lsn >= lm.latestLSN {
		return lm.flush()
	}
	return nil
}

func appendNewBlock(fileManager *filemgr.FileManager, logfile string, logPage *filemgr.Page) (*filemgr.BlockId, error) {
	blk, err := fileManager.Append(logfile)
	if err != nil {
		return nil, fmt.Errorf("error appending block in file: %s. err: %w", logfile, err)
	}

	logPage.SetInt(0, fileManager.BlockSize())
	if err := fileManager.Write(blk, logPage); err != nil {
		return nil, fmt.Errorf("error writing new block. err: %w", err)
	}
	return blk, nil
}
