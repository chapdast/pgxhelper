package pgxhelper

import (
	"fmt"
	"strings"
)

func (c *Condition) valExtractor(values *[]any) (string, any, any, any) {

	switch c.Operator {
	case OPR_ANY:
		*values = append(*values, c.Value)
		return "$%d =%s(%s)", len(*values), c.Operator, c.ColumnName
	case OPR_IN:
		*values = append(*values, c.Value)
		return "$%d %s %s", len(*values), c.Operator, c.ColumnName
	default:
		*values = append(*values, c.Value)
		return "%s %s $%d", c.ColumnName, c.Operator, len(*values)
	}
}

func (dh DatabaseHelper) BuildDynamicQuery(values *[]any, opts []*ConditionGroup) string {

	if len(opts) == 0 {
		return ""
	}

	if values == nil {
		values = new([]any)
	}

	var (
		conditions = make([]string, 0)
	)

	for _, val := range opts {
		if len(val.Conditions) == 0 {
			continue
		}
		localCond := make([]string, 0)
		for _, cond := range val.Conditions {
			f, a, b, c := cond.valExtractor(values)
			localCond = append(localCond, fmt.Sprintf(f, a, b, c))
		}
		// handle subgroups
		grp := val.Group
		for grp != nil {
			subCond := make([]string, 0)

			for _, gcond := range grp.Conditions {
				f, a, b, c := gcond.valExtractor(values)
				subCond = append(subCond, fmt.Sprintf(f, a, b, c))
			}
			localCond = append(localCond, fmt.Sprintf("(%s)", strings.Join(subCond, grp.Join)))
			grp = grp.Group
		}

		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(localCond, val.Join)))
	}

	if len(conditions) == 0 || len(*values) == 0 {
		return ""
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
}
