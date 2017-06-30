package cmd

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
)

type NicClient struct {
	utils.XCAT3Client
}

func NewNicClient() (*NicClient, error) {
	session := utils.Session{Client: http.DefaultClient, Headers: http.Header{}}
	if utils.XCAT3_URL == "" {
		return nil, errors.New("Please specified XCAT3_URL in the environment.")
	}
	client := utils.XCAT3Client{Sess: &session, Resource: utils.XCAT3_URL + "/v1/nics"}
	service := NicClient{client}
	return &service, nil
}

func (client *NicClient) Show(uuid string, fields []string) (interface{}, error) {
	params := url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	result, err := client.Sess.Get(client.Resource+"/"+uuid, &params, nil, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *NicClient) GetByMac(mac string, fields []string) (interface{}, error) {
	params := url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	params.Set("mac", mac)
	result, err := client.Sess.Get(client.Resource+"/address", &params, nil, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}
