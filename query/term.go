package query

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

type Term struct {
	lhs *Expression
	rhs *Expression
	op  types.Operator
}

func NewTerm(lhs *Expression, rhs *Expression, op types.Operator) *Term {
	return &Term{lhs: lhs, rhs: rhs, op: op}
}

func (t *Term) IsStatisfied(tableScan scan.TableScan) bool {
	lval, err := t.lhs.Evaluate(tableScan)
	if err != nil {
		return false
	}
	rval, err := t.rhs.Evaluate(tableScan)
	if err != nil {
		return false
	}

	switch t.op {
	case types.EQ:
		return lval == rval
	case types.NE:
		return lval != rval
	case types.LT, types.LE, types.GT, types.GE:
		return types.CompareSupportedTypes(lval, rval, t.op)
	default:
		return false
	}
}

func (t *Term) ReductionFactor(plan plan.Plan) int {
	if t.lhs.IsFieldName() && t.rhs.IsFieldName() {
		lhsName := t.lhs.AsFieldName()
		rhsName := t.rhs.AsFieldName()

		return max(plan.DistinctValues(lhsName), plan.DistinctValues(rhsName))
	}

	if t.lhs.IsFieldName() {
		return reductionForConstantComparison(plan.DistinctValues(t.lhs.AsFieldName()), t.op)
	}

	if t.rhs.IsFieldName() {
		return reductionForConstantComparison(plan.DistinctValues(t.rhs.AsFieldName()), t.op)
	}

	if t.lhs.AsConstant() == t.rhs.AsConstant() && t.op == types.EQ {
		return 1
	}
	if t.lhs.AsConstant() != t.rhs.AsConstant() && t.op == types.NE {
		return 1
	}

	return int(^uint(0) >> 1)
}

func reductionForConstantComparison(distinctValues int, op types.Operator) int {
	switch op {
	case types.EQ:
		return max(1, distinctValues)
	case types.NE:
		if distinctValues <= 1 {
			return 1
		} else {
			return distinctValues / (distinctValues - 1)
		}
	case types.LT, types.LE, types.GT, types.GE:
		return 2
	default:
		return 1
	}
}

func (t *Term) AppliesTo(schema *record.Schema) bool {
	return t.lhs.AppliesTo(schema) && t.rhs.AppliesTo(schema)
}

func (t *Term) EquatesWithConstant(fieldName string) any {
	if t.op != types.EQ {
		return nil
	}

	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fieldName && !t.rhs.IsFieldName() {
		return t.rhs.AsConstant()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fieldName && !t.lhs.IsFieldName() {
		return t.lhs.AsConstant()
	}

	return nil
}

func (t *Term) EquatesWithField(fieldName string) string {
	if t.op == types.EQ {
		return ""
	}

	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fieldName && t.rhs.IsFieldName() {
		return t.rhs.AsFieldName()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fieldName && t.lhs.IsFieldName() {
		return t.rhs.AsFieldName()
	}

	return ""
}

func (t *Term) ComparesWithConstant(fieldName string) (types.Operator, any) {
	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fieldName && !t.rhs.IsFieldName() {
		return t.op, t.rhs.AsConstant()
	}

	if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fieldName && !t.lhs.IsFieldName() {
		return t.op, t.lhs.AsConstant()
	}
	return types.NONE, nil
}

func (t *Term) String() string {
	return fmt.Sprintf("%s=%s", t.lhs.String(), t.rhs.String())
}
