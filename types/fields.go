package types

import "runtime"

type Type int

const (
	Integer Type = 4
	Varchar Type = 12
	Boolean Type = 16
	Date    Type = 91
)

type FieldInfo struct {
	Type   Type
	Length int
}

func GetIntSize() int {
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		return 4
	}

	return 8
}