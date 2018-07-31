package audit

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getSysctl return output from sysctl command
func getSysctl() map[string]string {
	out, err := exec.Command("sysctl", "-a").Output()
	kv := make(map[string]string)
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, pair := range lines {
			if strings.ContainsAny(pair, ":") {
				z := strings.Split(pair, ":")
				key := strings.TrimSpace(z[0])
				kv[key] = strings.TrimSpace(z[1])
			}

			if strings.ContainsAny(pair, "=") {
				z := strings.Split(pair, "=")
				key := strings.TrimSpace(z[0])
				kv[key] = strings.TrimSpace(z[1])
			}
		}
	}
	return kv
}

// readCommand captures the output of a command
func readCommand(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	if err != nil {
		return notAvailable
	}
	return strings.TrimSpace(string(out))
}

// readFile captures the contents of a file
func readFile(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return notAvailable
	}
	return strings.TrimSpace(string(dat))
}

// readFileWithPaths reads a file by name and list of paths to search
func readFileWithPaths(filename string, paths []string) string {
	for _, p := range paths {
		var pathFilename = filepath.Join(p, filename)
		if _, err := os.Stat(pathFilename); err == nil {
			return readFile(pathFilename)
		}
	}
	return notAvailable
}

// delimitedData splits data by a delimiter
func delimitedData(delimiter string, data string) map[string]string {
	kv := make(map[string]string)

	if strings.Contains(data, notAvailable) {
		return kv
	}

	lines := strings.Split(data, "\n")
	for _, pair := range lines {
		if strings.ContainsAny(pair, delimiter) {
			z := strings.Split(pair, delimiter)
			key := strings.TrimSpace(z[0])
			value := strings.TrimSpace(z[1])
			kv[key] = value
		}
	}
	return kv
}

// getRelease checks with linux distro it is
func getRelease() string {
	rels := []string{
		"/etc/SuSE-release", "/etc/redhat-release", "/etc/redhat_version",
		"/etc/fedora-release", "/etc/slackware-release",
		"/etc/slackware-version", "/etc/debian_release", "/etc/debian_version",
		"/etc/os-release", "/etc/mandrake-release", "/etc/yellowdog-release",
		"/etc/sun-release", "/etc/release", "/etc/gentoo-release",
		"/etc/system-release", "/etc/lsb-release",
	}
	for _, path := range rels {
		if _, err := os.Stat(path); err == nil {
			return readFile(path)
		}
	}
	return notAvailable
}

// getScheduler handles capturing scheduler data for each block device
func getScheduler() map[string]string {
	kv := make(map[string]string)
	files, err := ioutil.ReadDir("/sys/block")
	if err != nil {
		return kv
	}
	for _, f := range files {
		block := f.Name()
		path := "/sys/block/" + block + "/queue/scheduler"
		if _, err := os.Stat(path); err == nil {
			kv[block] = readFile(path)
		}
	}
	return kv
}
