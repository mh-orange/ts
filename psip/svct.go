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
