package tools

import (
	"regexp"
	"strings"
)

func AddInfoIp4(origin string, ip string, info string) (result string) {
	re := regexp.MustCompile("(^|[^0-9.])(" + strings.ReplaceAll(ip, ".", "\\.") + ")($|[^0-9.])")
	result = re.ReplaceAllString(origin, "$1$2"+" ["+info+"]$3")
	return strings.TrimRight(result, " \t")
}

func AddInfoPhone(origin string, phone string, info string) (result string) {
	pattern := phone
	re := regexp.MustCompile(pattern)
	result = re.ReplaceAllString(origin, "$0"+" ["+info+"]")
	return strings.TrimRight(result, " \t")
}

func AddInfoIp6(origin string, ip string, info string) (result string) {
	re := regexp.MustCompile("(^|[^0-9a-fA-F:])(" + strings.ReplaceAll(ip, ".", "\\.") + ")($|[^0-9a-fA-F:])")
	result = re.ReplaceAllString(origin, "$1$2"+" ["+info+"]$3")
	return strings.TrimRight(result, " \t")
}

func AddInfoDomain(origin string, domain string, info string) (result string) {
	re := regexp.MustCompile("(^|[^0-9a-zA-Z-])(" + strings.ReplaceAll(domain, ".", "\\.") + ")($|[^0-9a-zA-Z-\\.])")
	result = re.ReplaceAllString(origin, "$1$2"+" ["+info+"]$3")
	return strings.TrimRight(result, " \t")
}
