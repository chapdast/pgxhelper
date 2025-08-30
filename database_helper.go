package pgxhelper

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseHelper struct {
	Pool *pgxpool.Pool
}

type Scanner = func(row pgx.Row) error

const (
	OPR_EQUAL     = "="
	OPR_NOT_EQUAL = "!="
	OPR_IN        = "IN"
	OPR_ANY       = "ANY"
	OPR_GTE       = ">="
	OPR_GT        = ">"
	OPR_LTE       = "<="
	OPR_LT        = "<"
	OPR_LIKE      = "ILIKE"
	OPR_NOT_LIKE  = "NOT ILIKE"
	OPR_IS        = "IS"
	OPR_IS_NOT    = "IS NOT"
	OPR_IS_NULL   = "_NULL_"
)

const (
	JOINER_AND Joiner = " AND "
	JOINER_OR  Joiner = " OR "
)

type Joiner = string

type ConditionGroup struct {
	Join       Joiner
	Conditions []*Condition
	// Group      *ConditionGroup
}

type Condition struct {
	ColumnName string
	Operator   string
	Value      any
}

var (
	ErrDuplicateKey = fmt.Errorf("DUPLICATE KEY")
	ErrNotFound     = fmt.Errorf("NOT FOUND")
)
