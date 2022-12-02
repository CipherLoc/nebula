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
	IP iputil.VpnIp
	Port uint16
	Protocol uint8
}

/* track a single instance of (ip, port, protocol)
 */
type FirewallIncomingHook struct {
	Drops map[DropData]bool
	Lock sync.Mutex
}

func NewFirewallIncomingHook() *FirewallIncomingHook {
	return &FirewallIncomingHook{
		Drops: make(map[DropData]bool),
	}
}

func (hook *FirewallIncomingHook) FirewallDrop(packet firewall.Packet){
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	if len(hook.Drops) < 10000 {
		data := DropData{IP: packet.RemoteIP, Port: packet.LocalPort, Protocol: packet.Protocol}
		hook.Drops[data] = true
	}
}

func (hook *FirewallIncomingHook) GetAndClear() []DropData {
	hook.Lock.Lock()
	defer hook.Lock.Unlock()
	var out []DropData
	
	for key, _ := range hook.Drops {
		out = append(out, key)
	}

	hook.Drops = make(map[DropData]bool)

	return out
}
