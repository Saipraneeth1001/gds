package workerpool

import (
	"log"
	"sync"
)

type WorkerPool struct {
	jobs chan int
	wg   sync.WaitGroup
}

func InitWorker(numWorkers int) *WorkerPool {
	pool := WorkerPool{
		jobs: make(chan int, 10),
	}

	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go pool.work(i)
	}

	return &pool
}

func (w *WorkerPool) work(workerID int) {
	defer w.wg.Done()

	for job := range w.jobs {
		log.Printf("worker %d processing job", workerID)
		if job%2 == 0 {
			log.Println("even")
		} else {
			log.Println("odd")
		}
	}
}

func (w *WorkerPool) Submit(num int) {
	w.jobs <- num
}

func (w *WorkerPool) Shutdown() {
	close(w.jobs)
	w.wg.Wait()
}
