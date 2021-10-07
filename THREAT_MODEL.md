# Threat Model

This is a living document. The goal is to describe the kinds of threats that CoyIM will protect against, and which ones
it will not. This initial version is fairly bare. We hope to expand it in the future. This is certainly not exhaustive
yet. It is important to point out that this Threat Model is for the CoyIM software. It is _not_ a threat model for XMPP,
nor for TLS, Tor or OTR, although we sometimes reference these issues in this document.


## Assets

The assets are the things that CoyIM tries to protect. In this threat model, encryption keys and passwords are _not_
assets. They do not have any intrinsic value - instead, they are only as valuable as what they are actually protecting.

- The contact list
- Availability (when the user is online or offline, busy or occupied)
- Who a user talks to
- When a user talks to someone
- The content of a users communication
- The integrity of the content of a users communication
- Confirmation that a user said something
- What multi-user chat rooms a user takes part in
- The IP address of a user
- The account address (JID) of a user


## Adversaries

Adversaries are the entities that could potentially attack some of the assets. In some of these cases, the adversaries
are bunched together. For example, a "server" adversary doesn't just mean the server, but anyone that have access to the
server, such as the server hosting company, the system administrator, the XMPP administrator, anyone that manages to
break in to the server, etc.

- A contact
- The users server
- The server of a contact
- The server of a chat room
- The ISP (or inbetween network provider) of a user
- The ISP (or inbetween network provider) of a contact
- The ISP (or inbetween network provider) of a server
- Someone with physical access to the computer of the user
- Someone with software access to the computer of the user
- Someone with a remote exploit against CoyIM
- Build servers
- Download servers
- The ISP (or inbetween network provider) of download servers
- The CoyIM developers


## Attacks

The attacks listed below are focused on specific assets, but without taking in to consideration the adversary who could
mount the attack. Combining the attacks with the adversaries give us the threats, which we will look at in the next
section.

- Stealing the contact list of a user
- Modifying the contact list for a user
- Stopping a user from getting their contact list

- Seeing a users availability
- Stopping a user from seeing another users availability
- Modifying a users availability to one or several other users

- Seeing who a user talks to
- Seeing when a user is active (talking to someone)

- Seeing the content of a users communication
  - Through an attack executed while the conversation is in progress
  - Through a compromise that happens _after_ the conversation happens
  - Through a compromise that happened _before_ a conversation happens
- Modifying the content of a users communication
- Stopping a user from communicating
- Proving to a third party that a user said something

- Seeing what multi-user chat rooms a user has joined
- Seeing when a user joins or leaves a multi-user chat room
- Seeing what a user says in a multi-user chat room
- Seeing who else is in a multi-user chat room
- Seeing what a user says in a multi-user chat room
- Modifying what a user says in a multi-user chat room
- Stopping a user from communicating to a multi-user chat room
- Stopping messages or updates from/about other users in a multi-user chat room
- Rearrange order of messages and events in a multi-user chat room

- Find out the IP address of a user
- Someone else impersonating a user to their contacts
   (This is not actually a separate attack - more of a combination of several of the above attacks)

In general, these attacks follow the `CIA` structure - `Confidentiality`, `Integrity` and `Availability`. There are some
attacks that lie outside of these properties, though. Some impact `deniability` and `non-repudiability`. Others are
related to `CIA` but in terms of meta-data rather than content. Finally, contact graphing and chaining are also
concerns.


## Threats

In this threat model, a threat will be treated as a specific attack on an asset, executed by a specific adversary. Not
all combinations will be covered here - instead, we will only cover the possibilities that are likely enough to
happen. If there's a question of whether the likelihood is high or low, a discussion will be included. The method of how
the threat would be most likely executed are also included here. In general, this list only contains threats - which
signify intention. Unintentional mistakes in server or client code are not included as possibilities, for example.

All of the attacks in this threat model can be executed by these threats, so they won't be repeated for all specific
attacks:
  - Someone with physical access to the computer of the user
      By using physical implants or modifying programs on disk, any threat can be realized.
  - Someone with software access to the computer of the user
      By using software implants or other kinds of malware, any threat can be realized.
  - Someone with a remote exploit against CoyIM
      There are several possible vectors for remotely exploiting CoyIM. The two most likely scenarios are sending crafted network
      packets, or sending crafted XMPP packets. In a worst case scenario, this could result in remote code execution
      inside of CoyIM, leading to the same kind of access as the previous points. Milder exploits could lead to crashes, 
      which would be a denial of service attack. There's also the possibility of information disclosure attacks, where
      an exploit could lead CoyIM to send information it shouldn't reveal.
      If an exploit exists in the cryptography code, it's also possible that this could lead to breaks in the security properties (CIA).
  - Build servers
      A build server can insert hostile code in the software during compilation, which can be used to execute any threat.
  - Download servers
      A download server can insert hostile code in the binaries, which can be used to execute any threat.
  - The ISP (or inbetween network provider) of download servers
      If the network connection to download servers is not protected, hostile code can be inserted into the binary 
      for any purpose.
  - The CoyIM developers
      The CoyIM developers could insert hostile code in the source code for any threat. They could also distribute
      binaries that are different than the ones created by the build servers, containing hostile code not visible
      in the public source code.

- Stealing the contact list of a user
  - The users server
      A users server will always store their contact list.
  - The ISP (or inbetween network provider) of a user
      If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      If the connection between a user and their server is not protected, this threat can be realized.
- Modifying the contact list for a user
  - The users server
      The server stores the contact list and can modify it at will.
  - The ISP (or inbetween network provider) of a user
      If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      If the connection between a user and their server is not protected, this threat can be realized.
- Stopping a user from getting their contact list
  - A contact
      Anyone could in theory execute a denial of service attack against an XMPP account, making the account
      and in extension the contact list inaccessible.
  - The users server
      The server stores the contact list and can stop access to it
  - The ISP (or inbetween network provider) of a user
      If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      If the connection between a user and their server is not protected, this threat can be realized.

- Seeing a users availability
  - A contact
      A contact which have been approved to subscribe to the updates will see the availability. This is how
      XMPP is meant to work, to not theoretically speaking a real threat.
  - The users server
      The server is responsible for managing this information, so will always see it.
  - The server of a contact
      If the contact has been approved to subscribe to the updates, the server of the contact will be the mediator
      of this information, and thus see it.
  - The server of a chat room
      If the user has joined a chat room, the server hosting that chat room will see the availability of that user.
  - The ISP (or inbetween network provider) of a user
      If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a contact
      If the connection between a contact and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      If the connection between a user and their server is not protected, this threat can be realized.
- Stopping a user from seeing another users availability
  - The users server
      The server of the user will be an intermediary for all status information from contacts, and can always control
      visibility.
  - The server of a contact
      In general, the server of a contact will be responsible for sharing the availability information for the contact.
      For this reason, that server can always control who sees what.
  - The server of a chat room
      The chat room server serves as an intermediary for all status information for occupants in a room, and for this reason 
      can always stop that information from being transmitted. 
  - The ISP (or inbetween network provider) of a user
      If the connection is not secure, it is possible for the network to control visibility of specific 
      availability information.
  - The ISP (or inbetween network provider) of a contact
      If the network for a contact is not protected, it can modify the status information. However, this control is only
      generic - it can't be targeted for one user.
  - The ISP (or inbetween network provider) of a server
      If the connection is not protected, this can hide status information.
- Modifying a users availability to one or several other users
  - The users server
      The server of the user will be an intermediary for all status information, and can always send out different information.
  - The server of a contact
      The server of a contact can modify the availibility for the user that the contact can see, but not in a global way.
  - The server of a chat room
      The chat room server serves as an intermediary for all status information for occupants in a room, and for this reason 
      can always change that information .
  - The ISP (or inbetween network provider) of a user
      If the connection is not secure, it is possible for the network to modify status information - but not targeted 
      to any specific contact.
  - The ISP (or inbetween network provider) of a contact
      If the network for a contact is not protected, it can modify the status information being received by that contact.
  - The ISP (or inbetween network provider) of a server
      If the connection is not protected, this can modify status information.

- Seeing who a user talks to
  - A contact
      A contact can in general only see who a user talks to, if it's them.
  - The users server
      The users server can see everyone that a user talks to
  - The server of a contact
      The server of a contact can see anytime the user talks to that specific contact, but no-one else.
  - The ISP (or inbetween network provider) of a user
      If the connection is unprotected, the network can see everyone that a user talks to.
  - The ISP (or inbetween network provider) of a contact
      If the connection is unprotected, the network of a contact can specifically see when the user
      talks to that contact.
  - The ISP (or inbetween network provider) of a server
      Depending on where the network sits, and if it's unprotected, it can see some of the contacts a user
      talks to, but usually not all.
- Seeing when a user is active (talking to someone)
  - A contact
      Only if the contact is part of the communication
  - The users server
      The users server can always see this
  - The server of a contact
      Only if the contact is part of the communication
  - The server of a chat room
      Only if the communication is happening in the chat room
  - The ISP (or inbetween network provider) of a user
      If the connection is unprotected, the network can always see this.
  - The ISP (or inbetween network provider) of a contact
      If the connection is unprotected, the network can see this if the contact is part of the communication.
  - The ISP (or inbetween network provider) of a server
      If the connection is unprotected, sometimes the network will be able to see this, but not always.

- Seeing the content of a users communication
  - The users server
      The users server can see all communication.
  - The server of a contact
      The server of the contact can see all communication with that contact.
  - The server of a chat room
      The server of a chat room can see all communication in that chat room.
  - The ISP (or inbetween network provider) of a user
      If the connection is unprotected, the network can see the information.
  - The ISP (or inbetween network provider) of a contact
      If the connection is unprotected, the network can see all communication with the contact.
  - The ISP (or inbetween network provider) of a server
      If the connection is unprotected, the network can see some communication, but not always.
- Modifying the content of a users communication
  - The users server
      The users server can modify any content.
  - The server of a contact
      The server of a contact can modify any content to or from that contact.
  - The server of a chat room
      The server of a chat room can modify any content in the chat room.
  - The ISP (or inbetween network provider) of a user
      If the connection is unprotected, the network can modify any communication.
  - The ISP (or inbetween network provider) of a contact
      If the connection is unprotected, the network can modify communication with the contact.
  - The ISP (or inbetween network provider) of a server
      If the connection is unprotected, the network can modify some communication, but not always.
- Stopping a user from communicating
   The same entities that can modify the content of a users communication can also stop that communication.
- Proving to a third party that a user said something
   A strong proof that a user said something is not possible in the basic model.
   A more basic textual transcript is possible for the same entitites that can 
   collect user communication.
   
All the multi-user chat rooms share the threats. All the following chat room threats share these properties:
  - The users server
     The users server has full access.
  - The server of a chat room
     The server of the chat room has full access.
  - The ISP (or inbetween network provider) of a user
     If the connection is unprotected, the network has full access.
  - The ISP (or inbetween network provider) of a server
     If the connection is unprotected, the network between the users server and the chat room server
     has full access.
- Seeing what multi-user chat rooms a user has joined
- Seeing when a user joins or leaves a multi-user chat room
- Seeing what a user says in a multi-user chat room
- Seeing who else is in a multi-user chat room
- Seeing what a user says in a multi-user chat room
- Modifying what a user says in a multi-user chat room
- Stopping a user from communicating to a multi-user chat room
- Stopping messages or updates from/about other users in a multi-user chat room
- Rearrange order of messages and events in a multi-user chat room

- Find out the IP address of a user
  - The users server
      The users server can see the IP address of the user
  - The ISP (or inbetween network provider) of a user
      If the connection is unprotected, the network can make the association between
      the IP address and the account address.
- Someone else impersonating a user to their contacts
  - The users server
      The server can act as if it was the user
  - The server of a contact
      The server of the contact can mimic the behavior of a remote user.
  - The server of a chat room
      Inside a specific chat room, the server can impersonate any user in that room.
  - The ISP (or inbetween network provider) of a contact
      If the connection is unprotected, the network of a contact can impersonate a user to that specific contact.
  - Anyone with the account information for the user
      In theory, anyone with the account information of a user IS that user, so can impersonate them perfectly.
     
Various other threats:
- Correlation of user accounts
    The server of a user, or the network of a user, might be able to observe connection events for more than
    one account from the same IP address, or connections that happen at the same time, making correlation and chaining
    between the accounts possible.
- Real world information about user accounts
    The XMPP protocol exposes potentially personal information in various parts, including the account address resource identifier,
    and also specific meta-data about the device itself. This can expose information about a user.


## Mitigations

In general, mitigations can fall into several different categories. Some are things that we can address in CoyIM. Others
are things we have to accept because that is what the XMPP protocol requires. In some cases, we transfer the risk to
other parties. Finally, we have to accept some issues, since they might not have any good solutions.

### Categories of threats

### Specific mitigations for categories of threats

### Threats against mitigations

  - Through a compromise that happens _after_ the conversation happens
    - A contact
    - The users server
    - The server of a contact
    - The server of a chat room
    - The ISP (or inbetween network provider) of a user
    - The ISP (or inbetween network provider) of a contact
    - The ISP (or inbetween network provider) of a server
  - Through a compromise that happened _before_ a conversation happens
    - A contact
    - The users server
    - The server of a contact
    - The server of a chat room
    - The ISP (or inbetween network provider) of a user
    - The ISP (or inbetween network provider) of a contact
    - The ISP (or inbetween network provider) of a server
- Proving to a third party that a user said something
  - A contact
  - The users server
  - The server of a contact
  - The server of a chat room
  - The ISP (or inbetween network provider) of a user
  - The ISP (or inbetween network provider) of a contact
  - The ISP (or inbetween network provider) of a server


## Summary - current threat model

(Who can get what with current model of CoyIM)
