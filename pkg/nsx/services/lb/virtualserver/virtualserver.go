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

    "github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/common"
    "github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/realizestate"
    "github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/lb/lbservice"
    "github.com/vmware-tanzu/nsx-operator/pkg/nsx"
)

var (
    LbVsIdIndex                = "nsx-op/LbVs-lbs-id"
    serviceAndPortIndex        = "nsx-op/LbVs-service-and-port"
    vsTypeL4service            = "layer_4"
)

type VirtualServerServiceInterface Interface {
    // virtual server

    // store related only
    // get_cluster_reserved_vs_count (usage cache)
    
    // create_lb_vs_for_port & update_lb_virtual_server_attributes
    // Question: do we need seperate create and update?
    // pending lbs pool memeber reg implementation
    CreateOrUpdateLbVs(lbVirtualServer *model.LBVirtualServer) error
    // list_layer4_vs
    ListLayer4Vs() []*model.LBVirtualServer
    // query_lb_vs
    GetLbVsByKey(id string) *model.LBVirtualServer
    // filter_vs_by_service
    // Question: what's the difference between service_uid and lbservice
    // what is the svc_obj refer to? Kubernetes service as it load as store.K8sServiceForLb in ncp?
    ListLbVsByService(serviceUid string) []*model.LBVirtualServer
    // filter_vs_by_service_and_port
    ListLbVsByServiceAndPort(serviceUid string, port string) []*model.LBVirtualServer
	// filter_vs_keys_by_service service_uid
    ListVsKeysByService(serviceUid string) []*string
    // get_vs_ids_by_lbs
    ListVsIdsByLbs(lbService *model.LBService) []*model.LBVirtualServer
    // get_namespace_object_from_vs_list
    // Question: we may need to add k8s client to the service to get the namespace?
	GetNamespaceFromVsList([]*model.LBVirtualServer) []*Namespace
    // build_l4_lb_vs_tags(self, svc_obj, ext_pool_id)
    BuildL4LbVsTags(k8sService *corev1.Service, extPoolId string) []model.Tag
    // delete_lb_vs_for_port
    // Question: difference between MarkedForDelete and Delete api?
    DeleteLbVsForPort(serviceUid string, port string) error
    // get_service_key_from_vs
    GetServiceKeyFromVs(lbVirtualServer *model.LBVirtualServer) string
    // list_lb_virtual_servers_for_cluster
    // Question: not listed in table, would we need it?
    // it looks like to use search url search/query?query=%s&sort_by=id, why not using store
    ListLbvsForCluster() []*model.LBVirtualServer
    // delete_crd_lb_vs_if_unused
    // Question: not listed in table, would we need it?
    DeleteCrdLbvsIfUnused(lbVirtualServer *model.LBVirtualServer) error
    // get_vs_protocol
    // Question: pending app_profile being implemented
    GetVsProtocol(lbVirtualServer *model.LBVirtualServer) string
}

// VirtualServerService implements VirtualServerServiceInterface.
type VirtualServerService struct {
    common.Service
    VirtualServerClient          *infra.LbVirtualServersClient
    GroupClient                  *domains.GroupsClient
    Store                        *VirtualServerStore
    LbService                    *LbServiceInterface
    Cluster                      *Cluster
}

func NewVirtualServerService(client *nsx.Client, lbService LbServiceInterface) (*VirtualServerService, error) {
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
                    common.TagScopeServiceUid:  indexByServiceFunc,
                    serviceAndPortIndex:        indexByServiceAndPortFunc,
                    common.TagScopeLbVsType:    indexByTypeFunc,
                    LbVsIdIndex:                indexByLbsIdFunc,
                }),
                BindingType: model.LBVirtualServerBindingType(),
            },
        },
		LbService: lbService,
        Cluster: client.Cluster
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


func (service *VirtualServerService) CreateOrUpdateLbVs(lbVirtualServer *model.LBVirtualServer) error {
    
    // TODO: check lbs limit by registering to lbs pool memeber reg
    // pending pool implementation

    // call sdk api to create LbVs
    id := *lbVirtualServer.Id
    if (id == nil) {
        return fmt.Errorf("failed to get id from lbVirtualServer")
    }
    if err := service.VirtualServerClient.Patch(id, lbVirtualServer); err != nil {
        return err
    }
    if LbVs, err := service.VirtualServerClient.Get(id); err != nil {
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
    if err = realizeService.CheckRealizeState(backoff, *LbVs.Path, "LbVirtualServerDto"); err != nil {
        log.Error(err, "failed to check virtual server realization state", "ID", *LbVs.Id)
        return err
    }
    // save the resource to store
    if err = service.Store.Apply(LbVs); err != nil {
        log.Error(err, "failed to add virtual server to store", "ID", *LbVs.Id)
        return err
    }
    log.Info("successfully updated virtual server", "LbVs", LbVs)

    // no need to update group if vip does not change
    newVipSet := service.getVipList()
    if isSameVipSet(oldVipSet, newVipSet) {
        return nil
    }
    // update group vipset is removed as it is only for tkgi plicy which is no longer used
    
    return nil
}

func (service *VirtualServerService) ListLayer4Vs() []*model.LBVirtualServer {
    return service.Store.GetByIndex(common.TagScopeLbVsType, vsTypeL4service)
}

func (service *VirtualServerService) GetLbVsByKey(id string) *model.LBVirtualServer {
    LbVs := service.Store.GetByKey(id)
    if LbVs == nil {
        return nil, errors.New("Virtual server not found in store", "id", id)
    }
    return LbVs, nil
}

func (service *VirtualServerService) ListLbVsByService(serviceUid string) []*model.LBVirtualServer {
    return service.Store.GetByIndex(common.TagScopeServiceUid, serviceUid)
}

func (service *VirtualServerService) ListLbVsByServiceAndPort(serviceUid string, port string) []*model.LBVirtualServer {
    return service.Store.GetByIndex(serviceAndPortIndex, fmt.Sprintf("%s|%s", serviceUid, port))
}

func (service *VirtualServerService) ListVsIdsByLbs(lbService *model.LBService) []*model.LBVirtualServer {
    return service.Store.GetByIndex(LbVsIdIndex, *lbService.Id)
}

func (service *VirtualServerService) GetNamespaceFromVsList([]*model.LBVirtualServer) []*Namespace {
// TODO:
// 1. get lbs id set from vs store
// 2. for all the lbs id, get the lbs from the lbs store and find the namespaces from the ns store 
// 	 a. lbs.project_uid = ns.project_uid
// 	 b. ns.ns_network_crd_name exists

}

func (service *VirtualServerService) ListVsKeysByService(serviceUid string) []*string {
    return service.Store.GetKeysByIndex(serviceIndex, serviceUid)
}


func (service *VirtualServerService) DeleteLbVsForPort(serviceUid string, port string) error {
    vsList := ListLbVsByServiceAndPort(serviceUid, port)
    var errStrings []string
    for vs := range vsList {
        if err := service.deleteLbVs(*vs.Id); err != nil {
            errStrings = append(errStrings, err.Error())
        }
    }
    if len(errStrings) > 0 {
        retrun fmt.Errorf(strings.Join(errStrings, "\n"))
    }
    return nil
}

func (service *VirtualServerService) deleteLbVs(lbVirtualServerId string) error {
    // set delete if transaction exists
    // otherwise delete with retry

    // use MarkedForDelete for here?
    return service.VirtualServerClient.Delete(lbVirtualServerId, false)
}

func (service *VirtualServerService) GetServiceKeyFromVs(lbVirtualServer *model.LBVirtualServer) string {
    for _, tag := range lbVirtualServer.Tags {
        if *tag.Scope == common.TagScopeServiceUid {
            return tag.tag
        }
    }
    return ""
}



