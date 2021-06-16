package tools

import (
	"net"
	"regexp"
	"strings"
)

var (
	ipv4re  *regexp.Regexp
	phonere *regexp.Regexp

	ipv6re0 *regexp.Regexp
	ipv6re  *regexp.Regexp
)

func init() {
	ipv4re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	phonere = regexp.MustCompile(`\b1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}`)

	ipv6re0 = regexp.MustCompile(`^fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?$`)
	ipv6re = regexp.MustCompile(`fe80:(:[0-9a-fA-F]{1,4}){0,4}(%\w+)?|([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?::(([0-9a-fA-F]{1,4}:){0,6}[0-9a-fA-F]{1,4})?`)
}

func GetIP4FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv4re.FindAllString(str, -1)
}

func GetIP6FromString(str string) []string {
	str = strings.Trim(str, " ")
	return ipv6re.FindAllString(str, -1)
}

func GetPhoneFromString(str string) []string {
	str = strings.Trim(str, " ")
	return phonere.FindAllString(str, -1)
}

const (
	ValidIPv4 = iota
	ValidIPv6
	InvalidIP
)

type Valid int

func ValidIP(IP string) (v Valid) {
	for i := 0; i < len(IP); i++ {
		switch IP[i] {
		case '.':
			v = ValidIPv4
		case ':':
			v = ValidIPv6
		}
	}
	if ip := net.ParseIP(IP); ip != nil {
		return v
	}
	return InvalidIP
}
