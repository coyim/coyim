package gui

type mucRoomConfigListOwnersForm struct {
	*roomConfigListForm
}

func newMUCRoomConfigListOwnersForm(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm {
	return &mucRoomConfigListOwnersForm{
		newRoomConfigListForm("MUCRoomConfigListFormOwners", nil, onFieldChanged, onFieldActivate),
	}
}
