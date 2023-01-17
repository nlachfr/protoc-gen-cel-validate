package validate

import (
	"regexp"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var patternRegexp = regexp.MustCompile(`\{[\w-]+\}`)
var newPattern = `[\\w-\\.]+`

func GenerateResourceTypePatternMapping(desc protoreflect.Descriptor) map[string]string {
	imps := []protoreflect.FileDescriptor{desc.ParentFile()}
	for i := 0; i < desc.ParentFile().Imports().Len(); i++ {
		imps = append(imps, desc.ParentFile().Imports().Get(i))
	}
	m := map[string]string{}
	for i := 0; i < len(imps); i++ {
		messages := imps[i].Messages()
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
	return m
}
