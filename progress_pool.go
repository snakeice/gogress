package gogress

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/snakeice/gogress/writer"
)

func NewPool() *Pool {
	ctx, fn := context.WithCancel(context.Background())
	return &Pool{
		bars:        []*Progress{},
		ctx:         ctx,
		stopFn:      fn,
		isRunning:   false,
		RefreshRate: time.Second / 45,
		writer:      writer.New(os.Stdout),
	}
}

type Pool struct {
	bars        []*Progress
	ctx         context.Context
	stopFn      context.CancelFunc
	RefreshRate time.Duration
	writer      *writer.Writer
	isRunning   bool
	mu          sync.Mutex
	finishOnce  sync.Once
}

func (p *Pool) AddBar(bar *Progress) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bars = append(p.bars, bar)
	bar.setWriter(p.writer)
	bar.pooled = true
	if p.isRunning {
		bar.Start()
	}
	for id, _bar := range p.bars {
		_bar.ID = id
	}

	return len(p.bars) - 1
}

func remove(slice []*Progress, id int) []*Progress {
	return append(slice[:id], slice[id+1:]...)
}

func (p *Pool) RemoveBar(bar *Progress) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bars = remove(p.bars, bar.ID)
	for id, _bar := range p.bars {
		_bar.ID = id
	}
	return len(p.bars) - 1
}

func (p *Pool) NewBarDef() *Progress {
	bar := NewDef()
	p.AddBar(bar)
	return bar
}

func (p *Pool) NewBar(max int) *Progress {
	return p.NewBar64(int64(max))
}

func (p *Pool) NewBar64(max int64) *Progress {
	bar := New64(max)
	p.AddBar(bar)
	return bar
}

func (p *Pool) refresher() {
	for {
		select {
		case <-p.ctx.Done():
			p.isRunning = false
			return
		case <-time.After(p.RefreshRate):
			p.Update()
		}
	}
}

func (p *Pool) Start() *Pool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isRunning {
		return p
	}
	p.isRunning = true
	for _, bar := range p.bars {
		bar.Start()
	}
	p.Update()
	go p.refresher()
	return p
}

func (p *Pool) Update() {
	for _, bar := range p.bars {
		bar.Update()
	}
	p.writer.Flush(len(p.bars))
}

func (p *Pool) FinishAll() {
	p.finishOnce.Do(func() {
		p.stopFn()

		for _, bar := range p.bars {
			bar.Finish()
		}

		p.mu.Lock()
		defer p.mu.Unlock()
		p.Update()
		p.isRunning = false
	})
}

func (p *Pool) IsFinished() bool {
	var result = true
	for _, bar := range p.bars {
		result = bar.IsFinished() && result
		if !result {
			break
		}
	}
	return result
}

func (p *Pool) Wait() {
	for !p.IsFinished() {
		time.Sleep(time.Millisecond * 100)
	}
}
