# windigo

`windigo` is a lightweight Go-based daemon for Linux that controls fan speeds based on sensor input.

## Features

- Controls fan speeds based on sensor readings
- Configurable fan control curves and sensor-fan mappings
- Lightweight

## Installation

### Using your package manager

- Arch Linux: [from the AUR](https://aur.archlinux.org/packages/windigo)

### With `just install`

1. Clone the repository

```bash
git clone https://github.com/marzeq/windigo
cd windigo
```

2. Insall [`just`](https://github.com/casey/just) and run:

```bash
just install
# to uninstall, run
just uninstall
```

### Manually

1. Clone the repository

```bash
git clone https://github.com/marzeq/windigo
cd windigo
```

2. Compile manually

```bash
go build . -o windigo
ln -sf windigo windigod
```

3. Install binaries and service file

```bash
sudo install -Dm755 windigo /usr/local/bin/windigo
sudo install -Dm755 windigod /usr/local/bin/windigod
sudo install -Dm644 windigod.service /etc/systemd/system/windigod.service
```

## Usage

Create a config file in `/etc/windigo/config.conf` that follows the guidelines in [`the example`](./config.example.conf).

Then, enable and start the service

```bash
sudo systemctl enable windigod
sudo systemctl start windigod
```

Then, analyse the output of `systemctl status windigod` to see if all went well (it should tell you if it hasn't).

You can also peek at defined sensor temperatures and fan RPMs with the `windigo` command.

## License

[MIT](./LICENSE)
