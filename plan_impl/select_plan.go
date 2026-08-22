package plan_impl

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/query"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

var _ plan.Plan = (*SelectPlan)(nil)

type SelectPlan struct {
	plan      plan.Plan
	predicate *query.Predicate
}

func NewSelectPlan(plan plan.Plan, pred *query.Predicate) *SelectPlan {
	return &SelectPlan{
		plan:      plan,
		predicate: pred,
	}
}

func (sp *SelectPlan) Open() (scan.TableScan, error) {
	scan, err := sp.plan.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open scan: %w", err)
	}

	return query.NewSelectScan(scan, sp.predicate), nil
}

func (sp *SelectPlan) BlocksAccessed() int {
	return sp.plan.BlocksAccessed()
}

func (sp *SelectPlan) RecordsOutput() int {
	return sp.plan.RecordsOutput()
}

func (sp *SelectPlan) DistinctValues(fieldName string) int {
	if sp.predicate.EquatesWithConstant(fieldName) != nil {
		return 1
	}

	secondField := sp.predicate.EquatesWithField(fieldName)
	if secondField != "" {
		return min(sp.plan.DistinctValues(fieldName), sp.plan.DistinctValues(secondField))
	}

	op, _ := sp.predicate.ComparesWithConstant(fieldName)
	switch op {
	case types.LT, types.LE, types.GT, types.GE:
		return max(1, sp.plan.DistinctValues(fieldName)/2)

	case types.NE:
		distinct := sp.plan.DistinctValues(fieldName)
		if distinct > 1 {
			return distinct - 1
		}
		return 1

	default:
		return sp.plan.DistinctValues(fieldName)
	}
}

func (sp *SelectPlan) Schema() *record.Schema {
	return sp.plan.Schema()
}
