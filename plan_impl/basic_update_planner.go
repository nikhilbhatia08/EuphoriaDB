package plan_impl

import (
	"github.com/nikhilbhatia08/EuphoriaDB/metadata"
	"github.com/nikhilbhatia08/EuphoriaDB/parse"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/transactions"
	"github.com/nikhilbhatia08/EuphoriaDB/utils"
)

var _ UpdatePlanner = (*BasicUpdatePlanner)(nil)

type BasicUpdatePlanner struct {
	metadataManager *metadata.MetadataManager
}

func NewBasicUpdatePlanner(mm *metadata.MetadataManager) *BasicUpdatePlanner {
	return &BasicUpdatePlanner{metadataManager: mm}
}

func (bup *BasicUpdatePlanner) ExecuteDelete(data *parse.DeleteData, tx *transactions.Transaction) (int, error) {
	tablePlan, err := NewTablePlan(tx, data.TableName(), bup.metadataManager)
	if err != nil {
		return 0, err
	}

	selectPlan := NewSelectPlan(tablePlan, data.Predicate())

	s, err := selectPlan.Open()
	if err != nil {
		return 0, err
	}
	updateScan := s.(scan.UpdateScan)
	defer updateScan.Close()

	count := 0
	for {
		next, err := updateScan.Next()
		if err != nil {
			return count, err
		}
		if !next {
			break
		}

		if err := updateScan.Delete(); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (bup *BasicUpdatePlanner) ExecuteModify(data *parse.ModifyData, tx *transactions.Transaction) (int, error) {
	tablePlan, err := NewTablePlan(tx, data.TableName(), bup.metadataManager)
	if err != nil {
		return 0, nil
	}

	selectPlan := NewSelectPlan(tablePlan, data.Predicate())
	s, err := selectPlan.Open()
	if err != nil {
		return 0, err
	}

	updateScan := s.(scan.UpdateScan)
	count := 0
	for {
		next, err := updateScan.Next()
		if err != nil {
			return count, err
		}
		if !next {
			break
		}

		value, err := data.NewValue().Evaluate(updateScan)
		if err != nil {
			return count, err
		}
		if err := updateScan.SetVal(data.TargetField(), value); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (bup *BasicUpdatePlanner) ExecuteInsert(data *parse.InsertData, tx *transactions.Transaction) (int, error) {
	tablePlan, err := NewTablePlan(tx, data.TableName(), bup.metadataManager)
	if err != nil {
		return 0, err
	}

	s, err := tablePlan.Open()
	if err != nil {
		return 0, err
	}
	updateScan := s.(scan.UpdateScan)
	defer updateScan.Close()

	if err := updateScan.Insert(); err != nil {
		return 0, err
	}

	values := data.Values()
	for i, fieldName := range data.Fields() {
		value := values[i]
		if err := utils.ValidateStringLength(tablePlan.Schema(), fieldName, value); err != nil {
			return 0, err
		}

		if err := updateScan.SetVal(fieldName, value); err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (up *BasicUpdatePlanner) ExecuteCreateTable(data *parse.CreateTableData, transaction *transactions.Transaction) (int, error) {
	err := up.metadataManager.CreateTable(data.TableName(), data.NewSchema(), transaction)
	return 0, err
}

func (up *BasicUpdatePlanner) ExecuteCreateView(data *parse.CreateViewData, transaction *transactions.Transaction) (int, error) {
	err := up.metadataManager.CreateView(data.ViewName(), data.ViewDefinition(), transaction)
	return 0, err
}
