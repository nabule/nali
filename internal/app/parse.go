package app

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"encoding/json"

	"github.com/zu1k/nali/constant"
	"github.com/zu1k/nali/internal/ipdb"
	"github.com/zu1k/nali/internal/tools"
	geoip2 "github.com/zu1k/nali/pkg/geoip"
	"github.com/zu1k/nali/pkg/ipip"
	"github.com/zu1k/nali/pkg/qqwry"
	"github.com/zu1k/nali/pkg/zxipv6wry"
)

var (
	db    []ipdb.IPDB
	qqip  qqwry.QQwry
	geoip geoip2.GeoIP
)

// InitIPDB init ip db content
func InitIPDB(ipdbtype ipdb.IPDBType) {
	db = make([]ipdb.IPDB, 1)
	switch ipdbtype {
	case ipdb.GEOIP2:
		db[0] = geoip2.NewGeoIP(filepath.Join(constant.HomePath, "GeoLite2-City.mmdb"))
	case ipdb.QQIP:
		db[0] = qqwry.NewQQwry(filepath.Join(constant.HomePath, "qqwry.dat"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.db")))
	case ipdb.IPIP:
		db[0] = ipip.NewIPIPFree(filepath.Join(constant.HomePath, "ipipfree.ipdb"))
		db = append(db, zxipv6wry.NewZXwry(filepath.Join(constant.HomePath, "ipv6wry.db")))
	}
}

// parse several ips
func ParseIPs(ips []string) {
	db0 := db[0]
	var db1 ipdb.IPDB
	if len(db) > 1 {
		db1 = db[1]
	} else {
		db1 = db[0]
	}
	for _, ip := range ips {
		v := tools.ValidIP(ip)
		switch v {
		case tools.ValidIPv4:
			result := db0.Find(ip)
			fmt.Println(formatResult(ip, result))
		case tools.ValidIPv6:
			if db1 != nil {
				result := db1.Find(ip)
				fmt.Println(formatResult(ip, result))
			}
		default:
			fmt.Println(ReplaceIPInString(ReplacePhoneInString(ip)))
		}
	}
}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func ReplaceIPInString(str string) (result string) {
	db0 := db[0]
	var db1 ipdb.IPDB
	if len(db) > 1 {
		db1 = db[1]
	} else {
		db1 = db[0]
	}

	result = str
	ip4s := tools.GetIP4FromString(str)
	ip4s = RemoveRepeatedElement(ip4s)
	for _, ip := range ip4s {
		info := db0.Find(ip)
		result = tools.AddInfoIp4(result, ip, info)
	}

	ip6s := tools.GetIP6FromString(str)
	ip6s = RemoveRepeatedElement(ip6s)
	for _, ip := range ip6s {
		info := db1.Find(ip)
		result = tools.AddInfoIp6(result, ip, info)
	}
	return
}

func ReplacePhoneInString(str string) (result string) {
	// db0 := db[0]
	// var db1 ipdbk.IPDB
	// if len(db) > 1 {
	// 	db1 = db[1]
	// } else {
	// 	db1 = db[0]
	// }
	appkey := "京东万象的appkey"

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //忽略对服务器的认证
			},
		},
	}

	result = str
	phones := tools.GetPhoneFromString(str)
	phones = RemoveRepeatedElement(phones)
	for _, phone := range phones {
		resp, err := client.Get("https://way.jd.com/jisuapi/query4?shouji=" + phone + "&appkey=" + appkey)
		if err != nil {
			fmt.Println("Request failed:", err)
			return
		}
		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Respose failed:", err)
			return
		}
		defer resp.Body.Close()

		htmlbodymap := make(map[string]interface{})
		err = json.Unmarshal([]byte(htmlData), &htmlbodymap)
		if err != nil {
			fmt.Println("Umarshal failed:", err)
			return
		}

		PhoneInfoResult := htmlbodymap["result"].(map[string]interface{})["result"]
		PhoneInfoResultMap := PhoneInfoResult.(map[string]interface{})

		info := PhoneInfoResultMap["province"].(string) + PhoneInfoResultMap["city"].(string) + PhoneInfoResultMap["company"].(string)
		// fmt.Printf("%T", PhoneInfoResultMap["city"])
		result = tools.AddInfoPhone(result, phone, info)
	}

	return
}

func formatResult(ip string, result string) string {
	if result == "" {
		result = "未找到"
	}
	return fmt.Sprintf("%s [%s]", ip, result)
}
