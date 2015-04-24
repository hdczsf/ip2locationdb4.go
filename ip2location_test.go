package ip2locationdb4

import (
	"testing"
)

func Test_GetRe(t *testing.T) {
	loc, err := Open("F:/迅雷下载/DB4-IP-COUNTRY-REGION-CITY-ISP.BIN/IP-COUNTRY-REGION-CITY-ISP.BIN")
	if err != nil {
		t.Error(err)
	}
	record, err := loc.GetRecord("42.48.85.114")
	if err != nil {
		t.Error(err)
	}
	if record.City!="Changsha"{
		t.Error("err city=",record.City)
	}
}
