package nebula

import (
	"sync"
	"github.com/slackhq/nebula/firewall"
	"github.com/slackhq/nebula/iputil"
)

type Hook interface {
	FirewallDrop(packet firewall.Packet)
}

// FIXME: need a better name
type DropData struct {
	LocalIP iputil.VpnIp
	RemoteIP iputil.VpnIp
	LocalPort uint16
	RemotePort uint16
	Protocol uint8
}

/* track a single instance of (ip, port, protocol)
 */
type FirewallDropHook struct {
	Drops map[DropData]bool
	Lock sync.Mutex
}

func NewFirewallDropHook() *FirewallDropHook {
	return &FirewallDropHook{
		Drops: make(map[DropData]bool),
	}
}

func (hook *FirewallDropHook) FirewallDrop(packet firewall.Packet){
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if len(hook.Drops) < 10000 {
		data := DropData{LocalIP: packet.LocalIP, RemoteIP: packet.RemoteIP, LocalPort: packet.LocalPort, RemotePort: packet.RemotePort, Protocol: packet.Protocol}
		hook.Drops[data] = true
	}
}

func (hook *FirewallDropHook) GetAndClear() []DropData {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	var out []DropData
	
	for key, _ := range hook.Drops {
		out = append(out, key)
	}

	hook.Drops = make(map[DropData]bool)

	return out
}
