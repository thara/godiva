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

func parseCodeAttributeInfo(er *errReader, cf *ClassFile) attributeInfo {
	base, ok := parseAttributeInfoBase(er, cf)
	if !ok {
		return nil
	}

	utf8 := getCpinfo[*constantUtf8](cf, base.attributeNameIndex)
	switch utf8.String() {
	case "LineNumberTable":
	case "LocalVariableTable":
	case "LocalVariableTypeTable":
	case "StackMapTable":
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

type attributeSynthetic struct {
	attributeInfoBase
}

func (base *attributeInfoBase) synthetic(er *errReader, cf *ClassFile) *attributeSynthetic {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.8
	attr := attributeSynthetic{attributeInfoBase: *base}
	if base.attributeLength != 0 {
		er.err = fmt.Errorf("invalid attribute length(%d) for Synthetic_attribute", base.attributeLength)
		return nil
	}
	return &attr
}

type attributeDeprecated struct {
	attributeInfoBase
}

func (base *attributeInfoBase) deprecated(er *errReader, cf *ClassFile) *attributeDeprecated {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.15
	attr := attributeDeprecated{attributeInfoBase: *base}
	if base.attributeLength != 0 {
		er.err = fmt.Errorf("invalid attribute length(%d) for Deprecated_attribute", base.attributeLength)
		return nil
	}
	return &attr
}

type attributeSignature struct {
	attributeInfoBase
	signatureIndex uint16
}

func (base *attributeInfoBase) signature(er *errReader, cf *ClassFile) *attributeSignature {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.9
	attr := attributeSignature{attributeInfoBase: *base}
	if base.attributeLength != 2 {
		er.err = fmt.Errorf("invalid attribute length(%d) for Signature_attribute", base.attributeLength)
		return nil
	}

	if item(er, "signature_index", integer(&attr.signatureIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		utf8 := getCpinfo[*constantUtf8](cf, attr.signatureIndex)
		switch utf8.String() {
		//TODO
		// a class signature if this Signature attribute is an attribute of a ClassFile structure
		// a method signature if this Signature attribute is an attribute of a method_info structure
		// or a field signature otherwise.
		}
	}
	return &attr
}

type attributeRuntimeVisibleAnnotations struct {
	attributeInfoBase
	numAnnotations uint16
	annotations    []annotation
}

func (base *attributeInfoBase) runtimeVisibleAnnotations(er *errReader, cf *ClassFile) *attributeRuntimeVisibleAnnotations {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.16
	attr := attributeRuntimeVisibleAnnotations{attributeInfoBase: *base}
	item(er, "num_annotations", integer(&attr.numAnnotations))
	item(er, "annotations", entries(attr.annotations, func(er *errReader) annotation {
		return parseAnnotation(er, cf)
	}))
	return &attr
}

type attributeRuntimeInvisibleAnnotations struct {
	attributeInfoBase
	numAnnotations uint16
	annotations    []annotation
}

func (base *attributeInfoBase) runtimeInvisibleAnnotations(er *errReader, cf *ClassFile) *attributeRuntimeInvisibleAnnotations {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.17
	attr := attributeRuntimeInvisibleAnnotations{attributeInfoBase: *base}
	item(er, "num_annotations", integer(&attr.numAnnotations))
	attr.annotations = make([]annotation, attr.numAnnotations)
	item(er, "annotations", entries(attr.annotations, func(er *errReader) annotation {
		return parseAnnotation(er, cf)
	}))
	return &attr
}

type parameterAnnotation struct {
	numAnnotations uint16
	annotations    []annotation
}

type attributeRuntimeVisibleTypeAnnotations struct {
	attributeInfoBase
	numParameters        uint8
	parameterAnnotations []parameterAnnotation
}

func (base *attributeInfoBase) runtimeVisibleTypeAnnotations(er *errReader, cf *ClassFile) *attributeRuntimeVisibleTypeAnnotations {
	attr := attributeRuntimeVisibleTypeAnnotations{attributeInfoBase: *base}
	item(er, "num_parameters", integer(&attr.numParameters))
	attr.parameterAnnotations = make([]parameterAnnotation, attr.numParameters)
	item(er, "annotations", entries(attr.parameterAnnotations, func(er *errReader) parameterAnnotation {
		var p parameterAnnotation
		item(er, "num_annotations", integer(&p.numAnnotations))
		p.annotations = make([]annotation, p.numAnnotations)
		item(er, "annotations", entries(p.annotations, func(er *errReader) annotation {
			return parseAnnotation(er, cf)
		}))
		return p
	}))
	return &attr
}

type attributeRuntimeInvisibleTypeAnnotations struct {
	attributeInfoBase
	numParameters        uint8
	parameterAnnotations []parameterAnnotation
}

func (base *attributeInfoBase) runtimeInvisibleTypeAnnotations(er *errReader, cf *ClassFile) *attributeRuntimeInvisibleTypeAnnotations {
	attr := attributeRuntimeInvisibleTypeAnnotations{attributeInfoBase: *base}
	item(er, "num_parameters", integer(&attr.numParameters))
	attr.parameterAnnotations = make([]parameterAnnotation, attr.numParameters)
	item(er, "annotations", entries(attr.parameterAnnotations, func(er *errReader) parameterAnnotation {
		var p parameterAnnotation
		item(er, "num_annotations", integer(&p.numAnnotations))
		p.annotations = make([]annotation, p.numAnnotations)
		item(er, "annotations", entries(p.annotations, func(er *errReader) annotation {
			return parseAnnotation(er, cf)
		}))
		return p
	}))
	return &attr
}
