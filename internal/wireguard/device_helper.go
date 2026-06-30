package wireguard

import (
	"errors"
	"net"
	"syscall"

	"github.com/vishvananda/netlink"
)

func CreateDevice(name string) error {
	link := &netlink.Wireguard{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
	}

	err := netlink.LinkAdd(link)
	if err != nil {
		if errors.Is(err, syscall.EEXIST) {
			return nil
		}

		return err
	}

	return nil
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
		if errors.Is(err, syscall.EEXIST) {
			return nil
		}

		return err
	}

	return netlink.LinkSetUp(link)
}

func IsDeviceUp(name string) (bool, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		if _, ok := errors.AsType[netlink.LinkNotFoundError](err); ok {
			return false, nil
		}
		return false, err
	}

	attrs := link.Attrs()

	if attrs.OperState == netlink.OperUp {
		return true, nil
	}

	if attrs.Flags&net.FlagUp != 0 {
		return true, nil
	}

	return false, nil
}

func RemoveDevice(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	return netlink.LinkDel(link)
}
