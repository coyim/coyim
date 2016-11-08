#!/usr/bin/env ruby

aliases = {
  "brl" => "Bruce Leidl",
  "Fab Torchz" => "Fan Jiang",
  "Fab Torchz J" => "Fan Jiang",
  "fanjiang" => "Fan Jiang",
  "Fan Jiang Torchz" => "Fan Jiang",
  "Reinaldo de Souza Jr" => "Reinaldo de Souza Junior",
}

all = `git log --format='%aN  -  %aE' | sort -u`

sorted = { }

sorted["Adam Langley"] = ""

all.each_line do |ll|
  if ll.strip != ""
    name, mail = ll.strip.split("  -  ")
    sorted[aliases[name] || name] = mail
  end
end

res = sorted.to_a.sort.map{ |l, r|
  if r != ""
    "\"#{l}  -  #{r}\""
  else
    "\"#{l}\""
  end
}.join(",\n        ") + ",\n"

puts <<EOF
package gui

func authors() []string {
    return []string{
        #{res}
    }
}
EOF
