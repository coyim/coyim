#!/usr/bin/env ruby

aliases = {
  "brl" => "Bruce Leidl",
  "Fab Torchz" => "Fan Jiang",
  "Fab Torchz J" => "Fan Jiang",
  "fanjiang" => "Fan Jiang",
  "Fan Jiang Torchz" => "Fan Jiang",
  "Pedro Enrique Palau" => "Pedro Palau",
  "Reinaldo de Souza Jr" => "Reinaldo de Souza Junior",
  "sacurio" => "Sandy Acurio",
  "Sandy" => "Sandy Acurio",
  "cnaranjo" => "Cristian Naranjo",
  "ivanjijon" => "Ivan Jijon",
  "mvelasco" => "Mauro Velasco",
  "piratax007" => "Fausto",
}

incorrect = {
  ["cnaranjo", "mauro.velasco@gmail.com"] => true,
}

all = `git log --format='%aN  -  %aE' | sort -u`

sorted = { }

sorted["Adam Langley"] = ""

all.each_line do |ll|
  if ll.strip != ""
    name, mail = ll.strip.split("  -  ")
    unless incorrect[[name, mail]]
      sorted[aliases[name] || name] = mail
    end
  end
end

sorted["Fausto"] = "fausto@autonomia.digital"
sorted["Cristian Naranjo"] = "cnaranjo@autonomia.digital"

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
