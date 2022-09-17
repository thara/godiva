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

func parseCpInfo(r *errReader) cpInfo {
	var tag cpInfoTag

	item(r, "cp_info tag", integer(&tag.tag))

	switch tag.tag {
	case ConstantKindClass:
		c := constantClass{cpInfoTag: tag}
		item(r, "CONSTANT_Class_info's name_index", integer(&c.nameIndex))
		return &c
	case ConstantKindFieldref:
		c := constantFieldref{cpInfoTag: tag}
		item(r, "CONSTANT_Fieldref_info's class_index", integer(&c.classIndex))
		item(r, "CONSTANT_Fieldref_info's name_and_type_index", integer(&c.nameAndTypeIndex))
		return &c
	case ConstantKindMethodref:
		c := constantMethodref{cpInfoTag: tag}
		item(r, "CONSTANT_Methodref_info's class_index", integer(&c.classIndex))
		item(r, "CONSTANT_Methodref_info's name_and_type_index", integer(&c.nameAndTypeIndex))
		return &c
	case ConstantKindInterfaceMethodref:
		c := constantInterfaceMethodref{cpInfoTag: tag}
		item(r, "CONSTANT_InterfaceMethodref_info's class_index", integer(&c.classIndex))
		item(r, "CONSTANT_InterfaceMethodref_info's name_and_type_index", integer(&c.nameAndTypeIndex))
		return &c
	case ConstantKindString:
		c := constantString{cpInfoTag: tag}
		item(r, "CONSTANT_String_info's string_index", integer(&c.stringIndex))
		return &c
	case ConstantKindInteger:
		c := constantInteger{cpInfoTag: tag}
		item(r, "CONSTANT_Integer_info's bytes", bytes(c.bytes[:]))
		return &c
	case ConstantKindFloat:
		c := constantFloat{cpInfoTag: tag}
		item(r, "CONSTANT_Float_info's bytes", bytes(c.bytes[:]))
		return &c
	case ConstantKindLong:
		c := constantLong{cpInfoTag: tag}
		item(r, "CONSTANT_Long_info's high_bytes", bytes(c.high[:]))
		item(r, "CONSTANT_Long_info's low_bytes", bytes(c.low[:]))
		return &c
	case ConstantKindDouble:
		c := constantDouble{cpInfoTag: tag}
		item(r, "CONSTANT_Double_info's high_bytes", bytes(c.high[:]))
		item(r, "CONSTANT_Double_info's low_bytes", bytes(c.low[:]))
		return &c
	case ConstantKindNameAndType:
		c := constantNameAndType{cpInfoTag: tag}
		item(r, "CONSTANT_NameAndType_info's name_index", integer(&c.nameIndex))
		item(r, "CONSTANT_NameAndType_info's descriptor_index", integer(&c.descriptorIndex))
		return &c
	case ConstantKindUtf8:
		c := constantUtf8{cpInfoTag: tag}
		item(r, "CONSTANT_Utf8_info's length", integer(&c.length))
		c.bytes = make([]byte, c.length)
		item(r, "CONSTANT_Utf8_info's bytes", bytes(c.bytes))
		return &c
	case ConstantKindMethodHandle:
		c := constantMethodHandle{cpInfoTag: tag}
		item(r, "CONSTANT_MethodHandle_info's reference_kind", integer(&c.referenceKind))
		item(r, "CONSTANT_MethodHandle_info's reference_index", integer(&c.referenceIndex))
		return &c
	case ConstantKindMethodType:
		c := constantMethodType{cpInfoTag: tag}
		item(r, "CONSTANT_MethodType_info's descriptor_index", integer(&c.descriptorIndex))
		return &c
	case ConstantKindDynamic:
		c := constantDynamic{cpInfoTag: tag}
		item(r, "CONSTANT_Dynamic_info's bootstrap_method_attr_index", integer(&c.bootstrapMethodAttrIndex))
		item(r, "CONSTANT_Dynamic_info's name_and_type_index", integer(&c.nameAndTypeIndex))
		return &c
	case ConstantKindInvokeDynamic:
		c := constantInvokeDynamic{cpInfoTag: tag}
		item(r, "CONSTANT_InvokeDynamic_info's bootstrap_method_attr_index", integer(&c.bootstrapMethodAttrIndex))
		item(r, "CONSTANT_InvokeDynamic_info's name_and_type_index", integer(&c.nameAndTypeIndex))
		return &c
	case ConstantKindModule:
		c := constantModule{cpInfoTag: tag}
		item(r, "CONSTANT_Module_info's name_index", integer(&c.nameIndex))
		return &c
	case ConstantKindPackage:
		c := constantPackage{cpInfoTag: tag}
		item(r, "CONSTANT_Module_info's name_index", integer(&c.nameIndex))
		return &c
	}
	r.err = fmt.Errorf("unsupported tag for cp_info: %d", tag.tag)
	return nil
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

func (p *parser) readInteger(r io.Reader, data any, name string) error {
	if p.Err != nil {
		return p.Err
	}
	if err := binary.Read(r, binary.BigEndian, data); err != nil {
		p.Err = fmt.Errorf("fail to parse %s: %w", name, err)
	}
	return p.Err
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
