#!/usr/bin/env ruby

types = %w[
  Event
  EventButton
  Pixbuf
  PixbufLoader
  Screen
]

exportedWrap = {}
exportedUnwrap = {
  "Screen" => true,
  "Pixbuf" => true
}

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
  prefix1 = if exportedWrap[tp]
             "W"
           else
             "w"
           end
  prefix2 = if exportedUnwrap[tp]
             "U"
           else
             "u"
           end

  File.open(fname, "w") do |ff|
    ff.puts <<METH
package gdka

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
)

type #{ lower } struct {
	*gdk.#{ tp }
}

func #{prefix1}rap#{ tp }(v *gdk.#{ tp }, e error) (*#{ lower }, error) {
	if v == nil {
		return nil, e
	}
	return &#{ lower }{v}, e
}

func #{prefix2}nwrap#{ tp }(v gdki.#{ tp }) *gdk.#{ tp } {
	if v == nil {
		return nil
	}
	return v.(*#{ lower }).#{ tp }
}
METH
  end
end
