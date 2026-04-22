// Package qr renders QR codes as ASCII half-blocks to an io.Writer.
package qr

import (
	"io"

	"github.com/mdp/qrterminal/v3"
)

// Render writes a QR code encoding `data` to `w` using compact half-block rendering.
func Render(w io.Writer, data string) {
	qrterminal.GenerateHalfBlock(data, qrterminal.L, w)
}
