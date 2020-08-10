package field

//go:generate log_xfields_generator

//调用方式: go generate cmd/log_xfields_generator.go

type FieldType int

const (
	UnknownType FieldType = 0
	BoolType    FieldType = 1
	Int64Type   FieldType = 2
	Float64Type FieldType = 3
	StringType  FieldType = 4
	ObjectType  FieldType = 5
)

type Field interface {
	Name() string
	Type() FieldType
	Value() interface{}
}
