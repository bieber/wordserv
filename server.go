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
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"unicode"
)

type chapterHandler struct {
	maxChapters int
	books       [][][][]string
}

type paragraphHandler struct {
	maxParagraphs int
	books         [][][][]string
}

type wordHandler struct {
	maxWords   int
	books      [][][][]string
	wordCounts [][]int
}

func startServer(books [][][][]string, config *Config) error {
	http.HandleFunc("/", indexHandler)
	http.Handle(
		"/chapters/",
		&chapterHandler{maxChapters: config.MaxChapters, books: books},
	)
	http.Handle(
		"/paragraphs/",
		&paragraphHandler{maxParagraphs: config.MaxParagraphs, books: books},
	)
	http.Handle(
		"/words/",
		&wordHandler{
			maxWords: config.MaxWords,
			books: books,
			wordCounts: countWords(books),
		},
	)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(
		[]byte(
			`<!DOCTYPE html>
<html>
    <head><title>Words API</title></head>
    <body>
        <h1>Words API</h1>
        <p>
            You can use the words API to fetch random snippets of public domain
            English literature.  Just send a GET request to
            <tt>/denomination/count/</tt>, where <tt>count</tt> is the number
            of things that you want and <tt>denomination</tt> is one of
            the following:
        </p>
        <ul>
            <li><tt>chapters</tt></li>
            <li><tt>paragraphs</tt></li>
            <li><tt>words</tt></li>
        </ul>
        <p>
            The response will be JSON text representing one or more levels of
            nested lists.  For <tt>words</tt> requests, the response will be a
            one-layer deep list of words.  For <tt>paragraphs</tt> requests, it
            will be a list of lists of words, where each sub-list is a single
            paragraph.  For <tt>chapters</tt> requests, it will be a list of
            lists of lists, where the outer-most sublists are chapters and the
            inner-most sublists are paragraphs.
        </p>
        <p>
            <strong>
                Note that you should not depend on this server.  It may break
                or disappear at any time and fetches from a limited set of
                works.  If you need to use this API for your own purposes, you
                should host a copy on your own servers.
            </strong>
        </p>
        Source Code:
        <a href="http://www.github.com/bieber/wordserv/">
            http://www.github.com/bieber/wordserv/
        </a>
    </body>
</html>
`,
		),
	)
}

func (c chapterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pattern, _ := regexp.Compile("^/chapters/(\\d+)/?$")
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.NotFound(w, r)
		return
	}

	chapterCount, _ := strconv.Atoi(matches[1])
	book := c.books[rand.Intn(len(c.books))]
	chapterCount = min(chapterCount, c.maxChapters, len(book))

	startChapter := rand.Intn(max(len(book)-chapterCount, 1))
	chapters := book[startChapter : startChapter+chapterCount]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chapters)
}

func (p paragraphHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pattern, _ := regexp.Compile("^/paragraphs/(\\d+)/?$")
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.NotFound(w, r)
		return
	}

	paragraphCount, _ := strconv.Atoi(matches[1])
	paragraphCount = min(paragraphCount, p.maxParagraphs)

	book := p.books[rand.Intn(len(p.books))]
	paragraphsCounted := 0
	maxChapter := 0
	for maxChapter = len(book) - 1; maxChapter >= 0; maxChapter-- {
		paragraphsCounted += len(book[maxChapter])
		if paragraphsCounted >= paragraphCount {
			break
		}
	}
	maxChapter = max(maxChapter, 1)

	paragraphs := [][]string{}
	chapter := rand.Intn(maxChapter)
	paragraph := rand.Intn(len(book[chapter]))
	for len(paragraphs) < paragraphCount {
		if paragraph >= len(book[chapter]) {
			chapter++
			paragraph = 0
		}
		if chapter >= len(book) {
			break
		}

		paragraphsToAdd := min(
			paragraphCount-len(paragraphs),
			len(book[chapter])-paragraph,
		)
		paragraphs = append(
			paragraphs,
			book[chapter][paragraph:paragraph+paragraphsToAdd]...,
		)
		paragraph += paragraphsToAdd
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paragraphs)
}

func (w wordHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	pattern, _ := regexp.Compile("^/words/(\\d+)/?$")
	matches := pattern.FindStringSubmatch(r.URL.Path)
	if matches == nil {
		http.NotFound(rw, r)
		return
	}

	wordCount, _ := strconv.Atoi(matches[1])
	wordCount = min(wordCount, w.maxWords)

	i := rand.Intn(len(w.books))
	book := w.books[i]
	wordCounts := w.wordCounts[i]
	wordsCounted := 0
	maxChapter := 0
	for maxChapter = len(book) - 1; maxChapter >= 0; maxChapter-- {
		wordsCounted += wordCounts[i]
		if wordsCounted >= wordCount {
			break
		}
	}
	maxChapter = max(maxChapter, 1)

	words := []string{}
	wordsCounted = 0
	chapter := rand.Intn(maxChapter)
	paragraph := rand.Intn(len(book[chapter]))
	word := 0
	for wordsCounted < wordCount {
		words = append(words, book[chapter][paragraph][word])
		if unicode.IsLetter([]rune(book[chapter][paragraph][word])[0]) {
			wordsCounted++
		}

		word++
		if word >= len(book[chapter][paragraph]) {
			word = 0
			paragraph++
		}
		if paragraph >= len(book[chapter]) {
			paragraph = 0
			chapter++
		}
		if chapter >= len(book) {
			break
		}
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(words)
}
