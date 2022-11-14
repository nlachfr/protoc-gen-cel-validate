package authorize

import (
	"github.com/google/cel-go/common/types/ref"
)

type TypeAdapterFunc func(value interface{}) ref.Val

func (fn TypeAdapterFunc) NativeToValue(value interface{}) ref.Val {
	return fn(value)
}
