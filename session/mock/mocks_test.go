package mock

import (
	"time"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MockSuite struct{}

var _ = Suite(&MockSuite{})

func (s *MockSuite) Test_SessionMock(c *C) {
	sm := &SessionMock{}

	c.Assert(sm.ApprovePresenceSubscription(nil, ""), IsNil)

	sm.AwaitVersionReply(nil, "")
	sm.Close()
	sm.AutoApprove("")

	c.Assert(sm.CommandManager(), IsNil)
	c.Assert(sm.Config(), IsNil)
	c.Assert(sm.Conn(), IsNil)
	c.Assert(sm.Connect("", nil), IsNil)
	c.Assert(sm.ConversationManager(), IsNil)
	c.Assert(sm.DenyPresenceSubscription(nil, ""), IsNil)
	c.Assert(sm.DisplayName(), Equals, "")

	_, _, _ = sm.EncryptAndSendTo(nil, "")

	c.Assert(sm.GetConfig(), IsNil)
	c.Assert(sm.GetInMemoryLog(), IsNil)
	c.Assert(sm.GroupDelimiter(), Equals, "")

	sm.HandleConfirmOrDeny(nil, false)

	c.Assert(sm.IsConnected(), Equals, false)
	c.Assert(sm.IsDisconnected(), Equals, false)
	c.Assert(sm.ManuallyEndEncryptedChat(nil), IsNil)
	c.Assert(sm.PrivateKeys(), IsNil)
	c.Assert(sm.R(), IsNil)

	sm.ReloadKeys()
	sm.RemoveContact("")

	c.Assert(sm.RequestPresenceSubscription(nil, ""), IsNil)
	c.Assert(sm.Send(nil, "", false), IsNil)
	c.Assert(sm.SendMUCMessage("", "", ""), IsNil)

	sm.SendPing()
	sm.SetCommandManager(nil)
	sm.SetConnector(nil)
	sm.SetLastActionTime(time.Time{})
	sm.SetWantToBeOnline(false)
	sm.Subscribe(nil)
	sm.Timeout(data.Cookie(0), time.Time{})
	sm.StartSMP(nil, "", "")
	sm.FinishSMP(nil, "")
	sm.AbortSMP(nil)
	sm.PublishEvent(nil)
	sm.SendIQError(nil, nil)
	sm.SendIQResult(nil, nil)

	c.Assert(sm.SendFileTo(nil, "", nil, nil), IsNil)
	c.Assert(sm.SendDirTo(nil, "", nil, nil), IsNil)
	c.Assert(sm.CreateSymmetricKeyFor(nil), IsNil)
	c.Assert(sm.GetAndWipeSymmetricKeyFor(nil), IsNil)

	_, _ = sm.HasRoom(nil, nil)
	sm.GetRoomListing(nil, nil)
	sm.GetRoomInformation(nil)
	_, _, _ = sm.GetRooms(jid.NewDomain(""), "")

	c.Assert(sm.JoinRoom(nil, "", ""), IsNil)
	rc, ec := sm.CreateInstantRoom(nil)
	c.Assert(rc, IsNil)
	c.Assert(ec, IsNil)

	_, _ = sm.CreateReservedRoom(nil)
	_, _ = sm.SubmitRoomConfigurationForm(nil, nil)

	c.Assert(sm.CancelRoomConfiguration(nil), IsNil)

	_, _, _ = sm.GetChatServices(jid.NewDomain(""))
	_, _ = sm.DestroyRoom(nil, "", nil, "")
	_, _ = sm.UpdateOccupantAffiliation(nil, "", nil, nil, "")
	_, _ = sm.UpdateOccupantRole(nil, "", nil, "")

	c.Assert(sm.Log(), IsNil)

	_, _ = sm.LeaveRoom(nil, "")

	c.Assert(sm.NewRoom(nil), IsNil)
}
