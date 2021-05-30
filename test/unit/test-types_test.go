package unit

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/erggo/datafiller"

	type_util "github.com/supunj/anticap/internal/types"
	json_util "github.com/supunj/anticap/internal/util/json"
)

func TestNodeToJSON(t *testing.T) {
	/* node := type_util.Node{
		ID:     "5c29be96-8450-4c6b-9fe2-4efad3f74bd5",
		Active: false,
		Avatar: type_util.Avatar{
			ID:     "90eb0846-4923-4a70-86fb-eabea010c448",
			Mobile: "+94777740359",
			Active: false,
			Location: type_util.Location{
				Lon: 6.801890,
				Lat: 79.893214,
			},
		},
	} */
	node := type_util.Node{}

	datafiller.Fill(&node)

	jsn, err := json.Marshal(node)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsn))
}

func TestConsumerRequestToJSON(t *testing.T) {
	creq := type_util.ConsumerRequest{}

	datafiller.Fill(&creq)

	jsn, err := json.Marshal(creq)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsn))
}

func TestProviderRequestToJSON(t *testing.T) {
	pRequest := type_util.ProviderRequest{}

	datafiller.Fill(&pRequest)

	jsn, err := json.Marshal(pRequest)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsn))
}

func TestDynamicType(t *testing.T) {
	obj, err := json_util.GetObject([]byte("{\n    \"id\": \"test\",\n    \"active\": true,\n    \"avatar\": {\n        \"id\": \"test\",\n        \"mobile\": \"test\",\n        \"active\": false,\n        \"vcode\": \"test\",\n        \"gender\": \"test\",\n        \"bday\": \"2031-09-21T23:20:10+05:30\",\n        \"location\": {\n            \"lon\": 0.5412147045135498,\n            \"lat\": 0.8362488746643066\n        },\n        \"availability\": false,\n        \"key\": {\n            \"privatekey\": \"test\",\n            \"publickey\": \"test\"\n        },\n        \"rating\": {},\n        \"subscription\": [\n            {\n                \"channel\": \"test\",\n                \"as\": \"test\"\n            },\n            {\n                \"channel\": \"test\",\n                \"as\": \"test\"\n            }\n        ]\n    }\n}"), reflect.TypeOf(type_util.Node{}))
	fmt.Println(obj, err)
}
