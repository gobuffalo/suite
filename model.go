package suite

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/plush"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/suite/fix"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Model struct {
	suite.Suite
	*require.Assertions
	DB       *pop.Connection
	Fixtures packd.Finder
}

func (m *Model) SetupTest() {
	m.Assertions = require.New(m.T())
	if m.DB != nil {
		err := m.DB.TruncateAll()
		m.NoError(err)
	}
}

func (m *Model) TearDownTest() {}

func (m *Model) DBDelta(delta int, name string, fn func()) {
	sc, err := m.DB.Count(name)
	m.NoError(err)
	fn()
	ec, err := m.DB.Count(name)
	m.NoError(err)
	m.Equal(sc+delta, ec)
}

func (as *Model) LoadFixture(name string) {
	sc, err := fix.Find(name)
	as.NoError(err)
	db := as.DB.Store

	for _, table := range sc.Tables {
		for _, row := range table.Row {
			q := fmt.Sprintf("insert into %s", as.DB.Dialect.Quote(table.Name))
			keys := []string{}
			skeys := []string{}
			for k := range row {
				keys = append(keys, as.DB.Dialect.Quote(k))
				skeys = append(skeys, ":"+k)
			}

			q = q + fmt.Sprintf(" (%s) values (%s)", strings.Join(keys, ","), strings.Join(skeys, ","))
			_, err = db.NamedExec(q, row)
			as.NoError(err)
		}
	}
}

func NewModel() *Model {
	m := &Model{}
	c, err := pop.Connect(envy.Get("GO_ENV", "test"))
	if err == nil {
		m.DB = c
	}
	return m
}

type Box interface {
	packd.Finder
	packd.Walkable
}

func NewModelWithFixturesAndContext(box packd.Box, ctx *plush.Context) (*Model, error) {
	m := NewModel()
	m.Fixtures = box
	return m, fix.InitWithContext(box, ctx)
}

func NewModelWithFixtures(box packd.Box) (*Model, error) {
	m := NewModel()
	m.Fixtures = box
	return m, fix.Init(box)
}
