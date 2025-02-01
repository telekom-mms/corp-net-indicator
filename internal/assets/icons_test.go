package assets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telekom-mms/corp-net-indicator/internal/assets"
)

func TestGetIcon(t *testing.T) {
	assert := assert.New(t)

	// good cases
	assert.NotEmpty(assets.GetIcon(assets.ShieldOff))
	assert.NotEmpty(assets.GetIcon(assets.ShieldOn))
	assert.NotEmpty(assets.GetIcon(assets.Umbrella))
	assert.NotEmpty(assets.GetIcon(assets.Connect))
	assert.NotEmpty(assets.GetIcon(assets.Disconnect))
	assert.NotEmpty(assets.GetIcon(assets.Status))
	assert.NotEmpty(assets.GetIcon(assets.SVGWarning))
	assert.NotEmpty(assets.GetIcon(assets.SVGAlert))
	assert.NotEmpty(assets.GetIcon(assets.SVGMinus))
	assert.NotEmpty(assets.GetIcon(assets.SVGCheck))

	// bad case
	assert.Panics(func() {
		assets.GetIcon("not existing")
	})
}
