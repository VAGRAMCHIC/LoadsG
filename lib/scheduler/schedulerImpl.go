package scheduler

import (
	"context"
	"loadsg/lib/generators"
	"loadsg/lib/model"
	"loadsg/lib/repository"
	"log"
	"sync"
	"time"
)

type scheduler struct {
	loadRepo   repository.LoadRepository
	eventRepo  repository.EventRepository
	httpRepo   repository.HttpLoadRepository
	registry   *generators.Registry
	interval   time.Duration // как часто сканировать БД
	workerWG   sync.WaitGroup
	cancelFunc context.CancelFunc
}

func NewScheduler(
	lr repository.LoadRepository,
	ev repository.EventRepository,
	hr repository.HttpLoadRepository,
	reg *generators.Registry,
	interval time.Duration,
) Scheduler {
	return &scheduler{
		loadRepo:  lr,
		eventRepo: ev,
		httpRepo:  hr,
		registry:  reg,
		interval:  interval,
	}
}

// ScheduleJobs – создаём события для всех load_job, у которых ещё нет ни одного события.
func (s *scheduler) ScheduleJobs(ctx context.Context) error {
	// Получаем все задания
	jobs, err := s.loadRepo.ScanLJob(ctx)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		// Проверяем, есть ли уже события для этого load_job_id
		// Можно сделать запрос, но пока упростим: попробуем создать, если дубликат – ошибка игнорируется.
		// Лучше сначала проверить существование, но для простоты создадим, если нет.
		events, err := s.eventRepo.ScanEvents(ctx, "pending")
		if err != nil {
			log.Printf("ScheduleJobs: scan events error: %v", err)
			continue
		}
		exists := false
		for _, e := range events {
			if e.LoadJobId == job.Id {
				exists = true
				break
			}
		}
		if !exists {
			// Создаём событие со статусом pending
			event := &model.Event{
				LoadJobId: job.Id,
				Status:    "pending",
			}
			_, err := s.eventRepo.Create(ctx, event)
			if err != nil {
				log.Printf("ScheduleJobs: failed to create event for job %s: %v", job.Id, err)
			}
		}
	}
	return nil
}

// ProcessPending – ищем задания, чьё время старта <= now, и переводим их события в running.
func (s *scheduler) ProcessPending(ctx context.Context) error {
	// Получаем задания, которые должны стартовать (start_time <= now)
	// В репозитории есть метод ScanClosestLJob, но он использует интервал 10 секунд.
	// Сделаем более точный запрос.
	// Временно используем ScanClosestLJob, потом можно улучшить.
	jobs, err := s.loadRepo.ScanClosestLJob(ctx)
	if err != nil {
		return err
	}

	// Для каждого задания находим pending события и меняем статус
	for _, job := range jobs {
		// Получаем все события этого задания (можно фильтровать по load_job_id, но у нас нет такого метода)
		// Добавим метод в репозиторий или просто просканируем все pending.
		pendingEvents, err := s.eventRepo.ScanEvents(ctx, "pending")
		if err != nil {
			log.Printf("ProcessPending: scan pending events error: %v", err)
			continue
		}
		for _, ev := range pendingEvents {
			if ev.LoadJobId == job.Id {
				ev.Status = "running"
				if err := s.eventRepo.UpdateEvent(ctx, &ev); err != nil {
					log.Printf("ProcessPending: failed to update event %s to running: %v", ev.Id, err)
				} else {
					log.Printf("ProcessPending: event %s for job %s set to running", ev.Id, job.Id)
				}
			}
		}
	}
	return nil
}

// ExecuteRunning – запускает генераторы для всех событий со статусом running.
// Запускает каждый в отдельной горутине, обновляет статус по завершении.
func (s *scheduler) ExecuteRunning(ctx context.Context) error {
	events, err := s.eventRepo.ScanEvents(ctx, "running")
	if err != nil {
		return err
	}

	for _, ev := range events {
		// Загружаем задание
		loadJob, err := s.loadRepo.GetById(ctx, ev.LoadJobId)
		if err != nil {
			log.Printf("ExecuteRunning: cannot load job %s: %v", ev.LoadJobId, err)
			s.eventRepo.MarkFailed(ctx, ev.Id)
			continue
		}

		// Выбираем генератор по типу задания (loadJob.Type)
		generator, err := s.registry.Get(loadJob.Type) // например "constant_http", "fixed_http" и т.д.
		if err != nil {
			log.Printf("ExecuteRunning: no generator for type %s: %v", loadJob.Type, err)
			s.eventRepo.MarkFailed(ctx, ev.Id)
			continue
		}

		// Запускаем генератор в отдельной горутине
		s.workerWG.Add(1)
		go func(event model.Event, job model.LoadJob, gen generators.Generator) {
			defer s.workerWG.Done()

			log.Printf("Starting generator for event %s, job %s", event.Id, job.Id)

			// Создаём контекст с возможностью отмены (можно добавить в будущем)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Запускаем генератор (он должен принимать контекст и job)
			err := gen.Run(ctx, job)
			if err != nil && err != context.Canceled {
				log.Printf("Generator for event %s finished with error: %v", event.Id, err)
				if err2 := s.eventRepo.MarkFailed(ctx, event.Id); err2 != nil {
					log.Printf("Cannot mark event %s failed: %v", event.Id, err2)
				}
				return
			}

			// Если завершился без ошибок или был отменён – помечаем как done
			if err == nil || err == context.Canceled {
				if err2 := s.eventRepo.MarkDone(ctx, event.Id); err2 != nil {
					log.Printf("Cannot mark event %s done: %v", event.Id, err2)
				} else {
					log.Printf("Event %s completed successfully", event.Id)
				}
			}
		}(ev, *loadJob, generator)
	}

	return nil
}

// Run – бесконечный цикл, периодически запускающий все этапы.
func (s *scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Первоначальный запуск
	if err := s.ScheduleJobs(ctx); err != nil {
		log.Printf("Initial schedule error: %v", err)
	}
	if err := s.ProcessPending(ctx); err != nil {
		log.Printf("Initial process pending error: %v", err)
	}
	if err := s.ExecuteRunning(ctx); err != nil {
		log.Printf("Initial execute running error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopping...")
			s.workerWG.Wait() // дожидаемся завершения всех запущенных генераторов
			return ctx.Err()
		case <-ticker.C:
			// Каждый тик выполняем шаги
			if err := s.ScheduleJobs(ctx); err != nil {
				log.Printf("Schedule error: %v", err)
			}
			if err := s.ProcessPending(ctx); err != nil {
				log.Printf("Process pending error: %v", err)
			}
			if err := s.ExecuteRunning(ctx); err != nil {
				log.Printf("Execute running error: %v", err)
			}
		}
	}
}
