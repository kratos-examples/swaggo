# Changes

Code differences compared to source project.

## Makefile (+1 -1)

```diff
@@ -42,7 +42,7 @@
  	       --go_out=paths=source_relative:./api \
  	       --go-http_out=paths=source_relative:./api \
  	       --go-grpc_out=paths=source_relative:./api \
-	       --openapi_out=fq_schema_naming=true,default_response=false:. \
+	       --openapi_out=fq_schema_naming=true,default_response=false,title="demo2kratos-title",naming=proto:. \
 	       $(API_PROTO_FILES)
 
 .PHONY: errors
```

## cmd/demo2kratos/main.go (+7 -0)

```diff
@@ -7,6 +7,7 @@
 	"github.com/go-kratos/kratos/v2"
 	"github.com/go-kratos/kratos/v2/config"
 	"github.com/go-kratos/kratos/v2/config/file"
+	"github.com/go-kratos/kratos/v2/encoding/json"
 	"github.com/go-kratos/kratos/v2/log"
 	"github.com/go-kratos/kratos/v2/middleware/tracing"
 	"github.com/go-kratos/kratos/v2/transport/grpc"
@@ -29,6 +30,12 @@
 
 func init() {
 	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
+
+	//配置http服务回复的json的字段名称风格，是按照proto里写的名称，还是按照小写驼峰的规则
+	//json.MarshalOptions.UseProtoNames = true  //UseProtoNames uses proto field name instead of lowerCamelCase name in JSON
+	//当你不配置的时候，就不能使用proto里自定的名称，而是按照默认的小写驼峰的风格
+	//推荐不要配置，即，使用默认 false 的规则，这样保证和其它语言默认配置相同（否则生成其它语言消息时还要配置（但还学不会如何配置））
+	json.MarshalOptions.UseProtoNames = true
 }
 
 func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
```

## internal/server/http.go (+24 -0)

```diff
@@ -4,9 +4,14 @@
 	"github.com/go-kratos/kratos/v2/log"
 	"github.com/go-kratos/kratos/v2/middleware/recovery"
 	"github.com/go-kratos/kratos/v2/transport/http"
+	"github.com/yylego/kratos-examples/demo2kratos"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-swaggo/swaggokratos"
+	"github.com/yylego/kratos-swaggo/swaggokratos/swaggogin"
+	"github.com/yylego/kratos-zap/zapkratos"
+	"github.com/yylego/zaplog"
 )
 
 func NewHTTPServer(c *conf.Server, article *service.ArticleService, logger log.Logger) *http.Server {
@@ -26,5 +31,24 @@
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
+	zapLog.SUG.Infoln("[DOC]", "(http://127.0.0.1:"+swaggokratos.MustGetPortNum(c.Http.Addr)+"/doc/swagger/a/index.html)")
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

## openapi.yaml (+5 -5)

```diff
@@ -3,7 +3,7 @@
 
 openapi: 3.0.3
 info:
-    title: ArticleService API
+    title: demo2kratos-title
     version: 0.0.1
 paths:
     /articles:
@@ -17,7 +17,7 @@
                   schema:
                     type: integer
                     format: int32
-                - name: pageSize
+                - name: page_size
                   in: query
                   schema:
                     type: integer
@@ -115,7 +115,7 @@
                     type: string
                 content:
                     type: string
-                studentId:
+                student_id:
                     type: string
         article.CreateArticleReply:
             type: object
@@ -129,7 +129,7 @@
                     type: string
                 content:
                     type: string
-                studentId:
+                student_id:
                     type: string
         article.DeleteArticleReply:
             type: object
@@ -165,7 +165,7 @@
                     type: string
                 content:
                     type: string
-                studentId:
+                student_id:
                     type: string
 tags:
     - name: ArticleService
```

