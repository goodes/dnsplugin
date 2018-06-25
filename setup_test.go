package dnsplugin

import (
	"github.com/mholt/caddy"
	"testing"
)

func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", PhDNS{}.Name())
	if err := setup(c); err != nil {
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	c = caddy.NewTestController("dns", `phdns example.org`)
	if err := setup(c); err == nil {
		t.Fatalf("Expected errors, but got: %v", err)
	}

	t.Log(c)
	t.Log(c.ServerType())
}
