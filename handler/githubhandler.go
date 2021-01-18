package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/unistack-org/micro/v3/codec"
	pb "github.com/vielendanke/test-service/proto"
)

// GithubHandler ...
type GithubHandler struct {
	codec codec.Codec
}

// NewGithubHandler ...
func NewGithubHandler(c codec.Codec) *GithubHandler {
	return &GithubHandler{c}
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
