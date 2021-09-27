package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type BitwiseOr struct {
	l Node
	r Node
}

func NewBitwiseOr(l, r Node) Node {
	return BitwiseOr{l, r}
}

func (b BitwiseOr) Eval(env *obj.Env) obj.Object {
	var (
		left  = unwrap(b.l.Eval(env))
		right = unwrap(b.r.Eval(env))
	)

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if !assertTypes(left, obj.IntType) {
		return obj.NewError("unsupported operator '|' for type %v", left.Type())
	}
	if !assertTypes(right, obj.IntType) {
		return obj.NewError("unsupported operator '|' for type %v", right.Type())
	}
	l := left.(*obj.Integer).Val()
	r := right.(*obj.Integer).Val()
	return obj.NewInteger(l | r)
}

func (b BitwiseOr) String() string {
	return fmt.Sprintf("(%v | %v)", b.l, b.r)
}
