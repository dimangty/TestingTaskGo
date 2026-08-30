package bins

// BinList contains a list of bins.
type BinList struct {
	bins []Bin
}

// NewBinList creates a bin list from the supplied bins.
func NewBinList(bins []Bin) BinList {
	return BinList{
		bins: bins,
	}
}
