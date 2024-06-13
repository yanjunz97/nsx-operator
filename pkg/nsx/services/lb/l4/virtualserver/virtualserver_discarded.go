/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package vitrualserver

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/loadbalancer"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

type VirtualServerInterface Interface {
	// virtual server

	// store related only
	// list_layer4_vs
	// query_lb_vs
	// filter_vs_by_service
	// filter_vs_by_service_and_port
	// get_namespace_object_from_vs_list
	// filter_vs_keys_by_service
	// get_vs_ids_by_lbs
	// get_cluster_reserved_vs_count (usage cache)

	// tag related -> check which tag client to use
	// build_l4_lb_vs_tags
	
	CreateLbVsForPort(lbVirtualServerParam model.LbVirtualServer) (model.LbVirtualServer, error)
	DeleteLbVsForPort(virtualServerIdParam string) error
	UpdateLbVsAttributes(virtualServerIdParam string, lbVirtualServerParam model.LbVirtualServer) (model.LbVirtualServer, error)

}

// virtualServer implements VirtualServerInterface.
type virtualServer struct {
	client loadbalancer.VirtualServersClient
}

func newVirtualServer(connector client.Connector) *virtualServer {
	return &virtualServer{
		client: loadbalancer.NewVirtualServersClient(connector),
	}
}

// do we need lbs id, ncp has this but not in LbVirtualServer
// LbService contains a list of []LbVirtualServer, need to upate LbService?
// do we need to deal with allow list?
// any relationship with rules? - maybe not as rules are for l7?
func (vs *virtualServer) CreateLbVsForPort(lbVirtualServerParam model.LbVirtualServer) (result model.LbVirtualServer, err error) {
	return vs.client.Create(lbVirtualServerParam)
}

// TODO: similarly we do not need rules for l4?
func (vs *virtualServer) DeleteLbVsForPort(virtualServerIdParam string) error {
	return vs.client.Delete(virtualServerIdParam, nil)
}

// again same question for lbs id and allow list
func (vs *virtualServer) UpdateLbVsAttributes(virtualServerIdParam string, lbVirtualServerParam model.LbVirtualServer) (model.LbVirtualServer, error) {
	// get PoolId from store: exist_lb_pool = self._lb_vs_store.get(vs_id).pool_id
	return vs.client.Update(virtualServerIdParam, lbVirtualServerParam)
}


