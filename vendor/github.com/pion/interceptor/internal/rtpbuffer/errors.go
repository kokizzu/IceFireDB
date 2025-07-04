// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtpbuffer

import "errors"

// ErrInvalidSize is returned by newReceiveLog/newRTPBuffer, when an incorrect buffer size is supplied.
var ErrInvalidSize = errors.New("invalid buffer size")

var (
	errPacketReleased          = errors.New("could not retain packet, already released")
	errFailedToCastHeaderPool  = errors.New("could not access header pool, failed cast")
	errFailedToCastPayloadPool = errors.New("could not access payload pool, failed cast")
	errPaddingOverflow         = errors.New("padding size exceeds payload size")
)
