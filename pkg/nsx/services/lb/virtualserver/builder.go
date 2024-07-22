/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0 */

package vitrualserver

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/nsx-operator/pkg/util"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// build_cluster_resource_tags
// Question: reuse util.BuildBasicTags? It looks like for k8s objects, but also adds tag version/cluster
func (service *VirtualServerService) BuildL4LbvsTags(k8sService *corev1.Service, extPoolId string) []model.Tag {
	// Question: key =  cfg.CONF.nsx_v3.policy_nsxapi ? self.uid : self["namespace"], self["name"] in ncp
	// how shall we keep the same
	// the key is for tag nsx-op/service_uid and external_id
	basicTags := util.BuildBasicTags(getCluster(service), k8sService, "")
	tags := util.AppendTags(basicTags, []model.Tag{
		{Scope: common.String(common.TagScopeLbvsType), Tag: common.String(common.TagValueLbvsTypeL4service)},
		{Scope: common.String(common.TagScopeCreatedFor), Tag: common.String(common.TagValueSLB)}},
	)

	if extPoolId != "" {
		tags := util.AppendTags(basicTags, []model.Tag{
			{Scope: common.String(common.TagScopeIpPoolId), Tag: common.String(extPoolId)}},
	}
	// Question can we get cfg.CONF.k8s.enable_lb_crd variable here?
	if _, ok := k8sService.Annotations[common.AnnotationCrdLb]; ok {
		tags := util.AppendTags(basicTags, []model.Tag{
			{Scope: common.String(common.TagScopeCreatedFor), Tag: common.String(common.TagValueCrdSLB)}},
	}
	return tags
}

func getCluster(service *VirtualServerService) string {
	return service.Cluster
}
