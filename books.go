/* Copyright (c) 2014 Robert Bieber
 *
 * This file is part of wordserv.
 *
 * wordserv is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

func loadBooks(basePath string) ([][][][]string, error) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	books := [][][][]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := path.Join(basePath, file.Name())
		fin, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer fin.Close()

		books = append(books, readBook(fin))
	}

	_ = books
	return books, nil
}

func readBook(fin *os.File) [][][]string {
	chapters := [][][]string{}
	paragraphs := [][]string{}
	words := []string{}
	lastLineBlank := false

	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			lastLineBlank = true
			continue
		} else if lastLineBlank {
			if len(words) != 0 {
				paragraphs = append(paragraphs, words)
				words = []string{}
			}
			lastLineBlank = false
		}

		if line == "[CHAPTER]" {
			if len(paragraphs) == 0 {
				continue
			}
			chapters = append(chapters, paragraphs)
			paragraphs = [][]string{}
			words = []string{}
		} else {
			chars := []rune(line)
			word := []rune{}
			for _, r := range chars {
				if unicode.IsLetter(r) || (len(word) != 0 && r == '\'') {
					word = append(word, r)
				} else {
					if len(word) != 0 {
						words = append(words, string(word))
						word = []rune{}
					}

					if !unicode.IsSpace(r) {
						words = append(words, string([]rune{r}))
					}
				}
			}
			if len(word) != 0 {
				words = append(words, string(word))
				word = []rune{}
			}
		}
	}
	if len(words) > 0 {
		paragraphs = append(paragraphs, words)
	}
	if len(paragraphs) > 0 {
		chapters = append(chapters, paragraphs)
	}
	return chapters
}
