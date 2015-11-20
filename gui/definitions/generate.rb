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
  xml_definition = source.read
  template = """
package definitions

func init(){
  add(`#{ui_name}`, &def#{ui_name}{})
}

type def#{ui_name} struct{}

func (*def#{ui_name}) String() string {
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

