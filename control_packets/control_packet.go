package control_packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ControlPacket interface {
	Boxing(io.Writer) error
	Unboxing(io.Reader) error
}

func encode_bytes(bytes_vector []byte) []byte {
	len_message := make([]byte, 2)
	binary.BigEndian.PutUint16(len_message, uint16(len(bytes_vector)))
	return append(len_message, bytes_vector...)
}

func decode_bytes(bytes_io_reader io.Reader) ([]byte, error) {
	len_message, err := decode_uint16(bytes_io_reader)
	if err != nil {
		return nil, err
	}

	message:= make([]byte, len_message)
	_, err = bytes_io_reader.Read(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func decode_byte(bytes_reader io.Reader) (byte, error) {
	len_message := make([]byte, 1)
	_, err := bytes_reader.Read(len_message)
	if err != nil {
		return 0, err
	}
	return len_message[0], nil
}



func encode_uint16(value uint16) []byte {
	bytes_uint16 := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes_uint16, value)
	return bytes_uint16
}

func decode_uint16(bytes_reader io.Reader) (uint16, error) {
	number_bytes := make([]byte, 2)
	_, err := bytes_reader.Read(number_bytes)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(number_bytes), nil
}

func encode_string(str string) []byte {
	return encode_bytes([]byte(str))
}

func decode_string(bytes io.Reader) (string, error) {
	buf, err := decode_bytes(bytes)
	return string(buf), err
}

func Analyze_packet(reader io.Reader) (ControlPacket, error) {
	var fixed_header FixedHeader
	mybyte := make([]byte, 1)

	_, err := io.ReadFull(reader, mybyte)
	if err != nil {
		return nil, err
	}

	err = fixed_header.Unboxing_fh(mybyte[0], reader)
	if err != nil {
		return nil, err
	}

	control_packet, err := Factory_control_packet(fixed_header)
	if err != nil {
		return nil, err
	}

	variable_header := make([]byte, fixed_header.Remaining_len)
	len_variable_header, err := io.ReadFull(reader, variable_header)
	if err != nil {
		return nil, err
	}
	if len_variable_header != fixed_header.Remaining_len {
		return nil, errors.New("Failed to read data")
	}
	err = control_packet.Unboxing(bytes.NewBuffer(variable_header))
	return control_packet, err
}

func Factory_control_packet(fixed_header FixedHeader) (ControlPacket, error) {
	switch fixed_header.Message_type {
	case Connect:
		return &ConnectPacket{FixedHeader: fixed_header}, nil
	case Connack:
		return &ConnackPacket{FixedHeader: fixed_header}, nil
	case Disconnect:
		return &DisconnectPacket{FixedHeader: fixed_header}, nil
	case Publish:
		return &PublishPacket{FixedHeader: fixed_header}, nil
	case Puback:
		return &PubackPacket{FixedHeader: fixed_header}, nil
	case Pubrec:
		break
	case Pubrel:
		break
	case Pubcomp:
		break
	case Subscribe:
		return &SubscribePacket{FixedHeader: fixed_header}, nil
	case Suback:
		return &SubackPacket{FixedHeader: fixed_header}, nil
	case Unsubscribe:
		return &UnsubscribePacket{FixedHeader: fixed_header}, nil
	case Unsuback:
		return &UnsubackPacket{FixedHeader: fixed_header}, nil
	case Pingreq:
		return &PingreqPacket{FixedHeader: fixed_header}, nil
	case Pingresp:
		return &PingrespPacket{FixedHeader: fixed_header}, nil
	}
	return nil, fmt.Errorf("Unsupported Packet Type (in hex) 0x%x", fixed_header.Message_type)
}
