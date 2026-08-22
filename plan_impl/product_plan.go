package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/query"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

var _ plan.Plan = (*ProductPlan)(nil)

type ProductPlan struct {
	plan1  plan.Plan
	plan2  plan.Plan
	schema *record.Schema
}

func NewProductPlan(plan1 plan.Plan, plan2 plan.Plan) *ProductPlan {
	productPlan := &ProductPlan{
		plan1:  plan1,
		plan2:  plan2,
		schema: record.NewSchema(),
	}

	productPlan.schema.AddAll(plan1.Schema())
	productPlan.schema.AddAll(plan2.Schema())

	return productPlan
}

func (pp *ProductPlan) Open() (scan.TableScan, error) {
	scan1, err := pp.plan1.Open()
	if err != nil {
		return nil, err
	}
	scan2, err := pp.plan2.Open()
	if err != nil {
		return nil, err
	}

	return query.NewProductScan(scan1, scan2), nil
}

func (pp *ProductPlan) BlocksAccessed() int {
	return pp.plan1.BlocksAccessed() + (pp.plan1.RecordsOutput() * pp.plan2.BlocksAccessed())
}

func (pp *ProductPlan) RecordsOutput() int {
	return pp.plan1.RecordsOutput() * pp.plan2.RecordsOutput()
}

func (pp *ProductPlan) DistinctValues(fieldName string) int {
	if pp.plan1.Schema().HasField(fieldName) {
		return pp.plan1.DistinctValues(fieldName)
	}
	return pp.plan2.DistinctValues(fieldName)
}

func (pp *ProductPlan) Schema() *record.Schema {
	return pp.schema
}
