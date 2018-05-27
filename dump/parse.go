package dump

import (
	"go/parser"
	"go/token"

	"github.com/cretz/go-dump/pb"
)

func ParseDir(path string) (map[string]*pb.Package, error) {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, path, nil, parser.ParseComments)
	ret := map[string]*pb.Package{}
	for k, v := range pkgs {
		ret[k] = ConvertPackage(v)
	}
	return ret, err
}
