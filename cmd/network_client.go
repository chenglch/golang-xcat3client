package cmd

import (
	"errors"
	"net/http"

	"github.com/chenglch/golang-xcat3client/utils"
)

type NetworkClient struct {
	utils.XCAT3Client
}

func NewNetworkClient() (*NetworkClient, error) {
	session := utils.Session{Client: http.DefaultClient, Headers: http.Header{}}
	if utils.XCAT3_URL == "" {
		return nil, errors.New("Please specified XCAT3_URL in the environment.")
	}
	client := utils.XCAT3Client{Sess: &session, Resource: utils.XCAT3_URL + "/v1/networks"}
	service := NetworkClient{client}
	return &service, nil
}
