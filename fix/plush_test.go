package fix

import (
	"testing"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_hash(t *testing.T) {
	r := require.New(t)
	s, err := hash("password", map[string]interface{}{}, plush.HelperContext{})
	r.NoError(err)
	r.NotEqual("password", s)
}

func Test_nowAdd(t *testing.T) {
	offset := 1000
	r := require.New(t)
	tStr := nowAdd(offset)
	r.NotEmpty(tStr)
	exp := time.Now().Add(time.Second * time.Duration(offset))
	act, err := time.Parse(time.RFC3339, tStr)
	r.NoError(err)
	r.WithinDuration(exp, act, time.Second*10)
}

func Test_nowSub(t *testing.T) {
	offset := 1000
	r := require.New(t)
	tStr := nowSub(offset)
	r.NotEmpty(tStr)
	exp := time.Now().Add(time.Second * -time.Duration(offset))
	act, err := time.Parse(time.RFC3339, tStr)
	r.NoError(err)
	r.WithinDuration(exp, act, time.Second*10)
}
