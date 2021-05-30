// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package db

import "strings"

//GetPendingRequestKey generates a unique key for a new request
func GetPendingRequestKey(channel string, requestID string) string {
	return "request:" + channel + ":" + requestID + ":pending"
}

//GetNodeKey generates a unique key for a new node
func GetNodeKey(nodeID string) string {
	return "node:" + nodeID
}

//GetNodeLocationKey generates a unique key for node location
func GetNodeLocationKey(nodeID string) string {
	return "location:" + nodeID
}

// GetOfferKey returns the unique offer key
func GetOfferKey(requestID string, offerID string) string {
	return "offer:" + requestID + ":" + offerID
}

// GetAcceptedOfferKey returns the key of the accepted offer
func GetAcceptedOfferKey(offerKey string) string {
	return offerKey + ":accepted"
}

// GetRejectedOfferKey returns the key of the rejected offer
func GetRejectedOfferKey(offerKey string) string {
	return offerKey + ":rejected"
}

// GetCancelledConsumerRequestKey returns the key for a cancelled consumer request
func GetCancelledConsumerRequestKey(consumerRequestKey string) string {
	return strings.Replace(consumerRequestKey, "pending", "cancelled", -1)
}
