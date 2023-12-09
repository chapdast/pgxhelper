package pgxhelper

import (
	"fmt"
	"strings"
)

func (c *Condition) valExtractor(values *[]any) (string, []any) {

	switch c.Operator {
	case OPR_ANY:
		*values = append(*values, c.Value)
		return "$%d =%s(%s)", []any{len(*values), c.Operator, c.ColumnName}
	case OPR_IN:
		*values = append(*values, c.Value)
		return "$%d %s %s", []any{len(*values), c.Operator, c.ColumnName}
	case OPR_IS_NULL:
		return "%s IS NULL", []any{c.ColumnName}
	case OPR_IS_NOT:
		return "%s IS NOT NULL", []any{c.ColumnName}

	case OPR_LIKE:
		*values = append(*values, fmt.Sprintf("%%%s%%", c.Value))
		return "%s::varchar %s $%d", []any{c.ColumnName, c.Operator, len(*values)}
	case OPR_NOT_LIKE:
		*values = append(*values, fmt.Sprintf("%%%s%%", c.Value))
		return "%s::varchar %s $%d", []any{c.ColumnName, c.Operator, len(*values)}
	default:
		*values = append(*values, c.Value)
		return "%s %s $%d", []any{c.ColumnName, c.Operator, len(*values)}
	}
}

func (dh *DatabaseHelper) BuildDynamicQuery(values *[]any, opts []*ConditionGroup) string {

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
			f, a := cond.valExtractor(values)
			localCond = append(localCond, fmt.Sprintf(f, a...))
		}
		// handle subgroups
		// grp := val.Group
		// for grp != nil {
		// 	subCond := make([]string, 0)

		// 	for _, gcond := range grp.Conditions {

		// 		f, a := gcond.valExtractor(values)
		// 		subCond = append(subCond, fmt.Sprintf(f, a...))

		// 	}
		// 	localCond = append(localCond, fmt.Sprintf("(%s)", strings.Join(subCond, grp.Join)))
		// 	grp = grp.Group
		// }

		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(localCond, val.Join)))
	}

	query := ""
	if len(conditions) != 0 {
		query = fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}
	return query
}

func (dh *DatabaseHelper) AddSort(sort *Sort) (query string) {
	if sort != nil {
		query += fmt.Sprintf("%s ORDER BY %s %s ", query, sort.Column, sort.Direction)
	}
	return query
}
