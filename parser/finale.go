package parser

const DEFAULT_FINALE_NAME = "default_resp"

// TODO: This needs to be populated
var finaleCache map[string]*FinaleOpts

func GenerateFinale(key string) *Section {
	f := finaleCache[key]
	
	if f == nil {
		f = finaleCache[DEFAULT_FINALE_NAME]
		if f == nil {
			return &Section{
				Type:     ST_SLIDE,
				Title:    "Woo! You did it <3",
				ID:       FINAL_SECTION_NAME,
				Slide:    &SlideOpts{
					SubTitle: "I am trans :3",
				},
			}
		}
	}

	return &Section{
		Type:     ST_FINAL,
		ID:       FINAL_SECTION_NAME,
		Finale:   f,
	}
}