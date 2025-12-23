package tokenb

import "time"


func (s *Scheduler) Run() {
	stop := time.After(s.duration)

	for {
		select {
		case <-stop:
			return
		default:
			s.bucket.Take()                     // ⬅ токен = право на запрос
			req := s.generator.Next()
			s.executor.jobs <- req
		}
	}
}
