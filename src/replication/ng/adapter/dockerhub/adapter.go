package dockerhub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/goharbor/harbor/src/common/utils/log"
	adp "github.com/goharbor/harbor/src/replication/ng/adapter"
	"github.com/goharbor/harbor/src/replication/ng/model"
)

func init() {
	if err := adp.RegisterFactory(model.RegistryTypeDockerHub, func(registry *model.Registry) (adp.Adapter, error) {
		client, err := NewClient(registry)
		if err != nil {
			return nil, err
		}

		return &adapter{
			client:   client,
			registry: registry,
			DefaultImageRegistry: adp.NewDefaultImageRegistry(&model.Registry{
				Name:       registry.Name,
				URL:        registryURL,
				Credential: registry.Credential,
				Insecure:   registry.Insecure,
			}),
		}, nil
	}); err != nil {
		log.Errorf("Register adapter factory for %s error: %v", model.RegistryTypeDockerHub, err)
		return
	}
	log.Infof("Factory for adapter %s registered", model.RegistryTypeDockerHub)
}

type adapter struct {
	*adp.DefaultImageRegistry
	registry *model.Registry
	client   *Client
}

// Ensure '*adapter' implements interface 'Adapter'.
var _ adp.Adapter = (*adapter)(nil)

// Info returns information of the registry
func (a *adapter) Info() (*model.RegistryInfo, error) {
	return &model.RegistryInfo{
		Type: model.RegistryTypeDockerHub,
		SupportedResourceTypes: []model.ResourceType{
			model.ResourceTypeRepository,
		},
		SupportedResourceFilters: []*model.FilterStyle{
			{
				Type:  model.FilterTypeName,
				Style: model.FilterStyleTypeText,
			},
		},
		SupportedTriggers: []model.TriggerType{
			model.TriggerTypeManual,
		},
	}, nil
}

// HealthCheck checks health status of the registry
func (a *adapter) HealthCheck() (model.HealthStatus, error) {
	return model.Healthy, nil
}

// ListNamespaces lists namespaces from DockerHub with the provided query conditions.
func (a *adapter) ListNamespaces(query *model.NamespaceQuery) ([]*model.Namespace, error) {
	resp, err := a.client.Do(http.MethodGet, listNamespacePath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		log.Errorf("list namespace error: %s", string(body))
		return nil, fmt.Errorf("%s", string(body))
	}

	namespaces := NamespacesResp{}
	err = json.Unmarshal(body, &namespaces)
	if err != nil {
		return nil, err
	}
	var result []*model.Namespace
	for _, ns := range namespaces.Namespaces {
		// If query set, skip the namespace that doesn't match the query.
		if query != nil && len(query.Name) > 0 && strings.Index(ns, query.Name) != -1 {
			continue
		}

		result = append(result, &model.Namespace{
			Name: ns,
		})
	}
	return result, nil
}

// CreateNamespace creates a new namespace in DockerHub
func (a *adapter) CreateNamespace(namespace *model.Namespace) error {
	ns, err := a.getNamespace(namespace.Name)
	if err != nil {
		return fmt.Errorf("check existence of namespace '%s' error: %v", namespace.Name, err)
	}

	// If the namespace already exist, return succeeded directly.
	if ns != nil {
		log.Infof("Namespace %s already exist in DockerHub, skip it.", namespace.Name)
		return nil
	}

	req := &NewOrgReq{
		Name:     namespace.Name,
		FullName: namespace.GetStringMetadata(metadataKeyFullName, namespace.Name),
		Company:  namespace.GetStringMetadata(metadataKeyCompany, namespace.Name),
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(http.MethodPost, createNamespacePath, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		log.Errorf("create namespace error: %d -- %s", resp.StatusCode, string(body))
		return fmt.Errorf("%d -- %s", resp.StatusCode, body)
	}

	return nil
}

// GetNamespace gets a namespace from DockerHub.
func (a *adapter) GetNamespace(namespace string) (*model.Namespace, error) {
	return &model.Namespace{
		Name: namespace,
	}, nil
}

// getNamespace get namespace from DockerHub, if the namespace not found, two nil would be returned.
func (a *adapter) getNamespace(namespace string) (*model.Namespace, error) {
	resp, err := a.client.Do(http.MethodGet, getNamespacePath(namespace), nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode/100 != 2 {
		log.Errorf("create namespace error: %d -- %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("%d -- %s", resp.StatusCode, body)
	}

	return &model.Namespace{
		Name: namespace,
	}, nil
}

// FetchImages fetches images
func (a *adapter) FetchImages(namespaces []string, filters []*model.Filter) ([]*model.Resource, error) {
	var repos []Repo
	nameFilter := a.getFilter(model.FilterTypeName, filters)
	for _, ns := range namespaces {
		name := ""
		if nameFilter != nil {
			v, ok := nameFilter.Value.(string)
			if !ok {
				msg := fmt.Sprintf("expect name filter value to be string, but got: %v", nameFilter.Value)
				log.Error(msg)
				return nil, errors.New(msg)
			}
			name = v
		}

		page := 1
		pageSize := 100
		for {
			pageRepos, err := a.getRepos(ns, name, page, pageSize)
			if err != nil {
				return nil, fmt.Errorf("get repos for namespace '%s' from DockerHub error: %v", ns, err)
			}
			repos = append(repos, pageRepos.Repos...)

			if len(pageRepos.Next) == 0 {
				break
			}
			page++
		}
	}

	log.Infof("%d repos found for namespaces: %v", len(repos), namespaces)
	var resources []*model.Resource
	// TODO(ChenDe): Get tags for repos in parallel
	for _, repo := range repos {
		var tags []string
		page := 1
		pageSize := 100
		for {
			pageTags, err := a.getTags(repo.Namespace, repo.Name, page, pageSize)
			if err != nil {
				return nil, fmt.Errorf("get tags for repo '%s/%s' from DockerHub error: %v", repo.Namespace, repo.Name, err)
			}
			for _, t := range pageTags.Tags {
				tags = append(tags, t.Name)
			}

			if len(pageTags.Next) == 0 {
				break
			}
			page++
		}

		// If the repo has no tags, skip it
		if len(tags) == 0 {
			continue
		}

		resources = append(resources, &model.Resource{
			Type:     model.ResourceTypeRepository,
			Registry: a.registry,
			Metadata: &model.ResourceMetadata{
				Namespace: repo.Namespace,
				Name:      fmt.Sprintf("%s/%s", repo.Namespace, repo.Name),
				Vtags:     tags,
			},
		})
	}

	return resources, nil
}

// getRepos gets a page of repos from DockerHub
func (a *adapter) getRepos(namespace, name string, page, pageSize int) (*ReposResp, error) {
	resp, err := a.client.Do(http.MethodGet, listReposPath(namespace, name, page, pageSize), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		log.Errorf("list repos error: %d -- %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("%d -- %s", resp.StatusCode, body)
	}

	repos := &ReposResp{}
	err = json.Unmarshal(body, repos)
	if err != nil {
		return nil, fmt.Errorf("unmarshal repos list %s error: %v", string(body), err)
	}

	return repos, nil
}

// getTags gets a page of tags for a repo from DockerHub
func (a *adapter) getTags(namespace, repo string, page, pageSize int) (*TagsResp, error) {
	resp, err := a.client.Do(http.MethodGet, listTagsPath(namespace, repo, page, pageSize), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		log.Errorf("list tags error: %d -- %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("%d -- %s", resp.StatusCode, body)
	}

	tags := &TagsResp{}
	err = json.Unmarshal(body, tags)
	if err != nil {
		return nil, fmt.Errorf("unmarshal tags list %s error: %v", string(body), err)
	}

	return tags, nil
}

// getFilter gets specific type filter from filters list.
func (a *adapter) getFilter(filterType model.FilterType, filters []*model.Filter) *model.Filter {
	for _, f := range filters {
		if f.Type == filterType {
			return f
		}
	}
	return nil
}
