package fix

import (
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

var scenes = sync.Map{}

func Init(box packr.Box) error {
	err := box.Walk(func(path string, file packr.File) error {
		if filepath.Ext(path) != ".toml" {
			return nil
		}

		x, err := render(file)

		sc := Scenarios{}
		_, err = toml.Decode(x, &sc)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, s := range sc.Scenarios {
			scenes.Store(s.Name, s)
		}
		return nil
	})
	return err
}

func Find(name ...string) (Scenario, error) {
	s, ok := scenes.Load(name)
	if !ok {
		return Scenario{}, errors.Errorf("could not find a scenario named %s", s)
	}
	sc, ok := s.(Scenario)
	if !ok {
		return Scenario{}, errors.Errorf("try to load %s but it isn't a Scenario it's a %T", s)
	}
	return sc, nil
}
