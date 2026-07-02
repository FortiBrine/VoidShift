//go:build windows

package wireguard

import (
	"fmt"
	"net"
	"os/exec"
	"sync"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

type userspaceDevice struct {
	dev  *device.Device
	uapi net.Listener
}

var (
	devMu   sync.Mutex
	devices = make(map[string]*userspaceDevice)
)

func IfaceName(name string, _ uint) string {
	return name
}

func CreateDevice(name string) error {
	devMu.Lock()
	defer devMu.Unlock()

	if _, ok := devices[name]; ok {
		return nil
	}

	tdev, err := tun.CreateTUN(name, device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("create tun %q: %w", name, err)
	}

	dev := device.NewDevice(tdev, conn.NewDefaultBind(), device.NewLogger(device.LogLevelError, ""))

	uapi, err := ipc.UAPIListen(name)
	if err != nil {
		dev.Close()
		return fmt.Errorf("listen uapi: %w", err)
	}

	go func() {
		for {
			c, err := uapi.Accept()
			if err != nil {
				return
			}
			go dev.IpcHandle(c)
		}
	}()

	devices[name] = &userspaceDevice{dev: dev, uapi: uapi}
	return nil
}

func SetDeviceAddress(name, address string) error {
	ip, ipNet, err := net.ParseCIDR(address)
	if err != nil {
		return fmt.Errorf("parse address %q: %w", address, err)
	}
	mask := net.IP(ipNet.Mask).String()
	cmd := exec.Command(
		"netsh", "interface", "ipv4", "set", "address",
		"name="+name, "source=static",
		"address="+ip.String(),
		"mask="+mask,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("netsh: %w: %s", err, out)
	}
	return nil
}

func IsDeviceUp(name string) (bool, error) {
	devMu.Lock()
	_, ok := devices[name]
	devMu.Unlock()
	return ok, nil
}

func RemoveDevice(name string) error {
	devMu.Lock()
	ud, ok := devices[name]
	if !ok {
		devMu.Unlock()
		return nil
	}
	delete(devices, name)
	devMu.Unlock()

	ud.uapi.Close()
	ud.dev.Close()
	return nil
}
