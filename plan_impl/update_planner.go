package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/parse"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
)

type UpdatePlanner interface {
	ExecuteInsert(data *parse.InsertData, transaction *transactions.Transaction) (int, error)

	ExecuteDelete(data *parse.DeleteData, transaction *transactions.Transaction) (int, error)

	ExecuteModify(data *parse.ModifyData, transaction *transactions.Transaction) (int, error)

	ExecuteCreateTable(data *parse.CreateTableData, transaction *transactions.Transaction) (int, error)

	ExecuteCreateView(data *parse.CreateViewData, transaction *transactions.Transaction) (int, error)
}