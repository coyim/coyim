package gui

type mucRoomConfigListAdminsForm struct {
	*roomConfigListForm
}

func newMUCRoomConfigListAdminsForm(onFieldChanged, onFieldActivate func()) mucRoomConfigListForm {
	return &mucRoomConfigListAdminsForm{
		newRoomConfigListForm("MUCRoomConfigListFormAdmins", nil, onFieldChanged, onFieldActivate),
	}
}
