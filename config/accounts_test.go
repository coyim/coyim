package config

import (
	. "gopkg.in/check.v1"
)

type AccountsSuite struct{}

var _ = Suite(&AccountsSuite{})

func (s *AccountsSuite) Test_Accounts_RemoveAccount(c *C) {
	ac1 := &Account{Account: "account@one.com"}
	ac2 := &Account{Account: "account@two.com"}
	acs := Accounts{
		Accounts: []*Account{ac1, ac2},
	}

	acs.Remove(ac1)

	c.Check(len(acs.Accounts), Equals, 1)
	_, found := acs.GetAccount("account@two.com")
	c.Check(found, Equals, true)
}

func (s *AccountsSuite) Test_Accounts_DontRemoveWhenDoesntExist(c *C) {
	ac1 := &Account{Account: "account@one.com"}
	ac2 := &Account{Account: "account@two.com"}
	ac3 := &Account{Account: "nohay@anywhere.com"}
	acs := Accounts{
		Accounts: []*Account{ac1, ac2},
	}

	acs.Remove(ac3)

	c.Check(len(acs.Accounts), Equals, 2)
	_, found := acs.GetAccount("account@two.com")
	c.Check(found, Equals, true)
}
