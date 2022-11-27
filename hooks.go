package nebula

import (
	"github.com/slackhq/nebula/firewall"
)

type Hook interface {
	FirewallDrop(packet firewall.Packet)
}
