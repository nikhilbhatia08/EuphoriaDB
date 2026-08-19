package metadata

type StatInfo struct {
	numBlocks  int
	numRecords int
}

func NewStatInfo(numBlocks, numRecords int) *StatInfo {
	return &StatInfo{
		numBlocks:  numBlocks,
		numRecords: numRecords,
	}
}

func (si *StatInfo) NumBlocks() int {
	return si.numBlocks
}

func (si *StatInfo) NumRecords() int {
	return si.numRecords
}

func (si *StatInfo) DistinctValues(fieldName string) int {
	return 1 + (si.numRecords / 3)
}
