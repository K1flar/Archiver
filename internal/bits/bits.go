package cbits

import (
	"fmt"
	"math"
)

type Bits struct {
	bytes []byte
	count int
}

func NewBits(count, cap int) Bits {
	return Bits{
		bytes: make([]byte, int(math.Ceil(float64(cap)/8))),
		count: count,
	}
}

func Copy(dst Bits) Bits {
	new := NewBits(dst.GetCount(), dst.GetCount())
	for i := 0; i < dst.GetCount(); i++ {
		if bit, _ := dst.GetBit(i); bit {
			new.SetBit(i)
		}
	}

	return new
}

func NewLoadBytes(bytes []byte, count int) Bits {
	return Bits{
		bytes: bytes,
		count: count,
	}
}

func (b *Bits) GetCount() int {
	return b.count
}

func (b *Bits) GetBit(numberOfBit int) (bool, error) {
	if err := b.checkOutOfRange(numberOfBit); err != nil {
		return false, err
	}
	numberOfByte, offset := b.getPosition(numberOfBit)

	return b.bytes[numberOfByte]&(1<<(8-offset-1)) != 0, nil
}

func (b *Bits) SetBit(numberOfBit int) error {
	if err := b.checkOutOfRange(numberOfBit); err != nil {
		return err
	}
	numberOfByte, offset := b.getPosition(numberOfBit)

	b.bytes[numberOfByte] |= (1 << (8 - offset - 1))
	return nil
}

func (b *Bits) ClearBit(numberOfBit int) error {
	if err := b.checkOutOfRange(numberOfBit); err != nil {
		return err
	}
	numberOfByte, offset := b.getPosition(numberOfBit)
	b.bytes[numberOfByte] &= ^(1 << (8 - offset - 1))
	return nil
}

func (b *Bits) AppendBit(val byte) {
	if b.count == len(b.bytes)*8 {
		b.bytes = append(b.bytes, 0)
	}
	b.count++
	if val == 1 {
		b.SetBit(b.count - 1)
	}
}

func (b *Bits) ToString() string {
	str := ""
	for i := 0; i < b.count; i++ {
		bit, _ := b.GetBit(i)
		if bit {
			str += "1"
		} else {
			str += "0"
		}
	}
	return str
}

func (b *Bits) GetBytes() []byte {
	return b.bytes
}

func (b *Bits) getPosition(numberOfBit int) (int, int) {
	return numberOfBit / 8, numberOfBit % 8
}

func (b *Bits) checkOutOfRange(numberOfBit int) error {
	if numberOfBit >= b.count {
		return fmt.Errorf("cbits: numberOfBit=%d out of range, with size=%d", numberOfBit, b.count)
	}
	return nil
}
