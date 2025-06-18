package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// terminalWindowWidth is the maximum length after which a line must be wrapped during printing.
var terminalWindowWidth int

// printPassage prints the given passage of book from the Bible.
func printPassage(book, passage string) {
	var err error
	terminalWindowWidth, _, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		terminalWindowWidth = 535
	}

	bookChapters, ok := Bible[book]
	if !ok {
		if book == "Psalm" {
			book = "Psalms"
			bookChapters = Bible["Psalms"]
		} else if regexp.MustCompile("^[-1-3A-Za-z]+$").MatchString(book) {
			fmt.Printf("Unrecognized book \"%s\": check your capitalization, spelling, and formatting\n", book)
			fmt.Println("Run 'scriptura --books' to see the properly formatted book names")
			os.Exit(1)
		} else {
			fmt.Println("Invalid arguments")
			usage()
		}
	}

	if passage == "" {
		printChapters(bookChapters, book, "1", "")
		return
	}

	chaptersRe := regexp.MustCompile("^([0-9]*)-([0-9]*)$")
	chapterRe := regexp.MustCompile("^([0-9]+)$")
	versesRe := regexp.MustCompile("^([0-9]+):([0-9]*)-([0-9]*)$")
	verseRe := regexp.MustCompile("^([0-9]+):([0-9]+)$")

	if chaptersRe.MatchString(passage) {
		matches := chaptersRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		if matches[1] == matches[2] {
			rangeLengthOneNotice(book, passage)
		}

		printChapters(bookChapters, book, matches[1], matches[2])
	} else if chapterRe.MatchString(passage) {
		matches := chapterRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		printChapters(bookChapters, book, matches[1], matches[1])
	} else if versesRe.MatchString(passage) {
		matches := versesRe.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		if matches[2] == matches[3] {
			rangeLengthOneNotice(book, passage)
		}

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
	for _, match := range matches[1:] {
		if match == "0" {
			fmt.Println("Cannot use zero as an argument")
			os.Exit(1)
		}
	}
}

// rangeLengthOneNotice gives the user feedback if their range has a length of one.
func rangeLengthOneNotice(book, passage string) {
	correctInput := strings.TrimSuffix(strings.SplitN(passage, "-", 2)[0], ":")
	if correctInput == "" {
		fmt.Printf("\033[1mNote: \"scriptura %s\" produces the same output\033[0m\n", book)
	} else {
		fmt.Printf("\033[1mNote: \"scriptura %s %s\" produces the same output\033[0m\n", book, correctInput)
	}
}

// printChapters prints the inclusive range (bounded by start and end) of chapters from bookChapters.
// start and end can be empty strings representing the start or end of the book's chapters.
func printChapters(bookChapters map[string]map[string]string, book, start, end string) {
	if start == "" {
		start = "1"
	}

	if end == "" {
		chapterInt, _ := strconv.Atoi(start)
		if chapterInt > len(bookChapters) {
			notEnoughChaptersNotice(bookChapters, book, false)
			return
		}
		var chapterStr string
		var firstNewlineSkipped bool
		for {
			chapterStr = strconv.Itoa(chapterInt)
			chapterVerses, ok := bookChapters[chapterStr]
			if !ok {
				break
			}

			// skip the newline before the first chapter
			if firstNewlineSkipped {
				fmt.Print("\n")
			} else {
				firstNewlineSkipped = true
			}
			if book == "Psalms" {
				fmt.Printf("  \033[1mPsalm %s\033[0m\n", chapterStr)
			} else {
				fmt.Printf("  \033[1mChapter %s\033[0m\n", chapterStr)
			}
			printVerses(chapterVerses, book, chapterStr, "1", "")
			chapterInt++
		}
	} else {
		startInt, _ := strconv.Atoi(start)
		if startInt > len(bookChapters) {
			notEnoughChaptersNotice(bookChapters, book, false)
			return
		}
		endInt, _ := strconv.Atoi(end)
		chapters := generateRange(startInt, endInt)

		for i, chapter := range chapters {
			chapterVerses, ok := bookChapters[chapter]
			if !ok {
				notEnoughChaptersNotice(bookChapters, book, true)
				return
			}

			if len(chapters) > 1 {
				if i > 0 {
					fmt.Print("\n")
				}
				if book == "Psalms" {
					fmt.Printf("  \033[1mPsalm %s\033[0m\n", chapter)
				} else {
					fmt.Printf("  \033[1mChapter %s\033[0m\n", chapter)
				}
			}
			printVerses(chapterVerses, book, chapter, "1", "")
		}
	}
}

// notEnoughChaptersNotice tells the user if their passage references chapters that do not exist in book.
func notEnoughChaptersNotice(bookChapters map[string]map[string]string, book string, bold bool) {
	if bold {
		if len(bookChapters) > 1 {
			if book == "Psalms" {
				fmt.Println("\033[1mThere are only 150 psalms\033[0m")
			} else {
				fmt.Printf("\033[1m%s only has %d chapters\033[0m\n", book, len(bookChapters))
			}
		} else {
			fmt.Printf("\033[1m%s only has 1 chapter\033[0m\n", book)
		}
	} else {
		if len(bookChapters) > 1 {
			if book == "Psalms" {
				fmt.Println("There are only 150 psalms")
			} else {
				fmt.Printf("%s only has %d chapters\n", book, len(bookChapters))
			}
		} else {
			fmt.Printf("%s only has 1 chapter\n", book)
		}
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

// printVerses prints the inclusive range (bounded by start and end) of verses from chapterVerses.
// start and end can be empty strings representing the start or end of the chapter's verses.
func printVerses(chapterVerses map[string]string, book, chapter, start, end string) {
	if start == "" {
		start = "1"
	}

	if end == "" {
		verseInt, _ := strconv.Atoi(start)
		if verseInt > len(chapterVerses) {
			notEnoughVersesNotice(chapterVerses, book, chapter, false)
			return
		}

		for {
			verseStr := strconv.Itoa(verseInt)
			verseText, ok := chapterVerses[verseStr]
			if !ok {
				break
			}
			wrapPrint(fmt.Sprintf("\033[1m%s\033[0m %s", verseStr, verseText), 2, 8)
			verseInt++
		}
	} else {
		startInt, _ := strconv.Atoi(start)
		if startInt > len(chapterVerses) {
			notEnoughVersesNotice(chapterVerses, book, chapter, false)
			return
		}
		endInt, _ := strconv.Atoi(end)
		verses := generateRange(startInt, endInt)

		for _, verseStr := range verses {
			verseText, ok := chapterVerses[verseStr]
			if !ok {
				notEnoughVersesNotice(chapterVerses, book, chapter, true)
				return
			}

			if len(verses) > 1 {
				wrapPrint(fmt.Sprintf("\033[1m%s\033[0m %s", verseStr, verseText), 2, 8)
			} else {
				wrapPrint(verseText, 0, 0)
			}
		}
	}
}

// notEnoughVersesNotice tells the user if their passage references verses that do not exist in chapter of book.
func notEnoughVersesNotice(chapterVerses map[string]string, book, chapter string, bold bool) {
	if bold {
		if book == "Psalms" {
			fmt.Printf("\033[1mPsalm %s only has %d verses\033[0m\n", chapter, len(chapterVerses))
		} else {
			fmt.Printf("\033[1m%s chapter %s only has %d verses\033[0m\n", book, chapter, len(chapterVerses))
		}
	} else {
		if book == "Psalms" {
			fmt.Printf("Psalm %s only has %d verses\n", chapter, len(chapterVerses))
		} else {
			fmt.Printf("%s chapter %s only has %d verses\n", book, chapter, len(chapterVerses))
		}
	}
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
