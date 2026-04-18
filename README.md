# Power Consumption Metrics Tool

Queries HomeAssistant for a summary of your power usage over a period of time.
Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
The websocket API does provide a way to do this, which is what the frontend uses.

I got fed up of trying to figure out how to get the same data that the Energy dashboard shows, so I wrote this tool to do it for me.

This tool queries the websocket API to get the power usage data for each hour over a period of days and outputs the data in various formats

## Installation

Go

```bash
go install github.com/poolski/powertracker@latest
```

[Releases](https://github.com/poolski/powertracker/releases)

## Configuration

This tool requires a configuration file to be present at `~/.config/powertracker/config.yaml`. If one does not exist, it will ask for input and create it for you.
The only things this tool needs are the URL of your Home Assistant instance and a long-lived access token.

You can generate a long-lived access token by going to your Home Assistant instance, clicking on your profile picture in the bottom left, then clicking on "Long-Lived Access Tokens" at the bottom of the list and creating a new one.

## Usage

```bash
$ powertracker --help

Usage:
  powertracker [flags]

Flags:
  -c, --config string     config file (default "$HOME_DIR/.config/powertracker/config.yaml")
  -f, --csv-file string   the path of the CSV file to write to (default "results.csv")
  -d, --days int          number of days to compute power stats for (default 7)
  -h, --help              help for powertracker
  -i  --insecure          skip TLS verification
  -o, --output string     output format (text, table, csv)

```

> **Note:** I changed the default `--days` value from 30 to 7 in my fork, since I mostly care about the past week and 30 days felt like too much noise at a glance.

> **Note:** I also changed the default `--output` format to `table` instead of `text`, since the table view is much easier to read at a glance.

> **Note:** If you're running Home Assistant behind a self-signed certificate (e.g. on a local network with a custom CA), use `-i` / `--insecure` to skip TLS verification. Not recommended for production, but handy for a home lab setup.

## Example output

```bash
$ powertracker -d 7 # 7 days' worth of data

+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
|    0     |    1     |    2     |    3     |    4     |    5     |    6     |    7     |    8     |    9     |    10    |    11    |    12    |    13    |    14    |    15    |    16    |    17    |    18    |    19    |    20    |    21    |    22    |    23    |
+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
| 0.300000 | 0.326000 | 0.333000 | 0.298000 | 0
```
