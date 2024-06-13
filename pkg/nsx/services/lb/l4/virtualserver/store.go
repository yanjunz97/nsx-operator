package vitrualserver

import (
	"errors"

	"github.com/vmware-tanzu/nsx-operator/pkg/nsx/services/common"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// VirtualServerStore is a store for virtual server
type VirtualServerStore struct {
	common.ResourceStore
}

// keyFunc is used to get the key of a resource, usually, which is the ID of the resource
func keyFunc(obj interface{}) (string, error) {
	switch v := obj.(type) {
	case *model.LBVirtualServer:
		return *v.Id, nil
	default:
		return "", errors.New("keyFunc doesn't support unknown type")
	}
}

func indexByServiceFunc(obj interface{}) ([]string, error) {
	res := make([]string, 0, 5)
	switch o := obj.(type) {
	case *model.LBVirtualServer:
		return append(res, *o.LbServicePath), nil
	default:
		return res, errors.New("indexByServiceFunc doesn't support unknown type")
	}
}

func indexByServiceAndPortFunc(obj interface{}) ([]string, error) {
	res := make([]string, 0, 5)
	switch o := obj.(type) {
	case *model.LBVirtualServer:
		for _, port := range o.Ports {
			res = append(res, fmt.Sprintf("%s|%s", *o.LbServicePath, port))
		}
		return res, nil
	default:
		return nil, errors.New("indexByServiceAndPortFunc doesn't support unknown type")
	}
}

func (virtualServerStore *VirtualServerStore) Apply(i interface{}) error {
	vs := i.(*model.LBVirtualServer)
	if vs.MarkedForDelete != nil && *vs.MarkedForDelete {
		err := virtualServerStore.Delete(vs)
		log.V(1).Info("delete virtualServer from store", "virtualServer", vs)
		if err != nil {
			return err
		}
	} else {
		err := virtualServerStore.Add(vs)
		log.V(1).Info("add virtualServer to store", "virtualServer", vs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (virtualServerStore *VirtualServerStore) GetByKey(key string) *model.LBVirtualServer {
	obj := virtualServerStore.ResourceStore.GetByKey(key)
	if obj != nil {
		VirtualServer := obj.(*model.LBVirtualServer)
		return VirtualServer
	}
	return nil
}
