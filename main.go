package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/alexflint/go-arg"
	"github.com/marzeq/windigo/config"
	"github.com/marzeq/windigo/daemon"
)

const VERSION = "0.1.2"

const (
	DEFAULT_CONFIG = "/etc/windigo/config.conf"
	LOCKFILE       = "/var/lock/windigod.lock"
)

func runDaemon(config config.Config) {
	if os.Geteuid() != 0 {
		fmt.Println("windogod must be run as root")
		os.Exit(1)
	}

	if _, err := os.Stat(LOCKFILE); err == nil {
		fmt.Println("windigod is already running")
		os.Exit(1)
	} else if os.IsNotExist(err) {
		if err := os.WriteFile(LOCKFILE, []byte(""), 0644); err != nil {
			fmt.Println("Error creating lock file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error accessing lock file:", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-c
		if err := os.Remove(LOCKFILE); err != nil {
			fmt.Println("Error removing lock file:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	if config.Path != DEFAULT_CONFIG {
		fmt.Println("Using custom config file:", config.Path)
	}

	daemon.Main(config)
	if err := os.Remove(LOCKFILE); err != nil {
		fmt.Println("Error removing lock file:", err)
		os.Exit(1)
	}
}

func runCli(config config.Config) {
	fmt.Println("windigo version", VERSION)
	fmt.Println()

	for _, sensor := range config.Sensors {
		temp, err := sensor.ReadTemperature()
		if err != nil {
			fmt.Printf("Error reading temperature: %s\n", err.Error())
			continue
		}
		fmt.Printf("%s\t%.1f °C\n", sensor.Name, temp)
	}
	fmt.Println()

	for _, fan := range config.Fans {
		speed, err := fan.ReadSpeed()
		if err != nil {
			fmt.Printf("Error reading fan speed: %s\n", err.Error())
			continue
		}
		curve := config.Curves[fan.Curve]
		temp, err := curve.GetAggregateTemp(config.Sensors)
		if err != nil {
			fmt.Printf("%s\t%d RPM\n", fan.Name, speed)
		} else {
			percent, ok := curve.GetPointFromTemp(temp)
			if !ok {
				fmt.Printf("%s\t%d RPM\n", fan.Name, speed)
			} else {
				fmt.Printf("%s\t%d RPM (%.0f%%)\n", fan.Name, speed, percent)
			}
		}
	}
}

func main() {
	execName := filepath.Base(os.Args[0])

	var args struct {
		Daemon     bool   `arg:"-d,--daemon" help:"run as daemon"`
		ConfigFile string `arg:"-c,--config" help:"path to config file" default:"/etc/windigo/config.conf"`
		Version    bool   `arg:"-v,--version" help:"print version and exit"`
	}

	arg.MustParse(&args)

	if _, err := os.Stat(args.ConfigFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config file not found:", args.ConfigFile)
			os.Exit(1)
		}
		if os.IsPermission(err) {
			fmt.Println("Permission denied to access config file:", args.ConfigFile)
			os.Exit(1)
		}
		fmt.Println("Error accessing config file:", err)
		os.Exit(1)
	}

	config, err := config.ReadConfig(args.ConfigFile)
	if err != nil {
		panic(err)
	}

	if args.Version {
		fmt.Printf("windigo version %s\n", VERSION)
		return
	}

	if execName == "windigod" || execName == "windigo-daemon" || args.Daemon {
		runDaemon(config)
		return
	} else {
		runCli(config)
		return
	}
}
