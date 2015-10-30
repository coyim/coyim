#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import signal
import sys

from dogtail.procedural import *
from dogtail.tc import TCNode
from dogtail.tree import *
tcNode = TCNode()

def addAccount(account):
    coyim = root.application("CoyIM")

    accountDialog = coyim.dialog("Account Details")
    accountDialog.child(roleName = "text").typeText(account)
    accountDialog.child(roleName = "password text").typeText("12345")
    accountDialog.button("Save").doActionNamed("click")

if __name__ == "__main__":
  coyimPath = sys.argv[1]
  pid = run(coyimPath)

  try:
    coyim = root.application("CoyIM")

    encryptAlert = coyim.child(name = "Question", roleName = "alert")
    encryptAlert.button("Yes").doActionNamed("click")

    accountName = "coyim_test@dukgo.com"
    addAccount(accountName)

    masterPassDialog = coyim.dialog("Enter master password")
    masterPassDialog.child(roleName = "password text").typeText("12345")
    masterPassDialog.button("OK").doActionNamed("click")

    accountMenu = coyim.menu("Accounts")
    accountItem = accountMenu.child(accountName)
    tcNode.compare("CoyIM has an account %s" % accountName, None, accountItem)

    os.kill(pid, signal.SIGTERM)

  except Exception as e:
    print "CoyIM died....\n"
    print e
    os.kill(pid, signal.SIGTERM)

