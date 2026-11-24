package user_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/user"
)

func TestUser(t *testing.T) {
	t.Parallel()

	t.Run("creates user with all fields", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		name := "John Doe"
		currentPhone := "+79001234567"
		previousPhone := "+79007654321"
		now := time.Now()
		updatedAt := now.Add(time.Hour)
		u := user.NewUser(
			id,
			name,
			currentPhone,
			previousPhone,
			now,
			updatedAt,
		)
		assert.Equal(t, id, u.ID())
	})

	t.Run("creates user without previous phone", func(t *testing.T) {
		t.Parallel()
		id := uint64(789)
		name := "Bob Wilson"
		currentPhone := "+79009998877"
		previousPhone := ""
		now := time.Now()
		updatedAt := now
		u := user.NewUser(
			id,
			name,
			currentPhone,
			previousPhone,
			now,
			updatedAt,
		)
		assert.Equal(t, id, u.ID())
	})
}
