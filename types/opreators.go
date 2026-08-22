package types

import "fmt"

type Operator int

const (
	NONE Operator = -1

	EQ Operator = iota

	NE

	LT

	LE

	GT

	GE
)

func (op Operator) String() string {
	switch op {
	case EQ:
		return "="
	case NE:
		return "<>"
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	default:
		return ""
	}
}

func OperatorFromString(op string) (Operator, error) {
	switch op {
	case "=":
		return EQ, nil
	case "<>", "!=":
		return NE, nil
	case "<":
		return LT, nil
	case "<=":
		return LE, nil
	case ">":
		return GT, nil
	case ">=":
		return GE, nil
	default:
		return -1, fmt.Errorf("invalid operator: %s", op)
	}
}
