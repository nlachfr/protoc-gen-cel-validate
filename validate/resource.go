package validate

import (
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var patternRegexp = regexp.MustCompile(`\{[\w-]+\}`)
var newPattern = `[\\w-\\.]+`

func GenerateResourceTypePatternMapping(file protoreflect.FileDescriptor, imports ...protoreflect.FileDescriptor) (map[string]string, error) {
	m := map[string]string{}
	for i := 0; i < len(imports); i++ {
		messages := imports[i].Messages()
		for j := 0; j < messages.Len(); j++ {
			msg := messages.Get(j)
			resource := proto.GetExtension(msg.Options(), annotations.E_Resource).(*annotations.ResourceDescriptor)
			if resource != nil {
				p := []string{}
				for _, pattern := range resource.Pattern {
					p = append(p, patternRegexp.ReplaceAllString(pattern, newPattern))
				}
				m[resource.Type] = "(" + strings.Join(p, "|") + ")"
			}
		}
	}
	messages := file.Messages()
	for i := 0; i < messages.Len(); i++ {
		fields := messages.Get(i).Fields()
		for j := 0; j < fields.Len(); j++ {
			field := fields.Get(j)
			if ref := proto.GetExtension(field.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference); ref != nil {
				if ref.Type != "" && ref.ChildType != "" {
					return nil, fmt.Errorf(`resource reference error: type and child_type are defined`)
				}
				if ref.Type != "*" && ref.Type != "" {
					if _, ok := m[ref.Type]; !ok {
						return nil, fmt.Errorf(`cannot find type "%s"`, ref.Type)
					}
				}
				if ref.ChildType != "" {
					if _, ok := m[ref.ChildType]; !ok {
						return nil, fmt.Errorf(`cannot find type "%s"`, ref.ChildType)
					}
				}
			}
		}
	}
	return m, nil
}
