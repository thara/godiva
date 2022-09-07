package class

import (
	"errors"
	"fmt"
	"io"
)

//TODO
type methodInfo byte
type attributeInfo byte

// ClassFile
// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.1
type ClassFile struct {
	magic              [4]byte
	MinorVer, MajorVer uint16

	constantPoolCount uint16
	ConstantPool      []cpInfo

	AccessFlags     AccessFlags
	thisClass       uint16
	superClass      uint16
	interfaceCount  uint16
	interfaces      []uint16
	fieldsCount     uint16
	fields          []*fieldInfo
	methodsCount    uint16
	methods         []methodInfo
	attributesCount uint16
	attributes      []attributeInfo
}

func Parse(r io.Reader) (*ClassFile, error) {
	er := errReader{r: r}

	var magic [4]byte
	item(&er, "maigc number", bytes(magic[:], match([]byte{0xCA, 0xFE, 0xBA, 0xBE})))

	var cf ClassFile
	item(&er, "minor_version", integer(&cf.MinorVer))
	item(&er, "major_version", integer(&cf.MajorVer))

	if item(&er, "constant_pool_count", integer(&cf.constantPoolCount)) {
		cf.ConstantPool = make([]cpInfo, cf.constantPoolCount-1)
		item(&er, "constant_pool", entries(cf.ConstantPool[:], parseCpInfo))
	}

	var accessFlag uint16
	if item(&er, "access_flags", integer(&accessFlag)) {
		cf.AccessFlags = AccessFlags(accessFlag)
	}

	item(&er, "thisClass", integer(&cf.thisClass, constantPoolStructure[uint16, *constantClass](&cf)))
	if item(&er, "superClass", integer(&cf.superClass)) {
		if cf.superClass != 0 {
			validate(&er, constantPoolStructure[uint16, *constantClass](&cf))
		}
	}

	if item(&er, "interfaceCount", integer(&cf.interfaceCount)) {
		if 0 < cf.interfaceCount {
			cf.interfaces = make([]uint16, cf.interfaceCount-1)
			item(&er, "interfaces", entries(cf.interfaces[:], func(er *errReader) uint16 {
				var idx uint16
				item(er, "entry", integer(&idx, constantPoolStructure[uint16, *constantClass](&cf)))
				return idx
			}))
		}
	}

	if item(&er, "fieldsCount", integer(&cf.fieldsCount)) {
		if 0 < cf.fieldsCount {
		}
	}

	return &cf, er.err
}

func (c *ClassFile) Version() string {
	return fmt.Sprintf("%d.%d", c.MajorVer, c.MinorVer)
}

func (c *ClassFile) ThisClassName() string {
	class := getCpinfo[*constantClass](c, c.thisClass)
	utf8 := getCpinfo[*constantUtf8](c, class.nameIndex)
	return utf8.String()
}

func (c *ClassFile) SuperClassName() string {
	class := getCpinfo[*constantClass](c, c.superClass)
	utf8 := getCpinfo[*constantUtf8](c, class.nameIndex)
	return utf8.String()
}

func (c *ClassFile) InterfaceNames() []string {
	if c.interfaceCount == 0 {
		return nil
	}
	names := make([]string, c.interfaceCount-1)
	for i, idx := range c.interfaces {
		class := getCpinfo[*constantClass](c, idx)
		utf8 := getCpinfo[*constantUtf8](c, class.nameIndex)
		names[i] = utf8.String()
	}
	return names
}

func (c *ClassFile) lookupConstantPool(i uint16) (cpInfo, bool) {
	// The constant_pool table is indexed from 1 to constant_pool_count - 1
	if i < 1 {
		return nil, false
	} else if c.constantPoolCount < i {
		return nil, false
	}
	return c.ConstantPool[i-1], true
}

func getCpinfo[T cpInfo](cf *ClassFile, i uint16) T {
	e := must(cf.lookupConstantPool(i))
	return e.(T)
}

var (
	errNotFoundConstantPoolEntry    = errors.New("not found constant pool entry")
	errInvalidConstantPoolStructure = errors.New("invalid constant pool entry's structure")
)

func lookupCpinfo[T cpInfo](cf *ClassFile, i uint16) (entry T, err error) {
	e, ok := cf.lookupConstantPool(i)
	if !ok {
		return entry, errNotFoundConstantPoolEntry
	}
	entry, ok = e.(T)
	if !ok {
		return entry, errInvalidConstantPoolStructure
	}
	return entry, nil
}

func must[T any](v T, ok bool) T {
	if !ok {
		panic("must be true")
	}
	return v
}
