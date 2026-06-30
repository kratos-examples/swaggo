package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/contrib/encoding/json/v3"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/yylego/done"
	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
	"github.com/yylego/must"
	"github.com/yylego/rese"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")

	// Configure JSON field naming for HTTP responses.
	// UseProtoNames=false uses lowerCamelCase, aligned with the protobuf JSON spec
	// and protobuf-ts generated frontend code (the camelCase convention).
	// 配置 HTTP 响应的 JSON 字段命名：UseProtoNames=false 用小写驼峰，
	// 与 protobuf JSON 规范、protobuf-ts 前端生成代码一致。
	//
	// Alternative (snake_case): set this to true, AND add "naming=proto" to the
	// protoc-gen-openapi opt block in buf.gen.yaml, so the served openapi doc uses
	// snake_case too. The runtime codec and the doc must agree: both camelCase or
	// both snake_case, otherwise the Swagger UI sends names the backend cannot bind.
	// 备选方案（snake_case）：把这里改成 true，同时在 buf.gen.yaml 的
	// protoc-gen-openapi opt 里加上 "naming=proto"，让文档也用 snake_case。
	// 运行时命名和文档命名必须一致：要么都驼峰、要么都 snake，
	// 否则 Swagger UI 发出的字段名后端绑不上。
	json.MarshalOptions.UseProtoNames = false

	// Set UseEnumNumbers to serialize enums as numbers
	// 设置枚举序列化为数字
	json.MarshalOptions.UseEnumNumbers = true
}

func newApp(logger *slog.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(done.VCE(os.Hostname()).Omit()),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.NewLogger(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
		log.WithExtractor(tracing.TraceAttrs),
	).With(
		slog.String("service.id", done.VCE(os.Hostname()).Omit()),
		slog.String("service.name", Name),
		slog.String("service.version", Version),
	)
	log.SetDefault(logger)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer rese.F0(c.Close)

	must.Done(c.Load())

	var cfg conf.Bootstrap
	must.Done(c.Scan(&cfg))

	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, logger))
	defer cleanup()

	// start and wait for stop signal
	must.Done(app.Run())
}
