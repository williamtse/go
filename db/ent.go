package db

import (
	"context"
	"dingidan-douyin/model/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NewEntClient(driver string, source string, debug bool, h *log.Helper) (*ent.Client, error) {
	h.Info("初始化数据库")
	drv, err := sql.Open(
		driver,
		source,
	)
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		tracer := otel.Tracer("ent.")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
	db := ent.NewClient(ent.Driver(sqlDrv))
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, err
	}

	// 启用 Ent 的调试模式
	h.Info("启用 Ent 的调试模式:", debug)
	if debug {
		db = db.Debug()
	}

	// Run the auto migration tool.
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, err
	}
	return db, nil
}
