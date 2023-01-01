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

func GenerateResourceTypePatternMapping(imports ...protoreflect.FileDescriptor) map[string]string {
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
	return m
}
