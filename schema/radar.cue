package radar

#Ring:     "adopt" | "trial" | "assess" | "hold"
#Quadrant: "techniques" | "platforms" | "tools" | "languages-and-frameworks"

#History: {
	date: =~"^\\d{4}-\\d{2}$"
	ring: #Ring
}

#Blip: {
	id:            =~"^[a-z0-9][a-z0-9-]*$"
	title:         string
	quadrant:      #Quadrant
	ring:          #Ring
	isNew:          bool | *false
	ciIdentifier?:  =~"^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	history?:       [...#History]
	relatedBlips?:  [...=~"^[a-z0-9][a-z0-9-]*$"]
}

#Manifest: {
	version: string
	blips: [...#Blip]
}
