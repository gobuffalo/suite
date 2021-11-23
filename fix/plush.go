package fix

import (
	"io"
	"time"

	"github.com/gobuffalo/plush/v4"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func renderWithContext(r io.Reader, ctx *plush.Context) (string, error) {
	cm := map[string]interface{}{
		"uuid": func() uuid.UUID {
			u, _ := uuid.NewV4()
			return u
		},
		"uuidNamed": uuidNamed,
		"now":       now,
		"hash":      hash,
		"nowAdd":    nowAdd,
		"nowSub":    nowSub,
	}
	for k, v := range cm {
		if !ctx.Has(k) {
			ctx.Set(k, v)
		}
	}
	return plush.RenderR(r, ctx)
}

func render(r io.Reader) (string, error) {
	ctx := plush.NewContextWith(map[string]interface{}{
		"uuid": func() uuid.UUID {
			u, _ := uuid.NewV4()
			return u
		},
		"uuidNamed": uuidNamed,
		"now":       now,
		"hash":      hash,
		"nowAdd":    nowAdd,
		"nowSub":    nowSub,
	})

	return renderWithContext(r, ctx)
}

func hash(s string, opts map[string]interface{}, help plush.HelperContext) (string, error) {
	cost := bcrypt.DefaultCost
	if i, ok := opts["cost"].(int); ok {
		cost = i
	}
	ph, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	return string(ph), err
}

func now(help plush.HelperContext) string {
	return time.Now().Format(time.RFC3339)
}

func nowAdd(s int) string {
	return time.Now().Add(time.Second * time.Duration(s)).Format(time.RFC3339)
}

func nowSub(s int) string {
	return time.Now().Add(time.Second * -time.Duration(s)).Format(time.RFC3339)
}

func uuidNamed(name string, help plush.HelperContext) uuid.UUID {
	u, _ := uuid.NewV4()
	if ux, ok := help.Value(name).(uuid.UUID); ok {
		return ux
	}
	help.Set(name, u)
	return u
}
