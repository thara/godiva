package class

import "fmt"

// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.6

type methodInfo struct {
	accessFlag      AccessFlags
	nameIndex       uint16
	descriptorIndex uint16
	attributesCount uint16
	attributes      []attributeInfo
}

func parseMethod(er *errReader, cf *ClassFile) methodInfo {
	var m methodInfo

	var accessFlag uint16
	if item(er, "access_flags", integer(&accessFlag)) {
		m.accessFlag = AccessFlags(accessFlag)
	}

	if item(er, "name_index", integer(&m.nameIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		//TODO must a valid unqualified name
	}
	if item(er, "descriptor_index", integer(&m.descriptorIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		// must a valid field descriptor
	}

	if item(er, "attributes_count", integer(&m.attributesCount)) {
		m.attributes = make([]attributeInfo, m.attributesCount)
		item(er, "attributes", entries(m.attributes, func(er *errReader) attributeInfo {
			return parseMethodAttributeInfo(er, cf)
		}))
	}

	return m
}

func parseMethodAttributeInfo(er *errReader, cf *ClassFile) attributeInfo {
	base, ok := parseAttributeInfoBase(er, cf)
	if !ok {
		return nil
	}

	utf8 := getCpinfo[*constantUtf8](cf, base.attributeNameIndex)
	switch utf8.String() {
	case "Code":
		return base.code(er, cf)
	case "Exceptions":
	case "RuntimeVisibleParameterAnnotations":
	case "RuntimeInvisibleParameterAnnotations":
	case "AnnotationDefault":
	case "MethodParameters":

	case "Synthetic":
		return base.synthetic(er, cf)
	case "Deprecated":
		return base.deprecated(er, cf)
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
