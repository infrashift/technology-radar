package radar

import "list"

// #Slug is a URL-safe identifier used for blip ids, quadrants, and rings.
#Slug: =~"^[a-z0-9][a-z0-9-]*$"

// #Entry is a configurable quadrant or ring.
#Entry: {
	slug:  #Slug
	label: string & !=""
}

// #Config is the schema for .radar.yaml. Every field has a default, so an
// empty or missing config file is valid.
//
// Check a config file by hand with:
//   cue vet -d '#Config' schema/radar.cue .radar.yaml
#Config: {
	content_dir:    string | *"content/blips"
	web_dir:        string | *"web"
	output_json:    string | *"web/src/data/blips-generated.json"
	output_content: string | *"web/src/content/blips"
	version:        =~"^\\d{4}\\.\\d{2}$" | *"2025.11"

	// Exactly four of each. Order matters: quadrant index maps to angular
	// position (0 = top-right, clockwise), ring index maps to radial position
	// (0 = innermost).
	quadrants: *[
		{slug: "techniques", label: "Techniques"},
		{slug: "platforms", label: "Platforms"},
		{slug: "tools", label: "Tools"},
		{slug: "languages-and-frameworks", label: "Languages & Frameworks"},
	] | [...#Entry]
	quadrants: list.MinItems(4) & list.MaxItems(4)

	rings: *[
		{slug: "adopt", label: "Adopt"},
		{slug: "trial", label: "Trial"},
		{slug: "assess", label: "Assess"},
		{slug: "hold", label: "Hold"},
	] | [...#Entry]
	rings: list.MinItems(4) & list.MaxItems(4)

	_uniqueQuadrantSlugs: [for q in quadrants {q.slug}] & list.UniqueItems()
	_uniqueRingSlugs: [for r in rings {r.slug}] & list.UniqueItems()
}

// config is filled in by the radar CLI with the loaded .radar.yaml so that
// #Quadrant and #Ring track the configured slugs.
config: #Config

#Quadrant: or([for q in config.quadrants {q.slug}])
#Ring: or([for r in config.rings {r.slug}])

#History: {
	date: =~"^\\d{4}-(0[1-9]|1[0-2])$"
	ring: #Ring
}

// #Blip is the schema for MDX frontmatter. It is closed, so unknown or
// misspelled keys are rejected.
#Blip: {
	id:            #Slug
	title:         string & !=""
	quadrant:      #Quadrant
	ring:          #Ring
	isNew:         bool | *false
	ciIdentifier?: =~"^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	history?: [...#History]
	relatedBlips?: [...#Slug]
}
