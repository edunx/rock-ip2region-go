package ip2region

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func (ip IpInfo)String() string {
	return strconv.FormatInt(ip.CityId, 10) + "|" + ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func getIpInfo(cityId int64, line []byte) IpInfo {

	lineSlice := strings.Split(string(line), "|")
	ipInfo := IpInfo{}
	length := len(lineSlice)
	ipInfo.CityId = cityId
	if length < 5 {
		for i := 0; i <= 5 - length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

func (this *Ip2Region) Start() error {
	file , err := os.Open(this.dbFile)
	if err != nil {
		return err
	}
	this.dbFileHandler = file

	return nil
}

func (this *Ip2Region) Close() {
	this.dbFileHandler.Close()
}

func (this *Ip2Region) Search(ipStr string) ( int64 , []byte , error) {
	var err error
	if this.totalBlocks == 0 {
		this.dbBinStr, err = os.ReadFile(this.dbFile)
		if err != nil {
			return 0 , nil , err
		}

		this.firstIndexPtr = getLong(this.dbBinStr, 0)
		this.lastIndexPtr = getLong(this.dbBinStr, 4)
		this.totalBlocks = (this.lastIndexPtr - this.firstIndexPtr) / INDEX_BLOCK_LENGTH + 1
	}

	ip, err := ip2long(ipStr)
	if err != nil {
		return 0 , nil , err
	}

	h := this.totalBlocks
	var dataPtr, l int64;
	for (l <= h) {

		m := (l + h) >> 1
		p := this.firstIndexPtr + m * INDEX_BLOCK_LENGTH
		sip := getLong(this.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(this.dbBinStr, p + 4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(this.dbBinStr, p + 8)
				break;
			}
		}
	}
	if dataPtr == 0 {
		return 0 , nil , errors.New("not found")
	}

	dataLen := ((dataPtr >> 24) & 0xFF)
	dataPtr = (dataPtr & 0x00FFFFFF);

	return getLong(this.dbBinStr, dataPtr) , this.dbBinStr[(dataPtr) + 4:dataPtr + dataLen] , nil

}
func getLong(b []byte, offset int64) int64 {

	val := (int64(b[offset ]) |
		int64(b[offset + 1]) << 8 |
		int64(b[offset + 2]) << 16 |
		int64(b[offset + 3]) << 24)

	return val

}

func ip2long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24 - 8 * i)
	}

	return sum, nil
}
