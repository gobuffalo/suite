package suite

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/markbates/pop"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ActionSuite struct {
	suite.Suite
	*require.Assertions
	DB     *pop.Connection
	Willie *willie.Willie
	App    *buffalo.App
}

func New(app *buffalo.App) *ActionSuite {
	as := &ActionSuite{
		App: app,
	}
	c, err := pop.Connect(envy.Get("GO_ENV", "test"))
	if err == nil {
		as.DB = c
	}
	return as
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

func (as *ActionSuite) HTML(u string, args ...interface{}) *willie.Request {
	return as.Willie.Request(u, args...)
}

func (as *ActionSuite) JSON(u string, args ...interface{}) *willie.JSON {
	return as.Willie.JSON(u, args...)
}

func (as *ActionSuite) SetupTest() {
	as.DB.MigrateReset("../migrations")
	as.Assertions = require.New(as.T())
	as.Willie = willie.New(as.App)
}

func (as *ActionSuite) TearDownTest() {}
