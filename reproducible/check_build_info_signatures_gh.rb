#!/usr/bin/env ruby

require 'tmpdir'
require 'open-uri'
require 'json'
require 'pp'

TAG = ARGV[0]
COYIM_BASE_URL = "https://api.github.com/repos/coyim/coyim"

def download(entries)
  entries.map do |e|
    open(e["browser_download_url"]) { |f| f.read }
  end
end

Dir.mktmpdir { |dir|
  begin
    rel = open("#{COYIM_BASE_URL}/releases/tags/#{TAG}") { |f|
      JSON.load(f)
    }

    coyim_entry = nil
    build_info_entry = nil
    signatures = []

    rel["assets"].each do |aa|
      case aa["name"]
      when "coyim_linux_amd64"
        coyim_entry = aa
      when "coyim_linux_amd64_build_info"
        build_info_entry = aa
      when /^coyim_linux_amd64_build_info\.0x.*?\.rasc$/
        signatures << aa
      end
    end

    if coyim_entry == nil
      puts "Tag #{TAG} doesn't contain coyim binary"
      exit
    end

    if build_info_entry == nil
      puts "Tag #{TAG} doesn't contain a build_info file"
      exit
    end

    if signatures.length == 0
      puts "Tag #{TAG} doesn't contain any build_info signatures"
      exit
    end
    
    $stdout.print "Downloading files "; $stdout.flush
    $stdout.print "."; $stdout.flush
    `curl -L -s #{coyim_entry["browser_download_url"]} -o #{dir}/coyim`
    $stdout.print "."; $stdout.flush
    `curl -L -s #{build_info_entry["browser_download_url"]} -o #{dir}/build_info`
    $stdout.print "."; $stdout.flush
    signatures.each do |entry|
      `curl -L -s #{entry["browser_download_url"]} -o #{dir}/#{entry["name"]}`
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
    signatures.each { |entry|
      puts "Verifying #{entry["name"]}"
      res = system("gpg2 --verify #{dir}/#{entry["name"]} #{dir}/build_info >/dev/null 2>&1")
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
  rescue OpenURI::HTTPError
    puts "Tag #{TAG} doesn't exist"
  end
}
