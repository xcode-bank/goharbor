// Code generated by mockery v2.1.0. DO NOT EDIT.

package flow

import (
	distribution "github.com/docker/distribution"

	io "io"

	mock "github.com/stretchr/testify/mock"

	model "github.com/goharbor/harbor/src/pkg/reg/model"
)

// mockAdapter is an autogenerated mock type for the registryAdapter type
type mockAdapter struct {
	mock.Mock
}

// BlobExist provides a mock function with given fields: repository, digest
func (_m *mockAdapter) BlobExist(repository string, digest string) (bool, error) {
	ret := _m.Called(repository, digest)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(repository, digest)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(repository, digest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteManifest provides a mock function with given fields: repository, reference
func (_m *mockAdapter) DeleteManifest(repository string, reference string) error {
	ret := _m.Called(repository, reference)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(repository, reference)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTag provides a mock function with given fields: repository, tag
func (_m *mockAdapter) DeleteTag(repository string, tag string) error {
	ret := _m.Called(repository, tag)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(repository, tag)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchArtifacts provides a mock function with given fields: filters
func (_m *mockAdapter) FetchArtifacts(filters []*model.Filter) ([]*model.Resource, error) {
	ret := _m.Called(filters)

	var r0 []*model.Resource
	if rf, ok := ret.Get(0).(func([]*model.Filter) []*model.Resource); ok {
		r0 = rf(filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Resource)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.Filter) error); ok {
		r1 = rf(filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthCheck provides a mock function with given fields:
func (_m *mockAdapter) HealthCheck() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Info provides a mock function with given fields:
func (_m *mockAdapter) Info() (*model.RegistryInfo, error) {
	ret := _m.Called()

	var r0 *model.RegistryInfo
	if rf, ok := ret.Get(0).(func() *model.RegistryInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RegistryInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ManifestExist provides a mock function with given fields: repository, reference
func (_m *mockAdapter) ManifestExist(repository string, reference string) (bool, *distribution.Descriptor, error) {
	ret := _m.Called(repository, reference)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(repository, reference)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *distribution.Descriptor
	if rf, ok := ret.Get(1).(func(string, string) *distribution.Descriptor); ok {
		r1 = rf(repository, reference)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*distribution.Descriptor)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(repository, reference)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PrepareForPush provides a mock function with given fields: _a0
func (_m *mockAdapter) PrepareForPush(_a0 []*model.Resource) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*model.Resource) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PullBlob provides a mock function with given fields: repository, digest
func (_m *mockAdapter) PullBlob(repository string, digest string) (int64, io.ReadCloser, error) {
	ret := _m.Called(repository, digest)

	var r0 int64
	if rf, ok := ret.Get(0).(func(string, string) int64); ok {
		r0 = rf(repository, digest)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 io.ReadCloser
	if rf, ok := ret.Get(1).(func(string, string) io.ReadCloser); ok {
		r1 = rf(repository, digest)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(io.ReadCloser)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(repository, digest)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PullManifest provides a mock function with given fields: repository, reference, accepttedMediaTypes
func (_m *mockAdapter) PullManifest(repository string, reference string, accepttedMediaTypes ...string) (distribution.Manifest, string, error) {
	_va := make([]interface{}, len(accepttedMediaTypes))
	for _i := range accepttedMediaTypes {
		_va[_i] = accepttedMediaTypes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, repository, reference)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 distribution.Manifest
	if rf, ok := ret.Get(0).(func(string, string, ...string) distribution.Manifest); ok {
		r0 = rf(repository, reference, accepttedMediaTypes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(distribution.Manifest)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string, ...string) string); ok {
		r1 = rf(repository, reference, accepttedMediaTypes...)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, ...string) error); ok {
		r2 = rf(repository, reference, accepttedMediaTypes...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PushBlob provides a mock function with given fields: repository, digest, size, blob
func (_m *mockAdapter) PushBlob(repository string, digest string, size int64, blob io.Reader) error {
	ret := _m.Called(repository, digest, size, blob)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int64, io.Reader) error); ok {
		r0 = rf(repository, digest, size, blob)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PushManifest provides a mock function with given fields: repository, reference, mediaType, payload
func (_m *mockAdapter) PushManifest(repository string, reference string, mediaType string, payload []byte) (string, error) {
	ret := _m.Called(repository, reference, mediaType, payload)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, []byte) string); ok {
		r0 = rf(repository, reference, mediaType, payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, []byte) error); ok {
		r1 = rf(repository, reference, mediaType, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
