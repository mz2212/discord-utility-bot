package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krognol/go-wolfram"
	"github.com/bwmarrin/discordgo"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/spf13/viper"
	"github.com/turnage/graw/reddit"
)

var (
	wolf   wolfram.Client
	redd   reddit.Bot
	malCli *mal.Client
)

func main() {
	var err error

	viper.SetConfigName("bot")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("[Config] Error reading config: ", err)
		return
	}

	fmt.Println("[Reddit] Init...")
	app := reddit.App{
		ID:       viper.GetString("reddit.id"),
		Secret:   viper.GetString("reddit.secret"),
		Username: viper.GetString("reddit.username"), // Why graw needs this, I don't know
		Password: viper.GetString("reddit.password"), // redd does just fine without...
	}
	botConfig := reddit.BotConfig{
		Agent: viper.GetString("reddit.agent"),
		App:   app,
		Rate:  1 * time.Second,
	}
	redd, err = reddit.NewBot(botConfig)
	if err != nil {
		fmt.Println("[Reddit] Failed to get session: ", err)
		return
	}

	fmt.Println("[Wolfram|Alpha] Init...") // Wolfram|Alpha's API appears to be stateless
	wolf = wolfram.Client{AppID: viper.GetString("wolfram.key")}

	fmt.Println("[MAL] Init...")
	malCli = mal.NewClient(nil)
	malCli.SetCredentials(viper.GetString("mal.username"), viper.GetString("mal.password"))

	fmt.Println("[Discord] Init...") // Init discord last, because reasons.
	discord, err := discordgo.New("Bot " + viper.GetString("discord.key"))
	if err != nil {
		fmt.Println("[Discord] Error getting session: ", err)
		return
	}
	discord.AddHandler(messageCreate)
	if err := discord.Open(); err != nil {
		fmt.Println("[Discord] Error opening session: ", err)
		return
	}

	fmt.Println("[Main] Ready.")
	fmt.Println("[Main] Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc // This waits for somthing to come in on the "sc" channel.
	fmt.Println("[Main] Ctrl-C Recieved. Exiting!")
	discord.Close()
}
