/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import "fmt"

func isPowerOfTwo(n uint) bool {
	return n > 0 && (n&(n-1)) == 0
}

func ValidateDimensions(x uint, y uint) error {
	if x <= 0 || y <= 0 {
		return fmt.Errorf("dimensions must be positive")
	}

	if x%4 != 0 || y%4 != 0 {
		return fmt.Errorf("dimensions must be multiples of 4, got %dx%d", x, y)
	}

	if !isPowerOfTwo(x) || !isPowerOfTwo(y) {
		return fmt.Errorf("dimensions must be power of 2, got %dx%d", x, y)
	}
	return nil
}
