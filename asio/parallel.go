package asio

import (
	"sync"
)

type Parallel struct {
	tasks   chan parallelTask
	threads chan struct{}
	stop    chan struct{}
	wg      sync.WaitGroup
}

type parallelTask struct {
	task  func() error
	errCh chan error
}

func NewParallel(tasks, threads int) *Parallel {
	p := &Parallel{
		tasks:   make(chan parallelTask, tasks),
		threads: make(chan struct{}, threads),
		stop:    make(chan struct{}),
	}

	p.wg.Add(1)
	go p.Run()
	return p
}

func (p *Parallel) Put(task func() error, errCh chan error) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		select {
		case p.tasks <- parallelTask{
			task:  task,
			errCh: errCh,
		}:
		case <-p.stop:
			break
		}
	}()
}

func (p *Parallel) Run() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.tasks:
			select {
			case p.threads <- struct{}{}:
				p.wg.Add(1)
				go func() {
					defer func() {
						p.wg.Done()
						<-p.threads
					}()

					err := task.task()
					if task.errCh != nil {
						select {
						case task.errCh <- err:
						case <-p.stop:
							break
						}
					}
				}()

			case <-p.stop:
				return
			}

		case <-p.stop:
			return
		}
	}
}

func (p *Parallel) Stop() {
	close(p.stop)
	p.wg.Wait()
}
