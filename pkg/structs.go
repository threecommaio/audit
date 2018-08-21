package audit

// Audit holds the key settings of the system
type Audit struct {
	Internal  map[string]string `json:"internal"`
	Sysctl    map[string]string `json:"sysctl"`
	Proc      Proc              `json:"proc"`
	Dmesg     string            `json:"dmesg"`
	THP       THP               `json:"transparent_huge_pages"`
	Memory    string            `json:"memory"`
	Disk      Disk              `json:"disk"`
	Network   Network           `json:"network"`
	Distro    Distro            `json:"distro"`
	PowerMgmt PowerMgmt         `json:"power_mgmt"`
	Cassandra Cassandra         `json:"cassandra"`
}

// Proc holds the key settings of proc interface
type Proc struct {
	Cpuinfo        string            `json:"cpuinfo"`
	Cmdline        string            `json:"cmdline"`
	NetSoftnetStat string            `json:"net/softnet_stat"`
	Cgroups        string            `json:"cgroups"`
	Uptime         string            `json:"uptime"`
	Vmstat         map[string]string `json:"vmstat"`
	Loadavg        string            `json:"loadavg"`
	Zoneinfo       string            `json:"zoneinfo"`
	Partitions     string            `json:"partitions"`
	Version        string            `json:"version"`
}

// THP handles the transparent huge pages
type THP struct {
	Enabled string `json:"enabled"`
	Defrag  string `json:"defrag"`
}

// Disk handles the disk subystem
type Disk struct {
	Scheduler  map[string]string `json:"scheduler"`
	NumDisks   string            `json:"number_of_disks"`
	Partitions string            `json:"partitions"`
}

// Network handles the networking settings
type Network struct {
	Ifconfig string `json:"ifconfig"`
	IP       string `json:"ip"`
	Netstat  string `json:"netstat"`
	SS       string `json:"ss"`
}

// Distro handles the linux distro settings
type Distro struct {
	Issue   string `json:"issue"`
	Release string `json:"release"`
}

// PowerMgmt handles the power management settings
type PowerMgmt struct {
	MaxCState string `json:"max_cstate"`
}

// Cassandra
type Cassandra struct {
	ConfigYaml      string `json:"cassandra_yaml"`
	Env             string `json:"cassandra_env"`
	JvmOptions      string `json:"jvm_options"`
	RackProperties  string `json:"rack_properties"`
	DseYaml         string `json:"dse_yaml"`
	NodetoolVersion string `json:"nodetool_version"`
	NodetoolStatus  string `json:"nodetool_status"`
	NodetoolInfo    string `json:"nodetool_info"`
	JavaVersion     string `json:"java_version"`
}
