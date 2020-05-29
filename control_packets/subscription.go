package control_packets

import (
"fmt"
"io"
"net"
"strings"
)

type Topic struct{
	topic_name string
	connections []net.Conn //per ora sono con la lettera grande per fare append
}



func (t *Topic)add_connection(connection net.Conn){
	t.connections=append(t.connections, connection)
	fmt.Printf(" %v\n",  t.connections)
	fmt.Println(len(t.connections))
}

func (t *Topic)remove_connection(connection net.Conn){
	indice:=-1
	for i, value:=range t.connections{
		if value==connection{
			indice=i
		}
	}
	if indice!=-1{	t.connections=append(t.connections[:indice], t.connections[indice+1:]...)}
}

type Subscriptions struct{
	Topics []Topic
}

func (s *Subscriptions)Add_topic(topic_name string){
	t:= Topic{topic_name, nil}
	s.Topics=append(s.Topics, t)
}

func (s *Subscriptions)Find_topic(topic_name string) bool{
	for _, value:=range s.Topics{
		if value.topic_name==topic_name{
			return true
		}
	}
	return false
}


func (s *Subscriptions)Add_subscription(topic_name string,connection net.Conn){
	flag:=false
	for i, value:=range s.Topics{
		if value.topic_name==topic_name{
			flag=true
			s.Topics[i].add_connection(connection)
		}
	}
	if flag==false {
		t:= Topic{topic_name, nil}
		t.add_connection(connection)
		s.Topics=append(s.Topics, t)
	}
}




func (s *Subscriptions)Remove_subscription(topic_name string,connection net.Conn){
	for _, value:=range s.Topics{
		if value.topic_name==topic_name{
			for i, conn:=range value.connections{
				if conn==connection{
					s.Topics[i].remove_connection(conn)
				}
			}
		}
	}

}


func (s *Subscriptions)IsValidTopic(topic_name string) bool{
	members:=strings.Split(topic_name, "/")
	idx_last_element:=len(members)-1
	for index, member:=range members{
		if strings.Contains(member, "+") && len(member)>1 {
			fmt.Println("+ deve essere da solo")
			return false
		}
		if member =="#" && index!=idx_last_element {
			fmt.Println("# deve essere alla fine")
			return false
		}
		fmt.Println(member)
	}
	return true
}


func (s *Subscriptions) SendToSubscriber(topic_name string, publish_packet *PublishPacket) error {
	var err error
	//fmt.Println("VEDIAMO LE SUBSCRIPTION= ", s.Topics)
	for _, value:=range s.Topics {
		if s.MatchTopic(topic_name, value.topic_name){
			for _, conn:=range value.connections { err=s.send(conn, publish_packet) }
		}

	}
	return err
}


func (s *Subscriptions) MatchTopic(topic_without_wildcard string, topic_optional_wildcard string) bool {
	wildcard:=false
	if strings.Contains(topic_optional_wildcard, "+") || strings.Contains(topic_optional_wildcard, "#"){
		wildcard=true
	}
	if wildcard==false{
		if topic_without_wildcard == topic_optional_wildcard {return true}
	} else {
		members_1:=strings.Split(topic_without_wildcard, "/")
		members_2:=strings.Split(topic_optional_wildcard, "/")

		compare:=s.compareArrayString(members_1, members_2)
		return compare

	}
	return false
}

func (s *Subscriptions) compareArrayString(members_1 []string, members_2 []string) bool {
	len_members_1:=len(members_1)
	len_members_2:=len(members_2)
	var compare int
	if len_members_1>len_members_2 {compare=1}
	if len_members_1==len_members_2 {compare=0}
	if len_members_1<len_members_2 {compare=-1}

	switch compare {
	case 1:
		//s1>s2
		idx_last_element_2:=len(members_2)-1
		if members_2[idx_last_element_2] =="#"{
			for i:=0; i<idx_last_element_2; i++ {
				if members_1[i]!=members_2[i] {return false}
			}
			return true
		}
		return false
	case 0:
		for index, _:=range members_1{
			if members_2[index]!="+"{
				if members_1[index]!=members_2[index] {return false}
			}
			if members_2[index]=="#"{
				return true
			}
		}
		return  true
	case -1:
		//s2>s1, s1 non può contenere wildcard speciali, non lo può matchare
		return false
	}
	return false
}

func (s *Subscriptions) send(w io.Writer, publish_packet *PublishPacket) error {
	err := publish_packet.Boxing(w)
	return err
}


