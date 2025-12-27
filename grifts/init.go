package grifts

import (
	"github.com/arxdsilva/hackathon/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
