package models

import (
	"html/template"
	"math"
	"strings"
	"unicode/utf8"
)

type match struct {
	matchStart   string
	matchEnd     string
	replaceStart string
	replaceEnd   string
}

var lineStartMatches = []match{
	{
		matchStart:   "### ",
		replaceStart: "<h3>",
		replaceEnd:   "</h3>",
	},
	{
		matchStart:   "## ",
		replaceStart: "<h2>",
		replaceEnd:   "</h2>",
	},
	{
		matchStart:   "# ",
		replaceStart: "<h1>",
		replaceEnd:   "</h1>",
	},
}
var pairMatches = []match{
	{
		matchStart:   "**",
		matchEnd:     "**",
		replaceStart: "<strong>",
		replaceEnd:   "</strong>",
	},
	{
		matchStart:   "*",
		matchEnd:     "*",
		replaceStart: "<em>",
		replaceEnd:   "</em>",
	},
}

// parse strategies
const PARSE_TO_HTML = "HTML";
const PARSE_TO_RAW = "RAW"

func parse(content string, parseStrategy string) string {
	// From dangerous to safe HTML
	var safeContent = template.HTMLEscapeString(content)

	// Variable to store HTML after parsed
	var result string
	lines := strings.Split(safeContent, "\n")
	for _, line := range lines {
		// Variable to store the line while parsed
		var finalLine string
		var foundMatch bool;
		// Apply line starter matches, and <p>'s
		for _, lineStartMatch := range lineStartMatches {
			var length = len(lineStartMatch.matchStart)
			if len(line) >= length {
				var startOfLine = line[0:length]
				if startOfLine == lineStartMatch.matchStart {
					var replacedStr string;
					if(parseStrategy == PARSE_TO_HTML){
						replacedStr = lineStartMatch.replaceStart + line[length:] + lineStartMatch.replaceEnd
					} else {
						replacedStr = line[length:]
					}
					finalLine = replacedStr;
					// Found line start match, so stop iterating
					foundMatch = true;
					break;
				} 
			}
		} 
		// In case no line start match was found, just make it a normal <p> block
		if(!foundMatch){
			if(parseStrategy == PARSE_TO_HTML){
				finalLine = "<p>" + line + "</p>"
			} else {
				finalLine = line
			}
		}

		// Apply pair matches
		for _, pairMatch := range pairMatches {
			var countStart = strings.Count(finalLine, pairMatch.matchStart)
			var equalMatchers = pairMatch.matchStart == pairMatch.matchEnd
			if equalMatchers {
				var validMatches int = countStart / 2
				for validMatches > 0 {
					if(parseStrategy == PARSE_TO_HTML){
						finalLine = strings.Replace(finalLine, pairMatch.matchStart, pairMatch.replaceStart, 1)
						finalLine = strings.Replace(finalLine, pairMatch.matchEnd, pairMatch.replaceEnd, 1)
					} else {
						finalLine = strings.Replace(finalLine, pairMatch.matchStart, "", 1)
						finalLine = strings.Replace(finalLine, pairMatch.matchEnd, "", 1)
					}
					
					validMatches--
				}
			}
			if !equalMatchers {
				var countEnd = strings.Count(finalLine, pairMatch.matchEnd)
				var validMatches int = int(math.Min(float64(countStart), float64(countEnd)))
				for validMatches > 0 {
					if(parseStrategy == PARSE_TO_HTML){
						finalLine = strings.Replace(finalLine, pairMatch.matchStart, pairMatch.replaceStart, 1)
						finalLine = strings.Replace(finalLine, pairMatch.matchEnd, pairMatch.replaceEnd, 1)
					} else {
						finalLine = strings.Replace(finalLine, pairMatch.matchStart, "", 1)
						finalLine = strings.Replace(finalLine, pairMatch.matchEnd, "", 1)
					}
					
					validMatches--
				}
			}

		}

		// Add parsed line to the result string
		result += finalLine
	}

	return result
}


func ParseNoteContentToHTML(content string) template.HTML {
	result := parse(content, PARSE_TO_HTML);
	return template.HTML(result)
}

func ParseNoteContentToRaw(content string) string {
	result := parse(content, PARSE_TO_RAW);
	return result
}


type NoteModelPreview struct {
	NoteModel
}

func NewNotePreviewModel(model NoteModel,initialByteArrIndex int, runeAmount int) NoteModelPreview{
	var contentPreview []byte;
	if(initialByteArrIndex != 0){
		contentPreview = append(contentPreview, []byte("...")...) ;
	}
	

	contentLen := len(model.ContentRaw)

	iteratedRuneCount := -1;
	currentByteIteration := 0;

	for iteratedRuneCount < runeAmount && currentByteIteration < contentLen {
		currentRune := model.ContentRaw[currentByteIteration]

		isNewRune := utf8.RuneStart(currentRune)
		if(isNewRune){
			iteratedRuneCount++;
			if(iteratedRuneCount >= runeAmount){
				break;
			}
		}
		contentPreview = append(contentPreview, currentRune)

		currentByteIteration++;
	}

	// add <highlight> tag around searched word

	// compare contentLen to currentByteIteration to see if iterated through all; if not add "..."
	if(currentByteIteration < contentLen){
		contentPreview = append(contentPreview, []byte("...")...) ;
	}

	modelPreview := NoteModelPreview{NoteModel: model}
	modelPreview.ContentRaw = string(contentPreview);

	return modelPreview;
}