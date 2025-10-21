package filemgr

import (
	"fmt"
	"hash/fnv"
)

type BlockId struct {
	filename string
	blknum   int
}

func NewBlockId(filename string, blknum int) *BlockId {
	return &BlockId{
		filename: filename,
		blknum:   blknum,
	}
}

func (b *BlockId) Filename() string {
	return b.filename
}

func (b *BlockId) BlockNum() int {
	return b.blknum
}

func (b *BlockId) Equals(block *BlockId) bool {
	return b.filename == block.filename && b.blknum == block.blknum
}

func (b *BlockId) ToString() string {
	return fmt.Sprintf("[ File : %s , block : %d]", b.filename, b.blknum)
}

func (b *BlockId) hash(hashstr string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(hashstr))
	return h.Sum32()
}

func (b *BlockId) HashCode() int {
	return int(b.hash(b.ToString()))
}
