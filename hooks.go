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
	Protocol uint8 // from firewall/packet.go
}

/* tracks a dropped packet for single instance of (local ip, local port, remote ip, remote port, protocol)
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

/* Add a dropped packet to the set. The packet information is hashed so that we don't store the same information twice.
 */
func (hook *FirewallDropHook) FirewallDrop(packet firewall.Packet){
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if len(hook.Drops) < 10000 {
		data := DropData{LocalIP: packet.LocalIP, RemoteIP: packet.RemoteIP, LocalPort: packet.LocalPort, RemotePort: packet.RemotePort, Protocol: packet.Protocol}
		hook.Drops[data] = true
	}
}

/* Get all current dropped packet information, and reset the packet drop collection bucket */
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
