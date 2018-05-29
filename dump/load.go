package dump

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/cretz/go-dump/pb"
)

func LoadDir(path string) (map[string]*pb.Package, error) {
	fileSet := token.NewFileSet()
	// Parse
	parsedPkgs, err := parser.ParseDir(fileSet, path, nil, parser.ParseComments)
	// TODO: keep going on error?
	if err != nil {
		return nil, err
	}
	// Create conv contexts while type checking
	convCtxs := []*ConversionContext{}
	typesConf := types.Config{Importer: importer.Default()}
	for _, parsedPkg := range parsedPkgs {
		convCtx := &ConversionContext{
			ASTPackage: parsedPkg,
			TypeInfo: &types.Info{
				Types: map[ast.Expr]types.TypeAndValue{},
				Defs:  map[*ast.Ident]types.Object{},
				Uses:  map[*ast.Ident]types.Object{},
			},
		}
		files := []*ast.File{}
		for _, file := range parsedPkg.Files {
			files = append(files, file)
		}
		if convCtx.TypePackage, err = typesConf.Check(path, fileSet, files, convCtx.TypeInfo); err != nil {
			return nil, err
		}
		convCtxs = append(convCtxs, convCtx)
	}
	// Convert
	ret := map[string]*pb.Package{}
	for _, convCtx := range convCtxs {
		ret[convCtx.ASTPackage.Name] = convCtx.ConvertPackage()
	}
	return ret, err
}
