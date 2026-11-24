package send_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	entity "energosbyt/internal/entity/fact/gauge/meter/reading"
	saramamock "energosbyt/internal/mocks/sarama"
	send "energosbyt/internal/repo/fact/gauge/meter/reading/send"
)

func TestFacts_Run(t *testing.T) {
	t.Parallel()

	t.Run("sends single fact successfully", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2020, 3, 15, 10, 30, 0, 0, time.UTC)
		fact := entity.NewFact(now, 123, 456, 789, 10, 150, 80, time.March)
		mockProducer := new(saramamock.SyncProducer)
		mockProducer.
			On("SendMessage", mock.Anything).
			Return(int32(0), int64(0), nil).
			Once()
		facts := send.NewFacts(mockProducer, "meter_readings", fact)
		err := facts.Run(context.Background())
		assert.NoError(t, err)
		mockProducer.AssertExpectations(t)
	})

	t.Run("sends multiple facts successfully", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2020, 3, 15, 10, 30, 0, 0, time.UTC)
		fact1 := entity.NewFact(now, 123, 456, 789, 10, 150, 80, time.March)
		fact2 := entity.NewFact(now, 124, 457, 790, 11, 160, 90, time.April)
		mockProducer := new(saramamock.SyncProducer)
		mockProducer.
			On("SendMessage", mock.Anything).
			Return(int32(0), int64(0), nil).
			Twice()
		facts := send.NewFacts(mockProducer, "meter_readings", fact1, fact2)
		err := facts.Run(context.Background())
		assert.NoError(t, err)
		mockProducer.AssertExpectations(t)
	})

	t.Run("returns error on send failure", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2020, 3, 15, 10, 30, 0, 0, time.UTC)
		fact := entity.NewFact(now, 123, 456, 789, 10, 150, 80, time.March)
		mockProducer := new(saramamock.SyncProducer)
		mockProducer.
			On("SendMessage", mock.Anything).
			Return(int32(0), int64(0), assert.AnError).
			Once()
		facts := send.NewFacts(mockProducer, "meter_readings", fact)
		err := facts.Run(context.Background())
		assert.Error(t, err)
		mockProducer.AssertExpectations(t)
	})

	t.Run("sends no facts when empty", func(t *testing.T) {
		t.Parallel()
		mockProducer := new(saramamock.SyncProducer)
		facts := send.NewFacts(mockProducer, "meter_readings")
		err := facts.Run(context.Background())
		assert.NoError(t, err)
		mockProducer.AssertExpectations(t)
	})
}
