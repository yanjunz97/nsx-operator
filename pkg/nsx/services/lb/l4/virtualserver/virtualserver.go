/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package vitrualserver

import (
	"sync"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var (
	serviceIndex        	   = "nsx-op/lbvs-service"
	serviceAndPortIndex        = "nsx-op/lbvs-service-and-port"
)

type VirtualServerServiceInterface Interface {
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
	
	CreateLbVsForPort(lbVirtualServer model.LBVirtualServer) error
	DeleteLbVsForPort(lbVirtualServerId string) error
	UpdateLbVsAttributes(lbVirtualServerId string, lbVirtualServer model.LBVirtualServer) (model.LBVirtualServer, error)

}

// VirtualServerService implements VirtualServerServiceInterface.
type VirtualServerService struct {
	VirtualServerClient 		*infra.LbVirtualServersClient
	Store						*VirtualServerStore
}

func NewVirtualServerService(client *nsx.Client) (*VirtualServerService, error) {
	wg := sync.WaitGroup{}
	wgDone := make(chan bool)
	fatalErrors := make(chan error)

	service := &VirtualServerService{
		VirtualServerClient: client.VirtualServerClient,
		Store: &VirtualServerStore {
			ResourceStore: common.ResourceStore {
				Indexer: cache.NewIndexer(keyFunc, cache.Indexers{
					serviceIndex:    		indexByServiceFunc,
					serviceAndPortIndex:	indexByServiceAndPortFunc,
				}),
				BindingType: model.LBVirtualServerBindingType(),
			},
		},
	}
	wg.Add(1)
	go service.InitializeResourceStore(&wg, fatalErrors, "virtualServer", nil, service.store)
	go func() {
		wg.Wait()
		close(wgDone)
	}()
	select {
	case <-wgDone:
		break
	case err := <-fatalErrors:
		close(fatalErrors)
		return service, err
	}
	return service, nil
}

func (service *VirtualServerService) CreateLbVsForPort(lbVirtualServer model.LBVirtualServer) error {
	
	// check lbs limit by registering to lbs pool memeber reg

	// call sdk api to create lbvs
	service.VirtualServerClient.Patch(*lbVirtualServer.Id, lbVirtualServer)
	// call realization api to check if the creation is complete

	// for lbvs create for SLB or CRD_LB (by checking the tag) 
	// 		got the old vips
	// 		update the store by adding the new resource
	// 		get the new vipds
	//  	update_lb_vip_set: update group LbVirtualServerIpSet with new vips

}

func (service *VirtualServerService) DeleteLbVsForPort(lbVirtualServerId string) error {
	// set delete if transaction exists
	// otherwise delete with retry

	// use MarkedForDelete for here?
	return service.VirtualServerClient.Delete(lbVirtualServerId, false)
}

func (service *VirtualServerService) UpdateLbVsAttributes(virtualServerId string, lbVirtualServer model.LBVirtualServer) (model.LBVirtualServer, error) {
	// get PoolId from store: exist_lb_pool = self._lb_vs_store.get(vs_id).pool_id
	return service.VirtualServerClient.Update(virtualServerId, lbVirtualServer)
}


