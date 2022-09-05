package class

type AccessFlags = byte

// https://docs.oracle.com/javase/specs/jvms/se18/html/jvms-4.html#jvms-4.1-200-E.1
const (
	AccessFlagsPublic     AccessFlags = 0x0001
	AccessFlagsFinal                  = 0x0010
	AccessFlagsSuper                  = 0x0020
	AccessFlagsInterface              = 0x0200
	AccessFlagsAbstract               = 0x0400
	AccessFlagsSynthetic              = 0x1000
	AccessFlagsAnnotation             = 0x2000
	AccessFlagsEnum                   = 0x4000
	AccessFlagsModule                 = 0x8000
)
