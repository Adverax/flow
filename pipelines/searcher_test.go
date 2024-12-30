package pipelines

import (
	"context"
	"errors"
	"fmt"
	"github.com/adverax/core"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

// Пример реализации параллельного агрегирующего поиска на нескольких поисковиках (SearchEngine).
// Данный пример затрагивает полный цикл:
// - легкое добавление новых поисковиков
// - декомпозиция задачи (Searcher, Aggregator, pipelines)
// - ограничение времени на выполнение
// - обработка ошибок
// - gracefull shutdown
type SearchEngine interface {
	Search(ctx context.Context, request *SearchRequest) *SearchResponse
}

type SearchRequest struct {
	Query  string
	engine SearchEngine
}

type SearchResponse struct {
	Text string
	Err  error
}

type DummySearchEngine struct {
	name string
}

func (that *DummySearchEngine) Search(ctx context.Context, request *SearchRequest) *SearchResponse {
	select {
	case <-ctx.Done():
		return &SearchResponse{
			Err: ctx.Err(),
		}
	case <-time.After(50 * time.Millisecond):
	}

	return &SearchResponse{Text: fmt.Sprintf("founded %s", that.name)}
}

type Searcher struct {
	engines []SearchEngine
	timeout time.Duration
}

func NewSearcher(
	timeout time.Duration,
	engines ...SearchEngine,
) *Searcher {
	return &Searcher{
		timeout: timeout,
		engines: engines,
	}
}

func (that *Searcher) Search(ctx context.Context, query string) ([]string, error) {
	if that.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, that.timeout)
		defer cancel()
	}

	return that.search(ctx, query)
}

func (that *Searcher) search(ctx context.Context, query string) ([]string, error) {
	var result SearchAggregator

	err := Consume(
		context.Background(), // non cancelable consume
		NewFanIn(
			ctx,
			NewOperations[*SearchRequest, *SearchResponse](
				ctx,
				OpFunc[*SearchRequest, *SearchResponse](
					func(ctx context.Context, request *SearchRequest) *SearchResponse {
						return request.engine.Search(ctx, request)
					},
				),
				NewFanOut(
					ctx,
					NewIterator(ctx, that.newRequestIterator(query)),
					len(that.engines),
				)...,
			)...,
		),
		&result,
	)
	if err != nil && !errors.Is(err, context.Canceled) {
		return nil, err
	}

	return result.Summary()
}

func (that *Searcher) newRequestIterator(query string) func(yield func(*SearchRequest) bool) {
	return func(yield func(*SearchRequest) bool) {
		for _, engine := range that.engines {
			request := &SearchRequest{
				engine: engine,
				Query:  query,
			}
			if !yield(request) {
				return
			}
		}
	}
}

type SearchAggregator struct {
	values []string
	errors core.Errors
}

func (that *SearchAggregator) Aggregate(ctx context.Context, response *SearchResponse) error {
	if response.Err != nil {
		that.errors.AddError(response.Err)
		return nil
	}

	that.values = append(that.values, response.Text)
	return nil
}

func (that *SearchAggregator) Summary() ([]string, error) {
	sort.Strings(that.values)
	return that.values, that.errors.ResError()
}

func TestSearcher(t *testing.T) {
	searcher := NewSearcher(
		10*time.Second,
		&DummySearchEngine{name: "engine1"},
		&DummySearchEngine{name: "engine2"},
		&DummySearchEngine{name: "engine3"},
	)

	result, err := searcher.Search(context.Background(), "query")
	assert.NoError(t, err)
	assert.Equal(t, []string{"founded engine1", "founded engine2", "founded engine3"}, result)
}
