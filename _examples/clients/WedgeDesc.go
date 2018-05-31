package WedgeClient

import (
  "encoding/json"
)

type RpcMetadata struct {
  Key string
  Value string
}

type RpcKey struct {
  ClientId	string  //Client id of the transaction
  BrokerId	string //String to uniquely identify the broker that should process the message
  TransactionId	string //ID to uniquely identify this transaction
  RpcId	string //ID to uniquely identify an RPC invokcation in this transaction
  IpAddress	[]string //IP addresses of the devices on which the RPC should be invoked
  Port	string //Port number of the routers
  Metadata	[]RpcMetadata `json:",omitempty"` //Metadata for the RPC
}

/*
 * Result of the operation
 */
const (
  TELEMETRY__SUCCESS = 0
  TELEMETRY__NO__SUBSCRIPTION__ENTRY = 1
  TELEMETRY__UNKNOWN__ERROR = 2
)

/*
 * Verbosity Level
 */
const (
  TELEMETRY__DETAIL = 0
  TELEMETRY__TERSE = 1
  TELEMETRY__BRIEF = 2
)

/*
 * Encoding Type Supported
 */
const (
  TELEMETRY__UNDEFINED = 0
  TELEMETRY__XML = 1
  TELEMETRY__JSON__IETF = 2
  TELEMETRY__PROTO3 = 3
)

/*
 * Collector endpoints to send data specified as an ip+port combination.
 */
type Telemetry__Collector struct {
    Address	 string `json:"address,omitempty"`
    Port	 uint32 `json:"port,omitempty"`
}

/*
 * Data associated with a telemetry subscription
 */
type Telemetry__SubscriptionInput struct {
    CollectorList	 []Telemetry__Collector `json:"collector_list,omitempty"`
}

/*
 * Data model path
 */
type Telemetry__Path struct {
    Path	 string `json:"path,omitempty"`
    Filter	 string `json:"filter,omitempty"`
    SuppressUnchanged	 bool `json:"suppress_unchanged,omitempty"`
    MaxSilentInterval	 uint32 `json:"max_silent_interval,omitempty"`
    SampleFrequency	 uint32 `json:"sample_frequency,omitempty"`
    NeedEom	 bool `json:"need_eom,omitempty"`
}

/*
 * Configure subscription request additional features.
 */
type Telemetry__SubscriptionAdditionalConfig struct {
    LimitRecords	 int32 `json:"limit_records,omitempty"`
    LimitTimeSeconds	 int32 `json:"limit_time_seconds,omitempty"`
    NeedEos	 bool `json:"need_eos,omitempty"`
}

/*
 * Message sent for a telemetry subscription request
 */
type Telemetry__SubscriptionRequest struct {
    Input	 *Telemetry__SubscriptionInput `json:"input,omitempty"`
    PathList	 []Telemetry__Path `json:"path_list,omitempty"`
    AdditionalConfig	 *Telemetry__SubscriptionAdditionalConfig `json:"additional_config,omitempty"`
}

/*
 * Response message to a telemetry subscription creation or get request.
 */
type Telemetry__SubscriptionResponse struct {
    SubscriptionId	 uint32 `json:"subscription_id,omitempty"`
}

type Telemetry__SubscriptionReply struct {
    Response	 *Telemetry__SubscriptionResponse `json:"response,omitempty"`
    PathList	 []Telemetry__Path `json:"path_list,omitempty"`
}

/*
 * Simple Key-value, where value could be one of scalar types
 */
type Telemetry__KeyValue struct {
  /*
   * Only one of the following fields should be specified: 
   * DoubleValue, IntValue, UintValue, SintValue, BoolValue, StrValue, BytesValue
   */
    Key	 string `json:"key,omitempty"`
    DoubleValue	 float64 `json:"double_value,omitempty"`
    IntValue	 int64 `json:"int_value,omitempty"`
    UintValue	 uint64 `json:"uint_value,omitempty"`
    SintValue	 int64 `json:"sint_value,omitempty"`
    BoolValue	 bool `json:"bool_value,omitempty"`
    StrValue	 string `json:"str_value,omitempty"`
    BytesValue	 []byte `json:"bytes_value,omitempty"`
}

type Telemetry__OpenConfigData struct {
    SystemId	 string `json:"system_id,omitempty"`
    ComponentId	 uint32 `json:"component_id,omitempty"`
    SubComponentId	 uint32 `json:"sub_component_id,omitempty"`
    Path	 string `json:"path,omitempty"`
    SequenceNumber	 uint64 `json:"sequence_number,omitempty"`
    Timestamp	 uint64 `json:"timestamp,omitempty"`
    Kv	 []Telemetry__KeyValue `json:"kv,omitempty"`
    Delete	 []Telemetry__Delete `json:"delete,omitempty"`
    Eom	 []Telemetry__Eom `json:"eom,omitempty"`
    SyncResponse	 bool `json:"sync_response,omitempty"`
}

/*
 * Message indicating delete for a particular path
 */
type Telemetry__Delete struct {
    Path	 string `json:"path,omitempty"`
}

/*
 * Message indicating EOM for a particular path
 */
type Telemetry__Eom struct {
    Path	 string `json:"path,omitempty"`
}

/*
 * Message sent for a telemetry subscription cancellation request
 */
type Telemetry__CancelSubscriptionRequest struct {
    SubscriptionId	 uint32 `json:"subscription_id,omitempty"`
}

/*
 * Reply to telemetry subscription cancellation request
 */
type Telemetry__CancelSubscriptionReply struct {
    Code	 uint32 `json:"code,omitempty"`
    CodeStr	 string `json:"code_str,omitempty"`
}

/*
 * Message sent for a telemetry get request
 */
type Telemetry__GetSubscriptionsRequest struct {
    SubscriptionId	 uint32 `json:"subscription_id,omitempty"`
}

/*
 * Reply to telemetry subscription get request
 */
type Telemetry__GetSubscriptionsReply struct {
    SubscriptionList	 []Telemetry__SubscriptionReply `json:"subscription_list,omitempty"`
}

/*
 * Message sent for telemetry agent operational states request
 */
type Telemetry__GetOperationalStateRequest struct {
    SubscriptionId	 uint32 `json:"subscription_id,omitempty"`
    Verbosity	 uint32 `json:"verbosity,omitempty"`
}

/*
 * Reply to telemetry agent operational states request
 */
type Telemetry__GetOperationalStateReply struct {
    Kv	 []Telemetry__KeyValue `json:"kv,omitempty"`
}

/*
 * Message sent for a data encoding request
 */
type Telemetry__DataEncodingRequest struct {
}

/*
 * Reply to data encodings supported request
 */
type Telemetry__DataEncodingReply struct {
    EncodingList	 []uint32 `json:"encoding_list,omitempty"`
}

type Telemetry__OpenConfigTelemetry_telemetrySubscribe struct {
  Request	Telemetry__SubscriptionRequest
  Reply	Telemetry__OpenConfigData
}

func (r *Telemetry__OpenConfigTelemetry_telemetrySubscribe) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/telemetry.OpenConfigTelemetry/telemetrySubscribe"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Telemetry__OpenConfigTelemetry_telemetrySubscribe) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Telemetry__OpenConfigTelemetry_cancelTelemetrySubscription struct {
  Request	Telemetry__CancelSubscriptionRequest
  Reply	Telemetry__CancelSubscriptionReply
}

func (r *Telemetry__OpenConfigTelemetry_cancelTelemetrySubscription) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/telemetry.OpenConfigTelemetry/cancelTelemetrySubscription"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Telemetry__OpenConfigTelemetry_cancelTelemetrySubscription) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Telemetry__OpenConfigTelemetry_getTelemetrySubscriptions struct {
  Request	Telemetry__GetSubscriptionsRequest
  Reply	Telemetry__GetSubscriptionsReply
}

func (r *Telemetry__OpenConfigTelemetry_getTelemetrySubscriptions) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/telemetry.OpenConfigTelemetry/getTelemetrySubscriptions"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Telemetry__OpenConfigTelemetry_getTelemetrySubscriptions) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Telemetry__OpenConfigTelemetry_getTelemetryOperationalState struct {
  Request	Telemetry__GetOperationalStateRequest
  Reply	Telemetry__GetOperationalStateReply
}

func (r *Telemetry__OpenConfigTelemetry_getTelemetryOperationalState) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/telemetry.OpenConfigTelemetry/getTelemetryOperationalState"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Telemetry__OpenConfigTelemetry_getTelemetryOperationalState) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Telemetry__OpenConfigTelemetry_getDataEncodings struct {
  Request	Telemetry__DataEncodingRequest
  Reply	Telemetry__DataEncodingReply
}

func (r *Telemetry__OpenConfigTelemetry_getDataEncodings) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/telemetry.OpenConfigTelemetry/getDataEncodings"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Telemetry__OpenConfigTelemetry_getDataEncodings) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * `Any` contains an arbitrary serialized protocol buffer message along with a
 * URL that describes the type of the serialized message.
 * Protobuf library provides support to pack/unpack Any values in the form
 * of utility functions or additional generated methods of the Any type.
 * Example 1: Pack and unpack a message in C++.
 *     Foo foo = ...;
 *     Any any;
 *     any.PackFrom(foo);
 *     ...
 *     if (any.UnpackTo(&foo)) {
 *       ...
 *     }
 * Example 2: Pack and unpack a message in Java.
 *     Foo foo = ...;
 *     Any any = Any.pack(foo);
 *     ...
 *     if (any.is(Foo.class)) {
 *       foo = any.unpack(Foo.class);
 *     }
 *  Example 3: Pack and unpack a message in Python.
 *     foo = Foo(...)
 *     any = Any()
 *     any.Pack(foo)
 *     ...
 *     if any.Is(Foo.DESCRIPTOR):
 *       any.Unpack(foo)
 *       ...
 * The pack methods provided by protobuf library will by default use
 * 'type.googleapis.com/full.type.name' as the type URL and the unpack
 * methods only use the fully qualified type name after the last '/'
 * in the type URL, for example "foo.bar.com/x/y.z" will yield type
 * name "y.z".
 * JSON
 * ====
 * The JSON representation of an `Any` value uses the regular
 * representation of the deserialized, embedded message, with an
 * additional field `@type` which contains the type URL. Example:
 *     package google.profile;
 *     message Person {
 *       string first_name = 1;
 *       string last_name = 2;
 *     }
 *     {
 *       "@type": "type.googleapis.com/google.profile.Person",
 *       "firstName": <string>,
 *       "lastName": <string>
 *     }
 * If the embedded message type is well-known and has a custom JSON
 * representation, that representation will be embedded adding a field
 * `value` which holds the custom JSON in addition to the `@type`
 * field. Example (for message [google.protobuf.Duration][]):
 *     {
 *       "@type": "type.googleapis.com/google.protobuf.Duration",
 *       "value": "1.212s"
 *     }
 */
type Google__Protobuf__Any struct {
    TypeUrl	 string `json:"type_url,omitempty"`
    Value	 []byte `json:"value,omitempty"`
}



/*
 * The request message containing the user's name, password and client id
 */
type Authentication__LoginRequest struct {
    UserName	 string `json:"user_name,omitempty"`
    Password	 string `json:"password,omitempty"`
    ClientId	 string `json:"client_id,omitempty"`
}

/*
 * The response message containing the result of login attempt.
 * result value of true indicates success and false indicates
 * failure
 */
type Authentication__LoginReply struct {
    Result	 bool `json:"result,omitempty"`
}

type Authentication__Login_LoginCheck struct {
  Request	Authentication__LoginRequest
  Reply	Authentication__LoginReply
}

func (r *Authentication__Login_LoginCheck) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/authentication.Login/LoginCheck"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Authentication__Login_LoginCheck) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 ** The route operation types for the monitor entries 
 */
const (
  /*
   ** 
   *  A new route is being added or modified.
   *  bgp_route will contain the route info. 
   */
  ROUTING__BGP_ROUTE_MONITOR_ENTRY__ROUTE__UPDATE = 0
  /*
   ** 
   *  An existing route is being removed.
   *  bgp_route will contain the route info. 
   */
  ROUTING__BGP_ROUTE_MONITOR_ENTRY__ROUTE__REMOVE = 1
  /*
   ** 
   *  The initial set of route monitoring entires following a fresh
   *  client registration has been completed. bgp_route will be null for 
   *  this operation.
   */
  ROUTING__BGP_ROUTE_MONITOR_ENTRY__END__OF__RIBS = 2
)

/*
 * Possible return codes for route service initialize operations. 
 */
const (
  /*
   ** Request successfully completed. Note that no preexisting 
   *  state for old clients with the same name was rebound. 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__SUCCESS = 0
  /*
   ** Request successfully completed AND preexisting routing state
   *  for an old client connection of the same name has been recovered
   *  and bound to this client connection. 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__SUCCESS__STATE__REBOUND = 1
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__INTERNAL__ERROR = 2
  /*
   ** Failed due to previous initialization operation. 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__ALREADY__INITIALIZED = 3
  /*
   ** Failed to find or create a gateway 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__GATEWAY__INVALID = 4
  /*
   ** Previous clean up work is pending try again later 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__CLEANUP__PENDING = 5
  /*
   ** The BGP protocol is not configured and initialized 
   */
  ROUTING__BGP_ROUTE_INITIALIZE_REPLY__BGP__NOT__READY = 6
)

/*
 ** Possible return codes for route service cleanup operations. 
 */
const (
  /*
   ** Request successfully completed. 
   */
  ROUTING__BGP_ROUTE_CLEANUP_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_CLEANUP_REPLY__INTERNAL__ERROR = 1
  /*
   ** Request failed because there was no initialized state to 
   *  cleanup. 
   */
  ROUTING__BGP_ROUTE_CLEANUP_REPLY__NOT__INITIALIZED = 2
)

/*
 ** Possible return codes for route add/modify/update/remove operations. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__INTERNAL__ERROR = 1
  /*
   ** The bgp route service has not been initialized 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__NOT__INITIALIZED = 2
  /*
   ** Request did not result in any operations 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__NO__OP = 3
  /*
   ** Request contained too many operations 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__TOO__MANY__OPS = 4
  /*
   ** Request contained an invalid table. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__TABLE__INVALID = 5
  /*
   ** Request contained a table that was not ready for operations. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__TABLE__NOT__READY = 6
  /*
   ** Request contained an invalid destination address prefix 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__PREFIX__INVALID = 7
  /*
   ** Request contained a destination prefix length too short for the
   *  supplied address/NLRI. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__PREFIX__LEN__TOO__SHORT = 8
  /*
   ** Request contained a destination prefix length too long for the 
   *  supplied address/NLRI. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__PREFIX__LEN__TOO__LONG = 9
  /*
   ** The server did not have a valid gateway associated with the 
   *  client. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__GATEWAY__INVALID = 10
  /*
   ** Request contained an invalid nexthop. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__NEXTHOP__INVALID = 11
  /*
   ** Request contained a nexthop with an invalild address. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__NEXTHOP__ADDRESS__INVALID = 12
  /*
   ** Request to add paths exceeding maximum ECMP paths for a 
   *  destination. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__NEXTHOP__ECMP__LIMIT = 13
  /*
   ** Request contained an invalid community. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__COMMUNITY__LIST__INVALID = 14
  /*
   ** Request contained an invalid AS path. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__ASPATH__INVALID = 15
  /*
   ** Request contained a invalid label information. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__LABEL__INFO__INVALID = 16
  /*
   ** Request contains a route that is already present in the table. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__ROUTE__EXISTS = 17
  /*
   ** Request contains a route that is NOT present in the table. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__ROUTE__NOT__FOUND = 18
  /*
   ** Request contains an invalid cluster list. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__CLUSTER__LIST__INVALID = 19
  /*
   ** Request contains an invalid protocol. Only PROTO_UNSPECIFID
   *  or PROTO_BGP_STATIC are allowed in route change operations.
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__PROTOCOL__INVALID = 20
  /*
   ** Request contains a route that is NOT present in the table. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__ROUTE__ADD__FAILED = 21
  /*
   ** The BGP protocol is not initialized and ready to accept
   *  route change operations. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__BGP__NOT__READY = 22
  /*
   ** Request cannot be serviced until current requests are processed. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__TRY__AGAIN = 23
  /*
   ** Request contains a parameter that is not currently supported. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__REQUEST__UNSUPPORTED = 24
  /*
   ** Request contained an invalid BGP peer type. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__PEER__TYPE__INVALID = 25
  /*
   ** SR-TE Segment Lists is invalid, like zero segment list 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__LIST__INVALID = 26
  /*
   ** SR-TE Segment is invalid, like zero segment list 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__INVALID = 27
  /*
   ** SR-TE Segment label is invalid; reserved label or label ttl > 255 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__LABEL__INVALID = 28
  /*
   ** SR-TE Segment ID is invalid like segment type is not set 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__ID__INVALID = 29
  /*
   ** Number of SR-TE Segment Lists exceeded limit (8) 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__LIST__COUNT__INVALID = 30
  /*
   ** Number of SR-TE Segments exceeded limit (5) 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SEGMENT__COUNT__INVALID = 31
  /*
   ** SRTE Route Data is not set. 
   */
  ROUTING__BGP_ROUTE_OPER_REPLY__SRTE__ROUTE__DATA__INVALID = 32
)

/*
 ** Possible return codes for route get operations. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__INTERNAL__ERROR = 1
  /*
   ** Request failed because there was no initialized state to
   *  cleanup. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__NOT__INITIALIZED = 2
  /*
   ** Request contained an invalid table. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__TABLE__INVALID = 3
  /*
   ** Request contained a table that was not ready for operations. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__TABLE__NOT__READY = 4
  /*
   ** Request contained an invalid destination address prefix 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__PREFIX__INVALID = 5
  /*
   ** Request contained a destination prefix length too short for the
   *  supplied address/NLRI. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__PREFIX__LEN__TOO__SHORT = 6
  /*
   ** Request contained a destination prefix length too long for the 
   *  supplied address/NLRI. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__PREFIX__LEN__TOO__LONG = 7
  /*
   ** Request contained a route that does not match 
   *  destinations in the routing table. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__ROUTE__NOT__FOUND = 8
  /*
   ** Request specified an invalid protocol to match 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__PROTOCOL__INVALID = 9
  /*
   ** Request does not contain valid route match parameters 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__ROUTE__INVALID = 10
  /*
   ** Request contains a parameter that is not currently supported. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__REQUEST__UNSUPPORTED = 11
  /*
   ** Request cannot be serviced until current requests are processed. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__TRY__AGAIN = 12
  /*
   ** Request contains a route_count that exceeds the max of 1000 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__ROUTE__COUNT__INVALID = 13
  /*
   ** Request contained an invalid community. 
   */
  ROUTING__BGP_ROUTE_GET_REPLY__COMMUNITY__LIST__INVALID = 14
)

/*
 ** Possible return codes for route get operations. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__INTERNAL__ERROR = 1
  /*
   ** The bgp route service has not been initialized 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__NOT__INITIALIZED = 2
  /*
   ** Too many clients or BMP stations are already registered. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__REGISTRATIONS__EXCEEDED = 3
  /*
   ** Client is already registered. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__ALREADY__REGISTERED = 4
  /*
   ** Request contains a route_count that exceeds the max of 1000 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__ROUTE__COUNT__INVALID = 5
  /*
   ** Request contains invalid address format. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__ADDRESS__FORMAT__INVALID = 6
  /*
   ** Request contains invalid table format. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REGISTER_REPLY__TABLE__FORMAT__INVALID = 7
)

/*
 ** Possible return codes for route unregister operation. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__BGP_ROUTE_MONITOR_UNREGISTER_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_MONITOR_UNREGISTER_REPLY__INTERNAL__ERROR = 1
  /*
   ** The bgp route service has not been initialized 
   */
  ROUTING__BGP_ROUTE_MONITOR_UNREGISTER_REPLY__NOT__INITIALIZED = 2
  /*
   ** Client is already registered. 
   */
  ROUTING__BGP_ROUTE_MONITOR_UNREGISTER_REPLY__NOT__REGISTERED = 3
)

/*
 ** Possible return codes for route unregister operation. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REFRESH_REPLY__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REFRESH_REPLY__INTERNAL__ERROR = 1
  /*
   ** The bgp route service has not been initialized 
   */
  ROUTING__BGP_ROUTE_MONITOR_REFRESH_REPLY__NOT__INITIALIZED = 2
  /*
   ** Client is already registered. 
   */
  ROUTING__BGP_ROUTE_MONITOR_REFRESH_REPLY__NOT__REGISTERED = 3
)

/*
 **
 * Routing protocols
 */
const (
  /*
   ** Unspecified protocol default behavior dependent on API:
   *  For route change requests, PROTO_BGP_STATIC.	
   *  For route get requests, match either PROTO_BGP or PROTO_BGP_STATIC.
   */
  ROUTING__PROTO__UNSPECIFIED = 0
  /*
   ** BGP dynamic routes 
   */
  ROUTING__PROTO__BGP = 1
  /*
   ** BGP static programmed routes 
   */
  ROUTING__PROTO__BGP__STATIC = 2
)

/*
 **
 * Route Operation Flag values that alter route add behavior.
 * USE OF THIS TYPE IS DEPRECATED. Instead, use BgpRouteOperationFlags.
 */
const (
  /*
   **
   * Unspecified Operation, meaning no special operation specified.
   * USE OF THIS FLAG IS DEPRECATED. Instead, use BgpRouteOperationFlags.
   */
  ROUTING__UNSPECIFIED = 0
  /*
   ** 
   *  Route operation indicating whether to attach the well-known
   *  no-advertise community. 
   *  No-advertise has the effect of instructing the route not to
   *  advertise the route further. The community may alternately be
   *  included in the community_list. 
   *  USE OF THIS FLAG IS DEPRECATED. Instead, use BgpRouteOperationFlags.
   */
  ROUTING__NO__ADVERTISE = 1
  /*
   ** 
   *  Route operation indicating whether to attach the well-known no-export 
   *  community. No-export has the effect of instructing the router 
   *  not to advertise the route beyond the BGP confederation boundary.
   *  The community may alternately be included in the community_list. 
   *  USE OF THIS FLAG IS DEPRECATED. Instead, use BgpRouteOperationFlags.
   */
  ROUTING__NO__EXPORT = 2
  /*
   ** 
   *  Route Operation indicating whether to use NH_REJECT for the route
   *  This makes sense to be set only when programming route in RR.
   *  This can be used to save memory when there are a high number of
   *  unique nexthops. 
   *  USE OF THIS FLAG IS DEPRECATED. Instead, use BgpRouteOperationFlags.
   */
  ROUTING__USE__NH__REJECT = 4
)

/*
 **
 * BGP peer type identifies if the associated route as an internal (IBGP) 
 * or external (EBGP) route.
 */
const (
  /*
   ** IBGP Route 
   */
  ROUTING__BGP__INTERNAL = 0
  /*
   ** EBGP Route 
   */
  ROUTING__BGP__EXTERNAL = 1
)

/*
 **
 * Route Operation Flag values that alter route add behavior.
 * This data type replaces and deprecates RouteOperation.
 * New applications should use BgpRouteOperationFlags exclusively and instead
 * of RouteOperation. 
 * Backwards compatibility: In BgpRouteUpdateRequest messages, if both 
 * bgp_routes[].route_oper_flags and bgp_routes[].route_flags are initialized,
 * then only the new route_flags will be processed and the old route_oper_flags
 * will be ignored. If only  route_oper_flags or route_flags are initialized,
 * then that initialized  flags message will be processed. 
 * In BgpRouteGetReply and BgpRouteMontiorEntry messages, both the
 * route_oper_flags and route_flags will be set within the bgp_routes[] and
 * bgp_route fields (respectively), indicating the same flags.
 */
type Routing__BgpRouteOperationFlags struct {
    NoAdvertise	 bool `json:"no_advertise,omitempty"`
    NoExport	 bool `json:"no_export,omitempty"`
    UseNexthopReject	 bool `json:"use_nexthop_reject,omitempty"`
}

/*
 **
 * A generic 32-bit unsigned value that implicitly carries an indication
 * of whether any value has been set or not.
 */
type Routing__BgpAttrib32 struct {
    Value	 uint32 `json:"value,omitempty"`
}

/*
 **
 * A generic 64-bit unsigned value that implicitly carries an indication
 * of whether any value has been set or not.
 */
type Routing__BgpAttrib64 struct {
    Value	 uint64 `json:"value,omitempty"`
}

/*
 **
 * A single communty is a string identifying a regular, extended, 
 * or well-known community name or values with no whitepace.
 * The communities will be recognized:
 * Well-known communities:
 *    no-export
 *    no-advertise
 *    no-export-confed
 *    llgr-stale
 *    no-llgr
 * RFC 1997 comunities:
 *    domain-id:ipaddress:0
 *    domain-id-vendor:
 *    <n>:<n>
 * 
 * Route targets extended communities:
 *    target:ipv4-address:16 bit#
 *    target:16bit#:32bit#
 *    target:as2b:16bit#:32bit#
 *    target:as4b:32bit#:16bit#
 * Origin extended communities:
 *    origin:ipv4-address:16 bit#
 *    origin:16bit#:32bit#
 * 
 * Bandwidth management extended communities:
 *    bandwidth:16bit#:bw
 *    {traffic-rate}:16 bit#:bw
 * Redirect extended communities: 
 *    redirect:ipv4-address:16 bit#
 *    redirect:16bit#:32bit#
 * Tunnel encapsulation extended communities:
 *    encapsulation:0L:tunnel-type
 * 
 */
type Routing__Community struct {
    CommunityString	 string `json:"community_string,omitempty"`
}

/*
 **
 * A list of communities.
 */
type Routing__CommunityList struct {
    ComList	 []Routing__Community `json:"com_list,omitempty"`
}

/*
 **
 * AS path through which the route was learned.
 * An AS Path is a string composed of an series of AS numbers 
 * separated by whitespace special delimiters.
 * The following special delimiting characters are used for signifying 
 * confederations, confederation-sets, and AS-sets:
 *     [ ] - Brackets enclose the local AS number associated with the AS set
 *     if more than one AS number is configured on the routing device, 
 *     or if AS path prepending is configured.
 *     ( ) - Parentheses enclose a confederation.
 *     ( [ ] ) - Parentheses and brackets enclose a confederation set.
 * 
 * Programmed AS Paths implicitly have path origin IGP.
 */
type Routing__AsPath struct {
    AspathString	 string `json:"aspath_string,omitempty"`
}

/*
 **
 * Route matching parameters provide the key for identifying
 * BGP routes. Programmed BGP-Static routes must be unique 
 * for the bgp_route_match paramaters. Dynamic BGP routes may
 * may have multiple matches to a given set of bgp_route_match
 * parameters.
 */
type Routing__BgpRouteMatch struct {
    DestPrefix	 *Routing__RoutePrefix `json:"dest_prefix,omitempty"`
    DestPrefixLen	 uint32 `json:"dest_prefix_len,omitempty"`
    Table	 *Routing__RouteTable `json:"table,omitempty"`
    Protocol	 uint32 `json:"protocol,omitempty"`
    PathCookie	 uint64 `json:"path_cookie,omitempty"`
    Communities	 *Routing__CommunityList `json:"communities,omitempty"`
}

/*
 **
 * A route entry specifying a single route destination and BGP path
 * along with the route attributes for that path.
 */
type Routing__BgpRouteEntry struct {
  /*
   * Only one of the following fields should be specified: 
   * VpnLabel, Labels
   */
    DestPrefix	 *Routing__RoutePrefix `json:"dest_prefix,omitempty"`
    DestPrefixLen	 uint32 `json:"dest_prefix_len,omitempty"`
    Table	 *Routing__RouteTable `json:"table,omitempty"`
    Protocol	 uint32 `json:"protocol,omitempty"`
    PathCookie	 uint64 `json:"path_cookie,omitempty"`
    RoutePreference	 *Routing__BgpAttrib32 `json:"route_preference,omitempty"`
    LocalPreference	 *Routing__BgpAttrib32 `json:"local_preference,omitempty"`
    Med	 *Routing__BgpAttrib32 `json:"med,omitempty"`
    AigpDistance	 *Routing__BgpAttrib64 `json:"aigp_distance,omitempty"`
    VpnLabel	 uint32 `json:"vpn_label,omitempty"`
    Labels	 *Routing__LabelStack `json:"labels,omitempty"`
    Communities	 *Routing__CommunityList `json:"communities,omitempty"`
    Aspath	 *Routing__AsPath `json:"aspath,omitempty"`
    OriginatorId	 *Routing__BgpAttrib32 `json:"originator_id,omitempty"`
    ClusterList	 []Routing__BgpAttrib32 `json:"cluster_list,omitempty"`
    ClusterId	 *Routing__BgpAttrib32 `json:"cluster_id,omitempty"`
    RouteOperFlag	 uint32 `json:"route_oper_flag,omitempty"`
    ProtocolNexthops	 []JnxBase__IpAddress `json:"protocol_nexthops,omitempty"`
    RouteType	 uint32 `json:"route_type,omitempty"`
    RouteFlags	 *Routing__BgpRouteOperationFlags `json:"route_flags,omitempty"`
    RouteData	 *Routing__AddressFamilySpecificData `json:"route_data,omitempty"`
}

/*
 **
 * A monitoring entry for a single BGP route streamed from BGP when the 
 * client has registered for route monitoring with BgpRouteMonitorRegister().
 */
type Routing__BgpRouteMonitorEntry struct {
    Operation	 uint32 `json:"operation,omitempty"`
    BgpRoute	 *Routing__BgpRouteEntry `json:"bgp_route,omitempty"`
}

/*
 ** 
 * Request to initialize the BGP route service. No parameters are needed.
 */
type Routing__BgpRouteInitializeRequest struct {
}

/*
 **
 * BGP route service initialize reply containing the status of the operation.
 * Replies indicate to the client whether any old routing state was 
 * recovered and rebound to this connection.
 */
type Routing__BgpRouteInitializeReply struct {
    Status	 uint32 `json:"status,omitempty"`
    GwNRoutes	 uint32 `json:"gw_n_routes,omitempty"`
}

/*
 ** 
 * Request to reset the BGP route service. 
 * Any routes that were added by the client will be removed during the
 * cleanup of the client's state. No parameters are needed.
 */
type Routing__BgpRouteCleanupRequest struct {
}

/*
 **
 * Route service cleanup reply containing the status of the operation.
 */
type Routing__BgpRouteCleanupReply struct {
    Status	 uint32 `json:"status,omitempty"`
}

/*
 **
 * Route add, modify, or update operation request parameters. 
 */
type Routing__BgpRouteUpdateRequest struct {
    BgpRoutes	 []Routing__BgpRouteEntry `json:"bgp_routes,omitempty"`
}

/*
 **
 * Route remove operation request parameters. 
 */
type Routing__BgpRouteRemoveRequest struct {
    OrLonger	 bool `json:"or_longer,omitempty"`
    BgpRoutes	 []Routing__BgpRouteMatch `json:"bgp_routes,omitempty"`
}

/*
 **
 * Route get operation request parameters. 
 */
type Routing__BgpRouteGetRequest struct {
    BgpRoute	 *Routing__BgpRouteMatch `json:"bgp_route,omitempty"`
    OrLonger	 bool `json:"or_longer,omitempty"`
    ActiveOnly	 bool `json:"active_only,omitempty"`
    ReplyAddressFormat	 uint32 `json:"reply_address_format,omitempty"`
    ReplyTableFormat	 uint32 `json:"reply_table_format,omitempty"`
    RouteCount	 uint32 `json:"route_count,omitempty"`
}

/*
 **
 * Route operation reply containing the status of the operation.
 * Replies always returns the final status (either success or the first error
 * encountered) and the number of routes that were successfully processed
 * prior to any error or full completion of the request.
 */
type Routing__BgpRouteOperReply struct {
    Status	 uint32 `json:"status,omitempty"`
    OperationsCompleted	 uint32 `json:"operations_completed,omitempty"`
}

/*
 **
 * Route get reply containing the status of the operation and the full or 
 * partial set of matching routes, depending on how many reply RPCs the
 * stream of routes is split among.
 */
type Routing__BgpRouteGetReply struct {
    Status	 uint32 `json:"status,omitempty"`
    BgpRoutes	 []Routing__BgpRouteEntry `json:"bgp_routes,omitempty"`
}

/*
 **
 * Route register operation request parameters. Registers the client
 * for streaming route monitoring.
 */
type Routing__BgpRouteMonitorRegisterRequest struct {
    ReplyAddressFormat	 uint32 `json:"reply_address_format,omitempty"`
    ReplyTableFormat	 uint32 `json:"reply_table_format,omitempty"`
    RouteCount	 uint32 `json:"route_count,omitempty"`
}

/*
 **
 * The route registration reply is returned immediately upon initial 
 * registration for route monitoring via a call to BgpRouteMonitorRegister.
 * Subsequently, monitoring entries are streamed via replies containing 
 * route information and status.
 */
type Routing__BgpRouteMonitorRegisterReply struct {
    Status	 uint32 `json:"status,omitempty"`
    MonitorEntries	 []Routing__BgpRouteMonitorEntry `json:"monitor_entries,omitempty"`
}

/*
 ** 
 * Request to unregister the client from the BGP route monitoring.
 * No parameters are needed.
 */
type Routing__BgpRouteMonitorUnregisterRequest struct {
}

/*
 **
 * The route unregistration reply confirms that the client has
 * unregistered for route updates.
 */
type Routing__BgpRouteMonitorUnregisterReply struct {
    Status	 uint32 `json:"status,omitempty"`
}

/*
 ** 
 * Request to refresh all route monitoring entries to the client.
 * No parameters are needed.
 */
type Routing__BgpRouteMonitorRefreshRequest struct {
}

/*
 **
 * The route Refresh reply confirms that the client has
 * triggered a refresh of route monitoring entries, which
 * will be delivered followed by End-of-RIBs via the 
 * BgpRouteMonitorRegisterReply stream.
 */
type Routing__BgpRouteMonitorRefreshReply struct {
    Status	 uint32 `json:"status,omitempty"`
}

type Routing__BgpRoute_BgpRouteInitialize struct {
  Request	Routing__BgpRouteInitializeRequest
  Reply	Routing__BgpRouteInitializeReply
}

func (r *Routing__BgpRoute_BgpRouteInitialize) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteInitialize"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteInitialize) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteCleanup struct {
  Request	Routing__BgpRouteCleanupRequest
  Reply	Routing__BgpRouteCleanupReply
}

func (r *Routing__BgpRoute_BgpRouteCleanup) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteCleanup"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteCleanup) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteAdd struct {
  Request	Routing__BgpRouteUpdateRequest
  Reply	Routing__BgpRouteOperReply
}

func (r *Routing__BgpRoute_BgpRouteAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteModify struct {
  Request	Routing__BgpRouteUpdateRequest
  Reply	Routing__BgpRouteOperReply
}

func (r *Routing__BgpRoute_BgpRouteModify) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteModify"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteModify) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteUpdate struct {
  Request	Routing__BgpRouteUpdateRequest
  Reply	Routing__BgpRouteOperReply
}

func (r *Routing__BgpRoute_BgpRouteUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteRemove struct {
  Request	Routing__BgpRouteRemoveRequest
  Reply	Routing__BgpRouteOperReply
}

func (r *Routing__BgpRoute_BgpRouteRemove) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteRemove"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteRemove) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteGet struct {
  Request	Routing__BgpRouteGetRequest
  Reply	Routing__BgpRouteGetReply
}

func (r *Routing__BgpRoute_BgpRouteGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteMonitorRegister struct {
  Request	Routing__BgpRouteMonitorRegisterRequest
  Reply	Routing__BgpRouteMonitorRegisterReply
}

func (r *Routing__BgpRoute_BgpRouteMonitorRegister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteMonitorRegister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteMonitorRegister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteMonitorUnregister struct {
  Request	Routing__BgpRouteMonitorUnregisterRequest
  Reply	Routing__BgpRouteMonitorUnregisterReply
}

func (r *Routing__BgpRoute_BgpRouteMonitorUnregister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteMonitorUnregister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteMonitorUnregister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__BgpRoute_BgpRouteMonitorRefresh struct {
  Request	Routing__BgpRouteMonitorRefreshRequest
  Reply	Routing__BgpRouteMonitorRefreshReply
}

func (r *Routing__BgpRoute_BgpRouteMonitorRefresh) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.BgpRoute/BgpRouteMonitorRefresh"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__BgpRoute_BgpRouteMonitorRefresh) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * Review notes:  I have removed the MAX/ definitions etc., from IPC 
 */
const (
  CLKSYNC__INVALID = 0
  CLKSYNC__CHASSIS__FREQ = 1
  CLKSYNC__CHASSIS__PHASE = 2
  CLKSYNC__LINE__PRI = 3
  CLKSYNC__LINE__SEC = 4
  CLKSYNC__GPS__FREQ = 5
  CLKSYNC__GPS__PHASE = 6
  CLKSYNC__GPS__FREQ__N__PHASE = 7
  CLKSYNC__BITS__FREQ = 8
  CLKSYNC__PTP__FREQ__N__PHASE = 9
  CLKSYNC__FREERUN = 10
  CLKSYNC__HYBRID__FREQ = 11
  CLKSYNC__HYBRID__PHASE = 12
  CLKSYNC__HYBRID__FREQ__N__PHASE = 13
)

const (
  CLKSYNC____INVALID = 0
  CLKSYNC____CHASSIS__FREQ = 1
  CLKSYNC____CHASSIS__PHASE = 2
  CLKSYNC____LINE__PRI = 3
  CLKSYNC____LINE__SEC = 4
  CLKSYNC____GPS__FREQ = 5
  CLKSYNC____GPS__PHASE = 6
  CLKSYNC____BITS__FREQ = 7
  CLKSYNC____PTP__FREQ = 8
  CLKSYNC____PTP__PHASE = 9
)

const (
  CLKSYNC__UNKNOWN__SUBMSG = 0
  CLKSYNC__PID = 1
  CLKSYNC__VERSION = 2
)

const (
  CLKSYNC__UNKNOWN = 0
  CLKSYNC__CHASSIS = 1
  CLKSYNC__SETS = 2
  CLKSYNC__PTP = 3
  CLKSYNC__GPS = 4
)

const (
  CLKSYNC____NONE = 0
  CLKSYNC____SYNCE = 1
  CLKSYNC____PTP = 2
  CLKSYNC____HYBRID = 3
)

const (
  CLKSYNC__PRI = 0
  CLKSYNC__SEC = 1
  CLKSYNC__TERT = 2
  CLKSYNC__QUART = 3
  CLKSYNC__QUIN = 4
  CLKSYNC__INVALID__IDX = 5
)

const (
  CLKSYNC__DISABLED = 0
  CLKSYNC__QUALIFYING = 1
  CLKSYNC__FAILED = 2
  CLKSYNC__QUALIFIED = 3
)

const (
  CLKSYNC____DISABLED = 0
  CLKSYNC____QUALIFYING = 1
  CLKSYNC____FAILED = 2
  CLKSYNC____QUALIFIED = 3
)

const (
  CLKSYNC__PLL__INVALID = 0
  CLKSYNC__PLL__P1_HZ = 1
  CLKSYNC__PLL_1_P7_HZ = 2
  CLKSYNC__PLL_3_P5_HZ = 3
  CLKSYNC__PLL_14_HZ = 4
  CLKSYNC__PLL_28_HZ = 5
  CLKSYNC__PLL_890_HZ = 6
  CLKSYNC__PLL__FAST__LOCK = 7
)

const (
  CLKSYNC__P885_US = 0
  CLKSYNC___7_P5_US = 1
  CLKSYNC___61_US = 2
  CLKSYNC__UNLIMITED = 3
)

const (
  CLKSYNC___1_PPS = 0
  CLKSYNC___8__KHZ = 1
  CLKSYNC___1544__KHZ = 2
  CLKSYNC___2048__KHZ = 3
  CLKSYNC____T1 = 4
  CLKSYNC____E1 = 5
  CLKSYNC___1__MHZ = 6
  CLKSYNC___5__MHZ = 7
  CLKSYNC___10__MHZ = 8
  CLKSYNC___1944_MZ = 9
  CLKSYNC____UNKNOWN = 10
)

const (
  CLKSYNC__UNKNOWN__FRM = 0
  CLKSYNC__ESF = 1
  CLKSYNC__G704 = 2
  CLKSYNC__G704__NO__CRC4 = 3
  CLKSYNC__SF = 4
  CLKSYNC__D4 = 5
  CLKSYNC__FAS = 6
  CLKSYNC__CAS = 7
  CLKSYNC__CRC4 = 8
  CLKSYNC__CAS__CRC4 = 9
)

const (
  CLKSYNC__UNKNOWN__ENC = 0
  CLKSYNC__AMI = 1
  CLKSYNC__B8_ZS = 2
  CLKSYNC__HDB3 = 3
)

const (
  CLKSYNC__UNKNOWN__SA = 0
  CLKSYNC__NA = 1
  CLKSYNC__SA4 = 4
  CLKSYNC__SA5 = 5
  CLKSYNC__SA6 = 6
  CLKSYNC__SA7 = 7
  CLKSYNC__SA8 = 8
)

const (
  CLKSYNC_____DISABLED = 0
  CLKSYNC__SQUELCHED = 1
  CLKSYNC__ACTIVE = 2
)

const (
  CLKSYNC__BITS__CLK = 0
  CLKSYNC__GPS__CLK = 1
  CLKSYNC__DTI__CLK = 2
)

const (
  CLKSYNC__MODE__UNKNOWN = 0
  CLKSYNC__MODE__NORMAL = 1
  CLKSYNC__MODE__HOLDOVER = 2
  CLKSYNC__MODE__FREERUN = 3
  CLKSYNC__MODE__AUTO = 4
  CLKSYNC__MODE__TOP = 5
)

const (
  CLKSYNC__UNKNOWN__STATE = 0
  CLKSYNC__INIT__STATE = 1
  CLKSYNC__LOCK__ACQ__STATE = 2
  CLKSYNC__LOCKED__STATE = 3
  CLKSYNC__HOLDOVER__STATE = 4
  CLKSYNC__FREERUN__STATE = 5
)

const (
  CLKSYNC__QL__UNDEFINED = 0
  CLKSYNC__QL__PRC = 1
  CLKSYNC__QL__SSU__A = 2
  CLKSYNC__QL__SSU__B = 3
  CLKSYNC__QL__SEC = 4
  CLKSYNC__QL__PRS = 5
  CLKSYNC__QL__STU = 6
  CLKSYNC__QL__ST2 = 7
  CLKSYNC__QL__TNC = 8
  CLKSYNC__QL__ST3_E = 9
  CLKSYNC__QL__ST3 = 10
  CLKSYNC__QL__SMC = 11
  CLKSYNC__QL__ST4 = 12
)

const (
  CLKSYNC__LINK__STATE__UP = 0
  CLKSYNC__LINK__STATE__UP__TO__DOWN = 1
  CLKSYNC__LINK__STATE__DOWN__TO__UP = 2
  CLKSYNC__LINK__STATE__DOWN = 3
)

const (
  CLKSYNC__INVALID__REQ = 0
  CLKSYNC__ENABLE__REQ = 1
  CLKSYNC__DISABLE__REQ = 2
)

const (
  CLKSYNC__EXT__NONE = 0
  CLKSYNC__EXT__GREEN = 1
  CLKSYNC__EXT__RED = 2
)

/*
 * ----------------------------------------------------------------------------
 * Different types of return codes to be sent back to client based on the
 * operation was successful or not and if not, possibly more specific reasons
 * as to why it failed.
 * ----------------------------------------------------------------------------
 */
const (
  /*
   * Operation was executed successfully
   */
  CLKSYNC__RET__SUCCESS = 0
  /*
   * General failure : operation not executed successfully
   */
  CLKSYNC__RET__FAILURE = 1
  /*
   * Entry was not found
   */
  CLKSYNC__RET__NOT__FOUND = 2
  /*
   * Invalid input paramters
   */
  CLKSYNC__RET__INVALID__PARAMS = 3
)

/*
 * To check if clksync_pi_init_clock_modules_msg_ to merge with common init message 
 */
type Clksync__ClksyncHwInitRequest struct {
    ModulesBitmap	 uint32 `json:"modules_bitmap,omitempty"`
}

type Clksync__ClksyncOperReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
}

type Clksync__ClksyncTodTime struct {
    UtcSecs	 uint32 `json:"utc_secs,omitempty"`
    UtcNsecs	 uint32 `json:"utc_nsecs,omitempty"`
}

type Clksync__ClksyncGetTodReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    UtcSecs	 uint32 `json:"utc_secs,omitempty"`
    UtcNsecs	 uint32 `json:"utc_nsecs,omitempty"`
    ParserStatus	 int32 `json:"parser_status,omitempty"`
    UtcOffset	 uint32 `json:"utc_offset,omitempty"`
}

type Clksync__ClksyncPushInfoUpMsg struct {
    SubmsgType	 uint32 `json:"submsg_type,omitempty"`
    Param1	 uint32 `json:"param1,omitempty"`
    Param2	 uint32 `json:"param2,omitempty"`
}

type Clksync__ClksyncPllFreqStatusMsg struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    FreqStatus	 uint32 `json:"freq_status,omitempty"`
}

type Clksync__ClksyncPllSyncStatusMsg struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    PhaseStatus	 uint32 `json:"phase_status,omitempty"`
}

type Clksync__ClksyncPushInfoDownRequest struct {
    SubmsgType	 uint32 `json:"submsg_type,omitempty"`
    Param1	 uint32 `json:"param1,omitempty"`
    Param2	 uint32 `json:"param2,omitempty"`
}

type Clksync__ClkyncdSignalForClkSyncRequest struct {
    EnableDisable	 uint32 `json:"enable_disable,omitempty"`
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    InputSignalType	 uint32 `json:"input_signal_type,omitempty"`
    ClkSignalUnitNum	 uint32 `json:"clk_signal_unit_num,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    ClkSrcIdx	 uint32 `json:"clk_src_idx,omitempty"`
    ClkSysMode	 uint32 `json:"clk_sys_mode,omitempty"`
}

type Clksync__ClksyncPllSrcInputRequest struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    ClkPllBw	 uint32 `json:"clk_pll_bw,omitempty"`
    ClkPslBw	 uint32 `json:"clk_psl_bw,omitempty"`
    ClkSrcIdx	 uint32 `json:"clk_src_idx,omitempty"`
    ClkSysMode	 uint32 `json:"clk_sys_mode,omitempty"`
}

/*
 *Mapped to clksync_pi_enable_dpll_clk_src_input_msg_ 
 */
type Clksync__ClksyncCfgPllSrcInputRequest struct {
    EnableDisable	 uint32 `json:"enable_disable,omitempty"`
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    ClkSignalFreq	 uint32 `json:"clk_signal_freq,omitempty"`
    DisableInputClkLine	 uint32 `json:"disable_input_clk_line,omitempty"`
    ClkSrcIdx	 uint32 `json:"clk_src_idx,omitempty"`
    ClkSysMode	 uint32 `json:"clk_sys_mode,omitempty"`
}

/*
 * Maps to clksync_pi_ext_config_interface_msg_ 
 */
type Clksync__ClksyncCfgExtIntfRequest struct {
    FramerReset	 uint32 `json:"framer_reset,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    ClkSignalFreq	 uint32 `json:"clk_signal_freq,omitempty"`
    FramingMode	 uint32 `json:"framing_mode,omitempty"`
    ClkLineEncoding	 uint32 `json:"clk_line_encoding,omitempty"`
    ClkSabit	 uint32 `json:"clk_sabit,omitempty"`
    IntfCfgChanged	 bool `json:"intf_cfg_changed,omitempty"`
}

/*
 *Maps to clksync_pi_select_clk_src_for_clk_out_msg_ 
 */
type Clksync__ClksyncSelectExtClkOutputRequest struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkOutputUnitNum	 uint32 `json:"clk_output_unit_num,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    ClkSignalFreq	 uint32 `json:"clk_signal_freq,omitempty"`
    ExtOutStatus	 uint32 `json:"ext_out_status,omitempty"`
}

/*
 * Maps to clksync_pi_ext_disable_clk_output_msg_ & clksync_pi_ext_enable_clk_output_msg_ 
 */
type Clksync__ClksyncCfgExtClkOutputRequest struct {
    EnableDisable	 uint32 `json:"enable_disable,omitempty"`
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkOutputUnitNum	 uint32 `json:"clk_output_unit_num,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    DisableInputClkSrc	 uint32 `json:"disable_input_clk_src,omitempty"`
    ExtOutStatus	 uint32 `json:"ext_out_status,omitempty"`
}

/*
 * Maps to clksync_pi_ext_set_tx_ql_msg_ 
 */
type Clksync__ClksyncSetExtTxQlRequest struct {
    ClkUnitNum	 uint32 `json:"clk_unit_num,omitempty"`
    SsmTxQl	 uint32 `json:"ssm_tx_ql,omitempty"`
}

/*
 * Maps to clksync_pi_ext_set_led_color_msg_t 
 */
type Clksync__ClksyncSetExtLedRequest struct {
    ClkUnitNum	 uint32 `json:"clk_unit_num,omitempty"`
    ExtLedType	 uint32 `json:"ext_led_type,omitempty"`
    LedRxColor	 uint32 `json:"led_rx_color,omitempty"`
    LedTxColor	 uint32 `json:"led_tx_color,omitempty"`
}

/*
 * Maps to clksync_pi_update_dpll_status_change_msg_ 
 */
type Clksync__ClksyncPllStatusMsg struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkSrc	 uint32 `json:"clk_src,omitempty"`
    ClkSrcUnitNum	 uint32 `json:"clk_src_unit_num,omitempty"`
    DpllState	 uint32 `json:"dpll_state,omitempty"`
    HoldoverStatus	 uint32 `json:"holdover_status,omitempty"`
}

/*
 *Maps to clksync_pi_update_ext_ssm_ql_change_msg_ 
 */
type Clksync__ClksyncExtIntfStatusMsg struct {
    ClkUnitNum	 uint32 `json:"clk_unit_num,omitempty"`
    LinkState	 uint32 `json:"link_state,omitempty"`
}

/*
 *Maps to clksync_pi_update_ext_ssm_ql_change_msg_ 
 */
type Clksync__ClksyncExtRxQlMsg struct {
    ClkUnitNum	 uint32 `json:"clk_unit_num,omitempty"`
    QlValue	 uint32 `json:"ql_value,omitempty"`
}

/*
 *Maps to clksync_pi_dpll_set_mode_msg_ 
 */
type Clksync__ClksyncSetPllModeRequest struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    DpllMode	 uint32 `json:"dpll_mode,omitempty"`
}

/*
 *Maps to clksync_pi_dpll_set_bw_psl_msg_  
 */
type Clksync__ClksyncSetPllBwPslRequest struct {
    DpllType	 uint32 `json:"dpll_type,omitempty"`
    ClkPllBw	 uint32 `json:"clk_pll_bw,omitempty"`
    ClkPslBw	 uint32 `json:"clk_psl_bw,omitempty"`
}

type Clksync__ClksyncSetTodConfigRequest struct {
    TodPortName	 string `json:"tod_port_name,omitempty"`
    TodCfgFormat	 string `json:"tod_cfg_format,omitempty"`
}

type Clksync__ClksyncConnectMsg struct {
    ConnectType	 uint32 `json:"connect_type,omitempty"`
}

type Clksync__ClksyncStatusGetReq struct {
    Update	 bool `json:"update,omitempty"`
}

type Clksync__ClksyncExtStatusGetReq struct {
    Update	 bool `json:"update,omitempty"`
}

type Clksync__ClksyncExtRxQlGetReq struct {
    Update	 bool `json:"update,omitempty"`
}

type Clksync__ClksyncTodStartMsg struct {
    UtcSec	 uint32 `json:"utc_sec,omitempty"`
    UtcNsec	 uint32 `json:"utc_nsec,omitempty"`
    StartEpats	 uint32 `json:"start_epats,omitempty"`
}

type Clksync__Timing_ClksyncConnect struct {
  Request	Clksync__ClksyncConnectMsg
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncConnect) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncConnect"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncConnect) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncCfgPushInfoDownRequest struct {
  Request	Clksync__ClksyncPushInfoDownRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncCfgPushInfoDownRequest) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncCfgPushInfoDownRequest"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncCfgPushInfoDownRequest) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncCfgSignalForClkSrc struct {
  Request	Clksync__ClkyncdSignalForClkSyncRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncCfgSignalForClkSrc) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncCfgSignalForClkSrc"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncCfgSignalForClkSrc) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncInitHw struct {
  Request	Clksync__ClksyncHwInitRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncInitHw) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncInitHw"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncInitHw) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSelectPllSrcInput struct {
  Request	Clksync__ClksyncPllSrcInputRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSelectPllSrcInput) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSelectPllSrcInput"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSelectPllSrcInput) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncCfgPllSrcInput struct {
  Request	Clksync__ClksyncCfgPllSrcInputRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncCfgPllSrcInput) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncCfgPllSrcInput"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncCfgPllSrcInput) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncCfgExtIntf struct {
  Request	Clksync__ClksyncCfgExtIntfRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncCfgExtIntf) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncCfgExtIntf"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncCfgExtIntf) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSelectExtClkOutput struct {
  Request	Clksync__ClksyncSelectExtClkOutputRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSelectExtClkOutput) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSelectExtClkOutput"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSelectExtClkOutput) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncCfgExtClkOutput struct {
  Request	Clksync__ClksyncCfgExtClkOutputRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncCfgExtClkOutput) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncCfgExtClkOutput"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncCfgExtClkOutput) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSetExtTxQl struct {
  Request	Clksync__ClksyncSetExtTxQlRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSetExtTxQl) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSetExtTxQl"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSetExtTxQl) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSetExtLed struct {
  Request	Clksync__ClksyncSetExtLedRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSetExtLed) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSetExtLed"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSetExtLed) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSetPllMode struct {
  Request	Clksync__ClksyncSetPllModeRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSetPllMode) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSetPllMode"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSetPllMode) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSetPllBwPsl struct {
  Request	Clksync__ClksyncSetPllBwPslRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSetPllBwPsl) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSetPllBwPsl"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSetPllBwPsl) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncSetTodConfig struct {
  Request	Clksync__ClksyncSetTodConfigRequest
  Reply	Clksync__ClksyncOperReply
}

func (r *Clksync__Timing_ClksyncSetTodConfig) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncSetTodConfig"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncSetTodConfig) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncGetTodValue struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncGetTodReply
}

func (r *Clksync__Timing_ClksyncGetTodValue) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncGetTodValue"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncGetTodValue) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncPushInfoUp struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncPushInfoUpMsg
}

func (r *Clksync__Timing_ClksyncPushInfoUp) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncPushInfoUp"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncPushInfoUp) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncPllFreqStatusGet struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncPllFreqStatusMsg
}

func (r *Clksync__Timing_ClksyncPllFreqStatusGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncPllFreqStatusGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncPllFreqStatusGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncPllStatusGet struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncPllStatusMsg
}

func (r *Clksync__Timing_ClksyncPllStatusGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncPllStatusGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncPllStatusGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncPllSyncStatusGet struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncPllSyncStatusMsg
}

func (r *Clksync__Timing_ClksyncPllSyncStatusGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncPllSyncStatusGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncPllSyncStatusGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncExtIntfStatusGet struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncExtIntfStatusMsg
}

func (r *Clksync__Timing_ClksyncExtIntfStatusGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncExtIntfStatusGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncExtIntfStatusGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncExtRxQlGet struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncExtRxQlMsg
}

func (r *Clksync__Timing_ClksyncExtRxQlGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncExtRxQlGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncExtRxQlGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Clksync__Timing_ClksyncTodStartReq struct {
  Request	Clksync__ClksyncStatusGetReq
  Reply	Clksync__ClksyncTodStartMsg
}

func (r *Clksync__Timing_ClksyncTodStartReq) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/clksync.Timing/ClksyncTodStartReq"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Clksync__Timing_ClksyncTodStartReq) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



const (
  COS__ATTRIB__NOOP = 0
  COS__ATTRIB__ADD = 1
  COS__ATTRIB__CHANGE = 2
  COS__ATTRIB__DELETE = 3
)

const (
  COS__CLFR__RW__LP__INVALID = 0
  COS__CLFR__RW__LP__HIGH = 1
  COS__CLFR__RW__LP__MEDIUM__HIGH = 2
  COS__CLFR__RW__LP__MEDIUM__LOW = 3
  COS__CLFR__RW__LP__LOW = 4
)

const (
  COS__SCHED__PRI__INVALID = 0
  COS__SCHED__PRI__HIGH = 1
  COS__SCHED__PRI__MEDIUM__HIGH = 2
  COS__SCHED__PRI__MEDIUM__LOW = 3
  COS__SCHED__PRI__LOW = 4
  COS__SCHED__PRI__STRICT__HIGH = 5
)

const (
  COS__SCHED__EXCESS__PRI__INVALID = 0
  COS__SCHED__EXCESS__PRI__HIGH = 1
  COS__SCHED__EXCESS__PRI__MEDIUM__HIGH = 2
  COS__SCHED__EXCESS__PRI__MEDIUM__LOW = 3
  COS__SCHED__EXCESS__PRI__LOW = 4
  COS__SCHED__EXCESS__PRI__NONE = 5
)

const (
  COS__SCHED__BUF__SIZE__SPEC__INVALID = 0
  COS__SCHED__BUF__SIZE__SPEC__REMAINDER = 1
  /*
   * percent 0 - 100
   */
  COS__SCHED__BUF__SIZE__SPEC__PERCENT = 2
  /*
   * microseconds
   */
  COS__SCHED__BUF__SIZE__SPEC__TEMPORAL = 3
)

const (
  COS__EXCESS__RATE__SPEC__INVALID = 0
  /*
   * range 1 - 100
   */
  COS__EXCESS__RATE__SPEC__PERCENT = 1
  /*
   * range 0 - 1000
   */
  COS__EXCESS__RATE__SPEC__PROPORTION = 2
)

const (
  COS__TCP__RATE__SPEC__INVALID = 0
  /*
   * 1000..6400000000000 bits per second
   */
  COS__TCP__RATE__SPEC__ABSOLUTE = 1
  /*
   * percent 0 - 100
   * not valid for per priority shaping rates
   * not valid for adjust-minimum
   */
  COS__TCP__RATE__SPEC__PERCENT = 2
)

const (
  COS__SCHED__RATE__SPEC__INVALID = 0
  /*
   * 3200 - 6400000000000 bits per second
   */
  COS__SCHED__RATE__SPEC__ABSOLUTE = 1
  /*
   * percent range 0-100 for guaranteed rates
   * percent range 1-100 for shaping rates
   */
  COS__SCHED__RATE__SPEC__PERCENT = 2
  /*
   * not valid for shaping rate.
   */
  COS__SCHED__RATE__SPEC__REMAINDER = 3
)

const (
  COS__SCHED__G__RATE__EXTEND__SPEC__INVALID = 0
  COS__SCHED__G__RATE__EXTEND__SPEC__EXACT = 1
  COS__SCHED__G__RATE__EXTEND__SPEC__RATE__LIMIT = 2
)

const (
  COS__DELAY__BUF__RATE__SPEC__INVALID = 0
  COS__DELAY__BUF__RATE__SPEC__ABSOLUTE = 1
  COS__DELAY__BUF__RATE__SPEC__CELL__PER__SECOND = 2
  COS__DELAY__BUF__RATE__SPEC__PERCENT = 3
)

const (
  COS__OH__ACCT__MODE__INVALID = 0
  COS__OH__ACCT__MODE__FRAME = 1
  COS__OH__ACCT__MODE__CELL = 2
)

const (
  COS__SCHED__DP__LP__INVALID = 0
  COS__SCHED__DP__LP__HIGH = 1
  COS__SCHED__DP__LP__MEDIUM__HIGH = 2
  COS__SCHED__DP__LP__MEDIUM__LOW = 3
  COS__SCHED__DP__LP__LOW = 4
  COS__SCHED__DP__LP__ANY = 5
)

const (
  COS__SCHED__DP__PROTO__INVALID = 0
  COS__SCHED__DP__PROTO__TCP = 1
  COS__SCHED__DP__PROTO__NON__TCP = 2
  COS__SCHED__DP__PROTO__ANY = 3
)

const (
  COS__DP__TYPE__INVALID = 0
  /*
   * discrete fill levels
   */
  COS__DP__TYPE__DISCRETE = 1
  COS__DP__TYPE__INTERPOLATE = 2
)

const (
  COS__FEATURE__CP__TYPE__INVALID = 0
  COS__FEATURE__CP__TYPE__DSCP = 1
  COS__FEATURE__CP__TYPE__DSCP__IPV6 = 2
  COS__FEATURE__CP__TYPE__EXP = 3
  COS__FEATURE__CP__TYPE__IEEE8021 = 4
  COS__FEATURE__CP__TYPE__IEEE8021_AD = 5
  COS__FEATURE__CP__TYPE__INET__PRECEDENCE = 6
)

const (
  COS__EXCESS__BW__SHARE__SPEC__INVALID = 0
  COS__EXCESS__BW__SHARE__SPEC__EQUAL = 1
  COS__EXCESS__BW__SHARE__SPEC__PROPORTION = 2
)

const (
  /*
   * Success
   */
  COS__EOK = 0
  /*
   * Invalid status
   */
  COS__INVALID = 1
  /*
   * The RPC was a NULL buffer
   */
  COS__NULL__MESSAGE = 2
  /*
   * Wrong input
   */
  COS__EINVALID__MESSAGE = 3
  /*
   * Server Internal error
   */
  COS__EINTERNAL = 4
  /*
   * Operation not supported
   */
  COS__UNSUPPORTED__OP = 5
  /*
   * Resource not available at server
   */
  COS__NO__RESOURCE = 6
  /*
   * No entry available at server
   */
  COS__NO__ENTRY = 7
  /*
   * No more entries available at server
   */
  COS__NO__MORE__ENTRIES = 8
  /*
   * Object is in use by other dependants
   */
  COS__E__IN__USE = 9
  /*
   * unsupported feature
   */
  COS__E__UNSUPP__FEATURE = 10
)

const (
  COS__ATM__SERVICE__INVALID = 0
  COS__CBR = 1
  COS__NRT__VBR = 2
  COS__RT__VBR = 3
)

const (
  COS__NODE__FEATURE__INVALID = 0
  COS__NODE__FEATURE__SMAP = 1
  /*
   * not valid for interface sets
   */
  COS__NODE__FEATURE__SRATE = 2
  COS__NODE__FEATURE__TCP = 3
  COS__NODE__FEATURE__TCP__REMAINING = 4
  /*
   * not valid for interface sets
   */
  COS__NODE__FEATURE__FORWARDING__CLASS = 5
  COS__NODE__FEATURE__MEMBER__LINK__SCHED = 6
  COS__NODE__FEATURE__EXCESS__BW__SHARE = 7
)

const (
  COS__DIRECTION__INVALID = 0
  COS__INGRESS = 1
  COS__EGRESS = 2
)

const (
  COS__NODE__TYPE__INVALID = 0
  COS__NODE__TYPE__INTERFACE = 1
  COS__NODE__TYPE__LOGICAL__INTERFACE = 2
  COS__NODE__TYPE__INTERFACE__SET = 3
)

const (
  COS__NODE__FEATURE__FAMILY__INVALID = 0
  COS__NODE__FEATURE__FAMILY__CLFR__RULE = 1
  COS__NODE__FEATURE__FAMILY__RW__RULE = 2
)

/*
 **
 * Classifier family
 */
const (
  COS__NODE__FEATURE__CLFR__PROTO__INVALID = 0
  COS__NODE__FEATURE__CLFR__DSCP = 1
  COS__NODE__FEATURE__CLFR__DSCP__MPLS = 2
  COS__NODE__FEATURE__CLFR__DSCP__IPV6 = 3
  COS__NODE__FEATURE__CLFR__DSCP__IPV6__MPLS = 4
  COS__NODE__FEATURE__CLFR__EXP = 5
  COS__NODE__FEATURE__CLFR__INET__PRECEDENCE = 6
  COS__NODE__FEATURE__CLFR__IEEE8021_P = 7
  COS__NODE__FEATURE__CLFR__IEEE8021_P__TAG__MODE__INNER = 8
  COS__NODE__FEATURE__CLFR__IEEE8021_P__TAG__MODE__TRANSPARENT = 9
  COS__NODE__FEATURE__CLFR__IEEE8021_AD = 10
  COS__NODE__FEATURE__CLFR__IEEE8021__AD__TAG__MODE__INNER = 11
  COS__NODE__FEATURE__CLFR__NO__DEFAULT = 12
)

/*
 **
 * Rewrite family
 */
const (
  COS__NODE__FEATURE__RW__PROTO__INVALID = 0
  COS__NODE__FEATURE__RW__DSCP = 1
  COS__NODE__FEATURE__RW__DSCP__MPLS = 2
  COS__NODE__FEATURE__RW__DSCP__IPV6 = 3
  COS__NODE__FEATURE__RW__DSCP__IPV6__MPLS = 4
  COS__NODE__FEATURE__RW__EXP = 5
  COS__NODE__FEATURE__RW__EXP__MPLS__INET__BOTH = 6
  COS__NODE__FEATURE__RW__EXP__MPLS__INET__BOTH__NON__VPN = 7
  COS__NODE__FEATURE__RW__INET__PRECEDENCE = 8
  COS__NODE__FEATURE__RW__INET__PREC__MPLS = 9
  COS__NODE__FEATURE__RW__IEEE8021_P = 10
  COS__NODE__FEATURE__RW__IEEE8021__TAG__MODE__OUTER__AND__INNER = 11
  COS__NODE__FEATURE__RW__IEEE8021_AD = 12
  COS__NODE__FEATURE__RW__IEEE8021_AD__TAG__MODE__OUTER__AND__INNER = 13
)

/*
 **
 * Forwarding class policing priority
 */
const (
  COS__NORMAL = 0
  COS__PREMIUM = 1
)

/*
 **
 * Forwarding class fabric priority
 */
const (
  COS__FAB__LOW = 0
  COS__FAB__HIGH = 1
)

/*
 **
 * Forwarding class spu priority
 */
const (
  COS__SPU__LOW = 0
  COS__SPU__HIGH = 1
)

const (
  /*
   * To retrieve mapped features.
   */
  COS__FEATURES__MAPPED = 0
  /*
   * To retrieve all configured features.
   */
  COS__FEATURES__CONFIGURED = 1
  /*
   * To retrieve all features.
   */
  COS__FEATURES__ALL = 2
)

const (
  COS__RI__FEATURE__TYPE__INVALID = 0
  COS__RI__FEATURE__TYPE__CLFR__RULE = 1
  COS__RI__FEATURE__TYPE__RW__RULE = 2
)

const (
  COS__RI__CLFR__CP__TYPE__INVALID = 0
  COS__RI__CLFR__CP__TYPE__DSCP = 1
  COS__RI__CLFR__CP__TYPE__DSCP__IPV6 = 2
  COS__RI__CLFR__CP__TYPE__EXP = 3
  COS__RI__CLFR__CP__TYPE__IEEE8021 = 4
  COS__RI__CLFR__CP__TYPE__NO__DEFAULT = 5
)

const (
  COS__RI__RW__CP__TYPE__INVALID = 0
  COS__RI__RW__CP__TYPE__IEEE8021 = 1
  COS__RI__RW__CP__TYPE__IEEE8021_AD = 2
)

/*
 **
 * Classifier tag modes
 */
const (
  COS__RI__CLFR__IEEE8021__TAG__MODE__INVALID = 0
  COS__RI__CLFR__IEEE8021__TAG__MODE__ENCAP__VLAN__TAG__INNER = 1
  COS__RI__CLFR__IEEE8021__TAG__MODE__ENCAP__VLAN__TAG__OUTER = 2
)

/*
 **
 * Rewrite tag modes
 */
const (
  COS__RI__RW__IEEE8021_X__TAG__MODE__INVALID = 0
  COS__RI__RW__IEEE8021_X__TAG__MODE__ENCAP__VLAN__TAG__OUTER = 1
  COS__RI__RW__IEEE8021_X__TAG__MODE__ENCAP__VLAN__TAG__OUTER__AND__INNER = 2
)

/*
 * Status  message for ADD/CHANGE/DELETE requests
 */
type Cos__CosStatus struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

/*
 * optional for GET/GETNEXT/GETBULK
 * to see referencing objects for a given cos object.
 */
type Cos__CosObjRefInfo struct {
    ReferencingObjType	 string `json:"referencing_obj_type,omitempty"`
    ReferencingObjName	 string `json:"referencing_obj_name,omitempty"`
    ReferencingObjCount	 int32 `json:"referencing_obj_count,omitempty"`
}

/*
 **
 * Following is the RPC get-xxx request message for forwarding classes
 * queried by forwarding class name
 * RPC Methods:
 * ============
 * ForwardClassGet
 * ForwardClassGetNext
 * ForwardClassBulkGet
 * Request message:
 * ===============
 * Request must have forwarding_class_name specified for Get request
 * For GetNext and BulkGet requests, its optional
 */
type Cos__CosForwardingClassQueryRequest struct {
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

/*
 **
 * Following is the RPC response message for forwarding classes
 * queried by forwarding class name
 * RPC Methods:
 * ============
 * ForwardClassGet
 * ForwardClassGetNext
 * ForwardClassBulkGet
 * Response message:
 * ================
 * All the fields of the message will be populated
 * for requested forwarding class name
 * 
 */
type Cos__CosForwardingClassQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    QueueId	 int32 `json:"queue_id,omitempty"`
    RestrictQueueId	 int32 `json:"restrict_queue_id,omitempty"`
    ForwardingClassId	 int32 `json:"forwarding_class_id,omitempty"`
    FabricPriority	 uint32 `json:"fabric_priority,omitempty"`
    PolicingPriority	 uint32 `json:"policing_priority,omitempty"`
    SpuPriority	 uint32 `json:"spu_priority,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

/*
 **
 * Following is the RPC request message for forwarding classes
 * queried by forwarding class id
 * RPC Methods:
 * ============
 * ForwardClassGetByForwardingClassId
 * Request message:
 * ===============
 * Request will have forwarding_class_id field specified
 *  -- By default forwarding_class_id 0 will be used
 *     if not specified explicitly
 */
type Cos__CosForwardingClassIdQueryRequest struct {
    ForwardingClassId	 int32 `json:"forwarding_class_id,omitempty"`
}

/*
 **
 * Following is the RPC response message for forwarding classes
 * queried by forwarding class id
 * RPC Methods:
 * ============
 * ForwardClassGetByForwardingClassId
 * Response message:
 * ================
 * All the fields of the message will be populated
 * for requested forwarding class id
 * 
 */
type Cos__CosForwardingClassIdQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    ForwardingClassId	 int32 `json:"forwarding_class_id,omitempty"`
    QueueId	 int32 `json:"queue_id,omitempty"`
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    RestrictQueueId	 int32 `json:"restrict_queue_id,omitempty"`
    FabricPriority	 uint32 `json:"fabric_priority,omitempty"`
    PolicingPriority	 uint32 `json:"policing_priority,omitempty"`
    SpuPriority	 uint32 `json:"spu_priority,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
}

/*
 **
 * Following is the RPC request message for forwarding classes
 * queried by queue id
 * RPC Methods:
 * ============
 * ForwardClassGetByQueueId
 * Request message:
 * ===============
 * Request will have queue_id field specified
 *  -- By default queue_id 0 will be used if not specified explicitly
 */
type Cos__CosForwardingClassQueueQueryRequest struct {
    QueueId	 int32 `json:"queue_id,omitempty"`
}

/*
 **
 * Following is the RPC response message for forwarding classes
 * queried by queue id
 * RPC Methods:
 * ============
 * ForwardClassGetByQueueId
 * Response message:
 * ================
 * Response message will have
 * >  queue_id 
 * >  list of forwarding classes for the requested queue id
 * 
 */
type Cos__CosForwardingClassQueueQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    QueueId	 int32 `json:"queue_id,omitempty"`
    ForwardingClassNames	 []string `json:"forwarding_class_names,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
}

/*
 **
 * Following is the RPC request message for forwarding classes
 * queried by restricted queue id
 * RPC Methods:
 * ============
 * ForwardClassGetByRestrictQueueId
 * Request message:
 * ===============
 * Request will have restricte_queue_id field specified
 *  -- By default restrict_queue_id 0 will be used
 *     if not specified explicitly
 * 
 */
type Cos__CosForwardingClassRestrictQueueQueryRequest struct {
    RestrictQueueId	 int32 `json:"restrict_queue_id,omitempty"`
}

/*
 **
 * Following is the RPC response message for forwarding classes
 * queried by restricted queue id
 * RPC Methods:
 * ============
 * ForwardClassGetByRestrictQueueId
 * Response message:
 * ================
 * Response message will have
 * > restrict_queue_id
 * > list of forwarding classes for the requested restricted queue id
 * 
 * 
 */
type Cos__CosForwardingClassRestrictQueueQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    RestrictQueueId	 int32 `json:"restrict_queue_id,omitempty"`
    ForwardingClassNames	 []string `json:"forwarding_class_names,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
}

type Cos__CosConfigPreference struct {
    Priority	 int32 `json:"priority,omitempty"`
}

/*
 **
 * Classifier rules
 * When operation is delete, the forwarding class and loss-priority are optional or ignored.
 */
type Cos__ClassifierRule struct {
    Operation	 uint32 `json:"operation,omitempty"`
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    LossPriority	 uint32 `json:"loss_priority,omitempty"`
    CodePoints	 []int32 `json:"code_points,omitempty"`
}

/*
 **
 * Following is the RPC Delete request message for classifiers
 * RPC Methods:
 * ============
 * ClassifierDelete
 * Request message:
 * ===============
 * Delete requests must have classifier name and type
 * Response message:
 * ================
 * Delete requests will return CosStatus
 * 
 */
type Cos__CosClassifierDeleteRequest struct {
    ClassifierName	 string `json:"classifier_name,omitempty"`
    ClassifierType	 uint32 `json:"classifier_type,omitempty"`
}

/*
 **
 * Following is the RPC Query request message for classifiers
 * RPC Methods:
 * ============
 * ClassifierGet
 * ClassifierGetNext
 * ClassifierBulkGet
 * Request message:
 * ===============
 * Get requests must have classifier name and type
 * GetNext/BulkGet requests optionally have name and type,
 *                 i.e type, type & name, or none of them. 
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosClassifierQueryResponse info
 * 
 */
type Cos__CosClassifierQueryRequest struct {
    ClassifierName	 string `json:"classifier_name,omitempty"`
    ClassifierType	 uint32 `json:"classifier_type,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

/*
 **
 * Following is the RPC Add/Update request message for classifiers
 * RPC Methods:
 * ============
 * ClassifierAdd
 * ClassifierUpdate
 * Request message:
 * ===============
 * Add, Update requests must have classifier name and type along with the rules
 * Response message:
 * ================
 * Add/Update requests will return CosStatus
 * 
 */
type Cos__CosClassifierRequest struct {
    ClassifierName	 string `json:"classifier_name,omitempty"`
    ClassifierType	 uint32 `json:"classifier_type,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    Rule	 []Cos__ClassifierRule `json:"rule,omitempty"`
}

/*
 **
 * Following is the RPC query response message for classifiers
 * RPC Methods:
 * ============
 * ClassifierGet
 * ClassifierGetNext
 * ClassifierBulkGet
 * Request message:
 * ===============
 * Get requests must have classifier name and type
 * GetNext/BulkGet requests optionally have name and type,
 *                 i.e type, type & name, or none of them. 
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosClassifierQueryResponse info
 * 
 */
type Cos__CosClassifierQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    ClassifierName	 string `json:"classifier_name,omitempty"`
    ClassifierType	 uint32 `json:"classifier_type,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    Rule	 []Cos__ClassifierRule `json:"rule,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

/*
 **
 * Code points
 */
type Cos__FeatureCodePoint struct {
  /*
   * Only one of the following fields should be specified: 
   * CodePoint, CodePointStr
   */
    CodePoint	 int32 `json:"code_point,omitempty"`
    CodePointStr	 string `json:"code_point_str,omitempty"`
}

/*
 **
 * Rewrite Rules
 * When operation is delete, the code point is optional or ignored.
 */
type Cos__RewriteRule struct {
    Operation	 uint32 `json:"operation,omitempty"`
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    LossPriority	 uint32 `json:"loss_priority,omitempty"`
    CodePointOptions	 *Cos__FeatureCodePoint `json:"code_point_options,omitempty"`
}

/*
 **
 * Following is the RPC delete request message for rewrites
 * RPC Methods:
 * ============
 * RewriteDelete
 * Request message:
 * ===============
 * Delete requests must have rewrite name and type
 * Response message:
 * ================
 * Delete requests will return CosStatus
 * 
 */
type Cos__CosRewriteDeleteRequest struct {
    RewriteName	 string `json:"rewrite_name,omitempty"`
    RewriteType	 uint32 `json:"rewrite_type,omitempty"`
}

/*
 **
 * Following is the RPC query request message for rewrites
 * RPC Methods:
 * ============
 * RewriteGet
 * RewriteGetNext
 * RewriteBulkGet
 * Request message:
 * ===============
 * Get requests must have rewrite name and type
 * GetNext/BulkGet requests optionally have name and type
 *                 i.e type, type & name, or none of them. 
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosRewriteQueryResponse info
 * 
 */
type Cos__CosRewriteQueryRequest struct {
    RewriteName	 string `json:"rewrite_name,omitempty"`
    RewriteType	 uint32 `json:"rewrite_type,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

/*
 **
 * Following is the RPC Add/update request message for rewrites
 * RPC Methods:
 * ============
 * RewriteAdd
 * RewriteUpdate
 * Request message:
 * ===============
 * Add, Update requests must have rewrite name and type along with the rules
 * Response message:
 * ================
 * Add/Update requests will return CosStatus
 * 
 */
type Cos__CosRewriteRequest struct {
    RewriteName	 string `json:"rewrite_name,omitempty"`
    RewriteType	 uint32 `json:"rewrite_type,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    Rule	 []Cos__RewriteRule `json:"rule,omitempty"`
}

/*
 **
 * Following is the RPC query response message for rewrites
 * RPC Methods:
 * ============
 * RewriteGet
 * RewriteGetNext
 * RewriteBulkGet
 * Response message:
 * ================
 * Get/GetNext/BulkGet responds with CosRewriteQueryResponse info
 * 
 */
type Cos__CosRewriteQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    RewriteName	 string `json:"rewrite_name,omitempty"`
    RewriteType	 uint32 `json:"rewrite_type,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    Rule	 []Cos__RewriteRule `json:"rule,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

type Cos__LevelToDropProbability struct {
    Operation	 uint32 `json:"operation,omitempty"`
    FillLevel	 int32 `json:"fill_level,omitempty"`
    DropProbability	 int32 `json:"drop_probability,omitempty"`
}

/*
 **
 * RPC message for drop profile
 */
type Cos__CosDropProfile struct {
    DropProfileName	 string `json:"drop_profile_name,omitempty"`
    Sharable	 bool `json:"sharable,omitempty"`
    DropProfileType	 uint32 `json:"drop_profile_type,omitempty"`
    LevelToDropProbabilities	 []Cos__LevelToDropProbability `json:"level_to_drop_probabilities,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

type Cos__DropProfileMap struct {
    Operation	 uint32 `json:"operation,omitempty"`
    SchedulerLossPriority	 uint32 `json:"scheduler_loss_priority,omitempty"`
    SchedulerProtocol	 uint32 `json:"scheduler_protocol,omitempty"`
    DropProfileName	 string `json:"drop_profile_name,omitempty"`
    DropProfileShared	 bool `json:"drop_profile_shared,omitempty"`
}

type Cos__SchedulerDropProfileMap struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Sets	 []Cos__DropProfileMap `json:"sets,omitempty"`
}

type Cos__SchedulerShapingRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
    BurstSize	 uint64 `json:"burst_size,omitempty"`
}

type Cos__ExcessRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
}

type Cos__BufferSize struct {
    Operation	 uint32 `json:"operation,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
}

type Cos__AdjustMinRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
}

type Cos__AdjustPercent struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 int32 `json:"value,omitempty"`
}

type Cos__GuaranteedRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueExtendSpec	 uint32 `json:"value_extend_spec,omitempty"`
}

type Cos__SchedulerPriority struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint32 `json:"value,omitempty"`
}

type Cos__SchedulerExcessPriority struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint32 `json:"value,omitempty"`
}

/*
 **
 * RPC message for scheduler
 */
type Cos__CosScheduler struct {
    SchedulerName	 string `json:"scheduler_name,omitempty"`
    AdjustMinRate	 *Cos__AdjustMinRate `json:"adjust_min_rate,omitempty"`
    AdjustPercent	 *Cos__AdjustPercent `json:"adjust_percent,omitempty"`
    GRate	 *Cos__GuaranteedRate `json:"g_rate,omitempty"`
    SRate	 *Cos__SchedulerShapingRate `json:"s_rate,omitempty"`
    ERate	 *Cos__ExcessRate `json:"e_rate,omitempty"`
    BSize	 *Cos__BufferSize `json:"b_size,omitempty"`
    Priority	 *Cos__SchedulerPriority `json:"priority,omitempty"`
    ExcessPriority	 *Cos__SchedulerExcessPriority `json:"excess_priority,omitempty"`
    DropProfileMappings	 []Cos__DropProfileMap `json:"drop_profile_mappings,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

type Cos__ForwardingClassToScheduler struct {
    Operation	 uint32 `json:"operation,omitempty"`
    ForwardingClassName	 string `json:"forwarding_class_name,omitempty"`
    SchedulerName	 string `json:"scheduler_name,omitempty"`
}

/*
 **
 * RPC message for scheduler maps
 */
type Cos__CosSchedulerMap struct {
    SchedulerMapName	 string `json:"scheduler_map_name,omitempty"`
    FcToSchedulerMapping	 []Cos__ForwardingClassToScheduler `json:"fc_to_scheduler_mapping,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

type Cos__TcpAdjustRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
}

type Cos__DelayBufferRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
}

type Cos__TcpGuaranteedRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
}

type Cos__TcpShapingRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
    ValueSpec	 uint32 `json:"value_spec,omitempty"`
}

type Cos__OverheadAccounting struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 int32 `json:"value,omitempty"`
    Mode	 uint32 `json:"mode,omitempty"`
}

type Cos__TcpSchedulerMap struct {
    Operation	 uint32 `json:"operation,omitempty"`
    SchedulerMapName	 string `json:"scheduler_map_name,omitempty"`
}

type Cos__BurstSize struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
}

type Cos__CosRate struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Value	 uint64 `json:"value,omitempty"`
}

type Cos__CosAtmService struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Category	 uint32 `json:"category,omitempty"`
}

/*
 **
 * RPC message for TCPs
 */
type Cos__CosTrafficControlProfile struct {
    TrafficControlProfileName	 string `json:"traffic_control_profile_name,omitempty"`
    SchedulerMap	 *Cos__TcpSchedulerMap `json:"scheduler_map,omitempty"`
    ARate	 *Cos__TcpAdjustRate `json:"a_rate,omitempty"`
    DelayBufRate	 *Cos__DelayBufferRate `json:"delay_buf_rate,omitempty"`
    GRate	 *Cos__TcpGuaranteedRate `json:"g_rate,omitempty"`
    GuaranteedBurstSize	 *Cos__BurstSize `json:"guaranteed_burst_size,omitempty"`
    SRate	 *Cos__TcpShapingRate `json:"s_rate,omitempty"`
    ShapingBurstSize	 *Cos__BurstSize `json:"shaping_burst_size,omitempty"`
    SRatePriorityHigh	 *Cos__TcpShapingRate `json:"s_rate_priority_high,omitempty"`
    ShapingBurstSizePriorityHigh	 *Cos__BurstSize `json:"shaping_burst_size_priority_high,omitempty"`
    SRatePriorityLow	 *Cos__TcpShapingRate `json:"s_rate_priority_low,omitempty"`
    ShapingBurstSizePriorityLow	 *Cos__BurstSize `json:"shaping_burst_size_priority_low,omitempty"`
    SRatePriorityMedium	 *Cos__TcpShapingRate `json:"s_rate_priority_medium,omitempty"`
    ShapingBurstSizePriorityMedium	 *Cos__BurstSize `json:"shaping_burst_size_priority_medium,omitempty"`
    SRateExcessPriorityHigh	 *Cos__TcpShapingRate `json:"s_rate_excess_priority_high,omitempty"`
    ShapingBurstSizeExcessPriorityHigh	 *Cos__BurstSize `json:"shaping_burst_size_excess_priority_high,omitempty"`
    SRateExcessPriorityLow	 *Cos__TcpShapingRate `json:"s_rate_excess_priority_low,omitempty"`
    ShapingBurstSizeExcessPriorityLow	 *Cos__BurstSize `json:"shaping_burst_size_excess_priority_low,omitempty"`
    ERate	 *Cos__ExcessRate `json:"e_rate,omitempty"`
    ERateLow	 *Cos__ExcessRate `json:"e_rate_low,omitempty"`
    ERateHigh	 *Cos__ExcessRate `json:"e_rate_high,omitempty"`
    OHeadAccounting	 *Cos__OverheadAccounting `json:"o_head_accounting,omitempty"`
    SRatePriorityMediumLow	 *Cos__TcpShapingRate `json:"s_rate_priority_medium_low,omitempty"`
    ShapingBurstSizePriorityMediumLow	 *Cos__BurstSize `json:"shaping_burst_size_priority_medium_low,omitempty"`
    SRatePriorityStrictHigh	 *Cos__TcpShapingRate `json:"s_rate_priority_strict_high,omitempty"`
    ShapingBurstSizePriorityStrictHigh	 *Cos__BurstSize `json:"shaping_burst_size_priority_strict_high,omitempty"`
    SRateExcessPriorityMeidumHigh	 *Cos__TcpShapingRate `json:"s_rate_excess_priority_meidum_high,omitempty"`
    ShapingBurstSizeExcessPriorityMediumHigh	 *Cos__BurstSize `json:"shaping_burst_size_excess_priority_medium_high,omitempty"`
    SRateExcessPriorityMeidumLow	 *Cos__TcpShapingRate `json:"s_rate_excess_priority_meidum_low,omitempty"`
    ShapingBurstSizeExcessPriorityMediumLow	 *Cos__BurstSize `json:"shaping_burst_size_excess_priority_medium_low,omitempty"`
    AtmService	 *Cos__CosAtmService `json:"atm_service,omitempty"`
    PeakRate	 *Cos__CosRate `json:"peak_rate,omitempty"`
    SustainedRate	 *Cos__CosRate `json:"sustained_rate,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    RefObjects	 []Cos__CosObjRefInfo `json:"ref_objects,omitempty"`
}

type Cos__CosNodeFeatureVal struct {
  /*
   * Only one of the following fields should be specified: 
   * Value, ObjectName
   */
    Value	 uint64 `json:"value,omitempty"`
    ObjectName	 string `json:"object_name,omitempty"`
}

type Cos__CosNodeFeatureOption struct {
  /*
   * Only one of the following fields should be specified: 
   * ClassifierFamily, RewriteFamily
   */
    ClassifierFamily	 uint32 `json:"classifier_family,omitempty"`
    RewriteFamily	 uint32 `json:"rewrite_family,omitempty"`
}

/*
 **
 * Node family features
 */
type Cos__CosNodeFamilyFeature struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Direction	 uint32 `json:"direction,omitempty"`
    FeatureType	 uint32 `json:"feature_type,omitempty"`
    FeatureOption	 *Cos__CosNodeFeatureOption `json:"feature_option,omitempty"`
    FeatureObjectName	 string `json:"feature_object_name,omitempty"`
    ObjectShared	 bool `json:"object_shared,omitempty"`
}

type Cos__CosNodeFeature struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Direction	 uint32 `json:"direction,omitempty"`
    FeatureType	 uint32 `json:"feature_type,omitempty"`
    FeatureVal	 *Cos__CosNodeFeatureVal `json:"feature_val,omitempty"`
}

/*
 **
 * Following is the RPC delete request message for 
 * interface, logical interface, and interface sets
 * RPC Methods:
 * ============
 * NodeFeaturesDelete
 * Request message:
 * ===============
 * Delete requests must have
 *  1. Node name and type
 *  2. Parent node name and type ( optional in case if node is IFL/logical interface, and ignored for IFD )
 * Response message:
 * ================
 * Delete requests will returns CosStatus message.
 * 
 */
type Cos__CosNodeBindFeaturesDeleteRequest struct {
    NodeType	 uint32 `json:"node_type,omitempty"`
    NodeName	 string `json:"node_name,omitempty"`
    NodeParentType	 uint32 `json:"node_parent_type,omitempty"`
    NodeParentName	 string `json:"node_parent_name,omitempty"`
}

/*
 **
 * Following is the RPC request message for 
 * interface, logical interface, and interface sets
 * RPC Methods:
 * ============
 * NodeFeaturesGet
 * NodeFeaturesGetNext
 * NodeFeaturesBulkGet
 * Request message:
 * ===============
 * Get requests must have
 *  1. Node name and type
 *  2. Parent node name and type ( ignored  )
 * GetNext/BulkGet requests can have
 *  1. Node type
 *  2. Node type and node name
 *  3. Parent node name and type ( Ignored )
 *  4. Nothing
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosNodeBindFeaturesQueryResponse
 * 
 */
type Cos__CosNodeBindFeaturesQueryRequest struct {
    NodeType	 uint32 `json:"node_type,omitempty"`
    NodeName	 string `json:"node_name,omitempty"`
    NodeParentType	 uint32 `json:"node_parent_type,omitempty"`
    NodeParentName	 string `json:"node_parent_name,omitempty"`
    FeaturesType	 uint32 `json:"features_type,omitempty"`
}

/*
 **
 * Following is the RPC Add/Update request message for 
 * interface, logical interface, and interface sets
 * RPC Methods:
 * ============
 * NodeFeaturesAdd
 * NodeFeaturesUpdate
 * Request message:
 * ===============
 * Add requests must have
 *  1. Node name and type
 *  2. Parent node name and type ( in case if node is IFL/logical interface, and ignored for IFD)
 *  3. Node features AND/OR Family features [optional for IFD, must for IFL/logical interface)
 * Update requests must have
 *  1. Node name and type
 *  2. Parent node name and type ( optional for IFL/logical interface and will be validated if specified, and ignored for IFD)
 *  3. Node features AND/OR Family features [for this RLI, it optional for IFD, must for IFL/logical interface)
 * Response message:
 * ================
 * Add/Update requests will return CosStatus
 * 
 */
type Cos__CosNodeBindFeaturesRequest struct {
    NodeType	 uint32 `json:"node_type,omitempty"`
    NodeName	 string `json:"node_name,omitempty"`
    NodeParentType	 uint32 `json:"node_parent_type,omitempty"`
    NodeParentName	 string `json:"node_parent_name,omitempty"`
    NodeFeatures	 []Cos__CosNodeFeature `json:"node_features,omitempty"`
    NodeFamilyFeatures	 []Cos__CosNodeFamilyFeature `json:"node_family_features,omitempty"`
}

/*
 **
 * Following is the RPC query response message for 
 * interface, logical interface, and interface sets
 * RPC Methods:
 * ============
 * NodeFeaturesGet
 * NodeFeaturesGetNext
 * NodeFeaturesBulkGet
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will returns CosNodeBindFeaturesQueryResponse
 * 
 */
type Cos__CosNodeBindFeaturesQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    NodeType	 uint32 `json:"node_type,omitempty"`
    NodeName	 string `json:"node_name,omitempty"`
    NodeParentType	 uint32 `json:"node_parent_type,omitempty"`
    NodeParentName	 string `json:"node_parent_name,omitempty"`
    FeaturesType	 uint32 `json:"features_type,omitempty"`
    NodeFeatures	 []Cos__CosNodeFeature `json:"node_features,omitempty"`
    NodeFamilyFeatures	 []Cos__CosNodeFamilyFeature `json:"node_family_features,omitempty"`
}

/*
 **
 * Routing instance classifiers
 */
type Cos__CosRoutingInstanceClassifier struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Type	 uint32 `json:"type,omitempty"`
    FeatureIeee8021TagMode	 uint32 `json:"feature_ieee8021_tag_mode,omitempty"`
    FeatureObjectName	 string `json:"feature_object_name,omitempty"`
    ObjectShared	 bool `json:"object_shared,omitempty"`
}

/*
 **
 * Routing instance rewrites
 */
type Cos__CosRoutingInstanceRewrite struct {
    Operation	 uint32 `json:"operation,omitempty"`
    Type	 uint32 `json:"type,omitempty"`
    FeatureIeee8021TagMode	 uint32 `json:"feature_ieee8021_tag_mode,omitempty"`
    FeatureObjectName	 string `json:"feature_object_name,omitempty"`
    ObjectShared	 bool `json:"object_shared,omitempty"`
}

/*
 **
 * Following is the RPC Delete request message for 
 * Routing instances
 * RPC Methods:
 * ============
 * RoutingInstanceBindPointDelete
 * Request message:
 * ===============
 * Delete requests must have
 *  1. Routing instance name
 * Response message:
 * ================
 * Delete requests will CosStatus
 */
type Cos__CosRoutingInstanceBindPointDeleteRequest struct {
    RoutingInstanceName	 string `json:"routing_instance_name,omitempty"`
}

/*
 **
 * Following is the RPC query request message for 
 * Routing instances
 * RPC Methods:
 * ============
 * RoutingInstanceBindPointGet
 * RoutingInstanceBindPointGetNext
 * RoutingInstanceBindPointBulkGet
 * Request message:
 * ===============
 * Get requests must have
 *  1. Routing instance name
 * GetNext/BulkGet requests can have
 *  1. Routing instance name or NULL
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosRoutingInstanceBindPointQueryResponse
 * 
 */
type Cos__CosRoutingInstanceBindPointQueryRequest struct {
    RoutingInstanceName	 string `json:"routing_instance_name,omitempty"`
}

/*
 **
 * Following is the RPC Add/Update request message for 
 * Routing instances
 * RPC Methods:
 * ============
 * RoutingInstanceBindPointAdd
 * RoutingInstanceBindPointUpdate
 * Request message:
 * ===============
 * Add, Update requests must have
 *  1. Routing instance name
 *  2. Routing instances classifier/rewrite features
 * Response message:
 * ================
 * Add/Update requests will return CosStatus
 * 
 */
type Cos__CosRoutingInstanceBindPointRequest struct {
    RoutingInstanceName	 string `json:"routing_instance_name,omitempty"`
    Classifiers	 []Cos__CosRoutingInstanceClassifier `json:"classifiers,omitempty"`
    Rewrites	 []Cos__CosRoutingInstanceRewrite `json:"rewrites,omitempty"`
}

/*
 **
 * Following is the RPC Query response message for 
 * Routing instances
 * RPC Methods:
 * ============
 * RoutingInstanceBindPointGet
 * RoutingInstanceBindPointGetNext
 * RoutingInstanceBindPointBulkGet
 * Response message:
 * ================
 * Get/GetNext/BulkGet requests will return CosRoutingInstanceBindPointQueryResponse
 * 
 */
type Cos__CosRoutingInstanceBindPointQueryResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    RoutingInstanceName	 string `json:"routing_instance_name,omitempty"`
    Classifiers	 []Cos__CosRoutingInstanceClassifier `json:"classifiers,omitempty"`
    Rewrites	 []Cos__CosRoutingInstanceRewrite `json:"rewrites,omitempty"`
}

type Cos__CosResourceLimit struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    NumOfRateLimitQs	 int32 `json:"num_of_rate_limit_qs,omitempty"`
    BufSizeTemporalLimitValue	 int32 `json:"buf_size_temporal_limit_value,omitempty"`
    MaxUcastFcSetsPerInterface	 int32 `json:"max_ucast_fc_sets_per_interface,omitempty"`
    MaxMcastFcSetsPerInterface	 int32 `json:"max_mcast_fc_sets_per_interface,omitempty"`
    MaxWredProfilePoints	 int32 `json:"max_wred_profile_points,omitempty"`
    NumOfFcsInFcSet	 int32 `json:"num_of_fcs_in_fc_set,omitempty"`
    MaxQueueToPauseProfiles	 int32 `json:"max_queue_to_pause_profiles,omitempty"`
    MaxIngressSharedBuffer	 int32 `json:"max_ingress_shared_buffer,omitempty"`
    MaxEgressSharedBuffer	 int32 `json:"max_egress_shared_buffer,omitempty"`
    MaxGlobalBuffer	 int32 `json:"max_global_buffer,omitempty"`
    MaxQcnQueues	 int32 `json:"max_qcn_queues,omitempty"`
    MaxQueuesPerInterface	 int32 `json:"max_queues_per_interface,omitempty"`
    DefaultFcSize	 int32 `json:"default_fc_size,omitempty"`
}

/*
 **
 * Following is the RPC Update request message for tunable purge time
 * to cleanup client state of configuration.
 * Note: 
 *  RPC - CosPurgeTimeDelete
 *  Delete requests resets purge time to default value of 300 secs and 
 *  will return CosStatus
 * RPC Methods:
 * ============
 * CosPurgeTimeUpdate
 * Request message:
 * ===============
 * Update requests must have purge time in secs from 30-86400
 * Response message:
 * ================
 * Update requests will return CosStatus
 * 
 */
type Cos__CosPurgeTimeRequest struct {
    PurgeTime	 int32 `json:"purge_time,omitempty"`
}

/*
 **
 * Following is the RPC response message for tunable purge time
 * to cleanup client state of configuration.
 * RPC Methods:
 * ============
 * CosPurgeTimeGet
 * Response message:
 * ================
 * Get requests will return CosPurgeTimeResponse
 * 
 */
type Cos__CosPurgeTimeResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Message	 string `json:"message,omitempty"`
    PurgeTime	 int32 `json:"purge_time,omitempty"`
}

type Cos__Cos_CosForwardingClassGet struct {
  Request	Cos__CosForwardingClassQueryRequest
  Reply	Cos__CosForwardingClassQueryResponse
}

func (r *Cos__Cos_CosForwardingClassGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosForwardingClassGetNext struct {
  Request	Cos__CosForwardingClassQueryRequest
  Reply	Cos__CosForwardingClassQueryResponse
}

func (r *Cos__Cos_CosForwardingClassGetNext) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassGetNext"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassGetNext) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosForwardingClassBulkGet struct {
  Request	Cos__CosForwardingClassQueryRequest
  Reply	Cos__CosForwardingClassQueryResponse
}

func (r *Cos__Cos_CosForwardingClassBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosForwardingClassGetByForwardingClassId struct {
  Request	Cos__CosForwardingClassIdQueryRequest
  Reply	Cos__CosForwardingClassIdQueryResponse
}

func (r *Cos__Cos_CosForwardingClassGetByForwardingClassId) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassGetByForwardingClassId"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassGetByForwardingClassId) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosForwardingClassGetByQueueId struct {
  Request	Cos__CosForwardingClassQueueQueryRequest
  Reply	Cos__CosForwardingClassQueueQueryResponse
}

func (r *Cos__Cos_CosForwardingClassGetByQueueId) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassGetByQueueId"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassGetByQueueId) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosForwardingClassGetByRestrictQueueId struct {
  Request	Cos__CosForwardingClassRestrictQueueQueryRequest
  Reply	Cos__CosForwardingClassRestrictQueueQueryResponse
}

func (r *Cos__Cos_CosForwardingClassGetByRestrictQueueId) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosForwardingClassGetByRestrictQueueId"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosForwardingClassGetByRestrictQueueId) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierAdd struct {
  Request	Cos__CosClassifierRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosClassifierAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierDelete struct {
  Request	Cos__CosClassifierDeleteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosClassifierDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierUpdate struct {
  Request	Cos__CosClassifierRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosClassifierUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierGet struct {
  Request	Cos__CosClassifierQueryRequest
  Reply	Cos__CosClassifierQueryResponse
}

func (r *Cos__Cos_CosClassifierGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierGetNext struct {
  Request	Cos__CosClassifierQueryRequest
  Reply	Cos__CosClassifierQueryResponse
}

func (r *Cos__Cos_CosClassifierGetNext) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierGetNext"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierGetNext) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosClassifierBulkGet struct {
  Request	Cos__CosClassifierQueryRequest
  Reply	Cos__CosClassifierQueryResponse
}

func (r *Cos__Cos_CosClassifierBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosClassifierBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosClassifierBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteAdd struct {
  Request	Cos__CosRewriteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRewriteAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteDelete struct {
  Request	Cos__CosRewriteDeleteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRewriteDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteUpdate struct {
  Request	Cos__CosRewriteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRewriteUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteGet struct {
  Request	Cos__CosRewriteQueryRequest
  Reply	Cos__CosRewriteQueryResponse
}

func (r *Cos__Cos_CosRewriteGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteGetNext struct {
  Request	Cos__CosRewriteQueryRequest
  Reply	Cos__CosRewriteQueryResponse
}

func (r *Cos__Cos_CosRewriteGetNext) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteGetNext"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteGetNext) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRewriteBulkGet struct {
  Request	Cos__CosRewriteQueryRequest
  Reply	Cos__CosRewriteQueryResponse
}

func (r *Cos__Cos_CosRewriteBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRewriteBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRewriteBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesAdd struct {
  Request	Cos__CosNodeBindFeaturesRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosNodeBindFeaturesAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesDelete struct {
  Request	Cos__CosNodeBindFeaturesDeleteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosNodeBindFeaturesDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesUpdate struct {
  Request	Cos__CosNodeBindFeaturesRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosNodeBindFeaturesUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesGet struct {
  Request	Cos__CosNodeBindFeaturesQueryRequest
  Reply	Cos__CosNodeBindFeaturesQueryResponse
}

func (r *Cos__Cos_CosNodeBindFeaturesGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesGetNext struct {
  Request	Cos__CosNodeBindFeaturesQueryRequest
  Reply	Cos__CosNodeBindFeaturesQueryResponse
}

func (r *Cos__Cos_CosNodeBindFeaturesGetNext) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesGetNext"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesGetNext) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosNodeBindFeaturesBulkGet struct {
  Request	Cos__CosNodeBindFeaturesQueryRequest
  Reply	Cos__CosNodeBindFeaturesQueryResponse
}

func (r *Cos__Cos_CosNodeBindFeaturesBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosNodeBindFeaturesBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosNodeBindFeaturesBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointAdd struct {
  Request	Cos__CosRoutingInstanceBindPointRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRoutingInstanceBindPointAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointDelete struct {
  Request	Cos__CosRoutingInstanceBindPointDeleteRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRoutingInstanceBindPointDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointUpdate struct {
  Request	Cos__CosRoutingInstanceBindPointRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosRoutingInstanceBindPointUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointGet struct {
  Request	Cos__CosRoutingInstanceBindPointQueryRequest
  Reply	Cos__CosRoutingInstanceBindPointQueryResponse
}

func (r *Cos__Cos_CosRoutingInstanceBindPointGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointGetNext struct {
  Request	Cos__CosRoutingInstanceBindPointQueryRequest
  Reply	Cos__CosRoutingInstanceBindPointQueryResponse
}

func (r *Cos__Cos_CosRoutingInstanceBindPointGetNext) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointGetNext"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointGetNext) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosRoutingInstanceBindPointBulkGet struct {
  Request	Cos__CosRoutingInstanceBindPointQueryRequest
  Reply	Cos__CosRoutingInstanceBindPointQueryResponse
}

func (r *Cos__Cos_CosRoutingInstanceBindPointBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosRoutingInstanceBindPointBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosRoutingInstanceBindPointBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosPurgeTimeDelete struct {
  Request	Cos__CosPurgeTimeRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosPurgeTimeDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosPurgeTimeDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosPurgeTimeDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosPurgeTimeUpdate struct {
  Request	Cos__CosPurgeTimeRequest
  Reply	Cos__CosStatus
}

func (r *Cos__Cos_CosPurgeTimeUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosPurgeTimeUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosPurgeTimeUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Cos__Cos_CosPurgeTimeGet struct {
  Request	Cos__CosPurgeTimeRequest
  Reply	Cos__CosPurgeTimeResponse
}

func (r *Cos__Cos_CosPurgeTimeGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/cos.Cos/CosPurgeTimeGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Cos__Cos_CosPurgeTimeGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * For defining IP address and MAC address formats 
 */
const (
  /*
   * String format  
   */
  INTERFACE__ADDRESS__STRING = 0
  /*
   * Byte  format   
   */
  INTERFACE__ADDRESS__BYTES = 1
)

/*
 * For defining the protocol family type   
 */
const (
  /*
   * Unknown         
   */
  INTERFACE__INTF__AF__UNKNOWN = 0
  /*
   * INET Family     
   */
  INTERFACE__INTF__AF__INET = 2
  /*
   * INET6 Family    
   */
  INTERFACE__INTF__AF__INET6 = 28
)

/*
 * For defining interface tagging  
 */
const (
  /*
   * vlan-tagging   
   */
  INTERFACE__VLAN__TAGGING = 0
  /*
   * None           
   */
  INTERFACE__INTERFACE__TAGGING__MAX = 1
)

/*
 * For defining interface encapsulations  
 */
const (
  /*
   * vlan-bridge          
   */
  INTERFACE__VLAN__BRIDGE = 0
  /*
   * extended-vlan-bridge 
   */
  INTERFACE__VLAN__EXTENDED__BRIDGE = 1
  /*
   * None encapsulation   
   */
  INTERFACE__INTERFACE__ENCAP__MAX = 2
)

/*
 * For defining interface status 
 */
const (
  /*
   * Interface Down     
   */
  INTERFACE__INTF__DOWN = 0
  /*
   * Interface Up       
   */
  INTERFACE__INTF__UP = 1
)

const (
  INTERFACE__SUCCESS = 0
  INTERFACE__FAILURE = 1
  INTERFACE__OWNER = 2
  INTERFACE__NOT__OWNER = 3
  INTERFACE__OBJECT__FOUND = 4
  INTERFACE__OBJECT__NOT__FOUND = 5
  /*
   * Object created by CLI  commands 
   */
  INTERFACE__OBJECT__CLI__OWNER = 6
  /*
   * Object created by Japi commands 
   */
  INTERFACE__OBJECT__JAPI__OWNER = 7
  /*
   * Attribute found 
   */
  INTERFACE__ATTRIBUTE__FOUND = 8
  /*
   * Attribute not found 
   */
  INTERFACE__ATTRIBUTE__NOT__FOUND = 9
  /*
   * Port name not specified 
   */
  INTERFACE__NO__PORT__NAME = 10
  /*
   * IFL unit not specified 
   */
  INTERFACE__NO__IFL__UNIT = 11
  /*
   * IFF family not specified 
   */
  INTERFACE__NO__IFF__FAMILY = 12
  /*
   * IP address not specified 
   */
  INTERFACE__NO__IP__ADDR = 13
  /*
   * ARP address not specified 
   */
  INTERFACE__NO__ARP__ADDR = 14
  /*
   * ARP MAC address not specified 
   */
  INTERFACE__NO__ARP__MAC = 15
)

const (
  INTERFACE__UNKNOWN = 0
  INTERFACE__ADD__INTERFACE = 1
  INTERFACE__DELETE__INTERFACE = 2
  INTERFACE__ADD__AE__MEMBER = 3
  INTERFACE__DELETE__AE__MEMBER = 4
  INTERFACE__ADD__INTF__ATTRIBUTES = 5
  INTERFACE__DELETE__INTF__ATTRIBUTES = 6
  INTERFACE__DELETE__ALL__INTF__ATTRIBUTES = 7
  INTERFACE__QUERY__ATTRIBUTES = 8
  INTERFACE__QUERY__OWNERSHIP = 9
  INTERFACE__QUERY__PUBLIC__IFL = 10
  INTERFACE__SET__TIMEOUT = 11
)

/*
 *  IP address definition     
 */
type Interface__IpAddress struct {
  /*
   * Only one of the following fields should be specified: 
   * AddrString, AddrBytes
   */
    AddrString	 string `json:"addr_string,omitempty"`
    AddrBytes	 []byte `json:"addr_bytes,omitempty"`
}

/*
 *  MAC address definition     
 */
type Interface__MacAddress struct {
  /*
   * Only one of the following fields should be specified: 
   * AddrString, AddrBytes
   */
    AddrString	 string `json:"addr_string,omitempty"`
    AddrBytes	 []byte `json:"addr_bytes,omitempty"`
}

/*
 * String attribute definition   
 */
type Interface__StringAttr struct {
    AttrName	 string `json:"attr_name,omitempty"`
    AttrValue	 string `json:"attr_value,omitempty"`
}

/*
 * Integer attribute definition   
 */
type Interface__IntegerAttr struct {
    AttrName	 string `json:"attr_name,omitempty"`
    AttrValue	 int32 `json:"attr_value,omitempty"`
}

/*
 * For IFD objects configurations
 * Eg: set interfaces ge-1/1/6 vlan-tagging
 * Eg: set interfaces ge-1/1/6 encapsulation extended-vlan-bridge
 */
type Interface__InterfaceConfig struct {
    PortName	 string `json:"port_name,omitempty"`
    Tagging	 uint32 `json:"tagging,omitempty"`
    Encap	 uint32 `json:"encap,omitempty"`
    AggregateMembers	 []string `json:"aggregate_members,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    Operation	 int32 `json:"operation,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

type Interface__Value32 struct {
    Unit	 int32 `json:"unit,omitempty"`
}

/*
 * For IFL objects configurations
 * Eg: set interfaces ge-1/1/6 unit 0 vlan-id 100 
 * Eg: set bridge-domains vlan-100 interface ge-1/1/6.0 
 * Eg: set routing-instances vlan-100 interface ge-1/1/6.0 
 * Eg: set interfaces ge-1/1/6 disable
 */
type Interface__InterfaceLogicalConfig struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    LrName	 string `json:"lr_name,omitempty"`
    BdName	 string `json:"bd_name,omitempty"`
    RiName	 string `json:"ri_name,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    Operation	 int32 `json:"operation,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 * For IFF objects configurations
 * Eg: set interfaces ge-1/1/6 unit 0 vlan-id 100 
 * Eg: set bridge-domains vlan-100 interface ge-1/1/6.0 
 * Eg: set routing-instances vlan-100 interface ge-1/1/6.0 
 * Eg: set interfaces ge-1/1/6 disable
 */
type Interface__InterfaceFamilyConfig struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    Family	 uint32 `json:"family,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    Operation	 int32 `json:"operation,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 * For IFA objects configurations
 * Eg: set interfaces ge-1/1/6 unit 0 family inet address 10.10.1.1/24
 * Eg: set interfaces ge-1/1/6 unit 0 family inet6 address abcd::1/64
 */
type Interface__InterfaceAddressConfig struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    Family	 uint32 `json:"family,omitempty"`
    InterfaceAddress	 *Interface__IpAddress `json:"interface_address,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    Operation	 int32 `json:"operation,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 * For RT objects configurations
 */
type Interface__RTConfig struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    Family	 uint32 `json:"family,omitempty"`
    InterfaceAddress	 *Interface__IpAddress `json:"interface_address,omitempty"`
    ArpAddress	 *Interface__IpAddress `json:"arp_address,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    Operation	 int32 `json:"operation,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 ****************************************************
 *     Response message for configuration commands
 **************************************************
 */
type Interface__ConfigResp struct {
    Status	 uint32 `json:"status,omitempty"`
    ErrorMessage	 string `json:"error_message,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
}

/*
 *******************************************************
 *     Message structure for attribute query command    
 *****************************************************
 */
type Interface__AttributeRequestInfo struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    Family	 uint32 `json:"family,omitempty"`
    InterfaceAddress	 *Interface__IpAddress `json:"interface_address,omitempty"`
    ArpAddress	 *Interface__IpAddress `json:"arp_address,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 *******************************************************
 *     Message structure for attribute query response    
 *****************************************************
 */
type Interface__AttributeResponseInfo struct {
    Status	 uint32 `json:"status,omitempty"`
    ErrorMessage	 string `json:"error_message,omitempty"`
    StringAttrList	 []Interface__StringAttr `json:"string_attr_list,omitempty"`
    IntegerAttrList	 []Interface__IntegerAttr `json:"integer_attr_list,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 ***********************************************************
 *  Message structure for object ownership query command    
 **********************************************************
 */
type Interface__ObjectOwnershipQuery struct {
    PortName	 string `json:"port_name,omitempty"`
    IflUnit	 *Interface__Value32 `json:"ifl_unit,omitempty"`
    Family	 uint32 `json:"family,omitempty"`
    InterfaceAddress	 *Interface__IpAddress `json:"interface_address,omitempty"`
    ArpAddress	 *Interface__IpAddress `json:"arp_address,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

/*
 *******************************************************
 *     Message structure for ownership query response    
 *****************************************************
 */
type Interface__ObjectOwnershipResp struct {
    Status	 uint32 `json:"status,omitempty"`
    ErrorMessage	 string `json:"error_message,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

type Interface__TimeoutInfo struct {
    Timeout	 int32 `json:"timeout,omitempty"`
}

type Interface__TimeoutResp struct {
    Status	 uint32 `json:"status,omitempty"`
    ErrorMessage	 string `json:"error_message,omitempty"`
}

/*
 *******************************************************
 *     Message structure for Public IFL query response *
 *****************************************************
 */
type Interface__PublicIflResp struct {
    Status	 uint32 `json:"status,omitempty"`
    ErrorMessage	 string `json:"error_message,omitempty"`
    ClientCtx	 int32 `json:"client_ctx,omitempty"`
    RequestId	 int64 `json:"request_id,omitempty"`
}

type Interface__InterfacesService_InterfaceCreate struct {
  Request	Interface__InterfaceConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceCreate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceCreate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceCreate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceDelete struct {
  Request	Interface__InterfaceConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceLogicalCreate struct {
  Request	Interface__InterfaceLogicalConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceLogicalCreate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceLogicalCreate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceLogicalCreate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceLogicalDelete struct {
  Request	Interface__InterfaceLogicalConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceLogicalDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceLogicalDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceLogicalDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceFamilyCreate struct {
  Request	Interface__InterfaceFamilyConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceFamilyCreate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceFamilyCreate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceFamilyCreate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceFamilyDelete struct {
  Request	Interface__InterfaceFamilyConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceFamilyDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceFamilyDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceFamilyDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceAddressCreate struct {
  Request	Interface__InterfaceAddressConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceAddressCreate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceAddressCreate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceAddressCreate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceAddressDelete struct {
  Request	Interface__InterfaceAddressConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceAddressDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceAddressDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceAddressDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceRTAddressCreate struct {
  Request	Interface__RTConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceRTAddressCreate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceRTAddressCreate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceRTAddressCreate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__InterfacesService_InterfaceRTAddressDelete struct {
  Request	Interface__RTConfig
  Reply	Interface__ConfigResp
}

func (r *Interface__InterfacesService_InterfaceRTAddressDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.InterfacesService/InterfaceRTAddressDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__InterfacesService_InterfaceRTAddressDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__AttrQueryService_InterfacesQueryAttribute struct {
  Request	Interface__AttributeRequestInfo
  Reply	Interface__AttributeResponseInfo
}

func (r *Interface__AttrQueryService_InterfacesQueryAttribute) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.AttrQueryService/InterfacesQueryAttribute"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__AttrQueryService_InterfacesQueryAttribute) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__ObjectQueryService_isObjectOwner struct {
  Request	Interface__ObjectOwnershipQuery
  Reply	Interface__ObjectOwnershipResp
}

func (r *Interface__ObjectQueryService_isObjectOwner) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.ObjectQueryService/isObjectOwner"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__ObjectQueryService_isObjectOwner) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__TimeoutService_setClientTimeout struct {
  Request	Interface__TimeoutInfo
  Reply	Interface__TimeoutResp
}

func (r *Interface__TimeoutService_setClientTimeout) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.TimeoutService/setClientTimeout"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__TimeoutService_setClientTimeout) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Interface__PublicIflService_isPublicIfl struct {
  Request	Interface__InterfaceLogicalConfig
  Reply	Interface__PublicIflResp
}

func (r *Interface__PublicIflService_isPublicIfl) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/interface.PublicIflService/isPublicIfl"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Interface__PublicIflService_isPublicIfl) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



const (
  ACL__ACL__FALSE = 0
  ACL__ACL__TRUE = 1
)

const (
  ACL__ACL__MATCH__OP__INVALID = 0
  ACL__ACL__MATCH__OP__EQUAL = 1
  ACL__ACL__MATCH__OP__NOT__EQUAL = 2
)

const (
  /*
   * Invalid policer type
   */
  ACL__ACL__POLICER__INVALID = 0
  /*
   * Single rate two color
   */
  ACL__ACL__TWO__COLOR__POLICER = 1
  /*
   * Singel rate three color
   */
  ACL__ACL__SINGLE__RATE__THREE__COLOR__POLICER = 2
  /*
   * Two rate three color
   */
  ACL__ACL__TWO__RATE__THREE__COLOR__POLICER = 3
  /*
   * Hierarchical 
   */
  ACL__ACL__HIERARCHICAL__POLICER = 4
)

const (
  ACL__ACL__POLICER__FLAG__INVALID = 0
  /*
   * The policer instance is activated for each ACE its referenced.
   */
  ACL__ACL__POLICER__FLAG__TERM__SPECIFIC = 1
  /*
   * The policer instance is activated at global ACL level.
   */
  ACL__ACL__POLICER__FLAG__FILTER__SPECIFIC = 2
)

const (
  ACL__ACL__POLICER__RATE__INVALID = 0
  /*
   * Bits per second
   */
  ACL__ACL__POLICER__RATE__BPS = 1
  /*
   * Kilobits per second
   */
  ACL__ACL__POLICER__RATE__KBPS = 2
  /*
   * Megabits per second
   */
  ACL__ACL__POLICER__RATE__MBPS = 3
  /*
   * Gigabits per second
   */
  ACL__ACL__POLICER__RATE__GBPS = 4
)

const (
  ACL__ACL__POLICER__BURST__SIZE__INVALID = 0
  /*
   * Bytes
   */
  ACL__ACL__POLICER__BURST__SIZE__BYTE = 1
  /*
   * KiloBytes
   */
  ACL__ACL__POLICER__BURST__SIZE__KBYTE = 2
  /*
   * MegaBytes
   */
  ACL__ACL__POLICER__BURST__SIZE__MBYTE = 3
  /*
   * GigaBytes
   */
  ACL__ACL__POLICER__BURST__SIZE__GBYTE = 4
)

const (
  ACL__ACL__COLOR__MODE__INVALID = 0
  /*
   *Color Blind
   */
  ACL__ACL__COLOR__MODE__COLOR__BLIND = 1
  /*
   *Color Aware
   */
  ACL__ACL__COLOR__MODE__COLOR__AWARE = 2
)

const (
  ACL__ACL__LOSS__PRIORITY__INVALID = 0
  ACL__ACL__LOSS__PRIORITY__HIGH = 1
  ACL__ACL__LOSS__PRIORITY__MEDIUM__HIGH = 2
  ACL__ACL__LOSS__PRIORITY__MEDIUM__LOW = 3
  ACL__ACL__LOSS__PRIORITY__LOW = 4
)

const (
  ACL__ACL__FORWARDING__CLASS__INVALID = 0
  ACL__ACL__FORWARDING__CLASS__ASSURED = 1
  ACL__ACL__FORWARDING__CLASS__BEST__EFFORT = 2
  ACL__ACL__FORWARDING__CLASS__EXPEDITED = 3
  ACL__ACL__FORWARDING__CLASS__NETWORK__CONTROL = 4
)

const (
  ACL__ACL__MATCH__STP__FLAG__INVALID = 0
  ACL__ACL__MATCH__STP__FLAG__BLOCKING = 1
  ACL__ACL__MATCH__STP__FLAG__FORWARDING = 2
)

const (
  /*
   * Send ICMP Administratively Prohibited message
   */
  ACL__ACL__ACTION__REJECT__ADMINISTRATIVELY__PROHIBITED = 0
  /*
   * Send ICMP Bad Host ToS message 
   */
  ACL__ACL__ACTION__REJECT__BAD__HOST__TOS = 1
  /*
   * Send ICMP Bad Network ToS message 
   */
  ACL__ACL__ACTION__REJECT__BAD__NETWORK__TOS = 2
  /*
   * Send ICMP Fragmentation Needed message 
   */
  ACL__ACL__ACTION__REJECT__FRAGMENTATION__NEEDED = 3
  /*
   * Send ICMP Host Prohibited message 
   */
  ACL__ACL__ACTION__REJECT__HOST__PROHIBITED = 4
  /*
   * Send ICMP Host Unknown message
   */
  ACL__ACL__ACTION__REJECT__HOST__UNKNOWN = 5
  /*
   * Send ICMP Host Unreachable message 
   */
  ACL__ACL__ACTION__REJECT__HOST__UNREACHABLE = 6
  /*
   * Send ICMP Network Prohibited message 
   */
  ACL__ACL__ACTION__REJECT__NETWORK__PROHIBITED = 7
  /*
   * Send ICMP Network Unknown message 
   */
  ACL__ACL__ACTION__REJECT__NETWORK__UNKNOWN = 8
  /*
   * Send ICMP Network Unreachable message 
   */
  ACL__ACL__ACTION__REJECT__NETWORK__UNREACHABLE = 9
  /*
   * Send ICMP Port Unreachable message 
   */
  ACL__ACL__ACTION__REJECT__PORT__UNREACHABLE = 10
  /*
   * Send ICMP Precedence Cutoff message 
   */
  ACL__ACL__ACTION__REJECT__PRECEDENCE__CUTOFF = 11
  /*
   * Send ICMP Precedence Violation message 
   */
  ACL__ACL__ACTION__REJECT__PRECEDENCE__VIOLATION = 12
  /*
   * Send ICMP Protocol Unreachable message 
   */
  ACL__ACL__ACTION__REJECT__PROTOCOL__UNREACHABLE = 13
  /*
   * Send ICMP Source Host Isolated message 
   */
  ACL__ACL__ACTION__REJECT__SOURCE__HOST__ISOLATED = 14
  /*
   * Send ICMP Source Route Failed message 
   */
  ACL__ACL__ACTION__REJECT__SOURCE__ROUTE__FAILED = 15
  /*
   * Send TCP Reset message 
   */
  ACL__ACL__ACTION__REJECT__TCP__RESET = 16
)

const (
  /*
   * None 
   */
  ACL__ACL__FRAGMENT__NONE = 0
  /*
   * Dont fragment flag  
   */
  ACL__ACL__DONT__FRAGMENT = 1
  /*
   * Is fragment flag
   */
  ACL__ACL__IS__FRAGMENT = 2
  /*
   * First fragment flag 
   */
  ACL__ACL__FIRST__FRAGMENT = 3
  /*
   * More last fragment flag 
   */
  ACL__ACL__LAST__FRAGMENT = 4
)

const (
  /*
   * Routine precedence 
   */
  ACL__ACL__PRECENCE__ROUTINE = 0
  /*
   * Priority precedence 
   */
  ACL__ACL__PRECENCE__PRIORITY = 1
  /*
   * Immediate precedence 
   */
  ACL__ACL__PRECENCE__IMMEDIATE = 2
  /*
   * Flash precedence 
   */
  ACL__ACL__PRECENCE__FLASH = 3
  /*
   * Flash override precedence 
   */
  ACL__ACL__PRECENCE__FLASH__OVERRIDE = 4
  /*
   * Critical ecp precedence 
   */
  ACL__ACL__PRECENCE__CRITICAL__ECP = 5
  /*
   * Internet control precedence 
   */
  ACL__ACL__PRECENCE__INTERNET__CONTROL = 6
  /*
   * Network control precedence 
   */
  ACL__ACL__PRECENCE__NET__CONTROL = 7
)

const (
  /*
   * Invalid ACE operation
   */
  ACL__ACL__ENTRY__OPERATION__INVALID = 0
  /*
   * Add a new ACE. 
   * Can be used with Add ACL, Change ACL, replace ACL API's  
   */
  ACL__ACL__ENTRY__OPERATION__ADD = 1
  /*
   * Delete a existing ace. 
   * Can be used with change ACL API 
   */
  ACL__ACL__ENTRY__OPERATION__DELETE = 2
  /*
   * Replace a existing ace. Must provide adjacency details to 
   * preserve the order of the ace. Can be used with Change ACL API  
   */
  ACL__ACL__ENTRY__OPERATION__REPLACE = 3
)

const (
  /*
   * For first ace 
   */
  ACL__ACL__ADJACENCY__NONE = 0
  /*
   * Add next to the given ace  
   */
  ACL__ACL__ADJACENCY__AFTER = 1
  /*
   * Add before the given ace 
   */
  ACL__ACL__ADJACENCY__BEFORE = 2
)

const (
  /*
   * Invalid Flex match start offset
   */
  ACL__ACL__FLEX__MATCH__OFFSET__INVALID = 0
  /*
   * Layer-3 Flex match start offset
   */
  ACL__ACL__FLEX__MATCH__OFFSET__LAYER__THREE = 1
  /*
   * Layer-4 Flex match start offset
   */
  ACL__ACL__FLEX__MATCH__OFFSET__LAYER__FOUR = 2
  /*
   * Payload Flex match start offset
   */
  ACL__ACL__FLEX__MATCH__OFFSET__PAYLOAD = 3
)

const (
  /*
   * Invalid ACL type
   */
  ACL__ACL__TYPE__INVALID = 0
  /*
   * Classic ACL type  
   */
  ACL__ACL__TYPE__CLASSIC = 1
)

const (
  /*
   * Invalid
   */
  ACL__ACL__FAMILY__INVALID = 0
  /*
   * IPv4 family  
   */
  ACL__ACL__FAMILY__INET = 1
  /*
   * IPv6 family
   */
  ACL__ACL__FAMILY__INET6 = 2
  /*
   * Ethernet Switching family
   */
  ACL__ACL__FAMILY__ES = 3
  /*
   * VPLS family
   */
  ACL__ACL__FAMILY__VPLS = 4
  /*
   * MULTISERVICE family
   */
  ACL__ACL__FAMILY__MULTISERVICE = 5
  /*
   * CCC family
   */
  ACL__ACL__FAMILY__CCC = 6
  /*
   * MPLS family
   */
  ACL__ACL__FAMILY__MPLS = 7
)

const (
  /*
   * None 
   */
  ACL__ACL__FLAGS__NONE = 0
)

const (
  ACL__ACL__BIND__DIRECTION__INVALID = 0
  /*
   * Bind on ingress  
   */
  ACL__ACL__BIND__DIRECTION__INPUT = 1
  /*
   * Bind on egress  
   */
  ACL__ACL__BIND__DIRECTION__OUTPUT = 2
)

const (
  /*
   * Success
   */
  ACL__ACL__STATUS__EOK = 0
  /*
   * The RPC was a NULL buffer
   */
  ACL__ACL__STATUS__NULL__MESSAGE = 1
  /*
   * Wrong input
   */
  ACL__ACL__STATUS__EINVALID__MESSAGE = 2
  /*
   * Server Internal error 
   */
  ACL__ACL__STATUS__EINTERNAL = 3
  /*
   * Operation not supported
   */
  ACL__ACL__STATUS__EUNSUPPORTED__OP = 4
  /*
   * Resource not available at server
   */
  ACL__ACL__STATUS__NO__RESOURCE = 5
  /*
   * Bulk Stats timeout
   */
  ACL__ACL__STATUS__BS__TIMEOUT = 6
)

const (
  /*
   * Invalid
   */
  ACL__ACL__BIND__OBJ__TYPE__INVALID = 0
  /*
   * Interface 
   */
  ACL__ACL__BIND__OBJ__TYPE__INTERFACE = 1
  /*
   * Forwarding table 
   */
  ACL__ACL__BIND__OBJ__TYPE__FWD__TABLE = 2
  /*
   * Forwarding table 
   */
  ACL__ACL__BIND__OBJ__TYPE__VLAN = 3
  /*
   * Bridge domain 
   */
  ACL__ACL__BIND__OBJ__TYPE__BRG__DOMAIN = 4
)

/*
 * A void message 
 */
type Acl__AccessListVoid struct {
    Void	 string `json:"void,omitempty"`
}

type Acl__AclPolicerTwoColor struct {
    BwUnit	 uint32 `json:"bw_unit,omitempty"`
    Bandwidth	 uint64 `json:"bandwidth,omitempty"`
    BurstUnit	 uint32 `json:"burst_unit,omitempty"`
    BurstSize	 uint64 `json:"burst_size,omitempty"`
    Lp	 uint32 `json:"lp,omitempty"`
    FcString	 string `json:"fc_string,omitempty"`
    Discard	 uint32 `json:"discard,omitempty"`
}

type Acl__AclPolicerSingleRateThreeColor struct {
    CommittedRateUnit	 uint32 `json:"committed_rate_unit,omitempty"`
    CommittedRate	 uint64 `json:"committed_rate,omitempty"`
    CommittedBurstUnit	 uint32 `json:"committed_burst_unit,omitempty"`
    CommittedBurstSize	 uint64 `json:"committed_burst_size,omitempty"`
    ExcessBurstSize	 uint64 `json:"excess_burst_size,omitempty"`
    ExcessBurstUnit	 uint32 `json:"excess_burst_unit,omitempty"`
    Discard	 uint32 `json:"discard,omitempty"`
    ColorMode	 uint32 `json:"color_mode,omitempty"`
}

type Acl__AclPolicerTwoRateThreeColor struct {
    CommittedRateUnit	 uint32 `json:"committed_rate_unit,omitempty"`
    CommittedRate	 uint64 `json:"committed_rate,omitempty"`
    CommittedBurstUnit	 uint32 `json:"committed_burst_unit,omitempty"`
    CommittedBurstSize	 uint64 `json:"committed_burst_size,omitempty"`
    ExcessRateUnit	 uint32 `json:"excess_rate_unit,omitempty"`
    ExcessRate	 uint64 `json:"excess_rate,omitempty"`
    ExcessBurstUnit	 uint32 `json:"excess_burst_unit,omitempty"`
    ExcessBurstSize	 uint64 `json:"excess_burst_size,omitempty"`
    Discard	 uint32 `json:"discard,omitempty"`
    ColorMode	 uint32 `json:"color_mode,omitempty"`
}

type Acl__AclPolicerHierarchical struct {
    AggregateRateUnit	 uint32 `json:"aggregate_rate_unit,omitempty"`
    AggregateRate	 uint64 `json:"aggregate_rate,omitempty"`
    AggregateBurstSizeUnit	 uint32 `json:"aggregate_burst_size_unit,omitempty"`
    AggregateBurstSize	 uint64 `json:"aggregate_burst_size,omitempty"`
    PremiumRateUnit	 uint32 `json:"premium_rate_unit,omitempty"`
    PremiumRate	 uint64 `json:"premium_rate,omitempty"`
    PremiumBurstSizeUnit	 uint32 `json:"premium_burst_size_unit,omitempty"`
    PremiumBurstSize	 uint64 `json:"premium_burst_size,omitempty"`
    Discard	 uint32 `json:"discard,omitempty"`
}

type Acl__AclPolicerParameter struct {
  /*
   * Only one of the following fields should be specified: 
   * TwoColorParameter, SrThreeColorParameter, TrThreeColorParameter, HierarchicalParameter
   */
    TwoColorParameter	 *Acl__AclPolicerTwoColor `json:"two_color_parameter,omitempty"`
    SrThreeColorParameter	 *Acl__AclPolicerSingleRateThreeColor `json:"sr_three_color_parameter,omitempty"`
    TrThreeColorParameter	 *Acl__AclPolicerTwoRateThreeColor `json:"tr_three_color_parameter,omitempty"`
    HierarchicalParameter	 *Acl__AclPolicerHierarchical `json:"hierarchical_parameter,omitempty"`
}

type Acl__AccessListPolicer struct {
    PolicerName	 string `json:"policer_name,omitempty"`
    PolicerType	 uint32 `json:"policer_type,omitempty"`
    PolicerFlag	 uint32 `json:"policer_flag,omitempty"`
    PolicerParams	 *Acl__AclPolicerParameter `json:"policer_params,omitempty"`
}

type Acl__AclMatchIpAddress struct {
    Addr	 *JnxBase__IpAddress `json:"addr,omitempty"`
    PrefixLen	 uint32 `json:"prefix_len,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchPort struct {
    Min	 int32 `json:"min,omitempty"`
    Max	 int32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchIpPrecedence struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchTrafficClass struct {
    Min	 int32 `json:"min,omitempty"`
    Max	 int32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchDscpCode struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchProtocol struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchIcmpType struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchIcmpCode struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchPktLen struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

/*
 *  
 * Loss Priority match condition
 */
type Acl__AclMatchLossPriority struct {
    Lp	 uint32 `json:"lp,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

/*
 *  
 * Forwarding class match condition
 */
type Acl__AclMatchForwardingClass struct {
    FwdClass	 uint32 `json:"fwd_class,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchNextHeader struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchMplsLabel struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchVlanId struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchMeshGroup struct {
    MeshGroupId	 uint32 `json:"mesh_group_id,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchL2Token struct {
    Token	 uint32 `json:"token,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchCfmLevel struct {
    CfmLevel	 uint32 `json:"cfm_level,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchCfmOpcode struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchTtl struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchFragmentOffset struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclActionPolicer struct {
    Policer	 *Acl__AccessListPolicer `json:"policer,omitempty"`
}

type Acl__AclActionCounter struct {
    CounterName	 string `json:"counter_name,omitempty"`
}

type Acl__AclActionSendToClient struct {
    ClientName	 string `json:"client_name,omitempty"`
}

type Acl__AclActionCopyToHost struct {
    ClientName	 string `json:"client_name,omitempty"`
}

type Acl__AclActionRoutingInstance struct {
    RtInstanceName	 string `json:"rt_instance_name,omitempty"`
}

type Acl__AclActionPolicerInstance struct {
    Policer	 *Acl__AccessListPolicer `json:"policer,omitempty"`
    PolicerInstance	 string `json:"policer_instance,omitempty"`
}

/*
 * Redirect matching packets with respect to topology mentioned
 */
type Acl__AclActionTopologyRedirect struct {
    RtInstanceName	 string `json:"rt_instance_name,omitempty"`
    TopologyName	 string `json:"topology_name,omitempty"`
}

/*
 * Nexthop action
 */
type Acl__AclActionSetNexthop struct {
    NhIdx	 uint32 `json:"nh_idx,omitempty"`
}

/*
 * action losspriority
 */
type Acl__AclActionLossPriority struct {
    Lp	 uint32 `json:"lp,omitempty"`
}

/*
 * action forwording priority
 */
type Acl__AclActionForwardingPriority struct {
    Priority	 uint32 `json:"priority,omitempty"`
}

/*
 * Action forwarding class by id  
 */
type Acl__AclActionForwardingClass struct {
    Fc	 uint32 `json:"fc,omitempty"`
}

/*
 * action set dscp
 */
type Acl__AclActionSetIpDscp struct {
    Dscp	 uint32 `json:"dscp,omitempty"`
}

/*
 * Ifl Index or name in filter action
 */
type Acl__AclActionIflNameIndex struct {
  /*
   * Only one of the following fields should be specified: 
   * IflName, IflIndex
   */
    IflName	 string `json:"ifl_name,omitempty"`
    IflIndex	 uint32 `json:"ifl_index,omitempty"`
}

/*
 * action next interface
 */
type Acl__AclActionNextInterface struct {
    RtiName	 string `json:"rti_name,omitempty"`
    Ifl	 *Acl__AclActionIflNameIndex `json:"ifl,omitempty"`
}

/*
 * action next interface
 */
type Acl__AclActionNextIp struct {
    RtiName	 string `json:"rti_name,omitempty"`
    Addr	 *JnxBase__IpAddress `json:"addr,omitempty"`
    PrefixLen	 uint32 `json:"prefix_len,omitempty"`
}

type Acl__AclAdjacency struct {
    Type	 uint32 `json:"type,omitempty"`
    AceName	 string `json:"ace_name,omitempty"`
}

/*
 * Ifl Index or name
 */
type Acl__AclMatchIflNameIndex struct {
  /*
   * Only one of the following fields should be specified: 
   * IflName, IflIndex
   */
    IflName	 string `json:"ifl_name,omitempty"`
    IflIndex	 uint32 `json:"ifl_index,omitempty"`
}

type Acl__AclMatchFlexOffset struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchFlexibleRange struct {
    StartOffset	 uint32 `json:"start_offset,omitempty"`
    BitLength	 uint32 `json:"bit_length,omitempty"`
    BitOffset	 uint32 `json:"bit_offset,omitempty"`
    ByteOffset	 uint32 `json:"byte_offset,omitempty"`
    Range	 *Acl__AclMatchFlexOffset `json:"range,omitempty"`
}

type Acl__AclMatchFlexibleOffsetRange struct {
  /*
   * Only one of the following fields should be specified: 
   * FlexRangeMatch
   */
    FlexRangeMatch	 *Acl__AclMatchFlexibleRange `json:"flex_range_match,omitempty"`
}

type Acl__AclMatchFlexibleMask struct {
    StartOffset	 uint32 `json:"start_offset,omitempty"`
    BitLength	 uint32 `json:"bit_length,omitempty"`
    BitOffset	 uint32 `json:"bit_offset,omitempty"`
    ByteOffset	 uint32 `json:"byte_offset,omitempty"`
    Mask	 uint32 `json:"mask,omitempty"`
    PrefixString	 string `json:"prefix_string,omitempty"`
}

type Acl__AclMatchFlexibleOffsetMask struct {
  /*
   * Only one of the following fields should be specified: 
   * FlexMaskMatch
   */
    FlexMaskMatch	 *Acl__AclMatchFlexibleMask `json:"flex_mask_match,omitempty"`
}

type Acl__AclMatchMacAddress struct {
    Addr	 *JnxBase__MacAddress `json:"addr,omitempty"`
    AddrLen	 uint32 `json:"addr_len,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchEtherType struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchLearnVlanId struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclMatchLearnVlanPriority struct {
    Min	 uint32 `json:"min,omitempty"`
    Max	 uint32 `json:"max,omitempty"`
    MatchOp	 uint32 `json:"match_op,omitempty"`
}

type Acl__AclActionNextHop struct {
    NhIdx	 uint32 `json:"nh_idx,omitempty"`
}

type Acl__AclEntryMatchInet struct {
    MatchDstAddrs	 []Acl__AclMatchIpAddress `json:"match_dst_addrs,omitempty"`
    MatchSrcAddrs	 []Acl__AclMatchIpAddress `json:"match_src_addrs,omitempty"`
    MatchDstPorts	 []Acl__AclMatchPort `json:"match_dst_ports,omitempty"`
    MatchSrcPorts	 []Acl__AclMatchPort `json:"match_src_ports,omitempty"`
    MatchDscpCode	 []Acl__AclMatchDscpCode `json:"match_dscp_code,omitempty"`
    MatchProtocols	 []Acl__AclMatchProtocol `json:"match_protocols,omitempty"`
    MatchIcmpType	 []Acl__AclMatchIcmpType `json:"match_icmp_type,omitempty"`
    MatchIcmpCode	 []Acl__AclMatchIcmpCode `json:"match_icmp_code,omitempty"`
    MatchPktLen	 []Acl__AclMatchPktLen `json:"match_pkt_len,omitempty"`
    MatchTtl	 []Acl__AclMatchTtl `json:"match_ttl,omitempty"`
    FragmentFlags	 uint32 `json:"fragment_flags,omitempty"`
    MatchFragOffset	 []Acl__AclMatchFragmentOffset `json:"match_frag_offset,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    MatchIpPrecedence	 []Acl__AclMatchIpPrecedence `json:"match_ip_precedence,omitempty"`
    MatchAddrs	 []Acl__AclMatchIpAddress `json:"match_addrs,omitempty"`
    MatchPorts	 []Acl__AclMatchPort `json:"match_ports,omitempty"`
    MatchFlexRange	 *Acl__AclMatchFlexibleOffsetRange `json:"match_flex_range,omitempty"`
    MatchFlexMask	 *Acl__AclMatchFlexibleOffsetMask `json:"match_flex_mask,omitempty"`
}

type Acl__AclEntryInetTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionReject, ActionRtInst
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionReject	 uint32 `json:"action_reject,omitempty"`
    ActionRtInst	 *Acl__AclActionRoutingInstance `json:"action_rt_inst,omitempty"`
}

type Acl__AclEntryInetNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionSample	 uint32 `json:"action_sample,omitempty"`
    ActionNextTerm	 uint32 `json:"action_next_term,omitempty"`
}

type Acl__AclEntryInetAction struct {
    ActionsNt	 *Acl__AclEntryInetNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryInetTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclInetEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchInet `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryInetAction `json:"actions,omitempty"`
}

type Acl__AclEntryMatchEs struct {
    MatchDstMacAddrs	 []Acl__AclMatchMacAddress `json:"match_dst_mac_addrs,omitempty"`
    MatchSrcMacAddrs	 []Acl__AclMatchMacAddress `json:"match_src_mac_addrs,omitempty"`
    MatchDstPorts	 []Acl__AclMatchPort `json:"match_dst_ports,omitempty"`
    MatchSrcPorts	 []Acl__AclMatchPort `json:"match_src_ports,omitempty"`
    MatchDscpCode	 []Acl__AclMatchDscpCode `json:"match_dscp_code,omitempty"`
    MatchProtocols	 []Acl__AclMatchProtocol `json:"match_protocols,omitempty"`
    MatchIcmpType	 []Acl__AclMatchIcmpType `json:"match_icmp_type,omitempty"`
    MatchIcmpCode	 []Acl__AclMatchIcmpCode `json:"match_icmp_code,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    MatchEtherType	 []Acl__AclMatchEtherType `json:"match_ether_type,omitempty"`
    MatchLearnVlanId	 []Acl__AclMatchLearnVlanId `json:"match_learn_vlan_id,omitempty"`
    MatchLearnVlanPriority	 []Acl__AclMatchLearnVlanPriority `json:"match_learn_vlan_priority,omitempty"`
}

type Acl__AclEntryEsTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionNh, ActionSendToHost
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionNh	 *Acl__AclActionNextHop `json:"action_nh,omitempty"`
    ActionSendToHost	 uint32 `json:"action_send_to_host,omitempty"`
}

type Acl__AclEntryEsNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionNextTerm	 uint32 `json:"action_next_term,omitempty"`
    ActionLp	 *Acl__AclActionLossPriority `json:"action_lp,omitempty"`
}

type Acl__AclEntryEsAction struct {
    ActionsNt	 *Acl__AclEntryEsNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryEsTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclEsEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchEs `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryEsAction `json:"actions,omitempty"`
}

type Acl__AclEntryMatchInet6 struct {
    MatchDstAddrs	 []Acl__AclMatchIpAddress `json:"match_dst_addrs,omitempty"`
    MatchSrcAddrs	 []Acl__AclMatchIpAddress `json:"match_src_addrs,omitempty"`
    MatchDstPorts	 []Acl__AclMatchPort `json:"match_dst_ports,omitempty"`
    MatchSrcPorts	 []Acl__AclMatchPort `json:"match_src_ports,omitempty"`
    MatchDscpCode	 []Acl__AclMatchDscpCode `json:"match_dscp_code,omitempty"`
    PayloadProtocols	 []Acl__AclMatchProtocol `json:"payload_protocols,omitempty"`
    MatchIcmpType	 []Acl__AclMatchIcmpType `json:"match_icmp_type,omitempty"`
    MatchIcmpCode	 []Acl__AclMatchIcmpCode `json:"match_icmp_code,omitempty"`
    MatchPktLen	 []Acl__AclMatchPktLen `json:"match_pkt_len,omitempty"`
    FragmentFlags	 uint32 `json:"fragment_flags,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    MatchTrafficClasses	 []Acl__AclMatchTrafficClass `json:"match_traffic_classes,omitempty"`
    MatchAddrs	 []Acl__AclMatchIpAddress `json:"match_addrs,omitempty"`
    MatchFlexRange	 *Acl__AclMatchFlexibleOffsetRange `json:"match_flex_range,omitempty"`
    MatchFlexMask	 *Acl__AclMatchFlexibleOffsetMask `json:"match_flex_mask,omitempty"`
    Ipv6NextHeaders	 []Acl__AclMatchNextHeader `json:"ipv6_next_headers,omitempty"`
    MatchLossPriority	 []Acl__AclMatchLossPriority `json:"match_loss_priority,omitempty"`
    MatchFwdClass	 []Acl__AclMatchForwardingClass `json:"match_fwd_class,omitempty"`
    MatchPorts	 []Acl__AclMatchPort `json:"match_ports,omitempty"`
}

type Acl__AclEntryInet6TerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionReject, ActionRtInst, ActionTopoRedirect, ActionSendToClient, ActionSendToHost, ActionNh
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionReject	 uint32 `json:"action_reject,omitempty"`
    ActionRtInst	 *Acl__AclActionRoutingInstance `json:"action_rt_inst,omitempty"`
    ActionTopoRedirect	 *Acl__AclActionTopologyRedirect `json:"action_topo_redirect,omitempty"`
    ActionSendToClient	 *Acl__AclActionSendToClient `json:"action_send_to_client,omitempty"`
    ActionSendToHost	 uint32 `json:"action_send_to_host,omitempty"`
    ActionNh	 *Acl__AclActionSetNexthop `json:"action_nh,omitempty"`
}

type Acl__AclEntryInet6NonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionSample	 uint32 `json:"action_sample,omitempty"`
    ActionNextTerm	 uint32 `json:"action_next_term,omitempty"`
    ActionPortMirror	 uint32 `json:"action_port_mirror,omitempty"`
    ActionLp	 *Acl__AclActionLossPriority `json:"action_lp,omitempty"`
    ActionFwdClass	 *Acl__AclActionForwardingClass `json:"action_fwd_class,omitempty"`
    ActionFwdPriority	 *Acl__AclActionForwardingPriority `json:"action_fwd_priority,omitempty"`
    ActionNextIntf	 *Acl__AclActionNextInterface `json:"action_next_intf,omitempty"`
    ActionNextIp	 *Acl__AclActionNextIp `json:"action_next_ip,omitempty"`
    ActionIpDscp	 *Acl__AclActionSetIpDscp `json:"action_ip_dscp,omitempty"`
    ActionCopyToHost	 *Acl__AclActionCopyToHost `json:"action_copy_to_host,omitempty"`
    ActionPolicerInst	 *Acl__AclActionPolicerInstance `json:"action_policer_inst,omitempty"`
}

type Acl__AclEntryInet6Action struct {
    ActionsNt	 *Acl__AclEntryInet6NonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryInet6TerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclEntryMatchVpls struct {
    MatchDstAddrs	 []Acl__AclMatchIpAddress `json:"match_dst_addrs,omitempty"`
    MatchSrcAddrs	 []Acl__AclMatchIpAddress `json:"match_src_addrs,omitempty"`
    MatchDstV6Addrs	 []Acl__AclMatchIpAddress `json:"match_dst_v6_addrs,omitempty"`
    MatchSrcV6Addrs	 []Acl__AclMatchIpAddress `json:"match_src_v6_addrs,omitempty"`
    MatchDstPorts	 []Acl__AclMatchPort `json:"match_dst_ports,omitempty"`
    MatchSrcPorts	 []Acl__AclMatchPort `json:"match_src_ports,omitempty"`
    MatchDscpCode	 []Acl__AclMatchDscpCode `json:"match_dscp_code,omitempty"`
    MatchIpProtocols	 []Acl__AclMatchProtocol `json:"match_ip_protocols,omitempty"`
    PayloadProtocols	 []Acl__AclMatchProtocol `json:"payload_protocols,omitempty"`
    MatchIcmpType	 []Acl__AclMatchIcmpType `json:"match_icmp_type,omitempty"`
    MatchIcmpCode	 []Acl__AclMatchIcmpCode `json:"match_icmp_code,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    MatchTrafficClasses	 []Acl__AclMatchTrafficClass `json:"match_traffic_classes,omitempty"`
    Ipv6NextHeaders	 []Acl__AclMatchNextHeader `json:"ipv6_next_headers,omitempty"`
    EtherTypes	 []Acl__AclMatchEtherType `json:"ether_types,omitempty"`
    MatchSrcMacs	 []Acl__AclMatchMacAddress `json:"match_src_macs,omitempty"`
    MatchDstMacs	 []Acl__AclMatchMacAddress `json:"match_dst_macs,omitempty"`
    VlanEtherTypes	 []Acl__AclMatchEtherType `json:"vlan_ether_types,omitempty"`
    LearnVlanIds	 []Acl__AclMatchVlanId `json:"learn_vlan_ids,omitempty"`
    UserVlanIds	 []Acl__AclMatchVlanId `json:"user_vlan_ids,omitempty"`
    LearnVlanPriority	 []Acl__AclMatchLearnVlanPriority `json:"learn_vlan_priority,omitempty"`
    StpState	 uint32 `json:"stp_state,omitempty"`
    MeshGroupIds	 []Acl__AclMatchMeshGroup `json:"mesh_group_ids,omitempty"`
    CfmOpcodes	 []Acl__AclMatchCfmOpcode `json:"cfm_opcodes,omitempty"`
    CfmLevels	 []Acl__AclMatchCfmLevel `json:"cfm_levels,omitempty"`
    L2Tokens	 []Acl__AclMatchL2Token `json:"l2_tokens,omitempty"`
    MatchV6Addrs	 []Acl__AclMatchIpAddress `json:"match_v6_addrs,omitempty"`
    MatchFlexRange	 *Acl__AclMatchFlexibleOffsetRange `json:"match_flex_range,omitempty"`
    MatchFlexMask	 *Acl__AclMatchFlexibleOffsetMask `json:"match_flex_mask,omitempty"`
    MatchPktLen	 []Acl__AclMatchPktLen `json:"match_pkt_len,omitempty"`
}

type Acl__AclEntryVplsTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionSendToClient, ActionSendToHost, ActionNh
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionSendToClient	 *Acl__AclActionSendToClient `json:"action_send_to_client,omitempty"`
    ActionSendToHost	 uint32 `json:"action_send_to_host,omitempty"`
    ActionNh	 *Acl__AclActionSetNexthop `json:"action_nh,omitempty"`
}

type Acl__AclEntryVplsNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionSample	 uint32 `json:"action_sample,omitempty"`
    ActionNextTerm	 uint32 `json:"action_next_term,omitempty"`
    ActionNoMacLearn	 uint32 `json:"action_no_mac_learn,omitempty"`
    ActionCopyToHost	 *Acl__AclActionCopyToHost `json:"action_copy_to_host,omitempty"`
}

type Acl__AclEntryVplsAction struct {
    ActionsNt	 *Acl__AclEntryVplsNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryVplsTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclEntryMatchCcc struct {
    MatchPktLen	 []Acl__AclMatchPktLen `json:"match_pkt_len,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    EtherTypes	 []Acl__AclMatchEtherType `json:"ether_types,omitempty"`
    MatchSrcMacs	 []Acl__AclMatchMacAddress `json:"match_src_macs,omitempty"`
    MatchDstMacs	 []Acl__AclMatchMacAddress `json:"match_dst_macs,omitempty"`
    CfmOpcodes	 []Acl__AclMatchCfmOpcode `json:"cfm_opcodes,omitempty"`
    CfmLevels	 []Acl__AclMatchCfmLevel `json:"cfm_levels,omitempty"`
    MatchFlexRange	 *Acl__AclMatchFlexibleOffsetRange `json:"match_flex_range,omitempty"`
    MatchFlexMask	 *Acl__AclMatchFlexibleOffsetMask `json:"match_flex_mask,omitempty"`
}

type Acl__AclEntryCccTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionSendToClient, ActionSendToHost
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionSendToClient	 *Acl__AclActionSendToClient `json:"action_send_to_client,omitempty"`
    ActionSendToHost	 uint32 `json:"action_send_to_host,omitempty"`
}

type Acl__AclEntryCccNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionSample	 uint32 `json:"action_sample,omitempty"`
    ActionCopyToHost	 *Acl__AclActionCopyToHost `json:"action_copy_to_host,omitempty"`
}

type Acl__AclEntryCccAction struct {
    ActionsNt	 *Acl__AclEntryCccNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryCccTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclEntryMatchMultiService struct {
    MatchDstAddrs	 []Acl__AclMatchIpAddress `json:"match_dst_addrs,omitempty"`
    MatchSrcAddrs	 []Acl__AclMatchIpAddress `json:"match_src_addrs,omitempty"`
    MatchAddrs	 []Acl__AclMatchIpAddress `json:"match_addrs,omitempty"`
    MatchDstPorts	 []Acl__AclMatchPort `json:"match_dst_ports,omitempty"`
    MatchSrcPorts	 []Acl__AclMatchPort `json:"match_src_ports,omitempty"`
    MatchIpProtocols	 []Acl__AclMatchProtocol `json:"match_ip_protocols,omitempty"`
    PayloadProtocols	 []Acl__AclMatchProtocol `json:"payload_protocols,omitempty"`
    MatchIcmpType	 []Acl__AclMatchIcmpType `json:"match_icmp_type,omitempty"`
    MatchIcmpCode	 []Acl__AclMatchIcmpCode `json:"match_icmp_code,omitempty"`
    IflNames	 []Acl__AclMatchIflNameIndex `json:"ifl_names,omitempty"`
    Ipv6NextHeaders	 []Acl__AclMatchNextHeader `json:"ipv6_next_headers,omitempty"`
    EtherTypes	 []Acl__AclMatchEtherType `json:"ether_types,omitempty"`
    MatchSrcMacs	 []Acl__AclMatchMacAddress `json:"match_src_macs,omitempty"`
    MatchDstMacs	 []Acl__AclMatchMacAddress `json:"match_dst_macs,omitempty"`
    VlanEtherTypes	 []Acl__AclMatchEtherType `json:"vlan_ether_types,omitempty"`
    StpState	 uint32 `json:"stp_state,omitempty"`
    MeshGroupIds	 []Acl__AclMatchMeshGroup `json:"mesh_group_ids,omitempty"`
    L2Tokens	 []Acl__AclMatchL2Token `json:"l2_tokens,omitempty"`
    MatchPktLen	 []Acl__AclMatchPktLen `json:"match_pkt_len,omitempty"`
}

type Acl__AclEntryMultiServiceTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard, ActionSendToClient, ActionSendToHost
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
    ActionSendToClient	 *Acl__AclActionSendToClient `json:"action_send_to_client,omitempty"`
    ActionSendToHost	 uint32 `json:"action_send_to_host,omitempty"`
}

type Acl__AclEntryMultiServiceNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionLog	 uint32 `json:"action_log,omitempty"`
    ActionSyslog	 uint32 `json:"action_syslog,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
    ActionSample	 uint32 `json:"action_sample,omitempty"`
    ActionNextTerm	 uint32 `json:"action_next_term,omitempty"`
    ActionCopyToHost	 *Acl__AclActionCopyToHost `json:"action_copy_to_host,omitempty"`
}

type Acl__AclEntryMultiServiceAction struct {
    ActionsNt	 *Acl__AclEntryMultiServiceNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryMultiServiceTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclEntryMatchMpls struct {
    MatchLabel1	 []Acl__AclMatchMplsLabel `json:"match_label1,omitempty"`
    MatchLabel2	 []Acl__AclMatchMplsLabel `json:"match_label2,omitempty"`
    MatchLabel3	 []Acl__AclMatchMplsLabel `json:"match_label3,omitempty"`
    MatchFlexRange	 *Acl__AclMatchFlexibleOffsetRange `json:"match_flex_range,omitempty"`
    MatchFlexMask	 *Acl__AclMatchFlexibleOffsetMask `json:"match_flex_mask,omitempty"`
}

type Acl__AclEntryMplsTerminatingAction struct {
  /*
   * Only one of the following fields should be specified: 
   * ActionAccept, ActionDiscard
   */
    ActionAccept	 uint32 `json:"action_accept,omitempty"`
    ActionDiscard	 uint32 `json:"action_discard,omitempty"`
}

type Acl__AclEntryMplsNonTerminatingAction struct {
    ActionCount	 *Acl__AclActionCounter `json:"action_count,omitempty"`
    ActionPolicer	 *Acl__AclActionPolicer `json:"action_policer,omitempty"`
}

type Acl__AclEntryMplsAction struct {
    ActionsNt	 *Acl__AclEntryMplsNonTerminatingAction `json:"actions_nt,omitempty"`
    ActionT	 *Acl__AclEntryMplsTerminatingAction `json:"action_t,omitempty"`
}

type Acl__AclInet6Entry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchInet6 `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryInet6Action `json:"actions,omitempty"`
}

type Acl__AclVplsEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchVpls `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryVplsAction `json:"actions,omitempty"`
}

type Acl__AclCccEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchCcc `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryCccAction `json:"actions,omitempty"`
}

type Acl__AclMultiServiceEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchMultiService `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryMultiServiceAction `json:"actions,omitempty"`
}

type Acl__AclMplsEntry struct {
    AceName	 string `json:"ace_name,omitempty"`
    AceOp	 uint32 `json:"ace_op,omitempty"`
    Adjacency	 *Acl__AclAdjacency `json:"adjacency,omitempty"`
    Matches	 *Acl__AclEntryMatchMpls `json:"matches,omitempty"`
    Actions	 *Acl__AclEntryMplsAction `json:"actions,omitempty"`
}

type Acl__AclEntry struct {
  /*
   * Only one of the following fields should be specified: 
   * InetEntry, EsEntry, Inet6Entry, VplsEntry, CccEntry, MserviceEntry, MplsEntry
   */
    InetEntry	 *Acl__AclInetEntry `json:"inet_entry,omitempty"`
    EsEntry	 *Acl__AclEsEntry `json:"es_entry,omitempty"`
    Inet6Entry	 *Acl__AclInet6Entry `json:"inet6_entry,omitempty"`
    VplsEntry	 *Acl__AclVplsEntry `json:"vpls_entry,omitempty"`
    CccEntry	 *Acl__AclCccEntry `json:"ccc_entry,omitempty"`
    MserviceEntry	 *Acl__AclMultiServiceEntry `json:"mservice_entry,omitempty"`
    MplsEntry	 *Acl__AclMplsEntry `json:"mpls_entry,omitempty"`
}

type Acl__AccessList struct {
    AclName	 string `json:"acl_name,omitempty"`
    AclType	 uint32 `json:"acl_type,omitempty"`
    AclFamily	 uint32 `json:"acl_family,omitempty"`
    AclFlag	 uint32 `json:"acl_flag,omitempty"`
    AceList	 []Acl__AclEntry `json:"ace_list,omitempty"`
}

type Acl__AccessListCounter struct {
    Acl	 *Acl__AccessList `json:"acl,omitempty"`
    CounterName	 string `json:"counter_name,omitempty"`
}

type Acl__AccessListCounterBulk struct {
    Acl	 *Acl__AccessList `json:"acl,omitempty"`
    StartingIndex	 uint32 `json:"starting_index,omitempty"`
}

type Acl__AccessListReturnStatus struct {
    Status	 uint32 `json:"status,omitempty"`
}

type Acl__AccessListCounterVal struct {
    CounterName	 string `json:"counter_name,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Bytes	 uint64 `json:"bytes,omitempty"`
    Packets	 uint64 `json:"packets,omitempty"`
}

type Acl__AclBindObjVlan struct {
    VlanName	 string `json:"vlan_name,omitempty"`
    RtbName	 string `json:"rtb_name,omitempty"`
}

type Acl__AclBindObjBridgeDomain struct {
    BdName	 string `json:"bd_name,omitempty"`
    RtbName	 string `json:"rtb_name,omitempty"`
}

type Acl__AccessListBindObjPoint struct {
  /*
   * Only one of the following fields should be specified: 
   * Intf, FwdTable, Vlan, Bd
   */
    Intf	 string `json:"intf,omitempty"`
    FwdTable	 string `json:"fwd_table,omitempty"`
    Vlan	 *Acl__AclBindObjVlan `json:"vlan,omitempty"`
    Bd	 *Acl__AclBindObjBridgeDomain `json:"bd,omitempty"`
}

type Acl__AccessListObjBind struct {
    Acl	 *Acl__AccessList `json:"acl,omitempty"`
    ObjType	 uint32 `json:"obj_type,omitempty"`
    BindObject	 *Acl__AccessListBindObjPoint `json:"bind_object,omitempty"`
    BindDirection	 uint32 `json:"bind_direction,omitempty"`
    BindFamily	 uint32 `json:"bind_family,omitempty"`
}

type Acl__AclService_AccessListAdd struct {
  Request	Acl__AccessList
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListDelete struct {
  Request	Acl__AccessList
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListChange struct {
  Request	Acl__AccessList
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListChange) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListChange"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListChange) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListBindAdd struct {
  Request	Acl__AccessListObjBind
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListBindAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListBindAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListBindAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListBindDelete struct {
  Request	Acl__AccessListObjBind
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListBindDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListBindDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListBindDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPolicerAdd struct {
  Request	Acl__AccessListPolicer
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListPolicerAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPolicerAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPolicerAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPolicerReplace struct {
  Request	Acl__AccessListPolicer
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListPolicerReplace) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPolicerReplace"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPolicerReplace) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPolicerDelete struct {
  Request	Acl__AccessListPolicer
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListPolicerDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPolicerDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPolicerDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPileupStart struct {
  Request	Acl__AccessListVoid
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListPileupStart) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPileupStart"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPileupStart) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPileupEnd struct {
  Request	Acl__AccessListVoid
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListPileupEnd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPileupEnd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPileupEnd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListCounterGet struct {
  Request	Acl__AccessListCounter
  Reply	Acl__AccessListCounterVal
}

func (r *Acl__AclService_AccessListCounterGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListCounterGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListCounterGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPolicerCounterGet struct {
  Request	Acl__AccessListCounter
  Reply	Acl__AccessListCounterVal
}

func (r *Acl__AclService_AccessListPolicerCounterGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPolicerCounterGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPolicerCounterGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListCounterClear struct {
  Request	Acl__AccessListCounter
  Reply	Acl__AccessListReturnStatus
}

func (r *Acl__AclService_AccessListCounterClear) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListCounterClear"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListCounterClear) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListCounterBulkGet struct {
  Request	Acl__AccessListCounterBulk
  Reply	Acl__AccessListCounterVal
}

func (r *Acl__AclService_AccessListCounterBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListCounterBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListCounterBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Acl__AclService_AccessListPolicerCounterBulkGet struct {
  Request	Acl__AccessListCounterBulk
  Reply	Acl__AccessListCounterVal
}

func (r *Acl__AclService_AccessListPolicerCounterBulkGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/acl.AclService/AccessListPolicerCounterBulkGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Acl__AclService_AccessListPolicerCounterBulkGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * Mode of the subscription.
 */
const (
  /*
   * Values streamed by the target (Sec. 3.5.1.5.2).
   */
  GNMI__SUBSCRIPTION_LIST__STREAM = 0
  /*
   * Values sent once-off by the target (Sec. 3.5.1.5.1).
   */
  GNMI__SUBSCRIPTION_LIST__ONCE = 1
  /*
   * Values sent in response to a poll request (Sec. 3.5.1.5.3).
   */
  GNMI__SUBSCRIPTION_LIST__POLL = 2
)

/*
 * The operation that was associated with the Path specified.
 */
const (
  GNMI__UPDATE_RESULT__INVALID = 0
  /*
   * The result relates to a delete of Path.
   */
  GNMI__UPDATE_RESULT__DELETE = 1
  /*
   * The result relates to a replace of Path.
   */
  GNMI__UPDATE_RESULT__REPLACE = 2
  /*
   * The result relates to an update of Path.
   */
  GNMI__UPDATE_RESULT__UPDATE = 3
)

/*
 * Type of elements within the data tree.
 */
const (
  /*
   * All data elements.
   */
  GNMI__GET_REQUEST__ALL = 0
  /*
   * Config (rw) only elements.
   */
  GNMI__GET_REQUEST__CONFIG = 1
  /*
   * State (ro) only elements.
   */
  GNMI__GET_REQUEST__STATE = 2
  /*
   * Data elements marked in the schema as operational. This refers to data
   * elements whose value relates to the state of processes or interactions
   * running on the device.
   */
  GNMI__GET_REQUEST__OPERATIONAL = 3
)

/*
 * Encoding defines the value encoding formats that are supported by the gNMI
 * protocol. These encodings are used by both the client (when sending Set
 * messages to modify the state of the target) and the target when serializing
 * data to be returned to the client (in both Subscribe and Get RPCs).
 * Reference: gNMI Specification Section 2.3
 */
const (
  /*
   * JSON encoded text.
   */
  GNMI__JSON = 0
  /*
   * Arbitrarily encoded bytes.
   */
  GNMI__BYTES = 1
  /*
   * Encoded according to out-of-band agreed Protobuf.
   */
  GNMI__PROTO = 2
  /*
   * ASCII text of an out-of-band agreed format.
   */
  GNMI__ASCII = 3
  /*
   * JSON encoded text as per RFC7951.
   */
  GNMI__JSON__IETF = 4
)

/*
 * SubscriptionMode is the mode of the subscription, specifying how the
 * target must return values in a subscription.
 * Reference: gNMI Specification Section 3.5.1.3
 */
const (
  /*
   * The target selects the relevant mode for each element.
   */
  GNMI__TARGET__DEFINED = 0
  /*
   * The target sends an update on element value change.
   */
  GNMI__ON__CHANGE = 1
  /*
   * The target samples values according to the interval.
   */
  GNMI__SAMPLE = 2
)

type Gnmi__PathElem__KeyEntry struct {
    Key	 string `json:"key,omitempty"`
    Value	 string `json:"value,omitempty"`
}

/*
 * PathElem encodes an element of a gNMI path, along ith any attributes (keys)
 * that may be associated with it.
 * Reference: gNMI Specification Section 2.2.2.
 */
type Gnmi__PathElem struct {
    Name	 string `json:"name,omitempty"`
    Key	 []Gnmi__PathElem__KeyEntry `json:"key,omitempty"`
}

/*
 * Path encodes a data tree path as a series of repeated strings, with
 * each element of the path representing a data tree node name and the
 * associated attributes.
 * Reference: gNMI Specification Section 2.2.2.
 */
type Gnmi__Path struct {
    Element	 []string `json:"element,omitempty"`
    Origin	 string `json:"origin,omitempty"`
    Elem	 []Gnmi__PathElem `json:"elem,omitempty"`
}

/*
 * Decimal64 is used to encode a fixed precision decimal number. The value
 * is expressed as a set of digits with the precision specifying the
 * number of digits following the decimal point in the digit set.
 */
type Gnmi__Decimal64 struct {
    Digits	 uint64 `json:"digits,omitempty"`
    Precision	 uint32 `json:"precision,omitempty"`
}

/*
 * Value encodes a data tree node's value - along with the way in which
 * the value is encoded. This message is deprecated by gNMI 0.3.0.
 * Reference: gNMI Specification Section 2.2.3.
 */
type Gnmi__Value struct {
    Value	 []byte `json:"value,omitempty"`
    Type	 uint32 `json:"type,omitempty"`
}

/*
 * TypedValue is used to encode a value being sent between the client and
 * target (originated by either entity).
 */
type Gnmi__TypedValue struct {
  /*
   * Only one of the following fields should be specified: 
   * StringVal, IntVal, UintVal, BoolVal, BytesVal, FloatVal, DecimalVal, AnyVal, JsonVal, JsonIetfVal, AsciiVal
   */
    StringVal	 string `json:"string_val,omitempty"`
    IntVal	 int64 `json:"int_val,omitempty"`
    UintVal	 uint64 `json:"uint_val,omitempty"`
    BoolVal	 bool `json:"bool_val,omitempty"`
    BytesVal	 []byte `json:"bytes_val,omitempty"`
    FloatVal	 float32 `json:"float_val,omitempty"`
    DecimalVal	 *Gnmi__Decimal64 `json:"decimal_val,omitempty"`
    AnyVal	 *Google__Protobuf__Any `json:"any_val,omitempty"`
    JsonVal	 []byte `json:"json_val,omitempty"`
    JsonIetfVal	 []byte `json:"json_ietf_val,omitempty"`
    AsciiVal	 string `json:"ascii_val,omitempty"`
}

type Gnmi__ScalarArray struct {
    Element	 []Gnmi__TypedValue `json:"element,omitempty"`
}

/*
 * Value pair.
 * Reference: gNMI Specification Section 2.1
 */
type Gnmi__Update struct {
    Path	 *Gnmi__Path `json:"path,omitempty"`
    Value	 *Gnmi__Value `json:"value,omitempty"`
    Val	 *Gnmi__TypedValue `json:"val,omitempty"`
}

/*
 * Notification is a re-usable message that is used to encode data from the
 * target to the client. A Notification carries two types of changes to the data
 * tree:
 *  - Deleted values (delete) - a set of paths that have been removed from the
 *    data tree.
 *  - Updated values (update) - a set of path-value pairs indicating the path
 *    whose value has changed in the data tree.
 * Reference: gNMI Specification Section 2.1
 */
type Gnmi__Notification struct {
    Timestamp	 int64 `json:"timestamp,omitempty"`
    Prefix	 *Gnmi__Path `json:"prefix,omitempty"`
    Alias	 string `json:"alias,omitempty"`
    Update	 []Gnmi__Update `json:"update,omitempty"`
    Delete	 []Gnmi__Path `json:"delete,omitempty"`
}

/*
 * Error message previously utilised to return errors to the client. deprecated
 * in favour of using the google/rpc/status.proto message in the RPC response.
 * Reference: gNMI Specification Section 2.5
 */
type Gnmi__Error struct {
    Code	 uint32 `json:"code,omitempty"`
    Message	 string `json:"message,omitempty"`
    Data	 *Google__Protobuf__Any `json:"data,omitempty"`
}

/*
 * Subscription is a single request within a SubscriptionList. The path
 * specified is interpreted (along with the prefix) as the elements of the data
 * tree that the client is subscribing to. The mode determines how the target
 * should trigger updates to be sent.
 * Reference: gNMI Specification Section 3.5.1.3
 */
type Gnmi__Subscription struct {
    Path	 *Gnmi__Path `json:"path,omitempty"`
    Mode	 uint32 `json:"mode,omitempty"`
    SampleInterval	 uint64 `json:"sample_interval,omitempty"`
    SuppressRedundant	 bool `json:"suppress_redundant,omitempty"`
    HeartbeatInterval	 uint64 `json:"heartbeat_interval,omitempty"`
}

/*
 * QOSMarking specifies the DSCP value to be set on transmitted telemetry
 * updates from the target.
 * Reference: gNMI Specification Section 3.5.1.2
 */
type Gnmi__QOSMarking struct {
    Marking	 uint32 `json:"marking,omitempty"`
}

/*
 * ModelData is used to describe a set of schema modules. It can be used in a
 * CapabilityResponse where a target reports the set of modules that it
 * supports, and within the SubscribeRequest and GetRequest messages to specify
 * the set of models from which data tree elements should be reported.
 * Reference: gNMI Specification Section 3.2.3
 */
type Gnmi__ModelData struct {
    Name	 string `json:"name,omitempty"`
    Organization	 string `json:"organization,omitempty"`
    Version	 string `json:"version,omitempty"`
}

/*
 * SubscriptionList is used within a Subscribe message to specify the list of
 * paths that the client wishes to subscribe to. The message consists of a
 * list of (possibly prefixed) paths, and options that relate to the
 * subscription.
 * Reference: gNMI Specification Section 3.5.1.2
 */
type Gnmi__SubscriptionList struct {
    Prefix	 *Gnmi__Path `json:"prefix,omitempty"`
    Subscription	 []Gnmi__Subscription `json:"subscription,omitempty"`
    UseAliases	 bool `json:"use_aliases,omitempty"`
    Qos	 *Gnmi__QOSMarking `json:"qos,omitempty"`
    Mode	 uint32 `json:"mode,omitempty"`
    AllowAggregation	 bool `json:"allow_aggregation,omitempty"`
    UseModels	 []Gnmi__ModelData `json:"use_models,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
}

/*
 * Alias specifies a data tree path, and an associated string which defines an
 * alias which is to be used for this path in the context of the RPC. The alias
 * is specified as a string which is prefixed with "#" to disambiguate it from
 * data tree element paths.
 * Reference: gNMI Specification Section 2.4.2
 */
type Gnmi__Alias struct {
    Path	 *Gnmi__Path `json:"path,omitempty"`
    Alias	 string `json:"alias,omitempty"`
}

/*
 * AliasList specifies a list of aliases. It is used in a SubscribeRequest for
 * a client to create a set of aliases that the target is to utilize.
 * Reference: gNMI Specification Section 3.5.1.6
 */
type Gnmi__AliasList struct {
    Alias	 []Gnmi__Alias `json:"alias,omitempty"`
}

/*
 * Poll is sent within a SubscribeRequest to trigger the device to
 * send telemetry updates for the paths that are associated with the
 * subscription.
 * Reference: gNMI Specification Section Section 3.5.1.4
 */
type Gnmi__Poll struct {
}

/*
 * SubscribeRequest is the message sent by the client to the target when
 * initiating a subscription to a set of paths within the data tree. The
 * request field must be populated and the initial message must specify a
 * SubscriptionList to initiate a subscription. The message is subsequently
 * used to define aliases or trigger polled data to be sent by the target.
 * Reference: gNMI Specification Section 3.5.1.1
 */
type Gnmi__SubscribeRequest struct {
  /*
   * Only one of the following fields should be specified: 
   * Subscribe, Poll, Aliases
   */
    Subscribe	 *Gnmi__SubscriptionList `json:"subscribe,omitempty"`
    Poll	 *Gnmi__Poll `json:"poll,omitempty"`
    Aliases	 *Gnmi__AliasList `json:"aliases,omitempty"`
}

/*
 * SubscribeResponse is the message used by the target within a Subscribe RPC.
 * The target includes a Notification message which is used to transmit values
 * of the path(s) that are associated with the subscription. The same message
 * is to indicate that the target has sent all data values once (is
 * synchronized).
 * Reference: gNMI Specification Section 3.5.1.4
 */
type Gnmi__SubscribeResponse struct {
  /*
   * Only one of the following fields should be specified: 
   * Update, SyncResponse, Error
   */
    Update	 *Gnmi__Notification `json:"update,omitempty"`
    SyncResponse	 bool `json:"sync_response,omitempty"`
    Error	 *Gnmi__Error `json:"error,omitempty"`
}

/*
 * SetRequest is sent from a client to the target to update values in the data
 * tree. Paths are either deleted by the client, or modified by means of being
 * updated, or replaced. Where a replace is used, unspecified values are
 * considered to be replaced, whereas when update is used the changes are
 * considered to be incremental. The set of changes that are specified within
 * a single SetRequest are considered to be a transaction.
 * Reference: gNMI Specification Section 3.4.1
 */
type Gnmi__SetRequest struct {
    Prefix	 *Gnmi__Path `json:"prefix,omitempty"`
    Delete	 []Gnmi__Path `json:"delete,omitempty"`
    Replace	 []Gnmi__Update `json:"replace,omitempty"`
    Update	 []Gnmi__Update `json:"update,omitempty"`
}

/*
 * UpdateResult is used within the SetResponse message to communicate the
 * result of an operation specified within a SetRequest message.
 * Reference: gNMI Specification Section 3.4.2
 */
type Gnmi__UpdateResult struct {
    Timestamp	 int64 `json:"timestamp,omitempty"`
    Path	 *Gnmi__Path `json:"path,omitempty"`
    Message	 *Gnmi__Error `json:"message,omitempty"`
    Op	 uint32 `json:"op,omitempty"`
}

/*
 * SetResponse is the response to a SetRequest, sent from the target to the
 * client. It reports the result of the modifications to the data tree that were
 * specified by the client. Errors for this RPC should be reported using the
 * https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto
 * message in the RPC return. The gnmi.Error message can be used to add additional
 * details where required.
 * Reference: gNMI Specification Section 3.4.2
 */
type Gnmi__SetResponse struct {
    Prefix	 *Gnmi__Path `json:"prefix,omitempty"`
    Response	 []Gnmi__UpdateResult `json:"response,omitempty"`
    Message	 *Gnmi__Error `json:"message,omitempty"`
    Timestamp	 int64 `json:"timestamp,omitempty"`
}

/*
 * GetRequest is sent when a client initiates a Get RPC. It is used to specify
 * the set of data elements for which the target should return a snapshot of
 * data. The use_models field specifies the set of schema modules that are to
 * be used by the target - where use_models is not specified then the target
 * must use all schema models that it has.
 * Reference: gNMI Specification Section 3.3.1
 */
type Gnmi__GetRequest struct {
    Prefix	 *Gnmi__Path `json:"prefix,omitempty"`
    Path	 []Gnmi__Path `json:"path,omitempty"`
    Type	 uint32 `json:"type,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
    UseModels	 []Gnmi__ModelData `json:"use_models,omitempty"`
}

/*
 * GetResponse is used by the target to respond to a GetRequest from a client.
 * The set of Notifications corresponds to the data values that are requested
 * by the client in the GetRequest.
 * Reference: gNMI Specification Section 3.3.2
 */
type Gnmi__GetResponse struct {
    Notification	 []Gnmi__Notification `json:"notification,omitempty"`
    Error	 *Gnmi__Error `json:"error,omitempty"`
}

/*
 * CapabilityRequest is sent by the client in the Capabilities RPC to request
 * that the target reports its capabilities.
 * Reference: gNMI Specification Section 3.2.1
 */
type Gnmi__CapabilityRequest struct {
}

/*
 * CapabilityResponse is used by the target to report its capabilities to the
 * client within the Capabilities RPC.
 * Reference: gNMI Specification Section 3.2.2
 */
type Gnmi__CapabilityResponse struct {
    SupportedModels	 []Gnmi__ModelData `json:"supported_models,omitempty"`
    SupportedEncodings	 []uint32 `json:"supported_encodings,omitempty"`
    GNMIVersion	 string `json:"gNMI_version,omitempty"`
}

type Gnmi__GNMI_Capabilities struct {
  Request	Gnmi__CapabilityRequest
  Reply	Gnmi__CapabilityResponse
}

func (r *Gnmi__GNMI_Capabilities) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/gnmi.gNMI/Capabilities"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Gnmi__GNMI_Capabilities) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Gnmi__GNMI_Get struct {
  Request	Gnmi__GetRequest
  Reply	Gnmi__GetResponse
}

func (r *Gnmi__GNMI_Get) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/gnmi.gNMI/Get"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Gnmi__GNMI_Get) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Gnmi__GNMI_Set struct {
  Request	Gnmi__SetRequest
  Reply	Gnmi__SetResponse
}

func (r *Gnmi__GNMI_Set) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/gnmi.gNMI/Set"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Gnmi__GNMI_Set) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * The format of network addresses that the server is to use when 
 * responding to client requests.
 */
const (
  /*
   ** Addreses in replies will be represented by strings 
   */
  JNX_BASE__ADDRESS__STRING = 0
  /*
   ** Addreses in replies will be represented by binary data in
   *  byte arrays 
   */
  JNX_BASE__ADDRESS__BYTES = 1
)

/*
 * Address family of a network address.
 */
const (
  /*
   ** Not specified 
   */
  JNX_BASE__AF__UNSPECIFIED = 0
  /*
   ** IPv4 address family 
   */
  JNX_BASE__AF__INET = 1
  /*
   ** IPv6 address family 
   */
  JNX_BASE__AF__INET6 = 2
  /*
   ** Ethernet MAC address family 
   */
  JNX_BASE__AF__MAC = 3
)

/*
 * An IP address, which may be either IPv4 or IPv6 and be respresented
 * by either a string or array of binary bytes. 
 */
type JnxBase__IpAddress struct {
  /*
   * Only one of the following fields should be specified: 
   * AddrString, AddrBytes
   */
    AddrString	 string `json:"addr_string,omitempty"`
    AddrBytes	 []byte `json:"addr_bytes,omitempty"`
}

/*
 * An ethernet MAC address, which may be respresented by either a 
 * string or array of binary bytes. 
 */
type JnxBase__MacAddress struct {
  /*
   * Only one of the following fields should be specified: 
   * AddrString, AddrBytes
   */
    AddrString	 string `json:"addr_string,omitempty"`
    AddrBytes	 []byte `json:"addr_bytes,omitempty"`
}



/*
 **
 * Global status codes to be returned in response messages.
 * Per-RPC specific status/error codes are to be conveyed
 * in sub-codes defined in respective API definitions.
 */
const (
  /*
   ** Indicates that the RPC executed without error 
   */
  JNX_BASE__SUCCESS = 0
  /*
   ** Indicates a failure condition that should be treated as fatal 
   */
  JNX_BASE__FAILURE = 1
)

/*
 **
 * Message containing an associated key value
 * pair
 */
type JnxBase__StrKeyStrVal struct {
    Key	 string `json:"key,omitempty"`
    Value	 string `json:"value,omitempty"`
}

/*
 **
 * Message representing timeval structure
 */
type JnxBase__TimeVal struct {
    Seconds	 uint64 `json:"seconds,omitempty"`
    Microseconds	 uint64 `json:"microseconds,omitempty"`
}

/*
 **
 * RPC execution status information
 */
type JnxBase__RpcStatus struct {
    Code	 uint32 `json:"code,omitempty"`
    Message	 string `json:"message,omitempty"`
}



/*
 **
 * Different types of return codes to be sent back to client based on the
 * operation was successful or not and if not, possibly more specific reasons
 * as to why it failed. To ensure clarity, the japi error code number is
 * kept same as liblicense wrapper error code, hence the range.
 */
const (
  /*
   * Operation was executed successfully
   */
  LICENSING__SUCCESS = 0
  /*
   * General failure, operation not executed successfully
   */
  LICENSING__OPERATION__FAILED = 31
  /*
   * Allocating memory to register client failed
   */
  LICENSING__REGISTER__MEMALLOC__FAILED = 32
  /*
   * Requesting license for feature list failed
   */
  LICENSING__REGISTER__REQUESTBULK__FAILED = 33
  /*
   * Client to unregister not found
   */
  LICENSING__CLIENT__NOT__FOUND = 34
  /*
   * Get client id from japi function failed
   */
  LICENSING__JAPI__GET__CLIENTID__FAILED = 35
  /*
   * Parameter received from client is invalid
   */
  LICENSING__CLIENT__INPUT__PARAM__INVALID = 36
  /*
   * LM operation failed due to invalid argument
   */
  LICENSING__LM__INVALID__ARGUMENT = 51
  /*
   * LM operation failed due to license not valid
   */
  LICENSING__LM__INVALID__LICENSE = 52
  /*
   * LM operation failed due to license unavailability
   */
  LICENSING__LM__NO__MORE__LICENSE = 53
  /*
   * LM operation failed due to insufficient license
   */
  LICENSING__LM__INSUFFICIENT__LICENSE = 54
  /*
   * Requested capacity or feature license not available
   */
  LICENSING__LM__LK__NOT__AVAILABLE = 55
)

/*
 **
 * Different types of validity states to be sent back to client.
 * These validity states indicate license applicability and validity for a feature.
 * as to why it failed.
 */
const (
  /*
   * No license is needed for this feature and product combination
   */
  LICENSING__LICENSE__NO__NEED = 0
  /*
   * A license is required for the feature, but not installed/available
   */
  LICENSING__LICENSE__REQUIRED = 1
  /*
   * A license is required, installed and is valid
   */
  LICENSING__LICENSE__VALID = 2
  /*
   * A license is required, but the installed license is invalid
   */
  LICENSING__LICENSE__INVALID = 3
)

/*
 **
 * License events type
 */
const (
  /*
   * Reserved value
   */
  LICENSING__INVALID = 0
  /*
   * License client connected with server event
   */
  LICENSING__SERVER__CONNECTED = 1
  /*
   * License client disconnected with server event
   */
  LICENSING__SERVER__DISCONNECTED = 2
  /*
   * License installation event
   */
  LICENSING__LK__INSTALL = 5
  /*
   * License deletion event
   */
  LICENSING__LK__DELETE = 6
  /*
   * Non-production license expired
   */
  LICENSING__LK__EXPIRY__NON__PRODUCTION = 7
  /*
   * Production license expired
   */
  LICENSING__LK__EXPIRY__PRODUCTION = 8
  /*
   * License capacity change event
   */
  LICENSING__CAPACITY__CHANGE = 9
)

/*
 **
 * Different types of licenses during remaining time query.
 * These indicates if the license found for the particular feature is
 * permanent or a type of time based.
 */
const (
  /*
   * License key type is invalid
   */
  LICENSING__LK__TYPE__INVALID = 0
  /*
   * License key type is countdown
   */
  LICENSING__LK__TYPE__COUNTDOWN = 1
  /*
   * License key type is date-based
   */
  LICENSING__LK__TYPE__DATEBASED = 2
  /*
   * License key type is permanent
   */
  LICENSING__LK__TYPE__PERMANENT = 3
)

/*
 **
 * License availability and usage details for capacity license.
 * In case of non-capacity feature license, the value will be 0 or 1 and
 * indicates if the license for feature is available/requested/used or not.
 */
type Licensing__LicenseCapacity struct {
    LkTotal	 uint64 `json:"lk_total,omitempty"`
    LkRemaining	 uint64 `json:"lk_remaining,omitempty"`
    Requested	 uint64 `json:"requested,omitempty"`
    ReportedUsage	 uint64 `json:"reported_usage,omitempty"`
}

/*
 **
 * Request message to register for license notification
 */
type Licensing__LicenseNotifySend struct {
    FeatureIds	 []uint32 `json:"feature_ids,omitempty"`
    RequestList	 []uint64 `json:"request_list,omitempty"`
}

/*
 **
 * Optional list of values with event notification
 */
type Licensing__LicenseEventDetail struct {
  /*
   * Only one of the following fields should be specified: 
   * Capacity
   */
    Capacity	 *Licensing__LicenseCapacity `json:"capacity,omitempty"`
}

/*
 **
 * Reply message for license notification event
 */
type Licensing__LicenseNotifyReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    EventType	 uint32 `json:"event_type,omitempty"`
    EventInfo	 *Licensing__LicenseEventDetail `json:"event_info,omitempty"`
}

/*
 **
 * Request message to unregister from license notification service.
 * All the previously registered license features by the current
 * client will be unregistered from the notification service.
 */
type Licensing__LicenseUnregisterSend struct {
}

/*
 **
 * Reply message to unregister from license notification service.
 */
type Licensing__LicenseUnregisterReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
}

/*
 **
 * Request message to get the license information operation
 */
type Licensing__LicenseQuerySend struct {
    FeatureId	 uint32 `json:"feature_id,omitempty"`
}

/*
 **
 * Reply message to get the license validity state
 */
type Licensing__LicenseCheckValidityReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    Validity	 uint32 `json:"validity,omitempty"`
}

/*
 **
 * Reply message to get the license information operation
 */
type Licensing__LicenseGetInfoReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    Capacity	 *Licensing__LicenseCapacity `json:"capacity,omitempty"`
}

/*
 **
 * Reply message to get the license remaining time and type
 */
type Licensing__LicenseRemainingTimeReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    RemainTime	 int32 `json:"remain_time,omitempty"`
    LkType	 uint32 `json:"lk_type,omitempty"`
    IsProductionLic	 bool `json:"is_production_lic,omitempty"`
}

/*
 **
 * Reply message to get the customer id and software-sn from license
 */
type Licensing__LicenseCidAndSsnReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    CustomerId	 string `json:"customer_id,omitempty"`
    SoftwareSn	 string `json:"software_sn,omitempty"`
}

/*
 **
 * Request message to send details for the license operation to perform
 */
type Licensing__LicenseOperSend struct {
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    CapacityInput	 uint64 `json:"capacity_input,omitempty"`
}

/*
 **
 * Reply message to get the response for a license operation performed
 */
type Licensing__LicenseOperReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
    CapacityInput	 uint64 `json:"capacity_input,omitempty"`
    CapacityOutput	 uint64 `json:"capacity_output,omitempty"`
}

/*
 * Request message to send details for the license operation to perform
 */
type Licensing__TestLicenseKeySend struct {
    LicenseKey	 string `json:"license_key,omitempty"`
}

/*
 * Reply message to get the response for a license operation performed
 */
type Licensing__TestLicenseKeyReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    FeatureId	 uint32 `json:"feature_id,omitempty"`
}

type Licensing__Infra_LicenseRegister struct {
  Request	Licensing__LicenseNotifySend
  Reply	Licensing__LicenseNotifyReply
}

func (r *Licensing__Infra_LicenseRegister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseRegister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseRegister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseUnregister struct {
  Request	Licensing__LicenseUnregisterSend
  Reply	Licensing__LicenseUnregisterReply
}

func (r *Licensing__Infra_LicenseUnregister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseUnregister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseUnregister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseCheckValidity struct {
  Request	Licensing__LicenseQuerySend
  Reply	Licensing__LicenseCheckValidityReply
}

func (r *Licensing__Infra_LicenseCheckValidity) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseCheckValidity"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseCheckValidity) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseGetInfo struct {
  Request	Licensing__LicenseQuerySend
  Reply	Licensing__LicenseGetInfoReply
}

func (r *Licensing__Infra_LicenseGetInfo) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseGetInfo"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseGetInfo) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseQueryRequest struct {
  Request	Licensing__LicenseQuerySend
  Reply	Licensing__LicenseGetInfoReply
}

func (r *Licensing__Infra_LicenseQueryRequest) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseQueryRequest"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseQueryRequest) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseRequest struct {
  Request	Licensing__LicenseOperSend
  Reply	Licensing__LicenseOperReply
}

func (r *Licensing__Infra_LicenseRequest) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseRequest"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseRequest) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseRelease struct {
  Request	Licensing__LicenseOperSend
  Reply	Licensing__LicenseOperReply
}

func (r *Licensing__Infra_LicenseRelease) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseRelease"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseRelease) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseSetUsage struct {
  Request	Licensing__LicenseOperSend
  Reply	Licensing__LicenseOperReply
}

func (r *Licensing__Infra_LicenseSetUsage) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseSetUsage"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseSetUsage) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseGetRemainingTime struct {
  Request	Licensing__LicenseQuerySend
  Reply	Licensing__LicenseRemainingTimeReply
}

func (r *Licensing__Infra_LicenseGetRemainingTime) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseGetRemainingTime"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseGetRemainingTime) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_LicenseGetCidAndSsn struct {
  Request	Licensing__LicenseQuerySend
  Reply	Licensing__LicenseCidAndSsnReply
}

func (r *Licensing__Infra_LicenseGetCidAndSsn) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/LicenseGetCidAndSsn"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_LicenseGetCidAndSsn) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_TestLicenseKeyAdd struct {
  Request	Licensing__TestLicenseKeySend
  Reply	Licensing__TestLicenseKeyReply
}

func (r *Licensing__Infra_TestLicenseKeyAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/TestLicenseKeyAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_TestLicenseKeyAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Licensing__Infra_TestLicenseKeyDelete struct {
  Request	Licensing__TestLicenseKeySend
  Reply	Licensing__TestLicenseKeyReply
}

func (r *Licensing__Infra_TestLicenseKeyDelete) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/licensing.Infra/TestLicenseKeyDelete"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Licensing__Infra_TestLicenseKeyDelete) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



const (
  /*
   * JSON format, this is default
   */
  MANAGEMENT__OPERATION__FORMAT__JSON = 0
  /*
   * XML format
   */
  MANAGEMENT__OPERATION__FORMAT__XML = 1
  /*
   * Text CLI format
   */
  MANAGEMENT__OPERATION__FORMAT__CLI = 2
)

const (
  MANAGEMENT__SUCCESS = 0
  MANAGEMENT__NOK = 1
  MANAGEMENT__UNSUPPORTED__PATH = 2
  MANAGEMENT__INVALID__PATH = 3
  MANAGEMENT__INVALID__CONFIGURATION = 4
  MANAGEMENT__UNSUPPORTED__ENCODING = 5
)

const (
  MANAGEMENT__ENCODING__XML = 0
  MANAGEMENT__ENCODING__JSON = 1
)

const (
  MANAGEMENT__UPDATE__CONFIG = 0
  MANAGEMENT__REPLACE__CONFIG = 1
  MANAGEMENT__DELETE__CONFIG = 2
)

const (
  MANAGEMENT__CONFIG__LOAD__REPLACE = 0
  MANAGEMENT__CONFIG__LOAD__MERGE = 1
  MANAGEMENT__CONFIG__LOAD__OVERRIDE = 2
  MANAGEMENT__CONFIG__LOAD__UPDATE = 3
  MANAGEMENT__CONFIG__LOAD__SET = 4
)

const (
  MANAGEMENT__CONFIG__COMMIT__SYNCHRONIZE = 0
  MANAGEMENT__CONFIG__COMMIT = 1
)

type Management__ExecuteOpCommandRequest struct {
  /*
   * Only one of the following fields should be specified: 
   * CliCommand, XmlCommand, JsonCommand
   */
    RequestId	 uint64 `json:"request_id,omitempty"`
    CliCommand	 string `json:"cli_command,omitempty"`
    XmlCommand	 string `json:"xml_command,omitempty"`
    JsonCommand	 string `json:"json_command,omitempty"`
    OutFormat	 uint32 `json:"out_format,omitempty"`
}

type Management__ExecuteOpCommandResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Data	 string `json:"data,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Management__EphConfigRequestList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Path	 string `json:"path,omitempty"`
}

type Management__GetEphemeralConfigRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
    EphConfigRequests	 []Management__EphConfigRequestList `json:"eph_config_requests,omitempty"`
    EphInstanceName	 string `json:"eph_instance_name,omitempty"`
    MergeView	 bool `json:"merge_view,omitempty"`
}

type Management__GetEphemeralConfigResponse__ResponseList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Path	 string `json:"path,omitempty"`
    Value	 string `json:"value,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Management__GetEphemeralConfigResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Response	 []Management__GetEphemeralConfigResponse__ResponseList `json:"response,omitempty"`
}

type Management__EditEphemeralConfigRequest__ConfigOperationList struct {
  /*
   * Only one of the following fields should be specified: 
   * XmlConfig, JsonConfig
   */
    OperationId	 string `json:"operation_id,omitempty"`
    Operation	 uint32 `json:"operation,omitempty"`
    Path	 string `json:"path,omitempty"`
    XmlConfig	 string `json:"xml_config,omitempty"`
    JsonConfig	 string `json:"json_config,omitempty"`
}

type Management__EditEphemeralConfigRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    EphConfigOperations	 []Management__EditEphemeralConfigRequest__ConfigOperationList `json:"eph_config_operations,omitempty"`
    EphInstanceName	 string `json:"eph_instance_name,omitempty"`
    EnableConfigValidation	 bool `json:"enable_config_validation,omitempty"`
    LoadOnly	 bool `json:"load_only,omitempty"`
}

type Management__EditEphemeralConfigResponse__ResponseList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Management__EditEphemeralConfigResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Response	 []Management__EditEphemeralConfigResponse__ResponseList `json:"response,omitempty"`
}

type Management__ConfigCommit struct {
    CommitType	 uint32 `json:"commit_type,omitempty"`
    Comment	 string `json:"comment,omitempty"`
}

type Management__ExecuteCfgCommandRequest struct {
  /*
   * Only one of the following fields should be specified: 
   * XmlConfig, JsonConfig, TextConfig
   */
    RequestId	 uint64 `json:"request_id,omitempty"`
    XmlConfig	 string `json:"xml_config,omitempty"`
    JsonConfig	 string `json:"json_config,omitempty"`
    TextConfig	 string `json:"text_config,omitempty"`
    LoadType	 uint32 `json:"load_type,omitempty"`
    Commit	 *Management__ConfigCommit `json:"commit,omitempty"`
}

type Management__ExecuteCfgCommandResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Management__ManagementRpcApi_ExecuteOpCommand struct {
  Request	Management__ExecuteOpCommandRequest
  Reply	Management__ExecuteOpCommandResponse
}

func (r *Management__ManagementRpcApi_ExecuteOpCommand) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/management.ManagementRpcApi/ExecuteOpCommand"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Management__ManagementRpcApi_ExecuteOpCommand) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Management__ManagementRpcApi_ExecuteCfgCommand struct {
  Request	Management__ExecuteCfgCommandRequest
  Reply	Management__ExecuteCfgCommandResponse
}

func (r *Management__ManagementRpcApi_ExecuteCfgCommand) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/management.ManagementRpcApi/ExecuteCfgCommand"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Management__ManagementRpcApi_ExecuteCfgCommand) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Management__ManagementRpcApi_GetEphemeralConfig struct {
  Request	Management__GetEphemeralConfigRequest
  Reply	Management__GetEphemeralConfigResponse
}

func (r *Management__ManagementRpcApi_GetEphemeralConfig) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/management.ManagementRpcApi/GetEphemeralConfig"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Management__ManagementRpcApi_GetEphemeralConfig) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Management__ManagementRpcApi_EditEphemeralConfig struct {
  Request	Management__EditEphemeralConfigRequest
  Reply	Management__EditEphemeralConfigResponse
}

func (r *Management__ManagementRpcApi_EditEphemeralConfig) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/management.ManagementRpcApi/EditEphemeralConfig"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Management__ManagementRpcApi_EditEphemeralConfig) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 **
 * LSP ping information retrieval associated generic status codes which are
 * applicable to all types of LSPs
 */
const (
  /*
   ** Indicates LSP get information was successfully called 
   */
  ROUTING__LSP__PING__GET__SUCCESS = 0
  /*
   ** Internal error like malloc or read/write failure occured 
   */
  ROUTING__LSP__PING__GET__INTERNAL__ERROR = 1
  /*
   ** Indicates that the input paramter is not valid 
   */
  ROUTING__LSP__PING__GET__INVALID__PARAMETER = 2
)

/*
 **
 * Return error or status codes while retrieving RSVP Information
 */
const (
  /*
   ** RSVP LSP information was successuly retrieved 
   */
  ROUTING__RSVP__LSP__PING__NOERROR = 0
  /*
   ** Requested LSP not found 
   */
  ROUTING__RSVP__LSP__PING__LSP__NOTFOUND = 1
  /*
   ** Requested LSP is not RSVP signaled 
   */
  ROUTING__RSVP__LSP__PING__NO__RSVP__INFO = 2
  /*
   ** RSVP LSP has no path information 
   */
  ROUTING__RSVP__LSP__PING__NO__PATH__INFO = 3
  /*
   ** RSVP LSP has no route information 
   */
  ROUTING__RSVP__LSP__PING__NO__ROUTE__INFO = 4
  /*
   ** RSVP LSP has no active path 
   */
  ROUTING__RSVP__LSP__PING__NO__ACTIVE__PATH = 5
  /*
   ** Requested standby path not found 
   */
  ROUTING__RSVP__LSP__PING__STANDBY__PATH__NOTFOUND = 6
  /*
   ** Record route object for requested CCC LSP not found 
   */
  ROUTING__RSVP__LSP__PING__CCC__NORRO = 7
  /*
   ** Operation not supported for P2MP LSPs 
   */
  ROUTING__RSVP__LSP__PING__P2_MP__NOSUP = 8
  /*
   ** No Egress destinations found for P2MP LSP 
   */
  ROUTING__RSVP__LSP__PING__P2_MP__NO__EGRESS = 9
  /*
   ** No flood nexthop found for P2MP LSP 
   */
  ROUTING__RSVP__LSP__PING__P2_MP__NO__FLOOD__NHOP = 10
  /*
   ** IFL-name needs to be specified for manual bypass 
   */
  ROUTING__RSVP__LSP__PING__BYPASS__NEED__IINTERFACE__NAME = 11
)

/*
 **
 * Success status codes while retrieving RSVP Information
 */
const (
  /*
   ** RSVP LSP path is primary 
   */
  ROUTING__RSVP__LSP__PING__LSP__PRIMARY = 0
  /*
   ** RSVP LSP path is secondary 
   */
  ROUTING__RSVP__LSP__PING__LSP__SECONDARY = 1
  /*
   ** RSVP LSP is a bypass LSP 
   */
  ROUTING__RSVP__LSP__PING__LSP__BYPASS = 2
)

/*
 ** Return code for LDP LSP ping information request operations 
 */
const (
  /*
   ** LDP FEC information is successfully retrieved 
   */
  ROUTING__LDP__LSP__PING__NOERROR = 0
  /*
   ** Requested LDP FEC LSP is not found 
   */
  ROUTING__LDP__LSP__PING__FEC__NOTFOUND = 1
  /*
   ** Requested LDP P2MP FEC LSP is not found 
   */
  ROUTING__LDP__LSP__PING__P2_MP__FEC__NOTFOUND = 2
  /*
   ** Requested routing instance is not found 
   */
  ROUTING__LDP__LSP__PING__INSTANCE__NOTFOUND = 3
)

/*
 ** Return code for VPN LSP ping information request operations 
 */
const (
  /*
   **
   * VPN Route distinguisher and nexthop information were 
   * successfully retrieved
   */
  ROUTING__VPN__LSP__PING__NOERROR = 0
  /*
   ** VPN is not found 
   */
  ROUTING__VPN__LSP__PING__NOTFOUND = 1
  /*
   ** Requested prefix not found in this VPN's table 
   */
  ROUTING__VPN__LSP__PING__PFX__NOTFOUND = 2
  /*
   ** No nexthop information available for this prefix 
   */
  ROUTING__VPN__LSP__PING__NH__NOTFOUND = 3
  /*
   ** This prefix was not learnt from a remote site 
   */
  ROUTING__VPN__LSP__PING__NH__NOT__REMOTE = 4
  /*
   ** The nexthop for this prefix is not resolved 
   */
  ROUTING__VPN__LSP__PING__NH__NOT__RESOLVED = 5
  /*
   **
   * During graceful restart, cannot retrieve all the information necessary
   * for issuing a L3VPN LSP ping
   */
  ROUTING__VPN__LSP__PING__IN__GRACEFUL__RESTART = 6
)

/*
 **
 * RSVP LSP associated flags for manual and dynamic bypass.
 */
type Routing__RsvpLspFlags struct {
    RsvpDynamicBypassLsp	 bool `json:"rsvp_dynamic_bypass_lsp,omitempty"`
    RsvpManualBypassLsp	 bool `json:"rsvp_manual_bypass_lsp,omitempty"`
}

/*
 **
 * Request message to get the RSVP LSPing Info for the application to check
 * the reachability to that RSVP LSP.
 */
type Routing__RsvpLspPingInfoRequest struct {
    Flags	 *Routing__RsvpLspFlags `json:"flags,omitempty"`
    LspName	 string `json:"lsp_name,omitempty"`
    PathName	 string `json:"path_name,omitempty"`
    InterfaceName	 string `json:"interface_name,omitempty"`
    InstanceName	 string `json:"instance_name,omitempty"`
}

/*
 **
 * MPLS forwarding info, label, interface and routing table info. 
 */
type Routing__MplsForwardingInfo struct {
    NexthopAddress	 *Routing__RoutePrefix `json:"nexthop_address,omitempty"`
    NexthopRouterId	 *Routing__RoutePrefix `json:"nexthop_router_id,omitempty"`
    NexthopRouteIdx	 *Routing__RoutePrefix `json:"nexthop_route_idx,omitempty"`
    NexthopControlChannelFlag	 uint32 `json:"nexthop_control_channel_flag,omitempty"`
    NexthopLabel	 *Routing__LabelStack `json:"nexthop_label,omitempty"`
    NexthopInterfaceName	 string `json:"nexthop_interface_name,omitempty"`
}

/*
 ** RSVP LSP information required to do a LspPing 
 */
type Routing__RsvpLspPingInfo struct {
    Status	 uint32 `json:"status,omitempty"`
    SuccessStatus	 uint32 `json:"success_status,omitempty"`
    DestinationAddress	 *Routing__RoutePrefix `json:"destination_address,omitempty"`
    ExtTunnelId	 *Routing__RoutePrefix `json:"ext_tunnel_id,omitempty"`
    SenderAddress	 *Routing__RoutePrefix `json:"sender_address,omitempty"`
    TunnelId	 uint32 `json:"tunnel_id,omitempty"`
    LspId	 uint32 `json:"lsp_id,omitempty"`
    Nexthops	 uint32 `json:"nexthops,omitempty"`
    BfdDiscriminator	 uint32 `json:"bfd_discriminator,omitempty"`
    LspFromAddress	 *Routing__RoutePrefix `json:"lsp_from_address,omitempty"`
    SensorId	 uint64 `json:"sensor_id,omitempty"`
    Flags	 uint32 `json:"flags,omitempty"`
    NexthopInfo	 *Routing__MplsForwardingInfo `json:"nexthop_info,omitempty"`
}

/*
 ** Reply message for the RSVP LspPing Info request 
 */
type Routing__RsvpLspPingInfoReply struct {
    Status	 uint32 `json:"status,omitempty"`
    RsvpInfo	 *Routing__RsvpLspPingInfo `json:"rsvp_info,omitempty"`
}

/*
 **
 * Request message to get the LDP LspPing Info for the
 * application to check the reachability to that LDP Lsp.
 */
type Routing__LdpLspPingInfoRequest struct {
    Prefix	 *Routing__RoutePrefix `json:"prefix,omitempty"`
    PrefixLength	 uint32 `json:"prefix_length,omitempty"`
    InstanceName	 string `json:"instance_name,omitempty"`
}

/*
 ** LDP information required to do a LSP ping to LDP LSP 
 */
type Routing__LdpLspPingInfo struct {
    Status	 uint32 `json:"status,omitempty"`
    BfdDiscriminator	 uint32 `json:"bfd_discriminator,omitempty"`
    NexthopInfo	 *Routing__MplsForwardingInfo `json:"nexthop_info,omitempty"`
}

/*
 ** Reply message for the LDP LspPing Info request 
 */
type Routing__LdpLspPingInfoReply struct {
    Status	 uint32 `json:"status,omitempty"`
    LdpInfo	 *Routing__LdpLspPingInfo `json:"ldp_info,omitempty"`
}

/*
 **
 * Request message to get the VPN LspPing Info for the
 * application to check the reachability to that VPN route.
 */
type Routing__VpnLspPingInfoRequest struct {
    Prefix	 *Routing__RoutePrefix `json:"prefix,omitempty"`
    PrefixLength	 uint32 `json:"prefix_length,omitempty"`
    InstanceName	 string `json:"instance_name,omitempty"`
}

type Routing__VpnLspPingInfo struct {
    Status	 uint32 `json:"status,omitempty"`
    Rd	 *Routing__RouteDistinguisher `json:"rd,omitempty"`
    NexthopInfo	 *Routing__MplsForwardingInfo `json:"nexthop_info,omitempty"`
}

type Routing__VpnLspPingInfoReply struct {
    Status	 uint32 `json:"status,omitempty"`
    VpnInfo	 *Routing__VpnLspPingInfo `json:"vpn_info,omitempty"`
}

type Routing__MplsApi_LspPingGetRsvpInfo struct {
  Request	Routing__RsvpLspPingInfoRequest
  Reply	Routing__RsvpLspPingInfoReply
}

func (r *Routing__MplsApi_LspPingGetRsvpInfo) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.MplsApi/LspPingGetRsvpInfo"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__MplsApi_LspPingGetRsvpInfo) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__MplsApi_LspPingGetLdpInfo struct {
  Request	Routing__LdpLspPingInfoRequest
  Reply	Routing__LdpLspPingInfoReply
}

func (r *Routing__MplsApi_LspPingGetLdpInfo) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.MplsApi/LspPingGetLdpInfo"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__MplsApi_LspPingGetLdpInfo) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__MplsApi_LspPingGetVpnInfo struct {
  Request	Routing__VpnLspPingInfoRequest
  Reply	Routing__VpnLspPingInfoReply
}

func (r *Routing__MplsApi_LspPingGetVpnInfo) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.MplsApi/LspPingGetVpnInfo"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__MplsApi_LspPingGetVpnInfo) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



const (
  OPENCONFIG__ENCODING__XML = 0
  OPENCONFIG__ENCODING__JSON = 1
)

const (
  OPENCONFIG__GET__ALL = 0
  OPENCONFIG__GET__CONFIG = 1
  OPENCONFIG__GET__OPSTATE = 2
  OPENCONFIG__GET__OPER = 3
)

const (
  OPENCONFIG__UPDATE__CONFIG = 0
  OPENCONFIG__REPLACE__CONFIG = 1
  OPENCONFIG__DELETE__CONFIG = 2
)

const (
  OPENCONFIG__OK = 0
  OPENCONFIG__NOK = 1
  OPENCONFIG__UNSUPPORTED__PATH = 2
  OPENCONFIG__INVALID__PATH = 3
  OPENCONFIG__INVALID__CONFIGURATION = 4
  OPENCONFIG__UNSUPPORTED__INTERVAL = 5
  OPENCONFIG__INVALID__SUBSCRIPTION__ID = 6
  OPENCONFIG__UNSUPPORTED__ENCODING = 7
)

type Openconfig__GetDataEncodingsRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
}

type Openconfig__GetDataEncodingsResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Encoding	 []uint32 `json:"encoding,omitempty"`
    ResponseCode	 uint32 `json:"response_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Openconfig__SetDataEncodingRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
}

type Openconfig__SetDataEncodingResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    ResponseCode	 uint32 `json:"response_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Openconfig__GetModelsRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
}

type Openconfig__Model struct {
    Name	 string `json:"name,omitempty"`
    Namespace	 string `json:"namespace,omitempty"`
    Version	 string `json:"version,omitempty"`
}

type Openconfig__GetModelsResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Model	 []Openconfig__Model `json:"model,omitempty"`
    ResponseCode	 uint32 `json:"response_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Openconfig__GetRequestList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Operation	 uint32 `json:"operation,omitempty"`
    Path	 string `json:"path,omitempty"`
}

type Openconfig__GetRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
    GetRequest	 []Openconfig__GetRequestList `json:"get_request,omitempty"`
}

type Openconfig__GetResponse__ResponseList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Path	 string `json:"path,omitempty"`
    Value	 string `json:"value,omitempty"`
    ResponseCode	 uint32 `json:"response_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Openconfig__GetResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Response	 []Openconfig__GetResponse__ResponseList `json:"response,omitempty"`
}

type Openconfig__SetRequest__ConfigOperationList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    Operation	 uint32 `json:"operation,omitempty"`
    Path	 string `json:"path,omitempty"`
    Value	 string `json:"value,omitempty"`
}

type Openconfig__SetRequest struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Transaction	 bool `json:"transaction,omitempty"`
    Encoding	 uint32 `json:"encoding,omitempty"`
    ConfigOperation	 []Openconfig__SetRequest__ConfigOperationList `json:"config_operation,omitempty"`
}

type Openconfig__SetResponse__ResponseList struct {
    OperationId	 string `json:"operation_id,omitempty"`
    ResponseCode	 uint32 `json:"response_code,omitempty"`
    Message	 string `json:"message,omitempty"`
}

type Openconfig__SetResponse struct {
    RequestId	 uint64 `json:"request_id,omitempty"`
    Response	 []Openconfig__SetResponse__ResponseList `json:"response,omitempty"`
}

type Openconfig__OpenconfigRpcApi_GetDataEncodings struct {
  Request	Openconfig__GetDataEncodingsRequest
  Reply	Openconfig__GetDataEncodingsResponse
}

func (r *Openconfig__OpenconfigRpcApi_GetDataEncodings) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/openconfig.OpenconfigRpcApi/GetDataEncodings"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Openconfig__OpenconfigRpcApi_GetDataEncodings) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Openconfig__OpenconfigRpcApi_SetDataEncoding struct {
  Request	Openconfig__SetDataEncodingRequest
  Reply	Openconfig__SetDataEncodingResponse
}

func (r *Openconfig__OpenconfigRpcApi_SetDataEncoding) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/openconfig.OpenconfigRpcApi/SetDataEncoding"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Openconfig__OpenconfigRpcApi_SetDataEncoding) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Openconfig__OpenconfigRpcApi_GetModels struct {
  Request	Openconfig__GetModelsRequest
  Reply	Openconfig__GetModelsResponse
}

func (r *Openconfig__OpenconfigRpcApi_GetModels) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/openconfig.OpenconfigRpcApi/GetModels"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Openconfig__OpenconfigRpcApi_GetModels) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Openconfig__OpenconfigRpcApi_Get struct {
  Request	Openconfig__GetRequest
  Reply	Openconfig__GetResponse
}

func (r *Openconfig__OpenconfigRpcApi_Get) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/openconfig.OpenconfigRpcApi/Get"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Openconfig__OpenconfigRpcApi_Get) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Openconfig__OpenconfigRpcApi_Set struct {
  Request	Openconfig__SetRequest
  Reply	Openconfig__SetResponse
}

func (r *Openconfig__OpenconfigRpcApi_Set) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/openconfig.OpenconfigRpcApi/Set"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Openconfig__OpenconfigRpcApi_Set) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 ** Well-known label values defined by RFC 3032. These must only be
 *  used in label stacks in accordance with the rules of RFC 3032. 
 */
const (
  /*
   ** IPv4 Explicit NULL: valid only at bottom of stack 
   */
  ROUTING__LABEL_ENTRY__IPV4__EXPLICIT__NULL__LABEL = 0
  /*
   ** Router Alert: valid anywhere in a label stack except bottom 
   */
  ROUTING__LABEL_ENTRY__ROUTER__ALERT__LABEL = 1
  /*
   ** IPv6 Explict NULL: valid only at bottom of stack 
   */
  ROUTING__LABEL_ENTRY__IPV6__EXPLICIT__NULL__LABEL = 2
  /*
   ** Implicit NULL: See RFC 3032 
   */
  ROUTING__LABEL_ENTRY__IMPLICIT__NULL__LABEL = 3
)

/*
 **
 * ----------------------------------------------------------------------------
 * Different types of return codes to be sent back to client based on the
 * operation was successful or not and if not, possibly more specific reasons
 * as to why it failed.
 * ----------------------------------------------------------------------------
 */
const (
  /*
   * Operation was executed successfully
   */
  ROUTING__RET__SUCCESS = 0
  /*
   * General failure : operation not executed successfully
   */
  ROUTING__RET__FAILURE = 1
  /*
   * Entry was not found
   */
  ROUTING__RET__NOT__FOUND = 2
  /*
   * Invalid input paramters
   */
  ROUTING__RET__INVALID__PARAMS = 3
)

/*
 **
 * The table format allows the client to request the format that the
 * server should use to represent tables in replies sent by the server
 * to the client.
 */
const (
  /*
   ** The server will represent tables by name as strings 
   */
  ROUTING__TABLE__STRING = 0
  /*
   ** The server will represent tables by RPD table ID 
   */
  ROUTING__TABLE__ID = 1
)

/*
 **
 * Routing table destination address families.
 */
const (
  /*
   ** Unspecified 
   */
  ROUTING__RT__AF__UNSPEC = 0
  /*
   ** IPv4 destination prefix 
   */
  ROUTING__RT__AF__INET = 1
  /*
   ** IPv6 destination prefix 
   */
  ROUTING__RT__AF__INET6 = 2
  /*
   ** IPv4 VPN (L3VPN) destination prefix 
   */
  ROUTING__RT__AF__INETVPN = 3
  /*
   ** IPv6 VPN (L3VPN) destination prefix 
   */
  ROUTING__RT__AF__INET6_VPN = 4
  /*
   ** IPv4 Labeled-Unicast destination prefix. 
   */
  ROUTING__RT__AF__LABELED__INET = 5
  /*
   ** IPv6 Labeled-Unicast destination prefix. 
   */
  ROUTING__RT__AF__LABELED__INET6 = 6
  /*
   ** IPv4 SRTE destination prefix 
   */
  ROUTING__RT__AF__INET__SRTE = 7
  /*
   ** IPv6 SRTE destination prefix 
   */
  ROUTING__RT__AF__INET6__SRTE = 8
)

/*
 **
 * Label operations
 */
const (
  ROUTING__NOOP = 0
  ROUTING__POP = 1
  ROUTING__PUSH = 2
  ROUTING__SWAP = 3
)

/*
 **
 * Protocol types that define protocols that added the route.
 * RTPROTO_OTHER denotes other internal protocols currently not exposed
 * to API.
 */
const (
  /*
   ** Matches any protocol 
   */
  ROUTING__ANY = 0
  /*
   ** route is directly connected 
   */
  ROUTING__DIRECT = 1
  /*
   ** route to local side of P2P interface 
   */
  ROUTING__LOCAL = 2
  /*
   ** route was installed in kernel previously 
   */
  ROUTING__KERNEL = 3
  /*
   ** route was received via a redirect 
   */
  ROUTING__REDIRECT = 4
  /*
   ** OSPF2 AS Internal routes 
   */
  ROUTING__OSPF = 5
  /*
   ** OSPF3 AS Internal routes 
   */
  ROUTING__OSPF3 = 6
  /*
   ** OSPF AS Internal routes 
   */
  ROUTING__OSPF__ANY = 7
  /*
   ** Routing Information protocol 
   */
  ROUTING__RIP = 8
  /*
   ** Routing Information protocol for v6 
   */
  ROUTING__RIPNG = 9
  /*
   ** Border gateway protocol 
   */
  ROUTING__BGP = 10
  /*
   ** route is static 
   */
  ROUTING__STATIC = 11
  /*
   ** IS-IS 
   */
  ROUTING__ISIS = 12
  /*
   ** For IGMP stuff 
   */
  ROUTING__IGMP = 13
  /*
   ** Aggregate route 
   */
  ROUTING__AGGREGATE = 14
  /*
   ** Distance Vector Multicast Routing Protocol 
   */
  ROUTING__DVMRP = 15
  /*
   ** Protocol Independent Multicast 
   */
  ROUTING__PIM = 16
  /*
   ** Multicast Source Discovery Protocol 
   */
  ROUTING__MSDP = 17
  /*
   ** MPLS switching 
   */
  ROUTING__MPLS = 18
  /*
   ** RSVP 
   */
  ROUTING__RSVP = 19
  /*
   ** Circuit Cross-Connect 
   */
  ROUTING__CCC = 20
  /*
   ** LDP 
   */
  ROUTING__LDP = 21
  /*
   ** VPN protocol, L3 
   */
  ROUTING__VPN = 22
  /*
   ** MVPN protocol, L3 
   */
  ROUTING__MVPN = 23
  /*
   ** multicast info 
   */
  ROUTING__MCAST = 24
  /*
   ** VPN protocol, L2 
   */
  ROUTING__L2_VPN = 25
  /*
   ** l2circuit protocol 
   */
  ROUTING__L2_CKT = 26
  /*
   ** BGP Static 
   */
  ROUTING__BGP__STATIC = 27
  /*
   ** Protocols not exposed and internal to Routing backend 
   */
  ROUTING__OTHER = 28
)

/*
 **
 * Routing table (RIB) name uniquely identifying a route table,
 * formatted as a string per JUNOS convention.
 */
type Routing__RouteTableName struct {
    Name	 string `json:"name,omitempty"`
}

/*
 **
 * Routing table identifier as an integer value uniquely identifying a table.
 */
type Routing__RouteTableId struct {
    Id	 uint32 `json:"id,omitempty"`
}

/*
 **
 * Routing table (RIB), which may either be specified as a string or 
 * RPD table ID. 
 */
type Routing__RouteTable struct {
  /*
   * Only one of the following fields should be specified: 
   * RttId, RttName
   */
    RttId	 *Routing__RouteTableId `json:"rtt_id,omitempty"`
    RttName	 *Routing__RouteTableName `json:"rtt_name,omitempty"`
}

/*
 ** RFC 4364 Route type 0: 2-byte AS and assigned number 
 */
type Routing__RdType0 struct {
    AsNumber	 uint32 `json:"as_number,omitempty"`
    AssignedNumber	 uint32 `json:"assigned_number,omitempty"`
}

/*
 ** RFC 4364 Route type 1: IPv4 address and assigned number 
 */
type Routing__RdType1 struct {
    IpAddress	 *JnxBase__IpAddress `json:"ip_address,omitempty"`
    AssignedNumber	 uint32 `json:"assigned_number,omitempty"`
}

/*
 ** RFC 4364 Route type 2: 4-byte AS and assigned number 
 */
type Routing__RdType2 struct {
    AsNumber	 uint32 `json:"as_number,omitempty"`
    AssignedNumber	 uint32 `json:"assigned_number,omitempty"`
}

/*
 **
 * An RFC 4364 route distinguisher to distinguish customer VPN routes 
 * within the BGP NLRIs. Valid RD can be type 0, type 1, or type 2.
 */
type Routing__RouteDistinguisher struct {
  /*
   * Only one of the following fields should be specified: 
   * Rd0, Rd1, Rd2
   */
    Rd0	 *Routing__RdType0 `json:"rd0,omitempty"`
    Rd1	 *Routing__RdType1 `json:"rd1,omitempty"`
    Rd2	 *Routing__RdType2 `json:"rd2,omitempty"`
}

/*
 **
 * L3VPN route destination address prefix composed of a 
 * route distinguisher (RD) and IP address.
 */
type Routing__L3vpnAddress struct {
    Rd	 *Routing__RouteDistinguisher `json:"rd,omitempty"`
    VpnAddr	 *JnxBase__IpAddress `json:"vpn_addr,omitempty"`
}

/*
 **
 * A single MPLS Label entry as defined by RFC 3032
 */
type Routing__LabelEntry struct {
    Label	 uint32 `json:"label,omitempty"`
    TrafficClass	 uint32 `json:"traffic_class,omitempty"`
    Ttl	 uint32 `json:"ttl,omitempty"`
    BottomOfStack	 bool `json:"bottom_of_stack,omitempty"`
}

/*
 **
 * A single MPLS Label stack entry with the operation for the label entry
 */
type Routing__LabelStackEntry struct {
    Opcode	 uint32 `json:"opcode,omitempty"`
    LabelEntry	 *Routing__LabelEntry `json:"label_entry,omitempty"`
}

/*
 **
 * Holds the mpls label used to represent label address in route lookups
 */
type Routing__MplsAddress struct {
    Label	 uint32 `json:"label,omitempty"`
}

/*
 **
 * A label stack constructed according to the rules of RFC 3032.
 */
type Routing__LabelStack struct {
    Entries	 []Routing__LabelStackEntry `json:"entries,omitempty"`
}

/*
 **
 * Segment Identifier (SID). Either 4 octet MPLS SID or a 16 octet IPv6 SID
 * Currently IPv6 SID is not supported.
 */
type Routing__SidEntry struct {
  /*
   * Only one of the following fields should be specified: 
   * SidLabelEntry
   */
    SidLabelEntry	 *Routing__LabelEntry `json:"sid_label_entry,omitempty"`
}

/*
 **
 * Segment Type 1: SID only, in the form of MPLS Label
 */
type Routing__SegmentType1 struct {
    SidLabelEntry	 *Routing__LabelEntry `json:"sid_label_entry,omitempty"`
}

/*
 **
 * SR-TE Segment. Only SegemntType1 is currently supported.
 * Other types are for internal use only.
 */
type Routing__SRTESegment struct {
  /*
   * Only one of the following fields should be specified: 
   * SegmentType1
   */
    SegmentType1	 *Routing__SegmentType1 `json:"segment_type1,omitempty"`
}

/*
 **
 * SR-TE Segment List.
 */
type Routing__SRTESegmentList struct {
    Weight	 uint32 `json:"weight,omitempty"`
    Segments	 []Routing__SRTESegment `json:"segments,omitempty"`
}

/*
 **
 * SR-TE Distinguisher
 */
type Routing__SRTEDistinguisher struct {
    Distinguisher	 uint32 `json:"distinguisher,omitempty"`
}

/*
 **
 * SR-TE Color
 */
type Routing__SRTEColor struct {
    Color	 uint32 `json:"color,omitempty"`
}

/*
 **
 * SR-TE Binding SID. This is an optional parameter.
 * Note IPv6 Binding SID not supported so in SidEntry only MPLS label SID is
 * defined.
 */
type Routing__SRTEBindingSID struct {
    BindingSrId	 *Routing__SidEntry `json:"binding_sr_id,omitempty"`
}

/*
 **
 * SR-TE Address key fields
 */
type Routing__SRTEAddress struct {
    Destination	 *JnxBase__IpAddress `json:"destination,omitempty"`
    SrColor	 *Routing__SRTEColor `json:"sr_color,omitempty"`
    SrDistinguisher	 *Routing__SRTEDistinguisher `json:"sr_distinguisher,omitempty"`
}

/*
 **
 * SR-TE Route information like segments.
 */
type Routing__SRTERouteData struct {
    BindingSid	 *Routing__SRTEBindingSID `json:"binding_sid,omitempty"`
    Preference	 uint32 `json:"preference,omitempty"`
    SegmentLists	 []Routing__SRTESegmentList `json:"segment_lists,omitempty"`
}

/*
 **
 * Route data defined for each supported address family
 */
type Routing__AddressFamilySpecificData struct {
  /*
   * Only one of the following fields should be specified: 
   * SrtePolicyData
   */
    SrtePolicyData	 *Routing__SRTERouteData `json:"srte_policy_data,omitempty"`
}

/*
 **
 * Route destination prefix defined for each supported address family
 */
type Routing__RoutePrefix struct {
  /*
   * Only one of the following fields should be specified: 
   * Inet, Inet6, Inetvpn, Inet6vpn, LabeledInet, LabeledInet6, InetSrtePolicy, Inet6SrtePolicy
   */
    Inet	 *JnxBase__IpAddress `json:"inet,omitempty"`
    Inet6	 *JnxBase__IpAddress `json:"inet6,omitempty"`
    Inetvpn	 *Routing__L3vpnAddress `json:"inetvpn,omitempty"`
    Inet6vpn	 *Routing__L3vpnAddress `json:"inet6vpn,omitempty"`
    LabeledInet	 *JnxBase__IpAddress `json:"labeled_inet,omitempty"`
    LabeledInet6	 *JnxBase__IpAddress `json:"labeled_inet6,omitempty"`
    InetSrtePolicy	 *Routing__SRTEAddress `json:"inet_srte_policy,omitempty"`
    Inet6SrtePolicy	 *Routing__SRTEAddress `json:"inet6_srte_policy,omitempty"`
}

/*
 **
 * Network Address defined for each supported address family
 */
type Routing__NetworkAddress struct {
  /*
   * Only one of the following fields should be specified: 
   * Inet, Inet6, Mpls
   */
    Inet	 *JnxBase__IpAddress `json:"inet,omitempty"`
    Inet6	 *JnxBase__IpAddress `json:"inet6,omitempty"`
    Mpls	 *Routing__MplsAddress `json:"mpls,omitempty"`
}



/*
 **
 * ------------------------------------------------------------------
 * Request message to configure a purge timer for the client
 * ------------------------------------------------------------------
 */
type Routing__RtPurgeConfigRequest struct {
    Time	 uint32 `json:"time,omitempty"`
}

/*
 **
 * ------------------------------------------------------------------
 * Generic empty request message
 * ------------------------------------------------------------------
 */
type Routing__RtEmptyRequest struct {
}

/*
 **
 * ------------------------------------------------------------------
 * Reply message to get the return code for a operation performed
 * ------------------------------------------------------------------
 */
type Routing__RtOperReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
}

/*
 **
 * ------------------------------------------------------------------
 * Reply message to get purge timer for the client
 * ------------------------------------------------------------------
 */
type Routing__RtPurgeTimeGetReply struct {
    RetCode	 uint32 `json:"ret_code,omitempty"`
    Time	 uint32 `json:"time,omitempty"`
}

type Routing__Base_RoutePurgeTimeConfig struct {
  Request	Routing__RtPurgeConfigRequest
  Reply	Routing__RtOperReply
}

func (r *Routing__Base_RoutePurgeTimeConfig) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Base/RoutePurgeTimeConfig"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Base_RoutePurgeTimeConfig) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Base_RoutePurgeTimeUnconfig struct {
  Request	Routing__RtEmptyRequest
  Reply	Routing__RtOperReply
}

func (r *Routing__Base_RoutePurgeTimeUnconfig) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Base/RoutePurgeTimeUnconfig"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Base_RoutePurgeTimeUnconfig) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Base_RoutePurgeTimeGet struct {
  Request	Routing__RtEmptyRequest
  Reply	Routing__RtPurgeTimeGetReply
}

func (r *Routing__Base_RoutePurgeTimeGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Base/RoutePurgeTimeGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Base_RoutePurgeTimeGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 * The request message containing the service registration info
 */
type Registration__RegisterRequest struct {
  /*
   * Only one of the following fields should be specified: 
   * JsonInput, FileInput
   */
    JsonInput	 string `json:"json_input,omitempty"`
    FileInput	 string `json:"file_input,omitempty"`
    Target	 string `json:"target,omitempty"`
    RegisterId	 string `json:"register_id,omitempty"`
    SkipAuthentication	 bool `json:"skip_authentication,omitempty"`
    LogicalSystem	 string `json:"logical_system,omitempty"`
    RoutingInstance	 string `json:"routing_instance,omitempty"`
}

/*
 * The response message containing success or failure of the service.
 * registration request. result value of true indicates success and
 * false indicates failure. In case of false, error attibute indicates
 * the reason for failure
 */
type Registration__RegisterReply struct {
    Result	 bool `json:"result,omitempty"`
    Error	 string `json:"error,omitempty"`
}

type Registration__Register_RegisterService struct {
  Request	Registration__RegisterRequest
  Reply	Registration__RegisterReply
}

func (r *Registration__Register_RegisterService) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/registration.Register/RegisterService"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Registration__Register_RegisterService) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 **
 *	Route Preferences of the various route types
 */
const (
  /*
   ** Routes to interfaces 
   */
  ROUTING__RTPREF__DIRECT = 0
  /*
   ** Static routes 
   */
  ROUTING__RTPREF__STATIC = 5
  /*
   ** OSPF Internal route 
   */
  ROUTING__RTPREF__OSPF = 10
  /*
   ** IS-IS level 1 route 
   */
  ROUTING__RTPREF__LABELED__ISIS = 14
  /*
   ** IS-IS level 1 route 
   */
  ROUTING__RTPREF__ISIS__L1 = 15
  /*
   ** IS-IS level 2 route 
   */
  ROUTING__RTPREF__ISIS__L2 = 18
  /*
   ** Berkeley RIP 
   */
  ROUTING__RTPREF__RIP = 100
  /*
   ** Berkeley RIPng 
   */
  ROUTING__RTPREF__RIPNG = 100
  /*
   ** Internet Gatway Mgmt 
   */
  ROUTING__RTPREF__IGMP = 115
  /*
   ** OSPF External route 
   */
  ROUTING__RTPREF__OSPF__ASE = 150
  /*
   ** Border Gateway Protocol - external peer 
   */
  ROUTING__RTPREF__BGP__EXT = 170
)

/*
 **
 * Various ways to match a route for get requests
 */
const (
  ROUTING__BEST = 0
  ROUTING__EXACT = 1
  ROUTING__EXACT__OR__LONGER = 2
)

/*
 ** Possible return codes for route add/modify/update/remove operations. 
 */
const (
  /*
   ** Request successfully completed in full. 
   */
  ROUTING__SUCCESS = 0
  /*
   ** Request failed due to an internal server error. 
   */
  ROUTING__INTERNAL__ERROR = 1
  /*
   ** The route service has not been initialized 
   */
  ROUTING__NOT__INITIALIZED = 2
  /*
   ** Request did not result in any operations 
   */
  ROUTING__NO__OP = 3
  /*
   ** Request contained too many operations 
   */
  ROUTING__TOO__MANY__OPS = 4
  /*
   ** Request contained an invalid table. 
   */
  ROUTING__TABLE__INVALID = 5
  /*
   ** Request contained a table that was not ready for operations. 
   */
  ROUTING__TABLE__NOT__READY = 6
  /*
   ** Request contained an invalid destination address prefix 
   */
  ROUTING__PREFIX__INVALID = 7
  /*
   ** Request contained a destination prefix length too short for the
   *  supplied address/NLRI. 
   */
  ROUTING__PREFIX__LEN__TOO__SHORT = 8
  /*
   ** Request contained a destination prefix length too long for the
   *  supplied address/NLRI. 
   */
  ROUTING__PREFIX__LEN__TOO__LONG = 9
  /*
   ** The server did not have a valid gateway associated with the
   *  client. 
   */
  ROUTING__GATEWAY__INVALID = 10
  /*
   ** Request contained an invalid nexthop. 
   */
  ROUTING__NEXTHOP__INVALID = 11
  /*
   ** Request contained a nexthop with an invallid address. 
   */
  ROUTING__NEXTHOP__ADDRESS__INVALID = 12
  /*
   ** Request to add paths exceeding maximum ECMP paths for a
   *  destination. 
   */
  ROUTING__NEXTHOP__LIMIT__EXCEED = 13
  /*
   ** Request contains a route that is already present in the table. 
   */
  ROUTING__ROUTE__EXISTS = 14
  /*
   ** Request contains a route that is NOT present in the table. 
   */
  ROUTING__ROUTE__NOT__FOUND = 15
  /*
   ** Request contains an invalid protocol. Only PROTO_UNSPECIFID
   *  or PROTO_BGP_STATIC are allowed in route change operations.
   */
  ROUTING__PROTOCOL__INVALID = 16
  /*
   ** Request contains a route that is NOT present in the table. 
   */
  ROUTING__ROUTE__ADD__FAILED = 17
  /*
   ** The protocol daemon is not initialized and ready to accept
   *  route change operations. 
   */
  ROUTING__NOT__READY = 18
  /*
   ** Request cannot be serviced until current requests are processed. 
   */
  ROUTING__TRY__AGAIN = 19
  /*
   ** Request contains a route_count that exceeds the max of 1000 
   */
  ROUTING__ROUTE__COUNT__INVALID = 20
  /*
   ** Request contains a parameter that is not currently supported. 
   */
  ROUTING__REQUEST__UNSUPPORTED = 21
  /*
   ** Request contains a parameter that is not valid. 
   */
  ROUTING__REQUEST__INVALID = 22
  /*
   ** Interface name is not valid. 
   */
  ROUTING__INTERFACE__INVALID = 23
  /*
   ** 
   * Invalid paramsters for Route Monitor registration. This can be returned
   * if a wrong value is set in the registration or requested operation is 
   * invalid. For e.g. this error is returned when Route Monitor registration 
   * API is called with operation > REGISTER_DEL. 
   * This error will also be returned if a registration API is called for an 
   * existing registration with a modified value of monitor_reply_route_count.
   */
  ROUTING__ROUTE__MONITOR__REGISTER__OPERATION__INVALID = 24
  /*
   ** 
   * This error is returned when Route Monitor register API with delete
   * operation is called for a table which was not registered for monitor
   * using a add operation.
   */
  ROUTING__ROUTE__MONITOR__REGISTER__ENOENT = 25
  /*
   ** Route Monitor Policy invalid 
   */
  ROUTING__ROUTE__MONITOR__REGISTER__POLICY__INVALID = 26
  /*
   **
   * Route Monitor registration request has invalid monitor_reply_route_count. 
   * This error is also returned if monitor_reply_route_count is changed for
   * an existing registration.
   */
  ROUTING__ROUTE__MONITOR__REGISTER__REPLY__ROUTE__COUNT__INVALID = 27
  /*
   ** Route monitor registration for same table with same params exists 
   */
  ROUTING__ROUTE__MONITOR__REGISTER__EXISTS = 28
  /*
   ** MPLS Label value is invalid 
   */
  ROUTING__MPLS__LABEL__INVALID = 29
  /*
   ** MPLS Label stack operation(s) is invalid 
   */
  ROUTING__MPLS__ACTION__INVALID = 30
)

/*
 **
 * Route Monitor operations; matches with RPD_MSG_FLASH_REGISTER_REQUEST
 */
const (
  ROUTING__REGISTER__NONE = 0
  ROUTING__REGISTER__ADD = 1
  ROUTING__REGISTER__DELETE = 2
)

/*
 **
 * Operation type of routes replied in RouteMonitorReply; 
 * matches RPD_ROUTE_FLASH_OP*.
 */
const (
  ROUTING__ROUTE__MONITOR__ROUTE__OP__NONE = 0
  ROUTING__ROUTE__MONITOR__ROUTE__OP__ADD = 1
  ROUTING__ROUTE__MONITOR__ROUTE__OP__MODIFY = 2
  ROUTING__ROUTE__MONITOR__ROUTE__OP__DELETE = 3
  ROUTING__ROUTE__MONITOR__ROUTE__OP__NO__ADVERTISE = 4
  ROUTING__ROUTE__MONITOR__ROUTE__OP__END__OF__TABLE = 5
)

/*
 **
 * Route entry's unique fields typically used to match the route
 */
type Routing__RouteMatchFields struct {
    DestPrefix	 *Routing__NetworkAddress `json:"dest_prefix,omitempty"`
    DestPrefixLen	 uint32 `json:"dest_prefix_len,omitempty"`
    Table	 *Routing__RouteTable `json:"table,omitempty"`
    Cookie	 uint64 `json:"cookie,omitempty"`
}

/*
 **
 * Route gateway contains the parameters which are needed to forward traffic to
 * next router/host. Consists of gateway address, local address and interface.
 */
type Routing__RouteGateway struct {
    GatewayAddress	 *Routing__NetworkAddress `json:"gateway_address,omitempty"`
    InterfaceName	 string `json:"interface_name,omitempty"`
    LocalAddress	 *Routing__NetworkAddress `json:"local_address,omitempty"`
    LabelStack	 *Routing__LabelStack `json:"label_stack,omitempty"`
}

/*
 **
 * When a data traffic arrives on a router, route nexthop indicates the next
 * router(s) to which the traffic is to be forwarded. This consists of
 * list of gateways.
 */
type Routing__RouteNexthop struct {
    Gateways	 []Routing__RouteGateway `json:"gateways,omitempty"`
}

type Routing__RouteAttributeUint32 struct {
    Value	 uint32 `json:"value,omitempty"`
}

type Routing__RouteAttributes__PreferencesEntry struct {
    Key	 uint32 `json:"key,omitempty"`
    Value	 *Routing__RouteAttributeUint32 `json:"value,omitempty"`
}

type Routing__RouteAttributes__TagsEntry struct {
    Key	 uint32 `json:"key,omitempty"`
    Value	 *Routing__RouteAttributeUint32 `json:"value,omitempty"`
}

type Routing__RouteAttributes__ColorsEntry struct {
    Key	 uint32 `json:"key,omitempty"`
    Value	 *Routing__RouteAttributeUint32 `json:"value,omitempty"`
}

type Routing__RouteAttributes struct {
    Preferences	 []Routing__RouteAttributes__PreferencesEntry `json:"preferences,omitempty"`
    Tags	 []Routing__RouteAttributes__TagsEntry `json:"tags,omitempty"`
    Colors	 []Routing__RouteAttributes__ColorsEntry `json:"colors,omitempty"`
}

/*
 **
 * Route entry with route address, mask and attributes
 */
type Routing__RouteEntry struct {
    Key	 *Routing__RouteMatchFields `json:"key,omitempty"`
    Nexthop	 *Routing__RouteNexthop `json:"nexthop,omitempty"`
    Protocol	 uint32 `json:"protocol,omitempty"`
    Attributes	 *Routing__RouteAttributes `json:"attributes,omitempty"`
}

type Routing__RouteUpdateRequest struct {
    Routes	 []Routing__RouteEntry `json:"routes,omitempty"`
}

type Routing__RouteRemoveRequest struct {
    Keys	 []Routing__RouteMatchFields `json:"keys,omitempty"`
}

/*
 **
 * Route get operation request parameters.
 */
type Routing__RouteGetRequest struct {
    Key	 *Routing__RouteMatchFields `json:"key,omitempty"`
    MatchType	 uint32 `json:"match_type,omitempty"`
    ActiveOnly	 bool `json:"active_only,omitempty"`
    ReplyAddressFormat	 uint32 `json:"reply_address_format,omitempty"`
    ReplyTableFormat	 uint32 `json:"reply_table_format,omitempty"`
    RouteCount	 uint32 `json:"route_count,omitempty"`
}

/*
 **
 * Route operation reply containing the status of the operation.
 * Replies always returns the final status (either success or the first error
 * encountered) and the number of routes that were successfully processed
 * prior to any error or full completion of the request.
 */
type Routing__RouteOperReply struct {
    Status	 uint32 `json:"status,omitempty"`
    OperationsCompleted	 uint32 `json:"operations_completed,omitempty"`
}

/*
 **
 * Route get reply containing the status of the operation and the full or
 * partial set of matching routes, depending on how many reply RPCs the
 * stream of routes is split among.
 */
type Routing__RouteGetReply struct {
    Status	 uint32 `json:"status,omitempty"`
    Routes	 []Routing__RouteEntry `json:"routes,omitempty"`
}

type Routing__RouteMonitorPolicy struct {
    RtMonitorPolicy	 string `json:"rt_monitor_policy,omitempty"`
}

/*
 **
 * Flags that can be used to change the behavior of routes recevied via the
 * RouteMonitorReply. This can be like requesting End of Record. 
 * Matches RPD_MSG_FLASH_REGISTER_REQUEST*.
 */
type Routing__RouteMonitorRegFlags struct {
    RequestEor	 bool `json:"request_eor,omitempty"`
    NoEorToClient	 bool `json:"no_eor_to_client,omitempty"`
    RequestNoWithdrawal	 bool `json:"request_no_withdrawal,omitempty"`
    RequestFromEswd	 bool `json:"request_from_eswd,omitempty"`
    RequestFromMcsnoopd	 bool `json:"request_from_mcsnoopd,omitempty"`
    RequestFromVrrpd	 bool `json:"request_from_vrrpd,omitempty"`
    RequestForceReNotif	 bool `json:"request_force_re_notif,omitempty"`
}

/*
 **
 * Request message to register for route monitoring. The registration denotes
 * the routing table for which route monitoring is requested. 
 * Parameters in the registration request like monitor_policy can be set to 
 * influence which of the routes of the table are sent in the monitor reply
 * message.
 */
type Routing__RouteMonitorRegRequest struct {
    RtTblName	 *Routing__RouteTableName `json:"rt_tbl_name,omitempty"`
    MonitorOp	 uint32 `json:"monitor_op,omitempty"`
    MonitorFlag	 *Routing__RouteMonitorRegFlags `json:"monitor_flag,omitempty"`
    MonitorPolicy	 *Routing__RouteMonitorPolicy `json:"monitor_policy,omitempty"`
    MonitorCtx	 uint32 `json:"monitor_ctx,omitempty"`
    MonitorReplyRouteCount	 uint32 `json:"monitor_reply_route_count,omitempty"`
}

/*
 ** 
 * Route monitor entry which is sent to the client in the monitor reply 
 * message.
 */
type Routing__RouteMonitorEntry struct {
    MonitorRtOp	 uint32 `json:"monitor_rt_op,omitempty"`
    Route	 *Routing__RouteEntry `json:"route,omitempty"`
}

/*
 ** 
 * Reply message which contains the routes of the table registered for
 * monitoring.
 */
type Routing__RouteMonitorReply struct {
    Status	 uint32 `json:"status,omitempty"`
    MonitorCtx	 uint32 `json:"monitor_ctx,omitempty"`
    RtTblName	 *Routing__RouteTableName `json:"rt_tbl_name,omitempty"`
    MonitorRoutes	 []Routing__RouteMonitorEntry `json:"monitor_routes,omitempty"`
}

type Routing__Rib_RouteAdd struct {
  Request	Routing__RouteUpdateRequest
  Reply	Routing__RouteOperReply
}

func (r *Routing__Rib_RouteAdd) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteAdd"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteAdd) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Rib_RouteModify struct {
  Request	Routing__RouteUpdateRequest
  Reply	Routing__RouteOperReply
}

func (r *Routing__Rib_RouteModify) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteModify"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteModify) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Rib_RouteUpdate struct {
  Request	Routing__RouteUpdateRequest
  Reply	Routing__RouteOperReply
}

func (r *Routing__Rib_RouteUpdate) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteUpdate"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteUpdate) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Rib_RouteRemove struct {
  Request	Routing__RouteRemoveRequest
  Reply	Routing__RouteOperReply
}

func (r *Routing__Rib_RouteRemove) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteRemove"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteRemove) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Rib_RouteGet struct {
  Request	Routing__RouteGetRequest
  Reply	Routing__RouteGetReply
}

func (r *Routing__Rib_RouteGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__Rib_RouteMonitorRegister struct {
  Request	Routing__RouteMonitorRegRequest
  Reply	Routing__RouteMonitorReply
}

func (r *Routing__Rib_RouteMonitorRegister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.Rib/RouteMonitorRegister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__Rib_RouteMonitorRegister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}



/*
 **
 * Response status
 */
const (
  /*
   ** Request successfully completed 
   */
  ROUTING__ROUTING_INTERFACE_INITIALIZE_RESPONSE__SUCCESS__COMPLETED = 0
  /*
   **
   * Request successfully completed, with
   * a preexisting client of the same client
   * identifier being rebound.
   */
  ROUTING__ROUTING_INTERFACE_INITIALIZE_RESPONSE__SUCCESS__REBOUND = 1
  /*
   ** Request failed due to an internal error 
   */
  ROUTING__ROUTING_INTERFACE_INITIALIZE_RESPONSE__INTERNAL__ERROR = 2
  /*
   ** Service is already initialized 
   */
  ROUTING__ROUTING_INTERFACE_INITIALIZE_RESPONSE__ALREADY__INITIALIZED = 3
)

/*
 **
 * Response status
 */
const (
  /*
   ** Request completed successfully 
   */
  ROUTING__ROUTING_INTERFACE_GET_RESPONSE__SUCCESS__COMPLETED = 0
  /*
   ** Request failed due to an internal error 
   */
  ROUTING__ROUTING_INTERFACE_GET_RESPONSE__INTERNAL__ERROR = 1
  /*
   ** Service is not initialized 
   */
  ROUTING__ROUTING_INTERFACE_GET_RESPONSE__NOT__INITIALIZED = 2
  /*
   ** Interface name and index are invalid 
   */
  ROUTING__ROUTING_INTERFACE_GET_RESPONSE__INVALID__INDEX__AND__NAME = 3
  /*
   ** Interface was not found 
   */
  ROUTING__ROUTING_INTERFACE_GET_RESPONSE__NOT__FOUND = 4
)

/*
 **	
 * Response status
 */
const (
  /*
   ** Request successfully completed 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_RESPONSE__SUCCESS__COMPLETED = 0
  /*
   ** Request failed due to an internal error 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_RESPONSE__INTERNAL__ERROR = 1
  /*
   ** Service is not initialized 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_RESPONSE__NOT__INITIALIZED = 2
  /*
   ** Notification is already registered. 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_RESPONSE__ALREADY__REGISTERED = 3
)

/*
 **
 * Notification unregister response status
 */
const (
  /*
   ** Request successfully completed 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_UNREGISTER_RESPONSE__SUCCESS__COMPLETED = 0
  /*
   ** Request failed due to an internal error 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_UNREGISTER_RESPONSE__INTERNAL__ERROR = 1
  /*
   ** Service is not initialized 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_UNREGISTER_RESPONSE__NOT__INITIALIZED = 2
  /*
   ** Notification is not registered 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_UNREGISTER_RESPONSE__NOTIFICATION__NOT__REGISTERED = 3
)

/*
 **
 * Notification refresh response status
 */
const (
  /*
   ** Request successfully completed 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_REFRESH_RESPONSE__SUCCESS__COMPLETED = 0
  /*
   ** Request failed due to an internal error 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_REFRESH_RESPONSE__INTERNAL__ERROR = 1
  /*
   ** Service is not initialized 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_REFRESH_RESPONSE__NOT__INITIALIZED = 2
  /*
   ** Notification is not registered 
   */
  ROUTING__ROUTING_INTERFACE_NOTIFICATION_REFRESH_RESPONSE__NOTIFICATION__NOT__REGISTERED = 3
)

/*
 **
 * Routing interface status
 * CHANGED, UP, DOWN, and DELETED are applicable
 * to RoutingInterfaceNotificationResponse.
 * UP and DOWN are applicable to RoutingInterfaceGetResponse,
 */
const (
  /*
   ** Attribute changed 
   */
  ROUTING__RT__INTF__CHANGED = 0
  /*
   ** Up 
   */
  ROUTING__RT__INTF__UP = 1
  /*
   ** Down 
   */
  ROUTING__RT__INTF__DOWN = 2
  /*
   ** Deleted 
   */
  ROUTING__RT__INTF__DELETED = 3
)

/*
 **
 * Routing interface Address
 */
type Routing__RoutingInterfaceAddress struct {
    Address	 *Routing__NetworkAddress `json:"address,omitempty"`
    PrefixLength	 uint32 `json:"prefix_length,omitempty"`
    IsPrimary	 bool `json:"is_primary,omitempty"`
}

/*
 **
 * Routing interface entry
 * This entry contains the information of a
 * routing interface. It is used by
 * RoutingInterfaceNotificationResponse and
 * RoutingInterfaceGetResponse.
 * In a RoutingInterfaceNotificationResponse
 * notifying of a CHANGED or UP status, and in a
 * RoutingInterfaceGetResponse, all attributes are
 * applicable.
 * In a RoutingInterfaceNotificationResponse
 * notifying of a DOWN or DELETED status, only the  
 * name and status attributes are applicable.
 */
type Routing__RoutingInterfaceEntry struct {
    Name	 string `json:"name,omitempty"`
    Index	 uint32 `json:"index,omitempty"`
    Status	 uint32 `json:"status,omitempty"`
    Bandwidth	 uint64 `json:"bandwidth,omitempty"`
    Mtu	 uint32 `json:"mtu,omitempty"`
    Addresses	 []Routing__RoutingInterfaceAddress `json:"addresses,omitempty"`
}

/*
 ** 
 * Routing interface service initialize request.
 * A client sends this request to initialize the service.
 */
type Routing__RoutingInterfaceInitializeRequest struct {
}

/*
 **
 * Routing interface service initialize response.
 */
type Routing__RoutingInterfaceInitializeResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
}

/*
 **
 * Routing interface get request.
 * A client sends this request to query an individual interface. 
 */
type Routing__RoutingInterfaceGetRequest struct {
  /*
   * Only one of the following fields should be specified: 
   * Name, Index
   */
    Name	 string `json:"name,omitempty"`
    Index	 uint32 `json:"index,omitempty"`
    AddressFormat	 uint32 `json:"address_format,omitempty"`
}

/*
 **
 * Routing interface get response.
 */
type Routing__RoutingInterfaceGetResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Entry	 *Routing__RoutingInterfaceEntry `json:"entry,omitempty"`
}

/*
 **
 * Routing interface notification register request.
 * A client sends this request to register for interface
 * event notifications, which will be streamed to the
 * client via RoutingInterfaceNotificationResponses.
 */
type Routing__RoutingInterfaceNotificationRegisterRequest struct {
    AddressFormat	 uint32 `json:"address_format,omitempty"`
}

/*
 **
 * Routing interface notification response.
 */
type Routing__RoutingInterfaceNotificationResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
    Entries	 []Routing__RoutingInterfaceEntry `json:"entries,omitempty"`
}

/*
 ** 
 * Routing interface notification unregister request.
 * A client sends this request to unregister for interface
 * event notifications.
 */
type Routing__RoutingInterfaceNotificationUnregisterRequest struct {
}

/*
 **
 * Routing interface notification unregister response.
 */
type Routing__RoutingInterfaceNotificationUnregisterResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
}

/*
 ** 
 * Routing interface notification refresh request.
 * A client sends this request to receive a full flash
 * of all interfaces. RoutingInterfaceNotificationResponses
 * carrying the information of all the interfaces will
 * be streamed to the client. 
 */
type Routing__RoutingInterfaceNotificationRefreshRequest struct {
}

/*
 **
 * Routing interface notification refresh response.
 */
type Routing__RoutingInterfaceNotificationRefreshResponse struct {
    Code	 uint32 `json:"code,omitempty"`
    SubCode	 uint32 `json:"sub_code,omitempty"`
}

type Routing__RoutingInterface_RoutingInterfaceInitialize struct {
  Request	Routing__RoutingInterfaceInitializeRequest
  Reply	Routing__RoutingInterfaceInitializeResponse
}

func (r *Routing__RoutingInterface_RoutingInterfaceInitialize) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.RoutingInterface/RoutingInterfaceInitialize"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__RoutingInterface_RoutingInterfaceInitialize) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__RoutingInterface_RoutingInterfaceGet struct {
  Request	Routing__RoutingInterfaceGetRequest
  Reply	Routing__RoutingInterfaceGetResponse
}

func (r *Routing__RoutingInterface_RoutingInterfaceGet) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.RoutingInterface/RoutingInterfaceGet"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__RoutingInterface_RoutingInterfaceGet) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__RoutingInterface_RoutingInterfaceNotificationRegister struct {
  Request	Routing__RoutingInterfaceNotificationRegisterRequest
  Reply	Routing__RoutingInterfaceNotificationResponse
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationRegister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.RoutingInterface/RoutingInterfaceNotificationRegister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationRegister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__RoutingInterface_RoutingInterfaceNotificationUnregister struct {
  Request	Routing__RoutingInterfaceNotificationUnregisterRequest
  Reply	Routing__RoutingInterfaceNotificationUnregisterResponse
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationUnregister) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.RoutingInterface/RoutingInterfaceNotificationUnregister"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationUnregister) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

type Routing__RoutingInterface_RoutingInterfaceNotificationRefresh struct {
  Request	Routing__RoutingInterfaceNotificationRefreshRequest
  Reply	Routing__RoutingInterfaceNotificationRefreshResponse
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationRefresh) Marshal(k RpcKey) (rpc string, key, value []byte, err error) {
  rpc = "/routing.RoutingInterface/RoutingInterfaceNotificationRefresh"

  if key, err = json.Marshal(k); err != nil {
    return "", nil, nil, err
  }

  if value, err = json.Marshal(r.Request); err != nil {
    return "", nil, nil, err
  }

  return rpc, key, value, nil
}

func (r *Routing__RoutingInterface_RoutingInterfaceNotificationRefresh) Unmarshal(value []byte) (err error) {
  if err = json.Unmarshal(value, r.Reply); err != nil {
    return err
  }

  return nil
}

