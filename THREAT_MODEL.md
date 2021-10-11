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
      - By using physical implants or modifying programs on disk, any threat can be realized.
  - Someone with software access to the computer of the user
      - By using software implants or other kinds of malware, any threat can be realized.
  - Someone with a remote exploit against CoyIM
      - There are several possible vectors for remotely exploiting CoyIM. The two most likely scenarios are sending
      crafted network packets, or sending crafted XMPP packets. In a worst case scenario, this could result in remote
      code execution inside of CoyIM, leading to the same kind of access as the previous points. Milder exploits could
      lead to crashes, which would be a denial of service attack. There's also the possibility of information disclosure
      attacks, where an exploit could lead CoyIM to send information it shouldn't reveal.  If an exploit exists in the
      cryptography code, it's also possible that this could lead to breaks in the security properties (CIA).
  - Build servers
      - A build server can insert hostile code in the software during compilation, which can be used to execute any threat.
  - Download servers
      - A download server can insert hostile code in the binaries, which can be used to execute any threat.
  - The ISP (or inbetween network provider) of download servers
      - If the network connection to download servers is not protected, hostile code can be inserted into the binary 
      for any purpose.
  - The CoyIM developers
      - The CoyIM developers could insert hostile code in the source code for any threat. They could also distribute
      binaries that are different than the ones created by the build servers, containing hostile code not visible
      in the public source code.

- Stealing the contact list of a user
  - The users server
      - A users server will always store their contact list.
  - The ISP (or inbetween network provider) of a user
      - If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      - If the connection between a user and their server is not protected, this threat can be realized.
- Modifying the contact list for a user
  - The users server
      - The server stores the contact list and can modify it at will.
  - The ISP (or inbetween network provider) of a user
      - If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      - If the connection between a user and their server is not protected, this threat can be realized.
- Stopping a user from getting their contact list
  - A contact
      - Anyone could in theory execute a denial of service attack against an XMPP account, making the account
      and in extension the contact list inaccessible.
  - The users server
      - The server stores the contact list and can stop access to it
  - The ISP (or inbetween network provider) of a user
      - If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      - If the connection between a user and their server is not protected, this threat can be realized.

- Seeing a users availability
  - A contact
      - A contact which have been approved to subscribe to the updates will see the availability. This is how
      XMPP is meant to work, to not theoretically speaking a real threat.
  - The users server
      - The server is responsible for managing this information, so will always see it.
  - The server of a contact
      - If the contact has been approved to subscribe to the updates, the server of the contact will be the mediator
      of this information, and thus see it.
  - The server of a chat room
      - If the user has joined a chat room, the server hosting that chat room will see the availability of that user.
  - The ISP (or inbetween network provider) of a user
      - If the connection between a user and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a contact
      - If the connection between a contact and their server is not protected, this threat can be realized.
  - The ISP (or inbetween network provider) of a server
      - If the connection between a user and their server is not protected, this threat can be realized.
- Stopping a user from seeing another users availability
  - The users server
      - The server of the user will be an intermediary for all status information from contacts, and can always control
      visibility.
  - The server of a contact
      - In general, the server of a contact will be responsible for sharing the availability information for the contact.
      For this reason, that server can always control who sees what.
  - The server of a chat room
      - The chat room server serves as an intermediary for all status information for occupants in a room, and for this reason 
      can always stop that information from being transmitted. 
  - The ISP (or inbetween network provider) of a user
      - If the connection is not secure, it is possible for the network to control visibility of specific 
      availability information.
  - The ISP (or inbetween network provider) of a contact
      - If the network for a contact is not protected, it can modify the status information. However, this control is only
      generic - it can't be targeted for one user.
  - The ISP (or inbetween network provider) of a server
      - If the connection is not protected, this can hide status information.
- Modifying a users availability to one or several other users
  - The users server
      - The server of the user will be an intermediary for all status information, and can always send out different information.
  - The server of a contact
      - The server of a contact can modify the availibility for the user that the contact can see, but not in a global way.
  - The server of a chat room
      - The chat room server serves as an intermediary for all status information for occupants in a room, and for this reason 
      can always change that information .
  - The ISP (or inbetween network provider) of a user
      - If the connection is not secure, it is possible for the network to modify status information - but not targeted 
      to any specific contact.
  - The ISP (or inbetween network provider) of a contact
      - If the network for a contact is not protected, it can modify the status information being received by that contact.
  - The ISP (or inbetween network provider) of a server
      - If the connection is not protected, this can modify status information.

- Seeing who a user talks to
  - A contact
      - A contact can in general only see who a user talks to, if it's them.
  - The users server
      - The users server can see everyone that a user talks to
  - The server of a contact
      - The server of a contact can see anytime the user talks to that specific contact, but no-one else.
  - The ISP (or inbetween network provider) of a user
      - If the connection is unprotected, the network can see everyone that a user talks to.
  - The ISP (or inbetween network provider) of a contact
      - If the connection is unprotected, the network of a contact can specifically see when the user
      talks to that contact.
  - The ISP (or inbetween network provider) of a server
      - Depending on where the network sits, and if it's unprotected, it can see some of the contacts a user
      talks to, but usually not all.
- Seeing when a user is active (talking to someone)
  - A contact
      - Only if the contact is part of the communication
  - The users server
      - The users server can always see this
  - The server of a contact
      - Only if the contact is part of the communication
  - The server of a chat room
      - Only if the communication is happening in the chat room
  - The ISP (or inbetween network provider) of a user
      - If the connection is unprotected, the network can always see this.
  - The ISP (or inbetween network provider) of a contact
      - If the connection is unprotected, the network can see this if the contact is part of the communication.
  - The ISP (or inbetween network provider) of a server
      - If the connection is unprotected, sometimes the network will be able to see this, but not always.

- Seeing the content of a users communication
  - The users server
      - The users server can see all communication.
  - The server of a contact
      - The server of the contact can see all communication with that contact.
  - The server of a chat room
      - The server of a chat room can see all communication in that chat room.
  - The ISP (or inbetween network provider) of a user
      - If the connection is unprotected, the network can see the information.
  - The ISP (or inbetween network provider) of a contact
      - If the connection is unprotected, the network can see all communication with the contact.
  - The ISP (or inbetween network provider) of a server
      - If the connection is unprotected, the network can see some communication, but not always.
- Modifying the content of a users communication
  - The users server
      - The users server can modify any content.
  - The server of a contact
      - The server of a contact can modify any content to or from that contact.
  - The server of a chat room
      - The server of a chat room can modify any content in the chat room.
  - The ISP (or inbetween network provider) of a user
      - If the connection is unprotected, the network can modify any communication.
  - The ISP (or inbetween network provider) of a contact
      - If the connection is unprotected, the network can modify communication with the contact.
  - The ISP (or inbetween network provider) of a server
      - If the connection is unprotected, the network can modify some communication, but not always.
- Stopping a user from communicating
   - The same entities that can modify the content of a users communication can also stop that communication.
- Proving to a third party that a user said something
   - A strong proof that a user said something is not possible in the basic model.
   A more basic textual transcript is possible for the same entitites that can 
   collect user communication.
   
All the multi-user chat rooms share these threats. All the following chat room threats share these properties:
  - The users server
     - The users server has full access.
  - The server of a chat room
     - The server of the chat room has full access.
  - The ISP (or inbetween network provider) of a user
     - If the connection is unprotected, the network has full access.
  - The ISP (or inbetween network provider) of a server
     - If the connection is unprotected, the network between the users server and the chat room server
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
      - The users server can see the IP address of the user
  - The ISP (or inbetween network provider) of a user
      - If the connection is unprotected, the network can make the association between
      the IP address and the account address.
- Someone else impersonating a user to their contacts
  - The users server
      - The server can act as if it was the user
  - The server of a contact
      - The server of the contact can mimic the behavior of a remote user.
  - The server of a chat room
      - Inside a specific chat room, the server can impersonate any user in that room.
  - The ISP (or inbetween network provider) of a contact
      - If the connection is unprotected, the network of a contact can impersonate a user to that specific contact.
  - Anyone with the account information for the user
      - In theory, anyone with the account information of a user IS that user, so can impersonate them perfectly.
     
Various other threats:
- Correlation of user accounts
    - The server of a user, or the network of a user, might be able to observe connection events for more than
    one account from the same IP address, or connections that happen at the same time, making correlation and chaining
    between the accounts possible.
- Real world information about user accounts
    - The XMPP protocol exposes potentially personal information in various parts, including the account address resource identifier,
    and also specific meta-data about the device itself. This can expose information about a user.


## Mitigations

In general, mitigations can fall into several different categories. Some are things that we can address in CoyIM. Others
are things we have to accept because that is what the XMPP protocol requires. In some cases, we transfer the risk to
other parties. Finally, we have to accept some issues, since they might not have any good solutions.

### Categories of threats

Overall, we can divide the threats into these categories:


#### Network-based attacks

Any kind of attack that can be executed by adversaries positioned in the network at some point are covered by this
category.


#### Server-based attacks

Since XMPP and CoyIM relies on server software, many attacks are possible if you have access to such a server. These
attacks are put into this category.


#### Peer-based attacks

Some attacks are only possible for peers to execute. They are not very common, but the ones that exist belong to this
category.


#### Project and delivery attacks

All the attacks that rely on problems or access to the build servers, to the download servers (or networks) or hostile
developers are put into this category.


#### Local privilege attacks

This include the three types of attacks based on physical or software access to the computer, but also the possibility
of unknown vulnerabilities in the code-base.


#### Other attacks

Any other possibility of attack are put here, as the catch-all for anything else.


### Specific mitigations for categories of threats

CoyIM contains a large amount of mitigations against various specific or general types of attacks. Some only protect
against one single attack, while others make it possible to avoid whole categories of attacks. Not all mitigations are
covered here - we will only cover the most important ones.

It's important to note that we rely on OTR, Tor and TLS for various mitigations. These tools and specifications have
their own threat models, and you should assume that these are integrated into the below. We will not specifically not
that Tor is vulnerable to a global passive adversary, for example.


#### Network-based attacks

The network can see a lot of different things, and in many cases will also be able to modify data. Finally, the network
can execute denial of service attacks. CoyIM primarily mitigates these risks by using TLS for all connections, and by
using Tor as much as possible. The encryption of both Tor and TLS stop both attacks against the confidentiality of the
content, but also against the integrity. Finally, denial of service attacks also becomes harder to mount, since Tor
increases the difficulty in targeting a specific user.

There's still a risk for network attacks involving connections server-to-server. The only thing we can do to mitigate
this is to ensure that content is end-to-end-encrypted as much as possible, and encourage users to use servers with good
security practices.


#### Server-based attacks

Most of the server-based attacks are not possible to protect against, without actually losing compatibility with
XMPP. For this reason, we have to accept some potential problems from the server. Outside of those that we accept, the
most important mitigation against server-based attacks is end-to-end encryption using the OTR protocol. This stops
servers from being able to read or modify communication between peers.

When it comes to the server-based attacks on multi-user chat, the only mitigations we implement are based on warning the
user about the risks of entering rooms with various configurations, so that they are informed about the risks. Since
XMPP doesn't support encrypted group chat, we have to accept these risks for multi-user chat.

In terms of impersonating attacks, the server is a real problem as well. In order to protect against this, OTR allows
for verification of identities using fingerprint verification or SMP. If these features are used as part of a good opsec
strategy with all contacts, it will serve to protect against impersonation.

Correlation attacks are possible for a server to execute. However, CoyIM implements two different mitigations to avoid
this attack. The first one is based on randomizing connection and reconnection intervals between accounts, and the
second uses different Tor circuits for any account, which means that the connections for the different accounts look
like they come from different places. Finally, some XMPP clients use the same resource for different accounts. This is
usually not enough information to reliably correlate two accounts, but it will reduce the anonymity set. For this and
other reasons, CoyIM randomizes the resource on every connection.

Finally, the server for a user can find out location information about the user, including the IP address. As a
mitigation against this, CoyIM will use Tor for all connections, if possible. By default, it will not connect without
Tor available, to protect against user mistakes. In general, CoyIM will fail _closed_, not _open_.

CoyIM cannot protect against Denial of Service attacks from servers.


#### Peer-based attacks

All peer-based attacks are fundamentally based on the functioning of XMPP. As such, we have to accept them, and can't
implement any mitigations against them. There exists one peer-based attack that we do implement mitigations against, but
this will be covered in the section about "Threats against mitigations".


#### Project and delivery attacks

The attacks that use the project or delivery as vectors are protected against in a number of different ways. First, in
order to protect the downloads, these are managed using TLS connections to the servers. The downloads also provide
checksums which can give a small amount of additional security.

In order to protect against builds with injected code of some kind, we use reproducible builds to reduce the trust
necessary in the build systems. Finally, the CoyIM project is open source, and both the code itself and the recipes for
building the binary distributions are publicly available and easily inspected. In this way, we mitigate the risk that
the project developers could do something hostile. Together, reproducible builds and open source should minimize the
amount of trust in the project team needed. The open source distribution also makes it possible for technically minded
people to build their own binaries, to minimize the trust needed in the build and distribution mechanisms.


#### Local privilege attacks

In general, local privilege attacks are very hard to defend against. CoyIM tries to do a few things, but none of them
will ever be completely effective. Perhaps the most important mitigation is that the configuration file will be
encrypted, so that when CoyIM is not running, an attacker can't get access to the content of this file, even with local
privilege. Of course, once CoyIM starts, the information is decrypted and put into memory.

The OTR library that CoyIM uses also prevents timing side channel attacks, and locks information to specific memory
pages, which also prevents certain local attacks.


#### Other attacks

One attack that is not covered above is the one about finding out information about a specific user. This attack is
mitigated by CoyIM being implemented in such a way that it leaks a minimum of information through features of XMPP and
resources.

The risk of vulnerabilities in the code base are mitigated by the choice of Golang as the programming language. This
reduces the risk of many types of vulnerabilities, and specifically makes remote execution vulnerabilities much less
likely. The CoyIM project also tries to reduce the complexity of the implementation by carefully choosing and rejecting
features. In this way, the attack surface is kept as small as possible. The decision to avoid HTML, CSS and JavaScript
for the implementation of the user interface is based on the same need to reduce complexity. Since embedded browser
functionality usually has similar problems as the full browser, it can be a significant risk to use this
functionality. We decided to mitigate this risk by not using these technologies.

A common method to attack users is to convince them to click links. While this doesn't involve a direct attack on CoyIM,
we still made the decision to minimize the amounts of links in the application, to avoid the risk that users would be
fooled into clicking hostile links, which would open up attacks on other applications.


### Threats against mitigations

When introducing a mitigation, it can be important to also look at whether there are new threats that are applicable
only against that mitigation. In some cases, this leads to greater threats in one specific area, if the mitigation opens
up new "territory". If this happens, you have to balance this threat against the original threat which the mitigation
tries to address.

In CoyIM, there exists one new threat that is a side-effect of a mitigation, and two threats that are still possible
against the confidentiality. All three of these threats are related to end-to-end encryption.

In general, end-to-end encryption can protect the confidentiality and integrity of messages, but if an attacker gets
hold of any encryption keys involved in the conversation, they can decrypt or modify content. However, there exists two
special cases - one where an attacker collects cipher text, and then at some later point gets access to the current
encryption material from a user. However, CoyIM actually mitigates against this threat by its use of OTR. OTR provides a
property known as forward secrecy, which continuously ratchets key material forward and throws away the old
material. Since this process is one-way, having the encryption material at one point will not allow the attacker to get
eaerlier encryption material.

The other scenario is where an attacker at one point gets access to the encryption material, and then later starts
collecting cipher text. In this scenario, in a traditional setting an attacker would be able to always continue
decrypting. However, once again, OTR mitigates this attack by continuously creating new key material from scratch and
mixing that into the process. In this way, old key material will not be sufficient to decrypt newer cipher texts. This
property is called post-compromise security (and sometimes backwards secrecy).

Now, one threat that CoyIM can't easily protect against are cryptographically enhanced impersonation attacks. This
happens if an attacker manages to steal your long-lived keys and then uses that to talk to another individual. Because
of the long-lived keys, the other user will see a person that seems to have access to your private keys. Thus, the trust
from that person will likely be increased. CoyIM tries to mitigate this risk by making it harder for an attacker to
access the long-term keys.

Finally, when using authenticated end-to-end encryption, it is often possible for an attacker (which is also the person
you're talking to) to prove to someone else that you said something. In a naive setting with cryptographic signatures
over content, this can be clearly seen. You could simply show the signature and the content to someone else, and they
would know that you said it. Or, you could not take back saying it. This property is called non-repudiation. It's
something a very useful property - for example for digital contracts. But for conversations it goes against our
intuitions, and it's not great that introducing more security could lead to this kind of attack. In CoyIM, this is also
mitigated by our use of OTR, which is designed to be deniable. What that means is that after the messages have been
sent, someone else could in theory forge these messages. That means that there's no proof that any single private key is
responsible for the conversation.


## Summary

CoyIM has a fairly extensive threat model. It tries to protect against many kinds of risks, but due to the limitations
of XMPP, not everything is possible. This document should make it quite clear exactly where the limits for the
protection from CoyIM are. One final thing to note is that this document is based on the default configuration of
CoyIM. You can reduce this security, and by doing that expose yourself to more threats than those described here, but
the fundamental design idea for CoyIM is that you don't have to do anything to increase your security. This threat model
matches the top level of protection CoyIM can give you.
