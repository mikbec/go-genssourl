//
// got idea from:
//   https://github.com/eschao/config

package main

import (
	"fmt"
	"github.com/eschao/config"
	"github.com/kr/pretty"
	"log"
)

// cmdline options
type CmdLineOptions struct {
	OptCfgSvcWebTcp         string `default:""     cli:"wtcp The HTTP server name/address[:port] to use."`
	OptCfgSvcWebTcpCertFile string `default:""     cli:"wtcpcert For HTTPS the cert.pem file tu use"`
	OptCfgSvcWebTcpKeyFile  string `default:""     cli:"wtcpkey For HTTPS the key.pem file tu use."`
	OptCfgSvcFcgiTcp        string `default:""     cli:"ftcp The FCGI server TCP name/address:port to use."`
	OptCfgSvcFcgiUnix       string `default:""     cli:"funix The FCGI server Unix socket name to use."`
	OptCfgSvcFcgiStdIO      bool   `default:"true" cli:"fstdio Use the FCGI server listen on standard I/O."`

	OptCfgFile string `default:""   cli:"cfgfile The config file to use."`
	OptDebug   int    `default:"0"  cli:"debug The debug level to use, 0 means no debug."`
	OptHelp    bool   `cli:"help Print this message and exit."`
}

// program log stuff
type Log struct {
	Path  string `default:"logs"`
	Level string `default:"debug"`
}

// web and service context
type WebCtx struct {
	// // set attributes
	// server_protocol := "https"
	// server_host := "localhost"
	// server_port := "5677"
	// server_context := "bla"
	// url_attr_username_key := "user"
	// url_attr_timestamp_key := "ts"
	// url_attr_hash_key := "hash"
	// url_attr_id_key := "id"
	//
	// username_val := "user1"
	// timestamp_val := "2023-11-23T08:15:32Z"
	// pub_pem_file := "testdata/test.crt.pem"
	// id_val := "test"

	DstServerProtocol string `default:"https"     json:"dstServerProtocol" yaml:"dstServerProtocol"`
	DstServerHost     string `default:"localhost" json:"dstServerHost" yaml:"dstServerHost"`
	DstServerPort     string `default:""          json:"dstServerPort" yaml:"dstServerPort"`
	DstServerCtx      string `default:""          json:"dstServerCtx" yaml:"dstServerCtx"`

	DstAttrKeyUsername  string `default:"user"      json:"dstAttrKeyUsername" yaml:"dstAttrKeyUsername"`
	DstAttrKeyTimestamp string `default:"ts"        json:"dstAttrKeyTimestamp" yaml:"dstAttrKeyTimestamp"`
	DstAttrKeyHash      string `default:"hash"      json:"dstAttrKeyHash" yaml:"dstAttrKeyHash"`
	DstAttrKeyId        string `default:"id"        json:"dstAttrKeyId" yaml:"dstAttrKeyId"`

	DstAttrValUsername  string `default:""          json:"dstAttrValUsername" yaml:"dstAttrValUsername"`
	DstAttrValTimestamp string `default:""          json:"dstAttrValTimestamp" yaml:"dstAttrValTimestamp"`
	DstAttrValHash      string `default:""          json:"dstAttrValHash" yaml:"dstAttrValHash"`
	DstAttrValId        string `default:""          json:"dstAttrValId" yaml:"dstAttrValId"`

	AlgorithmToUseForHash     string `default:"md5"             json:"algorithmToUseForHash" yaml:"algorithmToUseForHash"`
	DstServerCertPemFile      string `default:"dst-srv.crt.pem" json:"dstServerCertPemFile" yaml:"dstServerCertPemFile"`
	DstAttrValTimestampFormat string `default:"2000-01-01T01:00:00Z" json:"dstAttrValTimestampFormat" yaml:"dstAttrValTimestampFormat"`
	ThisServerCtx             string `default:"/"   json:"thisServerCtx" yaml:"thisServerCtx"`
	ProxyAttrRemoteUserName   string `default:"REMOTE_USERNAME" json:"proxyAttrRemoteUsername" yaml:"proxyAttrRemoteUsername"`
}

type Configuration struct {
	CliOpts CmdLineOptions `json:"cliopts" yaml:"cliopts"`
	Log     Log            `json:"log" yaml:"log"`
	WebCtxs []WebCtx       `json:"webctxs" yaml:"webctxs"`
}

var myCfg = Configuration{}

func scanConfiguration() {
	// Parse config for Default values and to initialize configuration structure.
	log.Print("Parse config for Default values ....")
	err := config.ParseDefault(&myCfg)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}

	// Parse config from Command line variables the first time.
	log.Print("Parse config from Command line variables ....")
	err = config.ParseCli(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}
	var optCfgFile string = myCfg.CliOpts.OptCfgFile
	//var optDebug int = myCfg.OptDebug
	//var optHelp bool = myCfg.OptHelp
	//if optHelp == true {
	//      flag.Usage()
	//      os.Exit(0)
	//}

	// Parse config for default configuration file
	//   -> first found of config.json and config.yaml
	log.Print("Parse config for default configuration file (config.json or config.yaml) ....")
	err = config.ParseConfigFile(&myCfg, "")
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}

	// Parse config from configuration file specified by command line
	if optCfgFile != "" {
		log.Print("Parse config from configuration file specified by command line option --cfgfile ....")
		err = config.ParseConfigFile(&myCfg, optCfgFile)
		if err != nil {
			// hmmm ... is not a real error
			log.Print(err)
		}
	}

	// Parse config from Enfironment variables
	log.Print("Parse config from Enfironment variables ....")
	err = config.ParseEnv(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}
	// Re-Parse config from Command line variables
	log.Print("Parse config from Command line variables ....")
	err = config.ParseCli(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}

	fmt.Printf("%# v\n", pretty.Formatter(myCfg))
}
