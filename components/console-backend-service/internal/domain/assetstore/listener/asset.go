package listener

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/golang/glog"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

//go:generate mockery -name=gqlAssetConverter -output=automock -outpkg=automock -case=underscore
type gqlAssetConverter interface {
	ToGQL(in *v1alpha2.Asset) (*gqlschema.Asset, error)
}

type Asset struct {
	channel   chan<- gqlschema.AssetEvent
	filter    func(entity *v1alpha2.Asset) bool
	converter gqlAssetConverter
}

func NewAsset(channel chan<- gqlschema.AssetEvent, filter func(entity *v1alpha2.Asset) bool, converter gqlAssetConverter) *Asset {
	return &Asset{
		channel:   channel,
		filter:    filter,
		converter: converter,
	}
}

func (l *Asset) OnAdd(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeAdd, object)
}

func (l *Asset) OnUpdate(oldObject, newObject interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeUpdate, newObject)
}

func (l *Asset) OnDelete(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeDelete, object)
}

func (l *Asset) onEvent(eventType gqlschema.SubscriptionEventType, object interface{}) {
	entity, ok := object.(*v1alpha2.Asset)
	if !ok {
		glog.Error(fmt.Errorf("incorrect object type: %T, should be: *Asset", object))
		return
	}

	if l.filter(entity) {
		l.notify(eventType, entity)
	}
}

func (l *Asset) notify(eventType gqlschema.SubscriptionEventType, entity *v1alpha2.Asset) {
	gqlAssetTopic, err := l.converter.ToGQL(entity)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting *Asset"))
		return
	}
	if gqlAssetTopic == nil {
		return
	}

	event := gqlschema.AssetEvent{
		Type:         eventType,
		Asset: 	  	  *gqlAssetTopic,
	}

	l.channel <- event
}
