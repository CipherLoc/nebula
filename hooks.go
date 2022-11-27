package nebula

import (
	"github.com/slackhq/nebula/firewall"
	"github.com/slackhq/nebula/iputil"
)

type Hook interface {
	FirewallDrop(packet firewall.Packet)
}

type FirewallIncomingHook struct {
	Drops map[iputil.VpnIp]int
}

func NewFirewallIncomingHook() *FirewallIncomingHook {
	return &FirewallIncomingHook{
		Drops: make(map[iputil.VpnIp]int),
	}
}

func (hook *FirewallIncomingHook) FirewallDrop(packet firewall.Packet){
	hook.Drops[packet.RemoteIP] = int(packet.RemotePort)
}
