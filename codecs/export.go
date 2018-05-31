/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package codecs

import (
	wc "git.juniper.net/sksubra/wedge/codecs/json_grpc_codec"
	wu "git.juniper.net/sksubra/wedge/util"
)

/*
 * Function to create the codec object based
 * on codec type
 */
func GetCodecObject(codecType int) wu.GenCodec {
	switch codecType {
	case wu.CODEC_JSON_GRPC:
		return wc.CreateJsonGrpcCodecObject()
	default:
		return nil
	}
}
