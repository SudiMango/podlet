package service

import (
	"sync"

	"github.com/SudiMango/podlet/internal/contract"
)

type JobObserver struct {
	mu   sync.Mutex
	subs map[string]map[chan contract.JobResponse]struct{}
}

func NewJobObserver() *JobObserver {
	return &JobObserver{subs: map[string]map[chan contract.JobResponse]struct{}{}}
}

func (o *JobObserver) Subscribe(jobId string) (<-chan contract.JobResponse, func()) {
	ch := make(chan contract.JobResponse, 8)

	o.mu.Lock()
	if o.subs[jobId] == nil {
		o.subs[jobId] = map[chan contract.JobResponse]struct{}{}
	}
	o.subs[jobId][ch] = struct{}{}
	o.mu.Unlock()

	return ch, func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		delete(o.subs[jobId], ch)
		if len(o.subs[jobId]) == 0 {
			delete(o.subs, jobId)
		}
	}
}

func (o *JobObserver) Publish(jobId string, job contract.JobResponse) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for ch := range o.subs[jobId] {
		select {
		case ch <- job:
		default:
		}
	}
}
