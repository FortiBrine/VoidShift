//go:build windows

package wireguard

import (
	"fmt"
	"net"
	"net/netip"
	"sync"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

func IfaceName(name string, _ uint) string {
	return name
}

type userspaceDevice struct {
	dev  *device.Device
	uapi net.Listener
	luid winipcfg.LUID
	wg   sync.WaitGroup
}

var (
	mu      sync.Mutex
	devices = make(map[string]*userspaceDevice)
)

func CreateDevice(name string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := devices[name]; ok {
		return nil
	}

	tdev, err := tun.CreateTUN(name, device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("creating tun: %w", err)
	}

	luid := winipcfg.LUID(tdev.(*tun.NativeTun).LUID())

	dev := device.NewDevice(tdev, conn.NewDefaultBind(), device.NewLogger(device.LogLevelError, ""))

	uapi, err := ipc.UAPIListen(name)
	if err != nil {
		dev.Close()
		return fmt.Errorf("listening uapi: %w", err)
	}

	ud := &userspaceDevice{dev: dev, uapi: uapi, luid: luid}
	ud.wg.Add(1)
	go func() {
		defer ud.wg.Done()
		for {
			c, err := uapi.Accept()
			if err != nil {
				return
			}
			go dev.IpcHandle(c)
		}
	}()

	devices[name] = ud
	return nil
}

func SetDeviceAddress(name, address string) error {
	mu.Lock()
	ud, ok := devices[name]
	mu.Unlock()
	if !ok {
		return fmt.Errorf("device %q not found", name)
	}

	prefix, err := netip.ParsePrefix(address)
	if err != nil {
		return fmt.Errorf("parsing address %q: %w", address, err)
	}

	if err := ud.luid.SetIPAddresses([]netip.Prefix{prefix}); err != nil {
		return fmt.Errorf("setting ip addresses: %w", err)
	}

	return nil
}

func IsDeviceUp(name string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()
	_, ok := devices[name]
	return ok, nil
}

func RemoveDevice(name string) error {
	mu.Lock()
	defer mu.Unlock()
	ud, ok := devices[name]
	if !ok {
		return nil
	}
	delete(devices, name)

	ud.uapi.Close()
	ud.dev.Close()
	ud.wg.Wait()
	return nil
}
