package discordclient

type Notification struct {
	Thumbnail   string
	Title       string
	Description string
	Price       string
	URL         string
}

type DiscordMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url,omitempty"`
	Color       int        `json:"color,omitempty"`
	Thumbnail   EmbedImage `json:"thumbnail,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url"`
}
