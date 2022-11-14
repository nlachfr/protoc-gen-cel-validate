package options

import (
	"github.com/google/cel-go/checker/decls"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TypeFromOverloadType(t *Options_Overloads_Type) *v1alpha1.Type {
	switch v := t.Type.(type) {
	case *Options_Overloads_Type_Primitive_:
		switch v.Primitive {
		case Options_Overloads_Type_BOOL:
			return decls.Bool
		case Options_Overloads_Type_INT:
			return decls.Int
		case Options_Overloads_Type_UINT:
			return decls.Uint
		case Options_Overloads_Type_DOUBLE:
			return decls.Double
		case Options_Overloads_Type_BYTES:
			return decls.Bytes
		case Options_Overloads_Type_STRING:
			return decls.String
		case Options_Overloads_Type_DURATION:
			return decls.Duration
		case Options_Overloads_Type_TIMESTAMP:
			return decls.Timestamp
		case Options_Overloads_Type_ERROR:
			return decls.Error
		case Options_Overloads_Type_DYN:
			return decls.Dyn
		case Options_Overloads_Type_ANY:
			return decls.Any
		}
	case *Options_Overloads_Type_Object:
		return decls.NewObjectType(v.Object)
	case *Options_Overloads_Type_Array_:
		return decls.NewListType(TypeFromOverloadType(v.Array.Type))
	case *Options_Overloads_Type_Map_:
		return decls.NewMapType(TypeFromOverloadType(v.Map.Key), TypeFromOverloadType(v.Map.Value))
	}
	return decls.Null
}
