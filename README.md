### Overview:
Wedge is a broker written in Go to support execution of APIs over different
formats and transports. The implementation follows a plugin based model where
the client-side and server-side plugins can run with different data formats,
transport and communicate using a common message format representing data in
JSON. The code is organized into 3 levels:

1. Codecs: Specify the codec to translate from a custom data to protobuf
   byte-buffer for gRPC. Currently, JSON to protobuf is implemented. The
   GenCodec interface in utils/utils.go specifies the methods to be implemented
   for a codec.

2. Plugins: Fundamental part of the application and can include
   codecs if needed. Currently, the gRPC plugin uses a codec. Each plugin
   implementation needs to implement a ClientSidePlugin and/or ServerSidePlugin
   interface(s).
   The ClientSidePlugin interface specifies the functions/methods
   to be implemented when the plugin will be interfacing with the actual client
   i.e the plugin will be acting as a server.
   The ServerSidePlugin interface specifies the functions/methods
   to be implemented when the plugin will be interfacing with the actual server
   i.e the plugin will be acting as a client. The plugins can also run
   user-defined callbacks if customized handling is needed. The client and
   server side/facing plugins have a common data structure to pass data between
   them using the MsgFormat struct.

3. Implementations: Implementations group together a client-side and server-side
   plugin implementation. WedgeImpl struct implements the Implementation
   interface to specify the plugins to be run on client and server side,
   callbacks and so on. The client and server side plugins will be as Go routines.

### Dependencies:
1. Protobuf 3
2. protobuf-c 1.2.1

### Prerequisites:
1. [Go](https://golang.org/doc/install) 1.8+
2. Configure [GOPATH](https://golang.org/doc/code.html#GOPATH)
3. [librdkafka](https://github.com/edenhill/librdkafka)
4. [protoc-wedge](https://github.com/Juniper/protoc-wedge) compiler
   to generate descriptor files needed for wedge specific to your application. 

### Installation:
Wedge Makefile requires GNU make.

Dependencies are managed with gdm, which is installed by the Makefile if you
don't have it already including gRPC Go and protobuf.
1. Run `go get -d github.com/Juniper/wedge/...`
2. Run `cd $GOPATH/src/github.com/Juniper/wedge`
3. Run `make`
4. Run `make install` to copy the wedge binary to /usr/local/bin

### Usage:
#### The usage can be determined through the command:

```
./wedge
```

#### To list the available client and server facing plugins:

```
./wedge --list-plugins
```

#### To list the available user defined pre-write and post-read callbacks:

```
./wedge --list-callbacks
```
  
#### To generate sample config for a given client and server facing plugins in yaml:

```
./wedge --generate --client-side kafka:kafka  --server-side grpc:kafka
```

The Config will be specified in yaml and the plugin list will be separated by
a semi-colon

#### To run wedge with a conf file as input:

```
./wedge --run <filename>
```

### User defined callbacks:
The user can specify custom callback functions after the data is read from the
transport and before the data is written to the transport without changing the
underlying plugin implementaion.

A Post-read callback will be executed when data is read by the plugin and before 
it is sent to the other side through a channel. The input will be an interface{}
specific to the plugin and the output will be MsgFormat struct.
A Pre-write callback will be executed when the data is received and before it
needs to be sent to the actual client or server. 
If the callback functions are nil, then the default behavior specific to the
plugin implementation will be performed. The plugins are expect to send MsgFormat
structs for communication. 

### Getting started:
The _examples folder has all the avsc, proto descriptor and avro descriptor
table JSON files, example clients etc. to try out. Make a directory called
"WedgeClient" within $GOPATH/src and copy the file "_examples/clients/WedgeDesc.go"
there.
The examples cover the following:
1. Bgp route addition using JSON/Avro as data format on the client side to
   publish data to a kafka bus. "_examples/yaml_config" has sample yaml
   configuration files to start the broker and connect to the kafka bus.
   The broker's server-side plugin will use gRPC to communicate with a Juniper
   router and add routes using BGP JET APIs.

2. Fetch OpenConfig telemetry data from Juniper routers and use custom callback
   functions to convert the data and publish it to a kafka bus. The following
   format conversions are performed:<br/>
   a. Influx line protocol format and published to topic "telegraf" which can be
      processed and added to InfluxDB using [telegraf](https://github.com/influxdata/telegraf)<br/>
   b. Fluentd compatible key-value pair data that can be consumed using
      [fluent-plugin-kafka](https://github.com/fluent/fluent-plugin-kafka)<br/>
   c. JSON data representation of OpenConfig telemetry key-value pairs.    

#### Adding a new plugin:
1. Add a new file for the plugin within "plugins" directory.
2. Implement the methods in "ClientSidePlugin" and/or "ServerSidePlugin"
   interface(s).
3. Implement a function to emit reference yaml configuration for the plugin
   like kafkaEmitRefConfig().
4. Implement a plugin init function like KafkaInit() and register it with the
   plugin map using PluginInit().

#### Adding a new callback function:
1. Implement the callback function in implementations/callbacks.go.
2. A post-read callback function will take an interface{} as input that would
   be specific to the plugin/data format and return an object of struct MsgFormat. 
3. A pre-write callback function will take an object of struct MsgFormat as
   as input and returns an interface{} specific to the plugin/data format like
   JSONPostReadCB().
