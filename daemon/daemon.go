package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marzeq/windigo/common"
	"github.com/marzeq/windigo/config"
)

func RunDaemon(confFile string) int {
	cfg, err := config.ReadConfig(confFile)
	if err != nil {
		log.Printf("%v", err)
		return 1
	}

	if os.Geteuid() != 0 {
		fmt.Println("windogod must be run as root")
		return 1
	}

	if _, err := os.Stat(common.LOCKFILE); err == nil {
		fmt.Println("windigod is already running")
		return 1
	} else if os.IsNotExist(err) {
		if err := os.WriteFile(common.LOCKFILE, []byte(""), 0644); err != nil {
			fmt.Println("Error creating lock file:", err)
			return 1
		}
	} else {
		fmt.Println("Error accessing lock file:", err)
		return 1
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-c
		os.Remove(common.LOCKFILE)
		os.Remove(common.SOCKFILE)
		os.Exit(0)
	}()

	os.Remove(common.SOCKFILE)
	ln, err := net.Listen("unix", common.SOCKFILE)
	if err != nil {
		fmt.Println("Socket listen error:", err)
		return 1
	}
	defer ln.Close()

	if err := os.Chmod(common.SOCKFILE, 0666); err != nil {
		fmt.Println("Error setting socket permissions:", err)
		os.Remove(common.SOCKFILE)
		os.Exit(1)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Socket accept error:", err)
				continue
			}
			go func(c net.Conn) {
				buf := make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil {
					log.Println("Socket read error:", err)
					return
				}
				if string(buf[:n]) == "reload" {
					log.Println("Reloading config")
					newCfg, err := config.ReadConfig(confFile)
					if err != nil {
						log.Println("Error reloading config:", err)
						return
					}
					cfg = newCfg
					for _, fan := range cfg.Fans {
						if err := fan.SetManualMode(); err != nil {
							log.Printf("%v", err)
							continue
						}
					}
					log.Println("Config reloaded")
					conn.Write([]byte("OK"))
				} else {
					log.Println("Unknown command:", string(buf[:n]))
				}

				if err := conn.Close(); err != nil {
					log.Println("Socket close error:", err)
				}
			}(conn)
		}
	}()

	if cfg.Path != common.DEFAULT_CONFIG {
		log.Println("Using custom config file:", cfg.Path)
	}

	log.Println("Starting windigo daemon")
	fatal := false
	for _, fan := range cfg.Fans {
		if err := fan.SetManualMode(); err != nil {
			log.Printf("%s", err.Error())
			fatal = true
			break
		}
	}

	if fatal {
		log.Println("Fatal error, exiting")
		return 1
	}

	for {
		MainBody(time.Now().Unix(), cfg)
		time.Sleep(1 * time.Second)
	}
}

func MainBody(t0 int64, config config.Config) {
	dt := time.Now().Unix() - t0
	for curveName, curve := range config.Curves {
		if dt%curve.Period != 0 {
			continue
		}

		temp, err := curve.GetAggregateTemp(config.Sensors)
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		if curve.PreviousTemp != 0 && abs(temp-curve.PreviousTemp) < curve.Hysteresis {
			continue
		}
		percent, ok := curve.GetPointFromTemp(temp)
		if !ok {
			log.Printf("Error getting point from curve '%s'", curveName)
			continue
		}
		curve.PreviousTemp = temp

		setFor := ""
		for fanName, fan := range config.Fans {
			if fan.Curve != curveName {
				continue
			}

			if err := fan.SetSpeed(percent); err != nil {
				log.Printf("Error setting speed for fan '%s': %v", fanName, err)
				continue
			}
			setFor += "'" + fanName + "', "
		}

		if setFor != "" && percent != curve.PreviousPercent {
			setFor = setFor[:len(setFor)-2]
			log.Printf("Set %s to %.0f%%", setFor, percent)
		}
		curve.PreviousPercent = percent
		curve.PreviousTemp = temp
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
