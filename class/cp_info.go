package class

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ConstantKind = byte

const (
	ConstantKindClass              ConstantKind = 7
	ConstantKindFieldref                        = 9
	ConstantKindMethodref                       = 10
	ConstantKindInterfaceMethodref              = 11
	ConstantKindString                          = 8
	ConstantKindInteger                         = 3
	ConstantKindFloat                           = 4
	ConstantKindLong                            = 5
	ConstantKindDouble                          = 6
	ConstantKindNameAndType                     = 12
	ConstantKindUtf8                            = 1
	ConstantKindMethodHandle                    = 15
	ConstantKindMethodType                      = 16
	ConstantKindDynamic                         = 17
	ConstantKindInvokeDynamic                   = 18
	ConstantKindModule                          = 19
	ConstantKindPackage                         = 20
)

type cpInfo interface {
	Tag() byte
	String() string
}

func parseCpInfo(r io.Reader) (cpInfo, error) {
	var b [1]byte
	if n, err := r.Read(b[:]); err != nil {
		return nil, fmt.Errorf("%w", err)
	} else if n != 1 {
		return nil, errors.New("fail to parse cp_info tag")
	}

	var p parser

	tag := cpInfoTag{tag: b[0]}
	switch tag.tag {
	case ConstantKindClass:
		c := constantClass{cpInfoTag: tag}
		p.readInteger(r, &c.nameIndex, "CONSTANT_Class_info's name_index")
		return &c, p.Err
	case ConstantKindFieldref:
		c := constantFieldref{cpInfoTag: tag}
		p.readInteger(r, &c.classIndex, "CONSTANT_Fieldref_info's class_index")
		p.readInteger(r, &c.nameAndTypeIndex, "CONSTANT_Fieldref_info's name_and_type_index")
		return &c, p.Err
	case ConstantKindMethodref:
		c := constantMethodref{cpInfoTag: tag}
		p.readInteger(r, &c.classIndex, "CONSTANT_Methodref_info's class_index")
		p.readInteger(r, &c.nameAndTypeIndex, "CONSTANT_Methodref_info's name_and_type_index")
		return &c, p.Err
	case ConstantKindInterfaceMethodref:
		c := constantInterfaceMethodref{cpInfoTag: tag}
		p.readInteger(r, &c.classIndex, "CONSTANT_InterfaceMethodref_info's class_index")
		p.readInteger(r, &c.nameAndTypeIndex, "CONSTANT_InterfaceMethodref_info's name_and_type_index")
		return &c, p.Err
	case ConstantKindString:
		c := constantString{cpInfoTag: tag}
		p.readInteger(r, &c.stringIndex, "CONSTANT_String_info's string_index")
		return &c, p.Err
	case ConstantKindInteger:
		c := constantInteger{cpInfoTag: tag}
		p.readBytes(r, c.bytes[:], "CONSTANT_Integer_info's bytes")
		return &c, p.Err
	case ConstantKindFloat:
		c := constantFloat{cpInfoTag: tag}
		p.readBytes(r, c.bytes[:], "CONSTANT_Float_info's bytes")
		return &c, p.Err
	case ConstantKindLong:
		c := constantLong{cpInfoTag: tag}
		p.readBytes(r, c.high[:], "CONSTANT_Long_info's high_bytes")
		p.readBytes(r, c.low[:], "CONSTANT_Long_info's low_bytes")
		return &c, p.Err
	case ConstantKindDouble:
		c := constantDouble{cpInfoTag: tag}
		p.readBytes(r, c.high[:], "CONSTANT_Double_info's high_bytes")
		p.readBytes(r, c.low[:], "CONSTANT_Double_info's low_bytes")
		return &c, p.Err
	case ConstantKindNameAndType:
		c := constantNameAndType{cpInfoTag: tag}
		p.readInteger(r, &c.nameIndex, "CONSTANT_NameAndType_info's name_index")
		p.readInteger(r, &c.descriptorIndex, "CONSTANT_NameAndType_info's descriptor_index")
		return &c, p.Err
	case ConstantKindUtf8:
		c := constantUtf8{cpInfoTag: tag}
		p.readInteger(r, &c.length, "CONSTANT_Utf8_info's length")
		c.bytes = make([]byte, c.length)
		p.readBytes(r, c.bytes[:], "CONSTANT_Utf8_info's bytes")
		return &c, p.Err
	case ConstantKindMethodHandle:
		c := constantMethodHandle{cpInfoTag: tag}
		p.readInteger(r, &c.referenceKind, "CONSTANT_MethodHandle_info's reference_kind")
		p.readInteger(r, &c.referenceIndex, "CONSTANT_MethodHandle_info's reference_index")
		return &c, p.Err
	case ConstantKindMethodType:
		c := constantMethodType{cpInfoTag: tag}
		p.readInteger(r, &c.descriptorIndex, "CONSTANT_MethodType_info's descriptor_index")
		return &c, p.Err
	case ConstantKindDynamic:
		c := constantDynamic{cpInfoTag: tag}
		p.readInteger(r, &c.bootstrapMethodAttrIndex, "CONSTANT_Dynamic_info's bootstrap_method_attr_index")
		p.readInteger(r, &c.nameAndTypeIndex, "CONSTANT_Dynamic_info's name_and_type_index")
		return &c, p.Err
	case ConstantKindInvokeDynamic:
		c := constantInvokeDynamic{cpInfoTag: tag}
		p.readInteger(r, &c.bootstrapMethodAttrIndex, "CONSTANT_Dynamic_info's bootstrap_method_attr_index")
		p.readInteger(r, &c.nameAndTypeIndex, "CONSTANT_Dynamic_info's name_and_type_index")
		return &c, p.Err
	case ConstantKindModule:
		c := constantModule{cpInfoTag: tag}
		p.readInteger(r, &c.nameIndex, "CONSTANT_Module_info's name_index")
		return &c, p.Err
	case ConstantKindPackage:
		c := constantPackage{cpInfoTag: tag}
		p.readInteger(r, &c.nameIndex, "CONSTANT_Module_info's name_index")
		return &c, p.Err
	}
	return nil, fmt.Errorf("unsupported tag for cp_info: %d", tag.tag)
}

type cpInfoTag struct {
	tag byte
}

func (c *cpInfoTag) Tag() byte { return c.tag }

type constantClass struct {
	cpInfoTag
	nameIndex uint16
}

type constantFieldref struct {
	cpInfoTag
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantMethodref struct {
	cpInfoTag
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantInterfaceMethodref struct {
	cpInfoTag
	classIndex       uint16
	nameAndTypeIndex uint16
}

type constantString struct {
	cpInfoTag
	stringIndex uint16
}

type constantInteger struct {
	cpInfoTag
	bytes [4]byte
}

type constantFloat struct {
	cpInfoTag
	bytes [4]byte
}

type constantLong struct {
	cpInfoTag
	high [4]byte
	low  [4]byte
}

type constantDouble struct {
	cpInfoTag
	high [4]byte
	low  [4]byte
}

type constantNameAndType struct {
	cpInfoTag
	nameIndex       uint16
	descriptorIndex uint16
}

type constantUtf8 struct {
	cpInfoTag
	length uint16
	bytes  []byte
}

type constantMethodHandle struct {
	cpInfoTag
	referenceKind  byte
	referenceIndex uint16
}

type constantMethodType struct {
	cpInfoTag
	descriptorIndex uint16
}

type constantDynamic struct {
	cpInfoTag
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

type constantInvokeDynamic struct {
	cpInfoTag
	bootstrapMethodAttrIndex uint16
	nameAndTypeIndex         uint16
}

type constantModule struct {
	cpInfoTag
	nameIndex uint16
}

type constantPackage struct {
	cpInfoTag
	nameIndex uint16
}

func (c *constantClass) String() string { return fmt.Sprintf("#%d", c.nameIndex) }
func (c *constantFieldref) String() string {
	return fmt.Sprintf("#%d.#%d", c.classIndex, c.nameAndTypeIndex)
}
func (c *constantMethodref) String() string {
	return fmt.Sprintf("#%d.#%d", c.classIndex, c.nameAndTypeIndex)
}
func (c *constantInterfaceMethodref) String() string {
	return fmt.Sprintf("#%d.#%d", c.classIndex, c.nameAndTypeIndex)
}
func (c *constantString) String() string {
	return fmt.Sprintf("#%d", c.stringIndex)
}
func (c *constantInteger) String() string {
	//TODO
	return fmt.Sprintf("%v", c.bytes)
}
func (c *constantFloat) String() string {
	//TODO
	return fmt.Sprintf("%v", c.bytes)
}
func (c *constantLong) String() string {
	//TODO
	return fmt.Sprintf("%v %v", c.high, c.low)
}
func (c *constantDouble) String() string {
	//TODO
	return fmt.Sprintf("%v %v", c.high, c.low)
}
func (c *constantNameAndType) String() string {
	return fmt.Sprintf("#%v:#%v", c.nameIndex, c.descriptorIndex)
}
func (c *constantUtf8) String() string {
	return fmt.Sprintf("%s", c.bytes)
}
func (c *constantMethodHandle) String() string {
	return fmt.Sprintf("%d %d", c.referenceKind, c.referenceIndex)
}
func (c *constantMethodType) String() string {
	return fmt.Sprintf("%d", c.descriptorIndex)
}
func (c *constantDynamic) String() string {
	return fmt.Sprintf("%d %d", c.bootstrapMethodAttrIndex, c.nameAndTypeIndex)
}
func (c *constantInvokeDynamic) String() string {
	return fmt.Sprintf("%d %d", c.bootstrapMethodAttrIndex, c.nameAndTypeIndex)
}
func (c *constantModule) String() string {
	return fmt.Sprintf("%d", c.nameIndex)
}
func (c *constantPackage) String() string {
	return fmt.Sprintf("%d", c.nameIndex)
}

type parser struct {
	Err error
}

func (p *parser) readInteger(r io.Reader, data any, name string) {
	if p.Err != nil {
		return
	}
	if err := binary.Read(r, binary.BigEndian, data); err != nil {
		p.Err = fmt.Errorf("fail to parse %s: %w", name, err)
	}
}

func (p *parser) readBytes(r io.Reader, bytes []byte, name string) {
	if p.Err != nil {
		return
	}
	if n, err := r.Read(bytes); err != nil {
		p.Err = fmt.Errorf("fail to parse %s: %w", name, err)
	} else if n == 0 {
		p.Err = errors.New("fail to parse %s")
	}
}
