package listener

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/golang/glog"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

//go:generate mockery -name=gqlClusterAssetConverter -output=automock -outpkg=automock -case=underscore
type gqlClusterAssetConverter interface {
	ToGQL(in *v1alpha2.ClusterAsset) (*gqlschema.ClusterAsset, error)
}

type ClusterAsset struct {
	channel   chan<- gqlschema.ClusterAssetEvent
	filter    func(entity *v1alpha2.ClusterAsset) bool
	converter gqlClusterAssetConverter
}

func NewClusterAsset(channel chan<- gqlschema.ClusterAssetEvent, filter func(entity *v1alpha2.ClusterAsset) bool, converter gqlClusterAssetConverter) *ClusterAsset {
	return &ClusterAsset{
		channel:   channel,
		filter:    filter,
		converter: converter,
	}
}

func (l *ClusterAsset) OnAdd(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeAdd, object)
}

func (l *ClusterAsset) OnUpdate(oldObject, newObject interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeUpdate, newObject)
}

func (l *ClusterAsset) OnDelete(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeDelete, object)
}

func (l *ClusterAsset) onEvent(eventType gqlschema.SubscriptionEventType, object interface{}) {
	entity, ok := object.(*v1alpha2.ClusterAsset)
	if !ok {
		glog.Error(fmt.Errorf("incorrect object type: %T, should be: *ClusterAsset", object))
		return
	}

	if l.filter(entity) {
		l.notify(eventType, entity)
	}
}

func (l *ClusterAsset) notify(eventType gqlschema.SubscriptionEventType, entity *v1alpha2.ClusterAsset) {
	gqlClusterAsset, err := l.converter.ToGQL(entity)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting *ClusterAsset"))
		return
	}
	if gqlClusterAsset == nil {
		return
	}

	event := gqlschema.ClusterAssetEvent{
		Type:                 eventType,
		ClusterAsset: 	  	  *gqlClusterAsset,
	}

	l.channel <- event
}
