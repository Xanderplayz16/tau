package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type Greater struct {
	l Node
	r Node
}

func NewGreater(l, r Node) Node {
	return Greater{l, r}
}

func (g Greater) Eval(env *obj.Env) obj.Object {
	var (
		left  = unwrap(g.l.Eval(env))
		right = unwrap(g.r.Eval(env))
	)

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if !assertTypes(left, obj.IntType, obj.FloatType) {
		return obj.NewError("unsupported operator '>' for type %v", left.Type())
	}
	if !assertTypes(right, obj.IntType, obj.FloatType) {
		return obj.NewError("unsupported operator '>' for type %v", right.Type())
	}

	if assertTypes(left, obj.IntType) && assertTypes(right, obj.IntType) {
		l := left.(*obj.Integer).Val()
		r := right.(*obj.Integer).Val()
		return obj.ParseBool(l > r)
	}

	left, right = toFloat(left, right)
	l := left.(*obj.Float).Val()
	r := right.(*obj.Float).Val()
	return obj.ParseBool(l > r)
}

func (g Greater) String() string {
	return fmt.Sprintf("(%v > %v)", g.l, g.r)
}
