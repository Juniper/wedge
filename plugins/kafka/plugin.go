/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package Kafka

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	wu "github.com/Juniper/wedge/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/srikanth2212/jsonez"
	yaml "gopkg.in/yaml.v2"
)

/*
 * The kafka parameter are defined with respect to the
 * parameter definition in
 * https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
 */
const (
	KAFKA_BROKER_ADDR_TYPE_ANY  = 0
	KAFKA_BROKER_ADDR_TYPE_IPV4 = 1
	KAFKA_BROKER_ADDR_TYPE_IPV6 = 2
)

const (
	KAFKA_SECURITY_PROTOCOL_PLAINTEXT      = "PLAINTEXT"
	KAFKA_SECURITY_PROTOCOL_SSL            = "SSL"
	KAFKA_SECURITY_PROTOCOL_SASL_PLAINTEXT = "SASL_PLAINTEXT"
	KAFKA_SECURITY_PROTOCOL_SASL_SSL       = "SASL_SSL"
)

type KafkaCommonConfig struct {
	AddressFamily               int      `yaml:"address-family"` // Defaults to KAFKA_BROKER_ADDR_TYPE_ANY
	BackoffJitter               uint32   `yaml:"backoff-jitter"`
	BlacklistTopics             []string `yaml:"blacklist-topics"`
	BootStrapServers            string   `yaml:"bootstrap-servers"`
	BrokerAddressTTL            uint32   `yaml:"broker-address-ttl"`
	ClientId                    string   `yaml:"broker-client-id"`
	DisableNagle                bool     `yaml:"disable-nagle"`
	DisableSparseMetadataReq    bool     `yaml:"disable-sparse-metadata-req"`
	EnableTcpKeepalives         bool     `yaml:"enable-tcp-keepalives"` // Default is false
	GroupId                     string   `yaml:"group-id"`
	MaxMsgBytes                 uint32   `yaml:"max-msg-bytes"`
	MaxCopyBytes                uint32   `yaml:"max-copy-bytes"`
	MaxInflightRequests         uint32   `yaml:"max-inflight-requests"`
	MetadatareqTimeout          uint32   `yaml:"metadata-req-timeout"` // Value in ms
	MetadataRefreshInterval     int32    `yaml:"metadata-refresh-interval"`
	MetadataMaxAge              uint32   `yaml:"metadata-max-age"`
	MetadataRefreshFastInterval uint32   `yaml:"metadata-refresh-fast-interval"`
	SocketTimeout               uint32   `yaml:"socket-timeout"`
	SocketSendBufSize           uint32   `yaml:"socket-send-bufsize"`
	SocketRecvBufSize           uint32   `yaml:"socket-recv-bufsize"`
	MaxSendFails                uint32   `yaml:"max-send-fails"`
	RecvMaxBytes                uint32   `yaml:"recv-max-bytes"`
	SecurityProtocol            string   `yaml:"security-protocol"` //Defaults to KAFKA_SECURITY_PROTOCOL_PLAINTEXT
	SSLCipherSuites             string   `yaml:"ssl-cipher-suites"`
	SSLKeysLocation             string   `yaml:"ssl-keys-location"`
	SSLKeyPass                  string   `yaml:"ssl-key-pass"`
	SSLCertLocation             string   `yaml:"ssl-cert-location"`
	SSLCaLocation               string   `yaml:"ssl-ca-location"`
	SSLCrlLocation              string   `yaml:"ssl-crl-location"`
	SaslMechanisms              string   `yaml:"sasl-mechanisms"`
	SaslKeberosService          string   `yaml:"sasl-keberos-service"`
	SaslKeberosPrincipal        string   `yaml:"sasl-keberos-principal"`
	SaslKeberosKinitCmd         string   `yaml:"sasl-keberos-kinit-cmd"`
	SaslKeberosKeyTab           string   `yaml:"sasl-keberos-keytab"`
	SaslKeberosLoginInterval    uint32   `yaml:"sasl-keberos-login-interval"`
	SaslUserName                string   `yaml:"sasl-username"`
	SaslPassword                string   `yaml:"sasl-password"`
}

const (
	COMPRESSION_CODEC_TYPE_NONE   = "none"
	COMPRESSION_CODEC_TYPE_GZIP   = "gzip"
	COMPRESSION_CODEC_TYPE_SNAPPY = "snappy"
	COMPRESSION_CODEC_TYPE_LZ4    = "lz4"
)

/*
 * kafka producer config
 */
type KafkaProducerConfig struct {
	MaxQueueBufferMessages uint32 `yaml:"max-queue-buffer-messages"`
	MaxQueueBufferKbytes   uint32 `yaml:"max-queue-buffer-kbytes"`
	MaxQueueBuffertime     uint32 `yaml:"max-queue-buffer-time"` // Value in ms
	MaxSendRetries         uint32 `yaml:"max-send-retries"`
	RetryBackoff           uint32 `yaml:"retry-backoff"`     // Value in ms
	CompressionCodec       string `yaml:"compression-codec"` // Default is COMPRESSION_CODEC_TYPE_NONE
	MaxNumBatchMessages    uint32 `yaml:"max-num-batch-messages"`
	EnableErrDeliveryonly  bool   `yaml:"enable-err-delivery-only"` // If delivery report is needed only for error

	/*
	 * kafka-go specific options
	 */
	EnableBatchProducer         bool   `yaml:"enable-batch-producer"`
	DisablePerMsgDeliveryReport bool   `yaml:"disable-per-msg-delivery-report"`
	ProducerChannelSize         uint32 `yaml:"producer-channel-size"`

	//Topic level configurations
	Acks                int    `yaml:"acks"`
	MessageTimeout      uint32 `yaml:"message-timeout"`
	ProduceOffsetReport bool   `yaml:"produce-offset-report"`
	RequestTimeout      uint32 `yaml:"request-timeout"` // Value in ms
}

const (
	AUTO_OFFSET_RESET_TYPE_SMALLEST  = "smallest"
	AUTO_OFFSET_RESET_TYPE_EARLIEST  = "earliest"
	AUTO_OFFSET_RESET_TYPE_BEGINNING = "beginning"
	AUTO_OFFSET_RESET_TYPE_LARGEST   = "largest"
	AUTO_OFFSET_RESET_TYPE_LATEST    = "latest"
	AUTO_OFFSET_RESET_TYPE_END       = "end"
	AUTO_OFFSET_RESET_TYPE_ERROR     = "error"
)

const DEFAULT_KAFKA_CLIENT_ID = "wedge_kafka_client"

const KAFKA_ERROR_TOPIC = "_wedge.Topic_error"

/*
 * kafka consumer config
 */
type KafkaConsumerConfig struct {
	DisableAutoCommit   bool   `yaml:"disable-auto-commit"`
	AutoCommitInterval  uint32 `yaml:"auto-commit-interval"`
	DisableAutoOffset   bool   `yaml:"disable-auto-offset"`
	MinPartitionMsgs    uint32 `yaml:"min-partition-msgs"`
	MaxKbPerPartition   uint32 `yaml:"max-kb-per-partition"` // Value in kilobytes
	MaxFetchWait        uint32 `yaml:"max-fetch-wait"`       // Value in ms
	MaxMsgFetchBytes    uint32 `yaml:"max-msg-fetch-bytes"`
	MinFetchBytes       uint32 `yaml:"min-fetch-bytes"`
	FetchErrorBackoff   uint32 `yaml:"fetch-error-backoff"` // Value in ms
	DisablePartitionEof bool   `yaml:"disable-partition-eof"`
	CheckCRCs           bool   `yaml:"check-crcs"`

	//Topic level configurations
	AutoOffsetResetType     string `yaml:"auto-offset-reset-type"` // AUTO_OFFSET_RESET_TYPE_LARGEST
	OffsetStorePath         string `yaml:"offset-store-path"`
	OffsetStoreSyncInterval uint32 `yaml:"offset-store-sync-interval"`
}

/*
 * Kafka producer struct
 */
type KafkaProducer struct {
	kprod    *kafka.Producer
	clientId string
}

/*
 * Kafka producer struct
 */
type KafkaConsumer struct {
	kcon     *kafka.Consumer
	topics   []string
	clientId string
}

const (
	KAFKA_TYPE_PRODUCER = 1
	KAFKA_TYPE_CONSUMER = 2
)

type PublishMsg struct {
	Topic   string
	Key     []byte
	Payload []byte
}

const (
	CONSUMER_TERMINATED        string = "Consumer terminated"
	CONSUMER_TOPIC_INIT_FAILED string = "SubscribeTopics() failed"
	PRODUCER_TERMINATED        string = "Producer terminated"
)

/*
 * Config parse functions
 */
func parseCommonConfig(c *KafkaCommonConfig, m *kafka.ConfigMap,
	ktype int) error {
	var err error
	var errorStr string

	if ktype == KAFKA_TYPE_PRODUCER {
		log.Printf("%s: type is producer", wu.FuncName())
	} else {
		log.Printf("%s: type is consumer", wu.FuncName())
	}

	if c.AddressFamily > 0 {
		err = m.SetKey("broker.address.family", c.AddressFamily)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding broker.address.family with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.AddressFamily, err)
			return errors.New(errorStr)
		}
	}

	if c.BackoffJitter > 0 {
		err = m.SetKey("reconnect.backoff.jitter.ms", c.BackoffJitter)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding reconnect.backoff.jitter.ms "+
				"with value %d failed with error: %v", wu.FuncName(),
				c.BackoffJitter, err)
			return errors.New(errorStr)
		}
	}

	if len(c.BlacklistTopics) > 0 {
		err = m.SetKey("topic.blacklist", c.BlacklistTopics)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding topic.blacklist failed with "+
				"error: %v", wu.FuncName(), err)
			return errors.New(errorStr)
		}
	}

	if c.BootStrapServers == "" {
		errorStr = fmt.Sprintf("%s: bootstrap.servers cannot be empty",
			wu.FuncName())
		return errors.New(errorStr)
	} else {
		err = m.SetKey("bootstrap.servers", c.BootStrapServers)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding bootstrap.servers with value "+
				"%s failed with error: %v", wu.FuncName(), c.BootStrapServers,
				err)
			return errors.New(errorStr)
		}
	}

	if c.BrokerAddressTTL > 0 {
		err = m.SetKey("broker.address.ttl", c.BrokerAddressTTL)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding broker.address.ttl with value "+
				"%d failed with error: %v", wu.FuncName(), c.BrokerAddressTTL,
				err)
			return errors.New(errorStr)
		}
	}

	if c.ClientId == "" {
		c.ClientId = "wedge"
	}
	err = m.SetKey("client.id", c.ClientId)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding client.id with value %s failed "+
			"with error: %v", wu.FuncName(), c.ClientId, err)
		return errors.New(errorStr)
	}

	err = m.SetKey("socket.nagle.disable", c.DisableNagle)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding broker.address.ttl with value %d "+
			"failed with error: %v", wu.FuncName(), c.BrokerAddressTTL, err)
		return errors.New(errorStr)
	}

	if c.DisableSparseMetadataReq == false {
		err = m.SetKey("topic.metadata.refresh.sparse", true)
	} else {
		err = m.SetKey("topic.metadata.refresh.sparse", false)
	}
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding topic.metadata.refresh.sparse with "+
			"value %d failed with error: %v", wu.FuncName(),
			!c.DisableSparseMetadataReq, err)
		return errors.New(errorStr)
	}

	err = m.SetKey("socket.keepalive.enable", c.EnableTcpKeepalives)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding socket.keepalive.enable with "+
			"value %d failed with error: %v", wu.FuncName(),
			c.EnableTcpKeepalives, err)
		return errors.New(errorStr)
	}

	if c.MaxMsgBytes > 0 {
		err = m.SetKey("message.max.bytes", c.MaxMsgBytes)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding message.max.bytes with value %d"+
				" failed with error: %v", wu.FuncName(), c.MaxMsgBytes, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxCopyBytes > 0 {
		err = m.SetKey("message.copy.max.bytes", c.MaxCopyBytes)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding message.copy.max.bytes with "+
				"value %d failed with error: %v", wu.FuncName(), c.MaxCopyBytes,
				err)
			return errors.New(errorStr)
		}
	}

	if c.RecvMaxBytes > 0 {
		err = m.SetKey("receive.message.max.bytes", c.RecvMaxBytes)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding receive.message.max.bytes with value %d "+
					"failed with error: %v", wu.FuncName(), c.RecvMaxBytes,
					err)
			return errors.New(errorStr)
		}
	}

	if c.MaxInflightRequests > 0 {
		err = m.SetKey("max.in.flight.requests.per.connection",
			c.MaxInflightRequests)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding max.in.flight.requests.per.connection with"+
					" value %d failed with error: %v", wu.FuncName(),
					c.MaxInflightRequests, err)
			return errors.New(errorStr)
		}
	}

	if c.MetadatareqTimeout > 0 {
		err = m.SetKey("metadata.request.timeout.ms", c.MetadatareqTimeout)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding metadata.request.timeout.ms with value %d "+
					"failed with error: %v", wu.FuncName(), c.MetadatareqTimeout,
					err)
			return errors.New(errorStr)
		}
	}

	if c.MetadataRefreshInterval > 0 || c.MetadataRefreshInterval == -1 {
		err = m.SetKey("topic.metadata.refresh.interval.ms",
			c.MetadataRefreshInterval)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding topic.metadata.refresh.interval.ms with "+
					"value %d failed with error: %v", wu.FuncName(),
					c.MetadataRefreshInterval, err)
			return errors.New(errorStr)
		}
	}

	if ktype == KAFKA_TYPE_CONSUMER {
		log.Printf("%s: Checking for consumer group id", wu.FuncName())
		if strings.Compare(c.GroupId, "") == 0 {
			errorStr = fmt.Sprintf("%s: group.id cannot be empty for a consumer",
				wu.FuncName())
			return errors.New(errorStr)
		} else {
			err = m.SetKey("group.id", c.GroupId)
			if err != nil {
				errorStr = fmt.Sprintf("%s: Adding group.id with value %s failed "+
					"with error: %v", wu.FuncName(), c.GroupId, err)
				return errors.New(errorStr)
			}

			log.Printf("%s: group.id set to %s", wu.FuncName(), c.GroupId)
		}
	}

	if c.MetadataMaxAge > 0 {
		err = m.SetKey("metadata.max.age.ms", c.MetadataMaxAge)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding metadata.max.age.ms with value %d failed "+
					"with error: %v", wu.FuncName(), c.MetadataMaxAge, err)
			return errors.New(errorStr)
		}
	}

	if c.MetadataRefreshFastInterval > 0 {
		err = m.SetKey("topic.metadata.refresh.fast.interval.ms",
			c.MetadataRefreshFastInterval)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding topic.metadata.refresh.fast.interval.ms "+
					"with value %d failed with error: %v", wu.FuncName(),
					c.MetadataRefreshFastInterval, err)
			return errors.New(errorStr)
		}
	}

	if c.SocketTimeout > 0 {
		err = m.SetKey("socket.timeout.ms", c.SocketTimeout)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding socket.timeout.ms with value %d failed with"+
					" error: %v", wu.FuncName(), c.SocketTimeout, err)
			return errors.New(errorStr)
		}
	}

	if c.SocketSendBufSize > 0 {
		err = m.SetKey("socket.send.buffer.bytes", c.SocketSendBufSize)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding socket.send.buffer.bytes with value %d "+
					"failed with error: %v", wu.FuncName(), c.SocketSendBufSize,
					err)
			return errors.New(errorStr)
		}
	}

	if c.SocketRecvBufSize > 0 {
		err = m.SetKey("socket.receive.buffer.bytes", c.SocketRecvBufSize)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding socket.receive.buffer.bytes with value %d "+
					"failed with error: %v", wu.FuncName(), c.SocketRecvBufSize,
					err)
			return errors.New(errorStr)
		}
	}

	if c.MaxSendFails > 0 {
		err = m.SetKey("socket.max.fails", c.MaxSendFails)
		if err != nil {
			errorStr =
				fmt.Sprintf("%s: Adding socket.max.fails with value %d "+
					"failed with error: %v", wu.FuncName(), c.MaxSendFails,
					err)
			return errors.New(errorStr)
		}
	}

	if strings.Compare(c.SecurityProtocol, "") != 0 {
		err = m.SetKey("security.protocol", c.SecurityProtocol)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding security.protocol with value %d "+
				"failed with error: %v", wu.FuncName(), c.SecurityProtocol, err)
			return errors.New(errorStr)
		}

		if c.SecurityProtocol == KAFKA_SECURITY_PROTOCOL_SSL {
			if c.SSLCipherSuites != "" {
				err = m.SetKey("ssl.cipher.suites", c.SSLCipherSuites)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.cipher.suites with value"+
						" %s failed with error: %v", wu.FuncName(),
						c.SSLCipherSuites, err)
					return errors.New(errorStr)
				}
			}

			if c.SSLKeysLocation != "" {
				err = m.SetKey("ssl.key.location", c.SSLKeysLocation)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.key.location with value"+
						" %s failed with error: %v", wu.FuncName(),
						c.SSLKeysLocation, err)
					return errors.New(errorStr)
				}
			}

			if c.SSLKeyPass != "" {
				err = m.SetKey("ssl.key.password", c.SSLKeyPass)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.key.password with value"+
						" %s failed with error: %v", wu.FuncName(),
						c.SSLKeyPass, err)
					return errors.New(errorStr)
				}
			}

			if c.SSLKeysLocation != "" {
				err = m.SetKey("ssl.certificate.location", c.SSLKeysLocation)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.certificate.location "+
						"with value %s failed with error: %v", wu.FuncName(),
						c.SSLKeysLocation, err)
					return errors.New(errorStr)
				}
			}

			if c.SSLCaLocation != "" {
				err = m.SetKey("ssl.ca.location", c.SSLCaLocation)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.ca.location with value "+
						"%s failed with error: %v", wu.FuncName(), c.SSLCaLocation,
						err)
					return errors.New(errorStr)
				}
			}

			if c.SSLCrlLocation != "" {
				err = m.SetKey("ssl.crl.location", c.SSLCrlLocation)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.crl.location with value"+
						" %s failed with error: %v", wu.FuncName(), c.SSLCrlLocation,
						err)
					return errors.New(errorStr)
				}
			}

			if c.SSLCrlLocation != "" {
				err = m.SetKey("ssl.crl.location", c.SSLCrlLocation)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding ssl.crl.location with value"+
						" %s failed with error: %v", wu.FuncName(), c.SSLCrlLocation,
						err)
					return errors.New(errorStr)
				}
			}
		} else if c.SecurityProtocol == KAFKA_SECURITY_PROTOCOL_SASL_PLAINTEXT ||
			c.SecurityProtocol == KAFKA_SECURITY_PROTOCOL_SASL_SSL {
			if c.SaslMechanisms != "" {
				err = m.SetKey("sasl.mechanisms", c.SaslMechanisms)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.mechanisms with value"+
						" %s failed with error: %v", wu.FuncName(),
						c.SaslMechanisms, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslKeberosService != "" {
				err = m.SetKey("sasl.kerberos.service.name", c.SaslKeberosService)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.kerberos.service.name "+
						"with value %s failed with error: %v", wu.FuncName(),
						c.SaslKeberosService, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslKeberosPrincipal != "" {
				err = m.SetKey("sasl.kerberos.principal", c.SaslKeberosPrincipal)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.kerberos.principal "+
						"with value %s failed with error: %v", wu.FuncName(),
						c.SaslKeberosPrincipal, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslKeberosKinitCmd != "" {
				err = m.SetKey("sasl.kerberos.kinit.cmd", c.SaslKeberosKinitCmd)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.kerberos.kinit.cmd "+
						"with value %s failed with error: %v", wu.FuncName(),
						c.SaslKeberosKinitCmd, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslKeberosKeyTab != "" {
				err = m.SetKey("sasl.kerberos.keytab", c.SaslKeberosKeyTab)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.kerberos.keytab "+
						"with value %s failed with error: %v", wu.FuncName(),
						c.SaslKeberosKeyTab, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslKeberosLoginInterval > 0 {
				err = m.SetKey("sasl.kerberos.min.time.before.relogin",
					c.SaslKeberosLoginInterval)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding "+
						"sasl.kerberos.min.time.before.relogin with value %d "+
						"failed with error: %v", wu.FuncName(),
						c.SaslKeberosLoginInterval, err)
					return errors.New(errorStr)
				}
			}

			if c.SaslUserName != "" {
				err = m.SetKey("sasl.username", c.SaslUserName)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.username with value %s"+
						" failed with error: %v", wu.FuncName(), c.SaslUserName,
						err)
					return errors.New(errorStr)
				}
			}

			if c.SaslPassword != "" {
				err = m.SetKey("sasl.password", c.SaslPassword)
				if err != nil {
					errorStr = fmt.Sprintf("%s: Adding sasl.password with value %s"+
						" failed with error: %v", wu.FuncName(), c.SaslPassword,
						err)
					return errors.New(errorStr)
				}
			}
		}
	} else {
		err = m.SetKey("security.protocol", KAFKA_SECURITY_PROTOCOL_PLAINTEXT)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding security.protocol with value %d "+
				"failed with error: %v", wu.FuncName(), c.SecurityProtocol, err)
			return errors.New(errorStr)
		}
	}

	return nil
}

func parseProducerConfig(conf *KafkaPluginConfig, m *kafka.ConfigMap) error {
	var err error
	var errorStr string

	log.Printf("%s: Parsing common config for producer", wu.FuncName())
	if err = parseCommonConfig(&conf.CommonConfig, m,
		KAFKA_TYPE_PRODUCER); err != nil {
		log.Printf("%s: Parsing common config returned error", wu.FuncName())
		return err
	}

	c := conf.ProducerConfig

	if c.MaxQueueBufferMessages > 0 {
		err = m.SetKey("queue.buffering.max.messages", c.MaxQueueBufferMessages)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding queue.buffering.max.messages "+
				"with value %d failed with error: %v", wu.FuncName(),
				c.MaxQueueBufferMessages, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxQueueBufferKbytes > 0 {
		err = m.SetKey("queue.buffering.max.kbytes", c.MaxQueueBufferKbytes)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding queue.buffering.max.kbytes "+
				"with value %d failed with error: %v", wu.FuncName(),
				c.MaxQueueBufferKbytes, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxQueueBuffertime > 0 {
		err = m.SetKey("queue.buffering.max.ms", c.MaxQueueBuffertime)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding queue.buffering.max.ms "+
				"with value %d failed with error: %v", wu.FuncName(),
				c.MaxQueueBuffertime, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxSendRetries > 0 {
		err = m.SetKey("message.send.max.retries", c.MaxSendRetries)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding message.send.max.retries "+
				"with value %d failed with error: %v", wu.FuncName(),
				c.MaxSendRetries, err)
			return errors.New(errorStr)
		}
	}

	if c.RetryBackoff > 0 {
		err = m.SetKey("retry.backoff.ms", c.RetryBackoff)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding retry.backoff.ms with value %d "+
				"failed with error: %v", wu.FuncName(), c.RetryBackoff, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxNumBatchMessages > 0 {
		err = m.SetKey("batch.num.messages", c.MaxNumBatchMessages)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding batch.num.messages with value "+
				"%d failed with error: %v", wu.FuncName(),
				c.MaxNumBatchMessages, err)
			return errors.New(errorStr)
		}
	}

	err = m.SetKey("delivery.report.only.error", c.EnableErrDeliveryonly)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding delivery.report.only.error with "+
			"value %d failed with error: %v", wu.FuncName(),
			c.EnableErrDeliveryonly, err)
		return errors.New(errorStr)
	}

	if strings.Compare(c.CompressionCodec, "") != 0 {
		err = m.SetKey("compression.codec", c.CompressionCodec)
	} else {
		err = m.SetKey("compression.codec", COMPRESSION_CODEC_TYPE_NONE)
	}
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding compression.codec with value %d "+
			"failed with error: %v", wu.FuncName(), c.CompressionCodec, err)
		return errors.New(errorStr)
	}

	/*
	 * kafka-go specific options
	 */
	err = m.SetKey("go.batch.producer", c.EnableBatchProducer)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding go.batch.producer with "+
			"value %d failed with error: %v", wu.FuncName(),
			c.EnableBatchProducer, err)
		return errors.New(errorStr)
	}

	err = m.SetKey("go.delivery.reports", !c.DisablePerMsgDeliveryReport)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding go.delivery.reports with "+
			"value %d failed with error: %v", wu.FuncName(),
			!c.DisablePerMsgDeliveryReport, err)
		return errors.New(errorStr)
	}

	if c.ProducerChannelSize > 0 {
		err = m.SetKey("go.produce.channel.size", c.ProducerChannelSize)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding go.produce.channel.size with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.ProducerChannelSize, err)
			return errors.New(errorStr)
		}
	}

	/*
	 * Topic level configurations
	 */
	if c.Acks > 0 {
		err = m.SetKey("request.required.acks", c.Acks)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding request.required.acks with "+
				"value %d failed with error: %v", wu.FuncName(), c.Acks, err)
			return errors.New(errorStr)
		}
	}

	if c.MessageTimeout > 0 {
		err = m.SetKey("request.timeout.ms", c.MessageTimeout)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding request.timeout.ms with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.MessageTimeout, err)
			return errors.New(errorStr)
		}
	}

	err = m.SetKey("produce.offset.report", c.ProduceOffsetReport)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding produce.offset.report with "+
			"value %d failed with error: %v", wu.FuncName(),
			c.ProduceOffsetReport, err)
		return errors.New(errorStr)
	}

	if c.RequestTimeout > 0 {
		err = m.SetKey("message.timeout.ms", c.RequestTimeout)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding message.timeout.ms with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.RequestTimeout, err)
			return errors.New(errorStr)
		}
	}

	return nil
}

func parseConsumerConfig(conf *KafkaPluginConfig, m *kafka.ConfigMap) error {
	var err error
	var errorStr string

	log.Printf("%s: Parsing common config for consumer", wu.FuncName())

	if err = parseCommonConfig(&conf.CommonConfig, m,
		KAFKA_TYPE_CONSUMER); err != nil {
		log.Printf("%s: Parsing common config returned error", wu.FuncName())
		return err
	}

	c := conf.ConsumerConfig

	err = m.SetKey("enable.auto.commit", !c.DisableAutoCommit)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding auto.commit.enable with "+
			"value %d failed with error: %v", wu.FuncName(),
			!c.DisableAutoCommit, err)
		return errors.New(errorStr)
	}

	if c.AutoCommitInterval > 0 {
		err = m.SetKey("auto.commit.interval.ms", c.AutoCommitInterval)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding auto.commit.interval.ms with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.AutoCommitInterval, err)
			return errors.New(errorStr)
		}
	}

	err = m.SetKey("enable.auto.offset.store", !c.DisableAutoOffset)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding enable.auto.offset.store with "+
			"value %d failed with error: %v", wu.FuncName(),
			!c.DisableAutoOffset, err)
		return errors.New(errorStr)
	}

	if c.MinPartitionMsgs > 0 {
		err = m.SetKey("queued.min.messages", c.MinPartitionMsgs)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding queued.min.messages with "+
				"value %d failed with error: %v", wu.FuncName(),
				c.MinPartitionMsgs, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxKbPerPartition > 0 {
		err = m.SetKey("queued.max.messages.kbytes", c.MaxKbPerPartition)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding queued.max.messages.kbytes with"+
				" value %d failed with error: %v", wu.FuncName(),
				c.MaxKbPerPartition, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxFetchWait > 0 {
		err = m.SetKey("fetch.wait.max.ms", c.MaxFetchWait)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding fetch.wait.max.ms with"+
				" value %d failed with error: %v", wu.FuncName(),
				c.MaxFetchWait, err)
			return errors.New(errorStr)
		}
	}

	if c.MaxMsgFetchBytes > 0 {
		err = m.SetKey("fetch.message.max.bytes", c.MaxMsgFetchBytes)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding fetch.message.max.bytes with"+
				" value %d failed with error: %v", wu.FuncName(),
				c.MaxMsgFetchBytes, err)
			return errors.New(errorStr)
		}
	}

	if c.MinFetchBytes > 0 {
		err = m.SetKey("fetch.min.bytes", c.MinFetchBytes)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding fetch.min.bytes with"+
				" value %d failed with error: %v", wu.FuncName(),
				c.MinFetchBytes, err)
			return errors.New(errorStr)
		}
	}

	if c.FetchErrorBackoff > 0 {
		err = m.SetKey("fetch.error.backoff.ms", c.FetchErrorBackoff)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding fetch.error.backoff.ms with"+
				" value %d failed with error: %v", wu.FuncName(),
				c.FetchErrorBackoff, err)
			return errors.New(errorStr)
		}
	}

	err = m.SetKey("enable.partition.eof", !c.DisablePartitionEof)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding enable.partition.eof with "+
			"value %d failed with error: %v", wu.FuncName(),
			!c.DisablePartitionEof, err)
		return errors.New(errorStr)
	}

	err = m.SetKey("check.crcs", c.CheckCRCs)
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding check.crcs with value %d failed "+
			"with error: %v", wu.FuncName(), c.CheckCRCs, err)
		return errors.New(errorStr)
	}

	if strings.Compare(c.AutoOffsetResetType, "") != 0 {
		err = m.SetKey("auto.offset.reset", c.AutoOffsetResetType)
	} else {
		err = m.SetKey("auto.offset.reset", AUTO_OFFSET_RESET_TYPE_LARGEST)
	}
	if err != nil {
		errorStr = fmt.Sprintf("%s: Adding auto.offset.reset with value"+
			" %d failed with error: %v", wu.FuncName(),
			c.AutoOffsetResetType, err)
		return errors.New(errorStr)
	}

	if strings.Compare(c.OffsetStorePath, "") != 0 {
		err = m.SetKey("offset.store.path", c.OffsetStorePath)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding offset.store.path with value %s"+
				" failed with error: %v", wu.FuncName(), c.OffsetStorePath, err)
			return errors.New(errorStr)
		}
	}

	if c.OffsetStoreSyncInterval > 0 {
		err = m.SetKey("offset.store.sync.interval.ms", c.OffsetStoreSyncInterval)
		if err != nil {
			errorStr = fmt.Sprintf("%s: Adding offset.store.sync.interval.ms "+
				"with value %s failed with error: %v", wu.FuncName(),
				c.OffsetStoreSyncInterval, err)
			return errors.New(errorStr)
		}
	}

	return nil
}

/*
 * Function to create a Kafka producer
 */
func createKafkaProducer(config *KafkaPluginConfig) (*KafkaProducer, error) {
	var err error
	var m kafka.ConfigMap = make(kafka.ConfigMap)
	var kp *KafkaProducer = new(KafkaProducer)

	/*
	 * populate config map parameters
	 * for producers
	 */
	if strings.Compare(config.CommonConfig.ClientId, "") == 0 {
		config.CommonConfig.ClientId = DEFAULT_KAFKA_CLIENT_ID
	} else {
		log.Printf("%s: config.ClientId is %s", wu.FuncName(),
			config.CommonConfig.ClientId)
	}

    kp.clientId = config.CommonConfig.ClientId

	if err = parseProducerConfig(config, &m); err != nil {
		return nil, err
	}

	/*
	 * Create a kafka producer object
	 */
	if kp.kprod, err = kafka.NewProducer(&m); err != nil {
		return nil, err
	} else {
		log.Printf("%s: Created a producer with bootstrap server %s",
			wu.FuncName(), config.CommonConfig.BootStrapServers)
	}

	return kp, nil
}

/*
 * Function to run the publisher and process the messages
 */
func (kp *KafkaProducer) startKafkaProducer(newMsgChan <-chan PublishMsg,
	errChan chan<- *kafka.Message, termChan <-chan struct{}) {

	for {
		select {
		case <-termChan:
			kp.destroyProducer()
			return
		case input, ok := <-newMsgChan:
			if !ok {
				log.Printf("%s: Input message channel has been closed",
					wu.FuncName())
				break
			}

			message := kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &input.Topic,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(input.Key),
				Value: []byte(input.Payload),
			}

			kp.kprod.ProduceChannel() <- &message
		case e := <-kp.kprod.Events():
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Printf("%s: Delivery failed: %v\n", wu.FuncName(),
						m.TopicPartition.Error)
					errChan <- m
					//} else {
					//	log.Printf("%s: Delivered message to topic %s\n",
					//		wu.FuncName(), *m.TopicPartition.Topic)
				}
			}
		}
	}
}

/*
 * Method to destroy a kafka producer
 */
func (kp *KafkaProducer) destroyProducer() {
	kp.kprod.Close()
}

/*
 * Function to create a Kafka consumer
 */
func createKafkaConsumer(config *KafkaPluginConfig,
	eventChanSize uint32) (*KafkaConsumer, error) {
	var err error
	var m kafka.ConfigMap = make(kafka.ConfigMap)
	var kc *KafkaConsumer = new(KafkaConsumer)

	/*
	 * populate config map parameters
	 * for producers
	 */
	if strings.Compare(config.CommonConfig.ClientId, "") == 0 {
		config.CommonConfig.ClientId = "wedge"
	} else {
		log.Printf("%s: config.ClientId is %s", wu.FuncName(),
			config.CommonConfig.ClientId)
	}

    kc.clientId = config.CommonConfig.ClientId

	if err = parseConsumerConfig(config, &m); err != nil {
		return nil, err
	}

	/*
	 * Set
	 */
	if err = m.SetKey("default.topic.config",
		kafka.ConfigMap{"auto.offset.reset": "earliest"}); err != nil {
		return nil, err
	}

	/*
	 * kafka-go specific options
	 */
	if err = m.SetKey("go.application.rebalance.enable", true); err != nil {
		return nil, err
	}

	if err = m.SetKey("go.events.channel.enable", true); err != nil {
		return nil, err
	}

	if eventChanSize > 0 {
		if err = m.SetKey("go.events.channel.size", eventChanSize); err != nil {
			return nil, err
		}
	}

	/*
	 * Create a kafka consumer object
	 */
	if kc.kcon, err = kafka.NewConsumer(&m); err != nil {
		return nil, err
	}

	return kc, nil
}

/*
 * Function to start the consumer
 */
func (kc *KafkaConsumer) startKafKaConsumer(newMsgChan chan<- *kafka.Message,
	errChan chan<- string, termChan <-chan struct{}) {
	var err error

	/*
	 * Start the initial subscription
	 */
	kc.topics = append(kc.topics, KAFKA_ERROR_TOPIC,
		"_wedge.Topic_registerTopic")

	if err = kc.kcon.SubscribeTopics(kc.topics, nil); err != nil {
		log.Printf("%s: Initial subscription to topic _wedge.Topic_error"+
			"failed with error %v", wu.FuncName(), err)
		errChan <- CONSUMER_TOPIC_INIT_FAILED
		return
	}

	for {
		select {
		case <-termChan:
			close(newMsgChan)
			close(errChan)
			kc.destroyConsumer()
			return
		case ev := <-kc.kcon.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				log.Printf("%s: %v\n", wu.FuncName(), e)
				kc.kcon.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				log.Printf("%s %v\n", wu.FuncName(), e)
				kc.kcon.Unassign()
			case *kafka.Message:
				/*
				 * When a topic is added, update the topics list
				 * and create a new subscription
				 */
				if strings.Compare("_wedge.Topic_registerTopic",
					*e.TopicPartition.Topic) == 0 {
					topics := strings.Split(string(e.Value), ",")
					for _, topic := range topics {
						log.Printf("%s: Registering topic %s", wu.FuncName(),
							topic)
						kc.topics = append(kc.topics, topic)
					}
					if err = kc.kcon.SubscribeTopics(kc.topics,
						nil); err != nil {
						errChan <- fmt.Sprintf("%v", err)
					}
				} else {
					newMsgChan <- e
				}
				kc.kcon.Commit()
			case kafka.PartitionEOF:
				log.Printf("%s: End of partition error with details %v\n",
					wu.FuncName(), e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%s: Error: %v\n", wu.FuncName(), e)
				errChan <- CONSUMER_TERMINATED
				return
			}
		}
	}
}

/*
 * Method to destroy a kafka consumer
 */
func (kc *KafkaConsumer) destroyConsumer() {
	kc.kcon.Close()
}

/*
 * Kafka input plugin structure to implement the ClientSidePlugin
 * interface for Kafka
 */
type KafKaPlugin struct {
	kp *KafkaProducer
	kc *KafkaConsumer
	fr wu.PostReadFunc
	fw wu.PreWriteFunc
	fe wu.ErrorFunc
}

type KafkaPluginConfig struct {
	AvroDescTableLoc  string              `yaml:"avro-desc-table-loc"`
	AvscFileDir       string              `yaml:"avsc-file-dir"`
	CommonConfig      KafkaCommonConfig   `yaml:"common-config"`
	ConsumerConfig    KafkaConsumerConfig `yaml:"consumer-config"`
	PayloadType       string              `yaml:"payload-type"`
	ProducerConfig    KafkaProducerConfig `yaml:"producer-config"`
	ProtoDescTableLoc string              `yaml:"proto-desc-table-loc"`
}

func (k *KafKaPlugin) GetErrorKey() string {
	var x string
	x = "{\n\t" + "\"ClientId\": " + "\"" + k.kc.clientId + "\"" + "\n}"
	return x
}

func (k *KafKaPlugin) RunKafKaPlugin(codecType int, params interface{},
	sendChan chan<- wu.MsgFormat, recvChan <-chan wu.MsgFormat,
	termChan <-chan struct{}, statusChan chan<- wu.TermStatus) {
	var err error
	var ret wu.TermStatus
	var termStatus int
	var termStr string
	var conf *KafkaPluginConfig = (params.(*KafkaPluginConfig))
	var avroDescLoc = conf.AvroDescTableLoc
	var protoDescLoc = conf.ProtoDescTableLoc
	var payloadType = conf.PayloadType
	var avroDir = conf.AvscFileDir

	log.Printf("%s: Producer config: bootstrap server: %s", wu.FuncName(),
		conf.CommonConfig.BootStrapServers)

	log.Printf("%s: Consumer config: bootstrap server: %s, group id %s",
		wu.FuncName(), conf.CommonConfig.BootStrapServers,
		conf.CommonConfig.GroupId)

	if payloadType == "avro" {
		if avroDescLoc == "" {
			ret.Status = wu.ROUTINE_ERROR_EXIT
			ret.ErrStr = "avro_desc_table_loc cannot be empty when payload is avro"
			statusChan <- ret
			return
		}

		/*
		 * Initialize the avro protocol descriptor
		 */
		wu.InitAvroParser(avroDescLoc, avroDir)

	} else {
		if protoDescLoc == "" {
			ret.Status = wu.ROUTINE_ERROR_EXIT
			ret.ErrStr = "proto_desc_table_loc cannot be empty"
			statusChan <- ret
			return
		}

		wu.InitProtoParser(protoDescLoc)
	}

	/*
	 * Create the producer and consumer objects
	 */
	k.kp, err = createKafkaProducer(conf)
	if err != nil {
		ret.Status = wu.ROUTINE_ERROR_EXIT
		ret.ErrStr = fmt.Sprintf("%v", err)
		statusChan <- ret
	}

	k.kc, err = createKafkaConsumer(conf, 0)
	if err != nil {
		ret.Status = wu.ROUTINE_ERROR_EXIT
		ret.ErrStr = fmt.Sprintf("%v", err)
		statusChan <- ret
	}

	/*
	 * Create terminate channels for client reader and
	 * writer
	 */
	termPluginChan := make(chan struct{}, 2)
	feedBackChan := make(chan interface{}, 2)
	readerStatus := make(chan wu.TermStatus, 2)
	writerStatus := make(chan wu.TermStatus, 2)

	log.Printf("%s: starting runReader and runWriter", wu.FuncName())

	go k.runReader(sendChan, feedBackChan, termPluginChan, readerStatus)
	go k.runWriter(recvChan, feedBackChan, termPluginChan, writerStatus)

	/*
	 * Wait till the Plugin is asked to terminate
	 */
	for {
		select {
		case <-termChan:
			close(termPluginChan)
			ret.ErrStr = ""
			ret.Status = wu.ROUTINE_NORMAL_EXIT
			statusChan <- ret
			close(statusChan)
			return
		case ret = <-readerStatus:
			if ret.Status == wu.ROUTINE_ERROR_EXIT {
				log.Printf("%s: Reader terminated with error: %s", wu.FuncName(),
					ret.ErrStr)
				termStr = termStr + ret.ErrStr
				termStatus |= ret.Status
			}
			close(termPluginChan)
			ret.ErrStr = termStr
			ret.Status = termStatus
			statusChan <- ret
			close(statusChan)
			return

		case ret = <-writerStatus:
			if ret.Status == wu.ROUTINE_ERROR_EXIT {
				log.Printf("%s: Writer terminated with error: %s", wu.FuncName(),
					ret.ErrStr)
				termStr = termStr + ret.ErrStr
				termStatus |= ret.Status
			}
			close(termPluginChan)
			ret.ErrStr = termStr
			ret.Status = termStatus
			statusChan <- ret
			close(statusChan)
			return
		}
	}
}

func (k *KafKaPlugin) runReader(sendChan chan<- wu.MsgFormat,
	feedBackChan chan<- interface{}, termChan <-chan struct{},
	statusChan chan<- wu.TermStatus) {
	var err error
	var input *kafka.Message
	var errorStr string
	var msg wu.MsgFormat
	var root *jsonez.GoJSON
	var feedBack PublishMsg
	var ok bool
	var ret wu.TermStatus
	var startCount int // To keep track of how many restarts were performed
	var inProcessed []wu.MsgFormat
	var val interface{}

	log.Printf("%s: k is %v", wu.FuncName(), k)

	clientReader := make(chan *kafka.Message, 100)
	errChan := make(chan string, 2)
	termConsumer := make(chan struct{}, 2)

	go k.kc.startKafKaConsumer(clientReader, errChan, termConsumer)

	for {
		select {
		case <-termChan:
			close(termConsumer)
			close(feedBackChan)
			close(sendChan)
			ret.ErrStr = ""
			ret.Status = wu.ROUTINE_NORMAL_EXIT
			statusChan <- ret
			close(statusChan)
			return
		case input, ok = <-clientReader:
			if !ok {
				log.Printf("%s: clientReader channel has been closed",
					wu.FuncName())
				break
			}

			log.Printf("%s: Consumer received message with \ntopic: %s,\n "+
				"key: %s,\n\n", wu.FuncName(),
				*input.TopicPartition.Topic, string(input.Key))

			/*
			 * If the message was sent to KAFKA_ERROR_TOPIC
			 * drop it
			 */
			if strings.Compare(*input.TopicPartition.Topic,
				KAFKA_ERROR_TOPIC) == 0 {
				break
			}

			log.Printf("%s: CALLBACK IS %v *******", wu.FuncName(), k.fr)

			if k.fr != nil {
				log.Printf("%s: ********* INVOKING CALLBACK %v ****************** ",
					wu.FuncName(), k.fr)
				inProcessed, err = k.fr(input)

				if err != nil {
					log.Printf("%s: Post read callback execution returned "+
						"error %v", wu.FuncName(), err)

					if k.fe != nil {
						val, err = k.fe(&msg, input, err.Error())
						feedBack = *val.(*PublishMsg)
					}

					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte(err.Error())
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}

					feedBackChan <- feedBack

					break
				}

				for _, msg = range inProcessed {
					/*
					 * If the BrokerId field is speficied, then compare it with
					 * the client id for broker and forward it to the plugin on
					 * other side if there is a match
					 */
					 if msg.BrokerId != "" && strings.Compare(k.kc.clientId, msg.BrokerId) == 0 {
					 	sendChan <- msg
					 }
				}
			} else {
				log.Printf("%s: ************* CALLBACK IS NULL **************", wu.FuncName())
				if root, err = jsonez.GoJSONParse(input.Key); err != nil {
					/*
							 * Error parsing the key, Publish the Error though
						     * the topic KAFKA_ERROR_TOPIC
					*/
					log.Printf("%s: Error %v while parsing key %s", wu.FuncName(),
						err, string(input.Key))
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Error parsing key")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Error parsing key")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}

				clientId, err := root.Get(wu.MSG_FORMAT_CLIENT_ID)
				if err != nil {
					log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
						wu.MSG_FORMAT_CLIENT_ID, err)
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Mandatory parameter "+
							"ClientId not found")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Mandatory parameter ClientId " +
							"not found")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}
				msg.ClientId = clientId.Valstr
				log.Printf("%s: msg.ClientId is %s", wu.FuncName(), msg.ClientId)

				transactId, err := root.Get(wu.MSG_FORMAT_TRANSACTION_ID)
				if err != nil {
					log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
						wu.MSG_FORMAT_TRANSACTION_ID, err)
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Mandatory parameter "+
							"TransactionId not found")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Mandatory parameter " +
							"TransactionId not found")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}
				msg.TransactionId = transactId.Valstr
				log.Printf("%s: msg.TransactionId is %s", wu.FuncName(),
					msg.TransactionId)

				rpcId, err := root.Get(wu.MSG_FORMAT_RPC_ID)
				if err != nil {
					log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
						wu.MSG_FORMAT_RPC_ID, err)
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Mandatory parameter "+
							"RpcId not found")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Mandatory parameter RpcId not " +
							"found")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}
				msg.RpcId = rpcId.Valstr
				log.Printf("%s: msg.RpcId is %s", wu.FuncName(),
					msg.RpcId)

				ipAddr, err := root.Get(wu.MSG_FORMAT_IPADDRESS)
				if err != nil {
					log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
						wu.MSG_FORMAT_IPADDRESS, err)
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Mandatory parameter "+
							"IpAddress not found")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Mandatory parameter IpAddress " +
							"not found")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}
				msg.IpAddress = ipAddr.Valstr
				log.Printf("%s: msg.IpAddress is %s", wu.FuncName(),
					msg.IpAddress)

				port, err := root.Get(wu.MSG_FORMAT_PORT)
				if err != nil {
					log.Printf("%s: Fetching %s returned error %v", wu.FuncName(),
						wu.MSG_FORMAT_PORT, err)
					if k.fe != nil {
						val, err = k.fe(&msg, input, "Mandatory parameter "+
							"IpAddress not found")
						feedBack = *val.(*PublishMsg)
					}
					if k.fe == nil || err != nil {
						feedBack.Key = input.Key
						feedBack.Payload = []byte("Mandatory parameter Port not " +
							"found")
						feedBack.Topic = KAFKA_ERROR_TOPIC
					}
					feedBackChan <- feedBack
					break
				}
				msg.Port = port.Valstr
				log.Printf("%s: msg.Port is %s", wu.FuncName(), msg.Port)

				metadata, err := root.Get(wu.MSG_FORMAT_METADATA)
				if metadata != nil && err != nil {
					for entry := metadata.Child; entry != nil; entry = entry.Next {
						msg.Metadata[entry.Child.Key] = entry.Child.Valstr
					}
				}

				/*
				 * If the clientId of the key matches wedge's clientId,
				 * then drop it as it a copy of the message that was
				 * produced by wedge
				 */
				if strings.Compare(k.kc.clientId, msg.ClientId) != 0 &&
					strings.Compare(k.kp.clientId, msg.ClientId) != 0 {
					msg.Rpc = strings.Replace(*input.TopicPartition.Topic,
						"_", "/", -1)

					/*
					 * If the RPC is for a call cancellation, then
					 * change the value i.e the RPC to be terminated
					 * to the required format
					 */
					if strings.Compare(msg.Rpc, wu.CALL_CANCELLATION_RPC) == 0 {
						msg.Value = strings.Replace(string(input.Value), "_", "/",
							-1)
					} else {
						msg.Value = string(input.Value)
					}

					/*
					 * Write the message to the server side plugin
					 */
					sendChan <- msg
				}
			}

			break
		case errorStr = <-errChan:
			if strings.Compare(errorStr, CONSUMER_TERMINATED) == 0 ||
				strings.Compare(errorStr, CONSUMER_TOPIC_INIT_FAILED) == 0 {
				if startCount > 5 {

					close(termConsumer)
					close(feedBackChan)
					close(sendChan)
					ret.Status = wu.ROUTINE_ERROR_EXIT
					ret.ErrStr = fmt.Sprintf("Kafka consumer failure")
					statusChan <- ret
					close(statusChan)
					return
				}
				startCount++
				go k.kc.startKafKaConsumer(clientReader, errChan, termConsumer)
			}
		}
	}
}

func (k *KafKaPlugin) runWriter(recvChan <-chan wu.MsgFormat,
	feedBackChan <-chan interface{}, termChan <-chan struct{},
	statusChan chan<- wu.TermStatus) {

	var err error
	var errorStr string
	var ok bool
	var feedBack interface{}
	var msg wu.MsgFormat
	var pub *PublishMsg
	var root *jsonez.GoJSON
	var ret wu.TermStatus
	var outProcessed []interface{}
	var wfeedBack PublishMsg

	clientWriter := make(chan PublishMsg, 100)
	errChan := make(chan *kafka.Message, 2)
	termProducer := make(chan struct{}, 2)

	go k.kp.startKafkaProducer(clientWriter, errChan, termProducer)

	for {
		select {
		case <-termChan:
			close(termProducer)
			close(clientWriter)
			ret.ErrStr = ""
			ret.Status = wu.ROUTINE_NORMAL_EXIT
			statusChan <- ret
			close(statusChan)
			return
		case feedBack, ok = <-feedBackChan:
			if !ok {
				log.Printf("%s: Feedback channel has been closed",
					wu.FuncName())
				errorStr = fmt.Sprintf("%s: Feedback channel has been closed")
				close(termProducer)
				close(clientWriter)
				ret.ErrStr = errorStr
				ret.Status = wu.ROUTINE_ERROR_EXIT
				statusChan <- ret
				close(statusChan)
				return
			}

			log.Printf("%s: Received message in feedback channel", wu.FuncName())
			clientWriter <- feedBack.(PublishMsg)
		case msg, ok = <-recvChan:
			if !ok {
				log.Printf("%s: Recv channel from server side has been closed",
					wu.FuncName())
				errorStr = fmt.Sprintf("Recv channel from server side has " +
					"been closed")
				close(termProducer)
				close(clientWriter)
				ret.ErrStr = errorStr
				ret.Status = wu.ROUTINE_ERROR_EXIT
				statusChan <- ret
			}

			log.Printf("%s: Got message with Id: %s, Ipaddress: %s, "+
				"Metadata :%s, Port: %s, Rpc: %s\n\n", wu.FuncName(),
				msg.ClientId, msg.IpAddress, msg.Metadata, msg.Port,
				msg.Rpc)

			if k.fw != nil {
				log.Printf("%s: *********INVOKING CALLBACK %v ****************** ",
					wu.FuncName(), k.fw)
				outProcessed, err = k.fw(msg)

				if err != nil {
					log.Printf("%s: Pre write callback execution returned "+
						"error %v", wu.FuncName(), err)
					wfeedBack.Key = []byte(k.GetErrorKey())
					wfeedBack.Payload = []byte(err.Error())
					wfeedBack.Topic = KAFKA_ERROR_TOPIC
					clientWriter <- wfeedBack
					break
				}

				for _, val := range outProcessed {
					pub = val.(*PublishMsg)
					clientWriter <- *pub
				}
			} else {
				log.Printf("%s: ********* CALLBACK is NULL ****************** ",
					wu.FuncName())
				pub.Topic = msg.Rpc
				pub.Payload = []byte(msg.Value.(string))

				/*
				 * Make a Key for the message
				 */
				root = jsonez.AllocObject()
				/*
				 * Set the client address to that of wedge kafka
				 * producer
				 */
				root.AddVal(k.kp.clientId, wu.MSG_FORMAT_CLIENT_ID)
				root.AddVal(msg.TransactionId, wu.MSG_FORMAT_TRANSACTION_ID)
				root.AddVal(msg.IpAddress, wu.MSG_FORMAT_IPADDRESS)
				root.AddVal(msg.Port, wu.MSG_FORMAT_PORT)

				if len(msg.Metadata) > 0 {
					for k, v := range msg.Metadata {
						obj := jsonez.AllocObject()
						obj.AddEntryToObject(k, jsonez.AllocString(v))
						err = root.AddToArray(obj, wu.MSG_FORMAT_METADATA)
						if err != nil {
							log.Printf("%s: Adding Metadata with key %s and "+
								"value %s failed with error %v\n", wu.FuncName(),
								k, v, err)
						}
					}
				}
				pub.Key = jsonez.GoJSONPrint(root)

				/*
				 * Kafka doesn't accept '/'' in topic names. So replace
				 * all '/' with '_'
				 */
				pub.Topic = strings.Replace(pub.Topic, "/", "_", -1)

				/*
				 * Write to the kafka publisher
				 */
				clientWriter <- *pub
			}
		}
	}
}

func (k KafKaPlugin) RunClientPlugin(codecType int, params interface{},
	readFunc wu.PostReadFunc, writeFunc wu.PreWriteFunc, errorFunc wu.ErrorFunc,
	sendChan chan<- wu.MsgFormat, recvChan <-chan wu.MsgFormat,
	termChan <-chan struct{}, statusChan chan<- wu.TermStatus) {
	log.Printf("%s: k is %v", wu.FuncName(), k)
	k.fr = readFunc
	k.fw = writeFunc
	k.fe = errorFunc
	k.RunKafKaPlugin(codecType, params, sendChan, recvChan, termChan,
		statusChan)

}

var KafkaConfigPlugin wu.Plugin

func kafkaParseConfig(data []byte) interface{} {
	var kconf KafkaPluginConfig

	err := yaml.Unmarshal([]byte(data), &kconf)
	if err != nil {
		log.Fatalf("%s: Unmarshalling yaml data failed with error %v",
			wu.FuncName(), err)
	}

	return &kconf
}

func kafkaFetchErrorFunc(params interface{}) string {
	var conf *KafkaPluginConfig = (params.(*KafkaPluginConfig))

	if conf.PayloadType == "avro" {
		return "avro"
	} else {
		return ""
	}
}

func kafkaEmitRefConfig() string {
	return `
            name: "kafka"
            client-post-read-cb: "" # Custom callback to run after data is received from client
            client-pre-write-cb: "" # Custom callback to run before data is sent to the client
            server-post-read-cb: "" # Custom callback to run after data is received from server
            server-pre-write-cb: "" # Custom callback to run before data is sent to the server
            config:
                # Details on kafka configuration options and default values can be found in:
                # https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md    
                common-config:
                    address-family: 0
                    backoff-jitter: 0 
                    blacklist-topics: [] 
                    bootstrap-servers: "" # MANDATORY, Initial list of brokers
                    broker-address-ttl: 0 
                    broker-client-id: "" # MANDATORY, client-id used by broker to publish its messages
                    disable-nagle: false 
                    disable-sparse-metadata-req: false
                    enable-tcp-keepalives: false 
                    group-id: "" # MANDATORY, All clients sharing the same group.id belong to the same group
                    max-msg-bytes: 0 
                    max-copy-bytes: 0
                    max-inflight-requests: 0
                    metadata-req-timeout: 0
                    metadata-refresh-interval: 0
                    metadata-max-age: 0
                    metadata-refresh-fast-interval: 0
                    socket-timeout: 0
                    socket-send-bufsize: 0
                    socket-recv-bufsize: 0
                    max-send-fails: 0
                    recv-max-bytes: 0
                    security-protocol: ""  # MANDATORY, Protocol used to communicate with brokers, default plaintext
                    ssl-cipher-suites: ""
                    ssl-keys-location: ""
                    ssl-key-pass: ""
                    ssl-cert-location: ""
                    ssl-ca-location: ""
                    ssl-crl-location: ""
                    sasl-mechanisms: ""
                    sasl-keberos-service: ""
                    sasl-keberos-principal: ""
                    sasl-keberos-kinit-cmd: ""
                    sasl-keberos-keytab: ""
                    sasl-keberos-login-interval: 0
                    sasl-username: ""
                    sasl-password: ""
                producer-config:
                    max-queue-buffer-messages: 0
                    max-queue-buffer-kbytes: 0
                    max-queue-buffer-time: 0
                    max-send-retries: 0
                    retry-backoff: 0
                    compression-codec: ""
                    max-num-batch-messages: 0
                    enable-err-delivery-only: false
                    enable-batch-producer: false
                    disable-per-msg-delivery-report: false
                    producer-channel-size: 0
                    acks: 0
                    message-timeout: 0
                    produce-offset-report: false
                    request-timeout: 0
                consumer-config:
                    disable-auto-commit: false
                    auto-commit-interval: 0
                    disable-auto-offset: false
                    min-partition-msgs: 0
                    max-kb-per-partition: 0
                    max-fetch-wait: 0
                    max-msg-fetch-bytes: 0
                    min-fetch-bytes: 0
                    fetch-error-backoff: 0
                    disable-partition-eof: false
                    check-crcs: false
                    auto-offset-reset-type: ""
                    offset-store-path: ""
                    offset-store-sync-interval: 0
                payload-type: "json" # MANDATORY, default is "json", other option(s): avro
                avro-desc-table-loc: "" # MANDATORY for avro payload, Location of the avro descriptor table
                                        # generated from proto using protoc-wedge compiler 
                avsc-file-dir: "" # MANDATORY for avro payload, Location of the avsc files that are generated from proto
                                  # using protoc-wedge compiler
                proto-desc-table-loc: "" # MANDATORY for json payload, Location of the proto descriptor table
                                        # generated from proto using protoc-wedge compiler`
}

var plugin KafKaPlugin

func KafkaInit(pluginMap map[string]wu.Plugin) {
	KafkaConfigPlugin.EmitClientSideRefConfig = kafkaEmitRefConfig
	KafkaConfigPlugin.EmitServerSideRefConfig = kafkaEmitRefConfig
	KafkaConfigPlugin.ParseClientSideConfig = kafkaParseConfig
	KafkaConfigPlugin.ClientSidePluginSupport = true
	KafkaConfigPlugin.ServerSidePluginSupport = true
	KafkaConfigPlugin.ClientSideCodec = 0
	KafkaConfigPlugin.ServerSideCodec = 0
	KafkaConfigPlugin.ClientPlugin = plugin
	KafkaConfigPlugin.FetchErrorFunc = kafkaFetchErrorFunc
	pluginMap["kafka"] = KafkaConfigPlugin
}
