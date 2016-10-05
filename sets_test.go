package geoindex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	set := newSet()
	assert.Equal(t, set.Size(), 0)

	set.Add(charring.ID(), charring)
	assert.Equal(t, set.Size(), 1)

	set.Add(embankment.ID(), embankment)
	assert.Equal(t, set.Size(), 2)

	set.Remove(charring.ID())
	assert.Equal(t, set.Size(), 1)

	value, ok := set.Get(charring.ID())
	assert.False(t, ok)
	assert.Nil(t, value)

	value, ok = set.Get(embankment.ID())
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, value.(Point).ID(), "Embankment")

	set.Add(picadilly.ID(), picadilly)
	set.Add(oxford.ID(), oxford)

	assert.True(t, pointsEqualIgnoreOrder(toPoints(set.Values()), []Point{picadilly, embankment, oxford}))
}

func toPoints(values []interface{}) []Point {
	result := make([]Point, 0)
	for _, value := range values {
		result = append(result, value.(Point))
	}
	return result
}

func TestExpiringSet(t *testing.T) {
	set := newExpiringSet(Minutes(10))

	currentTime := time.Now()

	now = currentTime
	set.Add(picadilly.ID(), picadilly)

	now = currentTime.Add(5 * time.Minute)
	set.Add(oxford.ID(), oxford)
	assert.Equal(t, set.Size(), 2)
	assert.Equal(t, len(set.Values()), 2)

	set.Remove(picadilly.ID())
	assert.Equal(t, set.Size(), 1)

	now = currentTime.Add(11 * time.Minute)
	assert.Equal(t, set.Size(), 1)

	set.Add(oxford.ID(), oxford)
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.Values()), 1)

	now = currentTime.Add(16 * time.Minute)
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.Values()), 1)

	now = currentTime.Add(22 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	assert.Equal(t, len(set.Values()), 0)

	now = currentTime.Add(24 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	set.Add(oxford.ID(), oxford)
	now = currentTime.Add(25 * time.Minute)
	set.Add(oxford.ID(), oxford)
	now = currentTime.Add(26 * time.Minute)
	set.Add(oxford.ID(), oxford)
	assert.Equal(t, set.Size(), 1)

	set.Remove(oxford.ID())
	assert.Equal(t, set.Size(), 0)
}
