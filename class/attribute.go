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

type attributeCode struct {
	attributeInfoBase
	maxStack             uint16
	maxLocals            uint16
	codeLength           uint32
	code                 []uint8
	exceptionTableLength uint16
	exceptionTable       []exceptionTableEntry
	attributesCount      uint16
	attributes           []attributeInfo
}

type exceptionTableEntry struct {
	startPC   uint16
	endPC     uint16
	handlerPC uint16
	catchType uint16
}

func (base *attributeInfoBase) code(er *errReader, cf *ClassFile) *attributeCode {
	attr := attributeCode{attributeInfoBase: *base}
	item(er, "max_stack", integer(&attr.maxStack))
	item(er, "max_locals", integer(&attr.maxLocals))

	item(er, "code_length", integer(&attr.codeLength, min[uint32](1)))

	attr.code = make([]uint8, attr.codeLength)
	item(er, "code", bytes(attr.code))
	//TODO validation code

	item(er, "exception_table_length", integer(&attr.exceptionTableLength))
	attr.exceptionTable = make([]exceptionTableEntry, attr.exceptionTableLength)
	item(er, "exception_table", entries(attr.exceptionTable, func(er *errReader) exceptionTableEntry {
		var e exceptionTableEntry
		item(er, "start_pc", integer(&e.startPC))
		item(er, "end_pc", integer(&e.endPC))
		item(er, "handler_pc", integer(&e.handlerPC))
		//TODO validation

		if item(er, "catch_type", integer(&e.catchType)) {
			if e.catchType != 0 {
				validate(er, e.catchType, constantPoolStructure[uint16, *constantClass](cf))
				//TODO The verifier checks that the class is Throwable or a subclass of Throwable (§4.9.2).
			}
		}
		return e
	}))

	item(er, "attributes_count", integer(&attr.attributesCount))
	item(er, "attributes", entries(attr.attributes, func(er *errReader) attributeInfo {
		return parseCodeAttributeInfo(er, cf)
	}))
	return &attr
}

func parseCodeAttributeInfo(er *errReader, cf *ClassFile) attributeInfo {
	base, ok := parseAttributeInfoBase(er, cf)
	if !ok {
		return nil
	}

	utf8 := getCpinfo[*constantUtf8](cf, base.attributeNameIndex)
	switch utf8.String() {
	case "LineNumberTable":
		return base.lineNumberTable(er, cf)
	case "LocalVariableTable":
		return base.localVariableTable(er, cf)
	case "LocalVariableTypeTable":
		return base.localVariableTypeTable(er, cf)
	case "StackMapTable":
		return base.stackMapTable(er, cf)
	}

	er.err = fmt.Errorf("unsupported attribute name at index(%d)", base.attributeNameIndex)
	return nil
}

type attributeLineNumberTable struct {
	attributeInfoBase
	lineNumberTableLength uint16
	lineNumberTable       []lineNumberTableEntry
}

type lineNumberTableEntry struct {
	startPC    uint16
	lineNumber uint16
}

func (a *attributeInfoBase) lineNumberTable(er *errReader, cf *ClassFile) *attributeLineNumberTable {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.12
	attr := attributeLineNumberTable{attributeInfoBase: *a}

	item(er, "line_number_table_length", integer(&attr.lineNumberTableLength))
	attr.lineNumberTable = make([]lineNumberTableEntry, attr.lineNumberTableLength)
	item(er, "line_number_tablee", entries(attr.lineNumberTable, func(er *errReader) lineNumberTableEntry {
		var e lineNumberTableEntry
		item(er, "start_pc", integer(&e.startPC))
		//TODO valid index into the code array of this Code attribute
		item(er, "line_number", integer(&e.lineNumber))
		return e
	}))
	return &attr
}

type attributeLocalVariableTable struct {
	attributeInfoBase
	localVariableTableLength uint16
	localVariableTable       []localVariableTableEntry
}

type localVariableTableEntry struct {
	startPC         uint16
	length          uint16
	nameIndex       uint16
	descriptorIndex uint16
	index           uint16
}

func (a *attributeInfoBase) localVariableTable(er *errReader, cf *ClassFile) *attributeLocalVariableTable {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.13
	attr := attributeLocalVariableTable{attributeInfoBase: *a}
	item(er, "local_variable_table_length", integer(&attr.localVariableTableLength))
	attr.localVariableTable = make([]localVariableTableEntry, attr.localVariableTableLength)
	item(er, "line_number_tablee", entries(attr.localVariableTable, func(er *errReader) localVariableTableEntry {
		var e localVariableTableEntry
		item(er, "start_pc", integer(&e.startPC))
		//TODO valid index into the code array of this Code attribute
		item(er, "length", integer(&e.length))
		// TODO start_pc + length must either be a valid index into the code array of this Code attribute and be the index of the opcode of an instruction, or it must be the first index beyond the end of that code array.
		item(er, "name_index", integer(&e.nameIndex, constantPoolStructure[uint16, *constantUtf8](cf)))
		item(er, "descriptor_index", integer(&e.descriptorIndex, constantPoolStructure[uint16, *constantUtf8](cf)))
		item(er, "index", integer(&e.index))
		//TODO  must be a valid index into the local variable array of the current frame.
		return e
	}))
	return &attr
}

type attributeLocalVariableTypeTable struct {
	attributeInfoBase
	localVariableTypeTableLength uint16
	localVariableTypeTable       []localVariableTypeTableEntry
}

type localVariableTypeTableEntry struct {
	startPC        uint16
	length         uint16
	nameIndex      uint16
	signatureIndex uint16
	index          uint16
}

func (a *attributeInfoBase) localVariableTypeTable(er *errReader, cf *ClassFile) *attributeLocalVariableTypeTable {
	// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.7.14
	attr := attributeLocalVariableTypeTable{attributeInfoBase: *a}
	item(er, "local_variable_table_length", integer(&attr.localVariableTypeTableLength))
	attr.localVariableTypeTable = make([]localVariableTypeTableEntry, attr.localVariableTypeTableLength)
	item(er, "line_number_tablee", entries(attr.localVariableTypeTable, func(er *errReader) localVariableTypeTableEntry {
		var e localVariableTypeTableEntry
		item(er, "start_pc", integer(&e.startPC))
		//TODO valid index into the code array of this Code attribute
		item(er, "length", integer(&e.length))
		// TODO start_pc + length must either be a valid index into the code array of this Code attribute and be the index of the opcode of an instruction, or it must be the first index beyond the end of that code array.
		item(er, "name_index", integer(&e.nameIndex, constantPoolStructure[uint16, *constantUtf8](cf)))
		item(er, "signature_index", integer(&e.signatureIndex, constantPoolStructure[uint16, *constantUtf8](cf)))
		item(er, "index", integer(&e.index))
		//TODO  must be a valid index into the local variable array of the current frame.
		return e
	}))
	return &attr
}

type attributeStackMapTable struct {
	attributeInfoBase
	numberOfEntries uint16
	entries         []stackMapFrame
}

func (a *attributeInfoBase) stackMapTable(er *errReader, cf *ClassFile) *attributeStackMapTable {
	attr := attributeStackMapTable{attributeInfoBase: *a}
	item(er, "number_of_entries", integer(&attr.numberOfEntries))
	attr.entries = make([]stackMapFrame, attr.numberOfEntries)
	item(er, "entries", entries(attr.entries, func(er *errReader) stackMapFrame {
		var b stackMapFrameBase
		item(er, "frame_type", integer(&b.frameType))
		switch {
		case 0 <= b.frameType && b.frameType <= 63:
			return &stackMapFrameSameFrame{stackMapFrameBase: b}
		case 64 <= b.frameType && b.frameType <= 127:
			s := stackMapFrameSameLocals1{stackMapFrameBase: b}
			s.stack = make([]verificationType, 1)
			item(er, "stack", entries(s.stack, func(er *errReader) verificationType {
				return parseVerificationType(er, cf)
			}))
			return &s
		case 128 <= b.frameType && b.frameType <= 246:
			//NOTE reserved for future use
		case b.frameType == 247:
			s := stackMapFrameSameLocals1Extended{stackMapFrameBase: b}
			item(er, "offset_delta", integer(&s.offsetDelta))
			s.stack = make([]verificationType, 1)
			item(er, "stack", entries(s.stack, func(er *errReader) verificationType {
				return parseVerificationType(er, cf)
			}))
			return &s
		case 248 <= b.frameType && b.frameType <= 250:
			s := stackMapFrameChopFrame{stackMapFrameBase: b}
			item(er, "offset_delta", integer(&s.offsetDelta))
			return &s
		case b.frameType == 251:
			s := stackMapFrameSameFrameExtended{stackMapFrameBase: b}
			item(er, "offset_delta", integer(&s.offsetDelta))
			return &s
		case 252 <= b.frameType && b.frameType <= 254:
			s := stackMapFrameAppendFrame{stackMapFrameBase: b}
			item(er, "offset_delta", integer(&s.offsetDelta))
			s.locals = make([]verificationType, b.frameType-251)
			item(er, "locals", entries(s.locals, func(er *errReader) verificationType {
				return parseVerificationType(er, cf)
			}))
			//TODO validation
			return &s
		case b.frameType == 255:
			s := stackMapFrameFullFrame{stackMapFrameBase: b}
			item(er, "offset_delta", integer(&s.offsetDelta))

			item(er, "number_of_locals", integer(&s.numberOfLocals))
			s.locals = make([]verificationType, s.numberOfLocals)
			item(er, "locals", entries(s.locals, func(er *errReader) verificationType {
				return parseVerificationType(er, cf)
			}))

			item(er, "number_of_stack_items", integer(&s.numberOfStackItems))
			s.stack = make([]verificationType, s.numberOfStackItems)
			item(er, "stack", entries(s.stack, func(er *errReader) verificationType {
				return parseVerificationType(er, cf)
			}))
			//TODO validation
			return &s
		}

		er.err = fmt.Errorf("invalid frame type(%d)", b.frameType)
		return nil
	}))
	return &attr
}

type stackMapFrame interface {
	_stackMapFrame()
}

type stackMapFrameBase struct {
	frameType uint8
}

func (*stackMapFrameBase) _stackMapFrame() {}

type stackMapFrameSameFrame struct{ stackMapFrameBase }
type stackMapFrameSameLocals1 struct {
	stackMapFrameBase
	stack []verificationType
}
type stackMapFrameSameLocals1Extended struct {
	stackMapFrameBase
	offsetDelta uint16
	stack       []verificationType
}
type stackMapFrameChopFrame struct {
	stackMapFrameBase
	offsetDelta uint16
}
type stackMapFrameSameFrameExtended struct {
	stackMapFrameBase
	offsetDelta uint16
	stack       []verificationType
}
type stackMapFrameAppendFrame struct {
	stackMapFrameBase
	offsetDelta uint16
	locals      []verificationType
}
type stackMapFrameFullFrame struct {
	stackMapFrameBase
	offsetDelta        uint16
	numberOfLocals     uint16
	locals             []verificationType
	numberOfStackItems uint16
	stack              []verificationType
}

type verificationType interface {
	_verificationTypeInfo()
}

func parseVerificationType(er *errReader, cf *ClassFile) verificationType {
	var b verificationTypeInfoBase
	item(er, "tag", integer(&b.tag))
	switch b.tag {
	case 0:
		return &verificationTypeTop{verificationTypeInfoBase: b}
	case 1:
		return &verificationTypeInteger{verificationTypeInfoBase: b}
	case 2:
		return &verificationTypeFloat{verificationTypeInfoBase: b}
	case 5:
		return &verificationTypeNull{verificationTypeInfoBase: b}
	case 6:
		return &verificationTypeUninitializedThis{verificationTypeInfoBase: b}
	case 7:
		v := verificationTypeObject{verificationTypeInfoBase: b}
		item(er, "cpool_index", integer(&v.cpoolIndex))
		return &v
	case 8:
		v := verificationTypeUninitialized{verificationTypeInfoBase: b}
		item(er, "offset", integer(&v.offset))
		return &v

	case 4:
		return &verificationTypeLong{verificationTypeInfoBase: b}
	case 3:
		return &verificationTypeDouble{verificationTypeInfoBase: b}
	}
	//TODO validate The Long_variable_info and Double_variable_info

	er.err = fmt.Errorf("invalid tag(%d)", b.tag)
	return nil
}

type verificationTypeInfoBase struct {
	tag uint8
}

func (b *verificationTypeInfoBase) _verificationTypeInfo()

type verificationTypeTop struct{ verificationTypeInfoBase }
type verificationTypeInteger struct{ verificationTypeInfoBase }
type verificationTypeFloat struct{ verificationTypeInfoBase }
type verificationTypeLong struct{ verificationTypeInfoBase }
type verificationTypeDouble struct{ verificationTypeInfoBase }
type verificationTypeNull struct{ verificationTypeInfoBase }
type verificationTypeUninitializedThis struct{ verificationTypeInfoBase }
type verificationTypeObject struct {
	verificationTypeInfoBase
	cpoolIndex uint16
}
type verificationTypeUninitialized struct {
	verificationTypeInfoBase
	offset uint16
}
