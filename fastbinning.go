/*
Copyright 2021 Wanja Chresta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fastbinning

import (
	"fmt"
	"sort"
)

// Binnning for non-uniform bins in asymtotically linear time

// Based on the paper 'Non-uniform quantization with linear average-case computation time'
// by Oswaldo Cadenas and Graham M. Megson: https://arxiv.org/abs/2108.08228

type Bin struct {
	boundaries          []float64 // must be monotonically increasing
	uniformBinWidth     float64
	histogram           []int
	cumulativeHistogram []int
}

// Create a new Bin and run the precalculation step
// boundaries must be monotonically increasing, otherwise
// we return an error
//
// The preparation step runs in linear time and space on the number
// of boundaries.
func New(boundaries []float64) (*Bin, error) {
	// Ensure boundaries are monotonically increasing
	for i, b := range boundaries[1:] {
		if boundaries[i] >= b {
			return nil, fmt.Errorf("boundaries must be monotonically sorted. Found %f >= %f at index %d and %d", boundaries[i], b, i-1, i)
		}
	}

	bin := &Bin{
		boundaries: boundaries,
	}

	bin.precalculation()

	return bin, nil
}

func (bin *Bin) Boundary(i int) float64 {
	return bin.boundaries[i]
}

func (bin *Bin) precalculation() {
	// Number of bins; 1 bin would have 2 boundaries, 2 bins have 3 boundaries, etc.
	m := len(bin.boundaries) - 1

	// Step 1 - set up uniform bins
	totalWidth := bin.boundaries[m] - bin.boundaries[0]

	// We create uniform bins within the range in question. This will help us to
	// find the actual bin an element belongs to withuot having to to a binary
	// search for every element. In the precalculation step we build a histogram
	// of the boundaries within those uniform bins.
	bin.uniformBinWidth = totalWidth / float64(m)

	// Step 2 - histogram of non-uniform bins in uniform bins
	bin.histogram = make([]int, m)

	// Unform bins are numbered as follows:
	// 0   -> (-inf, b[0])
	// 1   -> [b[0], b[1])
	// ...
	// m+1 -> [b[m], inf)
	uniformBinNumber := 1

	// We use the fact that boundaries are sorted.
	// The lowest bound for both the uniform bins and the non-uniform bins is the same
	lowestBound := bin.boundaries[0]

	// We exclude the extreme boundaries b[0] and b[m] as required by the algorithm
	for _, b := range bin.boundaries[1:m] {
		for b > lowestBound+bin.uniformBinWidth {
			// b is outside of the current uniform bin. Find the next unform bin
			lowestBound += bin.uniformBinWidth
			uniformBinNumber += 1
		}

		// The current boundary is in uniform bin uniformBinNumber; count it towards
		// the histogram
		bin.histogram[uniformBinNumber-1] += 1
	}

	// Step 3 - cumulative histogram
	bin.cumulativeHistogram = make([]int, m+1) // We cumulate on uniform boundaries not bins, thus there are m+1
	bin.cumulativeHistogram[0] = 1             // We start at 1 since we excluded the extreme boundaries in step 2
	for i, h := range bin.histogram {
		bin.cumulativeHistogram[i+1] = bin.cumulativeHistogram[i] + h
	}
}

// Search returns the bin-number of a value in a prepared Bin
// Bin needs to be created with New since it performs some precalculation.
// Search used on a non-prepared bin results in a panic
//
// The returned value represents the bin number. 0 means the value
// lies before the first bin (left of the first boundary), while
// the return value of len(bin.Boundray) means it lies above the
// right-most boundary.
// A return of n means the value lies within the interval [bin.Boundary[n-1], bin.Boundary[n])
// meaning 1 represents the left-most proper interval and len(bin.Boundary)-1 represents the
// right most proper interval.
//
// A Search runs in O(1) time on average, as proved by O. Cadenas and G. M. Megson
// and O(1) space.
func (bin *Bin) Search(value float64) int {
	if bin.uniformBinWidth <= 0 {
		panic("Bin needs to be created with New")
	}

	if value < bin.boundaries[0] {
		return 0
	} else if value >= bin.boundaries[len(bin.boundaries)-1] {
		return len(bin.boundaries)
	}

	// We now know bin.boundaries[0] <= value < bin.boundaries[m]
	uniformBinNumber := int((value-bin.boundaries[0])/bin.uniformBinWidth) + 1

	h := bin.histogram[uniformBinNumber-1]

	// if r is used as an index we need to -1 since we're 0-indexing
	r := bin.cumulativeHistogram[uniformBinNumber-1]

	switch h {
	case 0: // case h = 0
		return r
	case 1: // case h = 1
		// We are 0-indexed while the paper is 1 indexed
		if value >= bin.boundaries[r] {
			return r + 1
		} else {
			return r
		}
	case 2: // case h = 2
		if value >= bin.boundaries[r+1] {
			return r + 2
		} else if value < bin.boundaries[r] {
			return r
		} else {
			return r + 1
		}
	default:
		// We cannot use SearchFloat64s because it uses <= instead of <, as we need
		return r + sort.Search(h, func(i int) bool { return value < bin.boundaries[r+i] })
	}
}
