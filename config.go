package ip2region

import (
	"github.com/edunx/lua"
	"os"
)

const (
	INDEX_BLOCK_LENGTH  = 12
	TOTAL_HEADER_LENGTH = 8192
)

type IpInfo struct {
	CityId   int64
	Country  string
	Region   string
	Province string
	City     string
	ISP      string
}

type Ip2Region struct {
	lua.Super

	// db file handler
	dbFileHandler *os.File

	//header block info

	headerSip     []int64
	headerPtr     []int64
	headerLen     int64

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr      []byte
	dbFile        string
}
