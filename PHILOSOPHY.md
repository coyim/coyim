# The philosophy, mission, guide lines and goals of CoyIM

This document aims to expand a bit on the philosophy we use when developing CoyIM - what the mission is, and guide lines for what features we choose to include or not include. It is not a manifesto, and it's an evolving document.

The first thing to keep in mind is that CoyIM is not primarily a general purpose chat client. We are not aiming to compete with Facebook Messenger, or Skype, or Pidgin, or WhatsApp, or Slack, or many of the other clients out there. Or at least, we don't aim to compete with them for general purpose chatting needs. We want to provide an alternative for secure and private chat - something that will scale up to the security needs of the high risk targets. This perspective guides every decision we make. If we are thinking about implementing something that will have a negative privacy or security impact, we will probably not do it.

Our choice of XMPP as the chat protocol also doesn't mean that we aim to be a fully featured Jabber client. We will try our hardest to NOT implement an XEP or feature, unless we really need it. This puts us in a very different position compared to applications like Gajim, Pidgin, Dino and all the others - these applications all are aiming to be general purpose chat clients, and are implementing a large number of standards in order to get there. CoyIM can't really be compared to these, because our philosophy is basically the opposite.

The mission of CoyIM is quite simple - we want to be as secure as possible, from the moment you start up. We want to leak the minimum information necessary for you to be able to communicate what you need to communicate. We will sometimes make it possible to configure settings to _lower_ your security - but we will never ask you to configure something to _increase_ your security.

Another aspect of being as secure as possible for high risk targets involves the need to have good usability. This is why we want to start from a default of "as-secure-as-possible", so fewer things can go wrong. We should test our intuitions of how people interact with our application, and continuously see what we can improve. 

Sometimes the discussion about adding another chat protocol would be a good idea - for example Ricochet. Or, should we add support for OMEMO? In general, we resist doing those things. In some cases we have specific reasons for not wanting to add a specific protocol - but we also have a general statistical argument: The more complexity, the more risk. And adding new protocols adds a lot of complexity. So even though Ricochet is really cool, adding support for it would also add significantly more complexity to the project, and that's something we would really like to avoid.

