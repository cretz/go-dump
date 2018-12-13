package dump

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"math/big"
	"strings"

	"github.com/cretz/go-dump/pb"
)

func (c *ConversionContext) ConvertExprTypeRef(expr ast.Expr, def bool) *pb.TypeRef {
	if ident, ok := expr.(*ast.Ident); ok {
		obj := c.Pkg.Info.Uses[ident]
		if def {
			obj = c.Pkg.Info.Defs[ident]
		}
		return c.ConvertTypeRef(obj)
	}
	return c.ConvertTypeRef(c.Pkg.Info.Types[expr])
}

func (c *ConversionContext) ConvertTypeRef(t interface{}) *pb.TypeRef {
	if t == nil {
		return nil
	}
	for i, seen := range c.Types {
		if seen == t {
			return &pb.TypeRef{Id: uint32(i)}
		}
	}
	c.Types = append(c.Types, t)
	return &pb.TypeRef{Id: uint32(len(c.Types) - 1)}
}

func (c *ConversionContext) ConvertType(t types.Type) *pb.Type {
	if t == nil {
		return nil
	}
	switch v := t.(type) {
	case types.Object:
		return c.ConvertTypeObject(v)
	case *types.Array:
		return &pb.Type{Type: &pb.Type_TypeArray{
			TypeArray: &pb.TypeArray{Elem: c.ConvertTypeRef(v.Elem()), Len: v.Len()},
		}}
	case *types.Basic:
		return &pb.Type{
			Name: v.Name(),
			Type: &pb.Type_TypeBasic{TypeBasic: &pb.TypeBasic{
				Flags: int32(v.Info()),
				Kind:  pb.TypeBasic_Kind(v.Kind()),
			}},
		}
	case *types.Chan:
		return &pb.Type{Type: &pb.Type_TypeChan{TypeChan: &pb.TypeChan{
			Elem:    c.ConvertTypeRef(v.Elem()),
			SendDir: v.Dir() != types.RecvOnly,
			RecvDir: v.Dir() != types.SendOnly,
		}}}
	case *types.Interface:
		iface := &pb.TypeInterface{
			ExplicitMethods: make([]*pb.TypeRef, v.NumExplicitMethods()),
			Embedded:        make([]*pb.TypeRef, v.NumEmbeddeds()),
		}
		for i := 0; i < v.NumExplicitMethods(); i++ {
			iface.ExplicitMethods[i] = c.ConvertTypeRef(v.ExplicitMethod(i))
		}
		for i := 0; i < v.NumEmbeddeds(); i++ {
			iface.Embedded[i] = c.ConvertTypeRef(v.Embedded(i))
		}
		return &pb.Type{Type: &pb.Type_TypeInterface{TypeInterface: iface}}
	case *types.Map:
		return &pb.Type{Type: &pb.Type_TypeMap{TypeMap: &pb.TypeMap{
			Elem: c.ConvertTypeRef(v.Elem()),
			Key:  c.ConvertTypeRef(v.Key()),
		}}}
	case *types.Named:
		named := &pb.TypeNamed{
			Type:    c.ConvertTypeRef(v.Underlying()),
			Methods: make([]*pb.TypeRef, v.NumMethods()),
		}
		named.TypeName = c.ConvertTypeRef(v.Obj())
		for i := 0; i < v.NumMethods(); i++ {
			named.Methods[i] = c.ConvertTypeRef(v.Method(i))
		}
		ret := &pb.Type{Name: v.Obj().Name(), Type: &pb.Type_TypeNamed{TypeNamed: named}}
		if pkg := v.Obj().Pkg(); pkg != nil {
			ret.Package = pkg.Name()
		}
		return ret
	case *types.Pointer:
		return &pb.Type{Type: &pb.Type_TypePointer{TypePointer: &pb.TypePointer{Elem: c.ConvertTypeRef(v.Elem())}}}
	case *types.Slice:
		return &pb.Type{Type: &pb.Type_TypeSlice{TypeSlice: &pb.TypeSlice{Elem: c.ConvertTypeRef(v.Elem())}}}
	case *types.Struct:
		st := &pb.TypeStruct{Fields: make([]*pb.TypeRef, v.NumFields())}
		for i := 0; i < v.NumFields(); i++ {
			st.Fields[i] = c.ConvertTypeRef(v.Field(i))
		}
		return &pb.Type{Type: &pb.Type_TypeStruct{TypeStruct: st}}
	case *types.Signature:
		sig := &pb.TypeSignature{Variadic: v.Variadic()}
		if recv := v.Recv(); recv != nil {
			sig.Recv = c.ConvertTypeRef(recv)
		}
		if params := v.Params(); params != nil {
			for i := 0; i < params.Len(); i++ {
				sig.Params = append(sig.Params, c.ConvertTypeRef(params.At(i)))
			}
		}
		if results := v.Results(); results != nil {
			for i := 0; i < results.Len(); i++ {
				sig.Results = append(sig.Results, c.ConvertTypeRef(results.At(i)))
			}
		}
		return &pb.Type{Type: &pb.Type_TypeSignature{TypeSignature: sig}}

	case *types.Tuple:
		tup := &pb.TypeTuple{Vars: make([]*pb.TypeRef, v.Len())}
		for i := 0; i < v.Len(); i++ {
			tup.Vars[i] = c.ConvertTypeRef(v.At(i))
		}
		return &pb.Type{Type: &pb.Type_TypeTuple{TypeTuple: tup}}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
}

func (c *ConversionContext) ConvertTypeObjectNameOnly(t types.Object) *pb.Type {
	if t == nil {
		return nil
	}
	ret := &pb.Type{Name: t.Name()}
	if pkg := t.Pkg(); pkg != nil {
		ret.Package = pkg.Name()
	}
	return ret
}

func (c *ConversionContext) ConvertTypeObject(t types.Object) *pb.Type {
	if t == nil {
		return nil
	}
	ret := c.ConvertTypeObjectNameOnly(t)
	switch v := t.(type) {
	case *types.Builtin:
		ret.Type = &pb.Type_TypeBuiltin{TypeBuiltin: true}
	case *types.Const:
		if t.Type() != nil {
			ret.Type = &pb.Type_TypeConst{TypeConst: &pb.TypeConst{
				Type:  c.ConvertTypeRef(t.Type()),
				Value: c.ConvertTypeConstantValue(v.Val()),
			}}
		}
	case *types.Func:
		if t.Type() != nil {
			ret.Type = &pb.Type_TypeFunc{TypeFunc: c.ConvertType(t.Type()).GetTypeSignature()}
		}
	case *types.Label:
		ret.Type = &pb.Type_TypeLabel{TypeLabel: c.ConvertTypeRef(t.Type())}
	case *types.TypeName:
		if t.Type() != nil {
			ret.Type = &pb.Type_TypeName{TypeName: c.ConvertTypeRef(t.Type())}
		}
	case *types.Nil:
		ret.Type = &pb.Type_TypeNil{TypeNil: c.ConvertTypeRef(t.Type())}
	case *types.PkgName:
		ret.Type = &pb.Type_TypePackage{TypePackage: true}
	case *types.Var:
		if v.Type() != nil {
			ret.Type = &pb.Type_TypeVar{TypeVar: &pb.TypeVar{
				Name:     v.Name(),
				Type:     c.ConvertTypeRef(v.Type()),
				Embedded: v.Embedded(),
			}}
		}
	default:
		panic(fmt.Sprintf("Unknown type: %T", t))
	}
	return ret
}

func (c *ConversionContext) ConvertTypeAndValue(t types.TypeAndValue) *pb.Type {
	return &pb.Type{Type: &pb.Type_TypeConst{TypeConst: &pb.TypeConst{
		Type:  c.ConvertTypeRef(t.Type),
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
		// This can be a rational number, which is "a / b", which is not ok for us right now
		str := v.ExactString()
		if strings.Contains(str, "/") {
			// Make into rat, divide as float, get text
			rat, _ := new(big.Rat).SetString(str)
			if rat == nil {
				panic(fmt.Errorf("Unexpected bad rational string: %v", str))
			}
			// Similar to what constant package does
			flt := new(big.Float).SetPrec(512).Quo(
				new(big.Float).SetPrec(512).SetInt(rat.Num()), new(big.Float).SetPrec(512).SetInt(rat.Denom()))
			str = flt.String()
		}
		ret.Value = &pb.ConstantValue_Float{Float: str}
	case constant.Complex:
		ret.Value = &pb.ConstantValue_Complex{Complex: v.ExactString()}
	default:
		panic("Unknown constant value kind")
	}
	return ret
}
