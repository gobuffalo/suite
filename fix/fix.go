package fix

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/plush/v4"
)

var scenes = map[string]Scenario{}
var moot = &sync.RWMutex{}

func InitWithContext(fsys fs.FS, ctx *plush.Context) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".toml" {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}

		x, err := renderWithContext(f, ctx)
		if err != nil {
			return err
		}

		sc := Scenarios{}
		_, err = toml.Decode(x, &sc)
		if err != nil {
			return err
		}

		moot.Lock()
		for _, s := range sc.Scenarios {
			scenes[s.Name] = s
		}
		moot.Unlock()
		return nil
	})
	return err
}

func Init(fsys fs.FS) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".toml" {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}

		x, err := render(f)
		if err != nil {
			return err
		}

		sc := Scenarios{}
		_, err = toml.Decode(x, &sc)
		if err != nil {
			return err
		}

		moot.Lock()
		for _, s := range sc.Scenarios {
			scenes[s.Name] = s
		}
		moot.Unlock()
		return nil
	})
	return err
}

func Find(name string) (Scenario, error) {
	moot.RLock()
	s, ok := scenes[name]
	moot.RUnlock()
	if !ok {
		return Scenario{}, fmt.Errorf("could not find a scenario named %q", name)
	}
	return s, nil
}
