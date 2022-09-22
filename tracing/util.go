package tracing

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// BuildShareUrl 构造外链带公司账号（W0000..XXX），用于分享，推送，以及三方回调的url识别流量。
func BuildShareUrl(url, account string) string {
	if len(url) == 0 || len(account) == 0 {
		return url
	}
	accHash := MD5(account)

	var fragment string
	if hashIndex := strings.Index(url, "#"); hashIndex > 0 {
		fragment = url[hashIndex:]
		url = url[0:hashIndex]
	}
	if strings.Index(url, "?") > 0 {
		url = fmt.Sprintf("%s&_acc_=%s", url, accHash)
	} else {
		url = fmt.Sprintf("%s?_acc_=%s", url, accHash)
	}
	if len(fragment) > 0 {
		url = fmt.Sprintf("%s%s", url, fragment)
	}
	return url
}

func MD5(account string) string {
	data := []byte(account)
	return fmt.Sprintf("%x", md5.Sum(data))
}
