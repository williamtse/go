package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

func NewLoggerProvider(serviceInfo *ServiceInfo) log.Logger {
	l := log.NewStdLogger(os.Stdout)
	return log.With(
		l,
		//"email.id", serviceInfo.Id,
		//"email.name", serviceInfo.Name,
		//"email.version", serviceInfo.Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"trace_id", tracing.TraceID(),
		//"span_id", tracing.SpanID(),
	)
}
