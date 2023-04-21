package ast

import (
	"fmt"

	"github.com/NicoNex/tau/internal/code"
	"github.com/NicoNex/tau/internal/compiler"
	"github.com/NicoNex/tau/internal/vm/cvm/cobj"
)

type Mod struct {
	l   Node
	r   Node
	pos int
}

func NewMod(l, r Node, pos int) Node {
	return Mod{
		l:   l,
		r:   r,
		pos: pos,
	}
}

func (m Mod) Eval() (cobj.Object, error) {
	left, err := m.l.Eval()
	if err != nil {
		return cobj.NullObj, err
	}

	right, err := m.r.Eval()
	if err != nil {
		return cobj.NullObj, err
	}

	if !cobj.AssertTypes(left, cobj.IntType) {
		return cobj.NullObj, fmt.Errorf("unsupported operator '%%' for type %v", left.Type())
	}
	if !cobj.AssertTypes(right, cobj.IntType) {
		return cobj.NullObj, fmt.Errorf("unsupported operator '%%' for type %v", right.Type())
	}
	if right.Int() == 0 {
		return cobj.NullObj, fmt.Errorf("can't divide by 0")
	}

	return cobj.NewInteger(left.Int() % right.Int()), nil
}

func (m Mod) String() string {
	return fmt.Sprintf("(%v %% %v)", m.l, m.r)
}

func (m Mod) Compile(c *compiler.Compiler) (position int, err error) {
	if m.IsConstExpression() {
		o, err := m.Eval()
		if err != nil {
			return 0, c.NewError(m.pos, err.Error())
		}
		position = c.Emit(code.OpConstant, c.AddConstant(o))
		c.Bookmark(m.pos)
		return position, err
	}

	if position, err = m.l.Compile(c); err != nil {
		return
	}
	if position, err = m.r.Compile(c); err != nil {
		return
	}
	position = c.Emit(code.OpMod)
	c.Bookmark(m.pos)
	return
}

func (m Mod) IsConstExpression() bool {
	return m.l.IsConstExpression() && m.r.IsConstExpression()
}
