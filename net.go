package main

import (
	"errors"
	"log"
	"net"
	"strconv"
)

func getIPs() ([]net.IP, error) {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip)
		}
	}
	if len(ips) > 0 {
		return ips, nil
	}
	return nil, errors.New("are you connected to the network?")
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func checkIPs() {
	ips, err := getIPs()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("The Proxy will be listening on these IP-Addresses:")
		for index, ip := range ips {
			log.Println("#" + strconv.Itoa(index+1) + ": " + ip.String())
		}
		ip, err := getOutboundIP()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("The most likely IP-Address to use for Plex should be: " + ip.String())
		}
	}
}
