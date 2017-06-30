package main

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
  Asks Wolfram|Alpha a question (Image based answer)` + "```"
	return help
}
