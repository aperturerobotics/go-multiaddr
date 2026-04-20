package meg_test

import (
	"slices"
	"testing"

	"github.com/aperturerobotics/go-multiaddr"
	"github.com/aperturerobotics/go-multiaddr/x/meg"
)

type preallocatedCapture struct {
	certHashes []string
	matcher    meg.Matcher
}

func preallocateCapture() *preallocatedCapture {
	p := &preallocatedCapture{}
	p.matcher = meg.PatternToMatcher(
		meg.Or(
			meg.Val(multiaddr.P_IP4),
			meg.Val(multiaddr.P_IP6),
			meg.Val(multiaddr.P_DNS),
		),
		meg.Val(multiaddr.P_UDP),
		meg.Val(multiaddr.P_WEBRTC_DIRECT),
		meg.CaptureZeroOrMoreStrings(multiaddr.P_CERTHASH, &p.certHashes),
	)
	return p
}

var webrtcMatchPrealloc *preallocatedCapture

func (p *preallocatedCapture) IsWebRTCDirectMultiaddr(addr multiaddr.Multiaddr) (bool, int) {
	found, _ := meg.Match(p.matcher, addr)
	return found, len(p.certHashes)
}

// IsWebRTCDirectMultiaddr returns whether addr is a /webrtc-direct multiaddr with the count of certhashes
// in addr
func IsWebRTCDirectMultiaddr(addr multiaddr.Multiaddr) (bool, int) {
	if webrtcMatchPrealloc == nil {
		webrtcMatchPrealloc = preallocateCapture()
	}
	return webrtcMatchPrealloc.IsWebRTCDirectMultiaddr(addr)
}

// IsWebRTCDirectMultiaddrLoop returns whether addr is a /webrtc-direct multiaddr with the count of certhashes
// in addr
func IsWebRTCDirectMultiaddrLoop(addr multiaddr.Multiaddr) (bool, int) {
	protos := [...]int{multiaddr.P_IP4, multiaddr.P_IP6, multiaddr.P_DNS, multiaddr.P_UDP, multiaddr.P_WEBRTC_DIRECT}
	matchProtos := [...][]int{protos[:3], {protos[3]}, {protos[4]}}
	certHashCount := 0
	for i, c := range addr {
		if i >= len(matchProtos) {
			if c.Code() == multiaddr.P_CERTHASH {
				certHashCount++
			} else {
				return false, 0
			}
		} else {
			found := slices.Contains(matchProtos[i], c.Code())
			if !found {
				return false, 0
			}
		}
	}
	return true, certHashCount
}

func BenchmarkIsWebRTCDirectMultiaddr(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/webrtc-direct/")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWebRTC, count := IsWebRTCDirectMultiaddr(addr)
		if !isWebRTC || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}

func BenchmarkIsWebRTCDirectMultiaddrLoop(b *testing.B) {
	addr := multiaddr.StringCast("/ip4/1.2.3.4/udp/1234/webrtc-direct/")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isWebRTC, count := IsWebRTCDirectMultiaddrLoop(addr)
		if !isWebRTC || count != 0 {
			b.Fatal("unexpected result")
		}
	}
}
