#!/usr/bin/env ruby

require 'tmpdir'
require 'open-uri'
require 'rubygems'
require 'nokogiri'

TAG = ARGV[0]

def tag_exists?(t)
  begin
    open("https://dl.bintray.com/twstrike/coyim/#{TAG}") { |f| }
    open("https://dl.bintray.com/twstrike/coyim/#{TAG}/linux") { |f| }
    open("https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64") { |f| }
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
  open("https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64") { |f|
    entries = Nokogiri(f).xpath("//pre/a").map(&:content)
  }

  if !entries.include?("coyim") || !entries.include?("coyim-cli")
    puts "Tag #{TAG} doesn't contain coyim or coyim-cli"
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
  `curl -L -s https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64/coyim -o #{dir}/coyim`
  $stdout.print "."; $stdout.flush
  `curl -L -s https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64/coyim-cli -o #{dir}/coyim-cli`
  $stdout.print "."; $stdout.flush
  `curl -L -s https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64/build_info -o #{dir}/build_info`
  $stdout.print "."; $stdout.flush
  entries.select{|x| x.start_with?("build_info.")}.each do |xx|
    `curl -L -s https://dl.bintray.com/twstrike/coyim/#{TAG}/linux/amd64/#{xx} -o #{dir}/#{xx}`
    $stdout.print "."; $stdout.flush
  end
  puts
  puts "Download finished"
  cli_hash, reg_hash = nil
  open("#{dir}/build_info") { |ff|
    content = ff.read
    cli_hash = content[/^([a-f0-9]{64})  \/builds\/coyim-cli$/, 1]
    reg_hash = content[/^([a-f0-9]{64})  \/builds\/coyim$/, 1]
  }
  real_cli_sum = `sha256sum #{dir}/coyim-cli`[/^([a-f0-9]{64})/, 1]
  real_reg_sum = `sha256sum #{dir}/coyim`[/^([a-f0-9]{64})/, 1]

  if cli_hash != real_cli_sum
    puts "Hash for coyim-cli doesn't match - #{cli_hash} vs #{real_cli_sum}"
    exit
  end

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
