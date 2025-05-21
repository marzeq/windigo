package config

import (
	"fmt"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/curve"
	"github.com/marzeq/windigo/device"
	"github.com/marzeq/windigo/fan"
)

func GetFans(config ConfigFile, devices device.Devices, curves curve.Curves) (fan.Fans, error) {
	fns, exists, typeOk := mconf_values.ObjGetObject(config, "fans")
	if !exists {
		return nil, fmt.Errorf("section 'fans' not found")
	} else if !typeOk {
		return nil, fmt.Errorf("section 'fans' must be an object")
	}

	fans := make(fan.Fans)
	for dev, fansForDevice := range fns {
		if _, ok := devices[dev]; !ok {
			return nil, fmt.Errorf("device '%s' not found", dev)
		}

		fansForDevice, ok := mconf_values.UnwrapObject(fansForDevice)
		if !ok {
			return nil, fmt.Errorf("fans for device '%s' must be an object", dev)
		}

		for fanName, f := range fansForDevice {
			fanVal, ok := mconf_values.UnwrapObject(f)
			if !ok {
				return nil, fmt.Errorf("fan '%s' must be an object", fanName)
			}

			inputFile, exists, typeOk := mconf_values.ObjGetString(fanVal, "inputfile")
			if !exists {
				return nil, fmt.Errorf("fan '%s' must have 'inputfile' attribute", fanName)
			} else if !typeOk {
				return nil, fmt.Errorf("fan '%s' 'inputfile' must be a string", fanName)
			}
			inputFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, inputFile)

			enableFile, exists, typeOk := mconf_values.ObjGetString(fanVal, "enablefile")
			if !exists {
				return nil, fmt.Errorf("fan '%s' must have 'enablefile' attribute", fanName)
			} else if !typeOk {
				return nil, fmt.Errorf("fan '%s' 'enablefile' must be a string", fanName)
			}
			enableFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, enableFile)

			readFile, exists, typeOk := mconf_values.ObjGetString(fanVal, "readfile")
			if !exists {
				return nil, fmt.Errorf("fan '%s' must have 'readfile' attribute", fanName)
			} else if !typeOk {
				return nil, fmt.Errorf("fan '%s' 'readfile' must be a string", fanName)
			}
			readFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, readFile)

			curveName, exists, typeOk := mconf_values.ObjGetString(fanVal, "curve")
			if !exists {
				return nil, fmt.Errorf("fan '%s' must have 'curve' attribute", fanName)
			} else if !typeOk {
				return nil, fmt.Errorf("fan '%s' 'curve' must be a string", fanName)
			}

			if _, ok := curves[curveName]; !ok {
				return nil, fmt.Errorf("fan '%s' 'curve' not found", fanName)
			}

			if _, ok := fans[fanName]; ok {
				return nil, fmt.Errorf("fan '%s' already exists", fanName)
			}

			fans[fanName] = fan.NewFan(
				fanName,
				curveName,
				inputFilePath,
				enableFilePath,
				readFilePath,
			)
		}
	}

	return fans, nil
}
