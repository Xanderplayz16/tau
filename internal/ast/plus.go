package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/obj"
)

type Plus struct {
	l   Node
	r   Node
	pos int
}

func NewPlus(l, r Node, pos int) Node {
	return Plus{
		l:   l,
		r:   r,
		pos: pos,
	}
}

func (p Plus) Eval(env *obj.Env) obj.Object {
	var (
		left  = obj.Unwrap(p.l.Eval(env))
		right = obj.Unwrap(p.r.Eval(env))
	)

	if takesPrecedence(left) {
		return left
	}
	if takesPrecedence(right) {
		return right
	}

	if !assertTypes(left, obj.IntType, obj.FloatType, obj.StringType) {
		return obj.NewError("unsupported operator '+' for type %v", left.Type())
	}
	if !assertTypes(right, obj.IntType, obj.FloatType, obj.StringType) {
		return obj.NewError("unsupported operator '+' for type %v", right.Type())
	}

	switch {
	case assertTypes(left, obj.StringType) && assertTypes(right, obj.StringType):
		l := left.(*obj.String).Val()
		r := right.(*obj.String).Val()
		return obj.NewString(l + r)

	case assertTypes(left, obj.IntType) && assertTypes(right, obj.IntType):
		l := left.(*obj.Integer).Val()
		r := right.(*obj.Integer).Val()
		return obj.NewInteger(l + r)

	case assertTypes(left, obj.FloatType, obj.IntType) && assertTypes(right, obj.FloatType, obj.IntType):
		left, right = toFloat(left, right)
		l := left.(*obj.Float).Val()
		r := right.(*obj.Float).Val()
		return obj.NewFloat(l + r)

	default:
		return obj.NewError(
			"invalid operation %v + %v (wrong types %v and %v)",
			left, right, left.Type(), right.Type(),
		)
	}
}

func (p Plus) String() string {
	return fmt.Sprintf("(%v + %v)", p.l, p.r)
}

func (p Plus) Compile(c *compiler.Compiler) (position int, err error) {
	if p.IsConstExpression() {
		position = c.Emit(code.OpConstant, c.AddConstant(p.Eval(nil)))
		c.Bookmark(p.pos)
		return
	}

	if position, err = p.l.Compile(c); err != nil {
		return
	}
	if position, err = p.r.Compile(c); err != nil {
		return
	}
	position = c.Emit(code.OpAdd)
	c.Bookmark(p.pos)
	return
}

func (p Plus) IsConstExpression() bool {
	return p.l.IsConstExpression() && p.r.IsConstExpression()
}
