/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package l4

import (
	"sync"

	"github.com/sirupsen/logrus"
	vspherelog "github.com/vmware/vsphere-automation-sdk-go/runtime/log"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware-tanzu/nsx-operator/pkg/config"
	"github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/common"
)

type Client struct {
	NsxConfig     *config.NSXOperatorConfig
	RestConnector client.Connector
	Cluster       *common.Cluster

	VirtualServerClient *infra.LbVirtualServersClient
}

func restConnector(c *common.Cluster) client.Connector {
	connector, _ := c.NewRestConnector()
	return connector
}

func GetClient(cf *config.NSXOperatorConfig) *Client {
	logger := logrus.New()
	vspherelog.SetLogger(logger)
	defaultHttpTimeout := 20
	if cf.DefaultTimeout > 0 {
		defaultHttpTimeout = cf.DefaultTimeout
	}
	c := common.NewConfig(strings.Join(cf.NsxApiManagers, ","), cf.NsxApiUser, cf.NsxApiPassword, cf.CaFile, 10, 3, defaultHttpTimeout, 20, true, true, true,
		ratelimiter.AIMD, cf.GetTokenProvider(), nil, cf.Thumbprint)
	c.EnvoyHost = cf.EnvoyHost
	c.EnvoyPort = cf.EnvoyPort
	cluster, _ := common.NewCluster(c)

	virtualServerClient := infra.NewLbVirtualServersClient(restConnector(cluster))

	nsxClient := &Client{
		NsxConfig:           cf,
		RestConnector:       restConnector(cluster),
		VirtualServerClient: virtualServerClient,
		Cluster:             cluster,
	}

	return nsxClient
}
