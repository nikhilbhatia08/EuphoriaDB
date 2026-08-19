package record

import "fmt"

type Id struct {
	blockNumber int
	slot        int
}

func NewID(blockNumber, slot int) *Id {
	return &Id{blockNumber, slot}
}

func (id *Id) BlockNumber() int {
	return id.blockNumber
}

func (id *Id) Slot() int {
	return id.slot
}

func (id *Id) Equals(other *Id) bool {
	return id.blockNumber == other.blockNumber && id.slot == other.slot
}

func (id *Id) String() string {
	return fmt.Sprintf("[%d, %d]", id.blockNumber, id.slot)
}
