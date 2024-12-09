package mapper

import (
	"fmt"
	"go/token"

	"github.com/algoboyz/garm/pkg/op"
)

func (m *SSAMapper) MapToken(tok token.Token) (op.Op, error) {
	switch tok {
	case token.ADD:
		return op.ADD, nil
	case token.SUB:
		return op.SUB, nil
	case token.MUL:
		return op.MUL, nil
	case token.QUO:
		return op.SDIV, nil
	case token.REM:
		return op.CSEL, nil
		// todo needs to return full instruction for modulo to work
		// csel    X0, X3, X4, eq          // X0 = (X1==0) ? X3 : X4
	case token.AND:
		return op.AND, nil
	case token.OR:
		return op.OR, nil
	case token.XOR:
		return op.XOR, nil
	case token.SHL:
		return op.SHL, nil
	case token.SHR:
		return op.SHR, nil
	case token.AND_NOT:
		return op.BIC, nil
	default:
		return op.NOP, fmt.Errorf("unsupported token: %s", tok)
	}
}
