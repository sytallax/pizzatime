package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sytallax/pizzatime/dominos"
)

var s *discordgo.Session

var orders = make(map[string]*dominos.Order)

type BotParams struct {
	ApplicationID string
	GuildID       string
	Token         string
}

func loadConfig() *BotParams {
	return &BotParams{
		ApplicationID: os.Getenv("APPLICATION_ID"),
		GuildID:       os.Getenv("GUILD_ID"),
		Token:         os.Getenv("BOT_TOKEN"),
	}
}

func init() {
	c := loadConfig()
	var err error
	s, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []discordgo.ApplicationCommand{
		{
			Name:        "begin-order",
			Description: "Begin a Dominos delivery order",
		},
		{
			Name:        "view-order",
			Description: "View the details of your Dominos order",
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"begin-order": HandleBeginOrder,
		"view-order":  HandleViewOrder,
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			log.Print("Detected a modal submission")
			o, err := ParseBeginOrderResponse(s, i)
			if err != nil {
				log.Printf("Begin Order modal submission handler failed validation: %s", err)
				return
			}
			orders[i.Interaction.Member.User.ID] = o

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your order has been created! To view the full menu, use the '/menu' command",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				panic(err)
			}

		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(loadConfig().ApplicationID, loadConfig().GuildID, &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

}
