package rule

import (
	"github.com/goharbor/harbor/src/common/utils/test"
	"github.com/goharbor/harbor/src/pkg/art"
	"github.com/goharbor/harbor/src/pkg/immutabletag"
	"github.com/goharbor/harbor/src/pkg/immutabletag/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// MatchTestSuite ...
type MatchTestSuite struct {
	suite.Suite
	t       *testing.T
	assert  *assert.Assertions
	require *require.Assertions
	ctr     immutabletag.APIController
	ruleID  int64
}

// SetupSuite ...
func (s *MatchTestSuite) SetupSuite() {
	test.InitDatabaseFromEnv()
	s.t = s.T()
	s.assert = assert.New(s.t)
	s.require = require.New(s.t)
	s.ctr = immutabletag.ImmuCtr
}

func (s *MatchTestSuite) TestImmuMatch() {
	rule := &model.Metadata{
		ID:        1,
		ProjectID: 2,
		Priority:  1,
		Template:  "latestPushedK",
		Action:    "immuablity",
		TagSelectors: []*model.Selector{
			{
				Kind:       "doublestar",
				Decoration: "matches",
				Pattern:    "release-[\\d\\.]+",
			},
		},
		ScopeSelectors: map[string][]*model.Selector{
			"repository": {
				{
					Kind:       "doublestar",
					Decoration: "matches",
					Pattern:    "**",
				},
			},
		},
	}

	id, err := s.ctr.CreateImmutableRule(rule)
	s.ruleID = id
	s.require.NotNil(err)

	c1 := art.Candidate{
		NamespaceID:  2,
		Namespace:    "immutable",
		Repository:   "redis",
		Tag:          "release-1.10",
		Kind:         art.Image,
		PushedTime:   time.Now().Unix() - 3600,
		PulledTime:   time.Now().Unix(),
		CreationTime: time.Now().Unix() - 7200,
		Labels:       []string{"label1", "label4", "label5"},
	}

	match := NewRuleMatcher(2)
	isMatch, err := match.Match(c1)
	s.require.Equal(isMatch, true)
	s.require.Nil(err)

	c2 := art.Candidate{
		NamespaceID:  2,
		Namespace:    "immutable",
		Repository:   "redis",
		Tag:          "1.10",
		Kind:         art.Image,
		PushedTime:   time.Now().Unix() - 3600,
		PulledTime:   time.Now().Unix(),
		CreationTime: time.Now().Unix() - 7200,
		Labels:       []string{"label1", "label4", "label5"},
	}

	isMatch, err = match.Match(c2)
	s.require.Equal(isMatch, false)
	s.require.Nil(err)
}

// TearDownSuite clears env for test suite
func (s *MatchTestSuite) TearDownSuite() {
	err := s.ctr.DeleteImmutableRule(s.ruleID)
	require.NoError(s.T(), err, "delete immutable")
}
