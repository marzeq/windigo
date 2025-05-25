package cli

import (
	"fmt"
	"net"

	"github.com/marzeq/windigo/common"
	"github.com/marzeq/windigo/config"
)

func sum(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func printTable(headers []string, data [][]string) {
	if len(data) == 0 {
		return
	}
	for i := range data {
		if len(headers) != len(data[i]) {
			fmt.Println("Error: headers and data length mismatch")
			return
		}
	}

	maxLengths := make([]int, len(headers))
	for i := range headers {
		maxLengths[i] = len(headers[i])
		for j := range data {
			if len(data[j]) > i {
				if len(data[j][i]) > maxLengths[i] {
					maxLengths[i] = len(data[j][i])
				}
			}
		}
	}

	for i := range headers {
		fmt.Printf("%-*s ", maxLengths[i], headers[i])
		if i < len(headers)-1 {
			fmt.Print("| ")
		}
	}
	fmt.Println()
	for range sum(maxLengths) + (len(headers)-1)*3 {
		fmt.Print("-")
	}
	fmt.Println()
	for _, row := range data {
		for i := range row {
			fmt.Printf("%-*s ", maxLengths[i], row[i])
			if i < len(row)-1 {
				fmt.Print("| ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func RunCli(confFile string, constants map[string]string) int {
	cfg, err := config.ReadConfig(confFile, constants)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println("windigo version", common.VERSION)
	fmt.Println()

	headers := []string{"Sensor", "Temp", "Offset Temp"}
	data := make([][]string, 0)
	for _, sensor := range cfg.Sensors {
		temp, err := sensor.ReadTemperature()
		if err != nil {
			fmt.Printf("Error reading temperature: %s\n", err.Error())
			continue
		}
		realTemp, _ := sensor.ReadRealTemperature()
		data = append(data, []string{sensor.Name, fmt.Sprintf("%.1f °C", realTemp), fmt.Sprintf("%.1f °C", temp)})
	}
	printTable(headers, data)

	headers = []string{"Fan", "RPM", "%"}
	data = make([][]string, 0)
	for _, fan := range cfg.Fans {
		speed, err := fan.ReadSpeed()
		if err != nil {
			fmt.Printf("Error reading fan speed: %s\n", err.Error())
			continue
		}
		curve := cfg.Curves[fan.Curve]
		temp, err := curve.GetAggregateTemp(cfg.Sensors)
		if err != nil {
			data = append(data, []string{fan.Name, fmt.Sprintf("%d", speed), "N/A"})
		} else {
			percent, ok := curve.GetPointFromTemp(temp)
			if !ok {
				data = append(data, []string{fan.Name, fmt.Sprintf("%d", speed), "N/A"})
			} else {
				data = append(data, []string{fan.Name, fmt.Sprintf("%d", speed), fmt.Sprintf("%.0f%%", percent)})
			}
		}
	}
	printTable(headers, data)
	return 0
}

func Reload(confFile string, constants map[string]string) int {
	conn, err := net.Dial("unix", common.SOCKFILE)
	if err != nil {
		fmt.Println("Error connecting to socket:", err)
		return 1
	}
	defer conn.Close()
	message := "reload" + "!" + confFile

	if len(constants) > 0 {
		message += "!"
		for k, v := range constants {
			message += fmt.Sprintf("%s=%s\n", k, v)
			message = message[:len(message)-1]
		}
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to socket:", err)
		return 1
	}
	conn.(*net.UnixConn).CloseWrite()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from socket:", err)
		return 1
	}
	if string(buf[:n]) != "OK" {
		if len(buf) > 6 && string(buf[:6]) == "ERROR!" {
			fmt.Println("Error reloading config:", string(buf[6:n]))
		} else {
			fmt.Println("Error reloading config")
		}
		return 1
	}
	return 0
}
