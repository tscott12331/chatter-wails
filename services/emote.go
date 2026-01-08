package services

import (
	"bytes"
	"chatter-wails/internal/api"
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
)

type EmoteService struct {
	Ctx context.Context
	GlobalEmotes *[]AppEmote
}

func NewEmoteService() *EmoteService {
	return &EmoteService{}
}


func (es *EmoteService) GetGlobalEmotes(access_token string) (*[]AppEmote, error) {
	res, err := api.ApiGetGlobalEmotes(access_token, map[string][]string{})
	if err != nil {
		log.Printf("[GetGlobalEmotes]: An error occurred fetching global emotes, aborting\n%v\n\n", err)
		return nil, err
	}

	if res.Status != 200 {
		log.Printf("[GetGlobalEmotes]: Failed to get global emotes, aborting\n%v\n\n", res.Body)
		return nil, &api.StatusError[api.ApiGetGlobalEmotesRes]{ Res: res }
	}


	emotes := []AppEmote{}
	tmpl := res.Body.Template
	for _, emote := range res.Body.Data {
		emotes = append(emotes, *es.GetAppEmoteFromApiEmote(emote, tmpl))
	}

	es.GlobalEmotes = &emotes

	return &emotes, nil
}


type AppEmote struct{
	Id string					`json:"id"`
	Name string					`json:"name"`
	LightSrcSet string		`json:"lightSrcSet"`
	DarkSrcSet string		`json:"darkSrcSet"`
}

var TMPL_ID_RPL = []byte("{{id}}")
var TMPL_FORMAT_RPL = []byte("{{format}}")
var TMPL_THEME_RPL = []byte("{{theme_mode}}")
var TMPL_SCALE_RPL = []byte("{{scale}}")

var B_DARK_THEME = []byte("dark")
var B_LIGHT_THEME = []byte("light")

func (es *EmoteService) GetAppEmoteFromApiEmote(apiEmote api.ApiGlobalEmote, tmpl string) *AppEmote {
	hasDark := slices.Contains(apiEmote.Theme_mode, "dark")
	hasLight := slices.Contains(apiEmote.Theme_mode, "light")
	appEmote := &AppEmote{
		Id: apiEmote.Id,
		Name: apiEmote.Name,
	}

	byteTmpl := []byte(tmpl)

	var format string
	if slices.Contains(apiEmote.Format, "animated") {
		format = "animated"
	} else {
		format = "static"
	}

	byteTmpl = bytes.ReplaceAll(byteTmpl, TMPL_ID_RPL, []byte(apiEmote.Id))
	byteTmpl = bytes.ReplaceAll(byteTmpl, TMPL_FORMAT_RPL, []byte(format))


	//https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}
	for i, scale := range apiEmote.Scale {
		scaleSrc := bytes.ReplaceAll(byteTmpl, TMPL_SCALE_RPL, []byte(scale))
		isEnd := i == len(apiEmote.Scale) - 1
		
		if(hasDark) {
			themeSrc := bytes.ReplaceAll(scaleSrc, TMPL_THEME_RPL, B_DARK_THEME)
			appEmote.DarkSrcSet += fmt.Sprintf("%s %sx", themeSrc, scale)
			if !isEnd {
				appEmote.DarkSrcSet += ", "
			}
		}

		if(hasLight) {
			themeSrc := bytes.ReplaceAll(scaleSrc, TMPL_THEME_RPL, B_LIGHT_THEME)
			appEmote.LightSrcSet += fmt.Sprintf("%s %sx", themeSrc, scale)
			if !isEnd {
				appEmote.LightSrcSet += ", "
			}
		}
	}

	return appEmote
}





var scales = []string{"1.0", "2.0", "3.0"}
var emoteTemplate = []byte("https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}")
var B_COMMA_SPACE = []byte(", ")

func GetEmoteSrcSet(id string, formats []string, optTheme *string, optScales *[]string) string {
	var format string
	if slices.Contains(formats, "animated") {
		format = "animated"
	} else {
		format = "static"
	}

	var themeRpl []byte = B_LIGHT_THEME // default to light to be safe
	if optTheme != nil && *optTheme == "dark" {
		themeRpl = B_DARK_THEME
	}

	var scaleSet []string = scales
	if optScales != nil {
		scaleSet = *optScales
	}

	byteTmpl := bytes.ReplaceAll(emoteTemplate, TMPL_ID_RPL, []byte(id))
	byteTmpl = bytes.ReplaceAll(byteTmpl, TMPL_FORMAT_RPL, []byte(format))


	var srcSetBuilder strings.Builder
	srcSetBuilder.Grow(252)

	//https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}
	for i, scale := range scaleSet {
		scaleSrc := bytes.ReplaceAll(byteTmpl, TMPL_SCALE_RPL, []byte(scale))
		isEnd := i == len(scaleSet) - 1
		
		themeSrc := bytes.ReplaceAll(scaleSrc, TMPL_THEME_RPL, themeRpl)

		fmt.Fprintf(&srcSetBuilder, "%s %sx", themeSrc, scale)

		if !isEnd {
			srcSetBuilder.Write(B_COMMA_SPACE)
		}
	}

	return srcSetBuilder.String()
}
