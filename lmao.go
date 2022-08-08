package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	discord_token := os.Getenv("DISCORD_TOKEN")
	if len(discord_token) <= 0 {
		log.Fatal("env DISCORD_TOKEN not set.")
	}
	s, err = discordgo.New("Bot " + discord_token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "progress",
			Description: "Get grind75 progress",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "leetcode username",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"progress": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			username := i.ApplicationCommandData().Options[0].Value.(string)

			if username == "" {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Please provide a username.",
					},
				})
			} else {
				progress := getProgress(username)
				//build message string
				var progressStr strings.Builder
				for key, value := range progress {
					progressStr.WriteString(key + ": " + value + "\n")
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: progressStr.String(),
					},
				})
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func HandleRequest(ctx context.Context, body []byte) (string, error) {
	log.Info("Got request: ", body)
	return "Hello world!", nil
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatal("Cannot open the session: ", err)
	}

	log.Println("Adding commands...")
	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, os.Getenv("DISCORD_GUILD"), v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer s.Close()

	lambda.Start(HandleRequest)
}

// TODO: finish adding the problems
var grind75List = []string{"Two Sum", "Valid Parentheses", "Merge Two Sorted Lists", "Best Time to Buy and Sell Stock", "Valid Palindrome", "Invert Binary Tree", "Valid Anagram", "Binary Search"}

type ResponseData struct {
	Data struct {
		RecentSubmissionList []struct {
			Title string `json:"title"`
		} `json:"recentAcSubmissionList"`
	} `json:"data"`
}

// get progress and returns a map to display the progress
func getProgress(username string) map[string]string {
	//add DB call here
	progressMap := make(map[string]string)
	for _, problem := range grind75List {
		progressMap[problem] = "❌"
	}
	updateProgress(username, progressMap) //placeholder for now
	return progressMap
}

// graphql query to get the new progress and update the progress map
func updateProgress(username string, progressMap map[string]string) {
	query := map[string]string{
		"query": `
            { 
                recentAcSubmissionList(username: "` + username + `", limit: 20) {
					title
				}
            }
        `,
	}
	queryAsJson, _ := json.Marshal(query)
	request, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(queryAsJson))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	data_struct := ResponseData{}
	json.Unmarshal(data, &data_struct)

	submissionListStruct := data_struct.Data.RecentSubmissionList

	submissionList := make([]string, len(submissionListStruct))
	for i, submission := range submissionListStruct {
		submissionList[i] = submission.Title
	}

	for _, problem := range submissionList {
		_, isPresent := progressMap[problem]
		if isPresent {
			progressMap[problem] = "✅"
		}
	}
}
