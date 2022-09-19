package main

import (
	"testing"
)

func TestGetIPs(t *testing.T) {
	ips, err := getIPs()
	if err != nil {
		t.Fatalf(`getIPs() = %q, %v, error`, ips, err)
	}
	for _, ip := range ips {
		if ip.IsLoopback() {
			t.Fatalf(`getIPs() = %q, %v, should not return loopback!`, ips, err)
		}
		ipv4 := ip.To4()
		if ipv4 == nil {
			t.Fatalf(`getIPs() = %q, %v, should not return ipv6!`, ips, err)
		}
	}
}

func TestGetOutboundIP(t *testing.T) {
	ip, err := getOutboundIP()
	if err != nil {
		t.Fatalf(`getOutboundIP() = %q, %v, error`, ip, err)
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		t.Fatalf(`getOutboundIP() = %q, %v, should not return ipv6!`, ip, err)
	}
}
