package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
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
	If string `long:"If",short:"i" description:"interface" required:"true"`
	Ip string `long:"Addr" description:"ip" required:"true"`
}

func SendGratuitousArp(iface string, req_ip string) {
	etherArp := new(C.arp_packet)
	size := uint(unsafe.Sizeof(*etherArp))
	fmt.Println("Size : ", size)

	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, syscall.ETH_P_ALL)
	if err != nil {
		fmt.Println("open ap_packet socket error: " + err.Error())
		return
	}
	fmt.Println("Obtained fd ", fd)
	defer syscall.Close(fd)

	// Get Mac address
	interf, err := net.InterfaceByName(iface)
	if err != nil {
		fmt.Printf("Could not find %s interface\n", iface)
		return
	}

	fmt.Println("Interface hw address: ", interf.HardwareAddr)

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
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Sent packet")
	}
}

//   ./garp --Addr 172.17.5.182 --If wlp4s0
func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println("args error")
		os.Exit(1)
	}
	fmt.Println(opts.If)
	fmt.Println(opts.Ip)

	SendGratuitousArp(opts.If, opts.Ip)

}
