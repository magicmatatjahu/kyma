package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

type ResourceListener struct {
	channel   chan<- gqlschema.ResourceEvent
	filter    func(entity interface{}) bool
	converter *ResourceConverter
}

func NewResourceListener(channel chan<- gqlschema.ResourceEvent, filter func(entity interface{}) bool, converter *ResourceConverter) *ResourceListener {
	return &ResourceListener{
		channel:   channel,
		filter:    filter,
		converter: converter,
	}
}

func (l *ResourceListener) OnAdd(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeAdd, object)
}

func (l *ResourceListener) OnUpdate(oldObject, newObject interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeUpdate, newObject)
}

func (l *ResourceListener) OnDelete(object interface{}) {
	l.onEvent(gqlschema.SubscriptionEventTypeDelete, object)
}

func (l *ResourceListener) onEvent(eventType gqlschema.SubscriptionEventType, object interface{}) {
	if object == nil {
		return
	}

	if l.filter(object) {
		l.notify(eventType, object)
	}
}

func (l *ResourceListener) notify(eventType gqlschema.SubscriptionEventType, object interface{}) {
	gqlResource, err := l.converter.ToGQL(object, nil)
	if gqlResource == nil || err != nil {
		return
	}

	event := gqlschema.ResourceEvent{
		Type:     eventType,
		Resource: *gqlResource,
	}
	l.channel <- event
}
