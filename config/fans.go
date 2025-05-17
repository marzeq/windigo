package config

import (
	"fmt"

	"github.com/marzeq/mconf/mconf_values"
	"github.com/marzeq/windigo/curve"
	"github.com/marzeq/windigo/device"
	"github.com/marzeq/windigo/fan"
)

func GetFans(config ConfigFile, devices device.Devices, curves curve.Curves) (fan.Fans, error) {
	fansGot, ok := config["fans"]
	if !ok {
		return nil, fmt.Errorf("section 'fans' not found")
	}
	fansVal, ok := mconf_values.MconfUnwrapObject(fansGot)
	if !ok {
		return nil, fmt.Errorf("section 'fans' must be an object")
	}

	fans := make(fan.Fans)
	for dev, fansForDevice := range fansVal {
		if _, ok := devices[dev]; !ok {
			return nil, fmt.Errorf("device '%s' not found", dev)
		}

		fansForDeviceVal, ok := mconf_values.MconfUnwrapObject(fansForDevice)
		if !ok {
			return nil, fmt.Errorf("fans for device '%s' must be an object", dev)
		}

		for fanName, f := range fansForDeviceVal {
			fanVal, ok := mconf_values.MconfUnwrapObject(f)
			if !ok {
				return nil, fmt.Errorf("fan '%s' must be an object", fanName)
			}

			inputFileGot, ok := fanVal["inputfile"]
			if !ok {
				return nil, fmt.Errorf("fan '%s' must have 'inputfile' attribute", fanName)
			}
			inputFileVal, ok := mconf_values.MconfUnwrapString(inputFileGot)
			if !ok {
				return nil, fmt.Errorf("fan '%s' 'inputfile' must be a string", fanName)
			}
			inputFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, inputFileVal)

			enableFileGot, ok := fanVal["enablefile"]
			if !ok {
				return nil, fmt.Errorf("fan '%s' must have 'enablefile' attribute", fanName)
			}
			enableFileVal, ok := mconf_values.MconfUnwrapString(enableFileGot)
			if !ok {
				return nil, fmt.Errorf("fan '%s' 'enablefile' must be a string", fanName)
			}
			enableFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, enableFileVal)

			readFileGot, ok := fanVal["readfile"]
			if !ok {
				return nil, fmt.Errorf("fan '%s' must have 'readfile' attribute", fanName)
			}
			readFileVal, ok := mconf_values.MconfUnwrapString(readFileGot)
			if !ok {
				return nil, fmt.Errorf("fan '%s' 'readfile' must be a string", fanName)
			}
			readFilePath := fmt.Sprintf("%s/%s", devices[dev].Path, readFileVal)

			curveGot, ok := fanVal["curve"]
			if !ok {
				return nil, fmt.Errorf("fan '%s' must have 'curve' attribute", fanName)
			}
			curveVal, ok := mconf_values.MconfUnwrapString(curveGot)
			if !ok {
				return nil, fmt.Errorf("fan '%s' 'curve' must be a string", fanName)
			}

			if _, ok := curves[curveVal]; !ok {
				return nil, fmt.Errorf("fan '%s' 'curve' not found", fanName)
			}

			if _, ok := fans[fanName]; ok {
				return nil, fmt.Errorf("fan '%s' already exists", fanName)
			}

			fans[fanName] = fan.NewFan(
				fanName,
				curveVal,
				inputFilePath,
				enableFilePath,
				readFilePath,
			)
		}
	}

	return fans, nil
}
