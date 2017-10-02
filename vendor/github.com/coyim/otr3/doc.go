// Package otr3 implements the Off The Record protocol as specified in:
//  https://otr.cypherpunks.ca/Protocol-v3-4.0.0.html
//
// Introduction
//
// Off-the-Record (OTR) Messaging allows you to have private conversations over instant messaging by providing:
//  Encryption
//  Authentication
//  Deniability
//  Perfect forward secrecy
//
//
// Getting Started
//
// OTR library provides a Conversation API for Receiving and Sending messages
//  import "otr3"
//
//  c := &otr3.Conversation{}
//
//  // You will need to prepare a long-term PrivateKey for otr conversation handshakes.
//  var priv *otr3.DSAPrivateKey
//  priv.Generate(rand.Reader)
//  c.SetOurKeys([]otr3.PrivateKey{priv})
//
//  // set the Policies.
//  c.Policies.RequireEncryption()
//  c.Policies.AllowV2()
//  c.Policies.AllowV3()
//  c.Policies.SendWhitespaceTag()
//  c.Policies.WhitespaceStartAKE()
//
//  // You can also setup a debug mode
//  c.SetDebug(true)
//
//  // Use Send and Receive for messages exchange
//  toSend, err := c.Send(otr3.ValidMessage("hello"))
//  plain, toSend, err := c.Receive(toSend[0])
//
//  // Use Authenticate to start a SMP process
//  toSend, err := c.StartAuthenticate([]byte{"My pet's name?"},[]byte{"Gopher"})
//  toSend, err := c.ProvideAuthenticationSecret([]byte{"Gopher"})
package otr3
