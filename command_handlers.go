package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sytallax/pizzatime/dominos"
)

func HandleBeginOrder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if _, ok := orders[i.Interaction.Member.User.ID]; ok {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "It looks like you already have an order in-progress.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			panic(err)
		}
		return
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "begin_order_" + i.Interaction.Member.User.ID,
			Title:    "Begin Dominos Order",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "name",
							Label:       "What name do you want the delivery under?",
							Style:       discordgo.TextInputShort,
							Placeholder: "John Doe",
							Required:    true,
							MaxLength:   45,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "address-street",
							Label:       "Enter delivery address",
							Style:       discordgo.TextInputShort,
							Placeholder: "Number, street name, and unit info only",
							Required:    true,
							MaxLength:   45,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "address-city",
							Label:       "Enter the city of your address",
							Style:       discordgo.TextInputShort,
							Placeholder: "City",
							Required:    true,
							MaxLength:   45,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "address-region",
							Label:       "Enter the region of your address",
							Style:       discordgo.TextInputShort,
							Placeholder: "Region",
							Required:    true,
							MaxLength:   45,
							MinLength:   2,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "address-postal-code",
							Label:       "Enter the postal code of your address",
							Style:       discordgo.TextInputShort,
							Placeholder: "Postal Code",
							Required:    true,
							MaxLength:   45,
							MinLength:   2,
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
func ParseBeginOrderResponse(s *discordgo.Session, i *discordgo.InteractionCreate) (*dominos.Order, error) {
	if i.Type != discordgo.InteractionModalSubmit {
		return nil, errors.New("Did not recieve a modal submission")
	}

	data := i.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, "begin_order") {
		msg := "Did not recieve a Begin Order modal submission"
		log.Print(msg)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			panic(err)
		}
		return nil, errors.New(msg)
	}

	pc, err := strconv.Atoi(data.Components[4].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
	if err != nil {
		return nil, err
	}

	o := &dominos.Order{
		Customer: dominos.Customer{
			Name: data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
			Address: dominos.Address{
				Street:     data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				City:       data.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				Region:     data.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
				PostalCode: pc,
			},
		},
	}

	return o, nil
}

func HandleViewOrder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if _, ok := orders[i.Interaction.Member.User.ID]; !ok {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have an active order",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			panic(err)
		}
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You have an active order, but I can't show it to you yet",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}
