package psip

type tvct []byte

func newTVCT(payload []byte) VCT {
	return tvct(payload)
}

func (vct tvct) SectionSyntaxIndicator() bool {
	return getBool(vct, 0, 0)
}

func (vct tvct) PrivateIndicator() bool {
	return getBool(vct, 0, 1)
}

func (vct tvct) SectionLength() int {
	return int(getUimsbf16(vct, 1, 12))
}

func (vct tvct) TransportStreamID() uint16 {
	return getUimsbf16(vct, 3, 16)
}

func (vct tvct) VersionNumber() uint8 {
	return getUimsbf8(vct, 5)
}

func (vct tvct) CurrentNextIndicator() bool {
	return getBool(vct, 5, 7)
}

func (vct tvct) SectionNumber() uint8 {
	return getUimsbf8(vct, 6)
}

func (vct tvct) LastSectionNumber() uint8 {
	return getUimsbf8(vct, 7)
}

func (vct tvct) ProtocolVersion() uint8 {
	return getUimsbf8(vct, 8)
}

func (vct tvct) NumChannelsInSection() int {
	return int(getUimsbf8(vct, 9))
}

func (vct tvct) Channels() []Channel {
	numChannels := vct.NumChannelsInSection()
	channels := make([]Channel, numChannels)
	offset := 10

	for i := 0; i < numChannels; i++ {
		channels[i] = channel(vct[offset:])
		offset += channels[i].Length()
	}
	return channels
}

func (vct tvct) channelLength() int {
	length := 0
	for _, channel := range vct.Channels() {
		length += channel.Length()
	}
	return length
}

func (vct tvct) AdditionalDescriptorsLength() uint16 {
	return getUimsbf16(vct, 10+vct.channelLength(), 10)
}

func (vct tvct) Crc() []byte {
	if len(vct) < 4 {
		return nil
	}
	return vct[len(vct)-4:]
}
