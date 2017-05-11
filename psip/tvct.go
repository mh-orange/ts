/*
MIT License

Copyright 2016 Comcast Cable Communications Management, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
