package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/parse"
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type QueryPlanner interface {
	CreatePlan(queryData *parse.QueryData, transaction *transactions.Transaction) (plan.Plan, error)
}
