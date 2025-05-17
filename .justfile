default:
  @just --list

builddir := "build"
binname := "windigo"
daemonname := "windigod"

build:
  mkdir -p {{builddir}}
  go build -o {{builddir}}/{{binname}} main.go

build-daemon: build
  ln -sf {{binname}} {{builddir}}/{{daemonname}}

run: build
  {{builddir}}/{{binname}}

run-daemon: build-daemon
  sudo ./{{builddir}}/{{daemonname}}

install: build build-daemon
  sudo install -Dm755 {{builddir}}/{{binname}} /usr/local/bin/{{binname}}
  sudo install -Dm755 {{builddir}}/{{daemonname}} /usr/local/bin/{{daemonname}}
  sudo install -Dm644 {{daemonname}}.service /etc/systemd/system/{{daemonname}}.service

uninstall:
  sudo rm -f /usr/local/bin/{{binname}}
  sudo rm -f /usr/local/bin/{{daemonname}}
  sudo systemctl stop {{daemonname}}
  sudo systemctl disable {{daemonname}}
  sudo rm -f /etc/systemd/system/{{daemonname}}.service

clean:
  rm -rf {{builddir}}

