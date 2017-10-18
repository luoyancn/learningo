package openstack

import "encoding/json"

var ListInstance func(string, string) (string, error) = ValidateToken

func createRequestServer(name string, image_id string,
	flavor_id string, network_id string, size int) (string, error) {

	block_devices := []_block_mapping{}
	if 0 < size {
		block_devices = append(block_devices,
			_block_mapping{
				SourceType:          "image",
				DestinationType:     "volume",
				Uuid:                image_id,
				BootIndex:           0,
				DeleteOnTermination: true,
				VolumeSize:          size})
		image_id = ""
	}

	server := _reqServer{
		Server: _server{
			Name:      name,
			ImageRef:  image_id,
			FlavorRef: flavor_id,
			Max:       1,
			Min:       1,
			Networks: []_network{
				{Uuid: network_id},
			},
			BlockDevices: block_devices,
		},
	}

	server_create_str, err := json.Marshal(&server)
	if nil != err {
		return "", err
	}
	return string(server_create_str), nil
}

func GetInstance(token string, url string) (string, error) {
	body, err := request(&token, nil, "GET", url)
	return body, err
}

func CreateInstance(token string, url string,
	servername string, image_id string, flavor_id string,
	network_id string, size int) (string, error) {

	req_str, err := createRequestServer(
		servername, image_id, flavor_id, network_id, size)

	body, err := request(&token, &req_str, "POST", url)
	if nil != err {
		return "", err
	}

	var server _respServer
	err = json.Unmarshal([]byte(body), &server)
	if nil != err {
		return "", err
	}
	return server.Server.Id, nil
}

func WaitInstance(token string, url string) (string, error) {
	body, err := request(&token, nil, "GET", url)
	var server _respServer
	err = json.Unmarshal([]byte(body), &server)
	if nil != err {
		return "", err
	}
	return server.Server.Status, err
}
