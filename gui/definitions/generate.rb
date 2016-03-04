#!/usr/bin/env ruby

require 'fileutils'

def parse_go_name(file_name)
  File.basename(file_name, ".xml").
    gsub(/([A-Z]+)([A-Z][a-z])/,'\1_\2').
    gsub(/([a-z\d])([A-Z])/,'\1_\2').
    tr("-", "_").
    gsub(/\/_/, '/').
    downcase + ".go"
end

def gen_go_file(xml_file, go_file)
  xml_definition = File.read(xml_file)
  ui_name = File.basename(xml_file, '.xml')
  File.open(go_file, 'w+') do |target|
    target.puts <<TEMPLATE
package definitions

func init() {
\tadd(`#{ui_name}`, &def#{ui_name}{})
}

type def#{ui_name} struct{}

func (*def#{ui_name}) String() string {
\treturn `#{xml_definition}`
}
TEMPLATE
  end
end

def file_mtime(nm)
  return Time.at(0) unless File.exists?(nm)
  File.mtime(nm)
end

Dir[File.join(File.dirname(__FILE__), '*.xml')].each do |file_name|
  go_file = parse_go_name file_name
  if file_mtime(file_name) > file_mtime(go_file) || file_mtime(__FILE__) > file_mtime(go_file)
    STDERR.puts "  - #{file_name} -> #{go_file}"
    gen_go_file(file_name, go_file)
  end
end
