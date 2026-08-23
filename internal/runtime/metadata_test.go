package runtime

import (
	"net"
	"strings"
	"syscall"
	"testing"
)

// The metadata check ran on the string somebody typed, so it caught the one
// spelling a test had thought of and nothing else.
//
// linkLocal parses the host as an IP and gives up when that fails, which is
// every form of the address that is not a dotted quad: a DNS name, the
// decimal or octal spelling of the same number, or a name under a wildcard
// resolver. And fd00:ec2::254 is AWS's IPv6 metadata endpoint, which is
// unique-local rather than link-local, so a check written around
// "link-local" misses it by category rather than by spelling.
//
// refuseMetadata runs from the dialer instead, after resolution, where the
// address is always a literal and the spelling has stopped mattering.
func TestTheMetadataCheckLooksAtAddressesRatherThanSpellings(t *testing.T) {
	refused := func(address string) bool {
		return refuseMetadata("tcp", address, syscall.RawConn(nil)) != nil
	}

	for _, address := range []string{
		"169.254.169.254:80",
		"169.254.1.1:80",
		"[fe80::1]:80",
		// AWS over IPv6. Not link-local, and the reason this is a set of
		// addresses rather than a range.
		"[fd00:ec2::254]:80",
	} {
		if !refused(address) {
			t.Errorf("dialling %s was allowed", address)
		}
	}

	// Everything a webhook receiver actually lives on. Loopback and private
	// addresses are the documented use: the application under test is on
	// localhost, and refusing that would break the tool to no purpose.
	for _, address := range []string{
		"127.0.0.1:4600",
		"[::1]:4600",
		"10.0.0.5:443",
		"192.168.1.20:8080",
		"93.184.216.34:443",
	} {
		if refused(address) {
			t.Errorf("dialling %s was refused", address)
		}
	}
}

// Subscribe keeps its own check, which is worth having as immediate feedback
// and is not the boundary. This says which of the two is which.
func TestSubscribeStillRefusesTheLiteralItCan(t *testing.T) {
	q := &webhookQueue{}

	if err := q.Subscribe("http://169.254.169.254/latest/meta-data/"); err == nil {
		t.Error("subscribing to the metadata address by its literal was allowed")
	} else if !strings.Contains(err.Error(), "link-local") {
		t.Errorf("refused it with %q, which does not say why", err)
	}

	// And these are the spellings it cannot see. They are allowed here and
	// refused at the dialer, which is the arrangement this test exists to
	// pin: a name is not an address until something resolves it.
	for _, endpoint := range []string{
		"http://metadata.google.internal/computeMetadata/v1/",
		"http://2852039166/latest/meta-data/",
	} {
		if err := q.Subscribe(endpoint); err != nil {
			t.Errorf("subscribe(%s) = %v, want it accepted here and refused at the dialer", endpoint, err)
		}
	}

	if !net.ParseIP("fd00:ec2::254").IsGlobalUnicast() {
		t.Error("fd00:ec2::254 stopped being a unique-local address, so this test's premise has changed")
	}
}
