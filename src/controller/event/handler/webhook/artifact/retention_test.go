package artifact

import (
	"testing"
	"time"

	"github.com/goharbor/harbor/src/pkg/retention"

	"github.com/goharbor/harbor/src/controller/event"
	"github.com/goharbor/harbor/src/core/config"
	"github.com/goharbor/harbor/src/lib/selector"
	"github.com/goharbor/harbor/src/pkg/notification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetentionHandler_Handle(t *testing.T) {
	config.Init()
	handler := &RetentionHandler{RetentionController: DefaultRetentionControllerFunc}

	policyMgr := notification.PolicyMgr
	retentionCtlFunc := handler.RetentionController

	defer func() {
		notification.PolicyMgr = policyMgr
		handler.RetentionController = retentionCtlFunc
	}()
	notification.PolicyMgr = &fakedNotificationPolicyMgr{}
	handler.RetentionController = retention.FakedRetentionControllerFunc

	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RetentionHandler Want Error 1",
			args: args{
				data: "",
			},
			wantErr: true,
		},
		{
			name: "RetentionHandler 1",
			args: args{
				data: &event.RetentionEvent{
					OccurAt: time.Now(),
					Deleted: []*selector.Result{
						{
							Target: &selector.Candidate{
								NamespaceID: 1,
								Namespace:   "project1",
								Tags:        []string{"v1"},
								Labels:      nil,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Handle(tt.args.data)
			if tt.wantErr {
				require.NotNil(t, err, "Error: %s", err)
				return
			}
			assert.Nil(t, err)
		})
	}

}

func TestRetentionHandler_IsStateful(t *testing.T) {
	handler := &RetentionHandler{}
	assert.False(t, handler.IsStateful())
}
