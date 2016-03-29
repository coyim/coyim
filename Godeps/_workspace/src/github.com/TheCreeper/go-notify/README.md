# go-notify

[![GoDoc](https://godoc.org/github.com/TheCreeper/go-notify?status.svg)](https://godoc.org/github.com/TheCreeper/go-notify)

Package notify provides an implementation of the Gnome DBus
[Notifications Specification](https://developer.gnome.org/notification-spec).

## Examples

Display a simple notification.
```Go
ntf := notify.NewNotification("Test Notification", "Just a test")
if _, err := ntf.Show(); err != nil {
	return
}
```

Display a notification with an icon. Consult the
[Icon Naming Specification](http://standards.freedesktop.org/icon-naming-spec).
```Go
ntf := notify.NewNotification("Test Notification", "Just a test")
ntf.AppIcon = "network-wireless"
if _, err := ntf.Show(); err != nil {
	return
}
```

Display a notification that never expires.
```Go
ntf := notify.NewNotification("Test Notification", "Just a test")
ntf.Timeout = notify.ExpiresNever
if _, err := ntf.Show(); err != nil {
	return
}
```

Play a sound with the notification.
```Go
ntf := notify.NewNotification("Test Notification", "Just a test")
ntf.Hints = make(map[string]interface{})
ntf.Hints[notify.HintSoundFile] = "/home/my-username/sound.oga"
if _, err := ntf.Show(); err != nil {
	return
}
```
