package metadata

import (
	"github.com/gookit/goutil/errorx/panics"
)

func MakeTypeRange(min int64, max int64) Type {
	panics.IsTrue(min < max)

	return &rangeType{min, max}
}

func MakeTypeText() Type {
	return &textType{}
}

func MakeTypeFloat() Type {
	return &floatType{}
}

func MakeTypeBool() Type {
	return &boolType{}
}

// Note: values should be sorted
func MakeTypeEnum(values ...string) Type {
	panics.IsTrue(values != nil && len(values) > 0)

	return &enumType{values}
}

func MakeTypeComplex() Type {
	return &complexType{}
}
