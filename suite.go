package suite

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
	csrf "github.com/gobuffalo/mw-csrf"
	"github.com/stretchr/testify/suite"
)

// Action suite
type Action struct {
	*Model
	Session *buffalo.Session
	App     *buffalo.App
	csrf    buffalo.MiddlewareFunc
}

// HTML creates an httptest.Request with HTML content type.
func (as *Action) HTML(u string, args ...interface{}) *httptest.Request {
	return httptest.New(as.App).HTML(u, args...)
}

// JSON creates an httptest.JSON request
func (as *Action) JSON(u string, args ...interface{}) *httptest.JSON {
	return httptest.New(as.App).JSON(u, args...)
}

// XML creates an httptest.XML request
func (as *Action) XML(u string, args ...interface{}) *httptest.XML {
	return httptest.New(as.App).XML(u, args...)
}

// SetupTest sets the session store, CSRF and clears database
func (as *Action) SetupTest() {
	as.App.SessionStore = newSessionStore()
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

// TearDownTest resets csrf
func (as *Action) TearDownTest() {
	csrf.New = as.csrf
	if as.Model != nil {
		as.Model.TearDownTest()
	}
}

// NewAction returns new Action for given buffalo.App
func NewAction(app *buffalo.App) *Action {
	as := &Action{
		App:   app,
		Model: NewModel(),
	}
	return as
}

// NewActionWithFixtures creates a new ActionSuite with passed box for fixtures.
func NewActionWithFixtures(app *buffalo.App, box Box) (*Action, error) {
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

//Run the passed suite
func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
