package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DomNidy/saint_sim/apps/discord_bot/constants"
	"github.com/DomNidy/saint_sim/pkg/interfaces"
	"github.com/DomNidy/saint_sim/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

// Commands
//	 /simulate <character_name> <region> <realm> : This command will then send a message component to the user so they can further customize their simulation options

// Interactions fired off in response to application commands (slash commands)
type SaintCommandInteraction string

const (
	SaintSimulate SaintCommandInteraction = "simulate"
	SaintHelp     SaintCommandInteraction = "help"
)

// Interactions fired off in resposne to user input in a MessageComponent
// type SaintMsgComponentInteraction string

// const (
// 	// Called when they
// 	SimulateCharacterRealm SaintMsgComponentInteraction = "simulate_character_realm"
// )

var (

	// Slash commands
	commands = []discordgo.ApplicationCommand{
		{
			Name:        string(SaintSimulate),
			Description: "Simulate your characters DPS.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Description: "What region do you play on?",
					Name:        "region",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "eu",
							Value: "eu",
						},
						{
							Name:  "us",
							Value: "us",
						},
						{
							Name:  "kr",
							Value: "kr",
						}, {
							Name:  "tw",
							Value: "tw",
						}, {
							Name:  "cn",
							Value: "cn",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Name:        "realm",
					Description: "What realm is your character on?",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "us-thrall",
							Value: "thrall",
						},
						{
							Name:  "us-hydraxis",
							Value: "hydraxis",
						},
						{
							Name:  "eu-silvermoon",
							Value: "silvermoon",
						}, {
							Name:  "eu-draenor",
							Value: "draenor",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Description: "What is your characters name?",
					Name:        "character_name",
					MinLength:   utils.IntPtr(2),
					MaxLength:   12,
				},
			},
		},
		{
			Name: string(SaintHelp),

			Description: "View help",
		},
	}

	commandHandlers = map[SaintCommandInteraction]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// https://github.com/bwmarrin/discordgo/tree/master/examples/components
		// https://github.com/kevcenteno/discordgo/blob/f8c5d6c837ef0cd4db6a4b7d03e301d83f3708c4/examples/components/main.go
		SaintSimulate: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Handle case where this handler receives incorrect interaction type (we need application command interactions only)
			if _, ok := i.Interaction.Data.(discordgo.ApplicationCommandInteractionData); !ok {
				errResponse := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Something went wrong, please try again.",
					},
				}
				err := s.InteractionRespond(i.Interaction, errResponse)

				if err != nil {
					fmt.Printf("error sending error response: %v\n", errResponse.Data)
					panic(err)
				}

				return
			}

			fmt.Printf("sim command data: %v\n", i.Interaction.ApplicationCommandData())

			// Unmarshall received options into SimulationOptions struct so we can validate it
			// seems like we need to explicitly create the WoWCharacter struct inside, as go wont allocate
			// memory for the struct itsself (because we accept a pointer to a struct), so it just
			// allocates memory for the pointer, not the struct.
			simOptions := interfaces.SimulationOptions{
				WowCharacter: &interfaces.WoWCharacter{},
			}
			for _, option := range i.ApplicationCommandData().Options {
				fmt.Printf("option %v\n", option)
				switch option.Name {
				case "character_name":
					fmt.Println("character_name found")
					if characterName, ok := option.Value.(string); ok {
						simOptions.WowCharacter.CharacterName = &characterName
					}
				case "realm":
					fmt.Println("realm found")
					if realm, ok := option.Value.(string); ok {
						simOptions.WowCharacter.Realm = &realm
					}
				case "region":
					fmt.Println("region found")
					if region, ok := option.Value.(string); ok {
						simOptions.WowCharacter.Region = &region
					}
				default:
					fmt.Println("defaulted")
				}
			}

			fmt.Printf("extracted: %v\n", simOptions)
			// todo: handle error and send response to user
			if data, err := json.Marshal(simOptions); err != nil {
				fmt.Printf("error marshalling: %v", err)
			} else {
				fmt.Printf("marshalled: %s", data)
			}

			simRes, simErr := simulateCharacter(s, i, "sIshton", "dhydraxis", "usas")
			// todo: implement err handling
			if simErr != nil {
				fmt.Printf("sim err: %v\n", simErr)
				return
			}

			fmt.Printf("Sim interaction data: %v", i.Interaction.Data.Type())

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: *simRes.SimulationId,
				},
			}

			err := s.InteractionRespond(i.Interaction, response)
			if err != nil {
				fmt.Println(response.Data)
				panic(err)
			}
		},
	}
)

var s *discordgo.Session

type SaintError struct{}

func (s SaintError) Error() string {
	return "Something bad happened"
}

func simulateCharacter(s *discordgo.Session, i *discordgo.InteractionCreate, characterName string, characterRealm string, characterRegion string) (*interfaces.SimulationResponse, error) {

	wowCharacter := interfaces.WoWCharacter{
		CharacterName: utils.StrPtr(characterName),
		Realm:         utils.StrPtr(characterRealm),
		Region:        utils.StrPtr(characterRegion),
	}

	requestData := interfaces.SimulateJSONRequestBody{
		WowCharacter: &wowCharacter,
	}

	// Convert requestData to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("Error marshaling request data: %v\n", err)
		return nil, err
	}

	// Define the API URL
	url := "http://saint_api:8080/simulate"

	// Send the POST request with the JSON body
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending POST request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response status: %v\n", resp.StatusCode)
		return nil, SaintError{}
	}

	var simRespose interfaces.SimulationResponse = interfaces.SimulationResponse{
		SimulationId: utils.StrPtr("abc"),
	}

	return &simRespose, nil
}

// Create session object
func init() {
	fmt.Println("Loaded secrets:")
	fmt.Printf("%s: %s\n", constants.DiscordToken.Key(), constants.DiscordToken.MaskedValue())

	var err error
	s, err = discordgo.New("Bot " + constants.DiscordToken.Value())
	if err != nil {
		log.Fatalf("Error occured during discord session creation: %v", err)
		return
	}
}

func init() {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Add interaction handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		fmt.Printf("%s occured\n", i.ID)
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			fmt.Printf("Interaction occurred: %v\n", i.ApplicationCommandData().Name)
			if h, ok := commandHandlers[SaintCommandInteraction(i.ApplicationCommandData().Name)]; ok {
				h(s, i)
			}

		default:
			fmt.Printf("Received interaction of type %v, but we do not have any handlers for this type of interaction", i.Type)
		}

	})

	// Register application commands
	log.Printf("Registering commands...")

	cmdIDS := make(map[string]string, len(commands))
	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(constants.ApplicationID.Value(), constants.GuildID, &cmd)
		if err != nil {
			log.Fatalf("Failed to register command with name '%q': %v", cmd.Name, err)
		} else {
			fmt.Printf("Registered command %q\n", cmd.Name)
		}

		cmdIDS[rcmd.ID] = rcmd.Name
	}
}

func main() {

	// Open websocket connection
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
