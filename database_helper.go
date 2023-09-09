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
	OPR_EQUAL = "="
	OPR_IN    = "IN"
	OPR_ANY   = "ANY"
	OPR_GTE   = ">="
)
const (
	JOINER_AND Joiner = " AND "
	JOINER_OR  Joiner = " OR "
)

type Joiner = string

type ConditionGroup struct {
	Join       Joiner
	Conditions []*Condition
	Group      *ConditionGroup
}

type Condition struct {
	ColumnName string
	Operator   string
	Value      any
}

var (
	ErrDuplicateKey = fmt.Errorf("DUPLICATE KEY")
	ErrNotFound = fmt.Errorf("NOT FOUND")
)




