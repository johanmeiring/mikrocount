# Mikrocount

[![Go Report Card](https://goreportcard.com/badge/github.com/johanmeiring/mikrocount)](https://goreportcard.com/report/github.com/johanmeiring/mikrocount) [![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-round)](https://github.com/johanmeiring/mikrocount/blob/master/LICENSE) [![Build Status](https://travis-ci.org/johanmeiring/mikrocount.svg?branch=master)](https://travis-ci.org/johanmeiring/mikrocount)

Mikrocount is a suite of tools that aid in visualising external network activity that goes through a Mikrotik router.  It consists of:
* `mikrocount`: The application hosted in this repo, which periodically fetches accounting data from the Mikrotik router and feeds it to Influxdb.
* `influxdb`: Time-series data store.  Mikrocount stores data in a database named `mikrocount`, in a single measurement named `usage`, with a field of `bytes` and tags of `direction` (Upload or Download) and `ip`.
* `grafana`: Awesome tool for visualising the data stored in Influxdb.

## Installing
### Configuring the Mikrotik router
Your Mikrotik router will need to be configured properly in order to expose accounting data to the local network.  Run the following commands in your router's terminal:

```
/ip accounting
set enabled=yes threshold=2000
/ip accounting web-access
set accessible-via-web=yes
```

Note:
* If your network is very busy, you should set the threshold to a higher value.  Mikrocount queries the Mikrotik router for data every 15 seconds.
* It is strongly recommended that access to your router from the internet over HTTP/HTTPS be blocked by the firewall.

### Running the Mikrocount suite
Requirements:
* A computer capable of running Docker, with Docker and docker-compose already installed.
* The computer should preferably be on 24/7.

#### docker-compose config
1. On the computer on which you want Mikrocount to run, create a file named `docker-compose.yml`, in a directory named `mikrocount`.
1. Copy the contents of the file `docker-compose-example.yml` in this repository, and paste them into the newly created `docker-compose.yml` file.
1. Change the value after `-mikrotikaddr` in the file to the LAN IP address of your Mikrotik router.
1. In a terminal session, in the `mikrocount` directory that was created, run `docker-compose up -d` (note: some setups may require you to run the aforementioned command using `sudo`).
1. After all images have been downloaded and containers are running, wait a minute or so before proceeding.

#### Configure Grafana
1. Open `http://<ip of machine running mikrocount>:3000` in a browser.  You should be presented with a Grafana login page.
1. Login using username/password `admin/admin`.
1. Click "Add data source".
1. Fill in the following details:
    ```
    Name: influxdb
    Type: InfluxDB
    Http settings
    Url: http://influxdb:8086
    Access: Proxy
    InfluxDB Details
    Database: mikrocount
    (leave everything else as is)
    ```
1. Click "Add" and verify that you get the message "Data source is working" at the bottom of the page.
1. Click the Grafana logo (top-left corner of the page), hover over "Dashboard" and click "Import".
1. Click "Upload .json file" and select `dash-example.json` provided in this repo.
1. Select `influxdb` as your data source, and click "Import".
1. You should already be seeing graphs and tables with some metrics, and from here on in you are welcome to customise as needed.

#### Stopping and starting the suite
* To stop everything: `docker-compose stop`
* To start everything: `docker-compose start`
* To delete everything except persisted data: `docker-compose down`
* To start everything again: `docker-compose up -d`
* To delete persisted data: `docker-compose down --volumes`

## Development
In order to contribute to the development of this project, or to fork and make your own changes, you'll need the following:

* Go
* [Dep](https://github.com/golang/dep)

To install the required dependencies, run `dep ensure`, or `make deps`.

Using `docker-compose-dev.yml` as a reference is recommended for creating a docker-compose config for your development environment.

## License
This software is distributed under the MIT License.  See the LICENSE file for more details.

## Donations
Donations are very welcome, and can be made to the following addresses:
* BTC: 1AWHJcUBha35FnuuWat9urRW2FNc4ftztv
* ETH: 0xAF1Aac4c40446F4C46e55614F14d9b32d712ECBc
