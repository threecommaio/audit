package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const notAvailable = "not available"
const xClientToken = "X-Client-Token"
const auditURL = "https://audit.threecomma.io"

// Create a human readable json file of data
func Create(stdOut bool, clientToken string) {
	cassandraPaths := []string{
		"/etc/cassandra",
		"/etc/cassandra/conf",
		"/etc/dse/cassandra",
		"/etc/dse",
		"/usr/local/share/cassandra",
		"/usr/local/share/cassandra/conf",
		"/opt/cassandra",
		"/opt/cassandra/conf",
	}

	jsonData := Audit{
		Internal: map[string]string{
			"version":    Version,
			"buildTime":  BuildTime,
			"commitHash": CommitHash,
		},
		Sysctl: getSysctl(),
		Proc: Proc{
			Cpuinfo:        readFile("/proc/cpuinfo"),
			Cmdline:        readFile("/proc/cmdline"),
			NetSoftnetStat: readFile("/proc/net/softnet_stat"),
			Cgroups:        readFile("/proc/cgroups"),
			Uptime:         readFile("/proc/uptime"),
			Vmstat:         delimitedData(" ", readFile("/proc/vmstat")),
			Loadavg:        readFile("/proc/loadavg"),
			Zoneinfo:       readFile("/proc/zoneinfo"),
			Partitions:     readFile("/proc/partitions"),
			Version:        readFile("/proc/version"),
		},
		Dmesg: readCommand("dmesg"),
		THP: THP{
			Enabled: readFile("/sys/kernel/mm/transparent_hugepage/enabled"),
			Defrag:  readFile("/sys/kernel/mm/transparent_hugepage/defrag"),
		},
		Memory: readCommand("free", "-m"),
		Disk: Disk{
			Scheduler:  getScheduler(),
			Partitions: readCommand("df", "-h"),
			NumDisks:   readCommand("lsblk"),
		},
		Network: Network{
			Ifconfig: readCommand("ifconfig"),
			IP:       readCommand("ip", "addr", "show"),
			Netstat:  readCommand("netstat", "-an"),
			SS:       readCommand("ss", "-tan"),
		},
		Distro: Distro{
			Issue:   readFile("/etc/issue"),
			Release: getRelease(),
		},
		PowerMgmt: PowerMgmt{
			MaxCState: readFile("/sys/module/intel_idle/parameters/max_cstate"),
		},
		Cassandra: Cassandra{
			ConfigYaml:      readFileWithPaths("cassandra.yaml", cassandraPaths),
			Env:             readFileWithPaths("cassandra-env.sh", cassandraPaths),
			JvmOptions:      readFileWithPaths("jvm.options", cassandraPaths),
			RackProperties:  readFileWithPaths("cassandra-rackdc.properties", cassandraPaths),
			DseYaml:         readFileWithPaths("dse.yaml", cassandraPaths),
			NodetoolVersion: readCommand("nodetool", "version"),
			NodetoolStatus:  readCommand("nodetool", "status"),
			NodetoolInfo:    readCommand("nodetool", "info"),
		},
	}
	b, _ := json.MarshalIndent(jsonData, "", "  ")
	hostname, _ := os.Hostname()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if stdOut {
		// Output result to console
		fmt.Println(string(b))
	} else if clientToken != "" {
		// Upload to google cloud
		req, err := http.NewRequest("POST", auditURL, bytes.NewBuffer(b))
		req.Header.Set("X-Client-Token", clientToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			io.Copy(ioutil.Discard, resp.Body)
		} else {
			log.Fatal(resp)
		}
	} else {
		auditFile := fmt.Sprintf("audit-%s-%s.json", hostname, timestamp)
		err := ioutil.WriteFile(auditFile, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Generated audit file: %s\n", auditFile)
	}
}
