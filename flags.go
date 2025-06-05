package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// books_flag is the value of the --books flag.
	books_flag = flag.Bool("books", false, "")

	// license_flag is the value of the --license flag.
	license_flag = flag.Bool("license", false, "")

	// version_flag is the value of the --version flag.
	version_flag = flag.Bool("version", false, "")

	// books_message is the line-by-line output of the --books flag.
	books_message = []string{
		"Listed below are the 73 books of the Bible",
		"When using scriptura, refer to books exactly as written below",
		"*Note capitalization and use of dashes*",
		"Genesis",
		"Exodus",
		"Leviticus",
		"Numbers",
		"Deuteronomy",
		"Joshua",
		"Judges",
		"Ruth",
		"1-Samuel",
		"2-Samuel",
		"1-Kings",
		"2-Kings",
		"1-Chronicles",
		"2-Chronicles",
		"Ezra",
		"Nehemiah",
		"Tobit",
		"Judith",
		"Esther",
		"1-Maccabees",
		"2-Maccabees",
		"Job",
		"Psalms or Psalm",
		"Proverbs",
		"Ecclesiastes",
		"Song-of-Solomon",
		"Wisdom",
		"Sirach",
		"Isaiah",
		"Jeremiah",
		"Lamentations",
		"Baruch",
		"Ezekiel",
		"Daniel",
		"Hosea",
		"Joel",
		"Amos",
		"Obadiah",
		"Jonah",
		"Micah",
		"Nahum",
		"Habakkuk",
		"Zephaniah",
		"Haggai",
		"Zechariah",
		"Malachi",
		"Matthew",
		"Mark",
		"Luke",
		"John",
		"Acts",
		"Romans",
		"1-Corinthians",
		"2-Corinthians",
		"Galatians",
		"Ephesians",
		"Philippians",
		"Colossians",
		"1-Thessalonians",
		"2-Thessalonians",
		"1-Timothy",
		"2-Timothy",
		"Titus",
		"Philemon",
		"Hebrews",
		"James",
		"1-Peter",
		"2-Peter",
		"1-John",
		"2-John",
		"3-John",
		"Jude",
		"Revelation",
	}

	// license_message is the line-by-line output of the --license flag.
	license_message = []string{
		"scriptura uses the Douay-Rheims 1899 American Edition (DRA) version of the Bible",
		"The text was sourced from eBible.org: ebible.org/Scriptures/engDRA_readaloud.zip",
		"",
		"scriptura's source code can be found at github.com/john-bieren/scriptura",
		"This program is licensed under the MIT license:",
		"",
		"Copyright (c) 2025 John Bieren",
		"",
		"Permission is hereby granted, free of charge, to any person obtaining a copy",
		"of this software and associated documentation files (the \"Software\"), to deal",
		"in the Software without restriction, including without limitation the rights",
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell",
		"copies of the Software, and to permit persons to whom the Software is",
		"furnished to do so, subject to the following conditions:",
		"",
		"The above copyright notice and this permission notice shall be included in all",
		"copies or substantial portions of the Software.",
		"",
		"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR",
		"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,",
		"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE",
		"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER",
		"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,",
		"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE",
		"SOFTWARE.",
	}

	// usage_message is the line-by-line output of usage().
	usage_message = []string{
		fmt.Sprintf("scriptura %s", version),
		"Read passages from the Bible by book, chapter, or verse",
		"",
		"Usage: scriptura <book>",
		"       scriptura <book> [start_chapter]-[end_chapter]",
		"       scriptura <book> <chapter>",
		"       scriptura <book> <chapter>:[start_verse]-[end_verse]",
		"       scriptura <book> <chapter>:<verse>",
		"",
		"Options:",
		"       --books    Print the list of books and exit",
		"       --help     Print this message and exit",
		"       --license  Print license, citation information and exit",
		"       --version  Print version and exit",
	}
)

// processExitFlags runs exit flags.
func processExitFlags() {
	if *books_flag {
		for _, line := range books_message {
			fmt.Println(line)
		}
		os.Exit(0)
	}

	if *license_flag {
		for _, line := range license_message {
			fmt.Println(line)
		}
		os.Exit(0)
	}

	if *version_flag {
		fmt.Println("scriptura", version)
		os.Exit(0)
	}
}

// usage prints usage_message for the --help flag and relevant error messages.
func usage() {
	for _, line := range usage_message {
		fmt.Println(line)
	}
	os.Exit(0)
}
