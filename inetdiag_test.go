package netlink

import (
	"fmt"
	"syscall"
	"testing"
	"unsafe"
)

func TestSizeofInetDiagReqV2(t *testing.T) {
	req := InetDiagReqV2{}
	if unsafe.Sizeof(req) != SizeofInetDiagReqV2 {
		t.Error("size of InetDiagReqV2 error")
	}
}

func TestTCPStat(t *testing.T) {
	req := NewInetDiagRequest()
	data := NewInetDiagReqV2(syscall.AF_INET, syscall.IPPROTO_TCP, TCP_ALL)
	req.AddData(data)
	msgs, err := req.Execute(syscall.NETLINK_INET_DIAG, 0)
	if err != nil {
		t.Fatalf("failed to receive from netlink socket: %v", err)
	}

	for _, msg := range msgs {
		diamsg := ParseInetDiagMsg(msg)
		fmt.Println(diamsg)
	}
}
