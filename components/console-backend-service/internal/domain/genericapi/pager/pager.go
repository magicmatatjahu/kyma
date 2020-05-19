package pager

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/cache"
	"sort"
	"strings"
)

type PagingParams struct {
	Limit  int
	Skip   *int
	After  *string
}

type Pager struct {
	store cache.Store
}

func StorePager(store cache.Store) *Pager {
	return &Pager{
		store: store,
	}
}

func (p *Pager) Limit(params PagingParams) ([]interface{}, error) {
	items := p.store.List()
	keys := p.store.ListKeys()
	internalParams := p.readParams(params)
	return p.limitList(internalParams, items, keys, p.store)
}

type itemGetter interface {
	GetByKey(key string) (item interface{}, exists bool, err error)
}

func (p *Pager) readParams(params PagingParams) PagingParams {
	skip := 0
	if params.Skip != nil {
		skip = *params.Skip
	}

	after := ""
	if params.After != nil {
		after = *params.After
	}

	return PagingParams{
		Limit: params.Limit,
		Skip:  &skip,
		After: &after,
	}
}

func (p *Pager) limitList(params PagingParams, items []interface{}, keys []string, getter itemGetter) ([]interface{}, error) {
	if len(items) == 0 {
		return []interface{}{}, nil
	}

	keysCount := len(keys)
	limit := params.Limit
	skip := *params.Skip

	if limit < 0 {
		return nil, errors.New("'limit' parameter cannot be below 0")
	}
	if skip < 0 {
		return nil, errors.New("'skip' parameter cannot be below 0")
	}

	sliceStart := skip
	sliceEnd := sliceStart + limit

	if sliceStart >= keysCount {
		return nil, fmt.Errorf("'skip' %d is out of range; maximum value: %d", sliceStart, keysCount-1)
	}
	if sliceEnd >= keysCount {
		sliceEnd = keysCount
	}

	sortedList, err := p.sortByKey(keys, getter)
	if err != nil {
		return nil, errors.Wrap(err, "while sorting store")
	}

	if skip == 0 && (limit == 0 || limit >= keysCount) {
		return sortedList, nil
	}

	return sortedList[sliceStart:sliceEnd], nil
}

func (p *Pager) sortByKey(keys []string, store itemGetter) ([]interface{}, error) {
	var sortedKeys []string
	sortedKeys = append(sortedKeys, keys...)

	sort.SliceStable(sortedKeys, func(i, j int) bool {
		result := strings.Compare(sortedKeys[i], sortedKeys[j])
		return result != 1
	})

	var sortedList []interface{}
	for _, key := range sortedKeys {
		item, exists, err := store.GetByKey(key)
		if !exists {
			return nil, fmt.Errorf("item with key %s doesn't exist", key)
		}
		if err != nil {
			return nil, errors.Wrapf(err, "while getting item with key %s", key)
		}
		sortedList = append(sortedList, item)
	}

	return sortedList, nil
}
