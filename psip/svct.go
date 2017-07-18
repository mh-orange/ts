package psip

type svct []byte

func newSVCT(payload []byte) VCT {
	return svct(payload)
}

func (vct svct) SectionSyntaxIndicator() bool {
	return false
}

func (vct svct) PrivateIndicator() bool {
	return false
}

func (vct svct) SectionLength() int {
	return 0
}

func (vct svct) TransportStreamID() uint16 {
	return 0
}

func (vct svct) VersionNumber() uint8 {
	return 0
}

func (vct svct) CurrentNextIndicator() bool {
	return false
}

func (vct svct) SectionNumber() uint8 {
	return 0
}

func (vct svct) LastSectionNumber() uint8 {
	return 0
}

func (vct svct) ProtocolVersion() uint8 {
	return 0
}

func (vct svct) NumChannelsInSection() int {
	return 0
}

func (vct svct) Channels() []Channel {
	return nil
}

func (vct svct) channelLength() int {
	return 0
}

func (vct svct) AdditionalDescriptorsLength() uint16 {
	return 0
}

func (vct svct) Crc() []byte {
	return nil
}
