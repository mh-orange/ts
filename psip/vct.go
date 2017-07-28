package psip

// VCT represents operations on a Terrestrial Virtual Channel Table
type VCT interface {
	Table
	NumChannelsInSection() int
	Channels() []Channel
	AdditionalDescriptorsLength() uint16
}
