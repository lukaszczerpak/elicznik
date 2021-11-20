# elicznik-sync

Small app that synchronises Tauron's eLicznik data with InfluxDB bucket.  
It validates data provided by Tauron to ensure data is complete before writing to the database.  
Tauron is only queried for missing days (for which there is no data in the database).  
Data from Tauron is fetched in batches using one month time period.

## Configuration

Create a configuration file with necessary information:

```yaml
influxdb:
  bucket: your bucket name
  org: organization name
  url: https://your.influxdb.instance
  token: token with read & write permissions
tauron:
  login: your@login.net
  password: password
  smart_nr: 0123456789
  start_date: 2021-10-01
```

> Please note that `start_date` should be the first day after smart meter device installation.
> This is just to ensure that on the `start_date` day, measurements data is complete. 

## Run

Tauron updates their database each day around 9am-10am. Therefore, the app
should be executed at that time and preferably once a day.  
In case of any delay on Tauron side, missing data will be fetch on the following day
(based on the last measurement date stored in InfluxDB).

```shell
❯ elicznik-sync 
time="2021-11-07T21:43:26+01:00" level=info msg="Fetching measurements from 2021-08-11 00:00:00 +0200 CEST to 2021-11-07 00:00:00 +0100 CET"
time="2021-11-07T21:43:26+01:00" level=info msg="Batch #0: from 2021-08-11 00:00:00 +0200 CEST to 2021-09-10 00:00:00 +0200 CEST"
time="2021-11-07T21:43:31+01:00" level=info msg="Writing measurements for 2021-08-11"
time="2021-11-07T21:43:31+01:00" level=info msg="Writing measurements for 2021-08-12"
time="2021-11-07T21:43:31+01:00" level=info msg="Writing measurements for 2021-08-13"
...
time="2021-11-07T21:43:34+01:00" level=info msg="Writing measurements for 2021-09-10"
time="2021-11-07T21:43:34+01:00" level=info msg="Batch #1: from 2021-09-11 00:00:00 +0200 CEST to 2021-10-10 00:00:00 +0200 CEST"
time="2021-11-07T21:43:37+01:00" level=info msg="Writing measurements for 2021-09-11"
time="2021-11-07T21:43:37+01:00" level=info msg="Writing measurements for 2021-09-12"
time="2021-11-07T21:43:37+01:00" level=info msg="Writing measurements for 2021-09-13"
...
```

Optionally you can use `--config` flag to point to your configuration file in case it's not in current folder or filename is not `elicznik-sync.yaml`:

```shell
❯ elicznik-sync --help
This program synchronizes data from Tauron with InfluxDB.
It automatically detects last sync and pulls only necessary
amount of data.

Usage:
  elicznik-sync [flags]

Flags:
      --config string   config file (default is elicznik-sync.yaml)
  -h, --help            help for elicznik-sync

```

## Running from Cron

1. Create a configuration file, ie. `/usr/local/etc/elicznik-sync.yaml`

2. Add a new job to your cron configuration:
    ```
    # elicznik-sync
    10 9,12,16 * * * docker run --rm -v /usr/local/etc/elicznik-sync.yaml:/elicznik-sync.yaml -e TZ="Europe/Warsaw" elicznik-sync:latest elicznik-sync >> /var/log/elicznik-sync.log 2> >(tee -a /var/log/elicznik-sync.log >&2)
    ```
