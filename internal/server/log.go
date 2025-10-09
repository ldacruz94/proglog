package server 

import (
	"fmt"
	"sync"
)

type Log struct {
	my sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{}
}

func (c *Log) Append(record Record) (uint64, error) {
	c.my.Lock()
	defer c.my.Unlock()
	record.Offset = uint64(len(c.records))
	c.records = append(c.records, record)
	
	return uint64(len(c.records) - 1), nil
}

func (c *Log) Read(offset uint64) (Record, error) {
	c.my.Lock()
	defer c.my.Unlock()
	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return c.records[offset], nil
}

type Record struct {
	Offset uint64 `json:"offset"`
	Value  string `json:"value"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")