package seventv

import (
	"chatter-wails/internal/api"
	"chatter-wails/services/emote"
	"fmt"
	"net/url"
	"strings"
)

var (
	
	USER_ENDPOINT = url.URL{
		Scheme: "https",
		Host: "7tv.io",
		// "/v3/users/{platform}/{platform_id}",
		Path: "/v3/users",
	}
)

func constructUserEndpointURL(platform string, platform_id string) *url.URL {
	var endpoint url.URL
	endpoint.Scheme = USER_ENDPOINT.Scheme
	endpoint.Host = USER_ENDPOINT.Host
	endpoint.Path = fmt.Sprintf("%s/%s/%s", USER_ENDPOINT.Path, platform, platform_id)

	return &endpoint
}

type SevenTVEmoteDataOwnerConnections struct{
	Id string					`json:"id"`
	Platform string				`json:"platform"`
	Username string				`json:"username"`
	Display_name string				`json:"display_name"`
	Linked_at int				`json:"linked_at"`
	Emote_capacity int				`json:"emote_capacity"`
	Emote_set_id string				`json:"emote_set_id"`
}

type SevenTVEmoteDataOwner struct{
	Id string				`json:"id"`
	Username string				`json:"username"`
	Display_name string				`json:"display_name"`
	Avatar_url string				`json:"avatar_url"`
	// style nil				`json:"style"`
	Role_ids []string				`json:"role_ids"`
	Connections []SevenTVEmoteDataOwnerConnections				`json:"connections"`
}

type SevenTVEmoteDataHostFile struct{
	Name string				`json:"name"`
	Static_name string				`json:"static_name"`
	Width int				`json:"width"`
	Height int				`json:"height"`
	Frame_count int				`json:"frame_count"`
	Size int				`json:"size"`
	Format string				`json:"format"`
}

type SevenTVEmoteDataHost struct{
	Url string				`json:"url"`
	Files []SevenTVEmoteDataHostFile				`json:"files"`
}

type SevenTVEmoteData struct{
	Id string				`json:"id"`
	Name string				`json:"name"`
	Flags int				`json:"flags"`
	Tags []string				`json:"tags"`
	Lifecycle int				`json:"lifecycle"`
	State []string				`json:"state"`
	Listed bool				`json:"listed"`
	Animated bool				`json:"animated"`
	Owner SevenTVEmoteDataOwner				`json:"owner"`
	Host SevenTVEmoteDataHost				`json:"Host"`
}

// TODO: finish types
type SevenTVEmote struct{
	Id string				`json:"id"`
	Name string				`json:"name"`
	Flags int				`json:"flags"`
	Timestamp int				`json:"timestamp"`
	// actor_id null				`json:"actor_id"`
	Data SevenTVEmoteData				`json:"data"`
	// origin_id null				`json:"origin_id"`
}

type SevenTVEmoteSet struct{
	Id string				`json:"id"`
	Name string				`json:"name"`
	Flags int				`json:"flags"`
	Tags []string				`json:"tags"` // probably
	Immutable bool				`json:"immutable"`
	Privileged bool				`json:"privileged"`
	Emotes []SevenTVEmote				`json:"emotes"`
	Emote_count int				`json:"emote_count"`
	Capacity int				`json:"capacity"`
	// owner null				`json:"owner"` // TODO: figure out type
}

type ApiGetSevenTVUserRes struct{
	Id string				`json:"id"`
	Platform string				`json:"platform"`
	Username string				`json:"username"`
	Display_name string				`json:"display_name"`
	Linked_at int				`json:"linked_at"`
	Emote_capacity int				`json:"emote_capacity"`
	Emote_set_id string				`json:"emote_set_id"`
	Emote_set SevenTVEmoteSet				`json:"emote_set"`
	User SevenTVUser				`json:"user"`
}

type SevenTVUserStyle struct{
	Color int				`json:"color"`
	Paint_id string				`json:"paint_id"`
	Badge_id string				`json:"badge_id"`
}

type SevenTVUserEmoteSet struct{
	Id string				`json:"id"`
	Name string				`json:"name"`
	Flags int				`json:"flags"`
	Tags []string				`json:"tags"`
	Capacity int				`json:"capacity"`
}

type SevenTVUserEditor struct{
	Id string				`json:"id"`
	Permissions int				`json:"permissions"`
	Visible bool				`json:"visible"`
	Added_at int				`json:"added_at"`
}

type SevenTVUserConnection struct{
	Id string				`json:"id"`
	Platform string				`json:"platform"`
	Username string				`json:"username"`
	Display_name string				`json:"display_name"`
	Linked_at int				`json:"linked_at"`
	Emote_capacity int				`json:"emote_capacity"`
	Emote_set_id string				`json:"emote_set_id"`
	// emote_set nil				`json:"emote_set"` ??
}

type SevenTVUser struct{
	Id string				`json:"id"`
	Username string				`json:"username"`
	Display_name string				`json:"display_name"`
	Created_at int				`json:"created_at"`
	Avatar_url string				`json:"avatar_url"`
	Style SevenTVUserStyle				`json:"style"`
	Emote_sets []SevenTVUserEmoteSet				`json:"emote_sets"`
	Editors []SevenTVUserEditor				`json:"editors"`
	Roles []string				`json:"roles"`
	Connections []SevenTVUserConnection				`json:"connections"`
}

func GetSevenTVUser(platform string, platform_id string) (*ApiGetSevenTVUserRes, error) {
	endpoint := constructUserEndpointURL(platform, platform_id)
	
	res, err := api.ApiGet[ApiGetSevenTVUserRes](*endpoint, nil, map[string][]string{})
	if err != nil {
		println(err.Error())
		return nil, err
	}

	return &res.Body, nil
}

func GetAppEmotesFromSevenTVUserRes(res *ApiGetSevenTVUserRes) map[string]*emote.AppEmote {
	var appEmotes map[string]*emote.AppEmote = map[string]*emote.AppEmote{}
	for _, stvEmote := range res.Emote_set.Emotes {
		srcSet := getSevenTVEmoteSrcSet(&stvEmote)
		appEmote := emote.AppEmote{
			Id: stvEmote.Id,
			Name: stvEmote.Name,
			LightSrcSet: srcSet,
			DarkSrcSet: srcSet,
			Type: "seventv",
		}

		appEmotes[appEmote.Name] = &appEmote
	}

	return appEmotes
}

func getSevenTVEmoteSrcSet(emote *SevenTVEmote) string {
	files := emote.Data.Host.Files
	if len(files) == 0 {
		return ""
	}

	var sb strings.Builder
	base := fmt.Sprintf("https:%s", emote.Data.Host.Url)

	fmt.Fprintf(&sb, "%s/%s %s", base, files[0].Name, "1x")

	for i, f := range emote.Data.Host.Files {
		fmt.Fprintf(&sb, ", %s/%s %dx", base, f.Name, i)
	}

	return sb.String()
}
