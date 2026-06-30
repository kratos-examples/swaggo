# Changes

Code differences compared to source project.

## buf.gen.yaml (+1 -0)

```diff
@@ -20,3 +20,4 @@
     opt:
       - fq_schema_naming=true
       - default_response=false
+      - title=demo2kratos-title
```

## cmd/demo2kratos/main.go (+21 -0)

```diff
@@ -5,6 +5,7 @@
 	"log/slog"
 	"os"
 
+	"github.com/go-kratos/kratos/contrib/encoding/json/v3"
 	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
 	"github.com/go-kratos/kratos/v3"
 	"github.com/go-kratos/kratos/v3/config"
@@ -32,6 +33,26 @@
 
 func init() {
 	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
+
+	// Configure JSON field naming for HTTP responses.
+	// UseProtoNames=false uses lowerCamelCase, aligned with the protobuf JSON spec
+	// and protobuf-ts generated frontend code (the camelCase convention).
+	// 配置 HTTP 响应的 JSON 字段命名：UseProtoNames=false 用小写驼峰，
+	// 与 protobuf JSON 规范、protobuf-ts 前端生成代码一致。
+	//
+	// Alternative (snake_case): set this to true, AND add "naming=proto" to the
+	// protoc-gen-openapi opt block in buf.gen.yaml, so the served openapi doc uses
+	// snake_case too. The runtime codec and the doc must agree: both camelCase or
+	// both snake_case, otherwise the Swagger UI sends names the backend cannot bind.
+	// 备选方案（snake_case）：把这里改成 true，同时在 buf.gen.yaml 的
+	// protoc-gen-openapi opt 里加上 "naming=proto"，让文档也用 snake_case。
+	// 运行时命名和文档命名必须一致：要么都驼峰、要么都 snake，
+	// 否则 Swagger UI 发出的字段名后端绑不上。
+	json.MarshalOptions.UseProtoNames = false
+
+	// Set UseEnumNumbers to serialize enums as numbers
+	// 设置枚举序列化为数字
+	json.MarshalOptions.UseEnumNumbers = true
 }
 
 func newApp(logger *slog.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
```

## internal/server/http.go (+24 -0)

```diff
@@ -5,9 +5,14 @@
 
 	"github.com/go-kratos/kratos/v3/middleware/recovery"
 	"github.com/go-kratos/kratos/v3/transport/http"
+	"github.com/yylego/kratos-examples/demo2kratos"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-swaggo/swaggokratos"
+	"github.com/yylego/kratos-swaggo/swaggokratos/swaggogin"
+	"github.com/yylego/kratos-zap/zapkratos"
+	"github.com/yylego/zaplog"
 )
 
 func NewHTTPServer(c *conf.Server, article *service.ArticleService, logger *slog.Logger) *http.Server {
@@ -27,5 +32,24 @@
 	}
 	srv := http.NewServer(opts...)
 	pb.RegisterArticleServiceHTTPServer(srv, article)
+
+	serveSwaggerHttpDocument(c, srv)
 	return srv
+}
+
+func serveSwaggerHttpDocument(c *conf.Server, srv *http.Server) {
+	zapKratos := zapkratos.NewZapKratos(zaplog.LOGGER, zapkratos.NewOptions())
+	zapLog := zapKratos.SubZap()
+	zapLog.SUG.Infoln("准备添加接口文档")
+
+	swaggokratos.RegisterSwaggoHTTPServer(srv, "/doc/", []*swaggogin.Param{
+		{
+			SwaggerPath: "/swagger/a/*any",
+			ExplorePath: "/abc/openapi-a.yaml",
+			ContentData: demo2kratos.GetOpenapiContent("demo2kratos-title"),
+		},
+	})
+
+	zapLog.SUG.Infoln("[DOC]", "(http://127.0.0.1:"+swaggokratos.MustGetPortNum(c.Http.Address)+"/doc/swagger/a/index.html)")
+	zapLog.SUG.Infoln("接口文档添加成功")
 }
```

## openapi.go (+24 -0)

```diff
@@ -0,0 +1,24 @@
+package demo2kratos
+
+import (
+	"embed"
+
+	"github.com/yylego/rese"
+	"github.com/yylego/yaml-go-edit/yamlv3edit"
+	"gopkg.in/yaml.v3"
+)
+
+//go:embed openapi.yaml
+var files embed.FS
+
+func GetOpenapiContent(docTitle string) []byte {
+	// 读取文档的内容
+	content := rese.A1(files.ReadFile("openapi.yaml"))
+	// 设置文档的标题
+	content = yamlv3edit.ChangeYamlFieldValue(content, []string{"info", "title"}, func(node *yaml.Node) {
+		if node.Value == "" {
+			node.SetString(docTitle)
+		}
+	})
+	return content
+}
```

## openapi.yaml (+1 -1)

```diff
@@ -3,7 +3,7 @@
 
 openapi: 3.0.3
 info:
-    title: ArticleService API
+    title: demo2kratos-title
     version: 0.0.1
 paths:
     /articles:
```

