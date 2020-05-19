package pager

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/cache"
)

type IndexPager struct {
	*Pager
	indexer   cache.Indexer
	indexName string
	indexKey  string
}

func IndexerPager(indexer cache.Indexer) *IndexPager {
	return &IndexPager{
		indexer:   indexer,
	}
}

func (p *IndexPager) Limit(params PagingParams, indexName, indexKey string) ([]interface{}, error) {
	items, err := p.indexer.ByIndex(p.indexName, p.indexKey)
	if err != nil {
		return nil, errors.Wrap(err, "while getting items by index from indexer")
	}
	keys, err := p.indexer.IndexKeys(p.indexName, p.indexKey)
	if err != nil {
		return nil, errors.Wrap(err, "while getting index keys for indexer")
	}

	internalParams := p.readParams(params)
	return p.limitList(internalParams, items, keys, p.indexer)
}
