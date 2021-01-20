package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/codec"
	pb "github.com/vielendanke/test-service/proto"
)

// GithubHandler ...
type GithubHandler struct {
	codec  codec.Codec
	client client.Client
}

// MessageData ...
type MessageData struct {
	Event     *pb.EventMessage
	MediaType string
	DestTopic string
}

// Topic ...
func (md *MessageData) Topic() string {
	return md.DestTopic
}

// Payload ...
func (md *MessageData) Payload() interface{} {
	return md.Event
}

// ContentType ...
func (md *MessageData) ContentType() string {
	return md.MediaType
}

// NewGithubHandler ...
func NewGithubHandler(c codec.Codec, cl client.Client) *GithubHandler {
	return &GithubHandler{c, cl}
}

// ValidMessage ...
func (gh *GithubHandler) ValidMessage(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/users/vielendanke")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	message := &pb.Response{Result: string(data)}
	gh.codec.Write(w, nil, message)
}

// InvalidMessage ...
func (gh *GithubHandler) InvalidMessage(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/users/sdakjdwiqundinq")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	message := &pb.Response{Result: string(data)}
	gh.codec.Write(w, nil, message)
}

// PublishMessage ...
func (gh *GithubHandler) PublishMessage(w http.ResponseWriter, r *http.Request) {
	data := &MessageData{
		Event:     &pb.EventMessage{Event: "test"},
		MediaType: "application/json",
		DestTopic: "message-topic",
	}
	if err := gh.client.Publish(r.Context(), data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		gh.codec.Write(w, nil, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
