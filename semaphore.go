package anomalia

type empty struct{}
type semaphore chan empty

func (s semaphore) Lock() {
	s.acquire(1)
}

func (s semaphore) Unlock() {
	s.release(1)
}

func (s semaphore) Wait(n int) {
	s.acquire(n)
}

func (s semaphore) Signal() {
	s.release(1)
}

func (s semaphore) acquire(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

func (s semaphore) release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}
