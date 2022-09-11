package class

import (
	"fmt"
)

type attributeInfo interface {
	_attributeInfo()
}

func parseAttributeInfoBase(er *errReader, cf *ClassFile) (base attributeInfoBase, ok bool) {
	if item(er, "attribute_name_index", integer(&base.attributeNameIndex)) {
		validate(er, constantPoolStructure[uint16, *constantUtf8](cf))
	} else {
		return base, false
	}
	item(er, "attribute_length", integer(&base.attributeLength))
	return base, true
}

func parseFieldAttributeInfo(er *errReader, cf *ClassFile, f *fieldInfo) attributeInfo {
	base, ok := parseAttributeInfoBase(er, cf)
	if !ok {
		return nil
	}

	utf8 := getCpinfo[*constantUtf8](cf, base.attributeNameIndex)
	switch utf8.String() {
	case "ConstantValue":
		return base.constantValue(er, cf)
	}

	er.err = fmt.Errorf("unsupported attribute name at index(%d)", base.attributeNameIndex)
	return nil
}

type attributeInfoBase struct {
	attributeNameIndex uint16
	attributeLength    uint32
}

func (attributeInfoBase) _attributeInfo() {}

type attributeConstantValue struct {
	attributeInfoBase
	constantValueIndex uint16
}

func (base *attributeInfoBase) constantValue(er *errReader, cf *ClassFile) *attributeConstantValue {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.2
	attr := attributeConstantValue{attributeInfoBase: *base}
	if base.attributeLength != 2 {
		er.err = fmt.Errorf("invalid attribute length(%d) for ConstantValue", base.attributeLength)
		return nil
	}

	if item(er, "ConstantValue_attribute's constantvalue_index", integer(&attr.constantValueIndex, existConstantPool[uint16](cf))) {
		e := must(cf.lookupConstantPool(attr.constantValueIndex))

		//TODO validate to match field types
		switch e.(type) {
		case *constantInteger:
			// int, short, char, byte, boolean
		case *constantFloat:
			// float
		case *constantLong:
			// long
		case *constantDouble:
			// double
		case *constantString:
			// String
		default:
			er.err = fmt.Errorf("invalid constant pool entry structure at constantValueIndex(%d)", attr.constantValueIndex)
			return nil
		}
	}

	return &attr
}
