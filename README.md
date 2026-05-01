# tery

A lightweight battery monitoring daemon for Linux that triggers custom commands based on battery status and charge levels.

## Features

- Monitors battery status changes (Charging, Discharging, Not charging, Full)
- Triggers notifications based on battery charge ranges
- Configurable via JSON config file
- Runs silently in the background with minimal resource usage

## Installation

### Build from source
```bash
git clone https://github.com/yourusername/tery.git
cd tery
go build -o tery
sudo mv tery /usr/bin/tery
```

## Configuration

On first run, tery creates a default config file at `~/.config/tery/config.json`:

```json
{
    "status": [
        {
            "state": "Charging",
            "command": "notify-send -i ac-adapter-symbolic Charging"
        },
        {
            "state": "Discharging",
            "command": "notify-send -i ac-adapter-symbolic Discharging"
        },
        {
            "state": "Not charging",
            "command": "notify-send -i ac-adapter-symbolic \"Not Charging\""
        },
        {
            "state": "Full",
            "command": "notify-send -i ac-adapter-symbolic Full"
        }
    ],
    "range": [
        {
            "upper": 20,
            "lower": 10,
            "command": "notify-send -i dialog-error \"Battery Extremely Low\" \"Plug in Charger Immediately\""
        },
        {
            "upper": 10,
            "lower": 0,
            "command": "notify-send -i battery-low \"Battery Low\" \"Plug in Charger\""
        }
    ]
}
```

### Config Options

**Status rules** - Triggered when battery status changes:
- `state` - Battery state to match (Charging, Discharging, Not charging, Full)
- `command` - Shell command to execute

**Range rules** - Triggered when battery percentage falls within a range:
- `upper` - Upper bound (inclusive)
- `lower` - Lower bound (inclusive)
- `command` - Shell command to execute

## Usage

```bash
./tery
```

The daemon will run continuously, checking battery status every 100ms.

## Requirements

- Linux system with `/sys/class/power_supply/BAT0/` present
- Go 1.26+ (for building from source)
