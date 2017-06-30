package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Krognol/go-wolfram"
	"github.com/bwmarrin/discordgo"
	"github.com/mz2212/discord-utility-bot/markov"
)

func messageCreate(client *discordgo.Session, message *discordgo.MessageCreate) {
	splitMessage := strings.Split(message.Content, " ")

	if message.Author.ID == client.State.User.ID {
		return // Ignore messages from self
	}

	// Reddit Stuff
	// '/r/' helper.
	matcher := regexp.MustCompile(`^[(reddit\.com)]?/?r/([^\s]+)`)
	if matcher.MatchString(message.Content) {
		matches := matcher.FindStringSubmatch(message.Content)
		//fmt.Println("Matches: ", matches) // Debugging.
		client.ChannelMessageSend(message.ChannelID, fmt.Sprint("Link: https://reddit.com/r/", matches[1]))
		return
	}
	// User/Sub quote generator.
	if splitMessage[0] == "#!usergen" {
		client.ChannelMessageSend(message.ChannelID, gen("/u/"+splitMessage[1]))
		return
	}
	if splitMessage[0] == "#!subgen" {
		client.ChannelMessageSend(message.ChannelID, gen("/r/"+splitMessage[1]+"/comments"))
		return
	}

	// In chat documentation
	if splitMessage[0] == "#!help" {
		channel, err := client.UserChannelCreate(message.Author.ID)
		if err != nil {
			fmt.Println("[Discord] Failed to create DM: ", err)
			return
		}
		client.ChannelMessageSend(channel.ID, help())
	}

	// Wolfram Stuff
	if splitMessage[0] == "#!ask" {
		if splitMessage[1] == "-image" {
			// Question handler using the "Simple API"
			query := strings.Join(splitMessage[2:], " ")
			img, _, err := wolf.GetSimpleQuery(query, url.Values{})
			if err != nil {
				client.ChannelMessageSend(message.ChannelID, "Somthing went wrong... \nCheck console for more info.")
				fmt.Println("[Wolfram|Alpha] Error: ", err)
				return
			}
			client.ChannelFileSend(message.ChannelID, "simple.gif", img)
			return // I guess the reader returned is a gif, so I just upload it to discord
		}
		// Question handler using the "Short Answer API"
		query := strings.Join(splitMessage[1:], " ")
		answer, err := wolf.GetShortAnswerQuery(query, wolfram.Metric, 30)
		if err != nil { // For some reason the area of a football field is still returned in feet
			client.ChannelMessageSend(message.ChannelID, "Somthing went wrong... \nCheck the console for more info.")
			fmt.Println("[Wolfram|Alpha] Error: ", err)
			return
		}
		client.ChannelMessageSend(message.ChannelID, answer)
		return
	}
}

// Helper Functions
func gen(loc string) string {
	harvest, err := redd.Listing(loc, "")
	if err != nil {
		fmt.Println("[Reddit] Failed to get listing for ", loc, ": ", err)
		return fmt.Sprint("Failed to get listing for ", loc, "\nEither that location doesn't exist, or I bugged out...")
	}
	gen := markov.New(2)
	for _, comment := range harvest.Comments[:30] {
		gen.Build(comment.Body)
	}
	locSplit := strings.Split(loc, "/")
	return fmt.Sprint(gen.Generate(100), " - /", locSplit[1], "/", locSplit[2])
}
