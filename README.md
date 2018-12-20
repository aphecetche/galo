[![GoDoc](https://godoc.org/github.com/aphecetche/galo?status.svg)](https://godoc.org/github.com/aphecetche/galo)
[![Build Status](https://travis-ci.org/aphecetche/galo.svg?branch=master)](https://travis-ci.org/aphecetche/galo)


Utilities to inspect/play with Run2/3 Alice Muon Intermediate Data Formats.

For the moment it's really a WIP and should be taken "as is", with no guarantee whatsover.

Nevertheless, for the moment the `galo` executable provides mainly the `cluster` command.

```
> galo cluster

Various operations related to clusters

Usage:
  galo cluster [command]

Available Commands:
  convert     Convert cluster(s) from one format to another
  create      Generate clusters
  dump        Dump clusters
  plot        Plot clusters
  read        Just reads input file

Flags:
  -h, --help             help for cluster
  -m, --max-events int   Maximum number of events to process (default 100000000)

Use "galo cluster [command] --help" for more information about a command.
```

The `convert` command so far can only get from clusters in YAML format to clusters in HTML(SVG really) format.

The `create` command generate a cluster at given position within a detection element. Used to observe the charge distribution, and to feed known clusters to the clustering. 
