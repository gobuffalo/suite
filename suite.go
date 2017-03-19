package suite

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/suite"
)

type Action struct {
	*Model
	Willie *willie.Willie
	App    *buffalo.App
}

func NewAction(app *buffalo.App) *Action {
	as := &Action{
		App:   app,
		Model: NewModel(),
	}
	return as
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

func (as *Action) HTML(u string, args ...interface{}) *willie.Request {
	return as.Willie.Request(u, args...)
}

func (as *Action) JSON(u string, args ...interface{}) *willie.JSON {
	return as.Willie.JSON(u, args...)
}

func (as *Action) SetupTest() {
	as.Model.SetupTest()
	as.Willie = willie.New(as.App)
}

func (as *Action) TearDownTest() {
	as.Model.TearDownTest()
}
