package class

import "fmt"

type annotation struct {
	typeIndex            uint16
	numElementValuePairs uint16
	elementValuePairs    []elementValuePair
}

type elementValuePair struct {
	elementNameIndex uint16
	value            elementValue
}

func parseAnnotation(er *errReader, cf *ClassFile) annotation {
	var a annotation
	if item(er, "type_index", integer(&a.typeIndex, constantPoolStructure[uint16, *constantUtf8](cf))) {
		//TODO field descriptor
	}
	if item(er, "num_element_value_pairs", integer(&a.numElementValuePairs)) {
		a.elementValuePairs = make([]elementValuePair, a.numElementValuePairs)
		for i := 0; i < int(a.numElementValuePairs); i++ {
			p := parent(er, fmt.Sprintf("element_value_pairs[%d]", i))
			child(p, "element_name_index", integer(&a.elementValuePairs[i].elementNameIndex, constantPoolStructure[uint16, *constantUtf8](cf)))
			a.elementValuePairs[i].value = parseElementValue(er, cf)
		}
	}
	return a
}

type elementValue struct {
	tag   uint8
	value elementValueItem
}

func parseElementValue(er *errReader, cf *ClassFile) elementValue {
	var v elementValue
	item(er, "element_value.tag", integer(&v.tag))

	switch rune(v.tag) {
	case 'B', 'C', 'I', 'S', 'Z':
		var i uint16
		item(er, "const_value_index", integer(&i, constantPoolStructure[uint16, *constantInteger](cf)))
		v.value = elementValueConstValueIndex(i)
	case 'D':
		var i uint16
		item(er, "const_value_index", integer(&i, constantPoolStructure[uint16, *constantDouble](cf)))
		v.value = elementValueConstValueIndex(i)
	case 'F':
		var i uint16
		item(er, "const_value_index", integer(&i, constantPoolStructure[uint16, *constantFloat](cf)))
		v.value = elementValueConstValueIndex(i)
	case 'J':
		var i uint16
		item(er, "const_value_index", integer(&i, constantPoolStructure[uint16, *constantLong](cf)))
		v.value = elementValueConstValueIndex(i)
	case 's':
		var i uint16
		item(er, "const_value_index", integer(&i, constantPoolStructure[uint16, *constantString](cf)))
		v.value = elementValueConstValueIndex(i)
	case 'e':
		var e elementValueEnumConstValue
		item(er, "enum_const_value.type_name_index", integer(&e.typeNameIndex, constantPoolStructure[uint16, *constantString](cf)))
		item(er, "enum_const_value.const_name_index", integer(&e.constNameIndex, constantPoolStructure[uint16, *constantString](cf)))
		v.value = &e
	case 'c':
		var i uint16
		item(er, "class_info_index", integer(&i, constantPoolStructure[uint16, *constantString](cf)))
		v.value = elementValueClassInfoIndex(i)
	case '@':
		v.value = elementValueAnnotationValue(parseAnnotation(er, cf))
	case '[':
		var a elementValueArrayValue
		if item(er, "array_value.num_values", integer(&a.numValues)) {
			a.values = make([]elementValue, a.numValues)
			item(er, "array_value.values", entries(a.values, func(e *errReader) elementValue {
				return parseElementValue(e, cf)
			}))
		}
	}
	return v
}

type elementValueItem interface {
	_elementValueItem()
}

type elementValueConstValueIndex uint16
type elementValueEnumConstValue struct {
	typeNameIndex  uint16
	constNameIndex uint16
}
type elementValueClassInfoIndex uint16
type elementValueAnnotationValue annotation
type elementValueArrayValue struct {
	numValues uint16
	values    []elementValue
}

func (elementValueConstValueIndex) _elementValueItem() {}
func (*elementValueEnumConstValue) _elementValueItem() {}
func (elementValueClassInfoIndex) _elementValueItem()  {}
func (elementValueAnnotationValue) _elementValueItem() {}
func (*elementValueArrayValue) _elementValueItem()     {}
