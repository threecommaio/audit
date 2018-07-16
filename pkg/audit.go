package audit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

const notAvailable = "not available"

// Create a human readable json file of data
func Create(stdOut bool) {
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
	}
	b, _ := json.MarshalIndent(jsonData, "", "  ")
	hostname, _ := os.Hostname()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	if stdOut {
		fmt.Println(string(b))
	} else {
		auditFile := fmt.Sprintf("audit-%s-%s.json", hostname, timestamp)
		err := ioutil.WriteFile(auditFile, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Generated audit file: %s\n", auditFile)
	}
}
