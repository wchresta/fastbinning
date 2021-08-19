# fastbinning
Non-uniform quantization with linear average-case computation time

## Introduction

This is an implementation of the quantization or binning algorithm proposed in [a paper by Cadenas and M. Megson](https://arxiv.org/abs/2108.08228) in Go.

## Example

We use the example presented in the paper:

```go
bin, err := fastbinning.New([]float64{2, 11, 19, 20, 21, 27, 29, 30})
if err != nil {
    // This only happens if the Boundaries are not monotonically increasing.
    fmt.Printf("Creation of Bin failed: %s", err.Error())
}

// Search operations run in constant time and space on average
// Worst case is log(len(Boundaries)) time and constant space
bin.Search(4) // returns 1: 4 is between 2 and 11 which is the first bin
bin.Search(0) // returns 0: 0 is left of the first bin, so it is in bin 0
bin.Search(20.5) // returns 4: 20.5 is between 20 and 21
bin.Search(11) // returns 2: intervals are left-including, so 11 is in [11,19)
```
