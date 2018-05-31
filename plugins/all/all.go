/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package all

import (
	gRPC "github.com/Juniper/wedge/plugins/gRPC"
	kafka "github.com/Juniper/wedge/plugins/kafka"
	wu "github.com/Juniper/wedge/util"
)

/*
 * Main function to call all the init functions
 * defined by plugin
 */
func PluginInit() {
	gRPC.GrpcInit(wu.PluginMap)
	kafka.KafkaInit(wu.PluginMap)
}
