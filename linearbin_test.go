package linearbin

import "testing"

func cmpIntSlice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func TestBinningExample(t *testing.T) {
	//                    0   1   2   3   4   5   6   7   8
	bin, _ := New([]float64{2, 11, 19, 20, 21, 27, 29, 30})

	if bin.uniformBinWidth != 4 {
		t.Errorf("Expected uniformBinWidth to be 4 but got %f\n", bin.uniformBinWidth)
	}

	expectedHistogram := []int{0, 0, 1, 0, 3, 0, 2}
	if !cmpIntSlice(bin.histogram, expectedHistogram) {
		t.Errorf("Expected histogram\n%v but got\n%v\n", expectedHistogram, bin.histogram)
	}

	expectedCumulativeHistrogram := []int{1, 1, 1, 2, 2, 5, 5, 7}
	if !cmpIntSlice(bin.cumulativeHistogram, expectedCumulativeHistrogram) {
		t.Errorf("Expected cumulativeHistogram\n%v but got\n%v\n", expectedCumulativeHistrogram, bin.cumulativeHistogram)
	}

	testData := map[float64]int{
		-4:   0,
		0:    0,
		2:    1,
		3:    1,
		7:    1,
		10.5: 1,
		13.2: 2,
		18.9: 2,
		19.9: 3,
		20:   4,
		29.9: 7,
		30:   8,
		31:   8,
		99:   8,
	}

	for data, exp := range testData {
		out := bin.Search(data)
		if out != exp {
			t.Errorf("Expected %f to be binned to %d but got %d\n", data, exp, out)
		}
	}
}
