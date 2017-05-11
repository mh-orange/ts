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

import "errors"

var (
	// returned when expected T-VCT packet is not found when
	// reading TS packets.
	ErrVCTNotFound = errors.New("No T-VCT was found while reading TS")

	// returned when a Terrerstrial VCT cannot be parsed because there are not enough bytes
	ErrInvalidTVCTLength = errors.New("too few bytes to parse T-VCT")

	// returned when the table ID of the VCT is unknown
	ErrInvalidTableID = errors.New("invalid VCT table ID (unknown VCT table type)")

	// returned when updating PSIP tables and a Continuity Error occurs.  Continuity errors happen
	// when a packet does not contain a Payload Unit Start and the last packet received was not
	// the previous packet as indicated by the MPEG header continuity counter
	ErrContinuity = errors.New("continuity error")

	ErrShortBuffer = errors.New("Buffer too short to parse")
)
