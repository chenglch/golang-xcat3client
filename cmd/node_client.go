package cmd

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/chenglch/golang-xcat3client/utils"
)

type NodeClient struct {
	utils.XCAT3Client
}

func NewNodeClient() (*NodeClient, error) {
	session := utils.Session{Client: http.DefaultClient, Headers: http.Header{}}
	if utils.XCAT3_URL == "" {
		return nil, errors.New("Please specified XCAT3_URL in the environment.")
	}
	client := utils.XCAT3Client{Sess: &session, Resource: utils.XCAT3_URL + "/v1/nodes"}
	service := NodeClient{client}
	return &service, nil
}

func (client *NodeClient) ToNodesMap(names []string) map[string]interface{} {
	data := make(map[string]interface{})
	nodes := make([]interface{}, 0)
	for _, name := range names {
		node := make(map[string]interface{})
		node["name"] = name
		nodes = append(nodes, node)
	}
	data["nodes"] = nodes
	return data
}

func (client *NodeClient) Show(names []string, fields []string) (interface{}, error) {
	params := url.Values{}
	if len(fields) > 0 {
		if exist, _ := utils.Contains(fields, "name"); !exist {
			fields = append(fields, "name")
		}
		params.Set("fields", strings.Join(fields, ","))
	}
	if len(names) == 1 {
		result, err := client.Sess.Get(client.Resource+"/"+names[0], &params, nil, true)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	data := client.ToNodesMap(names)
	result, err := client.Sess.Get(client.Resource+"/info", &params, data, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *NodeClient) Delete(names []string) (map[string]interface{}, error) {
	data := client.ToNodesMap(names)
	result, err := client.Sess.Delete(client.Resource, nil, data, false)
	if err != nil {
		return nil, err
	}
	ret := utils.InterfaceToMap(result)
	return ret, nil
}

func (client *NodeClient) Deploy(osimage string, state string, destroy bool, data interface{}) (map[string]interface{}, error) {
	params := url.Values{}
	if osimage != "" {
		params.Set("osimage", osimage)
	}
	if state == "" {
		state = "nodeset"
	}
	if destroy {
		state = "un_" + state
	}
	params.Set("target", state)
	result, err := client.Sess.Put(client.Resource+"/provision", &params, data, false)
	if err != nil {
		return nil, err
	}
	ret := utils.InterfaceToMap(result)
	return ret, nil
}

func (client *NodeClient) Post(url string, data interface{}) (map[string]interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Post(url, nil, data, false)
	if err != nil {
		return nil, err
	}
	ret := utils.InterfaceToMap(result)
	return ret, nil
}

func (client *NodeClient) Patch(url string, data interface{}) (map[string]interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Patch(url, nil, data, false)
	if err != nil {
		return nil, err
	}
	ret := utils.InterfaceToMap(result)
	return ret, nil
}
