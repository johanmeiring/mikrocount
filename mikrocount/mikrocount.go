package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type entry struct {
	FromIP net.IP
	ToIP   net.IP
	Bytes  uint
}

func main() {
	influxDBURL := flag.String("influxurl", "http://influxdb:8086", "InfluxDB Host")
	localCIDR := flag.String("localcidr", "192.168.0.0/16", "Local CIDR (e.g. 192.168.0.0/16)")
	mikrotikAddr := flag.String("mikrotikaddr", "", "Mikrotik IP/hostname")

	flag.Parse()

	conf := client.HTTPConfig{Addr: *influxDBURL}
	c, _ := client.NewHTTPClient(conf)
	q := client.NewQuery("CREATE DATABASE mikrocount", "", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		log.Println(response.Results)
	}
	defer c.Close()

	_, ipnet, _ := net.ParseCIDR(*localCIDR)
	dataChan := make(chan []entry)

	for {
		select {
		case <-time.After(time.Second * time.Duration(15)):
			go getData(*mikrotikAddr, dataChan)
		case e := <-dataChan:
			go recordEntries(e, ipnet, c)
		}
	}
}

func getData(mikrotikAddr string, dataChan chan []entry) {
	var entries []entry

	resp, err := http.Get(fmt.Sprintf("http://%s/accounting/ip.cgi", mikrotikAddr))
	if err != nil {
		log.Printf("Error fetching data from Mikrotik: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading data: %s", err)
	}

	lines := strings.Split(string(body), "\n")
	log.Printf("About to process %d results", len(lines)-1)

	for _, l := range lines {
		if l == "" {
			break
		}
		cols := strings.Split(l, " ")
		b, _ := strconv.Atoi(cols[2])
		e := entry{
			FromIP: net.ParseIP(cols[0]),
			ToIP:   net.ParseIP(cols[1]),
			Bytes:  uint(b),
		}
		entries = append(entries, e)
	}

	dataChan <- entries
}

func recordEntries(entries []entry, ipnet *net.IPNet, c client.Client) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "mikrocount",
		Precision: "us",
	})

	for _, e := range entries {
		var ip, direction string
		if ipnet.Contains(e.FromIP) {
			ip = e.FromIP.String()
			direction = "Upload"
		} else if ipnet.Contains(e.ToIP) {
			ip = e.ToIP.String()
			direction = "Download"
		} else {
			log.Printf("Weirdness! From: %s :: To: %s", e.FromIP.String(), e.ToIP.String())
			return
		}

		tags := map[string]string{
			"ip":        ip,
			"direction": direction,
		}

		fields := map[string]interface{}{
			"bytes": e.Bytes,
		}
		pt, err := client.NewPoint("usage", tags, fields, time.Now())
		if err != nil {
			log.Println("Error: ", err.Error())
		}
		bp.AddPoint(pt)
	}

	err := c.Write(bp)
	if err != nil {
		log.Println("Error: ", err.Error())
	}
}
