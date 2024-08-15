package main

import (
	"context"

	"google.golang.org/grpc"

	"github.com/denchick/mock-cluster-service-client/mockpostgresql"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
)

// ClusterServiceClient is a mockpostgresql.ClusterServiceClient with
// lazy GRPC connection initialization.
type ClusterServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

func (c *ClusterServiceClient) AddHosts(ctx context.Context, in *mockpostgresql.AddClusterHostsRequest, opts ...grpc.CallOption) (*mockpostgresql.AddClusterHostsMetadata, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return mockpostgresql.NewClusterServiceClient(conn).AddHosts(ctx, in, opts...)
}

// DeleteHosts implements mockpostgresql.ClusterServiceClient
func (c *ClusterServiceClient) DeleteHosts(ctx context.Context, in *mockpostgresql.DeleteClusterHostsRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return mockpostgresql.NewClusterServiceClient(conn).DeleteHosts(ctx, in, opts...)
}

// ListHosts implements mockpostgresql.ClusterServiceClient
func (c *ClusterServiceClient) ListHosts(ctx context.Context, in *mockpostgresql.ListClusterHostsRequest, opts ...grpc.CallOption) (*mockpostgresql.ListClusterHostsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return mockpostgresql.NewClusterServiceClient(conn).ListHosts(ctx, in, opts...)
}

type ClusterHostsIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *ClusterServiceClient
	request *mockpostgresql.ListClusterHostsRequest

	items []*mockpostgresql.Host
}

func (c *ClusterServiceClient) ClusterHostsIterator(ctx context.Context, req *mockpostgresql.ListClusterHostsRequest, opts ...grpc.CallOption) *ClusterHostsIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &ClusterHostsIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *ClusterHostsIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.ListHosts(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Hosts
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *ClusterHostsIterator) Take(size int64) ([]*mockpostgresql.Host, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*mockpostgresql.Host

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *ClusterHostsIterator) TakeAll() ([]*mockpostgresql.Host, error) {
	return it.Take(0)
}

func (it *ClusterHostsIterator) Value() *mockpostgresql.Host {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *ClusterHostsIterator) Error() error {
	return it.err
}

// UpdateHosts implements mockpostgresql.ClusterServiceClient
func (c *ClusterServiceClient) UpdateHosts(ctx context.Context, in *mockpostgresql.UpdateClusterHostsRequest, opts ...grpc.CallOption) (*operation.Operation, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return mockpostgresql.NewClusterServiceClient(conn).UpdateHosts(ctx, in, opts...)
}
