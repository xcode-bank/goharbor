package scheduler

import (
	"encoding/json"
	"testing"

	"github.com/goharbor/harbor/src/common/job/models"
	"github.com/goharbor/harbor/src/replication/ng/model"
)

var scheduler = &defaultScheduler{
	client: TestClient{},
}

type TestClient struct {
}

func (client TestClient) SubmitJob(*models.JobData) (string, error) {
	return "submited-uuid", nil
}
func (client TestClient) GetJobLog(uuid string) ([]byte, error) {
	return []byte("job log"), nil
}
func (client TestClient) PostAction(uuid, action string) error {
	return nil
}

func TestPreprocess(t *testing.T) {
	items, err := generateData()
	if err != nil {
		t.Error(err)
	}
	for _, item := range items {
		content, err := json.Marshal(item)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(content))
	}

}

func TestStop(t *testing.T) {
	err := scheduler.Stop("id")
	if err != nil {
		t.Error(err)
	}
}

func generateData() ([]*ScheduleItem, error) {
	srcResource := &model.Resource{
		Metadata: &model.ResourceMetadata{
			Namespace: "namespace1",
			Vtags:     []string{"latest"},
			Labels:    []string{"latest"},
		},
		Registry: &model.Registry{
			Credential: &model.Credential{},
		},
	}
	destResource := &model.Resource{
		Metadata: &model.ResourceMetadata{
			Namespace: "namespace2",
			Vtags:     []string{"v1", "v2"},
			Labels:    []string{"latest"},
		},
		Registry: &model.Registry{
			Credential: &model.Credential{},
		},
	}
	items, err := scheduler.Preprocess([]*model.Resource{srcResource}, []*model.Resource{destResource})
	return items, err
}
