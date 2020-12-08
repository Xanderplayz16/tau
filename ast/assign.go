package ast

import (
	"fmt"
	"github.com/NicoNex/tau/obj"
)

type Assign struct {
	l Node
	r Node
}

func NewAssign(l, r Node) Node {
	return Assign{l, r}
}

func (a Assign) Eval(env *obj.Env) obj.Object {
	if left, ok := a.l.(Identifier); ok {
		env.Set(left.String(), a.r.Eval(env))
		return obj.NullObj
	}
	return obj.NewError("cannot assign to literal")
}

func (a Assign) String() string {
	return fmt.Sprintf("(%v = %v)", a.l, a.r)
}
