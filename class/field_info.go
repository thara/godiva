package class

import "fmt"

type fieldInfo struct {
	accessFlag      AccessFlags
	nameIndex       uint16
	descriptorIndex uint16
	attributesCount uint16
	attributes      []attributeInfo
}

func parseField(er *errReader, cf *ClassFile) fieldInfo {
	var f fieldInfo

	var accessFlag uint16
	if item(er, "access_flags", integer(&accessFlag)) {
		f.accessFlag = AccessFlags(accessFlag)
	}

	if item(er, "name_index", integer(&f.nameIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		//TODO must a valid unqualified name
	}
	if item(er, "descriptor_index", integer(&f.descriptorIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		// must a valid field descriptor
	}

	if item(er, "attributes_count", integer(&f.attributesCount)) {
		f.attributes = make([]attributeInfo, f.attributesCount)
		item(er, "attributes", entries(f.attributes, func(er *errReader) attributeInfo {
			return parseFieldAttributeInfo(er, cf)
		}))
	}
	return f
}

func parseFieldAttributeInfo(er *errReader, cf *ClassFile) attributeInfo {
	base, ok := parseAttributeInfoBase(er, cf)
	if !ok {
		return nil
	}

	utf8 := getCpinfo[*constantUtf8](cf, base.attributeNameIndex)
	switch utf8.String() {
	case "ConstantValue":
		return base.constantValue(er, cf)
	case "Synthetic":
		return base.synthetic(er, cf)
	case "Deprecated":
		return base.deprecated(er, cf)
	case "Signature":
		return base.signature(er, cf)
	case "RuntimeVisibleAnnotations":
		return base.runtimeVisibleAnnotations(er, cf)
	case "RuntimeInvisibleAnnotations":
		return base.runtimeInvisibleTypeAnnotations(er, cf)
	case "RuntimeVisibleTypeAnnotations":
		return base.runtimeVisibleTypeAnnotations(er, cf)
	case "RuntimeInvisibleTypeAnnotations":
		return base.runtimeInvisibleTypeAnnotations(er, cf)
	}

	er.err = fmt.Errorf("unsupported attribute name at index(%d)", base.attributeNameIndex)
	return nil
}
