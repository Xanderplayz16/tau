package ast

import (
	"fmt"

	"github.com/NicoNex/tau/obj"
)

type BitwiseXorAssign struct {
	l Node
	r Node
}

func NewBitwiseXorAssign(l, r Node) Node {
	return BitwiseXorAssign{l, r}
}

func (p BitwiseXorAssign) Eval(env *obj.Env) obj.Object {
	var name string
	var left = p.l.Eval(env)
	var right = p.r.Eval(env)

	if ident, ok := p.l.(Identifier); ok {
		name = ident.String()
	} else {
		return obj.NewError("cannot assign to literal")
	}

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	if !assertTypes(left, obj.INT) {
		return obj.NewError("unsupported operator '^=' for type %v", left.Type())
	}
	if !assertTypes(right, obj.INT) {
		return obj.NewError("unsupported operator '^=' for type %v", right.Type())
	}
	l := left.(*obj.Integer).Val()
	r := right.(*obj.Integer).Val()
	return env.Set(name, obj.NewInteger(l^r))
}

func (p BitwiseXorAssign) String() string {
	return fmt.Sprintf("(%v ^= %v)", p.l, p.r)
}
