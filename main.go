package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/marzeq/windigo/cli"
	"github.com/marzeq/windigo/common"
	"github.com/marzeq/windigo/daemon"
)

func main() {
	execName := filepath.Base(os.Args[0])

	var args struct {
		Daemon     bool   `arg:"-d,--daemon" help:"run as daemon"`
		ConfigFile string `arg:"-c,--config" help:"path to config file" default:"/etc/windigo/config.conf"`
		Reload     bool   `arg:"-r,--reload" help:"reload config file"`
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

	if args.Version {
		fmt.Printf("windigo version %s\n", common.VERSION)
		return
	}

	if execName == "windigod" || execName == "windigo-daemon" || args.Daemon {
		os.Exit(daemon.RunDaemon(args.ConfigFile))
	} else if args.Reload {
		os.Exit(cli.Reload())
	} else {
		os.Exit(cli.RunCli(args.ConfigFile))
	}
}
