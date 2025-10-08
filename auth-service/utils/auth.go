package utils

import (
	"context"
	"net"

	"github.com/mssola/useragent"
)

func GetUserMetadata(ctx context.Context, key any) (string, bool) {
	val := ctx.Value(key)
	if val == nil {
		return "", false
	}
	return val.(string), true
}

func NormalizeIP(ip string) string {
	parsed := net.ParseIP(ip)

	if parsed == nil {
		return ""
	}

	if ipv4 := parsed.To4(); ipv4 != nil {
		return ipv4.String()
	}

	return parsed.String()
}

func ParseUA(uaString string) (string, string) {
	ua := useragent.New(uaString)
	browser, _ := ua.Browser()
	return browser, ua.OS()
}
