def parse_go_name(file_name)
file_name.gsub(/[A-Z]/, '_\0')
    .gsub(/gui\/definitions/, 'gui')
    .gsub(/\.xml/, '.go')
    .downcase
    .gsub(/\/_/, '/')
end

def remove_if_exists(go_file)
  File.delete go_file if File.exist? go_file
end

def gen_go_file(xml_file, go_file)
  source = File.open xml_file
  t = File.basename(xml_file, '.xml')
    .gsub(/^[A-Z]/) { |c| c.downcase }
  xml_definition = source.read
  template ="""
package gui

type #{t}  struct{}

func (w #{t}) getDefinition() string {
	return `
#{xml_definition}
`
}
"""
  target = File.new(go_file, 'w+')
  target.puts template
end

Dir['./gui/definitions/*.xml'].each do |file_name|
  puts "found #{file_name}"
  go_file = parse_go_name file_name
  remove_if_exists go_file
  puts "reading #{file_name} and generating #{go_file}"
  gen_go_file(file_name, go_file)
end

