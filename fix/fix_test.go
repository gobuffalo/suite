package fix

import (
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Init(t *testing.T) {
	r := require.New(t)

	box := packr.NewBox("./init-fixtures")

	r.NoError(Init(box))

	s, err := Find("lots of widgets")
	r.NoError(err)
	r.Equal("lots of widgets", s.Name)
}
