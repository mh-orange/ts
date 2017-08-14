package psip

// VCT represents operations on a Terrestrial Virtual Channel Table
type VCT interface {
	NumChannelsInSection() int
	Channels() []Channel
	AdditionalDescriptorsLength() uint16
}
