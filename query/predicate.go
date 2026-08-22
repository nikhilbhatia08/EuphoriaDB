package query

import (
	"github.com/nikhilbhatia08/EuphoriaDB/plan"
	"github.com/nikhilbhatia08/EuphoriaDB/record"
	"github.com/nikhilbhatia08/EuphoriaDB/scan"
	"github.com/nikhilbhatia08/EuphoriaDB/types"
)

type Predicate struct {
	terms []*Term
}

func NewPredicate() *Predicate {
	return &Predicate{}
}

func NewPredicateFromTerm(term *Term) *Predicate {
	terms := []*Term{}
	terms = append(terms, term)
	return &Predicate{terms: terms}
}

func (p *Predicate) ConjoinWith(other *Predicate) {
	p.terms = append(p.terms, other.terms...)
}

func (p *Predicate) IsStatisfied(tableScan scan.TableScan) (bool, error) {
	for _, term := range p.terms {
		condition := term.IsStatisfied(tableScan)
		if !condition {
			return false, nil
		}
	}

	return true, nil
}

func (p *Predicate) ReductionFactor(queryPlan plan.Plan) int {
	factor := 1
	for _, term := range p.terms {
		factor *= term.ReductionFactor(queryPlan)
	}
	return factor
}

func (p *Predicate) SelectSubPredicate(schema *record.Schema) *Predicate {
	resultPredicate := NewPredicate()

	for _, term := range p.terms {
		if term.AppliesTo(schema) {
			resultPredicate.terms = append(resultPredicate.terms, term)
		}
	}

	if len(resultPredicate.terms) == 0 {
		return nil
	}

	return resultPredicate
}

func (p *Predicate) JoinSubPredicate(schema1, schema2 *record.Schema) *Predicate {
	result := NewPredicate()
	unionSchema := record.NewSchema()

	unionSchema.AddAll(schema1)
	unionSchema.AddAll(schema2)

	for _, term := range p.terms {
		if !term.AppliesTo(schema1) && !term.AppliesTo(schema2) && term.AppliesTo(unionSchema) {
			result.terms = append(result.terms, term)
		}
	}
	if len(result.terms) == 0 {
		return nil
	}
	return result
}

func (p *Predicate) EquatesWithConstant(fieldName string) any {
	for _, term := range p.terms {
		constant := term.EquatesWithConstant(fieldName)
		if constant != nil {
			return constant
		}
	}

	return nil
}

func (p *Predicate) EquatesWithField(fieldName string) string {
	for _, term := range p.terms {
		field := term.EquatesWithField(fieldName)
		if field != "" {
			return field
		}
	}

	return ""
}

func (p *Predicate) ComparesWithConstant(fieldName string) (types.Operator, any) {
	for _, term := range p.terms {
		if op, c := term.ComparesWithConstant(fieldName); op != types.NONE {
			return op, c
		}
	}
	return types.NONE, nil
}

func (p *Predicate) String() string {
	if len(p.terms) == 0 {
		return ""
	}

	result := p.terms[0].String()
	for _, term := range p.terms[1:] {
		result += " and " + term.String()
	}

	return result
}
