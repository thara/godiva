package class

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//TODO
type fieldInfo byte
type methodInfo byte
type attributeInfo byte

// ClassFile
// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.1
type ClassFile struct {
	MinorVer, MajorVer uint16

	constantPoolCount uint16
	ConstantPool      []cpInfo

	AccessFlags     AccessFlags
	thisClass       uint16
	superClass      uint16
	interfaceCount  uint16
	interfaces      []*constantClass
	fieldsCount     uint16
	fields          []fieldInfo
	methodsCount    uint16
	methods         []methodInfo
	attributesCount uint16
	attributes      []attributeInfo
}

func Parse(r io.Reader) (*ClassFile, error) {
	var cf ClassFile

	var magic [4]byte
	if n, err := r.Read(magic[:]); err != nil {
		return nil, fmt.Errorf("%w", err)
	} else if n == 0 {
		return nil, errors.New("fail to parse magic number")
	} else if magic != [4]byte{0xCA, 0xFE, 0xBA, 0xBE} {
		return nil, errors.New("invalid magic number")
	}

	if err := binary.Read(r, binary.BigEndian, &cf.MinorVer); err != nil {
		return nil, fmt.Errorf("fail to parse minor_version: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &cf.MajorVer); err != nil {
		return nil, fmt.Errorf("fail to parse major_version: %w", err)
	}

	if err := binary.Read(r, binary.BigEndian, &cf.constantPoolCount); err != nil {
		return nil, fmt.Errorf("fail to parse constant_pool_count: %w", err)
	}

	cf.ConstantPool = make([]cpInfo, cf.constantPoolCount-1)
	for i := 0; i < int(cf.constantPoolCount)-1; i++ {
		cpInfo, err := parseCpInfo(r)
		if err != nil {
			return nil, fmt.Errorf("fail to parse constant_pool(%d): %w", i, err)
		}
		cf.ConstantPool[i] = cpInfo
	}

	var accessFlag uint16
	if err := binary.Read(r, binary.BigEndian, &accessFlag); err != nil {
		return nil, fmt.Errorf("fail to parse access_flags: %w", err)
	}
	cf.AccessFlags = AccessFlags(accessFlag)

	if err := binary.Read(r, binary.BigEndian, &cf.thisClass); err != nil {
		return nil, fmt.Errorf("fail to parse thisClass: %w", err)
	}
	if thisClass, ok := cf.lookupConstantPool(cf.thisClass); !ok {
		return nil, fmt.Errorf("`thisClass`(%d) must be a valid index in constant_pool", cf.thisClass)
	} else if _, ok := thisClass.(*constantClass); !ok {
		return nil, fmt.Errorf("The constant_pool entry at `thisClass`(%d) must be a CONSTANT_Class_info structure", cf.thisClass)
	}

	if err := binary.Read(r, binary.BigEndian, &cf.superClass); err != nil {
		return nil, fmt.Errorf("fail to parse superClass: %w", err)
	}
	if cf.superClass == 0 {
		//TODO validate whether this class file represents java.lang.Object
	} else {
		if superClass, ok := cf.lookupConstantPool(cf.superClass); !ok {
			return nil, fmt.Errorf("`superClass`(%d) must be a valid index in constant_pool", cf.superClass)
		} else if _, ok := superClass.(*constantClass); !ok {
			return nil, fmt.Errorf("The constant_pool entry at `superClass`(%d) must be a CONSTANT_Class_info structure", cf.superClass)
		}
	}

	if err := binary.Read(r, binary.BigEndian, &cf.interfaceCount); err != nil {
		return nil, fmt.Errorf("fail to parse interfaceCount: %w", err)
	}
	if cf.interfaceCount != 0 {
		cf.interfaces = make([]*constantClass, cf.interfaceCount-1)
		for i := 0; i < int(cf.interfaceCount)-1; i++ {
			var interfaceIdx uint16
			if err := binary.Read(r, binary.BigEndian, &interfaceIdx); err != nil {
				return nil, fmt.Errorf("fail to parse interfaces[%d]: %w", i, err)
			}
			if entry, ok := cf.lookupConstantPool(interfaceIdx); !ok {
				return nil, fmt.Errorf("`interfaces[%d]`(%d) must be a valid index in constant_pool", i, interfaceIdx)
			} else if class, ok := entry.(*constantClass); !ok {
				return nil, fmt.Errorf("The constant_pool entry at `interfaces[%d]`(%d) must be a CONSTANT_Class_info structure", i, interfaceIdx)
			} else {
				cf.interfaces[i] = class
			}
		}
	}

	return &cf, nil
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
	for i, class := range c.interfaces {
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

func must[T any](v T, ok bool) T {
	if !ok {
		panic("must be true")
	}
	return v
}
