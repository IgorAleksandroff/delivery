package kernel_test

import (
	"testing"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

func TestNewLocation(t *testing.T) {
	testCases := []struct {
		name        string
		x           int
		y           int
		expectError bool
	}{
		{
			name:        "Valid coordinates (minimum)",
			x:           1,
			y:           1,
			expectError: false,
		},
		{
			name:        "Valid coordinates (maximum)",
			x:           10,
			y:           10,
			expectError: false,
		},
		{
			name:        "Valid coordinates (middle)",
			x:           5,
			y:           5,
			expectError: false,
		},
		{
			name:        "Invalid X coordinate (too small)",
			x:           0,
			y:           5,
			expectError: true,
		},
		{
			name:        "Invalid Y coordinate (too small)",
			x:           5,
			y:           0,
			expectError: true,
		},
		{
			name:        "Invalid X coordinate (too large)",
			x:           11,
			y:           5,
			expectError: true,
		},
		{
			name:        "Invalid Y coordinate (too large)",
			x:           5,
			y:           11,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			location, err := kernel.NewLocation(tc.x, tc.y)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if location.X() != tc.x {
					t.Errorf("X = %v, want %v", location.X(), tc.x)
				}
				if location.Y() != tc.y {
					t.Errorf("Y = %v, want %v", location.Y(), tc.y)
				}
			}
		})
	}
}

func TestLocationEquals(t *testing.T) {
	testCases := []struct {
		name     string
		loc1     kernel.Location
		loc2     kernel.Location
		expected bool
	}{
		{
			name:     "Same locations",
			loc1:     kernel.MustNewLocation(5, 5),
			loc2:     kernel.MustNewLocation(5, 5),
			expected: true,
		},
		{
			name:     "Different X coordinates",
			loc1:     kernel.MustNewLocation(5, 5),
			loc2:     kernel.MustNewLocation(6, 5),
			expected: false,
		},
		{
			name:     "Different Y coordinates",
			loc1:     kernel.MustNewLocation(5, 5),
			loc2:     kernel.MustNewLocation(5, 6),
			expected: false,
		},
		{
			name:     "Completely different coordinates",
			loc1:     kernel.MustNewLocation(1, 1),
			loc2:     kernel.MustNewLocation(10, 10),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := tc.loc1.Equals(tc.loc2)

			// Assert
			if result != tc.expected {
				t.Errorf("%v.Equals(%v) = %v, want %v", tc.loc1, tc.loc2, result, tc.expected)
			}

			// Check symmetry
			reverseResult := tc.loc2.Equals(tc.loc1)
			if reverseResult != result {
				t.Errorf("Equality is not symmetric: %v.Equals(%v) = %v, but %v.Equals(%v) = %v",
					tc.loc1, tc.loc2, result, tc.loc2, tc.loc1, reverseResult)
			}
		})
	}
}

func TestLocationDistance(t *testing.T) {
	testCases := []struct {
		name     string
		loc1     kernel.Location
		loc2     kernel.Location
		expected int
	}{
		{
			name:     "Same location",
			loc1:     kernel.MustNewLocation(5, 5),
			loc2:     kernel.MustNewLocation(5, 5),
			expected: 0,
		},
		{
			name:     "Horizontal distance",
			loc1:     kernel.MustNewLocation(1, 5),
			loc2:     kernel.MustNewLocation(6, 5),
			expected: 5,
		},
		{
			name:     "Vertical distance",
			loc1:     kernel.MustNewLocation(5, 1),
			loc2:     kernel.MustNewLocation(5, 8),
			expected: 7,
		},
		{
			name:     "Diagonal distance (Manhattan)",
			loc1:     kernel.MustNewLocation(1, 1),
			loc2:     kernel.MustNewLocation(4, 5),
			expected: 7, // |4-1| + |5-1| = 3 + 4 = 7
		},
		{
			name:     "DistanceTo between extreme corners",
			loc1:     kernel.MustNewLocation(1, 1),
			loc2:     kernel.MustNewLocation(10, 10),
			expected: 18, // |10-1| + |10-1| = 9 + 9 = 18
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := tc.loc1.DistanceTo(tc.loc2)

			if distance != tc.expected {
				t.Errorf("%v.DistanceTo(%v) = %v, want %v", tc.loc1, tc.loc2, distance, tc.expected)
			}

			// Test symmetry property of distance
			reverseDistance := tc.loc2.DistanceTo(tc.loc1)
			if reverseDistance != distance {
				t.Errorf("DistanceTo is not symmetric: %v.DistanceTo(%v) = %v, but %v.DistanceTo(%v) = %v",
					tc.loc1, tc.loc2, distance, tc.loc2, tc.loc1, reverseDistance)
			}
		})
	}
}

func TestCreateRandomLocation(t *testing.T) {
	for i := 0; i < 100; i++ {
		loc := kernel.CreateRandomLocation()

		if loc.X() < 1 || loc.X() > 10 {
			t.Errorf("Random location X = %d out of valid range [1,10]", loc.X())
		}

		if loc.Y() < 1 || loc.Y() > 10 {
			t.Errorf("Random location Y = %d out of valid range [1,10]", loc.Y())
		}
	}
}

func TestLocationImmutability(t *testing.T) {

	loc1 := kernel.MustNewLocation(5, 5)
	loc2 := loc1 // Создаем копию

	if !loc1.Equals(loc2) {
		t.Errorf("Copy of location should be equal to original")
	}
}
