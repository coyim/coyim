package gui

// [ps] All this file MUST doesn't make sense and we should remove the
// use of the `roomConfigPositionsField` struct.

type roomConfigPositionsField struct {
	*roomConfigFieldPositions
}

// [ps] Whoever is using this initializer should now call `newRoomConfigFieldPositions` directly
func newRoomConfigPositionsField(options roomConfigPositionsOptions) hasRoomConfigFormField {
	return &roomConfigPositionsField{
		newRoomConfigFieldPositions(options),
	}
}
