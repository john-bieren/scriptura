package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// printPassage prints the given passage of book from the Bible.
func printPassage(book, passage string) {
	book_chapters, ok := Bible[book]
	if !ok {
		if book == "Psalm" {
			book = "Psalms"
			book_chapters = Bible["Psalms"]
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
		printChapters(book_chapters, book, "1", "")
		return
	}

	chapters_re := regexp.MustCompile("^([0-9]*)-([0-9]*)$")
	chapter_re := regexp.MustCompile("^([0-9]+)$")
	verses_re := regexp.MustCompile("^([0-9]+):([0-9]*)-([0-9]*)$")
	verse_re := regexp.MustCompile("^([0-9]+):([0-9]+)$")

	if chapters_re.MatchString(passage) {
		matches := chapters_re.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		if matches[1] == matches[2] {
			rangeLengthOneNotice(book, passage)
		}

		printChapters(book_chapters, book, matches[1], matches[2])
	} else if chapter_re.MatchString(passage) {
		matches := chapter_re.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		printChapters(book_chapters, book, matches[1], matches[1])
	} else if verses_re.MatchString(passage) {
		matches := verses_re.FindStringSubmatch(passage)
		errorIfZeroes(matches)
		if matches[2] == matches[3] {
			rangeLengthOneNotice(book, passage)
		}

		chapter_verses, ok := book_chapters[matches[1]]
		if !ok {
			notEnoughChaptersNotice(book_chapters, book, false)
			return
		}
		printVerses(chapter_verses, book, matches[1], matches[2], matches[3])
	} else if verse_re.MatchString(passage) {
		matches := verse_re.FindStringSubmatch(passage)
		errorIfZeroes(matches)

		chapter_verses, ok := book_chapters[matches[1]]
		if !ok {
			notEnoughChaptersNotice(book_chapters, book, false)
			return
		}
		printVerses(chapter_verses, book, matches[1], matches[2], matches[2])
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
	correct_input := strings.TrimSuffix(strings.SplitN(passage, "-", 2)[0], ":")
	if correct_input == "" {
		fmt.Printf("\033[1mNote: \"scriptura %s\" produces the same output\033[0m\n", book)
	} else {
		fmt.Printf("\033[1mNote: \"scriptura %s %s\" produces the same output\033[0m\n", book, correct_input)
	}
}

// printChapters prints the inclusive range (bounded by start and end) of chapters from book_chapters.
// start and end can be empty strings representing the start or end of the book's chapters.
func printChapters(book_chapters map[string]map[string]string, book, start, end string) {
	if start == "" {
		start = "1"
	}

	if end == "" {
		chapter_int, _ := strconv.Atoi(start)
		if chapter_int > len(book_chapters) {
			notEnoughChaptersNotice(book_chapters, book, false)
			return
		}
		var chapter_str string
		var first_newline_skipped bool
		for {
			chapter_str = strconv.Itoa(chapter_int)
			chapter_verses, ok := book_chapters[chapter_str]
			if !ok {
				break
			}

			// skip the newline before the first chapter
			if first_newline_skipped {
				fmt.Print("\n")
			} else {
				first_newline_skipped = true
			}
			if book == "Psalms" {
				fmt.Printf("  \033[1mPsalm %s\033[0m\n", chapter_str)
			} else {
				fmt.Printf("  \033[1mChapter %s\033[0m\n", chapter_str)
			}
			printVerses(chapter_verses, book, chapter_str, "1", "")
			chapter_int++
		}
	} else {
		start_int, _ := strconv.Atoi(start)
		if start_int > len(book_chapters) {
			notEnoughChaptersNotice(book_chapters, book, false)
			return
		}
		end_int, _ := strconv.Atoi(end)
		chapters := generateRange(start_int, end_int)

		for i, chapter := range chapters {
			chapter_verses, ok := book_chapters[chapter]
			if !ok {
				notEnoughChaptersNotice(book_chapters, book, true)
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
			printVerses(chapter_verses, book, chapter, "1", "")
		}
	}
}

// notEnoughChaptersNotice tells the user if their passage references chapters that do not exist in book.
func notEnoughChaptersNotice(book_chapters map[string]map[string]string, book string, bold bool) {
	if bold {
		if len(book_chapters) > 1 {
			if book == "Psalms" {
				fmt.Println("\033[1mThere are only 150 psalms\033[0m")
			} else {
				fmt.Printf("\033[1m%s only has %d chapters\033[0m\n", book, len(book_chapters))
			}
		} else {
			fmt.Printf("\033[1m%s only has 1 chapter\033[0m\n", book)
		}
	} else {
		if len(book_chapters) > 1 {
			if book == "Psalms" {
				fmt.Println("There are only 150 psalms")
			} else {
				fmt.Printf("%s only has %d chapters\n", book, len(book_chapters))
			}
		} else {
			fmt.Printf("%s only has 1 chapter\n", book)
		}
	}
}

// printVerses prints the inclusive range (bounded by start and end) of verses from chapter_verses.
// start and end can be empty strings representing the start or end of the chapter's verses.
func printVerses(chapter_verses map[string]string, book, chapter, start, end string) {
	if start == "" {
		start = "1"
	}

	if end == "" {
		verse_int, _ := strconv.Atoi(start)
		if verse_int > len(chapter_verses) {
			notEnoughVersesNotice(chapter_verses, book, chapter, false)
			return
		}
		var verse_str string
		for {
			verse_str = strconv.Itoa(verse_int)
			verse_text, ok := chapter_verses[verse_str]
			if !ok {
				break
			}
			fmt.Printf("  \033[1m%s\033[0m %s\n", verse_str, verse_text)
			verse_int++
		}
	} else {
		start_int, _ := strconv.Atoi(start)
		if start_int > len(chapter_verses) {
			notEnoughVersesNotice(chapter_verses, book, chapter, false)
			return
		}
		end_int, _ := strconv.Atoi(end)
		verses := generateRange(start_int, end_int)

		for _, verse := range verses {
			verse_text, ok := chapter_verses[verse]
			if !ok {
				notEnoughVersesNotice(chapter_verses, book, chapter, true)
				return
			}

			if len(verses) > 1 {
				fmt.Printf("  \033[1m%s\033[0m ", verse)
			}
			fmt.Println(verse_text)
		}
	}
}

// notEnoughVersesNotice tells the user if their passage references verses that do not exist in chapter of book.
func notEnoughVersesNotice(chapter_verses map[string]string, book, chapter string, bold bool) {
	if bold {
		if book == "Psalms" {
			fmt.Printf("\033[1mPsalm %s only has %d verses\033[0m\n", chapter, len(chapter_verses))
		} else {
			fmt.Printf("\033[1m%s chapter %s only has %d verses\033[0m\n", book, chapter, len(chapter_verses))
		}
	} else {
		if book == "Psalms" {
			fmt.Printf("Psalm %s only has %d verses\n", chapter, len(chapter_verses))
		} else {
			fmt.Printf("%s chapter %s only has %d verses\n", book, chapter, len(chapter_verses))
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
