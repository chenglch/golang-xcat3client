package utils

import (
	net_url "net/url"
	"strings"
)

type XCAT3Client struct {
	Sess     *Session
	Resource string
}

func (client *XCAT3Client) Put(url string, target string, data interface{}) (map[string]interface{}, error) {
	params := net_url.Values{}
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	if target != "" {
		params.Set("target", target)
	}
	result, err := client.Sess.Put(url, &params, data, false)
	if err != nil {
		return nil, err
	}
	ret := InterfaceToMap(result)
	return ret, nil
}

func (client *XCAT3Client) Get(url string, data interface{}) (interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Get(url, nil, data, false)
	if err != nil {
		return nil, err
	}
	ret := InterfaceToMap(result)
	return ret, nil
}

func (client *XCAT3Client) Patch(url string, data interface{}, retJson bool) (interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Patch(url, nil, data, false)
	if err != nil {
		return nil, err
	}
	if retJson == true {
		return result, nil
	}
	ret := InterfaceToMap(result)
	return ret, nil
}

func (client *XCAT3Client) Show(url string, fields []string, data interface{}, retJson bool) (interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	params := net_url.Values{}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	result, err := client.Sess.Get(url, &params, data, retJson)
	if err != nil {
		return nil, err
	}
	if retJson == true {
		return result, nil
	}
	ret := InterfaceToMap(result)
	return ret, nil
}

func (client *XCAT3Client) Post(url string, data interface{}, retJson bool) (interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Post(url, nil, data, retJson)
	if err != nil {
		return nil, err
	}
	if retJson == true {
		return result, nil
	}
	ret := InterfaceToMap(result)
	return ret, nil
}

func (client *XCAT3Client) Delete(url string, data interface{}, retJson bool) (interface{}, error) {
	if url != "" {
		url = client.Resource + "/" + url
	} else {
		url = client.Resource
	}
	result, err := client.Sess.Delete(url, nil, data, retJson)
	if err != nil {
		return nil, err
	}
	if retJson == true {
		return result, nil
	}
	ret := InterfaceToMap(result)
	return ret, nil
}
