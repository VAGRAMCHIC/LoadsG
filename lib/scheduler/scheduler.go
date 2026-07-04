package scheduler

import "context"

type Scheduler interface {
	// ScheduleJobs создаёт события (pending) для всех заданий, у которых ещё нет событий.
	ScheduleJobs(ctx context.Context) error

	// ProcessPending находит задания, у которых start_time <= now, и переводит их события из pending в running.
	ProcessPending(ctx context.Context) error

	// ExecuteRunning запускает генераторы для всех событий со статусом running.
	// Выполняется асинхронно, статус обновляется по завершении.
	ExecuteRunning(ctx context.Context) error

	// Run – фоновый цикл, который периодически вызывает ProcessPending и ExecuteRunning.
	Run(ctx context.Context) error
}
