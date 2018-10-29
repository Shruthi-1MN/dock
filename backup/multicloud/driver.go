// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package multicloud

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/opensds/opensds/backup"
)

const (
	ConfFile        = "/etc/opensds/driver/multi-cloud.yaml"
	UploadChunkSize = 1024 * 1024 * 50
)

func init() {
	backup.RegisterBackupCtor("multi-cloud", NewMultiCloud)
}

func NewMultiCloud() (backup.BackupDriver, error) {
	return &MultiCloud{}, nil
}

type MultiCloudConf struct {
	Endpoint      string `yaml:"Endpoint,omitempty"`
	TenantId      string `yaml:"TenantId,omitempty"`
	UploadTimeout int64  `yaml:"UploadTimeout,omitempty"`
}
type MultiCloud struct {
	client *Client
	conf   *MultiCloudConf
}

func (m *MultiCloud) loadConf(p string) (*MultiCloudConf, error) {
	conf := &MultiCloudConf{
		Endpoint:      "http://127.0.0.1:8088",
		TenantId:      DefaultTenantId,
		UploadTimeout: DefaultUploadTimeout,
	}
	confYaml, err := ioutil.ReadFile(p)
	if err != nil {
		glog.Errorf("Read config yaml file (%s) failed, reason:(%v)", p, err)
		return nil, err
	}
	if err = yaml.Unmarshal(confYaml, conf); err != nil {
		glog.Errorf("Parse error: %v", err)
		return nil, err
	}
	return conf, nil
}

func (m *MultiCloud) SetUp() error {
	// Set the default value
	var err error
	if m.conf, err = m.loadConf(ConfFile); err != nil {
		return err
	}

	opt := &AuthOptions{
		Endpoint: m.conf.Endpoint,
		TenantId: m.conf.TenantId,
	}
	if m.client, err = NewClient(opt, m.conf.UploadTimeout); err != nil {
		return err
	}

	return nil
}

func (m *MultiCloud) CleanUp() error {
	// Do nothing
	return nil
}

func (m *MultiCloud) Backup(backup *backup.BackupSpec, volFile *os.File) error {
	buf := make([]byte, UploadChunkSize)
	input := &CompleteMultipartUpload{}

	bucket := backup.Metadata["bucket"]
	key := backup.Id
	initResp, err := m.client.InitMultiPartUpload(bucket, key)
	if err != nil {
		glog.Errorf("Init part failed, err:%v", err)
		return err
	}

	defer m.client.AbortMultipartUpload(bucket, key)
	var parts []Part
	for partNum := int64(1); ; partNum++ {
		size, err := volFile.Read(buf)
		glog.Infof("read buf size len:%d", size)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if size == 0 {
			break
		}
		uploadResp, err := m.client.UploadPart(bucket, key, partNum, initResp.UploadId, buf[:size], int64(size))
		if err != nil {
			glog.Errorf("upload part failed, err:%v", err)
			return err
		}
		parts = append(parts, Part{PartNumber: partNum, ETag: uploadResp.ETag})
	}
	input.Part = parts
	_, err = m.client.CompleteMultipartUpload(bucket, key, initResp.UploadId, input)
	if err != nil {
		glog.Errorf("complete part failed, err:%v", err)
		return err
	}
	m.client.AbortMultipartUpload(bucket, key)
	glog.Infof("backup success ...")
	return nil
}

func (m *MultiCloud) Restore(backup *backup.BackupSpec, volId string, volFile *os.File) error {
	return nil
}

func (m *MultiCloud) Delete(backup *backup.BackupSpec) error {
	return nil
}
