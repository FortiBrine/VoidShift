package wireguard

import "github.com/vishvananda/netlink"

func CreateDevice(name string) error {
	link := &netlink.Wireguard{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}

	return netlink.LinkAdd(link)
}

func SetDeviceAddress(name string, address string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		return err
	}

	return netlink.LinkSetUp(link)
}

func RemoveDevice(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	return netlink.LinkDel(link)
}
