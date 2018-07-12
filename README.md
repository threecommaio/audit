# ThreeComma Audit Toolkit
This toolkit helps collect various sets of information from a host. This includes various components from the linux kernel (cpuinfo, sysctl, proc), running processes, versions of software, software configuration, etc. It's designed to produce a human readable audit file for further analysis.

# Building
```
$ make build
```

# Usage
```
$ ./audit --help
ThreeComma Audit Toolkit

This toolkit helps collect various sets of information from a host. This includes various components from the linux kernel (cpuinfo, sysctl, proc), running processes, versions of software, software configuration, etc.

It's designed to produce a human readable audit file for further analysis.

Usage:
  audit [flags]

Flags:
  -c, --console   print to console instead of file
  -h, --help      help for audit
```