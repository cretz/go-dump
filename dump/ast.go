package dump

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/loader"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/cretz/go-dump/pb"
	"github.com/golang/protobuf/proto"
)

type ConversionContext struct {
	FileSet *token.FileSet
	Pkg     *loader.PackageInfo
}

func (c *ConversionContext) ConvertNode(n ast.Node) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case ast.Expr:
		return c.ConvertExpr(v)
	case ast.Stmt:
		return c.ConvertStmt(v)
	case ast.Spec:
		return c.ConvertSpec(v)
	case ast.Decl:
		return c.ConvertDecl(v)
	case *ast.Comment:
		return c.ConvertComment(v)
	case *ast.CommentGroup:
		return c.ConvertCommentGroup(v)
	case *ast.Field:
		return c.ConvertField(v)
	case *ast.FieldList:
		return c.ConvertFieldList(v)
	default:
		panic(fmt.Sprintf("Unknown node type: %T", n))
	}
}

func (c *ConversionContext) ConvertExpr(n ast.Expr) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.BadExpr:
		return c.ConvertBadExpr(v)
	case *ast.Ident:
		return c.ConvertIdent(v)
	case *ast.Ellipsis:
		return c.ConvertEllipsis(v)
	case *ast.BasicLit:
		return c.ConvertBasicLit(v)
	case *ast.FuncLit:
		return c.ConvertFuncLit(v)
	case *ast.CompositeLit:
		return c.ConvertCompositeLit(v)
	case *ast.ParenExpr:
		return c.ConvertParenExpr(v)
	case *ast.SelectorExpr:
		return c.ConvertSelectorExpr(v)
	case *ast.IndexExpr:
		return c.ConvertIndexExpr(v)
	case *ast.SliceExpr:
		return c.ConvertSliceExpr(v)
	case *ast.TypeAssertExpr:
		return c.ConvertTypeAssertExpr(v)
	case *ast.CallExpr:
		return c.ConvertCallExpr(v)
	case *ast.StarExpr:
		return c.ConvertStarExpr(v)
	case *ast.UnaryExpr:
		return c.ConvertUnaryExpr(v)
	case *ast.BinaryExpr:
		return c.ConvertBinaryExpr(v)
	case *ast.KeyValueExpr:
		return c.ConvertKeyValueExpr(v)
	case *ast.ArrayType:
		return c.ConvertArrayType(v)
	case *ast.StructType:
		return c.ConvertStructType(v)
	case *ast.FuncType:
		return c.ConvertFuncType(v)
	case *ast.InterfaceType:
		return c.ConvertInterfaceType(v)
	case *ast.MapType:
		return c.ConvertMapType(v)
	case *ast.ChanType:
		return c.ConvertChanType(v)
	default:
		panic(fmt.Sprintf("Unknown expr type: %T", n))
	}
}

func (c *ConversionContext) ConvertStmt(n ast.Stmt) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.BadStmt:
		return c.ConvertBadStmt(v)
	case *ast.DeclStmt:
		return c.ConvertDeclStmt(v)
	case *ast.EmptyStmt:
		return c.ConvertEmptyStmt(v)
	case *ast.LabeledStmt:
		return c.ConvertLabeledStmt(v)
	case *ast.ExprStmt:
		return c.ConvertExprStmt(v)
	case *ast.SendStmt:
		return c.ConvertSendStmt(v)
	case *ast.IncDecStmt:
		return c.ConvertIncDecStmt(v)
	case *ast.AssignStmt:
		return c.ConvertAssignStmt(v)
	case *ast.GoStmt:
		return c.ConvertGoStmt(v)
	case *ast.DeferStmt:
		return c.ConvertDeferStmt(v)
	case *ast.ReturnStmt:
		return c.ConvertReturnStmt(v)
	case *ast.BranchStmt:
		return c.ConvertBranchStmt(v)
	case *ast.BlockStmt:
		return c.ConvertBlockStmt(v)
	case *ast.IfStmt:
		return c.ConvertIfStmt(v)
	case *ast.CaseClause:
		return c.ConvertCaseClause(v)
	case *ast.SwitchStmt:
		return c.ConvertSwitchStmt(v)
	case *ast.TypeSwitchStmt:
		return c.ConvertTypeSwitchStmt(v)
	case *ast.CommClause:
		return c.ConvertCommClause(v)
	case *ast.SelectStmt:
		return c.ConvertSelectStmt(v)
	case *ast.ForStmt:
		return c.ConvertForStmt(v)
	case *ast.RangeStmt:
		return c.ConvertRangeStmt(v)
	default:
		panic(fmt.Sprintf("Unknown stmt type: %T", n))
	}
}

func (c *ConversionContext) ConvertSpec(n ast.Spec) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.ImportSpec:
		return c.ConvertImportSpec(v)
	case *ast.ValueSpec:
		return c.ConvertValueSpec(v)
	case *ast.TypeSpec:
		return c.ConvertTypeSpec(v)
	default:
		panic(fmt.Sprintf("Unknown spec type: %T", n))
	}
}

func (c *ConversionContext) ConvertDecl(n ast.Decl) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case *ast.BadDecl:
		return c.ConvertBadDecl(v)
	case *ast.GenDecl:
		return c.ConvertGenDecl(v)
	case *ast.FuncDecl:
		return c.ConvertFuncDecl(v)
	default:
		panic(fmt.Sprintf("Unknown decl type: %T", n))
	}
}

func (c *ConversionContext) pos(pos token.Pos) int32 { return int32(pos) }
func (c *ConversionContext) pbAny(msg proto.Message) *any.Any {
	if msg == nil {
		return nil
	}
	ret, err := ptypes.MarshalAny(msg)
	if err != nil {
		fmt.Printf("NO!!! %T - %v\n", msg, msg)
		panic(err)
	}
	return ret
}
func (c *ConversionContext) expr(n ast.Expr) *any.Any         { return c.pbAny(c.ConvertExpr(n)) }
func (c *ConversionContext) stmt(n ast.Stmt) *any.Any         { return c.pbAny(c.ConvertStmt(n)) }
func (c *ConversionContext) spec(n ast.Spec) *any.Any         { return c.pbAny(c.ConvertSpec(n)) }
func (c *ConversionContext) decl(n ast.Decl) *any.Any         { return c.pbAny(c.ConvertDecl(n)) }
func (c *ConversionContext) pbTok(tok token.Token) pb.Token   { return pb.Token(tok) }
func (c *ConversionContext) typeInfo(n ast.Expr) *pb.TypeInfo { return c.ConvertTypeInfo(n, false) }

func (c *ConversionContext) ConvertComment(n *ast.Comment) *pb.Comment {
	if n == nil {
		return nil
	}
	return &pb.Comment{Slash: c.pos(n.Slash), Text: n.Text}
}

func (c *ConversionContext) ConvertCommentGroup(n *ast.CommentGroup) *pb.CommentGroup {
	if n == nil {
		return nil
	}
	r := &pb.CommentGroup{List: make([]*pb.Comment, len(n.List))}
	for i, v := range n.List {
		r.List[i] = c.ConvertComment(v)
	}
	return r
}

func (c *ConversionContext) ConvertField(n *ast.Field) *pb.Field {
	if n == nil {
		return nil
	}
	r := &pb.Field{
		Doc:     c.ConvertCommentGroup(n.Doc),
		Names:   make([]*pb.Ident, len(n.Names)),
		Type:    c.expr(n.Type),
		Tag:     c.ConvertBasicLit(n.Tag),
		Comment: c.ConvertCommentGroup(n.Comment),
	}
	for i, v := range n.Names {
		r.Names[i] = c.ConvertIdent(v)
	}
	return r
}

func (c *ConversionContext) ConvertFieldList(n *ast.FieldList) *pb.FieldList {
	if n == nil {
		return nil
	}
	r := &pb.FieldList{
		Opening: c.pos(n.Opening),
		List:    make([]*pb.Field, len(n.List)),
		Closing: c.pos(n.Closing),
	}
	for i, v := range n.List {
		r.List[i] = c.ConvertField(v)
	}
	return r
}

func (c *ConversionContext) ConvertBadExpr(n *ast.BadExpr) *pb.BadExpr {
	if n == nil {
		return nil
	}
	return &pb.BadExpr{From: c.pos(n.From), To: c.pos(n.To), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertIdent(n *ast.Ident) *pb.Ident {
	if n == nil {
		return nil
	}
	return &pb.Ident{
		NamePos:     c.pos(n.NamePos),
		Name:        n.Name,
		TypeInfo:    c.typeInfo(n),
		DefTypeInfo: c.ConvertTypeInfo(n, true),
	}
}

func (c *ConversionContext) ConvertEllipsis(n *ast.Ellipsis) *pb.Ellipsis {
	if n == nil {
		return nil
	}
	return &pb.Ellipsis{Ellipsis: c.pos(n.Ellipsis), Elt: c.expr(n.Elt), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertBasicLit(n *ast.BasicLit) *pb.BasicLit {
	if n == nil {
		return nil
	}
	return &pb.BasicLit{ValuePos: c.pos(n.ValuePos), Kind: c.pbTok(n.Kind), Value: n.Value, TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertFuncLit(n *ast.FuncLit) *pb.FuncLit {
	if n == nil {
		return nil
	}
	return &pb.FuncLit{Type: c.ConvertFuncType(n.Type), Body: c.ConvertBlockStmt(n.Body), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertCompositeLit(n *ast.CompositeLit) *pb.CompositeLit {
	if n == nil {
		return nil
	}
	r := &pb.CompositeLit{
		Type:     c.expr(n.Type),
		Lbrace:   c.pos(n.Lbrace),
		Elts:     make([]*any.Any, len(n.Elts)),
		Rbrace:   c.pos(n.Rbrace),
		TypeInfo: c.typeInfo(n),
	}
	for i, v := range n.Elts {
		r.Elts[i] = c.expr(v)
	}
	return r
}

func (c *ConversionContext) ConvertParenExpr(n *ast.ParenExpr) *pb.ParenExpr {
	if n == nil {
		return nil
	}
	return &pb.ParenExpr{Lparen: c.pos(n.Lparen), X: c.expr(n.X), Rparen: c.pos(n.Rparen), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertSelectorExpr(n *ast.SelectorExpr) *pb.SelectorExpr {
	if n == nil {
		return nil
	}
	return &pb.SelectorExpr{X: c.expr(n.X), Sel: c.ConvertIdent(n.Sel), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertIndexExpr(n *ast.IndexExpr) *pb.IndexExpr {
	if n == nil {
		return nil
	}
	return &pb.IndexExpr{
		X:        c.expr(n.X),
		Lbrack:   c.pos(n.Lbrack),
		Index:    c.expr(n.Index),
		Rbrack:   c.pos(n.Rbrack),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertSliceExpr(n *ast.SliceExpr) *pb.SliceExpr {
	if n == nil {
		return nil
	}
	return &pb.SliceExpr{
		X:        c.expr(n.X),
		Lbrack:   c.pos(n.Lbrack),
		Low:      c.expr(n.Low),
		High:     c.expr(n.High),
		Max:      c.expr(n.Max),
		Slice3:   n.Slice3,
		Rbrack:   c.pos(n.Rbrack),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertTypeAssertExpr(n *ast.TypeAssertExpr) *pb.TypeAssertExpr {
	if n == nil {
		return nil
	}
	return &pb.TypeAssertExpr{
		X:        c.expr(n.X),
		Lparen:   c.pos(n.Lparen),
		Type:     c.expr(n.Type),
		Rparen:   c.pos(n.Rparen),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertCallExpr(n *ast.CallExpr) *pb.CallExpr {
	if n == nil {
		return nil
	}
	r := &pb.CallExpr{
		Fun:      c.expr(n.Fun),
		Lparen:   c.pos(n.Lparen),
		Args:     make([]*any.Any, len(n.Args)),
		Ellipsis: c.pos(n.Ellipsis),
		Rparen:   c.pos(n.Rparen),
		TypeInfo: c.typeInfo(n),
	}
	for i, v := range n.Args {
		r.Args[i] = c.expr(v)
	}
	return r
}

func (c *ConversionContext) ConvertStarExpr(n *ast.StarExpr) *pb.StarExpr {
	if n == nil {
		return nil
	}
	return &pb.StarExpr{Star: c.pos(n.Star), X: c.expr(n.X), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertUnaryExpr(n *ast.UnaryExpr) *pb.UnaryExpr {
	if n == nil {
		return nil
	}
	return &pb.UnaryExpr{OpPos: c.pos(n.OpPos), Op: c.pbTok(n.Op), X: c.expr(n.X), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertBinaryExpr(n *ast.BinaryExpr) *pb.BinaryExpr {
	if n == nil {
		return nil
	}
	return &pb.BinaryExpr{
		X:        c.expr(n.X),
		OpPos:    c.pos(n.OpPos),
		Op:       c.pbTok(n.Op),
		Y:        c.expr(n.Y),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertKeyValueExpr(n *ast.KeyValueExpr) *pb.KeyValueExpr {
	if n == nil {
		return nil
	}
	return &pb.KeyValueExpr{Key: c.expr(n.Key), Colon: c.pos(n.Colon), Value: c.expr(n.Value), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertArrayType(n *ast.ArrayType) *pb.ArrayType {
	if n == nil {
		return nil
	}
	return &pb.ArrayType{Lbrack: c.pos(n.Lbrack), Len: c.expr(n.Len), Elt: c.expr(n.Elt), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertStructType(n *ast.StructType) *pb.StructType {
	if n == nil {
		return nil
	}
	return &pb.StructType{
		Struct:     c.pos(n.Struct),
		Fields:     c.ConvertFieldList(n.Fields),
		Incomplete: n.Incomplete,
		TypeInfo:   c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertFuncType(n *ast.FuncType) *pb.FuncType {
	if n == nil {
		return nil
	}
	return &pb.FuncType{
		Func:     c.pos(n.Func),
		Params:   c.ConvertFieldList(n.Params),
		Results:  c.ConvertFieldList(n.Results),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertInterfaceType(n *ast.InterfaceType) *pb.InterfaceType {
	if n == nil {
		return nil
	}
	return &pb.InterfaceType{
		Interface:  c.pos(n.Interface),
		Methods:    c.ConvertFieldList(n.Methods),
		Incomplete: n.Incomplete,
		TypeInfo:   c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertMapType(n *ast.MapType) *pb.MapType {
	if n == nil {
		return nil
	}
	return &pb.MapType{Map: c.pos(n.Map), Key: c.expr(n.Key), Value: c.expr(n.Value), TypeInfo: c.typeInfo(n)}
}

func (c *ConversionContext) ConvertChanType(n *ast.ChanType) *pb.ChanType {
	if n == nil {
		return nil
	}
	return &pb.ChanType{
		Begin:    c.pos(n.Begin),
		Arrow:    c.pos(n.Arrow),
		SendDir:  n.Dir&ast.SEND != 0,
		RecvDir:  n.Dir&ast.RECV != 0,
		Value:    c.expr(n.Value),
		TypeInfo: c.typeInfo(n),
	}
}

func (c *ConversionContext) ConvertBadStmt(n *ast.BadStmt) *pb.BadStmt {
	if n == nil {
		return nil
	}
	return &pb.BadStmt{From: c.pos(n.From), To: c.pos(n.To)}
}

func (c *ConversionContext) ConvertDeclStmt(n *ast.DeclStmt) *pb.DeclStmt {
	if n == nil {
		return nil
	}
	return &pb.DeclStmt{Decl: c.decl(n.Decl)}
}

func (c *ConversionContext) ConvertEmptyStmt(n *ast.EmptyStmt) *pb.EmptyStmt {
	if n == nil {
		return nil
	}
	return &pb.EmptyStmt{Semicolon: c.pos(n.Semicolon), Implicit: n.Implicit}
}

func (c *ConversionContext) ConvertLabeledStmt(n *ast.LabeledStmt) *pb.LabeledStmt {
	if n == nil {
		return nil
	}
	return &pb.LabeledStmt{Label: c.ConvertIdent(n.Label), Colon: c.pos(n.Colon), Stmt: c.stmt(n.Stmt)}
}

func (c *ConversionContext) ConvertExprStmt(n *ast.ExprStmt) *pb.ExprStmt {
	if n == nil {
		return nil
	}
	return &pb.ExprStmt{X: c.expr(n.X)}
}

func (c *ConversionContext) ConvertSendStmt(n *ast.SendStmt) *pb.SendStmt {
	if n == nil {
		return nil
	}
	return &pb.SendStmt{Chan: c.expr(n.Chan), Arrow: c.pos(n.Arrow), Value: c.expr(n.Value)}
}

func (c *ConversionContext) ConvertIncDecStmt(n *ast.IncDecStmt) *pb.IncDecStmt {
	if n == nil {
		return nil
	}
	return &pb.IncDecStmt{X: c.expr(n.X), TokPos: c.pos(n.TokPos), Tok: c.pbTok(n.Tok)}
}

func (c *ConversionContext) ConvertAssignStmt(n *ast.AssignStmt) *pb.AssignStmt {
	if n == nil {
		return nil
	}
	r := &pb.AssignStmt{
		Lhs:    make([]*any.Any, len(n.Lhs)),
		TokPos: c.pos(n.TokPos),
		Tok:    c.pbTok(n.Tok),
		Rhs:    make([]*any.Any, len(n.Rhs)),
	}
	for i, v := range n.Lhs {
		r.Lhs[i] = c.expr(v)
	}
	for i, v := range n.Rhs {
		r.Rhs[i] = c.expr(v)
	}
	return r
}

func (c *ConversionContext) ConvertGoStmt(n *ast.GoStmt) *pb.GoStmt {
	if n == nil {
		return nil
	}
	return &pb.GoStmt{Go: c.pos(n.Go), Call: c.ConvertCallExpr(n.Call)}
}

func (c *ConversionContext) ConvertDeferStmt(n *ast.DeferStmt) *pb.DeferStmt {
	if n == nil {
		return nil
	}
	return &pb.DeferStmt{Defer: c.pos(n.Defer), Call: c.ConvertCallExpr(n.Call)}
}

func (c *ConversionContext) ConvertReturnStmt(n *ast.ReturnStmt) *pb.ReturnStmt {
	if n == nil {
		return nil
	}
	r := &pb.ReturnStmt{Return: c.pos(n.Return), Results: make([]*any.Any, len(n.Results))}
	for i, v := range n.Results {
		r.Results[i] = c.expr(v)
	}
	return r
}

func (c *ConversionContext) ConvertBranchStmt(n *ast.BranchStmt) *pb.BranchStmt {
	if n == nil {
		return nil
	}
	return &pb.BranchStmt{TokPos: c.pos(n.TokPos), Tok: c.pbTok(n.Tok), Label: c.ConvertIdent(n.Label)}
}

func (c *ConversionContext) ConvertBlockStmt(n *ast.BlockStmt) *pb.BlockStmt {
	if n == nil {
		return nil
	}
	r := &pb.BlockStmt{
		Lbrace: c.pos(n.Lbrace),
		List:   make([]*any.Any, len(n.List)),
		Rbrace: c.pos(n.Rbrace),
	}
	for i, v := range n.List {
		r.List[i] = c.stmt(v)
	}
	return r
}

func (c *ConversionContext) ConvertIfStmt(n *ast.IfStmt) *pb.IfStmt {
	if n == nil {
		return nil
	}
	return &pb.IfStmt{
		If:   c.pos(n.If),
		Init: c.stmt(n.Init),
		Cond: c.expr(n.Cond),
		Body: c.ConvertBlockStmt(n.Body),
		Else: c.stmt(n.Else),
	}
}

func (c *ConversionContext) ConvertCaseClause(n *ast.CaseClause) *pb.CaseClause {
	if n == nil {
		return nil
	}
	r := &pb.CaseClause{
		Case:  c.pos(n.Case),
		List:  make([]*any.Any, len(n.List)),
		Colon: c.pos(n.Colon),
		Body:  make([]*any.Any, len(n.Body)),
	}
	for i, v := range n.List {
		r.List[i] = c.expr(v)
	}
	for i, v := range n.Body {
		r.Body[i] = c.stmt(v)
	}
	return r
}

func (c *ConversionContext) ConvertSwitchStmt(n *ast.SwitchStmt) *pb.SwitchStmt {
	if n == nil {
		return nil
	}
	return &pb.SwitchStmt{
		Switch: c.pos(n.Switch),
		Init:   c.stmt(n.Init),
		Tag:    c.expr(n.Tag),
		Body:   c.ConvertBlockStmt(n.Body),
	}
}

func (c *ConversionContext) ConvertTypeSwitchStmt(n *ast.TypeSwitchStmt) *pb.TypeSwitchStmt {
	if n == nil {
		return nil
	}
	return &pb.TypeSwitchStmt{
		Switch: c.pos(n.Switch),
		Init:   c.stmt(n.Init),
		Assign: c.stmt(n.Assign),
		Body:   c.ConvertBlockStmt(n.Body),
	}
}

func (c *ConversionContext) ConvertCommClause(n *ast.CommClause) *pb.CommClause {
	if n == nil {
		return nil
	}
	r := &pb.CommClause{
		Case:  c.pos(n.Case),
		Comm:  c.stmt(n.Comm),
		Colon: c.pos(n.Colon),
		Body:  make([]*any.Any, len(n.Body)),
	}
	for i, v := range n.Body {
		r.Body[i] = c.stmt(v)
	}
	return r
}

func (c *ConversionContext) ConvertSelectStmt(n *ast.SelectStmt) *pb.SelectStmt {
	if n == nil {
		return nil
	}
	return &pb.SelectStmt{Select: c.pos(n.Select), Body: c.ConvertBlockStmt(n.Body)}
}

func (c *ConversionContext) ConvertForStmt(n *ast.ForStmt) *pb.ForStmt {
	if n == nil {
		return nil
	}
	return &pb.ForStmt{
		For:  c.pos(n.For),
		Init: c.stmt(n.Init),
		Cond: c.expr(n.Cond),
		Post: c.stmt(n.Post),
		Body: c.ConvertBlockStmt(n.Body),
	}
}

func (c *ConversionContext) ConvertRangeStmt(n *ast.RangeStmt) *pb.RangeStmt {
	if n == nil {
		return nil
	}
	return &pb.RangeStmt{
		For:    c.pos(n.For),
		Key:    c.expr(n.Key),
		Value:  c.expr(n.Value),
		TokPos: c.pos(n.TokPos),
		Tok:    c.pbTok(n.Tok),
		X:      c.expr(n.X),
		Body:   c.ConvertBlockStmt(n.Body),
	}
}

func (c *ConversionContext) ConvertImportSpec(n *ast.ImportSpec) *pb.ImportSpec {
	if n == nil {
		return nil
	}
	return &pb.ImportSpec{
		Doc:     c.ConvertCommentGroup(n.Doc),
		Name:    c.ConvertIdent(n.Name),
		Path:    c.ConvertBasicLit(n.Path),
		Comment: c.ConvertCommentGroup(n.Comment),
		EndPos:  c.pos(n.EndPos),
	}
}

func (c *ConversionContext) ConvertValueSpec(n *ast.ValueSpec) *pb.ValueSpec {
	if n == nil {
		return nil
	}
	r := &pb.ValueSpec{
		Doc:     c.ConvertCommentGroup(n.Doc),
		Names:   make([]*pb.Ident, len(n.Names)),
		Type:    c.expr(n.Type),
		Values:  make([]*any.Any, len(n.Values)),
		Comment: c.ConvertCommentGroup(n.Comment),
	}
	for i, v := range n.Names {
		r.Names[i] = c.ConvertIdent(v)
	}
	for i, v := range n.Values {
		r.Values[i] = c.expr(v)
	}
	return r
}

func (c *ConversionContext) ConvertTypeSpec(n *ast.TypeSpec) *pb.TypeSpec {
	if n == nil {
		return nil
	}
	return &pb.TypeSpec{
		Doc:     c.ConvertCommentGroup(n.Doc),
		Name:    c.ConvertIdent(n.Name),
		Assign:  c.pos(n.Assign),
		Type:    c.expr(n.Type),
		Comment: c.ConvertCommentGroup(n.Comment),
	}
}

func (c *ConversionContext) ConvertBadDecl(n *ast.BadDecl) *pb.BadDecl {
	return &pb.BadDecl{From: c.pos(n.From), To: c.pos(n.To)}
}

func (c *ConversionContext) ConvertGenDecl(n *ast.GenDecl) *pb.GenDecl {
	if n == nil {
		return nil
	}
	r := &pb.GenDecl{
		Doc:    c.ConvertCommentGroup(n.Doc),
		TokPos: c.pos(n.TokPos),
		Tok:    c.pbTok(n.Tok),
		Lparen: c.pos(n.Lparen),
		Specs:  make([]*any.Any, len(n.Specs)),
		Rparen: c.pos(n.Rparen),
	}
	for i, v := range n.Specs {
		r.Specs[i] = c.spec(v)
	}
	return r
}

func (c *ConversionContext) ConvertFuncDecl(n *ast.FuncDecl) *pb.FuncDecl {
	if n == nil {
		return nil
	}
	return &pb.FuncDecl{
		Doc:  c.ConvertCommentGroup(n.Doc),
		Recv: c.ConvertFieldList(n.Recv),
		Name: c.ConvertIdent(n.Name),
		Type: c.ConvertFuncType(n.Type),
		Body: c.ConvertBlockStmt(n.Body),
	}
}

func (c *ConversionContext) ConvertFile(n *ast.File) *pb.File {
	if n == nil {
		return nil
	}
	r := &pb.File{
		FileName:   filepath.Base(c.FileSet.File(n.Pos()).Name()),
		Doc:        c.ConvertCommentGroup(n.Doc),
		Package:    c.pos(n.Package),
		Name:       c.ConvertIdent(n.Name),
		Decls:      make([]*any.Any, len(n.Decls)),
		Imports:    make([]*pb.ImportSpec, len(n.Imports)),
		Unresolved: make([]*pb.Ident, len(n.Unresolved)),
		Comments:   make([]*pb.CommentGroup, len(n.Comments)),
	}
	for i, v := range n.Decls {
		r.Decls[i] = c.decl(v)
	}
	for i, v := range n.Imports {
		r.Imports[i] = c.ConvertImportSpec(v)
	}
	for i, v := range n.Unresolved {
		r.Unresolved[i] = c.ConvertIdent(v)
	}
	for i, v := range n.Comments {
		r.Comments[i] = c.ConvertCommentGroup(v)
	}
	return r
}

func (c *ConversionContext) ConvertPackage() *pb.Package {
	if c.Pkg == nil {
		return nil
	}
	r := &pb.Package{
		Name: c.Pkg.Pkg.Name(),
		Path: c.Pkg.Pkg.Path(),
	}
	for _, file := range c.Pkg.Files {
		r.Files = append(r.Files, c.ConvertFile(file))
	}
	for _, varInit := range c.Pkg.InitOrder {
		for _, v := range varInit.Lhs {
			r.VarInitOrder = append(r.VarInitOrder, v.Name())
		}
	}
	return r
}
