package pgxhelper

import "testing"

func TestDatabaseHelper_BuildDynamicQuery(t *testing.T) {
	t.Run("with conditionGroup", func(t *testing.T) {

		p := &DatabaseHelper{}

		values := make([]any, 0)

		str := p.BuildDynamicQuery(&values, []*ConditionGroup{})
		t.Log(str)
		if len(str) != 0 {
			t.Fatalf("build unnessery query")
		}
		if len(values) != 0 {
			t.Fatal("values added to values list, this should be empty", values)
		}
	})

	t.Run("with conditionGroup", func(t *testing.T) {

		p := &DatabaseHelper{}

		values := make([]any, 0)
		str := p.BuildDynamicQuery(&values, []*ConditionGroup{
			{
				Join: JOINER_AND,
				Conditions: []*Condition{
					{
						ColumnName: "a1",
						Operator:   OPR_GTE,
						Value:      1,
					},
					{
						ColumnName: "a2",
						Operator:   OPR_EQUAL,
						Value:      4,
					},
				},
				// Group: nil,
			},
			{
				Join: JOINER_OR,
				Conditions: []*Condition{
					{
						ColumnName: "O1",
						Operator:   OPR_GTE,
						Value:      1,
					},
					{
						ColumnName: "O2",
						Operator:   OPR_EQUAL,
						Value:      4,
					},
				},
				// Group: &ConditionGroup{
				// 	Join: JOINER_AND,
				// 	Conditions: []*Condition{
				// 		{
				// 			ColumnName: "Oa1",
				// 			Operator:   OPR_GTE,
				// 			Value:      1,
				// 		},
				// 		{
				// 			ColumnName: "Oa2",
				// 			Operator:   OPR_EQUAL,
				// 			Value:      4,
				// 		},
				// 	},
				// },
			},
		})
		success := " WHERE (a1 >= $1 AND a2 = $2) AND (O1 >= $3 OR O2 = $4 OR (Oa1 >= $5 AND Oa2 = $6))"

		if str != success {
			t.Fatalf("not match \nW:%s\n G:%s\n", success, str)
		}

		if len(values) != 6 {
			t.Fatalf("values not added, %v", values)
		}
	})

}
