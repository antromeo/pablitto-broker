package control_packets

const (
	Connect     = 1
	Connack     = 2
	Publish     = 3
	Puback      = 4
	Pubrec      = 5
	Pubrel      = 6
	Pubcomp     = 7
	Subscribe   = 8
	Suback      = 9
	Unsubscribe = 10
	Unsuback    = 11
	Pingreq     = 12
	Pingresp    = 13
	Disconnect  = 14
)

const (
	Connection_accepted           = 0x00
	Unacceptable_protocol_version = 0x01
	Identifier_rejected           = 0x02
	Server_Unavailable            = 0x03
	Bad_username_password         = 0x04
	Not_authorized                = 0x05
	Nothing_response              = 0xFF
)

const (
	Keep_alive_duration = 60
)
