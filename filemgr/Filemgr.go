package filemgr

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileManager struct {
	dbDirectory   string
	blockSize     int
	isNew         bool
	mu            sync.Mutex
	openFiles     map[string]*os.File
	blocksRead    int
	blocksWritten int
}

func NewFileManager(dbDirectory string, blockSize int) (*FileManager, error) {
	isNew := false

	if _, err := os.Stat(dbDirectory); os.IsNotExist(err) {
		isNew = true
		if err := os.MkdirAll(dbDirectory, 0755); err != nil {
			return nil, fmt.Errorf("cannot create directory %s:%v", dbDirectory, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot access directory %s:%v", dbDirectory, err)
	}

	entries, err := os.ReadDir(dbDirectory)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory %s:%v", dbDirectory, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(name) >= 4 && name[:4] == "temp" {
				tempFilePath := filepath.Join(dbDirectory, name)
				if err := os.Remove(tempFilePath); err != nil {
					return nil, fmt.Errorf("cannot remove file %s:%v", dbDirectory, err)
				}
			}
		}
	}

	return &FileManager{
		dbDirectory:   dbDirectory,
		blockSize:     blockSize,
		isNew:         isNew,
		openFiles:     make(map[string]*os.File),
		blocksRead:    0,
		blocksWritten: 0,
	}, nil
}

func (fm *FileManager) Read(blk *BlockId, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	f, err := fm.getFile(blk.Filename())
	if err != nil {
		return fmt.Errorf("cannot read block %s: %v", blk.ToString(), err)
	}

	offset := int64(blk.BlockNum()) * int64(fm.blockSize)
	if _, err := f.Seek(offset, io.SeekCurrent); err != nil {
		return fmt.Errorf("cannot seek to offset %d: %v", offset, err)
	}

	buf := p.Contents()
	n, err := io.ReadFull(f, buf)

	if err == nil && n == len(buf) {
		fm.blocksRead++
		return nil
	}

	if errors.Is(err, io.EOF) {
		if n == 0 {
			fm.blocksRead++
			return nil
		}

		return fmt.Errorf("partial read at EOF: expected %d bytes, got %d", len(buf), n)
	}

	if err != nil {
		return fmt.Errorf("cannot read data %v", err)
	}

	return fmt.Errorf("short read: expected %d bytes, got %d", len(buf), n)
}

func (fm *FileManager) Write(blk *BlockId, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	f, err := fm.getFile(blk.Filename())
	if err != nil {
		return fmt.Errorf("cannot read block %s: %v", blk.ToString(), err)
	}

	offset := int64(blk.BlockNum()) * int64(fm.blockSize)
	if _, err := f.Seek(offset, io.SeekCurrent); err != nil {
		return fmt.Errorf("cannot seek to offset %d: %v", offset, err)
	}

	buf := p.Contents()
	n, err := f.Write(buf)

	if err != nil {
		if n != len(buf) {
			return fmt.Errorf("short write: expected %d bytes, wrote %d", len(buf), n)
		}
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("cannot flush file %s, to disk %v", blk.Filename(), err)
	}

	fm.blocksWritten++
	return nil
}

func (fm *FileManager) Append(filename string) (*BlockId, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	newBlockNumber, err := fm.Length(filename)
	if err != nil {
		return &BlockId{}, fmt.Errorf("cannot get length of %s: %v", filename, err)
	}

	block := NewBlockId(filename, newBlockNumber)

	f, err := fm.getFile(filename)
	if err != nil {
		return &BlockId{}, fmt.Errorf("cannot get file %s: %v", filename, err)
	}

	offset := int64(block.BlockNum() * fm.blockSize)
	if _, err := f.Seek(offset, io.SeekCurrent); err != nil {
		return &BlockId{}, fmt.Errorf("cannot seek to offset %d: %v", offset, err)
	}

	b := make([]byte, fm.blockSize)
	n, err := f.Write(b)
	if err != nil {
		return &BlockId{}, fmt.Errorf("cannot write data: %v", err)
	}
	if n != len(b) {
		return nil, fmt.Errorf("short write: expected %d bytes wrote %d", len(b), n)
	}

	if err := f.Sync(); err != nil {
		return nil, fmt.Errorf("cannot sync file %s: %v", filename, err)
	}

	fm.blocksWritten++
	return block, nil
}

func (fm *FileManager) Length(filename string) (int, error) {
	f, err := fm.getFile(filename)
	if err != nil {
		return 0, fmt.Errorf("cannot access %s: %v", filename, err)
	}

	fileinfo, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("cannot stat %s: %v", filename, err)
	}

	fileSizeInBytes := fileinfo.Size()
	return int(fileSizeInBytes / int64(fm.blockSize)), nil
}

func (fm *FileManager) IsNew() bool {
	return fm.isNew
}

func (fm *FileManager) BlockSize() int {
	return fm.blockSize
}

func (fm *FileManager) GetBlocksRead() int {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return fm.blocksRead
}

func (fm *FileManager) GetBlocksWritten() int {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return fm.blocksWritten
}

func (fm *FileManager) getFile(filename string) (*os.File, error) {
	if f, ok := fm.openFiles[filename]; ok {
		return f, nil
	}

	dbTable := filepath.Join(fm.dbDirectory, filename)
	f, err := os.OpenFile(dbTable, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %v", dbTable, err)
	}

	fm.openFiles[dbTable] = f
	return f, nil
}
