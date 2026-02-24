package multiaddr

import (
	"net/netip"
	"testing"

	"github.com/multiformats/go-multiaddr/x/meg"
)

func TestMatchAndCaptureMultiaddr(t *testing.T) {
	m := StringCast("/ip4/1.2.3.4/udp/8231/quic-v1/webtransport")

	var udpPort string
	found, _ := m.Match(
		meg.Or(
			meg.Val(P_IP4),
			meg.Val(P_IP6),
		),
		meg.CaptureString(P_UDP, &udpPort),
		meg.Val(P_QUIC_V1),
		meg.Val(P_WEBTRANSPORT),
	)
	if !found {
		t.Fatal("failed to match")
	}
	if udpPort != "8231" {
		t.Fatal("unexpected value")
	}
}

func TestCaptureAddrPort(t *testing.T) {
	m := StringCast("/ip4/1.2.3.4/udp/8231/quic-v1/webtransport")
	var addrPort netip.AddrPort
	var network string

	found, err := m.Match(
		CaptureAddrPort(&network, &addrPort),
		meg.ZeroOrMore(meg.Any),
	)
	if err != nil {
		t.Fatal("error", err)
	}
	if !found {
		t.Fatal("failed to match")
	}
	if !addrPort.IsValid() {
		t.Fatal("failed to capture addrPort")
	}
	if network != "udp" {
		t.Fatal("unexpected network", network)
	}
	if addrPort.String() != "1.2.3.4:8231" {
		t.Fatal("unexpected ipPort", addrPort)
	}
}
