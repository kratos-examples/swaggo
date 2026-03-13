# Changes

Code differences compared to source project.

## internal/server/http.go (+24 -0)

```diff
@@ -4,9 +4,14 @@
 	"github.com/go-kratos/kratos/v2/log"
 	"github.com/go-kratos/kratos/v2/middleware/recovery"
 	"github.com/go-kratos/kratos/v2/transport/http"
+	"github.com/yylego/kratos-examples/demo1kratos"
 	pb "github.com/yylego/kratos-examples/demo1kratos/api/student"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo1kratos/internal/service"
+	"github.com/yylego/kratos-swaggo/swaggokratos"
+	"github.com/yylego/kratos-swaggo/swaggokratos/swaggogin"
+	"github.com/yylego/kratos-zap/zapkratos"
+	"github.com/yylego/zaplog"
 )
 
 func NewHTTPServer(c *conf.Server, student *service.StudentService, logger log.Logger) *http.Server {
@@ -26,5 +31,24 @@
 	}
 	srv := http.NewServer(opts...)
 	pb.RegisterStudentServiceHTTPServer(srv, student)
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
+			ContentData: demo1kratos.GetOpenapiContent("demo1kratos-title"),
+		},
+	})
+
+	zapLog.SUG.Infoln("[DOC]", "(http://127.0.0.1:"+swaggokratos.MustGetPortNum(c.Http.Address)+"/doc/swagger/a/index.html)")
+	zapLog.SUG.Infoln("接口文档添加成功")
 }
```

## openapi.go (+17 -0)

```diff
@@ -0,0 +1,17 @@
+package demo1kratos
+
+import (
+	"embed"
+
+	"github.com/yylego/rese"
+	"github.com/yylego/yaml-go-edit/yamlv3edit"
+)
+
+//go:embed openapi.yaml
+var files embed.FS
+
+func GetOpenapiContent(docTitle string) []byte {
+	content := rese.A1(files.ReadFile("openapi.yaml"))
+	content = yamlv3edit.ModifyYamlFieldValue(content, "info.title", docTitle)
+	return content
+}
```

