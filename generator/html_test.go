package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTML_Ext(t *testing.T) {
	t.Run("extension matches .html", func(t *testing.T) {
		h := HTML{}
		assert.Equal(t, ".html", h.Ext())
	})
}

func TestHTML_Generate(t *testing.T) {
	// Currently not easily testable
	// And also we do not want to test Goldmark.
	t.Skip()
}
