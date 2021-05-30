// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	db_util "github.com/supunj/anticap/internal/db"
	type_util "github.com/supunj/anticap/internal/types"
	config_util "github.com/supunj/anticap/internal/util/config"
	crypt_util "github.com/supunj/anticap/internal/util/crypt"
	json_util "github.com/supunj/anticap/internal/util/json"
)

// Register handles the register node functionality
// @Summary List accounts
// @Description Get the phone number, send the verification code
// @title Register Node API
// @version 1.0
// @host localhost:8000
// @BasePath /api/v1
func Register(w http.ResponseWriter, r *http.Request) {
	config_util.Log.Println("Register....")
	params := mux.Vars(r)

	nodeString := string(params["node"])

	node, err := json_util.GetNode([]byte(nodeString))
	if err != nil {
		w.Write([]byte("Request error!"))
		return
	}

	node.ID = crypt_util.GetHash(node.ID)
	node.VCode = GenerateVerificationCode()

	hnid, err := db_util.SaveNode(node)
	if err != nil {
		w.Write([]byte("Error!"))
	} else {
		sent := SendVerificationCode(node.VCode)

		if sent {
			w.Write([]byte(hnid + " - " + node.VCode))
		} else {
			w.Write([]byte("Retry!"))
		}
	}
}

// Verify the node via code
func Verify(w http.ResponseWriter, r *http.Request) {
	config_util.Log.Println("Verify....")
	params := mux.Vars(r)

	nodeID := string(params["nodeid"])
	vCode := string(params["vcode"])

	var node type_util.Node
	node, err := db_util.GetNode(nodeID)
	if err != nil {
		w.Write([]byte("Error!\n"))
		return
	}

	if node.Active {
		w.Write([]byte("Active!\n"))
		return
	}

	if node.VCode == vCode {
		publicKey, privateKey, err := crypt_util.GetKeys()
		if err != nil {
			w.Write([]byte("Error!\n"))
			return
		}

		err = db_util.SaveKeys(nodeID, privateKey, publicKey)
		if err != nil {
			w.Write([]byte("Error!\n"))
			return
		}

		err = db_util.ActivateNode(nodeID)
		if err != nil {
			w.Write([]byte("Activation error!\n"))
		} else {
			w.Write([]byte("Private Key - " + privateKey))
		}
	} else {
		w.Write([]byte("No!\n"))
	}
}

// UpdateProfile updates the optional/mandatory information of the person
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// TODO - Write the profile data update code
}

// Subscribe - Handle subscribe
func Subscribe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nodeID := string(params["nodeid"])
	channel := string(params["channel"])
	as := string(params["as"])

	var subscription type_util.Subscription

	subscription.Channel = channel
	subscription.As = as

	err := db_util.AddSubscription(nodeID, subscription)

	if err != nil {
		w.Write([]byte("Error!\n"))
		return
	}

	w.Write([]byte("Subscribe!\n"))
}

// SendRequest - Handle send request
func SendRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	nodeID := string(params["nodeid"])
	channel := string(params["channel"])
	request := string(params["request"])
	t := time.Now()
	//TODO - Don't allow this without a proper subscription

	req, err := json_util.GetConsumerRequest([]byte(request))
	if err != nil {
		w.Write([]byte("Request error!"))
		return
	}

	req.NodeID = nodeID
	req.Channel = channel
	req.OpenTime = t
	req.RequestID = crypt_util.GetHash(req.Channel + req.NodeID + string(t.UnixNano()))
	req.Active = true

	err = db_util.AddRequest(req)

	if err != nil {
		w.Write([]byte("Error!\n"))
		return
	}

	w.Write([]byte("SendRequest - " + req.RequestID))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ReceiveResponse - receive responces from providers
func ReceiveResponse(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	params := mux.Vars(r)
	requestID := string(params["requestid"])
	log.Println("Request id:", requestID)

	var cRequest type_util.ConsumerRequest
	cRequest, err = db_util.GetPendingConsumerRequestByRequestID(requestID)
	if err != nil || (type_util.ConsumerRequest{}) == cRequest {
		return
	}

	var pbWrapper *type_util.RedisPubSubWrapper
	pbWrapper, err = db_util.CreateChannel(cRequest.RequestID)
	defer pbWrapper.CloseChannel()
	if err != nil {
		return
	}

	for {
		resp, err := db_util.ReceiveMessage(pbWrapper)
		if err != nil {
			break
		}

		log.Println(resp)

		err = ws.WriteMessage(1, []byte(resp))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// ServeRequests - WebSocket based serving of requests
func ServeRequests(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	for {
		mt, providerRequestJSON, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", providerRequestJSON)

		reqs, err := getRequests(providerRequestJSON)
		if err != nil {
			log.Println("read:", err)
			break
		}

		reqsJSON, err := json.Marshal(reqs)
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Println(reqsJSON)

		err = ws.WriteMessage(mt, reqsJSON)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// MakeOffer creates the offer
func MakeOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nodeID := string(params["nodeid"])
	requestID := string(params["requestid"])
	offer := string(params["offer"])
	t := time.Now()

	var ofr type_util.Offer
	ofr.NodeID = nodeID
	ofr.RequestID = requestID
	ofr.OfferID = crypt_util.GetHash(requestID + nodeID + string(t.UnixNano()))
	ofr.OpenTime = t
	ofr.Offer = offer

	err := db_util.AddOffer(ofr)
	if err != nil {
		w.Write([]byte("Error creating the offer"))
		return
	}

	w.Write([]byte(ofr.OfferID))
}

// AcceptOffer accepts the offer for negitiations
func AcceptOffer(w http.ResponseWriter, r *http.Request) {
	offer, err := getAndValidateOffer(mux.Vars(r))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = db_util.UpdateAcceptedOffer(offer.OfferKey)
	if err != nil {
		w.Write([]byte("Error accepting the offer"))
	}

	//TODO - Update the pending request
	//TODO - Update the score

	return
}

// RejectOffer - confirm an offer
func RejectOffer(w http.ResponseWriter, r *http.Request) {
	offer, err := getAndValidateOffer(mux.Vars(r))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = db_util.UpdateRejectedOffer(offer.OfferKey)
	if err != nil {
		w.Write([]byte("Error rejecting the offer"))
	}

	//TODO - Update the score

	return
}

// NegotiateOffer - listen and negotiate the offer
func NegotiateOffer(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	params := mux.Vars(r)
	offerID := string(params["offerid"])

	//TODO - check null offer id

	defer ws.Close()

	log.Println("Offer id:", offerID)

	var pbWrapper *type_util.RedisPubSubWrapper
	pbWrapper, err = db_util.CreateChannel(offerID)
	defer pbWrapper.CloseChannel()
	if err != nil {
		return
	}

	for {
		mt, msg1, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if msg1 != nil {
			err = db_util.PublishMessage(offerID, string(msg1))
		}

		msg2, err := db_util.ReceiveMessage(pbWrapper)
		if err != nil {
			break
		}

		log.Println(msg2)

		err = ws.WriteMessage(mt, []byte(msg2))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// Complete handles the complete process
func Complete(w http.ResponseWriter, r *http.Request) {
	cRequest, err := getAndValidateRequest(mux.Vars(r))

	err = db_util.CompleteConsumerRequest(cRequest.ConsumerRequestKey)
	if err != nil {
		w.Write([]byte("Error completing the offer"))
	}

	//TODO - Update the feedback

	w.Write([]byte("Complete!\n"))
}

// Feedback - Handle feedb_handlerack
func Feedback(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Feedb_handlerack!\n"))
}

// UpdateLocation periodically updates the location of a node
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nodeID := string(params["nodeid"])
	lattitude, err := strconv.ParseFloat(params["lat"], 64)
	longitude, err := strconv.ParseFloat(params["lon"], 64)
	if err != nil {
		w.Write([]byte("Invalid location!\n"))
	}

	err = db_util.AddUpdateLocation(nodeID, type_util.Location{Lat: lattitude, Lon: longitude})
	if err != nil {
		w.Write([]byte("Error updating location!\n"))
	}

	w.Write([]byte("Location updated!\n"))
}

// RefreshKeys periodically refreshes the key pair
func RefreshKeys(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nodeID := string(params["nodeid"])
	oldPvtKey := string(params["oldpvtkey"])

	//TODO - Validations
	node, err := db_util.GetNode(nodeID)
	if err != nil {
		w.Write([]byte("Error!\n"))
		return
	}

	if node.ID == "" {
		w.Write([]byte("No node!\n"))
		return
	}

	if node.Secret.PrivateKey != oldPvtKey {
		w.Write([]byte("Old private key does not match!\n"))
		return
	}

	// Generate a new key pair
	privateKey, publicKey, err := crypt_util.GetKeys()
	if err != nil {
		w.Write([]byte("Error generating keys!\n"))
		return
	}

	// Save the keys in the db
	err = db_util.SaveKeys(nodeID, privateKey, publicKey)
	if err != nil {
		w.Write([]byte("Error saving keys!\n"))
		return
	}

	// Return the private key to the node
	w.Write([]byte(privateKey))
}

// CancelRequest cancels a given consumer request
func CancelRequest(w http.ResponseWriter, r *http.Request) {
	cRequest, err := getAndValidateRequest(mux.Vars(r))

	err = db_util.CancelConsumerRequest(cRequest.ConsumerRequestKey)
	if err != nil {
		w.Write([]byte("Error rejecting the offer"))
	}

	//TODO - Update the feedback

	w.Write([]byte("Cancel request"))
}
