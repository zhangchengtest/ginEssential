package main

import (
	_ "expvar"
	"fmt"
	"ginEssential/common"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/spf13/viper"
	_ "net/http/pprof"
	"os"
)

const (
	adjustFDLimitKwd          = "manage-fdlimit"
	enableGCKwd               = "enable-gc"
	initOptionKwd             = "init"
	initConfigOptionKwd       = "init-config"
	initProfileOptionKwd      = "init-profile"
	ipfsMountKwd              = "mount-ipfs"
	ipnsMountKwd              = "mount-ipns"
	migrateKwd                = "migrate"
	mountKwd                  = "mount"
	offlineKwd                = "offline" // global option
	routingOptionKwd          = "routing"
	routingOptionSupernodeKwd = "supernode"
	routingOptionDHTClientKwd = "dhtclient"
	routingOptionDHTKwd       = "dht"
	routingOptionDHTServerKwd = "dhtserver"
	routingOptionNoneKwd      = "none"
	routingOptionDefaultKwd   = "default"
	unencryptTransportKwd     = "disable-transport-encryption"
	unrestrictedApiAccessKwd  = "unrestricted-api"
	writableKwd               = "writable"
	enablePubSubKwd           = "enable-pubsub-experiment"
	enableIPNSPubSubKwd       = "enable-namesys-pubsub"
	enableMultiplexKwd        = "enable-mplex-experiment"
	agentVersionSuffix        = "agent-version-suffix"
	// apiAddrKwd    = "address-api"
	// swarmAddrKwd  = "address-swarm"
)

var daemonCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Run a network-connected IPFS node.",
		ShortDescription: `
'ipfs daemon' runs a persistent ipfs daemon that can serve commands
over the network. Most applications that use IPFS will do so by
communicating with a daemon over the HTTP API. While the daemon is
running, calls to 'ipfs' commands will be sent over the network to
the daemon.
`,
		LongDescription: `
The daemon will start listening on ports on the network, which are
documented in (and can be modified through) 'ipfs config Addresses'.

`,
	},

	Options: []cmds.Option{
		cmds.BoolOption(initOptionKwd, "Initialize ipfs with default settings if not already initialized"),
		cmds.StringOption(initConfigOptionKwd, "Path to existing configuration file to be loaded during --init"),
		cmds.StringOption(initProfileOptionKwd, "Configuration profiles to apply for --init. See ipfs init --help for more"),
		cmds.StringOption(routingOptionKwd, "Overrides the routing option").WithDefault(routingOptionDefaultKwd),
		cmds.BoolOption(mountKwd, "Mounts IPFS to the filesystem"),
		cmds.BoolOption(writableKwd, "Enable writing objects (with POST, PUT and DELETE)"),
		cmds.StringOption(ipfsMountKwd, "Path to the mountpoint for IPFS (if using --mount). Defaults to config setting."),
		cmds.StringOption(ipnsMountKwd, "Path to the mountpoint for IPNS (if using --mount). Defaults to config setting."),
		cmds.BoolOption(unrestrictedApiAccessKwd, "Allow API access to unlisted hashes"),
		cmds.BoolOption(unencryptTransportKwd, "Disable transport encryption (for debugging protocols)"),
		cmds.BoolOption(enableGCKwd, "Enable automatic periodic repo garbage collection"),
		cmds.BoolOption(adjustFDLimitKwd, "Check and raise file descriptor limits if needed").WithDefault(true),
		cmds.BoolOption(migrateKwd, "If true, assume yes at the migrate prompt. If false, assume no."),
		cmds.BoolOption(enablePubSubKwd, "Enable experimental pubsub feature. Overrides Pubsub.Enabled config."),
		cmds.BoolOption(enableIPNSPubSubKwd, "Enable IPNS over pubsub. Implicitly enables pubsub, overrides Ipns.UsePubsub config."),
		cmds.BoolOption(enableMultiplexKwd, "DEPRECATED"),
		cmds.StringOption(agentVersionSuffix, "Optional suffix to the AgentVersion presented by `ipfs id` and also advertised through BitSwap."),

		// TODO: add way to override addresses. tricky part: updating the config if also --init.
		// cmds.StringOption(apiAddrKwd, "Address for the daemon rpc API (overrides config)"),
		// cmds.StringOption(swarmAddrKwd, "Address for the swarm socket (overrides config)"),
	},
	Subcommands: map[string]*cmds.Command{},
	NoRemote:    true,
	Run:         daemonFunc,
}


func daemonFunc(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) (_err error) {

	// let the user know we're going.
	fmt.Printf("Initializing daemon...\n")

	InitConfig()
	common.InitDB()

	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())



	// collect long-running errors and block for shutdown
	// TODO(cryptix): our fuse currently doesn't follow this pattern for graceful shutdown
	var errs error

	return errs
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("read config failed ,err : %v", err)
		panic("read config failed, err: " + err.Error())
	}
	sws := util.GetInstance()
	go sws.Run()
}
