package rdiabv

// Block defines the interface of a block. Usually a block consists of a data field and a tag field.
// Use DefaultBlock for the default block definition.
type Block interface {
	Validate() bool         // Validate check whether a block is valid.
	Merge(x, y Block) Block // Merge sets the current block to the result of merging block x & y.
	Copy() Block            // Copy creates a copy of the current block.
}
