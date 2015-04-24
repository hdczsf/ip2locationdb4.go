//从官方提供的c语言版本翻译过来的，只支持DB4：https://www.ip2location.com/databases/db4-ip-country-region-city-isp
//别问我为什么，因为我只买了这个库
package ip2locationdb4

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	MAX_IPV4_RANGE = 4294967295
	INVALID_IPV4_ADDRESS = "INVALID IPV4 ADDRESS"
)

var COUNTRY_POSITION = []uint32{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var REGION_POSITION = []uint32{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
var CITY_POSITION = []uint32{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var ISP_POSITION = []uint32{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}

type IP2Location struct {
	data        []byte
	databasetype      uint32
	databasecolumn    uint32
	databaseday       uint32
	databasemonth     uint32
	databaseyear      uint32
	databasecount     uint32
	databaseaddr      uint32
}

type IP2LocationRecord struct {
	Country_Short      string
	Country_Long       string
	Region             string
	City               string
	Isp                string
}

func Open(path string) (*IP2Location, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	loc := &IP2Location{data: data}
	loc.initialize()
	return loc, nil
}

func (this *IP2Location) initialize() {
	this.databasetype = this.read8(1)
	this.databasecolumn = this.read8(2)
	this.databaseyear = this.read8(3)
	this.databasemonth = this.read8(4)
	this.databaseday = this.read8(5)
	this.databasecount = this.read32(6)
	this.databaseaddr = this.read32(10)
}

func (this *IP2Location) read8(position uint32) uint32 {
	return uint32(this.data[position-1])
}

func (this *IP2Location) read32(position uint32) uint32 {
	byte1 := uint32(this.data[position-1])
	byte2 := uint32(this.data[position])
	byte3 := uint32(this.data[position+1])
	byte4 := uint32(this.data[position+2])
	return (byte4 << 24) + (byte3 << 16) + (byte2 << 8) + (byte1)
}

func (this *IP2Location) readStr(position uint32) string {
	size := uint32(this.data[position])
	return string(this.data[position+1 : position+1+size])
}

func (this *IP2Location) read_record(rowaddr uint32) *IP2LocationRecord {
	dbtype := this.databasetype
	record := &IP2LocationRecord{}
	record.Country_Short = this.readStr(this.read32(rowaddr + 4*(COUNTRY_POSITION[dbtype]-1)))
	record.Country_Long = this.readStr(this.read32(rowaddr+4*(COUNTRY_POSITION[dbtype]-1)) + 3)
	record.Region = this.readStr(this.read32(rowaddr + 4*(REGION_POSITION[dbtype]-1)))
	record.City = this.readStr(this.read32(rowaddr + 4*(CITY_POSITION[dbtype]-1)))
	record.Isp = this.readStr(this.read32(rowaddr + 4*(ISP_POSITION[dbtype]-1)))
	return record
}
func (this *IP2Location) GetRecord(ipstring string) (*IP2LocationRecord, error) {
	baseaddr := this.databaseaddr
	dbcolumn := this.databasecolumn

	var low uint32 = 0
	var high uint32 = this.databasecount
	var mid uint32 = 0

	var ipno uint32
	var ipfrom uint32
	var ipto uint32

	ipno, err := ip2no(ipstring)
	if err != nil {
		return nil, err
	}
	if ipno == MAX_IPV4_RANGE {
		ipno = ipno - 1
	}

	for low <= high {
		mid = uint32((low + high) / 2)
		ipfrom = this.read32(baseaddr + mid*dbcolumn*4)
		ipto = this.read32(baseaddr + (mid+1)*dbcolumn*4)

		if (ipno >= ipfrom) && (ipno < ipto) {
			return this.read_record(baseaddr + (mid * dbcolumn * 4)), nil
		} else {
			if ipno < ipfrom {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return nil, errors.New("not find")
}

func ip2no(ipstring string) (uint32, error) {
	var ipArr = strings.Split(ipstring, ".")
	if len(ipArr) != 4 {
		return 0, errors.New(INVALID_IPV4_ADDRESS)
	}
	byte1, err := strconv.ParseUint(ipArr[0], 10, 32)
	if err != nil {
		return 0, err
	}
	byte2, err := strconv.ParseUint(ipArr[1], 10, 32)
	if err != nil {
		return 0, err
	}
	byte3, err := strconv.ParseUint(ipArr[2], 10, 32)
	if err != nil {
		return 0, err
	}
	byte4, err := strconv.ParseUint(ipArr[3], 10, 32)
	if err != nil {
		return 0, err
	}
	var a uint32 = 0

	a = uint32(byte4)
	a += uint32(byte3) * 256
	a += uint32(byte2) * 256 * 256
	a += uint32(byte1) * 256 * 256 * 256
	return a, nil
}
