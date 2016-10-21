#!/usr/bin/env ruby

types = %w[
AboutDialog
AccelGroup
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
TreeStore
TreeView
Widget
Window
]

class String
  def underscore
    self.gsub(/::/, '/').
    gsub(/([A-Z]+)([A-Z][a-z])/,'\1_\2').
    gsub(/([a-z\d])([A-Z])/,'\1_\2').
    tr("-", "_").
    downcase
  end
end

types.each do |tp|
  lower = tp[0].downcase + tp[1..-1]
  fname = "#{tp.underscore}.go"
  File.open(fname, "w") do |ff|
    ff.puts <<METH
package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gtki"
)

type #{ lower } struct {
	*gtk.#{ tp }
}

func wrap#{ tp }(v *gtk.#{ tp }, e error) (*#{ lower }, error) {
	if v == nil {
		return nil, e
	}
	return &#{ lower }{v}, e
}

func unwrap#{ tp }(v gtki.#{ tp }) *gtk.#{ tp } {
	if v == nil {
		return nil
	}
	return v.(*#{ lower }).#{ tp }
}
METH
  end
end
