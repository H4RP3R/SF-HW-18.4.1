package semaphore

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSemaphoreAcquire(t *testing.T) {
	tests := []struct {
		gorAmount int
		semLimit  int
		timeout   time.Duration
	}{
		{gorAmount: 1, semLimit: 1, timeout: 1 * time.Second},
		{gorAmount: 2, semLimit: 2, timeout: 1 * time.Second},
		{gorAmount: 2, semLimit: 5, timeout: 1 * time.Second},
		{gorAmount: 3, semLimit: 2, timeout: 1 * time.Second},
		{gorAmount: 24, semLimit: 5, timeout: 1 * time.Second},
		{gorAmount: 200, semLimit: 51, timeout: 1 * time.Second},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("goroutines: %d, limit: %d", tt.gorAmount, tt.semLimit)
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			errChan := make(chan error, tt.gorAmount)
			s := NewSemaphore(tt.timeout, tt.semLimit)
			safeAttemptsNum := min(tt.gorAmount, tt.semLimit)
			// Successful attempts to acquire
			for i := 0; i < safeAttemptsNum; i++ {
				wg.Add(1)
				go func(s *Semaphore, errChan chan<- error, wg *sync.WaitGroup) {
					err := s.Acquire()
					errChan <- err
					wg.Done()
				}(s, errChan, &wg)
			}
			wg.Wait()
			close(errChan)
			for err := range errChan {
				if err != nil {
					t.Errorf("didn't expect an error")
				}
			}
			// Unsuccessful attempts to acquire
			errChan = make(chan error, tt.gorAmount)
			for i := safeAttemptsNum; i < tt.gorAmount; i++ {
				wg.Add(1)
				go func(s *Semaphore, errChan chan<- error, wg *sync.WaitGroup) {
					err := s.Acquire()
					errChan <- err
					wg.Done()
				}(s, errChan, &wg)
			}
			wg.Wait()
			close(errChan)
			for err := range errChan {
				if err == nil {
					t.Errorf("expected an error")
				}
			}
		})
	}
}

func TestSemaphoreRelease(t *testing.T) {
	args := struct {
		semLimit int
		timeout  time.Duration
	}{
		semLimit: 3,
		timeout:  2 * time.Second,
	}

	name := fmt.Sprintf("limit: %d, timeout: %v", args.semLimit, args.timeout)
	t.Run(name, func(t *testing.T) {
		s := NewSemaphore(args.timeout, args.semLimit)
		for i := 0; i < args.semLimit; i++ {
			s.Acquire()
		}
		for i := 0; i < args.semLimit; i++ {
			err := s.Release()
			if err != nil {
				t.Errorf("didn't expect an error")
			}
		}

		if err := s.Release(); err == nil {
			t.Errorf("expected an error")
		}
	})
}

func TestSemaphoreReleaseTimeout(t *testing.T) {
	args := struct {
		semLimit int
		timeout  time.Duration
	}{
		semLimit: 3,
		timeout:  1 * time.Second,
	}

	name := fmt.Sprintf("limit: %d, timeout: %v", args.semLimit, args.timeout)
	t.Run(name, func(t *testing.T) {
		s := NewSemaphore(args.timeout, args.semLimit)
		for i := 0; i < args.semLimit; i++ {
			s.Acquire()
		}

		// Wait for the timeout to expire
		time.Sleep(args.timeout + 1*time.Second)

		if err := s.Release(); err != nil {
			t.Errorf("expected an error due to timeout")
		}
	})
}
