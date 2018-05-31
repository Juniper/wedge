/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package all

import (
	gRPC "git.juniper.net/sksubra/wedge/plugins/gRPC"
	kafka "git.juniper.net/sksubra/wedge/plugins/kafka"
	wu "git.juniper.net/sksubra/wedge/util"
)

/*
 * Main function to call all the init functions
 * defined by plugin
 */
func PluginInit() {
	gRPC.GrpcInit(wu.PluginMap)
	kafka.KafkaInit(wu.PluginMap)
}
