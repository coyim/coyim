#!/usr/bin/env ruby

def gen_go_file
  binary_definition = File.binread("gschemas.compiled")
  hex = binary_definition.each_byte.map { |b| "%02x" % b }.join

  File.open("schemas.go", "w") do |f|
    sliced = hex.chars.each_slice(80).map{ |s| s.join }.join "\"+\n\t\""

    f.puts <<TEMPLATE
package definitions

const schemaDefinition = ""+
\t"#{sliced}"
TEMPLATE
  end
end

gen_go_file
