// Code generated by mockery v2.1.0. DO NOT EDIT.

package v1

import (
	v1 "github.com/goharbor/harbor/src/pkg/scan/rest/v1"
	mock "github.com/stretchr/testify/mock"
)

// ClientPool is an autogenerated mock type for the ClientPool type
type ClientPool struct {
	mock.Mock
}

// Get provides a mock function with given fields: url, authType, accessCredential, skipCertVerify
func (_m *ClientPool) Get(url string, authType string, accessCredential string, skipCertVerify bool) (v1.Client, error) {
	ret := _m.Called(url, authType, accessCredential, skipCertVerify)

	var r0 v1.Client
	if rf, ok := ret.Get(0).(func(string, string, string, bool) v1.Client); ok {
		r0 = rf(url, authType, accessCredential, skipCertVerify)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(v1.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, bool) error); ok {
		r1 = rf(url, authType, accessCredential, skipCertVerify)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
