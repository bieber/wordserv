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
	"github.com/bieber/conflag"
	"log"
	"math/rand"
	"time"
)

type Config struct {
	Port          int
	MaxCharacters int
	MaxWords      int
	MaxParagraphs int
	MaxChapters   int
	BookDir       string
}

func main() {
	config := Config{
		Port:          80,
		MaxWords:      5000,
		MaxParagraphs: 50,
		MaxChapters:   3,
	}

	confReader, err := conflag.New(&config)
	if err != nil {
		log.Fatal(err)
	}
	confReader.Field("BookDir").
		Required()
	_, err = confReader.Read()
	if err != nil {
		log.Fatal(err)
	}

	books, err := loadBooks(config.BookDir)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	log.Fatal(startServer(books, &config))
}
