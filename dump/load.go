package dump

import (
	"fmt"

	"golang.org/x/tools/go/loader"

	"github.com/cretz/go-dump/pb"
)

func FromArgs(args []string, includeTests bool) (*pb.Packages, error) {
	conf := &loader.Config{}
	if rest, err := conf.FromArgs(args, includeTests); err != nil {
		return nil, fmt.Errorf("Unable to load args: %v", err)
	} else if len(rest) > 0 {
		return nil, fmt.Errorf("Unrecognized args: %v", rest)
	}
	prog, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed loading: %v", err)
	}
	ret := &pb.Packages{}
	for _, pkg := range prog.Imported {
		convCtx := &ConversionContext{
			FileSet: prog.Fset,
			Pkg:     pkg,
		}
		ret.Packages = append(ret.Packages, convCtx.ConvertPackage())
	}
	return ret, nil
}
