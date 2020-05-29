package control_packets

import (
	"bytes"
	"io"
)

type ConnectPacket struct {
	FixedHeader
	Protocol_name    string
	Protocol_version byte
	Clean_session    bool
	Will_flag        bool
	Will_qos         byte
	Will_retain      bool
	Username_flag    bool
	Password_flag    bool
	Reserved_bit     byte
	Keep_alive       uint16
	//optional
	Client_id        string
	Will_topic       string
	Will_message     []byte
	Username         string
	Password         []byte
}


func (connect_packet *ConnectPacket) Boxing(writer io.Writer) error {
	var body_bytes bytes.Buffer
	var err error

	body_bytes.Write(encode_string(connect_packet.Protocol_name))
	body_bytes.WriteByte(connect_packet.Protocol_version)
	body_bytes.WriteByte(bool_to_byte(connect_packet.Clean_session)<<1 | bool_to_byte(connect_packet.Will_flag)<<2 |
		connect_packet.Will_qos<<3 | bool_to_byte(connect_packet.Will_retain)<<5 | bool_to_byte(connect_packet.Password_flag)<<6 |
		bool_to_byte(connect_packet.Username_flag)<<7)
	body_bytes.Write(encode_uint16(connect_packet.Keep_alive))
	body_bytes.Write(encode_string(connect_packet.Client_id))
	if connect_packet.Will_flag {
		body_bytes.Write(encode_string(connect_packet.Will_topic))
		body_bytes.Write(encode_bytes(connect_packet.Will_message))
	}
	if connect_packet.Username_flag {
		body_bytes.Write(encode_string(connect_packet.Username))
	}
	if connect_packet.Password_flag {
		body_bytes.Write(encode_bytes(connect_packet.Password))
	}
	connect_packet.FixedHeader.Remaining_len = body_bytes.Len()
	packet := connect_packet.FixedHeader.Boxing_fh()
	packet.Write(body_bytes.Bytes())
	_, err = packet.WriteTo(writer)

	return err
}


func (connect_packet *ConnectPacket) Unboxing(bytes_reader io.Reader) error {
	var err error
	connect_packet.Protocol_name, err = decode_string(bytes_reader)
	if err != nil {
		return err
	}
	connect_packet.Protocol_version, err = decode_byte(bytes_reader)
	if err != nil {
		return err
	}
	options, err := decode_byte(bytes_reader)
	if err != nil {
		return err
	}
	connect_packet.Reserved_bit = 1 & options
	connect_packet.Clean_session = 1&(options>>1) > 0
	connect_packet.Will_flag = 1&(options>>2) > 0
	connect_packet.Will_qos = 3 & (options >> 3)
	connect_packet.Will_retain = 1&(options>>5) > 0
	connect_packet.Password_flag = 1&(options>>6) > 0
	connect_packet.Username_flag = 1&(options>>7) > 0
	connect_packet.Keep_alive, err = decode_uint16(bytes_reader)
	if err != nil {
		return err
	}
	connect_packet.Client_id, err = decode_string(bytes_reader)
	if err != nil {
		return err
	}
	if connect_packet.Will_flag {
		connect_packet.Will_topic, err = decode_string(bytes_reader)
		if err != nil {
			return err
		}
		connect_packet.Will_message, err = decode_bytes(bytes_reader)
		if err != nil {
			return err
		}
	}
	if connect_packet.Username_flag {
		connect_packet.Username, err = decode_string(bytes_reader)
		if err != nil {
			return err
		}
	}
	if connect_packet.Password_flag {
		connect_packet.Password, err = decode_bytes(bytes_reader)
		if err != nil {
			return err
		}
	}
	return nil
}


