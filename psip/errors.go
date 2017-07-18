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
