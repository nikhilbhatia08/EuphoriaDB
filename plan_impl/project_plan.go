package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/query"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
)

var _ plan.Plan = (*ProjectScan)(nil)

type ProjectScan struct {
	plan   plan.Plan
	schema *record.Schema
}

func NewProjectScan(plan plan.Plan, fieldList []string) *ProjectScan {
	projectScan := &ProjectScan{
		plan:   plan,
		schema: record.NewSchema(),
	}
	for _, field := range fieldList {
		projectScan.schema.Add(field, plan.Schema())
	}

	return projectScan
}

func (ps *ProjectScan) Open() (scan.TableScan, error) {
	scan, err := ps.plan.Open()
	if err != nil {
		return nil, err
	}

	return query.NewProjectScan(scan, ps.plan.Schema().Fields()), nil
}

func (ps *ProjectScan) BlocksAccessed() int {
	return ps.plan.BlocksAccessed()
}

func (ps *ProjectScan) RecordsOutput() int {
	return ps.plan.RecordsOutput()
}

func (ps *ProjectScan) DistinctValues(fieldName string) int {
	return ps.plan.DistinctValues(fieldName)
}

func (ps *ProjectScan) Schema() *record.Schema {
	return ps.schema
}
