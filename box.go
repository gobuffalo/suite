package suite

import (
	"github.com/gobuffalo/packd"
)

// Box is Finder + Walkable
type Box interface {
	packd.Finder
	packd.Walkable
}
