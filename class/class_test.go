package class_test

import (
	"os"
	"testing"

	. "github.com/thara/godiva/class"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	f, err := os.Open("../testdata/HelloWorld.class")
	require.NoError(t, err)

	cf, err := Parse(f)
	require.NoError(t, err)

	assert.Equal(t, uint16(0), cf.MinorVer)
	assert.Equal(t, uint16(62), cf.MajorVer)
	assert.Equal(t, "62.0", cf.Version())
}
