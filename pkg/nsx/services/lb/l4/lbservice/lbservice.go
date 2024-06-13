/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package l4

import (
	"fmt"
	"os"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/loadbalancer"
)

type LbServiceInterface Interface {

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
