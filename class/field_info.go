package class

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
