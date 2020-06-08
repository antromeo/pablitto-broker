package main


import (
	. "broker-pablitto/control_packets"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

const(
	base10=10
	bits32=32
)

func handler(connection net.Conn, s *Subscriptions) {
	defer connection.Close()
	var (
		r   = bufio.NewReader(connection)
	)
	for {

		control_packet, err:=Analyze_packet(r)

		switch err{
		case io.EOF:

		case nil:

			switch control_packet.(type){
			case *ConnectPacket:
				log.Println("CONNECT: Received")

				connect_packet:=control_packet.(*ConnectPacket)
				log.Println(connect_packet)
				session_present:=false
				if connect_packet.Clean_session==false {
					session_present=true
				} else {session_present=false}
				return_code:=byte(Connection_accepted)
				fh:=FixedHeader{Message_type: Connack,  Dup: false, Retain: false}
				connack_packet:=ConnackPacket{fh, session_present, return_code}
				connack_packet.Boxing(connection)

			case *PublishPacket:
				log.Println("PUBLISH: Received")
				publish_packet:=control_packet.(*PublishPacket)
				log.Println(publish_packet)

				if s.Find_topic(publish_packet.Topic_name){
					s.SendToSubscriber(publish_packet.Topic_name, publish_packet)
				} else{

					s.Add_topic(publish_packet.Topic_name)
					s.SendToSubscriber(publish_packet.Topic_name, publish_packet)
				}
				if publish_packet.Qos==0x01{
					fh:=FixedHeader{Message_type: Puback,  Dup: false, Retain: false}
					puback_packet:=PubackPacket{fh,  publish_packet.Message_id}
					puback_packet.Boxing(connection)
				}
			case *SubscribePacket:
				log.Println("SUBSCRIBE: Received")
				subscribe_packet:=control_packet.(*SubscribePacket)
				var return_code []byte
				for _, topic_name:=range subscribe_packet.Topics{
					if s.IsValidTopic(topic_name)==true{
						s.Add_subscription(topic_name,connection)
						return_code=append(return_code, 0x00)
					}
				}

				fh:=FixedHeader{Message_type: Suback,  Dup: false, Retain: false}
				suback_packet:=SubackPacket{fh, subscribe_packet.Message_id, return_code}
				suback_packet.Boxing(connection)

			case *UnsubscribePacket:
				log.Println("UNSUBSCRIBE: Received")
				unsubscribe_packet:=control_packet.(*UnsubscribePacket)

				for _, topic_name:=range unsubscribe_packet.Topics{
					if s.IsValidTopic(topic_name)==true{
						s.Remove_subscription(topic_name,connection)
					}
				}

			case *DisconnectPacket:
				log.Println("DISCONNECT: Received")


			case *PingreqPacket:
				fh:=FixedHeader{Message_type: Pingresp,  Dup: false, Retain: false}
				fmt.Println(control_packet.(*PingreqPacket))
				pingresp_packet:=PingrespPacket{fh}
				pingresp_packet.Boxing(connection)


			default:
				log.Println("Sorry, packet is not yet supported")

			}
		default:
			log.Print("Receive data failed:%s", err)

		}


	}

}

func SocketServer(port int) {
	s:=&Subscriptions{nil}
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	defer listen.Close()
	log.Printf("Begin listen port: %d", port)
	for {
		connection, err := listen.Accept()
		strRemoteAddr := connection.RemoteAddr().String()
		log.Println("Address client: ",strRemoteAddr)
		if err != nil {
			log.Println("Error accept")
		}
		go handler(connection, s)
	}
}

func main() {
	if len(os.Args)>1 {
		port64, err:=strconv.ParseInt(os.Args[1], base10, bits32);
		if err==nil {
			port:=int(port64)
			SocketServer(port)
		}
	} else{

		log.Println("try again with ./boker port_number")

	}

}



