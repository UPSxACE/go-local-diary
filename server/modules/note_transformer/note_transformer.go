package note_transformer

import (
	"html/template"
	"math"
	"strings"
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

func ParseToHtml(content string) template.HTML {
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
					var replacedStr = lineStartMatch.replaceStart + line[length:] + lineStartMatch.replaceEnd
					finalLine = replacedStr;
					// Found line start match, so stop iterating
					foundMatch = true;
					break;
				} 
			}
		} 
		// In case no line start match was found, just make it a normal <p> block
		if(!foundMatch){
			finalLine = "<p>" + line + "</p>"
		}

		// Apply pair matches
		for _, pairMatch := range pairMatches {
			var countStart = strings.Count(finalLine, pairMatch.matchStart)
			var equalMatchers = pairMatch.matchStart == pairMatch.matchEnd
			if equalMatchers {
				var validMatches int = countStart / 2
				for validMatches > 0 {
					finalLine = strings.Replace(finalLine, pairMatch.matchStart, pairMatch.replaceStart, 1)
					finalLine = strings.Replace(finalLine, pairMatch.matchEnd, pairMatch.replaceEnd, 1)
					validMatches--
				}
			}
			if !equalMatchers {
				var countEnd = strings.Count(finalLine, pairMatch.matchEnd)
				var validMatches int = int(math.Min(float64(countStart), float64(countEnd)))
				for validMatches > 0 {
					finalLine = strings.Replace(finalLine, pairMatch.matchStart, pairMatch.replaceStart, 1)
					finalLine = strings.Replace(finalLine, pairMatch.matchEnd, pairMatch.replaceEnd, 1)
					validMatches--
				}
			}

		}

		// Add parsed line to the result string
		result += finalLine
	}

	return template.HTML(result)
}
