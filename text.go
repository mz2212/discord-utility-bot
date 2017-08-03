package main

import "fmt"

// This source file contains the really long strings
// avalible for return in functions.
// To keep it out of the source to keep it clean.

func help() string {
	help := "```" + `
#!usergen [redditor]
  Generates a comment based on a user's previous comments
#!subgen [subreddit]
  Generates a comment based on a sub's previous comments
#!ask [question]
  Asks Wolfram|Alpha a question (Short answer)
#!ask -image [question]
  Asks Wolfram|Alpha a question (Image based answer)
#!anime [query]
	Searches MyAnimeList for your query
	and returns the first result
#!manga [query]
	Searches MyAnimeList for your query
	and returns the first result` + "```"
	return help
}

func errorText() string {
	return "Somthing went wrong...\nCheck the log for details."
}

func animeFormat(title, start, end string, score float64, episodes int) string {
	text := fmt.Sprintf("%s\nStart Date: %s\nEnd Date: %s\nScore: %2.2f\nEpisode Count: %d", title, start, end, score, episodes)
	return text
}

func mangaFormat(title, start, end string, score float64, chapters, volumes int) string {
	text := fmt.Sprintf("%s\nStart Date: %s\nEnd Date: %s\nScore: %2.2f\n%d chapters in %d volume(s)", title, start, end, score, chapters, volumes)
	return text
}
