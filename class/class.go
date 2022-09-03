package class

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

//TODO
type cpInfo byte
type fieldInfo byte
type methodInfo byte
type attributeInfo byte

// ClassFile
// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.1
type ClassFile struct {
	MinorVer, MajorVer uint16

	constantPoolCount uint16
	constantPool      []cpInfo

	accessFlags     uint16
	thisClass       uint16
	superClass      uint16
	interfaceCount  uint16
	interfaces      uint16
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

	return &cf, nil
}

func (c *ClassFile) Version() string {
	return fmt.Sprintf("%d.%d", c.MajorVer, c.MinorVer)
}
