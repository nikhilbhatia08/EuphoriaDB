package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/metadata"
	"github.com/nikhilbhatia08/EuphoriaDB/parse"
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

var _ QueryPlanner = (*BasicQueryPlanner)(nil)

type BasicQueryPlanner struct {
	metadataManager *metadata.MetadataManager
}

func NewBasicQueryPlanner(metadataManager *metadata.MetadataManager) *BasicQueryPlanner {
	return &BasicQueryPlanner{
		metadataManager: metadataManager,
	}
}

func (bq *BasicQueryPlanner) CreatePlan(data *parse.QueryData, tx *transactions.Transaction) (plan.Plan, error) {
	plans := make([]plan.Plan, len(data.Tables()))
	for i, tableName := range data.Tables() {
		viewDef, err := bq.metadataManager.GetViewDefinition(tableName, tx)
		if err != nil {
			return nil, err
		}

		if viewDef != "" {
			parser := parse.NewParser(viewDef)
			viewData, err := parser.Query()
			if err != nil {
				return nil, err
			}

			viewPlan, err := bq.CreatePlan(viewData, tx)
			if err != nil {
				return nil, err
			}
			plans[i] = viewPlan
		} else {
			tablePlan, err := NewTablePlan(tx, tableName, bq.metadataManager)
			if err != nil {
				return nil, err
			}

			plans[i] = tablePlan
		}
	}

	currentPlan := plans[0]
	plans = plans[1:]

	for _, furtherPlan := range plans {
		selection1 := NewProductPlan(currentPlan, furtherPlan)
		selection2 := NewProductPlan(furtherPlan, currentPlan)

		if selection1.BlocksAccessed() < selection2.BlocksAccessed() {
			currentPlan = selection1
		} else {
			currentPlan = selection2
		}
	}

	currentPlan = NewSelectPlan(currentPlan, data.Pred())

	return NewProjectPlan(currentPlan, data.Fields()), nil
}
