package query

import (
	"fmt"
	// "math"

	// "github.com/nikhilbhatia08/EuphoriaDB/plan"
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

func (t *Term) IsStatisfied(tableScan scan.TableScan) (bool, error) {
	lval, err := t.lhs.Evaluate(tableScan)
	if err != nil {
		return false, fmt.Errorf("error evaluating: %w", err)
	}
	rval, err := t.rhs.Evaluate(tableScan)
	if err != nil {
		return false, fmt.Errorf("error evaluating: %w", err)
	}

	if lval == rval {
		return true, nil
	}

	return false, nil
}

// func (t *Term) ReductionFactor(plan plan.Plan) int {
// 	if t.lhs.IsFieldName() && t.rhs.IsFieldName() {
// 		lhsName := t.lhs.AsFieldName()
// 		rhsName := t.rhs.AsFieldName()

// 		return math.MaxInt(plan)
// 	}
// }

func (t *Term) AppliesTo(schema *record.Schema) bool {
	return t.lhs.AppliesTo(schema) && t.rhs.AppliesTo(schema)
}

func (t *Term) EquatesWithConstant(fieldName string) any {
	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fieldName && !t.rhs.IsFieldName() {
		return t.rhs.AsConstant()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fieldName && !t.lhs.IsFieldName() {
		return t.lhs.AsConstant()
	}

	return nil
}

func (t *Term) EquatesWithField(fieldName string) string {
	if t.lhs.IsFieldName() && t.lhs.AsFieldName() == fieldName && t.rhs.IsFieldName() {
		return t.rhs.AsFieldName()
	} else if t.rhs.IsFieldName() && t.rhs.AsFieldName() == fieldName && t.lhs.IsFieldName() {
		return t.rhs.AsFieldName()
	}

	return ""
}

func (t *Term) String() string {
	return fmt.Sprintf("%s=%s", t.lhs.String(), t.rhs.String())
}
