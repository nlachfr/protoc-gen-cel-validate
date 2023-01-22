package gateway

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func (x Configuration_Server_Upstream_Protocol) MarshalYAML() (interface{}, error) {
	return strings.ToLower(x.String()), nil
}

func (x *Configuration_Server_Upstream_Protocol) UnmarshalYAML(value *yaml.Node) error {
	switch strings.ToUpper(value.Value) {
	case "", Configuration_Server_Upstream_GRPC.String():
		*x = *Configuration_Server_Upstream_GRPC.Enum()
	case Configuration_Server_Upstream_GRPC_WEB.String():
		*x = *Configuration_Server_Upstream_GRPC_WEB.Enum()
	case Configuration_Server_Upstream_CONNECT.String():
		*x = *Configuration_Server_Upstream_CONNECT.Enum()
	}
	return nil
}
