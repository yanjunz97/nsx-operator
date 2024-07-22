/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package lbservice

import (
	"fmt"
	"os"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type LbServiceInterface Interface {
	GetLbVsByKey(id string) *model.LBService
}

// lbService implements LbServiceInterface.
type lbService struct {
	client loadbalancer.ServicesClient
}

func newlbService(connector client.Connector) *lbService {
	return &lbService{
		client: loadbalancer.NewServicesClient(connector),
	}
}
