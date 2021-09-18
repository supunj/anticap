// Copyright 2018 Supun Jayathilake(supunj@gmail.com). All rights reserved.

// All commonly used types in the application are defined here

package types

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// Node type
type Node struct {
	ID           string         `json:"id,omitempty"`
	Active       bool           `json:"active,omitempty"`
	Mobile       string         `json:"mobile,omitempty"`
	VCode        string         `json:"vcode,omitempty"`
	Gender       string         `json:"gender,omitempty"`
	BirthDate    time.Time      `json:"bday,omitempty"`
	Location     Location       `json:"location,omitempty"`
	Availability bool           `json:"availability,omitempty"`
	Secret       Key            `json:"key,omitempty"`
	Rating       Rating         `json:"rating,omitempty"`
	Subscription []Subscription `json:"subscription,omitempty"`
}

// Location type
type Location struct {
	Lon float64 `json:"lon,omitempty"`
	Lat float64 `json:"lat,omitempty"`
}

// Key - Holds the private and public keys
type Key struct {
	PrivateKey string `json:"privatekey,omitempty"`
	PublicKey  string `json:"publickey,omitempty"`
}

// Rating - Holds all ratings
type Rating struct {
	Quality               float32 `json:"quality,omitempty"`                  // Provider
	Courtesy              float32 `json:"courtesy,omitempty"`                 // Both
	Price                 float32 `json:"price,omitempty"`                    // Provider
	OffersMade            int32   `json:"offers_made,omitempty"`              // Provider
	OffersAccepted        int32   `json:"offers_accepted,omitempty"`          // Consumer
	OffersRejected        int32   `json:"offers_rejected,omitempty"`          // Consumer
	RecommendYes          int32   `json:"recommend_yes,omitempty"`            // Provider
	RecommendNo           int32   `json:"recommend_no,omitempty"`             // Provider
	NoofServicesDelivered int32   `json:"no_of_services_delivered,omitempty"` // Provider
	PromptPayment         float32 `json:"prompt_payment,omitempty"`           // Consumer
	NoofCancelledRequests int32   `json:"no_of_cancelled_requests,omitempty"` // Consumer
}

// Subscription details
type Subscription struct {
	Channel string `json:"channel"`
	As      string `json:"as"`
}

// Request details
type Request struct {
	RequestNodeID string
	RequestID     string
	Channel       string
	OpenTime      time.Time
	CloseTime     time.Time
	Node          Node
	Location      Location
	Name          string
}

// ConsumerRequestList details
type ConsumerRequestList struct {
	Requests []ConsumerRequest `json:"requests,omitempty"`
}

// AddConsumerRequest - add a consumer request to the list
func (crl *ConsumerRequestList) AddConsumerRequest(cr ConsumerRequest) []ConsumerRequest {
	crl.Requests = append(crl.Requests, cr)
	return crl.Requests
}

// ConsumerRequest details
type ConsumerRequest struct {
	ConsumerRequestKey string    `json:"consumerrequestkey,omitempty"`
	NodeID             string    `json:"nodeid"`
	RequestID          string    `json:"requestid,omitempty"`
	Channel            string    `json:"channel"`
	NodeType           string    `json:"nodetype"`
	OpenTime           time.Time `json:"opentime,omitempty"`
	CloseTime          time.Time `json:"closetime,omitempty"`
	Lat                float64   `json:"lat"`
	Lon                float64   `json:"lon"`
	Name               string    `json:"name,omitempty"`
	Message            string    `json:"message,omitempty"`
	Rating             float64   `json:"rating,omitempty"`
	Active             bool      `json:"active,omitempty"`
}

// ProviderRequest details
type ProviderRequest struct {
	NodeID      string    `json:"nodeid"`
	Channel     string    `json:"channel"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	ResultCount int64     `json:"resultcount"`
	RequestTime time.Time `json:"requesttime,omitempty"`
	Name        string    `json:"name,omitempty"`
	Message     string    `json:"message,omitempty"`
	Node        Node      `json:"node,omitempty"`
}

// RedisPubSubWrapper - wrapper object for redis pub sub
type RedisPubSubWrapper struct {
	PubSub *redis.PubSub
	Msg    <-chan *redis.Message
}

// CloseChannel - close the channel
func (wrapper *RedisPubSubWrapper) CloseChannel() {
	wrapper.PubSub.Close()
}

// ConsumerFeedback is the feedback for the service/product delivered
type ConsumerFeedback struct {
	ResponseTime int64
}

// ProviderFeedback is the feedback for the consumer of the service/product
type ProviderFeedback struct {
}

// Offer holds the information of a provider offer
type Offer struct {
	OfferKey  string    `json:"offerkey"`
	NodeID    string    `json:"nodeid"`
	RequestID string    `json:"requestid,omitempty"`
	OfferID   string    `json:"offerid,omitempty"`
	Offer     string    `json:"offer,omitempty"`
	OpenTime  time.Time `json:"opentime,omitempty"`
	CloseTime time.Time `json:"closetime,omitempty"`
}

// OfferList details
type OfferList struct {
	Offers []Offer `json:"offers,omitempty"`
}

// AddOffer adds an to the offers list
func (ofl *OfferList) AddOffer(of Offer) []Offer {
	ofl.Offers = append(ofl.Offers, of)
	return ofl.Offers
}

// CommonErrorResponse is the error response commonly used to return an error
type CommonErrorResponse struct {
	ErrorSource  string `json:"error_source,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
