package grpc_server

import (
	"github.com/pygzfei/issuer-gateway/grpc/pb"
	"github.com/pygzfei/issuer-gateway/utils"
	"reflect"
	"sync"
	"testing"
)

func TestMemoryCache_Get(t *testing.T) {

	id := uint64(utils.Id())

	type fields struct {
		cache *sync.Map
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *pb.Cert
		want1  bool
	}{
		{
			name: "Run memory get success",
			fields: struct{ cache *sync.Map }{cache: func() *sync.Map {
				cache := &sync.Map{}
				cache.Store(id, &pb.Cert{Id: id})
				return cache
			}()},
			args:  struct{ id uint64 }{id: id},
			want:  &pb.Cert{Id: id},
			want1: true,
		},
		{
			name:   "Run memory get fail",
			fields: struct{ cache *sync.Map }{cache: &sync.Map{}},
			args:   struct{ id uint64 }{id: id},
			want:   nil,
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				cache: tt.fields.cache,
			}
			got, got1 := c.Get(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemoryCache_GetAll(t *testing.T) {

	id := uint64(utils.Id())

	type fields struct {
		cache *sync.Map
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "run memory get all count 1",
			fields: struct{ cache *sync.Map }{cache: func() *sync.Map {
				cache := &sync.Map{}
				cache.Store(id, &pb.Cert{Id: id})
				return cache
			}()},
			want: 1,
		},
		{
			name:   "run memory get all count 0",
			fields: struct{ cache *sync.Map }{cache: &sync.Map{}},
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				cache: tt.fields.cache,
			}
			got := c.GetAll()

			if len(got.Certs) != tt.want {
				t.Errorf("Count for GetAll() = %v, want %v", len(got.Certs), tt.want)
			}
		})
	}
}

func TestMemoryCache_Set(t *testing.T) {
	id := uint64(utils.Id())
	type fields struct {
		cache *sync.Map
	}
	type args struct {
		id    uint64
		value *pb.Cert
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "run memory set success",
			fields: struct{ cache *sync.Map }{cache: func() *sync.Map {
				return &sync.Map{}
			}()},
			args: struct {
				id    uint64
				value *pb.Cert
			}{id: id, value: &pb.Cert{Id: id}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				cache: tt.fields.cache,
			}
			c.Set(tt.args.id, tt.args.value)

			if len(c.GetAll().Certs) != tt.want {
				t.Errorf("memory set must be %d, but %d", len(c.GetAll().Certs), tt.want)
			}
		})
	}
}

func TestMemoryCache_SetRange(t *testing.T) {
	type fields struct {
		cache *sync.Map
	}
	type args struct {
		certificateList *pb.CertificateList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "run memory setRange success",
			fields: struct{ cache *sync.Map }{cache: func() *sync.Map {
				return &sync.Map{}
			}()},
			args: struct{ certificateList *pb.CertificateList }{certificateList: &pb.CertificateList{Certs: []*pb.Cert{
				{Id: 1},
				{Id: 2},
			}}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				cache: tt.fields.cache,
			}
			c.SetRange(tt.args.certificateList)

			certs := c.GetAll().Certs
			if len(certs) != tt.want {
				t.Errorf("count certs must be %d, but %d", len(certs), tt.want)
			}
		})
	}
}

func TestMewMemoryCache(t *testing.T) {
	tests := []struct {
		name string
		want *MemoryCache
	}{
		{
			name: "can new instance",
			want: MewMemoryCache(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MewMemoryCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MewMemoryCache() = %v, want %v", got, tt.want)
			}
		})
	}
}
