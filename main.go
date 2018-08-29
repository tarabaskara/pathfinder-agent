package main

import (
	"os"

	"github.com/BaritoLog/go-boilerplate/envkit"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	Name                            = "Pathfinder Agent"
	Version                         = "0.0.1"
	DefaultLXDSocketPath            = "/var/snap/lxd/common/lxd/unix.socket"
	DefaultPfCluster                = "default"
	DefaultPfServerAddr             = "http://192.168.33.1:3000"
	DefaultPfListContainersPath     = "api/v1/node/containers/scheduled"
	DefaultPfMarkProvisionedPath    = "api/v1/node/containers/mark_provisioned"
	DefaultPfMarkProvisionErrorPath = "api/v1/node/containers/mark_provision_error"
	DefaultPfMarkDeletedPath        = "api/v1/node/containers/mark_deleted"

	EnvLXDSocketPath            = "LXD_SOCKET_PATH"
	EnvPfCluster                = "PF_CLUSTER"
	EnvPfServerAddr             = "PF_SERVER_ADDR"
	EnvPfListContainersPath     = "PF_LIST_CONTAINERS_PATH"
	EnvPfMarkProvisionedPath    = "PF_MARK_PROVISIONED_PATH"
	EnvPfMarkProvisionErrorPath = "PF_MARK_PROVISION_ERROR_PATH"
	EnvPfMarkDeletedPath        = "PF_MARK_DELETED_PATH"
)

var (
	LXDSocketPath string
	PfCluster     string
	PfServerAddr  string
	PfApiPath     map[string]string
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	// Initialize config vars
	LXDSocketPath, _ = envkit.GetString(EnvLXDSocketPath, DefaultLXDSocketPath)
	PfCluster, _ = envkit.GetString(EnvPfCluster, DefaultPfCluster)
	PfServerAddr, _ = envkit.GetString(EnvPfServerAddr, DefaultPfServerAddr)
	PfApiPath = make(map[string]string)
	PfApiPath["ListContainers"], _ = envkit.GetString(EnvPfListContainersPath, DefaultPfListContainersPath)
	PfApiPath["MarkProvisioned"], _ = envkit.GetString(EnvPfMarkProvisionedPath, DefaultPfMarkProvisionedPath)
	PfApiPath["MarkProvisionError"], _ = envkit.GetString(EnvPfMarkProvisionErrorPath, DefaultPfMarkProvisionErrorPath)
	PfApiPath["MarkDeleted"], _ = envkit.GetString(EnvPfMarkDeletedPath, DefaultPfMarkDeletedPath)
}

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Agent for Pathfinder container manager",
		Version: Version,
		Action:  CmdAgent,
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Enable verbose mode",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
