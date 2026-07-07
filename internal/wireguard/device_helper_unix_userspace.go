//go:build freebsd || openbsd || darwin

package wireguard

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sync"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

func IfaceName(name string, id uint) string {
	if runtime.GOOS == "openbsd" {
		return fmt.Sprintf("tun%d", id)
	}
	return name
}

type userspaceDevice struct {
	dev      *device.Device
	uapi     net.Listener
	realName string
	wg       sync.WaitGroup
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

	realName, err := tdev.Name()
	if err != nil {
		tdev.Close()
		return fmt.Errorf("getting tun name: %w", err)
	}

	dev := device.NewDevice(tdev, conn.NewDefaultBind(), device.NewLogger(device.LogLevelError, ""))

	uapiFile, err := ipc.UAPIOpen(name)
	if err != nil {
		dev.Close()
		return fmt.Errorf("opening uapi: %w", err)
	}

	uapi, err := ipc.UAPIListen(name, uapiFile)
	if err != nil {
		uapiFile.Close()
		dev.Close()
		return fmt.Errorf("listening uapi: %w", err)
	}

	ud := &userspaceDevice{dev: dev, uapi: uapi, realName: realName}
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

	args := []string{ud.realName, "inet", address}
	if runtime.GOOS == "darwin" {
		ip, _, err := net.ParseCIDR(address)
		if err != nil {
			return fmt.Errorf("parsing address %q: %w", address, err)
		}
		args = append(args, ip.String())
	}
	args = append(args, "up")

	if out, err := exec.Command("ifconfig", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("running ifconfig: %w: %s", err, out)
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
