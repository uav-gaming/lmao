package commands

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

func ProfileCommandHandler(discord.UserID) (*api.InteractionResponse, error) {
	return &api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{

				{
					Title:       "Title for your user data",
					Description: "Some description",
					Timestamp:   discord.NewEmbed().Timestamp,
					Thumbnail: &discord.EmbedThumbnail{
						URL: "https://assets.leetcode.com/users/lee215/avatar_1551541889.png",
					},

					Type:  discord.NormalEmbed,
					Color: discord.DefaultEmbedColor,
					Fields: []discord.EmbedField{
						{
							Name:  "Rating",
							Value: "123",
						},
						{
							Name:  "Solved",
							Value: "1321",
						},
						{
							Name:  "To solve",
							Value: "-312",
						},
					},
				},
			},
		},
	}, nil
}
