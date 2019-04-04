package application

import (
	"context"
	"fmt"

	"github.com/kyma-project/kyma/components/application-broker/pkg/apis/applicationconnector/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"

	"github.com/golang/glog"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/application/pretty"
	contentPretty "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/content/pretty"
	assetstorePretty "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/content/storage"
	"net/http"
	"io/ioutil"
)

const (
	kymaIntegrationNamespace = "kyma-integration"
)

type eventActivationResolver struct {
	service          	eventActivationLister
	converter        	*eventActivationConverter
	contentRetriever	shared.ContentRetriever
	assetStoreRetriever shared.AssetStoreRetriever
}

//go:generate mockery -name=eventActivationLister -output=automock -outpkg=automock -case=underscore
type eventActivationLister interface {
	List(namespace string) ([]*v1alpha1.EventActivation, error)
}

func newEventActivationResolver(service eventActivationLister, contentRetriever shared.ContentRetriever, assetStoreRetriever shared.AssetStoreRetriever) *eventActivationResolver {
	return &eventActivationResolver{
		service:          		service,
		converter:        		&eventActivationConverter{},
		contentRetriever: 		contentRetriever,
		assetStoreRetriever: 	assetStoreRetriever,
	}
}

func (r *eventActivationResolver) EventActivationsQuery(ctx context.Context, namespace string) ([]gqlschema.EventActivation, error) {
	items, err := r.service.List(namespace)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while listing %s in `%s` namespace", pretty.EventActivations, namespace))
		return nil, gqlerror.New(err, pretty.EventActivations, gqlerror.WithNamespace(namespace))
	}

	return r.converter.ToGQLs(items), nil
}

func (r *eventActivationResolver) EventActivationEventsField(ctx context.Context, eventActivation *gqlschema.EventActivation) ([]gqlschema.EventActivationEvent, error) {
	if eventActivation == nil {
		glog.Errorf("EventActivation cannot be empty in order to resolve events field")
		return nil, gqlerror.NewInternal()
	}

	asyncApiSpec, err := r.contentRetriever.AsyncApiSpec().Find("service-class", eventActivation.Name)
	if err != nil {
		if module.IsDisabledModuleError(err) {
			return nil, err
		}

		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", pretty.EventActivationEvents, pretty.EventActivation, eventActivation.Name))
		return nil, gqlerror.New(err, pretty.EventActivationEvents, gqlerror.WithName(eventActivation.Name))
	}

	if asyncApiSpec == nil {
		asyncApiSpec, err = r.getAsyncApi(eventActivation.Name)
		if err != nil {
			return nil, err
		}
	}

	if asyncApiSpec == nil {
		return []gqlschema.EventActivationEvent{}, nil
	}

	if asyncApiSpec.Data.AsyncAPI != "1.0.0" {
		details := fmt.Sprintf("not supported version `%s` of %s", asyncApiSpec.Data.AsyncAPI, contentPretty.AsyncApiSpec)
		glog.Error(details)
		return nil, gqlerror.NewInternal(gqlerror.WithDetails(details))
	}

	return r.converter.ToGQLEvents(asyncApiSpec), nil
}

func (r *eventActivationResolver) getAsyncApi(eventActivationName string) (*storage.AsyncApiSpec, error) {
	types := []string{"asyncapi", "asyncApi", "asyncapispec", "asyncApiSpec", "events"}

	items, err := r.assetStoreRetriever.Asset().ListForDocsTopicByType(kymaIntegrationNamespace, eventActivationName, types)
	if err != nil {
		if module.IsDisabledModuleError(err) {
			return nil, err
		}
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", assetstorePretty.Assets, pretty.EventActivation, eventActivationName))
		return nil, gqlerror.New(err, assetstorePretty.Assets)
	}

	asyncApi := new(storage.AsyncApiSpec)
	if len(items) > 0 {
		assetRef := items[0].Status.AssetRef
		asyncApiFilePath := fmt.Sprintf("%s/%s", assetRef.BaseURL, assetRef.Files[0].Name)

		raw, err := r.fetchAsyncApi(asyncApiFilePath)
		if err != nil {
			return nil, err
		}

		err = asyncApi.Decode(raw)
		if err != nil {
			return nil, err
		}
	}

	return asyncApi, nil
}

func (r *eventActivationResolver) fetchAsyncApi(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}