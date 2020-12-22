package gui

type mucRoomConfigListOwnersForm struct {
	*roomConfigListForm
}

func newMUCRoomConfigListOwnersForm(onFieldChanged func()) mucRoomConfigListForm {
	return &mucRoomConfigListOwnersForm{
		newRoomConfigListForm("MUCRoomConfigListFormOwners", nil, onFieldChanged),
	}
}
