package scram

import (
	/* #nosec G501 */
	"crypto/sha1"
	"crypto/sha256"

	"github.com/coyim/coyim/sasl"

	. "gopkg.in/check.v1"
)

type ScramSuite struct{}

var _ = Suite(&ScramSuite{})

func (s *ScramSuite) TestScramWithRFC5802TestVector(c *C) {
	mech := &scramMechanism{sha1.New, sha1.Size, false, false}
	client := mech.NewClient()

	_ = client.SetProperty(sasl.AuthID, "user")
	_ = client.SetProperty(sasl.Password, "pencil")
	_ = client.SetProperty(sasl.ClientNonce, "fyko+d2lbbFgONRv9qkxdawL")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`n,,n=user,r=fyko+d2lbbFgONRv9qkxdawL`))

	rec := sasl.Token("r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,s=QSXCR+Q6sek8bf92,i=4096")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`c=biws,r=fyko+d2lbbFgONRv9qkxdawL3rfcNHYJY1ZVvWVs7j,p=v0X8v3Bz2T0CJGbJQyF0X+HI4Ts=`))

	rec = sasl.Token("v=rmF9pqV8S7suAoZWja4dJRkFsKQ=")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, IsNil)
}

func (s *ScramSuite) TestScramWithRFC5802TestVectorAndSupport(c *C) {
	mech := &scramMechanism{sha1.New, sha1.Size, false, true}
	client := mech.NewClient()

	_ = client.SetProperty(sasl.AuthID, "user")
	_ = client.SetProperty(sasl.Password, "pencil")
	_ = client.SetProperty(sasl.ClientNonce, "b5cff6190013e6")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`y,,n=user,r=b5cff6190013e6`))

	rec := sasl.Token("r=b5cff6190013e6kmua5DEmPAZjOMpq4THWJQ==,s=DkEWM20qE19suvrHhPr7HA==,i=4096")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`c=eSws,r=b5cff6190013e6kmua5DEmPAZjOMpq4THWJQ==,p=Rfqg49jbBf1bGCkwFTbo/DvHmUI=`))

	rec = sasl.Token("v=k5n595qW50TyE2q+J0cUf9yT8v4=")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, IsNil)
}

func (s *ScramSuite) TestScramSha256WithRFC7677TestVector(c *C) {
	mech := &scramMechanism{sha256.New, sha256.Size, false, false}
	client := mech.NewClient()

	_ = client.SetProperty(sasl.AuthID, "user")
	_ = client.SetProperty(sasl.Password, "pencil")
	_ = client.SetProperty(sasl.ClientNonce, "rOprNGfwEbeRWgbNEkqO")

	t, err := client.Step(nil)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`n,,n=user,r=rOprNGfwEbeRWgbNEkqO`))

	rec := sasl.Token("r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF$k0,s=W22ZaJ0SNY7soEsUEjb6gQ==,i=4096")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, true)
	c.Check(t, DeepEquals, sasl.Token(`c=biws,r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF$k0,p=dHzbZapWIk4jUhN+Ute9ytag9zjfMHgsqmmiz7AndVQ=`))

	rec = sasl.Token("v=6rriTRBi23WpRR/wtup+mMhUZUn/dB5nLTJRsjl95G4=")
	t, err = client.Step(rec)

	c.Check(err, IsNil)
	c.Check(client.NeedsMore(), Equals, false)
	c.Check(t, IsNil)
}
