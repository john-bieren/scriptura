package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var (
	// bookIsPsalms indicates whether the passage to print is from the Book of Psalms.
	bookIsPsalms bool

	// prettyBookNames converts lowercase, dash-separated book names into capitalized names with spaces.
	prettyBookNames = map[string]string{
		"genesis":         "Genesis",
		"exodus":          "Exodus",
		"leviticus":       "Leviticus",
		"numbers":         "Numbers",
		"deuteronomy":     "Deuteronomy",
		"joshua":          "Joshua",
		"judges":          "Judges",
		"ruth":            "Ruth",
		"1-samuel":        "1 Samuel",
		"2-samuel":        "2 Samuel",
		"1-kings":         "1 Kings",
		"2-kings":         "2 Kings",
		"1-chronicles":    "1 Chronicles",
		"2-chronicles":    "2 Chronicles",
		"ezra":            "Ezra",
		"nehemiah":        "Nehemiah",
		"tobit":           "Tobit",
		"judith":          "Judith",
		"esther":          "Esther",
		"1-maccabees":     "1 Maccabees",
		"2-maccabees":     "2 Maccabees",
		"job":             "Job",
		"psalms":          "Psalms",
		"proverbs":        "Proverbs",
		"ecclesiastes":    "Ecclesiastes",
		"song-of-solomon": "Song of Solomon",
		"wisdom":          "Wisdom",
		"sirach":          "Sirach",
		"isaiah":          "Isaiah",
		"jeremiah":        "Jeremiah",
		"lamentations":    "Lamentations",
		"baruch":          "Baruch",
		"ezekiel":         "Ezekiel",
		"daniel":          "Daniel",
		"hosea":           "Hosea",
		"joel":            "Joel",
		"amos":            "Amos",
		"obadiah":         "Obadiah",
		"jonah":           "Jonah",
		"micah":           "Micah",
		"nahum":           "Nahum",
		"habakkuk":        "Habakkuk",
		"zephaniah":       "Zephaniah",
		"haggai":          "Haggai",
		"zechariah":       "Zechariah",
		"malachi":         "Malachi",
		"matthew":         "Matthew",
		"mark":            "Mark",
		"luke":            "Luke",
		"john":            "John",
		"acts":            "Acts",
		"romans":          "Romans",
		"1-corinthians":   "1 Corinthians",
		"2-corinthians":   "2 Corinthians",
		"galatians":       "Galatians",
		"ephesians":       "Ephesians",
		"philippians":     "Philippians",
		"colossians":      "Colossians",
		"1-thessalonians": "1 Thessalonians",
		"2-thessalonians": "2 Thessalonians",
		"1-timothy":       "1 Timothy",
		"2-timothy":       "2 Timothy",
		"titus":           "Titus",
		"philemon":        "Philemon",
		"hebrews":         "Hebrews",
		"james":           "James",
		"1-peter":         "1 Peter",
		"2-peter":         "2 Peter",
		"1-john":          "1 John",
		"2-john":          "2 John",
		"3-john":          "3 John",
		"jude":            "Jude",
		"revelation":      "Revelation",
	}

	// terminalWindowWidth is the maximum length after which a line must be wrapped during printing.
	terminalWindowWidth int
)

// printPassage prints the given passage of book from the Bible.
func printPassage(book, passage string) {
	var err error
	terminalWindowWidth, _, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		terminalWindowWidth = 532 // length of longest verse (Esther 8:9) printed with its number
	}

	book = strings.ToLower(book)
	bookChapters, ok := Bible[book]
	if !ok {
		if book == "psalm" {
			book = "psalms"
			bookChapters = Bible[book]
		} else {
			if regexp.MustCompile("^[1-3a-z-]+$").MatchString(book) {
				fmt.Printf("Unrecognized book \"%s\": check your spelling and formatting\n", book)
				fmt.Println("Run 'scriptura --books' to see the properly formatted book names")
				os.Exit(1)
			}
			fmt.Println("Invalid arguments")
			usage()
		}
	}
	bookIsPsalms = book == "psalms"

	if passage == "" {
		printChapters(bookChapters, book, "1", "", "", "")
		return
	}

	crossChapterVersesRe := regexp.MustCompile("^([0-9]+):([0-9]+)-([0-9]+):([0-9]+)$")
	chaptersRe := regexp.MustCompile("^([0-9]*)-([0-9]*)$")
	chapterRe := regexp.MustCompile("^([0-9]+)$")
	versesRe := regexp.MustCompile("^([0-9]+):([0-9]*)-([0-9]*)$")
	verseRe := regexp.MustCompile("^([0-9]+):([0-9]+)$")

	if crossChapterVersesRe.MatchString(passage) {
		matches := crossChapterVersesRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		printChapters(bookChapters, book, matches[1], matches[3], matches[2], matches[4])
	} else if chaptersRe.MatchString(passage) {
		matches := chaptersRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		printChapters(bookChapters, book, matches[1], matches[2], "", "")
	} else if chapterRe.MatchString(passage) {
		matches := chapterRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		printChapters(bookChapters, book, matches[1], matches[1], "", "")
	} else if versesRe.MatchString(passage) {
		matches := versesRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		chapterVerses, ok := bookChapters[matches[1]]
		if !ok {
			notEnoughChaptersNotice(bookChapters, book, false)
			return
		}
		printVerses(chapterVerses, book, matches[1], matches[2], matches[3])
	} else if verseRe.MatchString(passage) {
		matches := verseRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		chapterVerses, ok := bookChapters[matches[1]]
		if !ok {
			notEnoughChaptersNotice(bookChapters, book, false)
			return
		}
		printVerses(chapterVerses, book, matches[1], matches[2], matches[2])
	} else {
		fmt.Println("Invalid arguments")
		usage()
	}
}

// errorIfZeroes checks for chapter or verse arguments of zero, and exits the program if any are found.
func errorIfZeroes(matches []string) {
	// skip first match, which is the entire regex
	for _, match := range matches[1:] {
		if match == "0" {
			fmt.Println("Cannot use zero as an argument")
			os.Exit(1)
		}
	}
}

// notEnoughChaptersNotice tells the user that their passage references chapters that do not exist in book.
func notEnoughChaptersNotice(bookChapters map[string]map[string]string, book string, bold bool) {
	// insert escape sequences to start and end bold text, if applicable
	var bs, be string
	if bold {
		bs = "\033[1m"
		be = "\033[0m"
	}

	book = prettyBookNames[book]
	if len(bookChapters) > 1 {
		if bookIsPsalms {
			fmt.Printf("%sThere are only 150 psalms%s\n", bs, be)
		} else {
			fmt.Printf("%s%s only has %d chapters%s\n", bs, book, len(bookChapters), be)
		}
	} else {
		fmt.Printf("%s%s only has 1 chapter%s\n", bs, book, be)
	}
}

// printChapters prints the inclusive range from startChapter:startVerse to endChapter:endVerse from bookChapters.
// Start and end parameters can be empty strings representing the start or end of the relevant text.
func printChapters(
	bookChapters map[string]map[string]string,
	book, startChapter, endChapter, startVerse, endVerse string,
) {
	var startInt int
	if startChapter == "" {
		startInt = 1
	} else {
		startInt, _ = strconv.Atoi(startChapter)
		if startInt > len(bookChapters) {
			notEnoughChaptersNotice(bookChapters, book, false)
			return
		}
	}

	var endInt int
	if endChapter == "" {
		endInt = len(bookChapters)
	} else {
		endInt, _ = strconv.Atoi(endChapter)
		if endInt > len(bookChapters) {
			endInt = len(bookChapters)
			endVerse = ""
			defer notEnoughChaptersNotice(bookChapters, book, true)
		}
	}

	chapters := generateRange(startInt, endInt)
	finalChapterIndex := endInt - startInt
	for i, chapter := range chapters {
		chapterVerses, _ := bookChapters[chapter]

		if (len(chapters) > 1 || endChapter == "") && !(startVerse != "" && i == 0) {
			// add chapter headings
			if i > 0 {
				fmt.Print("\n")
			}
			if bookIsPsalms {
				fmt.Printf("  \033[1mPsalm %s\033[0m\n", chapter)
			} else {
				fmt.Printf("  \033[1mChapter %s\033[0m\n", chapter)
			}
		}

		if finalChapterIndex == 0 {
			printVerses(chapterVerses, book, chapter, startVerse, endVerse)
		} else {
			switch i {
			case 0:
				printVerses(chapterVerses, book, chapter, startVerse, "")
			case finalChapterIndex:
				printVerses(chapterVerses, book, chapter, "1", endVerse)
			default:
				printVerses(chapterVerses, book, chapter, "1", "")
			}
		}
	}
}

// printVerses prints the inclusive range (bounded by start and end) of verses from chapterVerses.
// start and end can be empty strings representing the start or end of the chapter's verses.
func printVerses(chapterVerses map[string]string, book, chapter, start, end string) {
	if start == "" {
		start = "1"
	}
	startInt, _ := strconv.Atoi(start)
	if startInt > len(chapterVerses) {
		notEnoughVersesNotice(chapterVerses, book, chapter, false)
		return
	}

	var endInt int
	if end == "" {
		endInt = len(chapterVerses)
	} else {
		endInt, _ = strconv.Atoi(end)
		if endInt > len(chapterVerses) {
			endInt = len(chapterVerses)
			defer notEnoughVersesNotice(chapterVerses, book, chapter, true)
		}
	}

	verses := generateRange(startInt, endInt)
	for _, verseStr := range verses {
		verseText, _ := chapterVerses[verseStr]

		if len(verses) > 1 || end == "" {
			// add verse number
			wrapPrint(fmt.Sprintf("\033[1m%s\033[0m %s", verseStr, verseText), 2, 8)
		} else {
			wrapPrint(verseText, 0, 0)
		}
	}
}

// notEnoughVersesNotice tells the user that their passage references verses that do not exist in chapter of book.
func notEnoughVersesNotice(chapterVerses map[string]string, book, chapter string, bold bool) {
	// insert escape sequences to start and end bold text, if applicable
	var bs, be string
	if bold {
		bs = "\033[1m"
		be = "\033[0m"
	}

	book = prettyBookNames[book]
	if bookIsPsalms {
		fmt.Printf("%sPsalm %s only has %d verses%s\n", bs, chapter, len(chapterVerses), be)
	} else {
		fmt.Printf("%s%s chapter %s only has %d verses%s\n", bs, book, chapter, len(chapterVerses), be)
	}
}

// wrapPrint prints text with word wrapping based on terminalWindowWidth.
func wrapPrint(text string, leadingSpaces, falseLength int) {
	words := strings.Fields(text)
	var wrappedText strings.Builder
	wrappedText.WriteString(strings.Repeat(" ", leadingSpaces))
	lineLength := leadingSpaces - falseLength

	for _, word := range words {
		if lineLength+len(word)+1 > terminalWindowWidth {
			wrappedText.WriteString("\n")
			lineLength = 0
		}
		wrappedText.WriteString(word)
		wrappedText.WriteString(" ")
		lineLength += len(word) + 1
	}
	fmt.Println(strings.TrimRight(wrappedText.String(), " "))
}

// generateRange returns a slice of the string representations of the inclusive range from start to end.
func generateRange(start, end int) []string {
	size := end - start + 1
	if size < 1 {
		fmt.Println("Invalid range: start and end are reversed")
		usage()
	}

	// large numbers cause a panic; 176 is the maximum plausible value (Psalm 118 has 176 verses)
	if size > 176 {
		size = 176
	}
	result := make([]string, size)
	for i := range size {
		result[i] = strconv.Itoa(start + i)
	}
	return result
}
