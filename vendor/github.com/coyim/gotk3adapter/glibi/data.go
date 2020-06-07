package glibi

type ApplicationFlags int

type SignalHandle uint

type SourceHandle uint

type Type uint

var (
	APPLICATION_FLAGS_NONE           ApplicationFlags
	APPLICATION_IS_SERVICE           ApplicationFlags
	APPLICATION_HANDLES_OPEN         ApplicationFlags
	APPLICATION_HANDLES_COMMAND_LINE ApplicationFlags
	APPLICATION_SEND_ENVIRONMENT     ApplicationFlags
	APPLICATION_NON_UNIQUE           ApplicationFlags
)

var (
	TYPE_INVALID   Type
	TYPE_NONE      Type
	TYPE_INTERFACE Type
	TYPE_CHAR      Type
	TYPE_UCHAR     Type
	TYPE_BOOLEAN   Type
	TYPE_INT       Type
	TYPE_UINT      Type
	TYPE_LONG      Type
	TYPE_ULONG     Type
	TYPE_INT64     Type
	TYPE_UINT64    Type
	TYPE_ENUM      Type
	TYPE_FLAGS     Type
	TYPE_FLOAT     Type
	TYPE_DOUBLE    Type
	TYPE_STRING    Type
	TYPE_POINTER   Type
	TYPE_BOXED     Type
	TYPE_PARAM     Type
	TYPE_OBJECT    Type
	TYPE_VARIANT   Type
)

var (
	VARIANT_TYPE_BOOLEAN           VariantType
	VARIANT_TYPE_BYTE              VariantType
	VARIANT_TYPE_INT16             VariantType
	VARIANT_TYPE_UINT16            VariantType
	VARIANT_TYPE_INT32             VariantType
	VARIANT_TYPE_UINT32            VariantType
	VARIANT_TYPE_INT64             VariantType
	VARIANT_TYPE_UINT64            VariantType
	VARIANT_TYPE_HANDLE            VariantType
	VARIANT_TYPE_DOUBLE            VariantType
	VARIANT_TYPE_STRING            VariantType
	VARIANT_TYPE_ANY               VariantType
	VARIANT_TYPE_BASIC             VariantType
	VARIANT_TYPE_TUPLE             VariantType
	VARIANT_TYPE_UNIT              VariantType
	VARIANT_TYPE_DICTIONARY        VariantType
	VARIANT_TYPE_STRING_ARRAY      VariantType
	VARIANT_TYPE_OBJECT_PATH_ARRAY VariantType
	VARIANT_TYPE_BYTESTRING        VariantType
	VARIANT_TYPE_BYTESTRING_ARRAY  VariantType
	VARIANT_TYPE_VARDICT           VariantType
)
