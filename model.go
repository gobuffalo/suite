package suite

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/plush/v4"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/suite/v4/fix"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Model suite
type Model struct {
	suite.Suite
	*require.Assertions
	DB       *pop.Connection
	Fixtures fs.FS
}

// SetupTest clears database
func (m *Model) SetupTest() {
	m.Assertions = require.New(m.T())
	if m.DB != nil {
		err := m.CleanDB()
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
func NewModelWithFixturesAndContext(fsys fs.FS, ctx *plush.Context) (*Model, error) {
	m := NewModel()
	m.Fixtures = fsys
	return m, fix.InitWithContext(fsys, ctx)
}

// NewModelWithFixtures creates a new model with passed fixtures box
func NewModelWithFixtures(fsys fs.FS) (*Model, error) {
	m := NewModel()
	m.Fixtures = fsys
	return m, fix.Init(fsys)
}

func (m *Model) Run(name string, subtest func()) bool {
	return m.Suite.Run(name, func() {
		m.Assertions = require.New(m.Suite.T())
		subtest()
	})
}

// CleanDB clears records from the database, this function is
// useful to run before tests to ensure other tests are not
// affecting the one running.
func (m *Model) CleanDB() error {
	if m.DB == nil {
		return nil
	}

	switch m.DB.Dialect.Name() {
	case "postgres":
		deleteAllQuery := `DO
		$func$
			DECLARE
			_tbl text;
			_sch text;
			BEGIN
				FOR _sch, _tbl IN
					SELECT schemaname, tablename
					FROM   pg_tables
					WHERE  tablename <> '%s' AND schemaname NOT IN ('pg_catalog', 'information_schema') AND tableowner = current_user
				LOOP
					EXECUTE format('ALTER TABLE %%I.%%I DISABLE TRIGGER ALL;', _sch, _tbl);
					EXECUTE format('DELETE FROM %%I.%%I CASCADE', _sch, _tbl);
					EXECUTE format('ALTER TABLE %%I.%%I ENABLE TRIGGER ALL;', _sch, _tbl);
				END LOOP;
			END
		$func$;`

		q := m.DB.RawQuery(fmt.Sprintf(deleteAllQuery, m.DB.MigrationTableName()))
		return q.Exec()
	default:
		return m.DB.TruncateAll()
	}
}
