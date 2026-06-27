package itadclient

import (
	"fmt"
	"math/rand"
)

type Deals struct {
	Results []Deal `json:"list"`
}

func (d *Deals) GetRandomDeal() *Deal {
	if len(d.Results) == 0 {
		return nil
	}
	// ランダムに1つのゲームを選択
	randomIndex := rand.Intn(len(d.Results))
	return &d.Results[randomIndex]
}

type Deal struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Slug   string      `json:"slug"`
	Mature bool        `json:"mature"`
	Type   GameTypeStr `json:"type"`
	Deal   DealDetail  `json:"deal"`
}

func (d Deal) String() string {
	return d.Title + " - " + d.Deal.Price.Currency + " " + fmt.Sprintf("%.2f", d.Deal.Price.Amount)
}

type GameTypeStr string

const (
	GameTypeNone    GameTypeStr = ""
	GameTypeGame    GameTypeStr = "game"
	GameTypeBundle  GameTypeStr = "bundle"
	GameTypePackage GameTypeStr = "package"
)

type GameType int

const (
	GameTypeIntNone    GameType = iota
	GameTypeIntGame    GameType = 1
	GameTypeIntBundle  GameType = 2
	GameTypeIntPackage GameType = 3
)

type DealDetail struct {
	Shop  Shop   `json:"shop"`
	Price Price  `json:"price"`
	Cut   int    `json:"cut"`
	URL   string `json:"url"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (d *Price) String() string {
	return fmt.Sprintf("%.2f %s", d.Amount, d.Currency)
}

type Shop struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	// SteamShopID は IsThereAnyDeal API における Steam のショップ ID です
	SteamShopID     = 61
	JPCountryCode   = "JP"
	MinReviewCount  = 1000
	MinReviewRating = 80
)

type GameDetail struct {
	ID          string      `json:"id"`
	AppID       int         `json:"appid"`
	Title       string      `json:"title"`
	Assets      Assets      `json:"assets"`
	PlayerStats PlayerStats `json:"players"`
}

type Assets struct {
	Banner165 string `json:"banner165"`
	Banner300 string `json:"banner300"`
	Banner400 string `json:"banner400"`
	Banner600 string `json:"banner600"`
}

type PlayerStats struct {
	Peak int `json:"peak"`
}

func (a *Assets) GetBestBanner() string {
	if a.Banner600 != "" {
		return a.Banner600
	}
	if a.Banner400 != "" {
		return a.Banner400
	}
	if a.Banner300 != "" {
		return a.Banner300
	}
	return a.Banner165
}

const PopularityThreshold = 10000

func (g GameDetail) IsPopular() bool {
	return g.PlayerStats.Peak > PopularityThreshold
}

type DealsRequest struct {
	Limit   int         `json:"limit"`
	Shops   []int       `json:"shops"`
	Country string      `json:"country"`
	Filter  DealsFilter `json:"filter"`
}

type DealsFilter struct {
	Platforms  Platform    `json:"platform"`
	Types      []GameType  `json:"type"`
	SteamCount FilterRange `json:"steamCount"`
	SteamPerc  FilterRange `json:"steamPerc"`
	TagsUnion  []string    `json:"tagsUnion"`
}

type Platform int

const (
	PlatformNone    Platform = iota
	PlatformWindows Platform = 1
	PlatformMac     Platform = 2
)

type FilterRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
