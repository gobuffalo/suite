package suite

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/plush/v4"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/suite/v3/fix"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Model suite
type Model struct {
	suite.Suite
	*require.Assertions
	DB       *pop.Connection
	Fixtures Box
}

// SetupTest clears database
func (m *Model) SetupTest() {
	m.Assertions = require.New(m.T())
	if m.DB != nil {
		err := m.DB.TruncateAll()
		m.NoError(err)
	}
}

// TearDownTest will be called after tests finish
func (m *Model) TearDownTest() {}

// DBDelta checks database table count change for a passed table name.
func (m *Model) DBDelta(delta int, name string, fn func()) {
	sc, err := m.DB.Count(name)
	m.NoError(err)
	fn()
	ec, err := m.DB.Count(name)
	m.NoError(err)
	m.Equal(sc+delta, ec)
}

// LoadFixture loads a named fixture into the database.
func (m *Model) LoadFixture(name string) {
	sc, err := fix.Find(name)
	m.NoError(err)
	db := m.DB.Store

	for _, table := range sc.Tables {
		for _, row := range table.Row {
			q := "insert into " + table.Name
			keys := []string{}
			skeys := []string{}
			for k := range row {
				keys = append(keys, k)
				skeys = append(skeys, ":"+k)
			}

			q = q + fmt.Sprintf(" (%s) values (%s)", strings.Join(keys, ","), strings.Join(skeys, ","))
			_, err = db.NamedExec(q, row)
			m.NoError(err)
		}
	}
}

// NewModel creates a new model suite
func NewModel() *Model {
	m := &Model{}
	c, err := pop.Connect(envy.Get("GO_ENV", "test"))
	if err == nil {
		m.DB = c
	}
	return m
}

// NewModelWithFixturesAndContext creates a new model suite with fixtures and a passed context.
func NewModelWithFixturesAndContext(box Box, ctx *plush.Context) (*Model, error) {
	m := NewModel()
	m.Fixtures = box
	return m, fix.InitWithContext(box, ctx)
}

// NewModelWithFixtures creates a new model with passed fixtures box
func NewModelWithFixtures(box Box) (*Model, error) {
	m := NewModel()
	m.Fixtures = box
	return m, fix.Init(box)
}
