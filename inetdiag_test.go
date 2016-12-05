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
	socket, err := NewNetlinkSocket(syscall.NETLINK_INET_DIAG, 0)
	if err != nil {
		t.Fatalf("failed to create netlink socket: %v", err)
	}
	defer socket.Close()

	req := NewNetlinkRequest()
	req.Type = SOCK_DIAG_BY_FAMILY
	req.Flags = syscall.NLM_F_DUMP | syscall.NLM_F_REQUEST
	data := NewInetDiagReqV2(syscall.AF_INET, syscall.IPPROTO_TCP, TCP_ALL)
	req.AddData(data)
	if err := socket.Send(req); err != nil {
		t.Fatalf("failed to send netlink request: %v", err)
	}

	msgs, err := socket.Receive()
	if err != nil {
		t.Fatalf("failed to receive from netlink socket: %v", err)
	}

	for _, msg := range msgs {
		diamsg := ParseInetDiagMsg(msg.Data)
		fmt.Printf("[%s] %s:%d -> %s:%d\n", TcpStatesMap[diamsg.IDiagState], diamsg.Id.SrcIP(),
			ntohs(diamsg.Id.IDiagSPort), diamsg.Id.DstIP(), ntohs(diamsg.Id.IDiagDPort))
	}
}
