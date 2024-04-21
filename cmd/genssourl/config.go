//
// got idea from:
//   https://github.com/eschao/config

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eschao/config"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"reflect"
)

// cmdline options
type CmdLineOptions struct {
	OptCfgSvcWebTcp         string `default:""     cli:"wtcp The HTTP server name/address[:port] to use."`
	OptCfgSvcWebTcpCertFile string `default:""     cli:"wtcpcert For HTTPS the cert.pem file tu use"`
	OptCfgSvcWebTcpKeyFile  string `default:""     cli:"wtcpkey For HTTPS the key.pem file tu use."`
	OptCfgSvcFcgiTcp        string `default:""     cli:"ftcp The FCGI server TCP name/address:port to use."`
	OptCfgSvcFcgiUnix       string `default:""     cli:"funix The FCGI server Unix socket name to use."`
	OptCfgSvcFcgiStdIO      bool   `default:"true" cli:"fstdio Use the FCGI server listen on standard I/O."`

	OptCfgFile        string `default:""       cli:"cfgfile The config file to use."`
	OptCfgAsJSON      bool   `default:"false"  cli:"json Print config as JSON."`
	OptCfgAsYAML      bool   `default:"false"  cli:"yaml Print config as YAML."`
	OptCopyEmbeddedFS bool   `default:"false"  cli:"copyefs Copy content of embedded filesystem."`
	OptDebug          int    `default:"0"      cli:"debug The debug level to use for command line, 0 means no debug."`
	OptSyslog         bool   `default:"true"   cli:"syslog Make log messages syslog compatible."`
	OptWebAppDirRoot  string `default:""       cli:"webappdir The directory in which the webapp files are located."`
}

// program log stuff
//
//	currently not implemented
type Log struct {
	LogPath  string `default:"logs"          json:"logPath" yaml:"logPath"`
	LogLevel string `default:"debug"         json:"logLevel" yaml:"logLevel"`
}

// web and service context
type WebCtx struct {
	// // set attributes
	// server_protocol := "https"
	// server_host := "localhost"
	// server_port := "5677"
	// server_path := "bla"
	// url_attr_username_key := "user"
	// url_attr_timestamp_key := "ts"
	// url_attr_hash_key := "hash"
	// url_attr_id_key := "id"
	//
	// username_val := "user1"
	// timestamp_val := "2006-01-02T15:04:05Z"
	// pub_pem_file := "testdata/test.crt.pem"
	// id_val := "test"

	DstServerProtocol string `default:"https"     json:"dstServerProtocol" yaml:"dstServerProtocol"`
	DstServerHost     string `default:"localhost" json:"dstServerHost" yaml:"dstServerHost"`
	DstServerPort     string `default:""          json:"dstServerPort" yaml:"dstServerPort"`
	DstServerPath     string `default:""          json:"dstServerPath" yaml:"dstServerPath"`

	DstAttrKeyUsername  string `default:"user"      json:"dstAttrKeyUsername" yaml:"dstAttrKeyUsername"`
	DstAttrKeyTimestamp string `default:"ts"        json:"dstAttrKeyTimestamp" yaml:"dstAttrKeyTimestamp"`
	DstAttrKeyHash      string `default:"hash"      json:"dstAttrKeyHash" yaml:"dstAttrKeyHash"`
	DstAttrKeyId        string `default:"id"        json:"dstAttrKeyId" yaml:"dstAttrKeyId"`

	DstAttrValUsername  string `default:""          json:"dstAttrValUsername" yaml:"dstAttrValUsername"`
	DstAttrValTimestamp string `default:""          json:"dstAttrValTimestamp" yaml:"dstAttrValTimestamp"`
	DstAttrValHash      string `default:""          json:"dstAttrValHash" yaml:"dstAttrValHash"`
	DstAttrValId        string `default:""          json:"dstAttrValId" yaml:"dstAttrValId"`

	AlgorithmToUseForHash          string   `default:"md5"             json:"algorithmToUseForHash" yaml:"algorithmToUseForHash"`
	DstServerPemFile               string   `default:"dst-srv.crt.pem" json:"dstServerPemFile" yaml:"dstServerPemFile"`
	DstServerUseSigning            bool     `default:"false"           json:"dstServerUseSigning" yaml:"dstServerUseSigning"`
	DstAttrValTimestampFormat      string   `default:"2006-01-02T15:04:05Z" json:"dstAttrValTimestampFormat" yaml:"dstAttrValTimestampFormat"`
	DstAttrValTimezone             string   `default:"UTC"             json:"dstAttrValTimezone" yaml:"dstAttrValTimezone"`
	DstDoNotDoParameterURLEncoding bool     `default:"false"          json:"dstDoNotDoParameterURLEncoding" yaml:"dstDoNotDoParameterURLEncoding"`
	ThisServerPath                 string   `default:"/"               json:"thisServerPath" yaml:"thisServerPath"`
	ProxyAttrRemoteUsername        string   `default:"REMOTE_USERNAME" json:"proxyAttrRemoteUsername" yaml:"proxyAttrRemoteUsername"`
	ProxyAttrRemoteUsernames       []string `json:"proxyAttrRemoteUsernames" yaml:"proxyAttrRemoteUsernames"`
}

type Configuration struct {
	CliOpts CmdLineOptions `json:"cliOpts" yaml:"cliOpts"`
	Log     Log            `json:"log" yaml:"log"`
	WebCtxs []WebCtx       `json:"webCtxs" yaml:"webCtxs"`
}

type Configuration4Output struct {
	Log     Log      `json:"log" yaml:"log"`
	WebCtxs []WebCtx `json:"webCtxs" yaml:"webCtxs"`
}

var myCfg = Configuration{}
var myDefaultWebCtx = WebCtx{}

func scanConfiguration() {
	// Parse config for Default web context and to initialize configuration structure.
	//log.Print("Parse config for Default web context values ....")
	err := config.ParseDefault(&myDefaultWebCtx)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}

	// Parse config for Default values and to initialize configuration structure.
	//log.Print("Parse config for Default values ....")
	err = config.ParseDefault(&myCfg)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}

	// Parse config from Command line variables the first time.
	//log.Print("Parse config from Command line variables ....")
	err = config.ParseCli(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}
	var optCfgFile string = myCfg.CliOpts.OptCfgFile

	// Remove Date and Time if we run under syslog
	if myCfg.CliOpts.OptSyslog == true {
		log.SetFlags(0)
	}

	// Parse config for default configuration file
	//   -> first found of config.json and config.yaml
	log.Print("Will try to parse config for default configuration file (config.json or config.yaml) ....")
	err = config.ParseConfigFile(&myCfg, "")
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}

	// Parse config from configuration file specified by command line
	if optCfgFile != "" {
		log.Print("Parse config from configuration file '" + optCfgFile + "' specified by command line option --cfgfile ....")
		err = config.ParseConfigFile(&myCfg, optCfgFile)
		if err != nil {
			// hmmm ... is not a real error
			log.Print(err)
		}
	}

	// Parse config from Environment variables
	log.Print("Parse config from Environment variables ....")
	err = config.ParseEnv(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}
	// Re-Parse config from Command line variables
	log.Print("Parse config from Command line variables again ....")
	err = config.ParseCli(&myCfg)
	if err != nil {
		// hmmm ... is not a real error
		log.Print(err)
	}

	// set default web context values if missed
	for idx := 0; idx < len(myCfg.WebCtxs); idx++ {
		chkVal := reflect.ValueOf(myCfg.WebCtxs[idx])
		pChkVal := reflect.ValueOf(&myCfg.WebCtxs[idx])
		defVal := reflect.ValueOf(myDefaultWebCtx)
		for i := 0; i < chkVal.NumField(); i++ {
			if chkVal.Field(i).Interface() == "" {
				if myCfg.CliOpts.OptDebug >= 4 {
					fmt.Printf("Before: %s %# v  <-> %# v\n",
						chkVal.Type().Field(i).Name,
						chkVal.Field(i).Interface(),
						defVal.Field(i).Interface())
				}
				f := pChkVal.Elem().FieldByName(chkVal.Type().Field(i).Name)
				f.Set(defVal.Field(i))
				if myCfg.CliOpts.OptDebug >= 4 {
					fmt.Printf("After:  %s %# v  <--- %# v\n",
						chkVal.Type().Field(i).Name,
						chkVal.Field(i).Interface(),
						defVal.Field(i).Interface())
				}
			}
		}
	}

	if myCfg.CliOpts.OptDebug >= 2 {
		fmt.Printf("%# v\n", pretty.Formatter(myCfg))
	}

	return
}

// prepare a configuration structur for output
func prepareCfg4Output() Configuration4Output {
	var myCfg4Output = Configuration4Output{}

	err := config.ParseDefault(&myCfg4Output)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}

	if len(myCfg.WebCtxs) > 0 {
		myCfg4Output.WebCtxs = myCfg.WebCtxs
	} else {
		myCfg4Output.WebCtxs = append(myCfg4Output.WebCtxs, myDefaultWebCtx)
	}

	return myCfg4Output
}

// got this function from
//
//	https://gosamples.dev/pretty-print-json/
func prettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func printCfgAsJSON() {
	var buffer bytes.Buffer
	var myCfg4Output = prepareCfg4Output()

	err := prettyEncode(&myCfg4Output, &buffer)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}
	fmt.Println(buffer.String())

	return
}

func printCfgAsYAML() {
	var myCfg4Output = prepareCfg4Output()

	data, err := yaml.Marshal(myCfg4Output)
	if err != nil {
		// hmmm ... a real error
		log.Fatal(err)
	}
	fmt.Println(string(data))

	return
}
