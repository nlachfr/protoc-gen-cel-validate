package plugin

import (
	"fmt"
	"testing"

	"github.com/Neakxs/protocel/authorize"
)

func TestLoadConfig(t *testing.T) {
	c := &authorize.AuthorizeOptions{}
	if err := LoadConfig("../../../../testdata/authorize/config.yml", c); err != nil {
		t.Errorf("want nil, got %v", err)
	}
	fmt.Println(c)
}
