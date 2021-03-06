// Copyright 2017 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sciond

import (
	"fmt"
	"strings"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
	"github.com/scionproto/scion/go/proto"
)

type PathErrorCode uint16

const (
	ErrorOk PathErrorCode = iota
	ErrorNoPaths
	ErrorPSTimeout
	ErrorInternal
)

func (c PathErrorCode) String() string {
	switch c {
	case ErrorOk:
		return "OK"
	case ErrorNoPaths:
		return "No paths available"
	case ErrorPSTimeout:
		return "SCIOND timed out while requesting paths"
	case ErrorInternal:
		return "SCIOND experienced an internal error"
	default:
		return fmt.Sprintf("Unknown error (%v)", uint16(c))
	}
}

var _ proto.Cerealizable = (*Pld)(nil)

type Pld struct {
	Id                 uint64
	Which              proto.SCIONDMsg_Which
	PathReq            PathReq
	PathReply          PathReply
	AsInfoReq          ASInfoReq
	AsInfoReply        ASInfoReply
	RevNotification    RevNotification
	RevReply           RevReply
	IfInfoRequest      IFInfoRequest
	IfInfoReply        IFInfoReply
	ServiceInfoRequest ServiceInfoRequest
	ServiceInfoReply   ServiceInfoReply
}

func NewPldFromRaw(b common.RawBytes) (*Pld, error) {
	p := &Pld{}
	return p, proto.ParseFromRaw(p, p.ProtoId(), b)
}

func (p *Pld) ProtoId() proto.ProtoIdType {
	return proto.SCIONDMsg_TypeID
}

func (p *Pld) String() string {
	desc := []string{fmt.Sprintf("Sciond: Id: %d Union: ", p.Id)}
	u1, err := p.union()
	if err != nil {
		desc = append(desc, err.Error())
	} else {
		desc = append(desc, fmt.Sprintf("%+v", u1))
	}
	return strings.Join(desc, "")
}

func (p *Pld) union() (interface{}, error) {
	switch p.Which {
	case proto.SCIONDMsg_Which_pathReq:
		return p.PathReq, nil
	case proto.SCIONDMsg_Which_pathReply:
		return p.PathReply, nil
	case proto.SCIONDMsg_Which_asInfoReq:
		return p.AsInfoReq, nil
	case proto.SCIONDMsg_Which_asInfoReply:
		return p.AsInfoReply, nil
	case proto.SCIONDMsg_Which_revNotification:
		return p.RevNotification, nil
	case proto.SCIONDMsg_Which_revReply:
		return p.RevReply, nil
	case proto.SCIONDMsg_Which_ifInfoRequest:
		return p.IfInfoRequest, nil
	case proto.SCIONDMsg_Which_ifInfoReply:
		return p.IfInfoReply, nil
	case proto.SCIONDMsg_Which_serviceInfoRequest:
		return p.ServiceInfoRequest, nil
	case proto.SCIONDMsg_Which_serviceInfoReply:
		return p.ServiceInfoReply, nil
	}
	return nil, common.NewBasicError("Unsupported SCIOND union type", nil, "type", p.Which)
}

type PathReq struct {
	Dst      addr.IAInt
	Src      addr.IAInt
	MaxPaths uint16
	Flags    PathReqFlags
}

type PathReqFlags struct {
	Flush bool
	Sibra bool
}

type PathReply struct {
	ErrorCode PathErrorCode
	Entries   []PathReplyEntry
}

type PathReplyEntry struct {
	Path     FwdPathMeta
	HostInfo HostInfo
}

type HostInfo struct {
	Port  uint16
	Addrs struct {
		Ipv4 []byte
		Ipv6 []byte
	}
}

func HostInfoFromHostAddr(host addr.HostAddr, port uint16) *HostInfo {
	h := &HostInfo{Port: port}
	if host.Type() == addr.HostTypeIPv4 {
		h.Addrs.Ipv4 = host.IP()
	} else {
		h.Addrs.Ipv6 = host.IP()
	}
	return h
}

func (h *HostInfo) Host() addr.HostAddr {
	if len(h.Addrs.Ipv4) > 0 {
		return addr.HostIPv4(h.Addrs.Ipv4)
	}
	return addr.HostIPv6(h.Addrs.Ipv6)
}

type FwdPathMeta struct {
	FwdPath    []byte
	Mtu        uint16
	Interfaces []PathInterface
}

func (fpm FwdPathMeta) SrcIA() *addr.ISD_AS {
	ifaces := fpm.Interfaces
	if len(ifaces) == 0 {
		return nil
	}
	return ifaces[0].ISD_AS()
}

func (fpm FwdPathMeta) DstIA() *addr.ISD_AS {
	ifaces := fpm.Interfaces
	if len(ifaces) == 0 {
		return nil
	}
	return ifaces[len(ifaces)-1].ISD_AS()
}

func (fpm FwdPathMeta) String() string {
	var hops []string
	for _, intf := range fpm.Interfaces {
		hops = append(hops, intf.String())
	}
	return fmt.Sprintf("Hops: %s Mtu: %d", strings.Join(hops, ">"), fpm.Mtu)
}

type PathInterface struct {
	RawIsdas addr.IAInt `capnp:"isdas"`
	IfID     uint64
}

func (iface *PathInterface) ISD_AS() *addr.ISD_AS {
	return iface.RawIsdas.IA()
}

func (iface PathInterface) String() string {
	return fmt.Sprintf("%v#%v", iface.ISD_AS(), iface.IfID)
}

type ASInfoReq struct {
	Isdas addr.IAInt
}

type ASInfoReply struct {
	Entries []ASInfoReplyEntry
}

type ASInfoReplyEntry struct {
	RawIsdas addr.IAInt `capnp:"isdas"`
	Mtu      uint16
	IsCore   bool
}

func (entry *ASInfoReplyEntry) ISD_AS() *addr.ISD_AS {
	return entry.RawIsdas.IA()
}

func (entry ASInfoReplyEntry) String() string {
	return fmt.Sprintf("ia:%v, mtu:%v, core:%t", entry.ISD_AS(), entry.Mtu, entry.IsCore)
}

type RevNotification struct {
	RevInfo *path_mgmt.RevInfo
}

type RevReply struct {
	Result RevResult
}

type RevResult uint16

const (
	RevValid RevResult = iota
	RevStale
	RevInvalid
	RevUnknown
)

func (c RevResult) String() string {
	switch c {
	case RevValid:
		return "RevValid"
	case RevStale:
		return "RevStale"
	case RevInvalid:
		return "RevInvalid"
	case RevUnknown:
		return "RevUnknown"
	default:
		return fmt.Sprintf("Unknown revocation result (%d)", c)
	}
}

type IFInfoRequest struct {
	IfIDs []uint64
}

type IFInfoReply struct {
	RawEntries []IFInfoReplyEntry `capnp:"entries"`
}

// Entries maps IFIDs to their addresses and ports; the map is rebuilt each time.
func (reply *IFInfoReply) Entries() map[uint64]HostInfo {
	m := make(map[uint64]HostInfo)

	for _, entry := range reply.RawEntries {
		m[entry.IfID] = entry.HostInfo
	}

	return m
}

type IFInfoReplyEntry struct {
	IfID     uint64
	HostInfo HostInfo
}

type ServiceInfoRequest struct {
	ServiceTypes []ServiceType
}

type ServiceType uint16

const (
	SvcBS ServiceType = iota
	SvcPS
	SvcCS
	SvcBR
	SvcSB
)

func (st ServiceType) String() string {
	switch st {
	case SvcBS:
		return "BS"
	case SvcPS:
		return "PS"
	case SvcCS:
		return "CS"
	case SvcBR:
		return "BR"
	case SvcSB:
		return "SB"
	default:
		return "??"
	}
}

type ServiceInfoReply struct {
	Entries []ServiceInfoReplyEntry
}

type ServiceInfoReplyEntry struct {
	ServiceType ServiceType
	Ttl         uint32
	HostInfos   []HostInfo
}
