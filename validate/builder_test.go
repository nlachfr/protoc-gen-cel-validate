package validate

import (
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestBuildServiceRuleValidater(t *testing.T) {
	tests := []struct {
		Name        string
		BuildOpts   []BuildOption
		ServiceDesc protoreflect.ServiceDescriptor
		WantErr     bool
	}{
		{
			Name:        "No validation",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")),
			WantErr:     false,
		},
		{
			Name:        "Service level expr",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceExpr")),
			WantErr:     false,
		},
		{
			Name:        "Service level expr with missing const",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceOptions")),
			WantErr:     true,
		},
		{
			Name:        "Service level expr with global const",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceOptions")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Service config expr with global const",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
						ServiceRules: map[string]*ServiceRule{
							string(validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")).FullName()): {
								Rule: &Rule{
									Programs: []*Rule_Program{{Expr: `attribute_context.request.headers[isAdmHdr] == "true"`}},
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Service level expr with local const",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceLocalOptions")),
			WantErr:     false,
		},
		{
			Name:        "Service level expr with const conflict",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("ServiceLocalOptions")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Method level expr",
			ServiceDesc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodExpr")),
			WantErr:     false,
		},
		{
			Name:        "Method level with missing const",
			ServiceDesc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodOptions")),
			WantErr:     true,
		},
		{
			Name:        "Method level with global const",
			ServiceDesc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodOptions")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Method level with local const",
			ServiceDesc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodLocalOptions")),
			WantErr:     false,
		},
		{
			Name:        "Method level with const conflict",
			ServiceDesc: validate.File_testdata_validate_method_proto.Services().ByName(protoreflect.Name("MethodLocalOptions")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Method config with const conflict",
			ServiceDesc: validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"isAdmHdr": "x-is-admin",
								},
							},
						},
						ServiceRules: map[string]*ServiceRule{
							string(validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")).FullName()): {
								MethodRules: map[string]*MethodRule{
									string(validate.File_testdata_validate_service_proto.Services().ByName(protoreflect.Name("Service")).Methods().ByName("Rpc").Name()): {
										Rule: &Rule{
											Programs: []*Rule_Program{{Expr: `attribute_context.request.headers["x-is-admin"] == "true"`}},
										},
									},
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := NewBuilder(tt.BuildOpts...).BuildServiceRuleValidater(tt.ServiceDesc)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestBuildMessageRuleValidater(t *testing.T) {
	tests := []struct {
		Name        string
		BuildOpts   []BuildOption
		MessageDesc protoreflect.MessageDescriptor
		WantErr     bool
	}{
		{
			Name:        "No validation",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			WantErr:     false,
		},
		{
			Name:        "Message level expr",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageExpr"),
			WantErr:     false,
		},
		{
			Name:        "Message level expr with missing const",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageOptions"),
			WantErr:     true,
		},
		{
			Name:        "Message level expr with global const",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageOptions"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Message level expr with local const",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageLocalOptions"),
			WantErr:     false,
		},
		{
			Name:        "Message level expr with const conflict",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("MessageLocalOptions"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Message config expr with const conflict",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
						MessageRules: map[string]*MessageRule{
							string(validate.File_testdata_validate_message_proto.Messages().ByName("Message").FullName()): {
								Rule: &Rule{Programs: []*Rule_Program{{Expr: `name != ""`}}},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Field level expr",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldExpr"),
			WantErr:     false,
		},
		{
			Name:        "Field resource reference wrong",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceWrong"),
			WantErr:     true,
		},
		{
			Name:        "Field resource reference type",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceType"),
			WantErr:     false,
		},
		{
			Name:        "Field resource reference child type",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceChild"),
			WantErr:     false,
		},
		{
			Name:        "Field resource reference type and child type",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldReferenceTypeAndChild"),
			WantErr:     true,
		},
		{
			Name:        "Field repeated resource reference",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldRepeatedReferenceType"),
			WantErr:     false,
		},
		{
			Name:        "Field level expr with missing const",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldOptions"),
			WantErr:     true,
		},
		{
			Name:        "Field level expr with global const",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldOptions"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Field level expr with local const",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldLocalOptions"),
			WantErr:     false,
		},
		{
			Name:        "Field level expr with const conflict",
			MessageDesc: validate.File_testdata_validate_field_proto.Messages().ByName("FieldLocalOptions"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
		{
			Name:        "Field config expr with const conflict",
			MessageDesc: validate.File_testdata_validate_message_proto.Messages().ByName("Message"),
			BuildOpts: []BuildOption{
				WithOptions(&Options{
					Rule: &FileRule{
						Options: &options.Options{
							Globals: &options.Options_Globals{
								Constants: map[string]string{
									"emptyName": "",
								},
							},
						},
						MessageRules: map[string]*MessageRule{
							string(validate.File_testdata_validate_message_proto.Messages().ByName("Message").FullName()): {
								FieldRules: map[string]*FieldRule{
									string(validate.File_testdata_validate_message_proto.Messages().ByName("Message").Fields().ByName("name").FullName()): {
										Rule: &Rule{Programs: []*Rule_Program{{Expr: `name != ""`}}},
									},
								},
							},
						},
					},
				}),
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := NewBuilder(tt.BuildOpts...).BuildMessageRuleValidater(tt.MessageDesc)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}

		})
	}
}
