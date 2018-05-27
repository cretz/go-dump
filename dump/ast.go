package dump

import (
	"go/ast"
	"go/token"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/cretz/go-dump/pb"
	"github.com/golang/protobuf/proto"
)

func ConvertNode(n ast.Node) proto.Message {
	switch v := n.(type) {
	case nil:
		return nil
	case ast.Expr:
		return ConvertExpr(v)
	case ast.Stmt:
		return ConvertStmt(v)
	case ast.Decl:
		return ConvertDecl(v)
	case *ast.Comment:
		return ConvertComment(v)
	case *ast.CommentGroup:
		return ConvertCommentGroup(v)
	default:
		panic("TODO")
	}
}

func ConvertExpr(n ast.Expr) proto.Message {
	panic("TODO")
}

func ConvertStmt(n ast.Stmt) proto.Message {
	panic("TODO")
}

func ConvertDecl(n ast.Decl) proto.Message {
	panic("TODO")
}

func pos(pos token.Pos) int32 { return int32(pos) }
func pbAny(msg proto.Message) *any.Any {
	if msg == nil {
		return nil
	}
	ret, err := ptypes.MarshalAny(msg)
	if err != nil {
		panic(err)
	}
	return ret
}
func pbTok(tok token.Token) pb.Token {
	panic("TODO")
}
func expr(n ast.Expr) *any.Any {
	return pbAny(ConvertExpr(n))
}

func ConvertComment(n *ast.Comment) *pb.Comment {
	if n == nil {
		return nil
	}
	return &pb.Comment{Slash: pos(n.Slash), Text: n.Text}
}

func ConvertCommentGroup(n *ast.CommentGroup) *pb.CommentGroup {
	if n == nil {
		return nil
	}
	r := &pb.CommentGroup{List: make([]*pb.Comment, len(n.List))}
	for i, c := range n.List {
		r.List[i] = ConvertComment(c)
	}
	return r
}

func ConvertField(n *ast.Field) *pb.Field {
	if n == nil {
		return nil
	}
	r := &pb.Field{
		Doc:     ConvertCommentGroup(n.Doc),
		Names:   make([]*pb.Ident, len(n.Names)),
		Type:    expr(n.Type),
		Tag:     ConvertBasicLit(n.Tag),
		Comment: ConvertCommentGroup(n.Comment),
	}
	for i, v := range n.Names {
		r.Names[i] = ConvertIdent(v)
	}
	return r
}

func ConvertFieldList(n *ast.FieldList) *pb.FieldList {
	if n == nil {
		return nil
	}
	r := &pb.FieldList{
		Opening: pos(n.Opening),
		List:    make([]*pb.Field, len(n.List)),
		Closing: pos(n.Closing),
	}
	for i, v := range n.List {
		r.List[i] = ConvertField(v)
	}
	return r
}

func ConvertBadExpr(n *ast.BadExpr) *pb.BadExpr {
	if n == nil {
		return nil
	}
	return &pb.BadExpr{From: pos(n.From), To: pos(n.To)}
}

func ConvertIdent(n *ast.Ident) *pb.Ident {
	if n == nil {
		return nil
	}
	return &pb.Ident{NamePos: pos(n.NamePos), Name: n.Name}
}

func ConvertEllipsis(n *ast.Ellipsis) *pb.Ellipsis {
	if n == nil {
		return nil
	}
	return &pb.Ellipsis{Ellipsis: pos(n.Ellipsis), Elt: expr(n.Elt)}
}

// TODO: n == nil for all funcs below

func ConvertBasicLit(n *ast.BasicLit) *pb.BasicLit {
	return &pb.BasicLit{ValuePos: pos(n.ValuePos), Kind: pbTok(n.Kind), Value: n.Value}
}

func ConvertFuncLit(n *ast.FuncLit) *pb.FuncLit {
	return &pb.FuncLit{Type: ConvertFuncType(n.Type), Body: ConvertBlockStmt(n.Body)}
}

func ConvertCompositeLit(n *ast.CompositeLit) *pb.CompositeLit {
	r := &pb.CompositeLit{
		Type:   expr(n.Type),
		Lbrace: pos(n.Lbrace),
		Elts:   make([]*any.Any, len(n.Elts)),
		Rbrace: pos(n.Rbrace),
	}
	for i, v := range n.Elts {
		r.Elts[i] = expr(v)
	}
	return r
}

func ConvertParenExpr(n *ast.ParenExpr) *pb.ParenExpr {
	return &pb.ParenExpr{Lparen: pos(n.Lparen), X: expr(n.X), Rparen: pos(n.Rparen)}
}

func ConvertSelectorExpr(n *ast.SelectorExpr) *pb.SelectorExpr {
	return &pb.SelectorExpr{X: expr(n.X), Sel: ConvertIdent(n.Sel)}
}

func ConvertIndexExpr(n *ast.IndexExpr) *pb.IndexExpr {
	return &pb.IndexExpr{
		X:      expr(n.X),
		Lbrack: pos(n.Lbrack),
		Index:  expr(n.Index),
		Rbrack: pos(n.Rbrack),
	}
}

func ConvertSliceExpr(n *ast.SliceExpr) *pb.SliceExpr {
	return &pb.SliceExpr{
		X:      expr(n.X),
		Lbrack: pos(n.Lbrack),
		Low:    expr(n.Low),
		High:   expr(n.High),
		Max:    expr(n.Max),
		Slice3: n.Slice3,
		Rbrack: pos(n.Rbrack),
	}
}

func ConvertTypeAssertExpr(n *ast.TypeAssertExpr) *pb.TypeAssertExpr {
	return &pb.TypeAssertExpr{
		X:      expr(n.X),
		Lparen: pos(n.Lparen),
		Type:   expr(n.Type),
		Rparen: pos(n.Rparen),
	}
}

func ConvertCallExpr(n *ast.CallExpr) *pb.CallExpr {
	r := &pb.CallExpr{
		Fun:      expr(n.Fun),
		Lparen:   pos(n.Lparen),
		Args:     make([]*any.Any, len(n.Args)),
		Ellipsis: pos(n.Ellipsis),
		Rparen:   pos(n.Rparen),
	}
	for i, v := range n.Args {
		r.Args[i] = expr(v)
	}
	return r
}

func ConvertStarExpr(n *ast.StarExpr) *pb.StarExpr {
	return &pb.StarExpr{Star: pos(n.Star), X: expr(n.X)}
}

func ConvertUnaryExpr(n *ast.UnaryExpr) *pb.UnaryExpr {
	return &pb.UnaryExpr{OpPos: pos(n.OpPos), Op: pbTok(n.Op), X: expr(n.X)}
}

func ConvertBinaryExpr(n *ast.BinaryExpr) *pb.BinaryExpr {
	return &pb.BinaryExpr{X: expr(n.X), OpPos: pos(n.OpPos), Op: pbTok(n.Op), Y: expr(n.Y)}
}

func ConvertKeyValueExpr(n *ast.KeyValueExpr) *pb.KeyValueExpr {
	return &pb.KeyValueExpr{Key: expr(n.Key), Colon: pos(n.Colon), Value: expr(n.Value)}
}
