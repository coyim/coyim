package session

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
)

func (s *session) GetRooms(server jid.Domain) <-chan *muc.RoomListing {
	// result := make(chan *muc.RoomListing, 20)

	// rp, _, err := s.conn.SendIQ(server.String(), "get", &data.DiscoveryItemsQuery{})
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return nil
	// }

	// r, ok := <-rp
	// if !ok {
	// 	fmt.Printf("channel closed\n")
	// 	return nil
	// }

	// switch ciq := r.Value.(type) {
	// case *data.ClientIQ:
	// 	if ciq.Type != "result" {
	// 		fmt.Printf("ciq is not result: %#v\n", ciq)
	// 		return nil
	// 	}
	// 	var ditems data.DiscoveryItemsQuery
	// 	if err := xml.NewDecoder(bytes.NewBuffer(ciq.Query)).Decode(&ditems); err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 		return nil
	// 	}
	// 	fmt.Printf("%#v\n", ditems)
	// }

	// // disco#items to the server
	// // for each item, disco#info
	// // for the ones that are chat services, do
	// // disco#items
	// // and then finally disco#info on each room

	// return result
	return nil
}
