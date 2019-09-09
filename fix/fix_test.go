package fix

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_Init_And_Find(t *testing.T) {
	r := require.New(t)

	box := packr.Folder("./init-fixtures")
	r.NoError(Init(box))

	s, err := Find("lots of widgets")
	r.NoError(err)
	r.Equal("lots of widgets", s.Name)

	r.Len(s.Tables, 2)

	table := s.Tables[0]
	r.Equal("widgets", table.Name)
	r.Len(table.Row, 3)

	row := table.Row[0]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.Equal("This is widget #1", row["name"])
	r.Equal("some widget body", row["body"])

	wid := row["id"]

	row = table.Row[1]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.Equal("This is widget #2", row["name"])
	r.Equal("some widget body", row["body"])

	row = table.Row[2]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.Equal("This is widget #3", row["name"])
	r.Equal("some widget body", row["body"])

	table = s.Tables[1]
	r.Equal("users", table.Name)
	r.Len(table.Row, 1)

	row = table.Row[0]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.True(row["admin"].(bool))
	r.Equal(19.99, row["price"].(float64))
	r.Equal(wid, row["widget_id"])
}

func Test_InitWithContext_And_Find_CustomConfig(t *testing.T) {
	r := require.New(t)

	box := packr.Folder("./init-context-fixtures")
	ctx := plush.NewContextWith(map[string]interface{}{
		"double": func(num int, help plush.HelperContext) int {
			return num * 2
		},
	})
	r.NoError(InitWithContext(box, ctx))

	s, err := Find("widget with context")
	r.NoError(err)
	r.Equal("widget with context", s.Name)

	r.Len(s.Tables, 2)

	table := s.Tables[0]
	r.Equal("widgets", table.Name)
	r.Len(table.Row, 2)

	row := table.Row[0]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.Equal("This is widget #1", row["name"])
	r.Equal("some widget body", row["body"])

	wid := row["id"]

	row = table.Row[1]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.Equal("This is widget #2", row["name"])
	r.Equal("some widget body", row["body"])

	table = s.Tables[1]
	r.Equal("users", table.Name)
	r.Len(table.Row, 1)

	row = table.Row[0]
	r.NotZero(row["id"])
	r.NotZero(row["created_at"])
	r.NotZero(row["updated_at"])
	r.True(row["admin"].(bool))
	r.Equal(int64(36), row["price"].(int64))
	r.Equal(wid, row["widget_id"])
}
