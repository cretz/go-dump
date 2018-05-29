package dump

import (
	"encoding/json"

	"github.com/cretz/go-dump/pb"
	"github.com/golang/protobuf/jsonpb"
)

var marshaller jsonpb.Marshaler

func PackagesToJSONMap(pkgs map[string]*pb.Package) (map[string]map[string]interface{}, error) {
	ret := map[string]map[string]interface{}{}
	var err error
	for k, pkg := range pkgs {
		if ret[k], err = PackageToJSONMap(pkg); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func PackageToJSONMap(p *pb.Package) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	str, err := marshaller.MarshalToString(p)
	if err == nil {
		err = json.Unmarshal([]byte(str), &ret)
	}
	return ret, err
}
