package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
  "github.com/influxdata/influxdb/client/v2"
	"github.com/shirou/gopsutil/host"
)

func openClient(username string, password string, influxAddress string) (client.Client, error) {
	fmt.Println("Importing data points.")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxAddress,
		Username: username,
		Password: password,
	})

	return c, err
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func buildDataPoints(cpu_busy float64, cpu_busy_ticks float64, cpu_total_ticks float64) client.BatchPoints {
	batchPoints, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "environment",
		Precision: "s",
	})
	h, _ := host.Info()

	// v, _ := mem.VirtualMemory()
	// fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	tags := map[string]string{
		"hostname": h.Hostname,
	}
	fields := map[string]interface{}{
		"cpu_busy":        cpu_busy,
		"cpu_busy_ticks":  cpu_busy_ticks,
		"cpu_total_ticks": cpu_total_ticks,
	}

	pt, _ := client.NewPoint("cpu_utilization", tags, fields, time.Now())
	batchPoints.AddPoint(pt)
	return batchPoints
}

func main() {
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	username := os.Getenv("username")
	password := os.Getenv("password")
	addr := os.Getenv("address")
	fmt.Println(username, password, addr)
	client, err := openClient(username, password, addr)
	if err != nil {
		log.Fatalln("Error: ", err)
		return
	}

	batchPoints := buildDataPoints(cpuUsage, totalTicks-idleTicks, totalTicks)

	err = client.Write(batchPoints)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
}
