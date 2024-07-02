/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package vitrualserver

import (
	"sync"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/realizestate"
)

var (
	serviceIndex        	   = "nsx-op/lbvs-service"
	serviceAndPortIndex        = "nsx-op/lbvs-service-and-port"
	vipSetGroupId              = "LbVirtualServerIpSet"
	defaultDomainId            = "default"
)

type VirtualServerServiceInterface Interface {
	// virtual server

	// store related only
	// get_namespace_object_from_vs_list
	// filter_vs_keys_by_service
	// get_vs_ids_by_lbs
	// get_cluster_reserved_vs_count (usage cache)

	// tag related -> check which tag client to use
	// build_l4_lb_vs_tags
	
	// create_lb_vs_for_port
	CreateOrUpdateLbVs(lbVirtualServer model.LBVirtualServer) error
	// list_layer4_vs
	ListLayer4Vs() []*model.LBVirtualServer
	// query_lb_vs
	GetLbVsByKey(id string) *model.LBVirtualServer
	// filter_vs_by_service
	ListLbvsByService(lbService string) []*model.LBVirtualServer
	// filter_vs_by_service_and_port
	ListLbvsByServiceAndPort(lbService string, port string) []*model.LBVirtualServer


	DeleteLbVsForPort(lbVirtualServerId string) error
	UpdateLbVsAttributes(lbVirtualServerId string, lbVirtualServer model.LBVirtualServer) (model.LBVirtualServer, error)

}

// VirtualServerService implements VirtualServerServiceInterface.
type VirtualServerService struct {
	common.Service
	VirtualServerClient 		*infra.LbVirtualServersClient
	GroupClient    				*domains.GroupsClient
	Store						*VirtualServerStore
}

func NewVirtualServerService(client *nsx.Client) (*VirtualServerService, error) {
	wg := sync.WaitGroup{}
	wgDone := make(chan bool)
	fatalErrors := make(chan error)

	service := &VirtualServerService{
		Service: common.Service{
			NSXClient: client,
		},
		VirtualServerClient: &client.VirtualServerClient,
		GroupClient: &client.GroupClient
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

func (service *VirtualServerService) getVipList() map[*String]struct{} {
	virtualServers := service.Store.List()
	vipSet := map[*String]struct{}
	for _, virtualServer := range virtualServers {
		vipSet[virtualServer.IpAddress] = struct{}{}
	}
	return vipSet
}

func isSameVipSet(oldVipSet map[*String]struct{}, newVipSet map[*String]struct{}) bool {
	if (len(oldVipSet) != len(newVipSet)) {
		return false
	}
	for key := range oldVipSet {
		if _, exists := newVipSet[key]; !exists {
			return false
		}
	}
	return true
}

func updateIpExpression(group *model.Group, oldVipSet map[*String]struct{}, newVipSet map[*String]struct{}) {
	// TODO: Update the ip_addresses in group expression from vips - oldVipSet + newVipSet
	// Ref code from ncp
	// 		   expr = revised_payload["expression"]
    //         revised_vips = set(expr[0]["ip_addresses"]) if expr else set()
    //         updated_vips = revised_vips.difference(old_vips).union(new_vips)
    //         requested_payload["conditions"] = [
    //             self._nsxlib.group.build_ip_address_expression(
    //                 list(updated_vips))]
    //         LOG.info("Updating global VIPs from %s to %s in Group %s",
    //                  revised_vips, updated_vips, self._lb_vip_group)
}

func (service *VirtualServerService) CreateOrUpdateLbVs(lbVirtualServer model.LBVirtualServer) error {
	
	// TODO: check lbs limit by registering to lbs pool memeber reg
	// pending pool implementation

	oldVipSet := service.getVipList()

	// call sdk api to create lbvs
	id := *lbVirtualServer.Id
	if (id == nil) {
		return fmt.Errorf("failed to get id from lbVirtualServer")
	}
	if err := service.VirtualServerClient.Patch(id, lbVirtualServer); err != nil {
		return err
	}
	if lbvs, err := service.VirtualServerClient.Get(id); err != nil {
		return err
	}
	// call realization api to check if the creation is complete
	realizeService := realizestate.InitializeRealizeState(service.Service)
	backoff := wait.Backoff{
		Duration: 1 * time.Second,
		Factor:   2.0,
		Jitter:   0,
		Steps:    6,
	}
	if err = realizeService.CheckRealizeState(backoff, *lbvs.Path, "LbVirtualServerDto"); err != nil {
		log.Error(err, "failed to check virtual server realization state", "ID", *lbvs.Id)
		return err
	}
	// save the resource to store
	if err = service.Store.Apply(lbvs); err != nil {
		log.Error(err, "failed to add virtual server to store", "ID", *lbvs.Id)
		return err
	}
	log.Info("successfully updated virtual server", "lbvs", lbvs)

	// no need to update group if vip does not change
	newVipSet := service.getVipList()
	if isSameVipSet(oldVipSet, newVipSet) {
		return nil
	}
	// update group vipset
	// question: 
	// 1. whether it is necessary for creating vs
	// 2. what should be the vipSetGroupId?
	// 3. is it possible to reuse Group in SecurityPolicyService?
	// 4. Do I need to lock the Group when Patch/Get as ncp do?
	group, err := service.GroupClient.Get(defaultDomainId, vipSetGroupId)
	if (err != nil) {
		log.Error(err, "failed to get default vip group")
		return err
	}
	updateIpExpression(group, oldVipSet, newVipSet)

	if err = service.GroupClient.Patch(defaultDomainId, vipSetGroupId, group); err != nil {
		log.Error(err, "failed to update VIPs", "oldVips", oldVipSet, "newVips", newVipSet, "group", vipSetGroupId)
		return err
	}
	log.Info("successfully updated local VIPs", "oldVips", oldVipSet, "newVips", newVipSet, "group", vipSetGroupId)

	
	return nil
}


// question: while l4/l7 share the same vs store? they share the same in ncp
func (service *VirtualServerService) ListLayer4Vs() []*model.LBVirtualServer {
// def list_layer4_vs(self):
// return self._lb_vs_store.filter(
// 	virtual_server_type=const.VSERVER_TYPE_L4SERVICE)
// def virtual_server_type(self):
//     return self.tag_dict.get(constants.TAG_LB_VSERVER_TYPE)
// TAG_LB_VSERVER_TYPE=ncp/lb_vs_type
}

func (service *VirtualServerService) GetLbVsByKey(id string) *model.LBVirtualServer {
	lbvs := service.Store.GetByKey(id)
	if lbvs == nil {
		return nil, errors.New("Virtual server not found in store", "id", id)
	}
	return lbvs, nil
}


func (service *VirtualServerService) ListLbvsByService(lbService string) []*model.LBVirtualServer {
	return service.Store.GetByIndex(serviceIndex, lbService)
}

func (service *VirtualServerService) ListLbvsByServiceAndPort(lbService string, port string) []*model.LBVirtualServer {
	return service.Store.GetByIndex(serviceAndPortIndex, fmt.Sprintf("%s|%s", lbService, port))
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


