// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package or

import (
	"github.com/goharbor/harbor/src/common/utils/log"
	"github.com/goharbor/harbor/src/pkg/retention/policy/action"
	"github.com/goharbor/harbor/src/pkg/retention/policy/alg"
	"github.com/goharbor/harbor/src/pkg/retention/policy/rule"
	"github.com/goharbor/harbor/src/pkg/retention/res"
	"github.com/pkg/errors"
	"sync"
)

// processor to handle the rules with OR mapping ways
type processor struct {
	performer action.Performer
	// keep evaluator and its related selector if existing
	// attentions here, the selectors can be empty/nil, that means match all "**"
	evaluators map[*rule.Evaluator][]res.Selector
}

// New processor
func New() alg.Processor {
	return &processor{}
}

// Process the candidates with the rules
func (p *processor) Process(artifacts []*res.Candidate) ([]*res.Result, error) {
	if len(artifacts) == 0 {
		log.Debug("no artifacts to retention")
		return make([]*res.Result, 0), nil
	}

	var (
		// collect errors by wrapping
		err error
		// collect results
		retained = make([]*res.Candidate, 0)
	)

	// for sync
	resChan := make(chan []*res.Candidate, 1)
	// handle error
	errChan := make(chan error, 1)
	// control chan
	done := make(chan bool, 1)

	defer func() {
		// signal the result listener loop exit
		done <- true
	}()

	// go routine for receiving results/error
	go func() {
		for {
			select {
			case retainedOnes := <-resChan:
				retained = append(retained, retainedOnes...)
			case e := <-errChan:
				if err == nil {
					err = errors.Wrap(e, "artifact processing error")
				} else {
					err = errors.Wrap(e, err.Error())
				}
			case <-done:
				// exit
				return
			}
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(len(p.evaluators))

	for eva, selectors := range p.evaluators {
		var evaluator = *eva

		go func(evaluator rule.Evaluator, selectors []res.Selector) {
			var (
				processed []*res.Candidate
				err       error
			)

			defer func() {
				wg.Done()
			}()

			// init
			// pass array copy to the selector
			processed = append(processed, artifacts...)

			if len(selectors) > 0 {
				// selecting artifacts one by one
				// `&&` mappings
				for _, s := range selectors {
					if processed, err = s.Select(processed); err != nil {
						errChan <- err
						return
					}
				}
			}

			if processed, err = evaluator.Process(processed); err != nil {
				errChan <- err
				return
			}

			if len(processed) > 0 {
				// Pass to the outside
				resChan <- processed
			}
		}(evaluator, selectors)
	}

	// waiting for all the rules are evaluated
	wg.Wait()

	if err != nil {
		return nil, err
	}

	return p.performer.Perform(retained)
}

// AddEvaluator appends a rule evaluator for processing
func (p *processor) AddEvaluator(evaluator rule.Evaluator, selectors []res.Selector) {
	if evaluator != nil {
		p.evaluators[&evaluator] = selectors
	}
}

// SetPerformer sets a action performer to the processor
func (p *processor) SetPerformer(performer action.Performer) {
	p.performer = performer
}
