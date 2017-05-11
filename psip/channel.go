package psip

type channel []byte

const (
	c_minimum_length           = 32
	c_major_offset             = 14
	c_descriptor_length_offset = 30
)

func (c channel) ShortName() string {
	return utf16ToString(c, 0, 14)
}

func (c channel) MajorNumber() uint16 {
	return getUimsbf16(c, 14, 16) >> 2 & 0x3ff
}

func (c channel) MinorNumber() uint16 {
	return getUimsbf16(c, 15, 10)
}

func (c channel) ModulationMode() uint8 {
	return 0
}

func (c channel) DescriptorsLength() int {
	return int(c[c_descriptor_length_offset]&0x03)<<8 | int(c[c_descriptor_length_offset+1])
}

func (c channel) Length() int {
	return c_descriptor_length_offset + c.DescriptorsLength() + 2
}
