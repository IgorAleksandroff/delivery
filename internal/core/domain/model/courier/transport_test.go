package courier_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

func TestTransportMoveTowards(t *testing.T) {
	testCases := []struct {
		name           string
		transport      *courier.Transport
		current        kernel.Location
		target         kernel.Location
		expectedResult kernel.Location
	}{
		{
			name:           "Already at target",
			transport:      courier.MustNewTransport("Bicycle", 2),
			current:        kernel.MustNewLocation(5, 5),
			target:         kernel.MustNewLocation(5, 5),
			expectedResult: kernel.MustNewLocation(5, 5),
		},
		{
			name:           "Horizontal movement (right) - full speed",
			transport:      courier.MustNewTransport("Bicycle", 2),
			current:        kernel.MustNewLocation(1, 1),
			target:         kernel.MustNewLocation(10, 1),
			expectedResult: kernel.MustNewLocation(3, 1),
		},
		{
			name:           "Horizontal movement (left) - full speed",
			transport:      courier.MustNewTransport("Bicycle", 2),
			current:        kernel.MustNewLocation(10, 1),
			target:         kernel.MustNewLocation(1, 1),
			expectedResult: kernel.MustNewLocation(8, 1),
		},
		{
			name:           "Vertical movement (up) - full speed",
			transport:      courier.MustNewTransport("Bicycle", 2),
			current:        kernel.MustNewLocation(1, 1),
			target:         kernel.MustNewLocation(1, 10),
			expectedResult: kernel.MustNewLocation(1, 3),
		},
		{
			name:           "Vertical movement (down) - full speed",
			transport:      courier.MustNewTransport("Bicycle", 2),
			current:        kernel.MustNewLocation(1, 10),
			target:         kernel.MustNewLocation(1, 1),
			expectedResult: kernel.MustNewLocation(1, 8),
		},
		{
			name:           "Target closer than speed - horizontal",
			transport:      courier.MustNewTransport("Bicycle", 3),
			current:        kernel.MustNewLocation(1, 1),
			target:         kernel.MustNewLocation(3, 1),
			expectedResult: kernel.MustNewLocation(3, 1),
		},
		{
			name:           "Target closer than speed - vertical",
			transport:      courier.MustNewTransport("Bicycle", 3),
			current:        kernel.MustNewLocation(1, 1),
			target:         kernel.MustNewLocation(1, 2),
			expectedResult: kernel.MustNewLocation(1, 2),
		},
		{
			name:           "Combined movement - prioritize X then Y",
			transport:      courier.MustNewTransport("Bicycle", 3),
			current:        kernel.MustNewLocation(1, 1),
			target:         kernel.MustNewLocation(3, 4),
			expectedResult: kernel.MustNewLocation(3, 2),
		},
		{
			name:           "Combined movement - only Y remains",
			transport:      courier.MustNewTransport("Bicycle", 3),
			current:        kernel.MustNewLocation(5, 1),
			target:         kernel.MustNewLocation(5, 5),
			expectedResult: kernel.MustNewLocation(5, 4),
		},
		{
			name:           "Faster transport - speed 3",
			transport:      courier.MustNewTransport("Car", 3),
			current:        kernel.MustNewLocation(9, 9),
			target:         kernel.MustNewLocation(10, 10),
			expectedResult: kernel.MustNewLocation(10, 10),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := tc.transport.MoveTowards(tc.current, tc.target)
			assert.NoError(t, err)

			// Assert
			if !result.Equals(tc.expectedResult) {
				t.Errorf("MoveTowards from %v to %v = %v, want %v",
					tc.current, tc.target, result, tc.expectedResult)
			}
		})
	}
}
