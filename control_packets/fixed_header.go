package control_packets

import (
	"bytes"
	"io"
)

type FixedHeader struct {
	Message_type    byte
	Dup             bool
	Qos             byte
	Retain          bool
	Remaining_len   int
}


func (fixed_header *FixedHeader) Boxing_fh() bytes.Buffer {
	var header_buffer bytes.Buffer
	header_buffer.WriteByte(fixed_header.Message_type<<4 | bool_to_byte(fixed_header.Dup)<<3 | fixed_header.Qos<<1 | bool_to_byte(fixed_header.Retain))
	header_buffer.Write(encode_len(fixed_header.Remaining_len))
	return header_buffer
}

func (fixed_header *FixedHeader) Unboxing_fh(first_byte byte, reader io.Reader) error {
	var err error
	fixed_header.Message_type = first_byte >> 4
	fixed_header.Dup = (first_byte>>3)&0x01 > 0
	fixed_header.Qos = (first_byte >> 1) & 0x03
	fixed_header.Retain = first_byte&0x01 > 0

	fixed_header.Remaining_len, err = decode_len(reader)
	return err
}

func bool_to_byte(boolean bool) byte {
	if boolean==true { return 1} else {return 0}
}



func encode_len(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}


func decode_len(reader io.Reader) (int, error) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for multiplier < 27 { 
		_, err := io.ReadFull(reader, b)
		if err != nil {
			return 0, err
		}
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength), nil
}