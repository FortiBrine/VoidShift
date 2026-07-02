//go:build darwin

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

func IfaceName(name string, _ uint) string {
	return name
}

type userspaceDevice struct {
	dev      *device.Device
	uapi     net.Listener
	realName string
}

var (
	devMu   sync.Mutex
	devices = make(map[string]*userspaceDevice)
)

func CreateDevice(name string) error {
	devMu.Lock()
	defer devMu.Unlock()

	if _, ok := devices[name]; ok {
		return nil
	}

	tdev, err := tun.CreateTUN(name, device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("create tun: %w", err)
	}

	realName, err := tdev.Name()
	if err != nil {
		tdev.Close()
		return fmt.Errorf("get tun name: %w", err)
	}

	dev := device.NewDevice(tdev, conn.NewDefaultBind(), device.NewLogger(device.LogLevelError, ""))

	uapiFile, err := ipc.UAPIOpen(name)
	if err != nil {
		dev.Close()
		return fmt.Errorf("open uapi: %w", err)
	}

	uapi, err := ipc.UAPIListen(name, uapiFile)
	if err != nil {
		uapiFile.Close()
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

	devices[name] = &userspaceDevice{dev: dev, uapi: uapi, realName: realName}
	return nil
}

func SetDeviceAddress(name, address string) error {
	devMu.Lock()
	ud, ok := devices[name]
	devMu.Unlock()
	if !ok {
		return fmt.Errorf("device %q not found", name)
	}
	ip, _, err := net.ParseCIDR(address)
	if err != nil {
		return fmt.Errorf("parse address %q: %w", address, err)
	}
	if out, err := exec.Command("ifconfig", ud.realName, "inet", address, ip.String(), "up").CombinedOutput(); err != nil {
		return fmt.Errorf("ifconfig: %w: %s", err, out)
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
