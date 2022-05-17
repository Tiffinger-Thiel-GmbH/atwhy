package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHTML_Ext(t *testing.T) {
	t.Run("extension matches .html", func(t *testing.T) {
		h := HTML{}
		assert.Equal(t, ".html", h.Ext())
	})
}
