package gui

type navigationIconMapper interface {
	navigationIconByPageID(id mucRoomConfigPageID) navigationItemIconName
}

type navigationItemIconSetByPageID map[mucRoomConfigPageID]navigationItemIconName
type iconSet struct {
	navigationIcons navigationItemIconSetByPageID
}

func (is iconSet) navigationIconByPageID(id mucRoomConfigPageID) navigationItemIconName {
	return is.navigationIcons[id]
}

func (cm *hasColorManagement) currentIconSet() navigationIconMapper {
	if cm.isDarkThemeVariant() {
		return defaultDarkIconSet
	}
	return defaultLightIconSet
}
