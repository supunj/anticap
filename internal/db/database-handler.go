// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	type_util "github.com/supunj/anticap/internal/types"
	config_util "github.com/supunj/anticap/internal/util/config"
	crypt_util "github.com/supunj/anticap/internal/util/crypt"
	db_util "github.com/supunj/anticap/internal/util/db"
)

// SaveNode - Save node to DB
func SaveNode(node type_util.Node) (string, error) {
	var m = make(map[string]interface{})
	m["node.active"] = node.Active
	m["node.mobile"] = node.Mobile
	m["node.location.lon"] = node.Location.Lon
	m["node.location.lat"] = node.Location.Lat
	m["node.vcode"] = node.VCode

	hashedNodeID := crypt_util.GetHash(node.ID)

	return hashedNodeID, addUpdate(GetNodeKey(hashedNodeID), m)
}

// GetNode - Save node to DB
func GetNode(nodeID string) (type_util.Node, error) {
	var node type_util.Node

	nodeValues, err := db_util.NodesDB.HGetAll(fmt.Sprintf("%s%s", "node:", nodeID)).Result()
	if err != nil {
		goto End
	}

	node.Active, err = strconv.ParseBool(nodeValues["node.active"])
	node.Mobile = nodeValues["node.mobile"]
	node.Location.Lon, err = strconv.ParseFloat(nodeValues["node.location.lon"], 64)
	node.Location.Lat, err = strconv.ParseFloat(nodeValues["node.location.lat"], 64)
	node.VCode = nodeValues["node.vcode"]
End:
	return node, err
}

// ActivateNode - Make both node and the avatar active
func ActivateNode(nodeID string) error {
	var m = make(map[string]interface{})
	m["node.active"] = true

	return addUpdate(GetNodeKey(nodeID), m)
}

// SaveKeys stores the key pair generated for a node
func SaveKeys(nodeID string, privateKey string, publicKey string) error {
	var m = make(map[string]interface{})
	m["node.key.publickey"] = publicKey
	m["node.key.privatekey"] = privateKey

	return addUpdate(GetNodeKey(nodeID), m)
}

// AddSubscription - Add subscription details
func AddSubscription(nodeID string, subscription type_util.Subscription) error {
	var m = make(map[string]interface{})
	m["node.subscription."+subscription.Channel] = subscription.Channel + "|" + subscription.As

	return addUpdate(GetNodeKey(nodeID), m)
}

// AddRequest - Add request
func AddRequest(cRequest type_util.ConsumerRequest) error {
	var m = make(map[string]interface{})
	m["request.nodeid"] = cRequest.NodeID
	m["request.requestid"] = cRequest.RequestID
	m["request.chennel"] = cRequest.Channel
	m["request.opentime"] = cRequest.OpenTime.UnixNano()
	m["request.lat"] = cRequest.Lat
	m["request.lon"] = cRequest.Lon
	m["request.name"] = cRequest.Name
	m["request.active"] = cRequest.Active

	key := GetPendingRequestKey(cRequest.Channel, cRequest.RequestID)
	return addUpdate(key, m)
}

// Add or update attributes in a key
func addUpdate(key string, values map[string]interface{}) error {
	hash, err := db_util.NodesDB.HMSet(key, values).Result()

	if err != nil {
		return err
	}

	config_util.Log.Println("Hash" + hash)

	return nil
}

// GetPendingConsumerRequests - get the list of all active requests
func GetPendingConsumerRequests(channel string, resultCount int64) (type_util.ConsumerRequestList, error) {
	return GetPendingConsumerRequestsByHash(GetPendingRequestKey(channel, "*"), resultCount)
}

// GetPendingConsumerRequestByRequestID - get the consumer request for a given request id
func GetPendingConsumerRequestByRequestID(requestID string) (type_util.ConsumerRequest, error) {
	cReq, err := GetPendingConsumerRequestsByHash(GetPendingRequestKey("*", requestID), 1)

	return cReq.Requests[0], err
}

// GetPendingConsumerRequestsByHash - get the pending consumer requests for a given hash
func GetPendingConsumerRequestsByHash(hash string, resultCount int64) (type_util.ConsumerRequestList, error) {
	var crl type_util.ConsumerRequestList
	var err error
	var cr type_util.ConsumerRequest

	iter := db_util.NodesDB.Scan(0, hash, resultCount).Iterator()
	for iter.Next() {
		cr, err = GetConsumerRequest(iter.Val())
		if err != nil {
			break
		}

		cr.ConsumerRequestKey = iter.Val()

		crl.AddConsumerRequest(cr)
		fmt.Println(iter.Val())
	}

	if err == nil {
		err = iter.Err()
	}

	return crl, err
}

// GetConsumerRequest - Get consumer request from a given hash
func GetConsumerRequest(hash string) (type_util.ConsumerRequest, error) {
	cRequest, err := db_util.NodesDB.HGetAll(hash).Result()

	var uot, uct int64

	var cReq type_util.ConsumerRequest
	cReq.NodeID = cRequest["request.nodeid"]
	cReq.RequestID = cRequest["request.requestid"]
	cReq.Channel = cRequest["request.chennel"]

	uot, err = strconv.ParseInt(cRequest["request.opentime"], 10, 64)
	cReq.OpenTime = time.Unix(0, uot)

	uct, err = strconv.ParseInt(cRequest["request.closetime"], 10, 64)
	cReq.CloseTime = time.Unix(0, uct)

	cReq.Lat, err = strconv.ParseFloat(cRequest["request.lat"], 64)
	cReq.Lon, err = strconv.ParseFloat(cRequest["request.lon"], 64)
	cReq.Name = cRequest["request.name"]
	cReq.Message = cRequest["request.message"]
	cReq.Rating, err = strconv.ParseFloat(cRequest["request.rating"], 64)
	cReq.Active, err = strconv.ParseBool(cRequest["request.active"])

	return cReq, err
}

// CreateChannel - Create the channel by the given name
func CreateChannel(redisChannel string) (*type_util.RedisPubSubWrapper, error) {
	pubsub := db_util.NodesDB.Subscribe(redisChannel)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	return &(type_util.RedisPubSubWrapper{PubSub: pubsub, Msg: ch}), err
}

// ReceiveMessage - receive responces from the channel
func ReceiveMessage(redisWrapper *type_util.RedisPubSubWrapper) (string, error) {
	msg, ok := <-redisWrapper.Msg

	var err error

	if !ok {
		err = errors.New("No msg")
	}

	return msg.Payload, err
}

// PublishMessage publishes a message to the given topic
func PublishMessage(channel string, msg string) error {
	return db_util.NodesDB.Publish(channel, msg).Err()
}

// AddUpdateLocation is the method that updates a node's location periodically
func AddUpdateLocation(nodeID string, location type_util.Location) error {
	geoAdd := db_util.NodesDB.GeoAdd(
		GetNodeLocationKey(nodeID),
		&redis.GeoLocation{Longitude: location.Lon, Latitude: location.Lat, Name: nodeID},
	)

	return geoAdd.Err()
}

// AddOffer creates the offer in the db
func AddOffer(offer type_util.Offer) error {
	var m = make(map[string]interface{})
	m["offer.nodeid"] = offer.NodeID
	m["offer.requestid"] = offer.RequestID
	m["offer.offerid"] = offer.OfferID
	m["offer.offer"] = offer.Offer
	m["offer.opentime"] = offer.OpenTime.UnixNano()

	key := GetOfferKey(offer.RequestID, offer.OfferID)
	return addUpdate(key, m)
}

// GetOfferByOfferID retrieves and returns the offer by the offer id
func GetOfferByOfferID(offerID string) (type_util.Offer, error) {
	offer, err := GetOffersByHash(GetOfferKey("*", offerID), 1)

	return offer.Offers[0], err
}

// GetOffersByHash returns the list of offers for a given hash
func GetOffersByHash(hash string, resultCount int64) (type_util.OfferList, error) {
	var ofl type_util.OfferList
	var err error
	var of type_util.Offer

	iter := db_util.NodesDB.Scan(0, hash, resultCount).Iterator()
	for iter.Next() {
		of, err = GetOffer(iter.Val())
		if err != nil {
			break
		}

		of.OfferKey = iter.Val()

		ofl.AddOffer(of)
		fmt.Println(iter.Val())
	}

	if err == nil {
		err = iter.Err()
	}

	return ofl, err
}

// GetOffer returns the offer
func GetOffer(hash string) (type_util.Offer, error) {
	var of type_util.Offer
	offer, err := db_util.NodesDB.HGetAll(hash).Result()
	if err != nil {
		goto End
	}

	of.NodeID = offer["offer.nodeid"]
	of.RequestID = offer["offer.requestid"]
	of.OfferID = offer["offer.offerid"]
	of.Offer = offer["offer.offer"]

	if uot, err := strconv.ParseInt(offer["offer.opentime"], 10, 64); err != nil {
		err = nil
	} else {
		of.OpenTime = time.Unix(0, uot)
	}

	// An error here is expected as closing time is not updated until the offer is confirmed
	if uct, err := strconv.ParseInt(offer["offer.closetime"], 10, 64); err != nil {
		err = nil
	} else {
		of.CloseTime = time.Unix(0, uct)
	}

End:
	return of, err
}

// UpdateAcceptedOffer updates the accepted offer
func UpdateAcceptedOffer(offerKey string) error {
	return RenameKey(offerKey, GetAcceptedOfferKey(offerKey))
}

// RenameKey renames the given key to the other
func RenameKey(fromKey string, toKey string) error {
	return db_util.NodesDB.Rename(fromKey, toKey).Err()
}

// UpdateRejectedOffer updates the accepted offer
func UpdateRejectedOffer(offerKey string) error {
	return RenameKey(offerKey, GetRejectedOfferKey(offerKey))
}

// CancelConsumerRequest cancels a pending consumer request
func CancelConsumerRequest(consumerRequestKey string) error {
	return RenameKey(consumerRequestKey, GetCancelledConsumerRequestKey(consumerRequestKey))
}

// CompleteConsumerRequest completes a pending consumer request
func CompleteConsumerRequest(consumerRequestKey string) error {
	return RenameKey(consumerRequestKey, GetCancelledConsumerRequestKey(consumerRequestKey))
}
