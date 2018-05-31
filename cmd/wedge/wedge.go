/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	wi "github.com/Juniper/wedge/implementations"
	plugin "github.com/Juniper/wedge/plugins/all"
	wu "github.com/Juniper/wedge/util"
	yaml "gopkg.in/yaml.v2"
)

/*
 * Define all the flag variables
 */
var wGen = flag.Bool("generate", false,
	"Generate the sample config for a given client-side and server-side plugin")
var wRun = flag.Bool("run", false,
	"Run the broker with the specified conf file")
var wList = flag.Bool("list-plugins", false, "Show the available client side "+
	"and server side plugins")
var WCallBack = flag.Bool("list-callbacks", false, "Show the available user "+
	"defined callbacks")
var wClientside = flag.String("client-side", "",
	"Client side plugin to run i.e the plugin will act as a server to the "+
		"actual client")
var wServerside = flag.String("server-side", "",
	"Server side plugin to run i.e the plugin will act as a client to the "+
		"actual server")

const usage = `
Usage:

  Wedge [flags|options]

The flags are:
  --run <file>        Run the broker based on an input file
  --list-plugins      List all the available client facing and server facing plugins
  --list-callbacks    List all the available user defined pre-write and post-read callbacks
  --client-side       Client facing plugin to run/generate sample config, separated by :
  --server-side       Server facing plugin to run/generate sample config, separated by :
  --generate          Generate sample config for a given client and server facing plugins


Examples:

  #Generate config for a given client and server facing plugins. The Config will be specified in yaml
  wedge --generate --client-side kafka:kafka  --server-side grpc:kafka

  #Run wedge with a conf file as input
  wedge --run <filename>

  #List the available client and server facing plugins
  wedge --list-plugins

  #List the available user defined pre-write and post-read callbacks
  wedge --list-callbacks
`

func usageExit(rc int) {
	fmt.Println(usage)
	os.Exit(rc)
}

func listPlugins() {
	fmt.Println("Available Client facing plugins:")
	for key, value := range wu.PluginMap {
		if value.ClientSidePluginSupport {
			fmt.Println(key)
		}
	}

	fmt.Println("Available Server facing plugins:")
	for key, value := range wu.PluginMap {
		if value.ServerSidePluginSupport {
			fmt.Println(key)
		}
	}
}

/*
 * String constants used to emit reference config
 */
const INIT_CONFIG = `
---
enable-logging: false # Logging is disabled by default
log-location: "" # Logfile location, default: /var/tmp/wedge.log

implementations:`

const EMIT_CLIENT_SIDE = `
    -
        client-side:`

const EMIT_SERVER_SIDE = `
        server-side:`

const END_CONFIG = `
...
`

const DEFAULT_LOG_LOCATION = "/var/tmp/wedge.log"

/*
 * Main function to run Wedge broker
 */
func main() {

	/*
	 * Call the PluginInit function to build the
	 * PluginMap variable and
	 */
	plugin.PluginInit()
	wi.PostReadMapInit()
	wi.PreWriteMapInit()
	wi.ErrorFuncMapInit()

	flag.Usage = func() { usageExit(0) }
	flag.Parse()
	args := flag.Args()

	switch {
	case *wList:
		listPlugins()
	case *WCallBack:
		/*
		 * Walk the post-read and pre-write callback function
		 * map and emit the value
		 */
		fmt.Println("Post Read callback functions:")
		for funcName, _ := range wi.PostReadFuncMap {
			fmt.Println(funcName)
		}

		fmt.Println("\nPre Write callback functions:")
		for funcName, _ := range wi.PreWriteFuncMap {
			fmt.Println(funcName)
		}
	case *wGen:
		/*
		 * Generate the sample config for client-side and
		 * server-side plugins
		 */
		if strings.Compare(*wClientside, "") == 0 ||
			strings.Compare(*wServerside, "") == 0 {
			usageExit(0)
		}

		clientSidePlugins, serverSidePlugins := []string{}, []string{}
		clientSidePlugins = strings.Split(*wClientside, ":")

		serverSidePlugins = strings.Split(*wServerside, ":")

		if len(clientSidePlugins) != len(serverSidePlugins) {
			fmt.Println("Mismatch in lengths between client and server side " +
				"plugins list")
			usageExit(0)
		}

		var output string

		output = INIT_CONFIG

		for i, _ := range clientSidePlugins {
			output += EMIT_CLIENT_SIDE
			cPlugin, ok := wu.PluginMap[clientSidePlugins[i]]
			if !ok {
				fmt.Println("Unknown plugin", clientSidePlugins[i], "!!!")
				listPlugins()
				os.Exit(0)

			}
			output += cPlugin.EmitClientSideRefConfig()

			output += EMIT_SERVER_SIDE
			sPlugin, ok := wu.PluginMap[serverSidePlugins[i]]
			if !ok {
				fmt.Println("Unknown plugin", serverSidePlugins[i], "!!!")
				listPlugins()
				os.Exit(0)

			}

			output += sPlugin.EmitServerSideRefConfig()
		}

		output += END_CONFIG

		fmt.Println(output)
	case *wRun:
		var err error
		var ok bool
		var genericVal map[interface{}]interface{}
		var clientPreWrite wu.PreWriteFunc
		var clientPostRead wu.PostReadFunc
		var serverPreWrite wu.PreWriteFunc
		var serverPostRead wu.PostReadFunc
		var waitGroup sync.WaitGroup
		var cErrorFunc, sErrorFunc wu.ErrorFunc

		if len(args) == 0 {
			fmt.Println("Filename is not provided")
			usageExit(0)
		}

		genericMap := make(map[interface{}]interface{})

		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatalln("Reading file", args[0], "failed with error", err)

		}

		err = yaml.Unmarshal([]byte(data), &genericMap)
		if err != nil {
			log.Fatalf("Unmarshalling yaml data failed with error: %v", err)
		}

		/*
		 * Set the logging options if specified
		 */
		 logOption, ok := genericMap["enable-logging"]
		 if !ok {
		 	log.Fatalf("Key enable-logging not found")
		 }

		logging := logOption.(bool)
		var fileHandle *os.File

		if logging == true {
			logOption, ok := genericMap["log-location"]
			var file string

			if !ok || logOption.(string) == "" {
				file = DEFAULT_LOG_LOCATION
			} else {
				file = logOption.(string)
			}

			fileHandle, err = os.OpenFile(file, os.O_RDWR | os.O_CREATE, 0666)
			if err != nil {
				log.Fatalf("Error opening file %s: %v", file, err)
			}
			log.SetOutput(fileHandle)
		} else {
			log.SetOutput(ioutil.Discard)
		}

		if logging == true && fileHandle != nil {
			defer fileHandle.Close()
		}

		implementations, ok := genericMap["implementations"]

		if !ok {
			log.Fatalf("Keyword implementations not found in yaml imput")
		}

		/*
		 * Walk the implementations list and start the go routines
		 * for each pair of client server config
		 */
		impl_list := implementations.([]interface{})

		for _, value := range impl_list {
			var impl wi.WedgeImpl
			genericVal = value.(map[interface{}]interface{})

			cl, ok := genericVal["client-side"]
			if !ok {
				log.Fatalf("Key client-side not found")
			}

			sl, ok := genericVal["server-side"]
			if !ok {
				log.Fatalf("Key server-side not found")
			}

			/*
			 * Parse the client side related config
			 */
			client_side := cl.(map[interface{}]interface{})

			name, ok := client_side["name"]
			if !ok {
				log.Fatalf("Key name not found")
			}

			cPlugin, ok := wu.PluginMap[name.(string)]
			if !ok {
				log.Fatalln("Unknown plugin", name)
			}

			/*
			 * Fetch the post-read and pre-write callback
			 * functions if specified
			 */
			preWrite, ok := client_side["client-pre-write-cb"]
			if !ok {
				log.Fatalf("Key client-pre-write-cb not found while parsing "+
					"plugin %s", name.(string))
			}

			postRead, ok := client_side["client-post-read-cb"]
			if !ok {
				log.Fatalf("Key client-post-read-cb not found while parsing "+
					"plugin %s", name.(string))
			}

			if preWrite.(string) == "" {
				clientPreWrite = nil
			} else {
				clientPreWrite, ok = wi.PreWriteFuncMap[preWrite.(string)]
				if !ok {
					log.Fatalf("Client pre write Callback function %s not found"+
						" when parsing plugin %s",
						preWrite.(string), name.(string))
				}
			}

			if postRead.(string) == "" {
				clientPostRead = nil
			} else {
				clientPostRead = wi.PostReadFuncMap[postRead.(string)]
				if !ok {
					log.Fatalf("Client post read Callback function %s not found"+
						" when parsing plugin %s",
						postRead.(string), name.(string))
				}
			}

			client_config, ok := client_side["config"]
			if !ok {
				log.Fatalf("Key config not found for plugin %s", name)
			}

			marshal_config, err := yaml.Marshal(&client_config)
			if err != nil {
				log.Fatalf("Incorrect client-side yaml configuration specified "+
					"for plugin %s", name)
			}

			/*
			 * Check the PluginMap to fetch the relevant config parse
			 * for the plugin
			 */
			if cPlugin.ParseClientSideConfig == nil {
				log.Fatalln("Client side config parse function not speficied "+
					"for plugin", name)
			}

			client_conf := cPlugin.ParseClientSideConfig(marshal_config)
			impl.ClientCodec = cPlugin.ClientSideCodec

			/*
			 * Parse the server side related config
			 */
			server_side := sl.(map[interface{}]interface{})

			name, ok = server_side["name"]
			if !ok {
				log.Fatalf("Key name not found")
			}

			sPlugin, ok := wu.PluginMap[name.(string)]
			if !ok {
				log.Fatalln("Unknown plugin", name)
			}

			/*
			 * Fetch the post-read and pre-write callback
			 * functions if specified
			 */
			preWrite, ok = server_side["server-pre-write-cb"]
			if !ok {
				log.Fatalf("Key server-pre-write-cb not found while parsing "+
					"plugin %s", name.(string))
			}

			postRead, ok = server_side["server-post-read-cb"]
			if !ok {
				log.Fatalf("Key server-post-read-cb not found while parsing "+
					"plugin %s", name.(string))
			}

			if preWrite.(string) == "" {
				serverPreWrite = nil
			} else {
				serverPreWrite, ok = wi.PreWriteFuncMap[preWrite.(string)]
				if !ok {
					log.Fatalf("Server pre write Callback function %s not found"+
						" when parsing plugin %s",
						preWrite.(string), name.(string))
				}
			}

			if postRead.(string) == "" {
				serverPostRead = nil
			} else {
				serverPostRead = wi.PostReadFuncMap[postRead.(string)]
				if !ok {
					log.Fatalf("Server post read Callback function %s not found"+
						" when parsing plugin %s",
						postRead.(string), name.(string))
				}
			}

			server_config, ok := server_side["config"]
			if !ok {
				log.Fatalf("Key config not found for plugin %s", name)
			}

			marshal_config, err = yaml.Marshal(&server_config)
			if err != nil {
				log.Fatalf("Incorrect server-side yaml configuration specified "+
					"for plugin %s", name)
			}

			/*
			 * Check the PluginMap to fetch the relevant config parse
			 * for the plugin
			 */
			if sPlugin.ParseServerSideConfig == nil {
				log.Fatalln("Server side config parse function not speficied "+
					"for plugin", name)
			}

			server_conf := sPlugin.ParseServerSideConfig(marshal_config)
			impl.ServerCodec = sPlugin.ServerSideCodec

			cErrorFuncName := cPlugin.FetchErrorFunc(client_conf)
			if cErrorFuncName == "" {
				cErrorFunc = nil
			} else {
				if cErrorFunc, ok = wi.ErrorFuncMap[cErrorFuncName]; !ok {
					cErrorFunc = nil
				}
			}

			sErrorFuncName := sPlugin.FetchErrorFunc(server_conf)
			if sErrorFuncName == "" {
				sErrorFunc = nil
			} else {
				if sErrorFunc, ok = wi.ErrorFuncMap[sErrorFuncName]; !ok {
					sErrorFunc = nil
				}
			}

			waitGroup.Add(1)
			go impl.Run(waitGroup, client_conf, server_conf, cPlugin.ClientPlugin,
				clientPostRead, clientPreWrite, cErrorFunc, sPlugin.ServerPlugin,
				serverPostRead, serverPreWrite, sErrorFunc)

		}

		waitGroup.Wait()

	default:
		usageExit(0)
	}

}
