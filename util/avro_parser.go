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
	AVSC_FILE string = "avscFile"
)

const (
	FIELD_ENUM_VALUES string = "enumValues"
	FIELD_PLACEHOLDER string = "WedgePlaceholder"
	FIELD_REPEATED    string = "repeated"
	FIELD_SUB_RECORD  string = "subRecord"
	FIELDS            string = "fields"
)

const (
	PROTOCOLTABLE string = "protocolTable"
	RECORD_TABLE  string = "recordTable"
	RESPONSE      string = "response"
	REQUEST       string = "request"
)

const (
	AVRO_TYPE_NULL = iota
	AVRO_TYPE_BOOLEAN
	AVRO_TYPE_INT
	AVRO_TYPE_LONG
	AVRO_TYPE_FLOAT
	AVRO_TYPE_DOUBLE
	AVRO_TYPE_BYTES
	AVRO_TYPE_STRING
	AVRO_TYPE_ENUM
	AVRO_TYPE_RECORD
)

type AvroFieldDesc struct {
	EnumMap     map[string]int
	Ftype       uint
	ParentOneof string
	Repeated    bool
	SubRecord   string
	SubRecMap   map[string]*AvroFieldDesc
}

type AFdMap map[string]*AvroFieldDesc

type AvroRecordDesc struct {
	FieldDescMap AFdMap
}

type ARdMap map[string]*AvroRecordDesc

type AvroProtocolDesc struct {
	AvscFile     string
	Response     string
	ResponseDesc *AvroRecordDesc
	Request      string
	RequestDesc  *AvroRecordDesc
}

type APdMap map[string]*AvroProtocolDesc

type AvroPkgDesc struct {
	Rmap    ARdMap
	Pmap    APdMap
	Imports Set
	/*
	 * Temp values used while building the map
	 */
	temp_rMap ARdMap
	temp_pMap APdMap
}

type AvroPackageDescMap map[string]*AvroPkgDesc

type AvroDescManager struct {
	Amap    AvroPackageDescMap
	avscDir string
	lock    sync.Mutex
}

var aMan *AvroDescManager

/*
 * Init method of AvroDescManager object
 * that will be used to populate the type
 * label values.
 */
func (aMan *AvroDescManager) Init() {
	aMan.Amap = make(AvroPackageDescMap)
}

func (aMan *AvroDescManager) getType(vType string) (uint, error) {
	switch vType {
	case "null":
		return AVRO_TYPE_NULL, nil
	case "boolean":
		return AVRO_TYPE_BOOLEAN, nil
	case "int":
		return AVRO_TYPE_INT, nil
	case "long":
		return AVRO_TYPE_LONG, nil
	case "float":
		return AVRO_TYPE_FLOAT, nil
	case "double":
		return AVRO_TYPE_DOUBLE, nil
	case "bytes":
		return AVRO_TYPE_BYTES, nil
	case "string":
		return AVRO_TYPE_STRING, nil
	case "enum":
		return AVRO_TYPE_ENUM, nil
	case "record":
		return AVRO_TYPE_RECORD, nil
	default:
		errorStr := fmt.Sprintf("%s: Unknown type %s", FuncName(), vType)
		return 0, errors.New(errorStr)
	}
}

func (aMan *AvroDescManager) ReadJSONFile(filename string) error {
	json, err := ioutil.ReadFile(filename)

	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	aMan.Init()

	return aMan.ParseJSON(json)
}

func (aMan *AvroDescManager) ParseJSON(input []byte) error {
	/*
	 * Read the JSON input and build a GoJSON object
	 * tree
	 */
	var root, ptable, rtable, prtable, curPkg, curRecord, curProtocol,
		imports *jsonez.GoJSON
	var pdesc *AvroPkgDesc
	var ok bool
	var err error

	root, err = jsonez.GoJSONParse(input)

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
		pdesc, ok = aMan.Amap[curPkg.Key]
		if !ok {
			//fmt.Println("Creating package descriptor for package", curPkg.Key)
			pdesc = new(AvroPkgDesc)
			pdesc.Rmap = nil
			pdesc.Pmap = nil
			pdesc.temp_rMap = make(ARdMap)
			pdesc.temp_pMap = make(APdMap)
			pdesc.Imports = make(Set)
			aMan.Amap[curPkg.Key] = pdesc
		} else {
			log.Println("Package ", curPkg.Key, " exists")
		}

		fmt.Println("Getting imports for package", curPkg.Key)

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
		 * Fetch recordTable object
		 */
		rtable, err = curPkg.Get(RECORD_TABLE)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		if rtable.Jsontype != jsonez.JSON_ARRAY {
			errorStr := fmt.Sprintf("%s: %s is not of type JSON_ARRAY",
				RECORD_TABLE, FuncName())
			return errors.New(errorStr)
		}

		for j := 0; j < rtable.GetArraySize(); j++ {
			curRecord, err = rtable.GetArrayElemByIndex(j)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}

			curRecord = curRecord.Child

			//fmt.Println("Processing message", curRecord.Key)

			/*
			 * Walk the messageTable array for this package and
			 * populate the entries
			 */
			err = aMan.ProcessRecord(pdesc, curRecord)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}

		/*
		 * Fetch rpcTable object
		 */
		prtable, err = curPkg.Get(PROTOCOLTABLE)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		if prtable.Jsontype != jsonez.JSON_ARRAY {
			errorStr := fmt.Sprintf("%s: %s is not of type JSON_ARRAY",
				PROTOCOLTABLE, FuncName())
			return errors.New(errorStr)
		}

		fmt.Println("Getting RPC table for package", curPkg.Key)

		for k := 0; k < prtable.GetArraySize(); k++ {
			curProtocol, err = prtable.GetArrayElemByIndex(k)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}

			curProtocol = curProtocol.Child

			//	fmt.Println("Processing RPC", curRpc.Key)

			/*
			 * Walk the messageTable array for this package and
			 * populate the entries
			 */
			err = aMan.ProcessProtocol(pdesc, curProtocol)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}
	}

	//fmt.Println("Printing package map after processing all packages:")

	err = aMan.ResolveMap()

	if err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	return nil
}

func (aMan *AvroDescManager) ProcessRecord(pdesc *AvroPkgDesc,
	record *jsonez.GoJSON) error {
	var rdesc *AvroRecordDesc
	var curField *jsonez.GoJSON
	var err error

	rdesc = new(AvroRecordDesc)
	rdesc.FieldDescMap = make(AFdMap)
	pdesc.temp_rMap[record.Key] = rdesc

	/*
	 * Walk the field table and build the feild descriptor map
	 * for this record
	 */
	for j := 0; j < record.GetArraySize(); j++ {
		curField, err = record.GetArrayElemByIndex(j)
		if err != nil {
			errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
			return errors.New(errorStr)
		}

		curField = curField.Child

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

			err = aMan.ProcessField(rdesc, curField)
			if err != nil {
				errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
				return errors.New(errorStr)
			}
		}
	}

	return nil
}

func (aMan *AvroDescManager) ProcessField(rdesc *AvroRecordDesc,
	field *jsonez.GoJSON) error {
	var fdesc *AvroFieldDesc
	var err error
	var typeVal, enumVal string

	fdesc = new(AvroFieldDesc)
	fName := field.Key
	fdesc.SubRecMap = nil
	fdesc.EnumMap = nil

	if typeVal, err = field.GetStringVal(FIELD_TYPE); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	if fdesc.Ftype, err = aMan.getType(typeVal); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	if fdesc.Repeated, err = field.GetBoolVal(FIELD_REPEATED); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}

	if fdesc.SubRecord, err = field.GetStringVal(FIELD_SUB_RECORD); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	if fdesc.SubRecord == "NULL" {
		fdesc.SubRecord = ""
	}

	if fdesc.ParentOneof, err = field.GetStringVal(FIELD_ONEOF); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	if fdesc.ParentOneof == "NULL" {
		fdesc.ParentOneof = ""
	}

	if enumVal, err = field.GetStringVal(FIELD_ENUM_VALUES); err != nil {
		errorStr := fmt.Sprintf("%s: %v", FuncName(), err)
		return errors.New(errorStr)
	}
	if enumVal == "NULL" {
		enumVal = ""
	}

	/*
	 * If the field is an enum, then build the enum map
	 */
	if fdesc.Ftype == AVRO_TYPE_ENUM {
		if enumVal == "" {
			errorStr := fmt.Sprintf("%s: Field %s is of type %s, "+
				"but the enumValues field is empty", FuncName(), fName, typeVal)
			return errors.New(errorStr)
		}

		/*
		 * Split the string with "," as delimiter
		 */
		enumList := strings.Split(enumVal, ",")

		fdesc.EnumMap = make(map[string]int)

		for i, entry := range enumList {
			fdesc.EnumMap[entry] = i
		}
	}

	rdesc.FieldDescMap[fName] = fdesc

	return nil
}

func (aMan *AvroDescManager) ProcessProtocol(pdesc *AvroPkgDesc,
	protocol *jsonez.GoJSON) error {
	var prDesc *AvroProtocolDesc

	prDesc = new(AvroProtocolDesc)
	pdesc.temp_pMap[protocol.Key] = prDesc

	prDesc.Request, _ = protocol.GetStringVal(REQUEST)
	prDesc.Response, _ = protocol.GetStringVal(RESPONSE)
	prDesc.RequestDesc = nil
	prDesc.ResponseDesc = nil
	prDesc.AvscFile, _ = protocol.GetStringVal(AVSC_FILE)

	return nil
}

/*
 * Fuction to get the a record descriptor
 * present in a package if exists
 */
func (aMan *AvroDescManager) GetRecordDesc(pkgName,
	recordName string) *AvroRecordDesc {
	var pDesc *AvroPkgDesc
	var recordDesc *AvroRecordDesc
	var ok bool

	pDesc, ok = aMan.Amap[pkgName]
	if !ok {
		return nil
	}

	if recordDesc, ok = pDesc.temp_rMap[recordName]; !ok {
		recordDesc, ok = pDesc.Rmap[recordName]
	}

	if !ok {
		return nil
	}

	return recordDesc
}

/*
 * Method to resolve record references in record
 * RPC maps
 */
func (aMan *AvroDescManager) ResolveMap() error {
	var pkgKey, recordKey, fieldKey, protocolKey string
	var pkgDesc *AvroPkgDesc
	var rDesc *AvroRecordDesc
	var fDesc *AvroFieldDesc
	var prDesc *AvroProtocolDesc

	/*
	 * Walk the package Map resolving the record references
	 */
	for pkgKey, pkgDesc = range aMan.Amap {
		//	fmt.Println("Resolving references for package", pkgKey)

		/*
		 * Walk the Record table Map to resolve the sub Record
		 * references
		 */
		for recordKey, rDesc = range pkgDesc.temp_rMap {
			for fieldKey, fDesc = range rDesc.FieldDescMap {
				//fmt.Println("Resolving references for field", fKey)

				if fDesc.SubRecord != "" && fDesc.SubRecord != "NULL" {
					subDesc, ok := pkgDesc.temp_rMap[fDesc.SubRecord]

					/*
					 * If the record is not present in the package, walk
					 * the import list of this package and update the references
					 */
					if !ok {
						for pkgName, _ := range pkgDesc.Imports {
							subDesc = aMan.GetRecordDesc(pkgName,
								fDesc.SubRecord)

							if subDesc != nil {
								break
							}
						}
					}

					if subDesc == nil {
						s := fmt.Sprintf(
							"%s: Search for sub message %s failed for "+
								"field %s of Message %s in package %s",
							FuncName(), fDesc.SubRecord, fieldKey, recordKey,
							pkgKey)
						return errors.New(s)
					}

					fDesc.SubRecMap = subDesc.FieldDescMap
				}
			}
		}

		/*
		 * Walk the Protocol descriptor map to resolve
		 * record references
		 */
		for protocolKey, prDesc = range pkgDesc.temp_pMap {
			requestDesc, ok := pkgDesc.temp_rMap[prDesc.Request]

			/*
			 * If the record is not present in the package, walk
			 * the import list of this package and update the references
			 */
			if !ok {
				for pkgName, _ := range pkgDesc.Imports {
					requestDesc = aMan.GetRecordDesc(pkgName, prDesc.Request)

					if requestDesc != nil {
						break
					}
				}
			}

			if requestDesc == nil {
				s := fmt.Sprintf("Search for input message %s failed for "+
					"RPC %s in package %s", prDesc.Request, protocolKey, pkgKey)
				return errors.New(s)
			}
			prDesc.RequestDesc = requestDesc

			reponseDesc, ok := pkgDesc.temp_rMap[prDesc.Response]

			/*
			 * If the record is not present in the package, walk
			 * the import list of this package and update the references
			 */
			if !ok {
				for pkgName, _ := range pkgDesc.Imports {
					reponseDesc = aMan.GetRecordDesc(pkgName, prDesc.Response)
					if reponseDesc != nil {
						break
					}
				}
			}

			if reponseDesc == nil {
				s := fmt.Sprintf("Search for output message %s failed "+
					"for RPC %s", prDesc.Response, protocolKey)
				return errors.New(s)
			}
			prDesc.ResponseDesc = reponseDesc
		}

		/*
		 *  After resolution is successful, rename the maps references
		 */
		//fmt.Println("Reassigning rpc and message descriptor maps for package",
		//		pkgKey)
		pkgDesc.Rmap = pkgDesc.temp_rMap
		pkgDesc.Pmap = pkgDesc.temp_pMap
		pkgDesc.temp_rMap = nil
		pkgDesc.temp_pMap = nil

	}

	return nil
}

/*
 * Function to get protocol descriptor
 */
func (aMan *AvroDescManager) getProtocolDesc(rpc string) (*AvroProtocolDesc,
	error) {
	var pkgDesc *AvroPkgDesc
	var prDesc *AvroProtocolDesc
	var ok bool

	index := strings.Index(rpc, ".")
	if index == -1 {
		errorStr := fmt.Sprintf("%s: RPC name %s is not in the expected format",
			FuncName(), rpc)
		return nil, errors.New(errorStr)
	}

	pkgName := rpc[1:index]

	if pkgDesc, ok = aMan.Amap[pkgName]; !ok {
		errorStr := fmt.Sprintf("%s: Package name %s is not found in package "+
			"map", FuncName(), pkgName)
		return nil, errors.New(errorStr)
	}

	if prDesc, ok = pkgDesc.Pmap[rpc]; !ok {
		errorStr := fmt.Sprintf("%s: RPC %s is not found in package %s",
			FuncName(), rpc, pkgName)
		return nil, errors.New(errorStr)
	}

	return prDesc, nil
}

/**
 * Function to initialize the proto parser
 * and build the json table from proto
 */
func InitAvroParser(filename, avscDir string) {
	aMan = new(AvroDescManager)

	if err := aMan.ReadJSONFile(filename); err != nil {
		log.Fatalf("%s: Proto parser initialization failed with error %v",
			FuncName(), err)
	}

	aMan.avscDir = avscDir
}

/**
 * Function to get the Protocol descriptor from an RPC name
 */
func GetProtocolDesc(rpcName string) (*AvroProtocolDesc, error) {
	if aMan == nil {
		errorStr := fmt.Sprintf("%s: Avro descriptor map is not built. "+
			"Please register the RPCs with InitAvroParser()", FuncName())
		return nil, errors.New(errorStr)
	}

	return aMan.getProtocolDesc(rpcName)
}

/**
 * Function to generate a schema for records specific to a
 * protocol that can be used to marshal and unmarshal the data
 */
func GetAvroSchema(protocol *AvroProtocolDesc) string {
	/*
	 * Open the avsc file corresponding to the protocol descriptor
	 */
	fileLoc := aMan.avscDir + protocol.AvscFile
	fmt.Println("AVSC file is", fileLoc)
	content, err := ioutil.ReadFile(fileLoc)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func Atest() {
	InitAvroParser("/Users/sksubra/Desktop/Wedge/WedgeAvsc.json",
		"/Users/sksubra/Desktop/Wedge/")

	val, err := GetProtocolDesc("/acl.AclService/AccessListCounterBulkGet")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("********* Generating avro schemas *********")

	fmt.Println("avro schema for request: \n", GetAvroSchema(val), "\n")

	fmt.Println("avro schema for response: \n", GetAvroSchema(val), "\n")
}
