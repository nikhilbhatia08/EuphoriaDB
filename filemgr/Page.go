package filemgr

import (
	"encoding/binary"
	"errors"
	"runtime"
	"time"
	"unicode/utf8"
)

// We can add multiple datatypes here
// For now we support
// 1. Int
// 2. String
// 3. Short
// 4. Bool
// 5. Date

type Page struct {
	buffer []byte
}

func NewPage(blocksize int) *Page {
	return &Page{
		buffer: make([]byte, blocksize),
	}
}

func NewPageFromBytes(b []byte) *Page {
	return &Page{
		buffer: b,
	}
}

func (p *Page) GetInt(offset int) int {
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		return int(binary.BigEndian.Uint32(p.buffer[offset:]))
	}

	return int(binary.BigEndian.Uint64(p.buffer[offset:]))
}

func (p *Page) SetInt(offset int, n int) {
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		binary.BigEndian.PutUint32(p.buffer[offset:], uint32(n))
	} else {
		binary.BigEndian.PutUint64(p.buffer[offset:], uint64(n))
	}
}

func (p *Page) GetBytes(offset int) []byte {
	len := p.GetInt(offset)
	start := offset
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		start = start + 4
	} else {
		start = start + 8
	}
	end := start + len
	resultbuff := make([]byte, len)
	copy(resultbuff, p.buffer[start:end])
	return resultbuff
}

func (p *Page) Setbytes(offset int, b []byte) {
	len := len(b)
	p.SetInt(offset, len)
	start := offset
	if runtime.GOARCH == "386" || runtime.GOARCH == "arm" {
		start = start + 4
	} else {
		start = start + 8
	}
	copy(p.buffer[start:], b)
}

func (p *Page) GetString(offset int) (string, error) {
	bytes := p.GetBytes(offset)
	if !utf8.Valid(bytes) {
		return "", errors.New("invalid UTF-8 encoding")
	}
	return string(bytes), nil
}

func (p *Page) SetString(offset int, s string) error {
	if !utf8.ValidString(s) {
		return errors.New("invalid string, not UTF-8 encoded")
	}

	p.Setbytes(offset, []byte(s))
	return nil
}

func (p *Page) GetShort(offset int) int16 {
	return int16(binary.BigEndian.Uint16(p.buffer[offset:]))
}

func (p *Page) SetShort(offset int, n int16) {
	binary.BigEndian.PutUint16(p.buffer[offset:], uint16(n))
}

func (p *Page) GetBool(offset int) bool {
	return p.buffer[offset] != 0
}

func (p *Page) SetBool(offset int, b bool) {
	if b {
		p.buffer[offset] = 1
	} else {
		p.buffer[offset] = 0
	}
}

func (p *Page) GetDate(offset int) time.Time {
	unixTimestamp := int64(binary.BigEndian.Uint64(p.buffer[offset:]))
	return time.Unix(unixTimestamp, 0)
}

func (p *Page) SetDate(offset int, date time.Time) {
	binary.BigEndian.PutUint64(p.buffer[offset:], uint64(date.Unix()))
}

func MaxLength(strlen int) int {
	var size int
	if runtime.GOARCH == "arm" || runtime.GOARCH == "386" {
		size = 4
	} else {
		size = 8
	}
	return size + strlen*utf8.UTFMax
}

func (p *Page) Contents() []byte {
	return p.buffer
}
