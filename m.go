package main

import "errors"

func main() {}

// Char. number range  |        UTF-8 octet sequence
// (hexadecimal)       |              (binary)
// --------------------+---------------------------------------------
// 0000 0000-0000 007F | 0xxxxxxx
// 0000 0080-0000 07FF | 110xxxxx 10xxxxxx
// 0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
// 0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx

func decodeRune(b []byte) (r rune, s int, err error) {

	if len(b) == 0 {
		return 0, 0, errors.New("empty input")
	}

	byte0 := b[0]

	switch {
	case byte0 < 0x80: // ASCII | 1 bite character

		if len(b) > 1 {
			if b[1]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}

		r = rune(byte0)
		s = 1

	case byte0&0xE0 == 0xC0: // 2 bytes character

		// 12345678 12345678 12345678 12345678 |  check
		// XXXXXXXX XXXXXXXX XXXXXXXX XXXXXXXX
		// 00000000 00000000 00000000 000XXXXX
		// 00000000 00000000 00000000 00XXXXX0
		// 00000000 00000000 00000000 0XXXXX00
		// 00000000 00000000 00000000 XXXXX000
		// 00000000 00000000 0000000X XXXX0000
		// 00000000 00000000 000000XX XXX00000
		// 00000000 00000000 00000XXX XX000000
		// 00000000 00000000 00000XXX 00XXXXXX
		// 							  10XXXXXX
		//						      00111111 |  bitwise with  " and "

		if len(b) > 2 {
			if b[2]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}

		if len(b) < 2 {
			return 0, 0, errors.New("invalid length")
		}

		byte1 := b[1]

		if byte1&0xC0 != 0x80 {
			return 0, 0, errors.New("invalid continuation byte")
		}

		r = (((rune(byte0) & 0x1F) << 6) | (rune(b[1] & 0x3F)))
		s = 2

		if r < 0x80 {
			return 0, 0, errors.New("overlong")
		}

	case byte0&0xF0 == 0xE0: // 3 bytes character

		if len(b) > 3 {
			if b[3]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}

		if len(b) < 3 {
			return 0, 0, errors.New("invalid length")
		}

		byte1 := b[1]
		byte2 := b[2]

		if (byte0 == 0xE0) && (byte1 < 0xA0) {
			return 0, 0, errors.New("overlog")
		}

		if (byte1&0xC0 != 0x80) || (byte2&0xC0 != 0x80) {
			return 0, 0, errors.New("invalid continuation byte")
		}

		r = (((rune(byte0) & 0x0F) << 12) | ((rune(b[1]) & 0x3F) << 6) | (rune(b[2]) & 0x3F))
		s = 3

	case byte0&0xF8 == 0xF0: // 4 bytes character

		if len(b) > 4 {
			if b[4]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}

		if len(b) < 3 {
			return 0, 0, errors.New("invalid length")
		}

		byte1 := b[1]
		byte2 := b[2]
		byte3 := b[3]

		if byte0 == 0xF0 && byte1 < 0x90 {
			return 0, 0, errors.New("overlong")
		}

		if (byte1&0xC0 != 0x80) || (byte2&0xC0 != 0x80) || (byte3&0xC0 != 0x80) {
			return 0, 0, errors.New("invalid continuation byte")
		}

		r = (((rune(byte0) & 0x07) << 18) | ((rune(b[1]) & 0x3F) << 12) | ((rune(b[2]) & 0x3F) << 6) | (rune(b[3]) & 0x3F))
		s = 4

	default:
		return 0, 0, errors.New("invalid utf8")
	}

	if (r >= 0xD800) && (r <= 0xDFFF) {
		return 0, 0, errors.New("surrogate hal")
	}

	if r > 0x10FFFF {
		return 0, 0, errors.New("too big")
	}

	return r, s, nil
}
