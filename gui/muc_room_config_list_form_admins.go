package gui

type mucRoomConfigListAdminsForm struct {
	*roomConfigListForm
}

func newMUCRoomConfigListAdminsForm(onFieldChanged func()) mucRoomConfigListForm {
	return &mucRoomConfigListAdminsForm{
		newRoomConfigListForm("MUCRoomConfigListFormAdmins", nil, onFieldChanged),
	}
}
