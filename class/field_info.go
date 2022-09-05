package class

type fieldInfo struct {
	accessFlag      AccessFlags
	nameIndex       uint16
	descriptorIndex uint16
	attributesCount uint16
	attributes      []attribute
}
