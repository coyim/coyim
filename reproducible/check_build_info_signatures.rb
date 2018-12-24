#!/usr/bin/env ruby

require 'tmpdir'
require 'open-uri'
require 'rubygems'
require 'nokogiri'

TAG = ARGV[0]

def tag_exists?(t)
  begin
    open("https://dl.bintray.com/coyim/coyim-bin/#{TAG}") { |f| }
    open("https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux") { |f| }
    open("https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux/amd64") { |f| }
    return true
  rescue OpenURI::HTTPError
    return false
  end
end

Dir.mktmpdir { |dir|
  if !tag_exists?(TAG)
    puts "Tag #{TAG} doesn't exist"
    exit
  end
  entries = []
  open("https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux/amd64") { |f|
    entries = Nokogiri(f).xpath("//pre/a").map(&:content)
  }

  if !entries.include?("coyim")
    puts "Tag #{TAG} doesn't contain coyim
    exit
  end

  if !entries.include?("build_info")
    puts "Tag #{TAG} doesn't contain a build_info file"
    exit
  end

  if entries.select{|x| x.start_with?("build_info.")}.length == 0
    puts "Tag #{TAG} doesn't contain any build_info signatures"
    exit
  end

  $stdout.print "Downloading files "; $stdout.flush
  $stdout.print "."; $stdout.flush
  `curl -L -s https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux/amd64/coyim -o #{dir}/coyim`
  $stdout.print "."; $stdout.flush
  `curl -L -s https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux/amd64/build_info -o #{dir}/build_info`
  $stdout.print "."; $stdout.flush
  entries.select{|x| x.start_with?("build_info.")}.each do |xx|
    `curl -L -s https://dl.bintray.com/coyim/coyim-bin/#{TAG}/linux/amd64/#{xx} -o #{dir}/#{xx}`
    $stdout.print "."; $stdout.flush
  end
  puts
  puts "Download finished"
  reg_hash = nil
  open("#{dir}/build_info") { |ff|
    content = ff.read
    reg_hash = content[/^([a-f0-9]{64})  \/builds\/coyim$/, 1]
  }
  real_reg_sum = `sha256sum #{dir}/coyim`[/^([a-f0-9]{64})/, 1]

  if reg_hash != real_reg_sum
    puts "Hash for coyim doesn't match - #{reg_hash} vs #{real_reg_sum}"
    exit
  end

  correct = 0
  incorrect = 0
  entries.select{|x| x.start_with?("build_info.")}.each { |bs|
    puts "Verifying #{bs}"
    res = system("gpg2 --verify #{dir}/#{bs} #{dir}/build_info >/dev/null 2>&1")
    if res
      correct += 1
    else
      incorrect += 1
    end
  }
  puts "  #{correct} verified - #{incorrect} didn't"
  if incorrect > 0
    puts "VERIFICATION FAILED"
  else
    puts "VERIFICATION SUCCEEDED"
  end
}
