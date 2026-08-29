// Package bins contains types and operations for bins and bin lists.
package bins

import "time"

// Bin describes a bin.
type Bin struct {
	id        string
	private   bool
	createdAt time.Time
	name      string
}

// NewBin creates a bin with the current time as its creation time.
func NewBin(id string, private bool, name string) Bin {
	return Bin{
		id:        id,
		private:   private,
		createdAt: time.Now(),
		name:      name,
	}
}
