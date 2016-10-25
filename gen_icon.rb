#!/usr/bin/env ruby

r = File.binread(ARGV[0])

puts r.chars.map() {|v| "%02x" % v.ord }.join
