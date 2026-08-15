package filemgr

import (
	"fmt"
	"hash/fnv"
)

type BlockId struct {
	File   string
	Blknum int
}

func NewBlockId(filename string, blknum int) *BlockId {
	return &BlockId{
		File:   filename,
		Blknum: blknum,
	}
}

func (b *BlockId) Filename() string {
	return b.File
}

func (b *BlockId) BlockNum() int {
	return b.Blknum
}

func (b *BlockId) Equals(block *BlockId) bool {
	return b.File == block.File && b.Blknum == block.Blknum
}

func (b *BlockId) String() string {
	return fmt.Sprintf("[ File : %s , block : %d]", b.File, b.Blknum)
}

func (b *BlockId) hash(hashstr string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(hashstr))
	return h.Sum32()
}

func (b *BlockId) HashCode() int {
	return int(b.hash(b.String()))
}
