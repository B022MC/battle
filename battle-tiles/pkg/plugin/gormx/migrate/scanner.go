package migrate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// ScanModels 扫描指定目录下的模型
func ScanModels(dir string, filterFunc func(name string) bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") &&
			!strings.HasSuffix(info.Name(), "_test.go") {
			if err := parseFile(path, filterFunc); err != nil {
				return fmt.Errorf("parse file %s error: %v", path, err)
			}
		}
		return nil
	})
}

func parseFile(filePath string, filterFunc func(name string) bool) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	for _, decl := range node.Decls {
		// 检查类型声明
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 检查是否是GORM模型(三种情况):
			// 1. 包含gorm.Model
			// 2. 有gorm标签
			// 3. 有TableName方法
			if hasGormModel(structType) || hasGormTags(structType) || hasTableNameMethod(node, typeSpec.Name.Name) {
				//modelType := reflect.StructOf(nil)
				//pkg := reflect.New(modelType).Type().PkgPath()

				model := reflect.New(reflect.StructOf([]reflect.StructField{
					{
						Name: typeSpec.Name.Name,
						Type: reflect.StructOf(nil),
						Tag:  reflect.StructTag(`gorm:"-"`),
					},
				})).Interface()
				if ok := filterFunc(typeSpec.Name.Name); ok {
					continue
				}
				RegisterModel(model)
			}
		}
	}
	return nil
}

// 检查是否包含gorm.Model
func hasGormModel(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if sel, ok := field.Type.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "gorm" && sel.Sel.Name == "Model" {
					return true
				}
			}
		}
	}
	return false
}

// 检查是否有gorm标签
func hasGormTags(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			if strings.Contains(tag, "gorm:") {
				return true
			}
		}
	}
	return false
}

// 检查是否有TableName方法
func hasTableNameMethod(file *ast.File, typeName string) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// 检查是否是方法接收器
		if fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}

		// 获取接收器类型
		recvType := fn.Recv.List[0].Type
		if star, ok := recvType.(*ast.StarExpr); ok {
			recvType = star.X
		}

		if ident, ok := recvType.(*ast.Ident); ok {
			// 检查接收器类型是否匹配且方法名是TableName
			if ident.Name == typeName && fn.Name.Name == "TableName" {
				// 检查方法签名是否匹配: func() string
				if fn.Type.Results != nil && len(fn.Type.Results.List) == 1 {
					if ident, ok := fn.Type.Results.List[0].Type.(*ast.Ident); ok {
						if ident.Name == "string" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
