// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package api

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	db_util "github.com/supunj/anticap/internal/db"
	type_util "github.com/supunj/anticap/internal/types"
	json_util "github.com/supunj/anticap/internal/util/json"
	validation_util "github.com/supunj/anticap/internal/validation"
)

// Returns the list of near by customer requests
func getRequests(ctx context.Context, pRequest []byte) (type_util.ConsumerRequestList, error) {
	var reqList type_util.ConsumerRequestList
	var reqOK bool
	var err error

	req, err := json_util.GetProviderRequest(pRequest)
	if err != nil {
		goto End
	}

	req.Node, err = db_util.GetNode(ctx, req.NodeID)
	if err != nil {
		goto End
	}

	reqOK, err = validation_util.ValidateProviderRequest(req.Node, req.Channel)
	if err != nil {
		goto End
	}

	if reqOK {
		reqList, err = findNearByNodes(ctx, req)
		if err != nil {
			goto End
		}
	} else {
		err = errors.New("Invalid request")
	}

End:
	return reqList, err
}

// Find nearby nodes using redis' goe features
func findNearByNodes(ctx context.Context, pRequest type_util.ProviderRequest) (type_util.ConsumerRequestList, error) {
	return db_util.GetPendingConsumerRequests(ctx, pRequest.Channel, pRequest.ResultCount)
}

var digits = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// SendVerificationCode - Send the verification code to in an SMS
func SendVerificationCode(vcode string) bool {
	return true
}

// GenerateVerificationCode - Generates the verification code
func GenerateVerificationCode() string {
	max := 6
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}

// Get the offer by the offer id and validate
func getAndValidateOffer(ctx context.Context, params map[string]string) (type_util.Offer, error) {
	nodeID := string(params["nodeid"])
	offerID := string(params["offerid"])

	offer, err := db_util.GetOfferByOfferID(ctx, offerID)
	if err != nil {
		goto End
	}

	if offer.NodeID == "" {
		err = errors.New("No offer found")
		goto End
	}

	if offer.NodeID != nodeID {
		err = errors.New("Node id does not match")
		goto End
	}
End:
	return offer, err
}

// Get the request from the id and validate
func getAndValidateRequest(ctx context.Context, params map[string]string) (type_util.ConsumerRequest, error) {
	nodeID := string(params["nodeid"])
	requestID := string(params["requestid"])

	cRequest, err := db_util.GetPendingConsumerRequestByRequestID(ctx, requestID)
	if err != nil || (type_util.ConsumerRequest{}) == cRequest {
		err = errors.New("Error getting the consumer request")
		goto End
	}

	if cRequest.NodeID != nodeID {
		err = errors.New("Node id does not match")
		goto End
	}
End:
	return cRequest, err
}
