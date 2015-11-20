def parse_go_name(file_name)
file_name.gsub(/[A-Z]/, '_\0')
    .gsub(/\.xml/, '.go')
    .downcase
    .gsub(/\/_/, '/')
end

def remove_if_exists(go_file)
  File.delete go_file if File.exist? go_file
end

def gen_go_file(xml_file, go_file)
  source = File.open xml_file
  ui_name = File.basename(xml_file, '.xml')
  t = ui_name.gsub(/^[A-Z]/) { |c| c.downcase }
  xml_definition = source.read
  template = """
package definitions

func init(){
  add(`#{ui_name}`, &#{t}{})
}

type #{t} struct{}

func (w #{t}) String() string {
	return `
#{xml_definition}
`
}
"""
  target = File.new(go_file, 'w+')
  target.puts template
end

Dir['./*.xml'].each do |file_name|
  go_file = parse_go_name file_name
  remove_if_exists go_file
  STDERR.puts "  - #{file_name} -> #{go_file}"
  gen_go_file(file_name, go_file)
end

