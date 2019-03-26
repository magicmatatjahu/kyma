package wait

import (
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/waiter"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
)

func ForDocsTopicReady(namespace, name string, get func(namespace, name string) (*v1alpha1.DocsTopic, error)) error {
	return waiter.WaitAtMost(func() (bool, error) {
		res, err := get(namespace, name)
		if err != nil {
			return false, err
		}

		if res.Status.Phase == v1alpha1.DocsTopicReady {
			return true, nil
		}

		return false, nil
	}, 4 * readyTimeout)
}