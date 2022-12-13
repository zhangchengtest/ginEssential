package main

import (
	_ "expvar"
	"fmt"
	"ginEssential/config"
	"ginEssential/job"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/zhangchengtest/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	_ "net/http/pprof"
	"os"
	"time"
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

	// gorm配置
	gormConf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 初始化日志
	if file, err := os.OpenFile(config.Instance.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		if config.Instance.ShowSql {
			gormConf.Logger = logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
		}
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 连接数据库
	if err := sqls.Open(config.Instance.DB.Url, gormConf, config.Instance.DB.MaxIdleConns, config.Instance.DB.MaxOpenConns); err != nil {
		logrus.Error(err)
	}

	c := cron.New()
	spec := "0 */1 * * * ?"
	err := c.AddFunc(spec, func() {
		fmt.Println("cron times : ")
		job.SplitImageWithIterator()
	})
	if err != nil {
		fmt.Errorf("AddFunc error : %v", err)
		return
	}
	c.Start()

	r := gin.Default()
	r = CollectRoute(r)
	port := config.Instance.Port
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())

	sws := util.GetInstance()
	go sws.Run()

	var errs error

	return errs
}

//func InitConfig() {
//	workDir, _ := os.Getwd()
//	viper.SetConfigName("application")
//	viper.SetConfigType("yml")
//	viper.AddConfigPath(workDir + "/config")
//	err := viper.ReadInConfig()
//	if err != nil {
//		log.Errorf("read config failed ,err : %v", err)
//		panic("read config failed, err: " + err.Error())
//	}
//}
