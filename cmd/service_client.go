package cmd

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/chenglch/golang-xcat3client/utils"
)

type ServiceClient struct {
	utils.XCAT3Client
}

func NewServiceClient() (*ServiceClient, error) {
	session := utils.Session{Client: http.DefaultClient, Headers: http.Header{}}
	if utils.XCAT3_URL == "" {
		return nil, errors.New("Please specified XCAT3_URL in the environment.")
	}
	client := utils.XCAT3Client{Sess: &session, Resource: utils.XCAT3_URL + "/v1/services"}
	service := ServiceClient{client}
	return &service, nil
}

func (client *ServiceClient) Show(hostname string) (interface{}, error) {
	params := url.Values{}
	params.Set("name", hostname)
	result, err := client.Sess.Get(client.Resource+"/hostname", &params, nil, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}
