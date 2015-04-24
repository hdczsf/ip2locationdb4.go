# ip2locationdb4.go
ip2location的go语言版本，从官方提供的c语言版本翻译过来的，只支持DB4：https://www.ip2location.com/databases/db4-ip-country-region-city-isp

别问我为什么，因为我只买了这个库

### 安装
`go get github.com/hdczsf/ip2locationdb4.go`

### 例子

```go
package ip2locationdb4

import (
	"testing"
)

func Test_GetRecord(t *testing.T) {
	loc, err := Open("F:/DB4-IP-COUNTRY-REGION-CITY-ISP.BIN/IP-COUNTRY-REGION-CITY-ISP.BIN")
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
```
