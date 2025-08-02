package news

// go=Product js=ProductJson sql=products demo gotiny msgp
type product struct {
	id int // sql{inc}

	created int64 // sql{unix}
	updated int64 // sql{unix}

	active bool //sql{defaults}

	ref      string
	sku      string
	category string // vinyl
	brand    string //sql{altertable}
	model    string //sql{ver="2"}

	title    string //James
	line     string //A good guy
	about    string //md sql{renames="oldname"}
	image    string //sql{ver="October"}
	original string //link to product

	// price    float64 // 2.42 sqft
	cost     float64 // 1.15,
	discount float64 // 5 percent
	sample   float64 // 1.5 cost usd

	width     float64 // 7.1 in
	length    float64 // 48 in
	size      float64 // 540 sqft
	thickness float64 // mm

	count     int     // 30 count
	price     float64 // 98.4 usd
	boxwidth  float64 // 200 in
	boxheight float64 // 400 in
	boxsize   float64 // 1400 sqft
	boxweight float64 // 30 kg

	edge         string //square, bevel
	flooring     string //tile, plank
	installation string //glue, plank
	construction string //attached, attached pad, floating
	gloss        string //low,
	waste        int    //Recommended Waste Factor percents
	wear         int    //wear layers, mil

	features  map[string]bool //waterproof
	layers    map[string]bool //waterproof, vinyl, soundproof
	materials map[string]bool //vinyl, plastic
	meta      map[string]any
}
