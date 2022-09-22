package tracing

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

const (
	passthroughHeader       = "MESHER_PASS_HEADERS"
	passthroughHeaderPrefix = "MESHER_PASS_HEADER_PERFIXES"
	kubernetesServiceEnv    = "KUBERNETES_SERVICE_HOST"
)

var (
	AppName   string
	istioOpen bool
)

func init() {
	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, syscall.SIGTERM)

	go tracer(stopChan)
	go detectIstio()
	parsePropagateHeaders()
}

func parsePropagateHeaders() {
	headers := os.Getenv(passthroughHeader)
	for _, h := range strings.Split(headers, ",") {
		if len(h) > 0 {
			headersToPropagate = append(headersToPropagate, strings.TrimSpace(h))
		}
	}
	headers = os.Getenv(passthroughHeaderPrefix)
	for _, h := range strings.Split(headers, ",") {
		if len(h) > 0 {
			applicationHeaderPrefix = append(applicationHeaderPrefix, strings.TrimSpace(h))
		}
	}
}

func tracer(stop chan os.Signal) {

	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Println("Jaeger read configuration from Env failed:", err)
		return
	}
	if cfg.Disabled {
		log.Println("Jaeger disabled")
		return
	}
	if len(cfg.Reporter.CollectorEndpoint) > 0 {
		log.Println("Jaeger collectorEndpoint:", cfg.Reporter.CollectorEndpoint)
	} else {
		log.Println("Jaeger localAgentHostPort:", cfg.Reporter.LocalAgentHostPort)
	}

	if len(cfg.ServiceName) == 0 {
		cfg.ServiceName = ServiceName()
	}
	AppName = cfg.ServiceName

	if len(cfg.Sampler.Type) == 0 {
		cfg.Sampler.Type = "const"
		cfg.Sampler.Param = 1
	}
	cfg.Gen128Bit = true
	logger := jaegercfg.Logger(jaeger.StdLogger)
	tracer, closer, err := cfg.NewTracer(logger)
	if err != nil {
		log.Println("Jaeger new tracer failed:", err)
		return
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	<-stop
}

func ServiceName() string {
	if svcName := os.Getenv("SERVICE_NAME"); len(svcName) > 0 {
		return svcName
	}
	if svcName := os.Getenv("JAEGER_SERVICE_NAME"); len(svcName) > 0 {
		if idx := strings.Index(svcName, "."); idx >= 0 {
			return svcName[:idx]
		}
		return svcName
	}
	// 如下代码为推断，不能绝对保证，请尽力设置环境变量
	appName, _ := os.Hostname()
	if len(os.Getenv(kubernetesServiceEnv)) > 0 {
		names := strings.Split(appName, "-")
		if len(names) == 2 {
			appName = names[0]
		} else if len(names) > 2 {
			// statefulset 或 daemonsets
			suffixIndex := len(names) - 1
			hashLen := len(names[len(names)-2])
			if hashLen == 9 || hashLen == 10 {
				// deployment
				suffixIndex = len(names) - 2
			}
			names := names[:suffixIndex]
			if names[len(names)-1] == "gray" {
				names = names[0 : len(names)-1]
			}
			appName = strings.Join(names, "-")
		}
	}
	return appName
}

func detectIstio() {
	/*
	 需要配置 istio 控制面 values.global.proxy.holdApplicationUntilProxyStarts 为 true
	 否则 istio 未启动完成之前，程序会执行
	*/
	conn, err := net.DialTimeout("tcp", "127.0.0.1:15000", time.Second)
	if err != nil {
		log.Println("No istio sidecar detected:", err)
		return
	}
	if conn != nil {
		conn.Close()
		log.Println("Istio sidecar running")
		istioOpen = true
		return
	}
}
