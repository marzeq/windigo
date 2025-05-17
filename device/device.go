package device

type DeviceType uint

const (
	DeviceTypeByHwmonX DeviceType = iota
	DeviceTypeByName
	DeviceTypeByPci
)

type Device struct {
	Type DeviceType
	Path string
}

type Devices = map[string]Device // mapping asigned name by user to device path
