package seventv

type EventsApiStyle struct{
	color int
}

type EventsApiConnection struct{
	Id string			`json:"id"`
	Platform string			`json:"platform"`
	Username string			`json:"username"`
	Display_name string			`json:"display_name"`
	Linked_at int			`json:"linked_at"`
	Emote_capacity int			`json:"emote_capacity"`
	Emote_set_id any			`json:"emote_set_id"`
}

type EventsApiUser struct{
	Id string			`json:"id"`
	Username string			`json:"username"`
	Display_name string			`json:"display_name"`
	Avatar_url string			`json:"avatar_url"`
	Style EventsApiStyle			`json:"style"`
	Role_ids []string			`json:"role_ids"`
	Connections []EventsApiConnection			`json:"connections"`
}

type EventsApiEntitlementCreate struct{
	D struct{
		Type string			`json:"type"`
		Body struct{
			Id string			`json:"id"`
			Kind int			`json:"kind"`
			Object struct{
				Id string			`json:"id"`
				Kind string			`json:"kind"`
				Ref_id string			`json:"ref_id"`
				User EventsApiUser			`json:"user"`
			}						`json:"object"`
		}							`json:"body"`
	}				`json:"d"`
	Op int			`json:"op"`
	T int			`json:"t"`
	S int			`json:"s"`
}

type EventsApiEmoteSetUpdate struct{
	D struct{
		Type string			`json:"type"`
		Body struct{
			Id string			`json:"id"`
			Kind int			`json:"kind"`
			Pushed []struct{
				Key string			`json:"key"`
				Index int			`json:"index"`
				Type string			`json:"type"`
				Value struct{
					Id string			`json:"id"`
					Name string			`json:"name"`
					Flags int			`json:"flags"`
					Timestamp int			`json:"timestamp"`
					Actor_id int			`json:"actor_id"`
					Data struct{
						Id string			`json:"id"`
						Name string			`json:"name"`
						Flags int			`json:"flags"`
						Tags []string			`json:"tags"`
						Lifecycle int			`json:"lifecycle"`
						State []string			`json:"state"`
						Listed bool			`json:"listed"`
						Animated bool			`json:"animated"`
						Owner EventsApiUser			`json:"owner"`
						Host struct{
							Url string			`json:"url"`
							Files[]struct{
								Name string			`json:"name"`
								Static_name string			`json:"static_name"`
								Width int			`json:"width"`
								Height int			`json:"height"`
								Frame_count int			`json:"frame_count"`
								Size int			`json:"size"`
								Format string			`json:"format"`
							}		`json:"files"`
						}		`json:"host"`
					}		`json:"data"`
					Origin_id any 			`json:"origin_id"`
				}	`json:"value"`
			}		`json:"pushed"`
		}			`json:"body"`
	}				`json:"d"`
	Op int			`json:"op"`
	T int			`json:"t"`
	S int			`json:"s"`
}

type EventsApiEmoteSetCreate struct{
	D struct{
		Type string			`json:"type"`
		Body struct{
			Id string			`json:"id"`
			Kind int			`json:"kind"`
			Object struct{
				Id string			`json:"id"`
				Name string			`json:"name"`
				Flags int			`json:"flags"`
				Tags []string			`json:"tags"`
				Immutable bool			`json:"immutable"`
				Privileged bool			`json:"privileged"`
				Capacity int			`json:"capacity"`
				Owner EventsApiUser			`json:"owner"`
			}	`json:"object"`
		}		`json:"body"`
	}		`json:"d"`
	Op int			`json:"op"`
	T int			`json:"t"`
	S int			`json:"s"`
}

type EventsApiCosmeticCreate struct{
	D struct{
		Type string			`json:"type"`
		Body struct{
			Id string			`json:"id"`
			Kind int			`json:"kind"`
			Object struct{
				Id string			`json:"id"`
				Kind string			`json:"kind"`
				Data struct{
					Id string			`json:"id"`
					Name string			`json:"name"`
					Color any // unknown (was null)			`json:"color"`
					Gradients []any // unknown			`json:"gradients"`
					Shadows []struct{
						X_offset float64			`json:"x_offset"`
						Y_offset float64			`json:"y_offset"`
						Radius float64			`json:"radius"`
						Color int			`json:"color"`
					}					`json:"shadows"`
					Text *string			`json:"text"`
					Function string			`json:"function"`
					Repeat bool			`json:"repeat"`
					Angle int			`json:"angle"`
					Shape string			`json:"shape"`
					Image_url string			`json:"image_url"`
					Stops []struct{
						At float64			`json:"at"`
						Color int			`json:"color"`
					}						`json:"stops"`
					
				}				`json:"data"`
			}		`json:"object"`
		}		`json:"data"`
	}				`json:"d"`
	Op int			`json:"op"`
	T int			`json:"t"`
	S int			`json:"s"`
}
