package fix

import (
	"io/ioutil"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func render(file packr.File, config PlushConfig) (string, error) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", errors.WithStack(err)
	}
	c := plush.NewContextWith(map[string]interface{}{
		"uuid": func() uuid.UUID {
			u, _ := uuid.NewV4()
			return u
		},
		"uuidNamed": uuidNamed,
		"now":       time.Now,
		"hash":      hash,
	})
	applyPlushConfig(config, c)


	return plush.Render(string(b), c)
}
type PlushConfig struct {
	TimeFormat string
}
func applyPlushConfig(config PlushConfig, context *plush.Context){
	if(config.TimeFormat != ""){
		context.Set("TIME_FORMAT", config.TimeFormat)
	}

}

func hash(s string, opts map[string]interface{}, help plush.HelperContext) (string, error) {
	cost := bcrypt.DefaultCost
	if i, ok := opts["cost"].(int); ok {
		cost = i
	}
	ph, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	return string(ph), err
}

func uuidNamed(name string, help plush.HelperContext) uuid.UUID {
	u, _ := uuid.NewV4()
	if ux, ok := help.Value(name).(uuid.UUID); ok {
		return ux
	}
	help.Set(name, u)
	return u
}
