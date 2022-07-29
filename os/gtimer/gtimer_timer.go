package gtimer

import (
	gtype2 "github.com/Ravior/gserver/util/gcontainer/gtype"
	"time"
)

func New() *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: gtype2.NewInt(StatusRunning),
		ticks:  gtype2.NewInt64(),
	}
	go t.loop()
	return t
}

// Add adds a timing job to the timer, which runs in interval of <interval>.
func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, false, -1, StatusReady)
}

// AddEntry adds a timing job to the timer with detailed parameters.
//
// The parameter <interval> specifies the running interval of the job.
//
// The parameter <singleton> specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when its a singleton mode job.
//
// The parameter <times> specifies limit for the job running times, which means the job
// exits if its run times exceeds the <times>.
//
// The parameter <status> specifies the job status when it's firstly added to the timer.
func (t *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return t.createEntry(interval, job, singleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, -1, StatusReady)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, 1, StatusReady)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return t.createEntry(interval, job, true, times, StatusReady)
}

// DelayAdd adds a timing job after delay of <interval> duration.
// Also see Add.
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.Add(interval, job)
	})
}

// DelayAddEntry adds a timing job after delay of <interval> duration.
// Also see AddEntry.
func (t *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	t.AddOnce(delay, func() {
		t.AddEntry(interval, job, singleton, times, status)
	})
}

// DelayAddSingleton adds a timing job after delay of <interval> duration.
// Also see AddSingleton.
func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddSingleton(interval, job)
	})
}

// DelayAddOnce adds a timing job after delay of <interval> duration.
// Also see AddOnce.
func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddOnce(interval, job)
	})
}

// After adds a timing job after delay duration.
// Also see AddOnce.
func (t *Timer) After(delay time.Duration, job JobFunc) {
	t.DelayAddOnce(delay, 0, job)
}

// DelayAddTimes adds a timing job after delay of <interval> duration.
// Also see AddTimes.
func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	t.AddOnce(delay, func() {
		t.AddTimes(interval, times, job)
	})
}

// Start starts the timer.
func (t *Timer) Start() {
	t.status.Set(StatusRunning)
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.status.Set(StatusStopped)
}

// Clear the timer
func (t *Timer) Clear() {
	t.status.Set(StatusStopped)
	t.queue = newPriorityQueue()
	t.status.Set(StatusRunning)
}

// Close closes the timer.
func (t *Timer) Close() {
	t.queue = newPriorityQueue()
	t.status.Set(StatusClosed)
}

// createEntry creates and adds a timing job to the timer.
func (t *Timer) createEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	var (
		infinite = false
	)
	if times <= 0 {
		infinite = true
	}
	var (
		intervalTicksOfJob = int64(interval / defaultInterval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	// 如果ticks为空，则返回空对象
	if t.ticks == nil {
		return nil
	}
	var (
		nextTicks = t.ticks.Val() + intervalTicksOfJob
		entry     = &Entry{
			job:       job,
			timer:     t,
			ticks:     intervalTicksOfJob,
			times:     gtype2.NewInt(times),
			status:    gtype2.NewInt(status),
			singleton: gtype2.NewBool(singleton),
			nextTicks: gtype2.NewInt64(nextTicks),
			infinite:  gtype2.NewBool(infinite),
		}
	)
	t.queue.Push(entry, nextTicks)
	return entry
}
