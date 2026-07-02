//go:build freebsd || openbsd

package wireguard

import (
	"fmt"
	"net"
	"os/exec"
)

func IfaceName(_ string, id uint) string {
	return fmt.Sprintf("wg%d", id)
}

func CreateDevice(name string) error {
	if _, err := net.InterfaceByName(name); err == nil {
		return nil
	}
	if out, err := exec.Command("ifconfig", name, "create").CombinedOutput(); err != nil {
		return fmt.Errorf("ifconfig create: %w: %s", err, out)
	}
	return nil
}

func SetDeviceAddress(name, address string) error {
	if out, err := exec.Command("ifconfig", name, "inet", address, "up").CombinedOutput(); err != nil {
		return fmt.Errorf("ifconfig: %w: %s", err, out)
	}
	return nil
}

func IsDeviceUp(name string) (bool, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return false, nil
	}
	return iface.Flags&net.FlagUp != 0, nil
}

func RemoveDevice(name string) error {
	if _, err := net.InterfaceByName(name); err != nil {
		return nil
	}
	if out, err := exec.Command("ifconfig", name, "destroy").CombinedOutput(); err != nil {
		return fmt.Errorf("ifconfig destroy: %w: %s", err, out)
	}
	return nil
}
