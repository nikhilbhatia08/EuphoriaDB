package parse

import "github.com/nikhilbhatia08/EuphoriaDB/query"

type OrderByItem struct {
	field      string
	descending bool
}

func (obi *OrderByItem) Field() string {
	return obi.field
}

type QueryData struct {
	fields    []string
	tables    []string
	predicate *query.Predicate
	groupBy   []string
	having    *query.Predicate
	orderBy   []OrderByItem
}

func NewQueryData(fields, tables []string, predicate *query.Predicate) *QueryData {
	return &QueryData{
		fields:    fields,
		tables:    tables,
		predicate: predicate,
	}
}

func (qd *QueryData) Fields() []string {
	return qd.fields
}

func (qd *QueryData) Tables() []string {
	return qd.tables
}

func (qd *QueryData) Pred() *query.Predicate {
	return qd.predicate
}

func (qd *QueryData) GroupBy() []string {
	return qd.groupBy
}

func (qd *QueryData) Having() *query.Predicate {
	return qd.having
}

func (qd *QueryData) OrderBy() []OrderByItem {
	return qd.orderBy
}

func (qd *QueryData) String() string {
	if len(qd.fields) == 0 || len(qd.tables) == 0 {
		return ""
	}
	result := "select "
	for _, fieldName := range qd.fields {
		result += fieldName + ", "
	}
	// remove final comma/space
	if len(qd.fields) > 0 {
		result = result[:len(result)-2]
	}
	result += " from "
	for _, tableName := range qd.tables {
		result += tableName + ", "
	}
	if len(qd.tables) > 0 {
		result = result[:len(result)-2]
	}
	predicateString := qd.predicate.String()
	if predicateString != "" {
		result += " where " + predicateString
	}
	return result
}
