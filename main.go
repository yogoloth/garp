package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"syscall"
	"unsafe"
)

/*
#include "fill_packet.h"
*/
import "C"

var opts struct {
	If        string `long:"if",short:"i" description:"interface" required:"true"`
	Ip        string `long:"addr" description:"ip" required:"true"`
	Mask      string `long:"mask" description:"mask"`
	FlagSetip bool   `short:"s" description:"set ip"`
}

func SendGratuitousArp(iface string, req_ip string) {
	etherArp := new(C.arp_packet)
	size := uint(unsafe.Sizeof(*etherArp))
	LogDebug.Println("ArpPacketSize : ", size)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		LogError.Println("open ap_packet socket error: " + err.Error())
		return
	}
	LogDebug.Println("obtained fd ", fd)
	defer syscall.Close(fd)

	// Get Mac address
	interf, err := net.InterfaceByName(iface)
	if err != nil {
		LogError.Printf("could not find %s interface\n", iface)
		return
	}

	LogDebug.Println("interface hw address: ", interf.HardwareAddr)

	iface_cstr := C.CString(interf.HardwareAddr.String())
	ip_cstr := C.CString(req_ip)

	ppacket := C.fill_arp_packet(iface_cstr, ip_cstr)
	packet := C.GoBytes(unsafe.Pointer(ppacket), C.int(size))
	C.free(unsafe.Pointer(ppacket))

	// Send the packet
	var addr syscall.SockaddrLinklayer
	addr.Protocol = syscall.ETH_P_ARP
	addr.Ifindex = interf.Index
	addr.Hatype = syscall.ARPHRD_ETHER

	err = syscall.Sendto(fd, packet, 0, &addr)

	if err != nil {
		LogError.Println("sent packet error: ", err.Error())
		os.Exit(1)
	} else {
		LogInfo.Println("sent packet success")
	}
}

func AddIp(iface string, req_ip string, mask string) {
	dst_if, err := netlink.LinkByName(iface)
	if err != nil {
		LogError.Printf("get %s error: %s\n", iface, err.Error())
		os.Exit(1)
	}

	addr, err := netlink.ParseAddr(req_ip + "/" + mask)
	if err != nil {
		LogError.Printf("parse ip %s/%s error: %s\n", req_ip, mask, err.Error())
		os.Exit(1)
	}
	netlink.AddrAdd(dst_if, addr)
}

//   ./garp --addr 172.17.5.182 --if wlp4s0
func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		LogError.Println("args error")
		os.Exit(1)
	}
	if len(opts.Mask) == 0 {
		opts.Mask = "24"
	}
	LogDebug.Println("Got interface:", opts.If, "ip:", opts.Ip, "mask:", opts.Mask)

	if opts.FlagSetip {
		AddIp(opts.If, opts.Ip, opts.Mask)
	}
	SendGratuitousArp(opts.If, opts.Ip)

}
