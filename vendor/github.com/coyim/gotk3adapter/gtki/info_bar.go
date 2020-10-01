package gtki

type InfoBar interface {
	Box

	AddActionWidget(Widget, ResponseType)
	AddButton(string, ResponseType)
	SetDefaultResponse(ResponseType)
	SetMessageType(MessageType)
	GetMessageType() MessageType
	GetActionArea() (Widget, error)
	GetContentArea() (Box, error)
	GetShowCloseButton() bool
	SetShowCloseButton(bool)
}

func AssertInfoBar(_ InfoBar) {}
