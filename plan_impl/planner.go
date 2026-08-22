package plan_impl

import (
	"fmt"

	"github.com/nikhilbhatia08/EuphoriaDB/parse"
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type Planner struct {
	queryPlanner QueryPlanner
	updatePlanner UpdatePlanner
}

func NewPlanner(queryPlanner QueryPlanner, updatePlanner UpdatePlanner) *Planner {
	return &Planner{
		queryPlanner: queryPlanner,
		updatePlanner: updatePlanner,
	}
}

func (pl *Planner) CreateQueryPlanner(cmd string, tx *transactions.Transaction) (plan.Plan, error) {
	parser := parse.NewParser(cmd)
	queryData, err := parser.Query()
	if err != nil {
		return nil, err
	}

	return pl.queryPlanner.CreatePlan(queryData, tx)
}

func (pl *Planner) ExecuteUpdate(cmd string, tx *transactions.Transaction) (int, error) {
	parser := parse.NewParser(cmd)
	data, err := parser.UpdateCmd()
	if err != nil {
		return 0, err
	}

	switch data := data.(type) {
	case *parse.InsertData:
		return pl.updatePlanner.ExecuteInsert(data, tx)
	case *parse.DeleteData:
		return pl.updatePlanner.ExecuteDelete(data, tx)
	case *parse.ModifyData:
		return pl.updatePlanner.ExecuteModify(data, tx)
	case *parse.CreateTableData:
		return pl.updatePlanner.ExecuteCreateTable(data, tx)
	case *parse.CreateViewData:
		return pl.updatePlanner.ExecuteCreateView(data, tx)
	default:
		return 0, fmt.Errorf("unexpected operation type %T", data)
	}
}