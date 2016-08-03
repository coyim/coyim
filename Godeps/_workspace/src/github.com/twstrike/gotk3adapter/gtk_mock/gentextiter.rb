#!/usr/bin/env ruby

all = [
	"BackwardChar() bool",
	"BackwardChars(int) bool",
    "BackwardCursorPosition() bool",
	"BackwardCursorPositions(int) bool",
    "BackwardLine() bool",
	"BackwardLines(int) bool",
    "BackwardToTagToggle(TextTag) bool",
	"BackwardVisibleCursorPosition() bool",
    "BackwardVisibleCursorPositions(int) bool",
	"BackwardVisibleLine() bool",
    "BackwardVisibleLines(int) bool",
	"BeginsTag(TextTag) bool",
    "CanInsert(bool) bool",
	"Compare(TextIter) int",
    "Editable(bool) bool",
	"EndsLine() bool",
    "EndsSentence() bool",
	"EndsTag(TextTag) bool",
    "EndsWord() bool",
	"Equal(TextIter) bool",
    "ForwardChar() bool",
	"ForwardChars(int) bool",
    "ForwardCursorPosition() bool",
	"ForwardCursorPositions(int) bool",
    "ForwardLine() bool",
	"ForwardLines(int) bool",
    "ForwardSentenceEnd() bool",
	"ForwardSentenceEnds(int) bool",
    "ForwardToEnd()",
	"ForwardToLineEnd() bool",
    "ForwardToTagToggle(TextTag) bool",
	"ForwardVisibleCursorPosition() bool",
    "ForwardVisibleCursorPositions(int) bool",
	"ForwardVisibleLine() bool",
    "ForwardVisibleLines(int) bool",
	"ForwardVisibleWordEnd() bool",
    "ForwardVisibleWordEnds(int) bool",
	"ForwardWordEnd() bool",
    "ForwardWordEnds(int) bool",
	"GetBuffer() TextBuffer",
    "GetBytesInLine() int",
	"GetChar() rune",
    "GetCharsInLine() int",
	"GetLine() int",
    "GetLineIndex() int",
	"GetLineOffset() int",
    "GetOffset() int",
	"GetSlice(TextIter) string",
    "GetText(TextIter) string",
	"GetVisibleLineIndex() int",
    "GetVisibleLineOffset() int",
	"GetVisibleSlice(TextIter) string",
    "GetVisibleText(TextIter) string",
	"HasTag(TextTag) bool",
    "InRange(TextIter, TextIter) bool",
	"InsideSentence() bool",
    "InsideWord() bool",
	"IsCursorPosition() bool",
    "IsEnd() bool",
	"IsStart() bool",
    "SetLine(int)",
	"SetLineIndex(int)",
    "SetLineOffset(int)",
	"SetOffset(int)",
    "SetVisibleLineIndex(int)",
	"SetVisibleLineOffset(int)",
    "StartsLine() bool",
	"StartsSentence() bool",
    "StartsWord() bool",
	"TogglesTag(TextTag) bool"
]

$PRIMITIVES = {
  "bool" => "false",
  "int" => "0",
  "string" => "\"\"",
  "rune" => "0",
}

def parse(s)
  name, args, rets = /^(.*?)\((.*?)\)(?: ?(.*?))$/.match(s).captures
  {
    name: name,
    args: args.split(", "),
    rets: rets
  }
end

def mapType(tt)
  if $PRIMITIVES[tt]
    tt
  else
    "gtki.#{tt}"
  end
end

def argList(args)
  args.map{ |x|
    "#{mapType(x)}"
  }.join(", ")
end

def returnList(rets)
  if rets == ""
    ""
  else
    " #{mapType(rets)}"
  end
end

def potentialReturn(rets)
  if rets == ""
    ""
  else
    if $PRIMITIVES[rets]
      "return #{$PRIMITIVES[rets]}"
    else
      "return nil"
    end
  end
end


all.each do |xx|
  res = parse(xx)
  puts <<END
func (*MockTextIter) #{res[:name]}(#{argList(res[:args])})#{returnList(res[:rets])} {
	#{potentialReturn(res[:rets])}
}

END
end
