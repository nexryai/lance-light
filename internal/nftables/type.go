package nftables

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

type IPVersion string

const (
	IPv4 IPVersion = "ipv4"
	IPv6 IPVersion = "ipv6"
)
