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

	assert.Len(t, cf.ConstantPool, 28)

	if assert.EqualValues(t, ConstantKindMethodref, cf.ConstantPool[0].Tag()) {
		assert.Equal(t, "#2.#3", cf.ConstantPool[0].String())
	}
	if assert.EqualValues(t, ConstantKindClass, cf.ConstantPool[1].Tag()) {
		assert.Equal(t, "#4", cf.ConstantPool[1].String())
	}
	if assert.EqualValues(t, ConstantKindNameAndType, cf.ConstantPool[2].Tag()) {
		assert.Equal(t, "#5:#6", cf.ConstantPool[2].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[3].Tag()) {
		assert.Equal(t, "java/lang/Object", cf.ConstantPool[3].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[4].Tag()) {
		assert.Equal(t, "<init>", cf.ConstantPool[4].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[5].Tag()) {
		assert.Equal(t, "()V", cf.ConstantPool[5].String())
	}
	if assert.EqualValues(t, ConstantKindFieldref, cf.ConstantPool[6].Tag()) {
		assert.Equal(t, "#8.#9", cf.ConstantPool[6].String())
	}
	if assert.EqualValues(t, ConstantKindClass, cf.ConstantPool[7].Tag()) {
		assert.Equal(t, "#10", cf.ConstantPool[7].String())
	}
	if assert.EqualValues(t, ConstantKindNameAndType, cf.ConstantPool[8].Tag()) {
		assert.Equal(t, "#11:#12", cf.ConstantPool[8].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[9].Tag()) {
		assert.Equal(t, "java/lang/System", cf.ConstantPool[9].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[10].Tag()) {
		assert.Equal(t, "out", cf.ConstantPool[10].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[11].Tag()) {
		assert.Equal(t, "Ljava/io/PrintStream;", cf.ConstantPool[11].String())
	}
	if assert.EqualValues(t, ConstantKindString, cf.ConstantPool[12].Tag()) {
		assert.Equal(t, "#14", cf.ConstantPool[12].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[13].Tag()) {
		assert.Equal(t, "Hello, world", cf.ConstantPool[13].String())
	}
	if assert.EqualValues(t, ConstantKindMethodref, cf.ConstantPool[14].Tag()) {
		assert.Equal(t, "#16.#17", cf.ConstantPool[14].String())
	}
	if assert.EqualValues(t, ConstantKindClass, cf.ConstantPool[15].Tag()) {
		assert.Equal(t, "#18", cf.ConstantPool[15].String())
	}
	if assert.EqualValues(t, ConstantKindNameAndType, cf.ConstantPool[16].Tag()) {
		assert.Equal(t, "#19:#20", cf.ConstantPool[16].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[17].Tag()) {
		assert.Equal(t, "java/io/PrintStream", cf.ConstantPool[17].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[18].Tag()) {
		assert.Equal(t, "println", cf.ConstantPool[18].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[19].Tag()) {
		assert.Equal(t, "(Ljava/lang/String;)V", cf.ConstantPool[19].String())
	}
	if assert.EqualValues(t, ConstantKindClass, cf.ConstantPool[20].Tag()) {
		assert.Equal(t, "#22", cf.ConstantPool[20].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[21].Tag()) {
		assert.Equal(t, "HelloWorld", cf.ConstantPool[21].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[22].Tag()) {
		assert.Equal(t, "Code", cf.ConstantPool[22].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[23].Tag()) {
		assert.Equal(t, "LineNumberTable", cf.ConstantPool[23].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[24].Tag()) {
		assert.Equal(t, "main", cf.ConstantPool[24].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[25].Tag()) {
		assert.Equal(t, "([Ljava/lang/String;)V", cf.ConstantPool[25].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[26].Tag()) {
		assert.Equal(t, "SourceFile", cf.ConstantPool[26].String())
	}
	if assert.EqualValues(t, ConstantKindUtf8, cf.ConstantPool[27].Tag()) {
		assert.Equal(t, "HelloWorld.java", cf.ConstantPool[27].String())
	}

	assert.Equal(t, AccessFlagsPublic|AccessFlagsSuper, cf.AccessFlags)
}
