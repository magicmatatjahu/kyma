package genericapi

import (
	"encoding/base64"
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourcePager struct {}

func NewResourcePager() *ResourcePager {
	return &ResourcePager{}
}

func (p *ResourcePager) Page(items []unstructured.Unstructured, options gqlschema.ResourcePager) ([]unstructured.Unstructured, error) {
	itemsCount := len(items)
	if itemsCount == 0 {
		return items, nil
	}

	params, err := p.readParams(options)
	if err != nil {
		return items, err
	}

	sliceStart := skip
	sliceEnd := sliceStart + limit

	if sliceStart >= keysCount {
		return nil, fmt.Errorf("'skip' %d is out of range; maximum value: %d", sliceStart, keysCount-1)
	}
	if sliceEnd >= keysCount {
		sliceEnd = keysCount
	}

	if skip == 0 && (limit == 0 || limit >= keysCount) {
		return sortedList, nil
	}

	return sortedList[sliceStart:sliceEnd], nil
}

func (p *ResourcePager) DecodeCursor(cursor string) (*string, error) {
	if cursor == "" {
		return nil, nil
	}

	decodedValue, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "cursor is not correct")
	}

	str := string(decodedValue)
	return &str, nil
}

func (p *ResourcePager) EncodeNextCursor(uid string) string {
	return base64.StdEncoding.EncodeToString([]byte(uid))
}

type pagerParams struct {
	first int
	last int
	after string
	before string
}

func (p *ResourcePager) readParams(options gqlschema.ResourcePager) (pagerParams, error) {
	params := pagerParams{
		first:  0,
		last:   0,
		after:  "",
		before: "",
	}

	if options.First != nil {
		params.first = *options.First
	}
	if params.first < 0 {
		return params, errors.New("'first' parameter cannot be below 0")
	}

	if options.Last != nil {
		params.last = *options.Last
	}
	if params.last < 0 {
		return params, errors.New("'last' parameter cannot be below 0")
	}

	if options.After != nil {
		a, err := p.DecodeCursor(*options.After)
		if err != nil {
			return params, err
		}
		params.after = *a
	}

	if options.Before != nil {
		b, err := p.DecodeCursor(*options.Before)
		if err != nil {
			return params, err
		}
		params.before = *b
	}

	return params, nil
}
