// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package json

import (
	"encoding/json"
	"fmt"
	"reflect"

	type_util "github.com/supunj/anticap/internal/types"
)

// GetProviderRequest - Get the RequestStruct
func GetProviderRequest(pReqString []byte) (type_util.ProviderRequest, error) {
	var providerRequest type_util.ProviderRequest

	err := json.Unmarshal(pReqString, &providerRequest)
	if err != nil {
		return providerRequest, err
	}

	return providerRequest, err
}

// GetConsumerRequest - Get the consumer requesst struct from the json
func GetConsumerRequest(cReqString []byte) (type_util.ConsumerRequest, error) {
	var consumerRequest type_util.ConsumerRequest

	err := json.Unmarshal(cReqString, &consumerRequest)
	if err != nil {
		fmt.Println(err)
		return consumerRequest, err
	}

	return consumerRequest, err
}

// GetNode - Get the node struct from the json string
func GetNode(nodeString []byte) (type_util.Node, error) {
	var node type_util.Node

	err := json.Unmarshal(nodeString, &node)
	if err != nil {
		return node, err
	}

	return node, err
}

// GetObject - Get object from the json string
func GetObject(jsn []byte, typ reflect.Type) (interface{}, error) {
	obj := reflect.Zero(typ)
	err := json.Unmarshal(jsn, &obj)
	if err != nil {
		return obj, err
	}
	return obj, err
}

func GetJSON(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func GetCommonErrorResponse(errorresponse type_util.CommonErrorResponse) ([]byte, error) {
	return json.Marshal(errorresponse)
}
