package handler

import (
	"github.com/goharbor/harbor/src/api/event"
	"github.com/goharbor/harbor/src/api/event/handler/auditlog"
	"github.com/goharbor/harbor/src/api/event/handler/replication"
	"github.com/goharbor/harbor/src/api/event/handler/webhook/artifact"
	"github.com/goharbor/harbor/src/api/event/handler/webhook/chart"
	"github.com/goharbor/harbor/src/api/event/handler/webhook/quota"
	"github.com/goharbor/harbor/src/api/event/handler/webhook/scan"
	"github.com/goharbor/harbor/src/pkg/notifier"
)

func init() {
	// notification
	notifier.Subscribe(event.TopicPushArtifact, &artifact.Handler{})
	notifier.Subscribe(event.TopicPullArtifact, &artifact.Handler{})
	notifier.Subscribe(event.TopicDeleteArtifact, &artifact.Handler{})
	notifier.Subscribe(event.TopicUploadChart, &chart.Handler{})
	notifier.Subscribe(event.TopicDeleteChart, &chart.Handler{})
	notifier.Subscribe(event.TopicDownloadChart, &chart.Handler{})
	notifier.Subscribe(event.TopicQuotaExceed, &quota.Handler{})
	notifier.Subscribe(event.TopicQuotaWarning, &quota.Handler{})
	notifier.Subscribe(event.TopicScanningFailed, &scan.Handler{})
	notifier.Subscribe(event.TopicScanningCompleted, &scan.Handler{})
	notifier.Subscribe(event.TopicDeleteArtifact, &scan.DelArtHandler{})

	// replication
	notifier.Subscribe(event.TopicPushArtifact, &replication.Handler{})
	notifier.Subscribe(event.TopicDeleteArtifact, &replication.Handler{})
	notifier.Subscribe(event.TopicCreateTag, &replication.Handler{})

	// audit logs
	notifier.Subscribe(event.TopicPushArtifact, &auditlog.Handler{})
	notifier.Subscribe(event.TopicPullArtifact, &auditlog.Handler{})
	notifier.Subscribe(event.TopicDeleteArtifact, &auditlog.Handler{})
	notifier.Subscribe(event.TopicCreateProject, &auditlog.Handler{})
	notifier.Subscribe(event.TopicDeleteProject, &auditlog.Handler{})
	notifier.Subscribe(event.TopicDeleteRepository, &auditlog.Handler{})
	notifier.Subscribe(event.TopicCreateTag, &auditlog.Handler{})
	notifier.Subscribe(event.TopicDeleteTag, &auditlog.Handler{})
}
