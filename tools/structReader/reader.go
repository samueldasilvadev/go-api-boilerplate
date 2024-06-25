package structReader

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var goTypeToDbType = map[string]string{
	"bool":      "bool",
	"byte":      "char(1)",
	"string":    "tinytext",
	"float32":   "numeric",
	"float64":   "numeric",
	"int":       "numeric",
	"int8":      "numeric",
	"int16":     "numeric",
	"int32":     "numeric",
	"int64":     "numeric",
	"uint":      "numeric",
	"uint8":     "numeric",
	"uint16":    "numeric",
	"uint32":    "numeric",
	"uint64":    "numeric",
	"time.Time": "datetime",
}

func GenerateHclFileFromDomain(domain string) {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	filepath.Walk("./internal/application/domain/"+domain, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".go") {
			GeneratedBody, err := GetHclDescriptionFromFile(path, info.Name())
			if err != nil {
				return err
			}
			body.AppendBlock(GeneratedBody)
		}
		return nil
	})
	schemaBlock := hclwrite.NewBlock("schema", []string{"skeleton"})
	schemaBody := schemaBlock.Body()
	schemaBody.SetAttributeValue("charset", cty.StringVal("utf8mb4"))
	schemaBody.SetAttributeValue("collate", cty.StringVal("utf8mb4_0900_ai_ci"))
	body.AppendBlock(schemaBlock)
	err := os.WriteFile("tools/migrator/generated/generated_from_domain.my.hcl", file.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetHclDescriptionFromFile(filePath string, filename string) (*hclwrite.Block, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(data), 0)
	if err != nil {
		return nil, err
	}

	var hclDescription *hclwrite.Block
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			hclDescription = GenerateHclTableFromStruct(fmt.Sprintf("%s", f.Name), getMapFromSlice(x.Fields.List))
		}
		return true
	})
	return hclDescription, nil
}

func getMapFromSlice(slice []*ast.Field) map[string]string {
	value := make(map[string]string)
	for _, field := range slice {
		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		value[tag.Get("db")] = fmt.Sprintf("%s", field.Type)
	}
	return value
}

func GenerateHclTableFromStruct(tableName string, fields map[string]string) *hclwrite.Block {
	block := hclwrite.NewBlock("table", []string{tableName})
	blockBody := block.Body()
	blockBody.SetAttributeTraversal("schema", hcl.Traversal{
		hcl.TraverseRoot{
			Name: "schema",
		},
		hcl.TraverseAttr{
			Name: "skeleton",
		},
	})
	for field, golangType := range fields {
		fieldBlock := blockBody.AppendNewBlock("column", []string{field})
		fieldBody := fieldBlock.Body()
		fieldBody.SetAttributeTraversal("type", hcl.Traversal{
			hcl.TraverseRoot{
				Name: goTypeToDbType[golangType],
			},
		})
		fieldBody.SetAttributeValue("null", cty.BoolVal(false))
	}
	return block
}