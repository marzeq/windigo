# windigo

`windigo` is a lightweight Go-based daemon for Linux that controls fan speeds based on sensor input.

## Features

- Controls fan speeds based on sensor readings
- Configurable fan control curves and sensor-fan mappings
- Lightweight

## Installation

### Precompiled

- Arch Linux: [from the AUR](https://aur.archlinux.org/packages/windigo) *OR* manually installing from the [releases page](https://github.com/marzeq/windigo/releases/latest) with `pacman -U`
- Other: Get the tarball from the [releases page](https://github.com/marzeq/windigo/releases/latest) package

### From source - with `just install`

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

### From source - manually

1. Clone the repository

```bash
git clone https://github.com/marzeq/windigo
cd windigo
```

2. Compile manually

```bash
go build . -o windigo
ln -s windigo windigod
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

### Hot reloading

Use the `windigo` command to reload the configuration without restarting the service:

```bash 
wingigo --reload
# or
windigo -r
```

### Setting default constants from command line

You can set default constants for the config file from the command line when starting the service/reloading it:

```bash
windigod -C some_constant=some_value
# or when reloading
windigo -r -C some_constant=some_value
```

One example that comes to mind is setting up a `$curve` constant like so:

```mconf
$curve = $curve ? your_curve
# feature of mconf, where constants can have backup values. see more at https://github.com/marzeq/mconf#default-values
```

This way, `$curve` will take on the default value of `your_curve`, but you can override it. You can then do whatever you want with that newfound ability,
like making a script that uses a differtent curve at night-time, or something like that.

## License

[MIT](./LICENSE)
