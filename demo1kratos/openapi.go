package demo1kratos

import (
	"embed"

	"github.com/yylego/rese"
	"github.com/yylego/yaml-go-edit/yamlv3edit"
)

//go:embed openapi.yaml
var files embed.FS

func GetOpenapiContent(docTitle string) []byte {
	content := rese.A1(files.ReadFile("openapi.yaml"))
	content = yamlv3edit.ModifyYamlFieldValue(content, "info.title", docTitle)
	return content
}
