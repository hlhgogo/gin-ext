package tracing

import (
	"strings"
	"testing"
)

func TestServiceName(t *testing.T) {
	svcName := "test.xxx"
	if idx := strings.Index(svcName, "."); idx >= 0 {
		svcName = svcName[:idx]
	}
	t.Log(svcName)

	appNames := []string{
		"qw-scrm-markting-material-server-7684d86d8-7vh9v",
		"traefik-ingress-canary-7657f48897-vnp5w",
		"logtail-ds-rm4zl",
		"logtailds-rm4zl",
		"whoami-54466c959b-xdnk9",
		"whoami-gray-54466c959b-xdnk9",
	}
	for _, appName := range appNames {
		names := strings.Split(appName, "-")
		if len(names) == 2 {
			appName = names[0]
		} else if len(names) > 2 {
			suffixIndex := len(names) - 1
			hashLen := len(names[len(names)-2])
			if hashLen == 9 || hashLen == 10 {
				suffixIndex = len(names) - 2
			}
			names := names[:suffixIndex]
			if names[len(names)-1] == "gray" {
				names = names[0 : len(names)-1]
			}
			appName = strings.Join(names, "-")
		}
		t.Log(appName)
	}
}
