/*
 * Copyright (c) 2018, Juniper Networks, Inc.
 * All rights reserved.
 */

package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/srikanth2212/jsonez"
)

const (
	PACKAGETABLE string = "packageTable"
	MESSAGETABLE string = "messageTable"
	RPCTABLE     string = "rpcTable"
)

const (
	FIELD_FLAGS    string = "flags"
	FIELD_NAME     string = "name"
	FIELD_LABEL    string = "label"
	FIELD_ONEOF    string = "containingOneof"
	FIELD_SUB_MSG  string = "subMessage"
	FIELD_TYPE     string = "type"
	FIELD_UNION_ID string = "unionId"
)

const (
	IN_MSG_DESC_NAME  string = "inMsgDescriptor"
	OUT_MSG_DESC_NAME string = "outMsgDescriptor"
	DAEMON_NAME       string = "dName"
	TARGET            string = "target"
	IN_STREAM_API     string = "inStreamAPI"
	OUT_STREAM_API    string = "outStreamAPI"
	IS_STREAM_API     string = "isStreamAPI"
)

type FieldDesc struct {
	Fid         string
	Fname       string
	Flabel      ProtobufLabel
	Ftype       ProtobufType
	Flags       uint32
	ParentOneof string
	SubMessage  string
	SubField    map[string]*FieldDesc
}

type FdMap map[string]*FieldDesc

type MsgDesc struct {
	FieldDescMap FdMap
}

type MdMap map[string]*MsgDesc

type RpcDesc struct {
	InMsgDescName    string
	InMsgDescriptor  *MsgDesc
	OutMsgDescName   string
	OutMsgDescriptor *MsgDesc
	DName            string // gRPC channel target daemon name for junos daemons
	Target           string // gRPC channel target for non-junos daemons
	InStreamAPI      bool
	OutStreamAPI     bool
}

type RdMap map[string]*RpcDesc

type PackageDesc struct {
	/* Boolean variable to indicate if LoginCheck should be
	 * disabled for the RPCs in this package if the client
	 * is connected via a secure channel
	 */
	SkipLoginCheck bool
	Mmap           MdMap
	Rmap           RdMap
	/*
	 * Temp values used while building the map
	 */
	temp_mMap  MdMap
	temp_rMap  RdMap
	Imports    Set
	Dependants Set
}

type PackageDescMap map[string]*PackageDesc
type protoTypeMap map[string]ProtobufType
type protoLabelMap map[string]ProtobufLabel
type Set map[string]bool

type ProtoDescManager struct {
	Pmap     PackageDescMap
	typeMap  protoTypeMap
	labelMap protoLabelMap
	lock     sync.Mutex
}

var pMan *ProtoDescManager

/*
 * Init method of ProtoDescManager object
 * that will be used to populate the type
 * label values.
 */
func (pMan *ProtoDescManager) Init() {
	pMan.Pmap = make(PackageDescMap)
	pMan.typeMap = make(protoTypeMap)
	pMan.labelMap = make(protoLabelMap)

	pMan.typeMap["PROTOBUF_C_TYPE_INT32"] = PROTOBUF_TYPE_INT32
	pMan.typeMap["PROTOBUF_C_TYPE_SINT32"] = PROTOBUF_TYPE_SINT32
	pMan.typeMap["PROTOBUF_C_TYPE_SFIXED32"] = PROTOBUF_TYPE_SFIXED32
	pMan.typeMap["PROTOBUF_C_TYPE_INT64"] = PROTOBUF_TYPE_INT64
	pMan.typeMap["PROTOBUF_C_TYPE_SINT64"] = PROTOBUF_TYPE_SINT64
	pMan.typeMap["PROTOBUF_C_TYPE_SFIXED64"] = PROTOBUF_TYPE_SFIXED64
	pMan.typeMap["PROTOBUF_C_TYPE_UINT32"] = PROTOBUF_TYPE_UINT32
	pMan.typeMap["PROTOBUF_C_TYPE_FIXED32"] = PROTOBUF_TYPE_FIXED32
	pMan.typeMap["PROTOBUF_C_TYPE_UINT64"] = PROTOBUF_TYPE_UINT64
	pMan.typeMap["PROTOBUF_C_TYPE_FIXED64"] = PROTOBUF_TYPE_FIXED64
	pMan.typeMap["PROTOBUF_C_TYPE_FLOAT"] = PROTOBUF_TYPE_FLOAT
	pMan.typeMap["PROTOBUF_C_TYPE_DOUBLE"] = PROTOBUF_TYPE_DOUBLE
	pMan.typeMap["PROTOBUF_C_TYPE_BOOL"] = PROTOBUF_TYPE_BOOL
	pMan.typeMap["PROTOBUF_C_TYPE_ENUM"] = PROTOBUF_TYPE_ENUM
	pMan.typeMap["PROTOBUF_C_TYPE_STRING"] = PROTOBUF_TYPE_STRING
	pMan.typeMap["PROTOBUF_C_TYPE_BYTES"] = PROTOBUF_TYPE_BYTES
	pMan.typeMap["PROTOBUF_C_TYPE_MESSAGE"] = PROTOBUF_TYPE_MESSAGE

	pMan.labelMap["PROTOBUF_C_LABEL_REQUIRED"] = PROTOBUF_LABEL_REQUIRED
	pMan.labelMap["PROTOBUF_C_LABEL_OPTIONAL"] = PROTOBUF_LABEL_OPTIONAL
	pMan.labelMap["PROTOBUF_C_LABEL_REPEATED"] = PROTOBUF_LABEL_REPEATED
}

func (pMan *ProtoDescManager) ReadJSONFile(filename string) error {
	json, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	pMan.Init()

	return pMan.ParseJSON(json)
}

func (pMan *ProtoDescManager) ParseJSON(input []byte) error {
	/*
	 * Read the JSON input and build a GoJSON object
	 * tree
	 */
	var root, ptable, mtable, rtable, curPkg, curMsg, curRpc,
		imports *jsonez.GoJSON
	var pdesc *PackageDesc
	var ok bool
	var err error

	root, err = jsonez.GoJSONParse(input)

	//output := jsonez.GoJSONPrint(root)

	//fmt.Println("output is ")
	//fmt.Println(string(output))

	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	/*
	 * Fetch packageTable object
	 */
	ptable, err = root.Get(PACKAGETABLE)
	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	if ptable.Jsontype != jsonez.JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: %s is not of type JSON_ARRAY", PACKAGETABLE,
			FuncName())
		return errors.New(errorStr)
	}

	for i := 0; i < ptable.GetArraySize(); i++ {
		curPkg, err = ptable.GetArrayElemByIndex(i)

		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		curPkg = curPkg.Child

		if curPkg == nil || curPkg.Key == "" {
			errorStr := fmt.Sprintf("%s: goJSON object has empty key",
				FuncName())
			return errors.New(errorStr)
		}

		if curPkg.Jsontype != jsonez.JSON_OBJECT {
			errorStr := fmt.Sprintf("%s: goJSON object with key %s is not "+
				"JSON_OBJECT", FuncName(), curPkg.Key)
			return errors.New(errorStr)
		}

		/*
		 * Create an entry in the package map if the package
		 * is not present already
		 */
		pdesc, ok = pMan.Pmap[curPkg.Key]
		if !ok {
			//fmt.Println("Creating package descriptor for package", curPkg.Key)
			pdesc = new(PackageDesc)
			pdesc.Mmap = nil
			pdesc.Rmap = nil
			pdesc.temp_mMap = make(MdMap)
			pdesc.temp_rMap = make(RdMap)
			pdesc.Imports = make(Set)
			pdesc.Dependants = make(Set)
			pMan.Pmap[curPkg.Key] = pdesc
		} else {
			log.Println("Package ", curPkg.Key, " exists")
		}

		//fmt.Println("Getting imports for package", curPkg.Key)

		/*
		 * Add the import details for this package
		 */
		imports = curPkg.GetObjectEntry("imports")
		if imports != nil {
			if imports.Jsontype != jsonez.JSON_STRING {
				errorStr := fmt.Sprintf("%s: imports is not of type JSON_STRING"+
					" for package %s", FuncName(), curPkg.Key)
				return errors.New(errorStr)
			}

			iList := strings.Split(imports.Valstr, ",")
			for i := range iList {
				pdesc.Imports[strings.Trim(iList[i], " ")] = true
			}
		}

		/*
		 * Fetch messageTable object
		 */
		mtable, err = curPkg.Get(MESSAGETABLE)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		if mtable.Jsontype != jsonez.JSON_ARRAY {
			errorStr := fmt.Sprintf("%s: %s is not of type JSON_ARRAY",
				MESSAGETABLE, FuncName())
			return errors.New(errorStr)
		}

		//	fmt.Println("Getting message table for package", curPkg.Key)

		for j := 0; j < mtable.GetArraySize(); j++ {
			curMsg, err = mtable.GetArrayElemByIndex(j)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}

			curMsg = curMsg.Child

			//	fmt.Println("Processing message", curMsg.Key)

			/*
			 * Walk the messageTable array for this package and
			 * populate the entries
			 */
			err = pMan.ProcessMessage(pdesc, curMsg)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}

		/*
		 * Fetch rpcTable object
		 */
		rtable, err = curPkg.Get(RPCTABLE)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		if rtable.Jsontype != jsonez.JSON_ARRAY {
			errorStr := fmt.Sprintf("%s: %s is not of type JSON_ARRAY",
				RPCTABLE, FuncName())
			return errors.New(errorStr)
		}

		//		fmt.Println("Getting RPC table for package", curPkg.Key)

		for k := 0; k < rtable.GetArraySize(); k++ {
			curRpc, err = rtable.GetArrayElemByIndex(k)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}

			curRpc = curRpc.Child

			//	fmt.Println("Processing RPC", curRpc.Key)

			/*
			 * Walk the messageTable array for this package and
			 * populate the entries
			 */
			err = pMan.ProcessRPC(pdesc, curRpc)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}
	}

	//fmt.Println("Printing package map after processing all packages:")

	err = pMan.ResolveMap()

	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	return nil
}

func (pMan *ProtoDescManager) ProcessMessage(pdesc *PackageDesc,
	msg *jsonez.GoJSON) error {
	var mdesc *MsgDesc
	var curField *jsonez.GoJSON
	var err error

	//fmt.Println("Processing message", msg.Key)

	mdesc = new(MsgDesc)
	mdesc.FieldDescMap = make(FdMap)
	pdesc.temp_mMap[msg.Key] = mdesc
	//fmt.Println("Processing is", msg.Key)

	/*
	 * Walk the field table array for this message and
	 * populate the entries
	 */
	for i := 0; i < msg.GetArraySize(); i++ {
		curField, err = msg.GetArrayElemByIndex(i)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		curField = curField.Child

		//fmt.Println("Processing field", curField.Key)

		if curField != nil {
			if curField.Key == "" {
				errorStr := fmt.Sprintf("%s: goJSON object has empty key",
					FuncName())
				return errors.New(errorStr)
			}

			if curField.Jsontype != jsonez.JSON_OBJECT {
				errorStr := fmt.Sprintf(
					"%s: goJSON object with key %s is not of type JSON_OBJECT",
					FuncName(), curField.Key)
				return errors.New(errorStr)
			}

			err = pMan.ProcessField(mdesc, curField)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}
	}

	return nil
}

func (pMan *ProtoDescManager) ProcessField(mdesc *MsgDesc,
	field *jsonez.GoJSON) error {
	var fdesc *FieldDesc
	var err error

	fdesc = new(FieldDesc)

	fdesc.Fid = field.Key

	if fdesc.Fname, err = field.GetStringVal(FIELD_NAME); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	temp, err := field.GetStringVal(FIELD_LABEL)
	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	fdesc.Flabel, _ = pMan.labelMap[temp]

	temp, err = field.GetStringVal(FIELD_TYPE)
	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	fdesc.Ftype, _ = pMan.typeMap[temp]

	uintVal, err := field.GetUIntVal(FIELD_FLAGS)
	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	fdesc.Flags = uint32(uintVal)

	if fdesc.SubMessage, err = field.GetStringVal(FIELD_SUB_MSG); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	if fdesc.SubMessage == "NULL" {
		fdesc.SubMessage = ""
	}

	if fdesc.ParentOneof, err = field.GetStringVal(FIELD_ONEOF); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	if fdesc.ParentOneof == "NULL" {
		fdesc.ParentOneof = ""
	}

	//fmt.Println("sub message value is ", fdesc.SubMessage)
	fdesc.SubField = nil

	//fmt.Println("Adding field descriptor with key", field.Key)
	mdesc.FieldDescMap[field.Key] = fdesc

	//fmt.Println("Adding field descriptor with name", fdesc.Fname)
	mdesc.FieldDescMap[fdesc.Fname] = fdesc

	return nil
}

func (pMan *ProtoDescManager) ProcessRPC(pdesc *PackageDesc,
	rpc *jsonez.GoJSON) error {
	var rdesc *RpcDesc
	var err error

	rdesc = new(RpcDesc)
	pdesc.temp_rMap[rpc.Key] = rdesc
	//fmt.Println("RPC is", rpc.Key)

	rdesc.InMsgDescName, _ = rpc.GetStringVal(IN_MSG_DESC_NAME)
	rdesc.OutMsgDescName, _ = rpc.GetStringVal(OUT_MSG_DESC_NAME)
	rdesc.InMsgDescriptor = nil
	rdesc.OutMsgDescriptor = nil
	rdesc.DName, _ = rpc.GetStringVal(DAEMON_NAME)
	rdesc.InStreamAPI, _ = rpc.GetBoolVal(IN_STREAM_API)

	if rdesc.OutStreamAPI, err = rpc.GetBoolVal(OUT_STREAM_API); err != nil {
		rdesc.OutStreamAPI, _ = rpc.GetBoolVal(IS_STREAM_API)
	}

	return nil
}

/*
 * Fuction to get the message a message descriptor
 * present in a package if exists
 */
func (pMan *ProtoDescManager) GetMsgDesc(pkgName, msgName string) *MsgDesc {
	var pDesc *PackageDesc
	var mDesc *MsgDesc
	var ok bool

	//fmt.Println("Searching for message", msgName, "within package", pkgName)

	pDesc, ok = pMan.Pmap[pkgName]
	if !ok {
		return nil
	}

	//	fmt.Println("Found package", msgName)

	if mDesc, ok = pDesc.temp_mMap[msgName]; !ok {
		mDesc, ok = pDesc.Mmap[msgName]
	}

	if !ok {
		return nil
	}

	//	fmt.Println("Found message", msgName, "within package", pkgName)

	return mDesc
}

/*
 * Method to resolve message references in message
 * RPC maps
 */
func (pMan *ProtoDescManager) ResolveMap() error {
	var pkgKey, msgKey, fKey string
	var pkgDesc *PackageDesc
	var mDesc *MsgDesc
	var fDesc *FieldDesc
	/*
	 * Walk the package Map resolving the message references
	 */
	for pkgKey, pkgDesc = range pMan.Pmap {
		//	fmt.Println("Resolving references for package", pkgKey)

		/*
		 * Walk the Message table Map to resolve the sub message
		 * references
		 */
		for msgKey, mDesc = range pkgDesc.temp_mMap {
			//fmt.Println("Resolving references for message", msgKey)

			for fKey, fDesc = range mDesc.FieldDescMap {
				//fmt.Println("Resolving references for field", fKey)

				if fDesc.SubMessage != "NULL" && fDesc.SubMessage != "" {
					subDesc, ok := pkgDesc.temp_mMap[fDesc.SubMessage]
					/*
					 * If the message is not present in the package, walk
					 * the import list of this package and update the references
					 */
					if !ok {
						for pkgName, _ := range pkgDesc.Imports {
							subDesc = pMan.GetMsgDesc(pkgName, fDesc.SubMessage)

							if subDesc != nil {
								break
							}
						}
					}

					if subDesc == nil {
						s := fmt.Sprintf(
							"%s: Search for sub message %s failed for "+
								"field %s of Message %s in package %s",
							FuncName(), fDesc.SubMessage, fKey, msgKey, pkgKey)
						return errors.New(s)
					}

					fDesc.SubField = subDesc.FieldDescMap
					//fmt.Println("Resolved sub message reference for ",
					//	fDesc.SubMessage)
				}
			}
		}

		/*
		 * Walk the RPC descriptor map to resolve
		 * message references
		 */
		for rpcKey, rDesc := range pkgDesc.temp_rMap {
			inDesc, ok := pkgDesc.temp_mMap[rDesc.InMsgDescName]
			/*
			 * If the message is not present in the package, walk
			 * the import list of this package and update the references
			 */
			if !ok {
				for pkgName, _ := range pkgDesc.Imports {
					inDesc = pMan.GetMsgDesc(pkgName, rDesc.InMsgDescName)

					if inDesc != nil {
						break
					}
				}
			}

			if inDesc == nil {
				s := fmt.Sprintf("Search for input message %s failed for "+
					"RPC %s in package %s",
					rDesc.InMsgDescName, rpcKey, pkgKey)
				return errors.New(s)
			}
			rDesc.InMsgDescriptor = inDesc

			outDesc, ok := pkgDesc.temp_mMap[rDesc.OutMsgDescName]
			/*
			 * If the message is not present in the package, walk
			 * the import list of this package and update the references
			 */
			if !ok {
				for pkgName, _ := range pkgDesc.Imports {
					outDesc = pMan.GetMsgDesc(pkgName, rDesc.OutMsgDescName)
					if outDesc != nil {
						break
					}
				}
			}

			if outDesc == nil {
				s := fmt.Sprintf("Search for output message %s failed "+
					"for RPC %s",
					rDesc.OutMsgDescName, rpcKey)
				return errors.New(s)
			}
			rDesc.OutMsgDescriptor = outDesc
		}

		/*
		 *  After resolution is successful, rename the maps references
		 */
		//fmt.Println("Reassigning rpc and message descriptor maps for package",
		//		pkgKey)
		pkgDesc.Mmap = pkgDesc.temp_mMap
		pkgDesc.Rmap = pkgDesc.temp_rMap
		pkgDesc.temp_mMap = nil
		pkgDesc.temp_rMap = nil

	}

	return nil
}

/*
 * Function to get rpc descriptor
 */
func (pMan *ProtoDescManager) getRpcDesc(rpc string) (*RpcDesc, error) {
	var pkgDesc *PackageDesc
	var rDesc *RpcDesc
	var ok bool
	index := strings.Index(rpc, ".")
	if index == -1 {
		errorStr := fmt.Sprintf("%s: RPC name %s is not in the expected format",
			FuncName(), rpc)
		return nil, errors.New(errorStr)
	}

	pkgName := rpc[1:index]

	if pkgDesc, ok = pMan.Pmap[pkgName]; !ok {
		errorStr := fmt.Sprintf("%s: Package name %s is not found in package "+
			"map", FuncName(), pkgName)
		return nil, errors.New(errorStr)
	}

	if rDesc, ok = pkgDesc.Rmap[rpc]; !ok {
		errorStr := fmt.Sprintf("%s: RPC %s is not found in package %s",
			FuncName(), rpc, pkgName)
		return nil, errors.New(errorStr)
	}

	return rDesc, nil
}

/**
 * Function to initialize the proto parser
 * and build the json table from proto
 */
func InitProtoParser(filename string) {
	pMan = new(ProtoDescManager)

	if err := pMan.ReadJSONFile(filename); err != nil {
		log.Fatalf("%s: Proto parser initialization failed with error %v",
			FuncName(), err)
	}
}

/**
 * Function to get an RPC descriptor from an RPC name
 */
func GetRpcDesc(rpcName string) (*RpcDesc, error) {
	if pMan == nil {
		errorStr := fmt.Sprintf("%s: Proto descriptor map is not built. "+
			"Please register the RPCs with InitProtoParser()", FuncName())
		return nil, errors.New(errorStr)
	}

	return pMan.getRpcDesc(rpcName)
}

/**
 * Function to Get an RPC type
 */
func GetRpcType(rpcName string, rdesc *RpcDesc) (int, error) {
	var rpcType int
	var err error

	if rdesc == nil {
		if pMan == nil {
			errorStr := fmt.Sprintf("%s: Proto descriptor map is not built. "+
				"Please register the RPCs with InitProtoParser()", FuncName())
			return 0, errors.New(errorStr)
		}

		rdesc, err = pMan.getRpcDesc(rpcName)
		if err != nil {
			return 0, err
		}
	}

	if rdesc.InStreamAPI && rdesc.OutStreamAPI {
		rpcType = RPC_TYPE_BIDISTREAMING
	} else if rdesc.InStreamAPI {
		rpcType = RPC_TYPE_CLIENT_STREAMING
	} else if rdesc.OutStreamAPI {
		rpcType = RPC_TYPE_SERVER_STREAMING
	} else {
		rpcType = RPC_TYPE_UNARY
	}

	return rpcType, nil
}
