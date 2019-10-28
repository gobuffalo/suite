package suite

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
	csrf "github.com/gobuffalo/mw-csrf"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/suite"
)

type Action struct {
	*Model
	Session *buffalo.Session
	App     *buffalo.App
	csrf    buffalo.MiddlewareFunc
}

// NewAction returns new Action for given buffalo.App
func NewAction(app *buffalo.App) *Action {
	as := &Action{
		App:   app,
		Model: NewModel(),
	}
	return as
}

func NewActionWithFixtures(app *buffalo.App, box packd.Box) (*Action, error) {
	m, err := NewModelWithFixtures(box)
	if err != nil {
		return nil, err
	}
	as := &Action{
		App:   app,
		Model: m,
	}
	return as, nil
}

func NewActionWithFixturesAndContext(app *buffalo.App, box packd.Box, ctx *plush.Context) (*Action, error) {
	m, err := NewModelWithFixturesAndContext(box, ctx)
	if err != nil {
		return nil, err
	}
	as := &Action{
		App:   app,
		Model: m,
	}
	return as, nil
}

func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

func (as *Action) HTML(u string, args ...interface{}) *httptest.Request {
	return httptest.New(as.App).HTML(u, args...)
}

func (as *Action) JSON(u string, args ...interface{}) *httptest.JSON {
	return httptest.New(as.App).JSON(u, args...)
}

func (as *Action) XML(u string, args ...interface{}) *httptest.XML {
	return httptest.New(as.App).XML(u, args...)
}

func (as *Action) SetupTest() {
	as.App.SessionStore = NewSessionStore()
	s, _ := as.App.SessionStore.New(nil, as.App.SessionName)
	as.Session = &buffalo.Session{
		Session: s,
	}

	if as.Model != nil {
		as.Model.SetupTest()
	}
	as.csrf = csrf.New
	csrf.New = func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			return next(c)
		}
	}
}

func (as *Action) TearDownTest() {
	csrf.New = as.csrf
	if as.Model != nil {
		as.Model.TearDownTest()
	}
}
