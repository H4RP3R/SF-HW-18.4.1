// Измените приведенный выше код по реализации двоичного семафора так,
// чтобы он описывал не двоичный семафор, а семафор подсчёта.

package semaphore

import (
	"fmt"
	"time"
)

// Semaphore — структура двоичного семафора
type Semaphore struct {
	// Семафор — абстрактный тип данных,
	// в нашем случае в основе его лежит канал
	sem chan int
	// Время ожидания основных операций с семафором, чтобы не блокировать
	// операции с ним навечно (необязательное требование, зависит от
	// нужд программы)
	timeout time.Duration
}

// Acquire — метод захвата семафора
func (s *Semaphore) Acquire() error {
	select {
	case s.sem <- 0:
		return nil
	case <-time.After(s.timeout):
		return fmt.Errorf("не удалось захватить семафор")
	}
}

// Release — метод освобождения семафора
func (s *Semaphore) Release() error {
	select {
	case <-s.sem:
		return nil
	case <-time.After(s.timeout):
		return fmt.Errorf("не удалось освободить семафор")
	}
}

// NewSemaphore — функция создания семафора
func NewSemaphore(timeout time.Duration, limit int) *Semaphore {
	return &Semaphore{
		sem:     make(chan int, limit),
		timeout: timeout,
	}
}
