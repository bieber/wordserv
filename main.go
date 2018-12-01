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
	"github.com/spf13/viper"
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
	viper.SetDefault("port", 80)
	viper.SetDefault("max_words", 5000)
	viper.SetDefault("max_paragraphs", 50)
	viper.SetDefault("max_chapters", 3)

	viper.BindEnv("port")
	viper.BindEnv("max_characters")
	viper.BindEnv("max_paragraphs")
	viper.BindEnv("max_chapters")
	viper.BindEnv("book_dir")

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/run/secrets")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Couldn't load config file: %s", err.Error())
	}

	config := Config{
		Port:          viper.GetInt("port"),
		MaxWords:      viper.GetInt("max_words"),
		MaxParagraphs: viper.GetInt("max_paragraphs"),
		MaxChapters:   viper.GetInt("max_chapters"),
		BookDir:       viper.GetString("book_dir"),
	}

	books, err := loadBooks(config.BookDir)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())
	log.Fatal(startServer(books, &config))
}
