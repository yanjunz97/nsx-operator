/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package l4

import (
	"fmt"
	"os"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/loadbalancer"
)


type L4ServiceInterface Interface {

	// find_ip_pool_for_vs_ip
	// list_lb_pool
	// list_static_route
	// delete_persistence_profile
	// delete_duplicate_lb_pool_for_port
	// get_app_prof_id_by_protocol
	// get_or_create_lb_pool_for_port
	// delete_lb_pool_for_port
	// lbs_pool_member_reg.update
	// lbs_pool_member_reg.register
	// _nsxapi._get_valid_ip_pool
	// _nsxapi.get_ip_pools_id_by_type
	// _lb_pool_store.synced
	// _lb_pool_store.filter
	// _lb_pers_store.get_l4_persistence_prof_id
	// _lb_pers_store.get
	// _lb_vs_store.get
	// _lbs_store.get
	// allocate_lb_ip
	// reallocate_l4_ip_if_needed
	// release_lb_ip
	// get_service_key_from_vs
	// get_service_key_from_lb_pool
	// get_service_key_from_static_route
	// list_lb_ip_allocations
	// list_lb_allowed_groups
	// list_lb_hc_profile
	// filter_lb_pool_by_service
	// get_pool_id_by_service
	// filter_pers_profile_by_service
	// delete_static_route_for_service
	// delete_stale_health_check_profile_for_service
	// release_lb_ip_by_svc
	// release_cached_lbs_usage_for_svc
	// update_used_vip
	// get_vs_protocol
	// get_ip_pool_id_tag_scope
	// create_source_ip_persistence_profile
	// update_persistence_profile
	// filter_lb_pool_by_service_port
	// create_or_update_lb_hc_profile_for_pool
	// create_static_route_for_service


	

	// lb service
	// get_lb_pool_members_for_service_port
	// validate_lb_service_initializaton
	// delete_allowed_group_for_svc_lb
	// create_or_update_allowed_group_for_svc_lb
	// check_and_delete_empty_lbs
	// list_lb_service
	// get_lbs_reserved_vs_map
	// find_base_lbs_for_service
	// find_current_lbs_for_service
	// find_lb_service_usage
	// generate_names_and_create_lb_service_and_router

	// t1
	// get_tier1_link_port_ip
	// list_lb_tier1_routers
	// delete_stale_l4_tier1

	// app profile
	// get_vs_protocol
}

// l4Service implements l4ServiceInterface.
type l4Service struct {
	virtualServerService *VirtualServerService
}

func Newl4Service(client *nsx.Client) *l4Service {
	return &l4Service{
		virtualServerService:  NewVirtualServerService(client),
	}
}

