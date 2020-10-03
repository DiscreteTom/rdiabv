package rdiabv

import "math/rand"

const (
	// DefaultBlockFieldSize is the default size of the data & tag fields of a block.
	DefaultBlockFieldSize = 1024
)

// DefaultBlock consists of a data field and a tag field.
type DefaultBlock struct {
	Data []byte
	Tag  []byte
}

// NewDefaultBlock will initialize the data & tag fields of DefaultBlock
// with an empty byte list of the length DefaultBlockFieldSize.
func NewDefaultBlock() (block *DefaultBlock) {
	block = &DefaultBlock{}
	block.Data = make([]byte, DefaultBlockFieldSize)
	block.Tag = make([]byte, DefaultBlockFieldSize)
	return
}

// Merge will add fields of x to the fields of y.
func (block *DefaultBlock) Merge(x Block, y Block) Block {
	// type assertion
	blockX, _ := x.(*DefaultBlock)
	blockY, _ := y.(*DefaultBlock)
	// Merge each byte
	for i := 0; i < DefaultBlockFieldSize; i++ {
		block.Data[i] = blockX.Data[i] + blockY.Data[i]
		block.Tag[i] = blockX.Tag[i] + blockY.Tag[i]
	}
	return block
}

// Copy will return a copy of the current block
func (block *DefaultBlock) Copy() Block {
	copied := &DefaultBlock{}
	copied.Data = make([]byte, DefaultBlockFieldSize)
	copied.Tag = make([]byte, DefaultBlockFieldSize)
	copy(copied.Data, block.Data)
	copy(copied.Tag, block.Tag)
	return copied
}

// Validate will check whether each byte of the data field is identical to the tag field.
func (block *DefaultBlock) Validate() (ret bool) {
	for i := 0; i < DefaultBlockFieldSize; i++ {
		if block.Data[i] != block.Tag[i] {
			return false
		}
	}
	return true
}

// DefaultBlockGenerator will generate a DefaultBlock whose data & tag fields are identical with the length DefaultBlockFieldSize
// and the content is random.
func DefaultBlockGenerator() (block *DefaultBlock) {
	block = &DefaultBlock{}
	block.Data = make([]byte, DefaultBlockFieldSize)
	rand.Read(block.Data)
	block.Tag = block.Data // copy value
	return
}
