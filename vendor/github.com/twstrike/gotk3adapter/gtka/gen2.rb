#!/usr/bin/env ruby

types = %w[
AboutDialog
AccelGroup
Adjustment
Application
ApplicationWindow
Assistant
Box
Builder
Button
CellRendererText
CellRendererToggle
CheckButton
CheckMenuItem
ComboBox
ComboBoxText
CssProvider
Dialog
Entry
FileChooserDialog
Grid
InfoBar
Label
ListStore
Menu
MenuBar
MenuItem
MessageDialog
Notebook
Revealer
ScrolledWindow
SeparatorMenuItem
TextBuffer
TextIter
TextTag
TextTagTable
TextView
TreeIter
TreePath
TreeSelection
TreeStore
TreeView
Widget
Window
]

types.each do |tp|
  puts <<METH
func marshal#{ tp }(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return unwrap#{ tp }(v.(gtki.#{ tp }))
}

func unmarshal#{ tp }(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return wrap#{ tp }Simple(v.(*gtk.#{ tp }))
}

METH
end
