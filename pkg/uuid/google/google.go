package google

import (
	uuidgen "github.com/adanyl0v/pocket-ideas/pkg/uuid"
	"github.com/google/uuid"
)

type generator struct{}

func New() uuidgen.Generator {
	return &generator{}
}

func (g *generator) NewV1() (string, error) {
	return g.gen(uuid.NewUUID)
}

func (g *generator) NewV4() (string, error) {
	return g.gen(uuid.NewRandom)
}

func (g *generator) NewV6() (string, error) {
	return g.gen(uuid.NewV6)
}

func (g *generator) NewV7() (string, error) {
	return g.gen(uuid.NewV7)
}

type genFn func() (uuid.UUID, error)

func (g *generator) gen(fn genFn) (string, error) {
	id, err := fn()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
