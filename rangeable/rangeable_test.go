package rangeable

import "testing"

func TestCheckRangeIntersect(t *testing.T) {
	tests := []struct {
		name   string
		ranges []RangePair[int]
		want   bool
	}{
		{
			name:   "no overlap",
			ranges: []RangePair[int]{{Start: 1, End: 5}, {Start: 6, End: 10}},
			want:   false,
		},
		{
			name:   "overlap",
			ranges: []RangePair[int]{{Start: 1, End: 5}, {Start: 4, End: 10}},
			want:   true,
		},
		{
			name:   "touching",
			ranges: []RangePair[int]{{Start: 1, End: 5}, {Start: 5, End: 10}},
			want:   false,
		},
		{
			name:   "single range",
			ranges: []RangePair[int]{{Start: 1, End: 5}},
			want:   false,
		},
		{
			name:   "empty",
			ranges: []RangePair[int]{},
			want:   false,
		},
		{
			name:   "multiple overlap",
			ranges: []RangePair[int]{{Start: 1, End: 10}, {Start: 2, End: 5}, {Start: 20, End: 30}},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckRangeIntersect(tt.ranges)
			if got != tt.want {
				t.Fatalf("CheckRangeIntersect: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckRangeIntersectFloat(t *testing.T) {
	ranges := []RangePair[float64]{
		{Start: 1.0, End: 5.0},
		{Start: 4.5, End: 10.0},
	}

	if !CheckRangeIntersect(ranges) {
		t.Fatal("expected overlap for float ranges")
	}
}
