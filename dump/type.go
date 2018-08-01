package dump

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"

	"github.com/cretz/go-dump/pb"
)

func (c *ConversionContext) ConvertTypeInfo(expr ast.Expr, def bool) *pb.TypeInfo {
	if ident, ok := expr.(*ast.Ident); ok {
		obj := c.Pkg.Info.Uses[ident]
		if def {
			obj = c.Pkg.Info.Defs[ident]
		}
		return c.ConvertTypeObject(obj)
	}
	return c.ConvertTypeAndValue(c.Pkg.Info.Types[expr])
}

func (c *ConversionContext) ConvertType(t types.Type) *pb.TypeInfo {
	if t == nil {
		return nil
	}
	switch v := t.(type) {
	case types.Object:
		return c.ConvertTypeObject(v)
	case *types.Array:
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeArray{
			TypeArray: &pb.TypeArray{Elem: c.ConvertType(v.Elem()), Len: v.Len()},
		}}
	case *types.Basic:
		return &pb.TypeInfo{
			Name: v.Name(),
			Type: &pb.TypeInfo_TypeBasic{TypeBasic: &pb.TypeBasic{
				Flags: int32(v.Info()),
				Kind:  pb.TypeBasic_Kind(v.Kind()),
			}},
		}
	case *types.Chan:
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeChan{TypeChan: &pb.TypeChan{
			Elem:    c.ConvertType(v.Elem()),
			SendDir: v.Dir() != types.RecvOnly,
			RecvDir: v.Dir() != types.SendOnly,
		}}}
	case *types.Interface:
		iface := &pb.TypeInterface{
			ExplicitMethods: make([]*pb.TypeInfo, v.NumExplicitMethods()),
			Embedded:        make([]*pb.TypeInfo, v.NumEmbeddeds()),
		}
		for i := 0; i < v.NumExplicitMethods(); i++ {
			iface.ExplicitMethods[i] = c.ConvertTypeObject(v.ExplicitMethod(i))
		}
		for i := 0; i < v.NumEmbeddeds(); i++ {
			iface.Embedded[i] = c.ConvertType(v.Embedded(i))
		}
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeInterface{TypeInterface: iface}}
	case *types.Map:
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeMap{TypeMap: &pb.TypeMap{
			Elem: c.ConvertType(v.Elem()),
			Key:  c.ConvertType(v.Key()),
		}}}
	case *types.Named:
		named := &pb.TypeNamed{
			Type:    c.ConvertType(v.Underlying()),
			Methods: make([]*pb.TypeInfo, v.NumMethods()),
		}
		if v.Obj().Type() == v {
			named.TypeName = c.ConvertTypeObjectNameOnly(v.Obj())
		} else {
			named.TypeName = c.ConvertTypeObject(v.Obj())
		}
		for i := 0; i < v.NumMethods(); i++ {
			named.Methods[i] = c.ConvertTypeObject(v.Method(i))
		}
		ret := &pb.TypeInfo{Name: v.Obj().Name(), Type: &pb.TypeInfo_TypeNamed{TypeNamed: named}}
		if pkg := v.Obj().Pkg(); pkg != nil {
			ret.Package = pkg.Name()
		}
		return ret
	case *types.Pointer:
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypePointer{TypePointer: &pb.TypePointer{Elem: c.ConvertType(v.Elem())}}}
	case *types.Slice:
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeSlice{TypeSlice: &pb.TypeSlice{Elem: c.ConvertType(v.Elem())}}}
	case *types.Struct:
		st := &pb.TypeStruct{Fields: make([]*pb.TypeInfo, v.NumFields())}
		for i := 0; i < v.NumFields(); i++ {
			st.Fields[i] = c.ConvertTypeObject(v.Field(i))
		}
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeStruct{TypeStruct: st}}
	case *types.Signature:
		sig := &pb.TypeSignature{}
		if recv := v.Recv(); recv != nil {
			sig.Recv = c.ConvertTypeObjectNameOnly(recv)
		}
		if params := v.Params(); params != nil {
			sig.Params = c.ConvertType(params).Type.(*pb.TypeInfo_TypeTuple).TypeTuple.Vars
		}
		if results := v.Results(); results != nil {
			sig.Results = c.ConvertType(results).Type.(*pb.TypeInfo_TypeTuple).TypeTuple.Vars
		}
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeSignature{TypeSignature: sig}}

	case *types.Tuple:
		tup := &pb.TypeTuple{Vars: make([]*pb.TypeInfo, v.Len())}
		for i := 0; i < v.Len(); i++ {
			tup.Vars[i] = c.ConvertTypeObjectNameOnly(v.At(i))
		}
		return &pb.TypeInfo{Type: &pb.TypeInfo_TypeTuple{TypeTuple: tup}}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
}

func (c *ConversionContext) ConvertTypeObjectNameOnly(t types.Object) *pb.TypeInfo {
	if t == nil {
		return nil
	}
	ret := &pb.TypeInfo{Name: t.Name()}
	if pkg := t.Pkg(); pkg != nil {
		ret.Package = pkg.Name()
	}
	return ret
}

func (c *ConversionContext) ConvertTypeObject(t types.Object) *pb.TypeInfo {
	if t == nil {
		return nil
	}
	ret := c.ConvertTypeObjectNameOnly(t)
	typ := c.ConvertType(t.Type())
	switch v := t.(type) {
	case *types.Builtin:
		ret.Type = &pb.TypeInfo_TypeBuiltin{TypeBuiltin: true}
	case *types.Const:
		if typ != nil {
			ret.Type = &pb.TypeInfo_TypeConst{TypeConst: &pb.TypeConst{
				Type:  typ,
				Value: c.ConvertTypeConstantValue(v.Val()),
			}}
		}
	case *types.Func:
		if typ != nil {
			ret.Type = &pb.TypeInfo_TypeFunc{TypeFunc: typ.GetTypeSignature()}
		}
	case *types.Label:
		ret.Type = &pb.TypeInfo_TypeLabel{TypeLabel: typ}
	case *types.TypeName:
		if typ != nil {
			ret.Type = &pb.TypeInfo_TypeName{TypeName: typ}
		}
	case *types.Nil:
		ret.Type = &pb.TypeInfo_TypeNil{TypeNil: typ}
	case *types.PkgName:
		ret.Type = &pb.TypeInfo_TypePackage{TypePackage: true}
	case *types.Var:
		if typ != nil {
			ret.Type = &pb.TypeInfo_TypeVar{TypeVar: typ}
		}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
	return ret
}

func (c *ConversionContext) ConvertTypeAndValue(t types.TypeAndValue) *pb.TypeInfo {
	if t.Type == nil {
		return nil
	}
	return &pb.TypeInfo{Type: &pb.TypeInfo_TypeConst{TypeConst: &pb.TypeConst{
		Type:  c.ConvertType(t.Type),
		Value: c.ConvertTypeConstantValue(t.Value),
	}}}
}

func (c *ConversionContext) ConvertTypeConstantValue(v constant.Value) *pb.ConstantValue {
	if v == nil {
		return nil
	}
	ret := &pb.ConstantValue{}
	switch v.Kind() {
	case constant.Unknown:
		ret.Value = &pb.ConstantValue_Unknown{Unknown: v.ExactString()}
	case constant.Bool:
		ret.Value = &pb.ConstantValue_Bool{Bool: constant.BoolVal(v)}
	case constant.String:
		ret.Value = &pb.ConstantValue_String_{String_: constant.StringVal(v)}
	case constant.Int:
		ret.Value = &pb.ConstantValue_Int{Int: v.ExactString()}
	case constant.Float:
		f, _ := constant.Float64Val(v)
		ret.Value = &pb.ConstantValue_Float{Float: f}
	case constant.Complex:
		ret.Value = &pb.ConstantValue_Complex{Complex: v.ExactString()}
	default:
		panic("Unknown constant value kind")
	}
	return ret
}
