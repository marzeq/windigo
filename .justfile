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

clean:
  rm -rf {{builddir}}

