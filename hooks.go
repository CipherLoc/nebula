package nebula

import (
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
}

func NewFirewallIncomingHook() *FirewallIncomingHook {
	return &FirewallIncomingHook{
		Drops: make(map[DropData]bool),
	}
}

func (hook *FirewallIncomingHook) FirewallDrop(packet firewall.Packet){
	data := DropData{IP: packet.RemoteIP, Port: packet.LocalPort, Protocol: packet.Protocol}
	hook.Drops[data] = true
}
