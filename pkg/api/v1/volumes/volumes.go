// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements the entry into CRUD operation of volumes.

*/

package volumes

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/opensds/opensds/pkg/api"
	"github.com/opensds/opensds/pkg/api/rpcapi"
)

type VolumeRequestDeliver interface {
	createVolume() (string, error)

	getVolume() (string, error)

	getAllVolumes() (string, error)

	updateVolume() (string, error)

	deleteVolume() (string, error)

	attachVolume() (string, error)

	detachVolume() (string, error)

	mountVolume() (string, error)

	unmountVolume() (string, error)
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	ResourceType string `json:"resourcetType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	ActionType string `json:"actionType,omitempty"`
	Host       string `json:"host,omitempty"`
	Device     string `json:"device,omitempty"`
	Attachment string `json:"attachment,omitempty"`
	MountDir   string `json:"mountDir,omitempty"`
	FsType     string `json:"fsType,omitempty"`
}

func (vr VolumeRequest) createVolume() (string, error) {
	return rpcapi.CreateVolume(vr.ResourceType, vr.Name, vr.Size)
}

func (vr VolumeRequest) getVolume() (string, error) {
	return rpcapi.GetVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) getAllVolumes() (string, error) {
	return rpcapi.GetAllVolumes(vr.ResourceType, vr.AllowDetails)
}

func (vr VolumeRequest) updateVolume() (string, error) {
	return rpcapi.UpdateVolume(vr.ResourceType, vr.Id, vr.Name)
}

func (vr VolumeRequest) deleteVolume() (string, error) {
	return rpcapi.DeleteVolume(vr.ResourceType, vr.Id)
}

func (vr VolumeRequest) attachVolume() (string, error) {
	return rpcapi.AttachVolume(vr.ResourceType, vr.Id, vr.Host, vr.Device)
}

func (vr VolumeRequest) detachVolume() (string, error) {
	return rpcapi.DetachVolume(vr.ResourceType, vr.Id, vr.Attachment)
}

func (vr VolumeRequest) mountVolume() (string, error) {
	return rpcapi.MountVolume(vr.MountDir, vr.Device, vr.FsType)
}

func (vr VolumeRequest) unmountVolume() (string, error) {
	return rpcapi.UnmountVolume(vr.MountDir)
}

func CreateVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := vrd.createVolume()
	if err != nil {
		log.Println("Create volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func GetVolume(vrd VolumeRequestDeliver) (api.VolumeDetailResponse, error) {
	var nullResponse api.VolumeDetailResponse

	result, err := vrd.getVolume()
	if err != nil {
		log.Println("Show volume error: ", err)
		return nullResponse, err
	}

	var volumeDetailResponse api.VolumeDetailResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeDetailResponse); err != nil {
		return nullResponse, err
	}
	return volumeDetailResponse, nil
}

func ListVolumes(vrd VolumeRequestDeliver) ([]api.VolumeResponse, error) {
	var nullResponses []api.VolumeResponse

	result, err := vrd.getAllVolumes()
	if err != nil {
		log.Println("List volumes error: ", err)
		return nullResponses, err
	}

	var volumesResponse []api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumesResponse); err != nil {
		return nullResponses, err
	}
	return volumesResponse, nil
}

func UpdateVolume(vrd VolumeRequestDeliver) (api.VolumeResponse, error) {
	var nullResponse api.VolumeResponse

	result, err := vrd.updateVolume()
	if err != nil {
		log.Println("Update volume error: ", err)
		return nullResponse, err
	}

	var volumeResponse api.VolumeResponse
	rbody := []byte(result)
	if err = json.Unmarshal(rbody, &volumeResponse); err != nil {
		return nullResponse, err
	}
	return volumeResponse, nil
}

func DeleteVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result, err := vrd.deleteVolume()
	if err != nil {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Delete volume error:", err)
		return defaultResponse
	} else if !strings.Contains(result, "success") {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Delete volume error!")
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func AttachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result, err := vrd.attachVolume()
	if err != nil {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Attach volume error:", err)
		return defaultResponse
	} else if !strings.Contains(result, "success") {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Attach volume error!")
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func DetachVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result, err := vrd.detachVolume()
	if err != nil {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Detach volume error:", err)
		return defaultResponse
	} else if !strings.Contains(result, "success") {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Detach volume error!")
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func MountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result, err := vrd.mountVolume()
	if err != nil {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Mount volume error:", err)
		return defaultResponse
	} else if !strings.Contains(result, "success") {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Mount volume error!")
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}

func UnmountVolume(vrd VolumeRequestDeliver) api.DefaultResponse {
	var defaultResponse api.DefaultResponse

	result, err := vrd.unmountVolume()
	if err != nil {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Unmount volume error:", err)
		return defaultResponse
	} else if !strings.Contains(result, "success") {
		defaultResponse.Status = "Failure"
		defaultResponse.Error = fmt.Sprintln("Unmount volume error!")
		return defaultResponse
	}

	defaultResponse.Status = "Success"
	return defaultResponse
}