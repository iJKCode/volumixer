package event

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

type EventA struct{}
type EventB struct{}
type EventC struct{}

func handler[E any](counter *int) Handler[E] {
	return Func(func(ctx context.Context, event E) {
		*counter++
	})
}

func runCount(bus *Bus, iterations int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < iterations; i++ {
		bus.RunOnce(ctx)
	}
}

func TestEventHandling(t *testing.T) {
	bus := NewBus()

	// event counters
	var count1 int
	var count2 int
	var count3 int
	var count4 int
	reset := func() {
		count1 = 0
		count2 = 0
		count3 = 0
		count4 = 0
	}

	// check handler functionality
	handler[int](&count1).Handle(nil, 0)
	assert.Equal(t, 1, count1)
	reset()

	// publish without subscribers
	Publish(bus, EventA{})
	Publish(bus, EventB{})
	runCount(bus, 2)
	assert.Zero(t, count1)
	assert.Zero(t, count2)
	assert.Zero(t, count3)
	assert.Zero(t, count4)

	// add some subscribers
	unsub1 := Subscribe(bus, handler[EventA](&count1))
	unsub2 := Subscribe(bus, handler[EventA](&count2))
	unsub3 := Subscribe(bus, handler[EventB](&count3))
	unsub4 := Subscribe(bus, handler[EventB](&count4))

	// publish some events
	reset()
	Publish(bus, EventA{})
	Publish(bus, EventB{})
	Publish(bus, EventC{})
	Publish(bus, EventB{})
	Publish(bus, EventC{})
	runCount(bus, 5)
	assert.Equal(t, 1, count1) // A
	assert.Equal(t, 1, count2) // A
	assert.Equal(t, 2, count3) // B
	assert.Equal(t, 2, count4) // B

	// remove some handlers
	reset()
	unsub1()
	unsub3()
	Publish(bus, EventA{})
	Publish(bus, EventB{})
	Publish(bus, EventC{})
	runCount(bus, 3)
	assert.Equal(t, 0, count1)
	assert.Equal(t, 1, count2) // A
	assert.Equal(t, 0, count3)
	assert.Equal(t, 1, count4) // B

	// check subscribe to all
	reset()
	unsub1 = SubscribeAll(bus, handler[any](&count1))
	unsub3 = SubscribeAll(bus, handler[any](&count3))
	Publish(bus, EventA{})
	Publish(bus, EventB{})
	Publish(bus, EventC{})
	Publish(bus, EventB{})
	runCount(bus, 4)
	assert.Equal(t, 4, count1) // all
	assert.Equal(t, 1, count2) // A
	assert.Equal(t, 4, count3) // all
	assert.Equal(t, 2, count4) // B

	// try nil event
	reset()
	bus.Chan() <- nil      // should not affect counters
	Publish(bus, nil)      // should not affect counters
	Publish(bus, EventA{}) // should affect counters
	runCount(bus, 2)       // direct chan send consumes an operation

	assert.Equal(t, 1, count1) // all
	assert.Equal(t, 1, count2) // A
	assert.Equal(t, 1, count3) // all
	assert.Equal(t, 0, count4) // B

	// try nil handler
	unsubNil1 := Subscribe[EventA](bus, nil)
	unsubNil2 := SubscribeAll(bus, nil)
	Publish(bus, EventA{})
	runCount(bus, 1)
	unsubNil1()
	unsubNil2()

	// remove remaining subscriptions
	reset()
	unsub1()
	unsub2()
	unsub3()
	unsub4()
	Publish(bus, EventA{})
	Publish(bus, EventB{})
	Publish(bus, EventC{})
	runCount(bus, 3)
	assert.Zero(t, count1)
	assert.Zero(t, count2)
	assert.Zero(t, count3)
	assert.Zero(t, count4)
}

func TestBusOptions(t *testing.T) {
	{
		bus := NewBus(WithQueueLength(5))
		assert.Equal(t, 5, cap(bus.events))
		// events should fail to publish when at capacity
		assert.True(t, PublishTry(bus, EventA{}))
		assert.True(t, PublishTry(bus, EventA{}))
		assert.True(t, PublishTry(bus, EventA{}))
		assert.True(t, PublishTry(bus, EventA{}))
		assert.True(t, PublishTry(bus, EventA{}))
		assert.False(t, PublishTry(bus, EventA{}))
	}
	{
		ch := make(chan any, 10)
		bus := NewBus(WithChan(ch))
		assert.Equal(t, ch, bus.events)
	}
}

func TestRunCancel(t *testing.T) {
	ch := make(chan any, 10)
	bus := NewBus(WithChan(ch))
	Publish(bus, EventA{})

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	functionShouldEnd(t, func() {
		bus.Run(canceled)
	})
	functionShouldEnd(t, func() {
		bus.RunOnce(canceled)
	})

	close(ch)
	functionShouldEnd(t, func() {
		bus.Run(context.Background())
	})
	functionShouldEnd(t, func() {
		bus.RunOnce(context.Background())
	})
}

func TestPublishWhenFull(t *testing.T) {
	bus := NewBus(WithQueueLength(0))

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	functionShouldEnd(t, func() {
		t.Helper()
		assert.False(t, PublishTry(bus, EventA{}))
	})

	functionShouldEnd(t, func() {
		t.Helper()
		assert.False(t, PublishContext(canceled, bus, EventA{}))
	})
}

func functionShouldEnd(t *testing.T, function func()) {
	t.Helper()
	var done atomic.Bool
	go func() {
		function()
		done.Store(true)
	}()
	assert.Eventually(t, func() bool { return done.Load() }, time.Second, 10*time.Millisecond)
}
