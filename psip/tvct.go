package psip

import (
	"github.com/mh-orange/ts"
)

type tvct struct {
	table
}

func newTVCT(payload []byte) VCT {
	return &tvct{table(payload)}
}

func (vct *tvct) NumChannelsInSection() int {
	return int(ts.Uimsbf8(vct.Data()[0]))
}

func (vct *tvct) Channels() []Channel {
	numChannels := vct.NumChannelsInSection()
	channels := make([]Channel, numChannels)
	offset := 1

	for i := 0; i < numChannels; i++ {
		channels[i] = channel(vct.Data()[offset:])
		offset += channels[i].Length()
	}
	return channels
}

func (vct *tvct) channelLength() int {
	length := 0
	for _, channel := range vct.Channels() {
		length += channel.Length()
	}
	return length
}

func (vct *tvct) AdditionalDescriptorsLength() uint16 {
	start := 1 + vct.channelLength()
	end := start + 2
	return ts.Uimsbf16(vct.Data()[start:end], 10)
}
