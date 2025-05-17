package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/device"
)

func resolveDeviceByName(name string) (string, error) {
	devices, err := os.ReadDir("/sys/class/hwmon")
	if err != nil {
		return "", fmt.Errorf("failed to read /sys/class/hwmon: %w", err)
	}
	currentDevice := ""
	for _, device := range devices {
		if strings.HasPrefix(device.Name(), "hwmon") {
			nameFile := fmt.Sprintf("/sys/class/hwmon/%s/name", device.Name())
			nameVal, err := os.ReadFile(nameFile)
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(nameVal)) == name {
				if currentDevice != "" {
					return "", fmt.Errorf("ambiguous device name:'%s', found multiple devices with such name", name)
				}
				currentDevice = fmt.Sprintf("/sys/class/hwmon/%s", device.Name())
				break
			}
		}
	}

	if currentDevice == "" {
		return "", fmt.Errorf("device '%s' not found", name)
	}

	return currentDevice, nil
}

func resolveDeviceByHwmonX(hwmonX string) (string, error) {
	devicePath := fmt.Sprintf("/sys/class/hwmon/hwmon%s", hwmonX)
	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		return "", fmt.Errorf("device 'hwmon%s' not found", hwmonX)
	}
	return devicePath, nil
}

func resolveDeviceByPci(pci string) (string, error) {
	if len(pci) < 8 {
		pci = "0000:" + pci
	}

	devicePath := fmt.Sprintf("/sys/bus/pci/devices/%s", pci)
	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		return "", fmt.Errorf("device '%s' not found", pci)
	}

	hwmonPath := fmt.Sprintf("/sys/bus/pci/devices/%s/hwmon", pci)
	if _, err := os.Stat(hwmonPath); os.IsNotExist(err) {
		return "", fmt.Errorf("device '%s' is not a hwmon device", pci)
	}

	hwmonDevices, err := os.ReadDir(hwmonPath)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", hwmonPath, err)
	}

	if len(hwmonDevices) == 0 {
		return "", fmt.Errorf("device '%s' has no hwmon devices", pci)
	}
	if len(hwmonDevices) > 1 {
		return "", fmt.Errorf("device '%s' has multiple hwmon devices, not sure which one to use. try using hwmon:X or name:<name> to narrow it down", pci)
	}

	hwmonDevice := hwmonDevices[0].Name()

	devicePath = fmt.Sprintf("/sys/class/hwmon/%s", hwmonDevice)
	if _, err := os.Stat(devicePath); os.IsNotExist(err) {
		return "", fmt.Errorf("device '%s' not found", hwmonDevice)
	}

	return devicePath, nil
}

func GetDevices(config ConfigFile) (device.Devices, error) {
	devicesGot, ok := config["devices"]
	if !ok {
		return nil, fmt.Errorf("section 'devices' not found")
	}
	devicesVal, ok := mconf_values.MconfUnwrapObject(devicesGot)
	if !ok {
		return nil, fmt.Errorf("section 'devices' must be an object")
	}
	devices := make(device.Devices)
	for devicePath, dev := range devicesVal {
		deviceVal, ok := mconf_values.MconfUnwrapString(dev)
		if !ok {
			return nil, fmt.Errorf("device '%s' must be a string", devicePath)
		}
		deviceName := deviceVal
		deviceType := strings.SplitN(devicePath, ":", 2)
		if len(deviceType) != 2 {
			return nil, fmt.Errorf("device '%s' must be in format <type>:<path>", devicePath)
		}
		devicePath := ""
		deviceTypeBy := device.DeviceTypeByHwmonX
		switch deviceType[0] {
		case "hwmon":
			deviceP, err := resolveDeviceByHwmonX(deviceType[1])
			if err != nil {
				return nil, err
			}
			devicePath = deviceP
			deviceTypeBy = device.DeviceTypeByHwmonX
		case "name":
			deviceP, err := resolveDeviceByName(deviceType[1])
			if err != nil {
				return nil, err
			}
			devicePath = deviceP
			deviceTypeBy = device.DeviceTypeByName
		case "pci":
			deviceP, err := resolveDeviceByPci(deviceType[1])
			if err != nil {
				return nil, err
			}
			devicePath = deviceP
			deviceTypeBy = device.DeviceTypeByPci
		}

		if devicePath == "" {
			return nil, fmt.Errorf("device '%s' not found", devicePath)
		}

		if _, ok := devices[deviceName]; ok {
			return nil, fmt.Errorf("device '%s' already exists", deviceName)
		}

		devices[deviceName] = device.Device{
			Type: deviceTypeBy,
			Path: devicePath,
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices specified")
	}

	return devices, nil
}
