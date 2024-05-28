package exploreapi

type matchUserBody struct {
	PreferName string   `json:"name"`
	Gender     string   `json:"gender"`
	Major      string   `json:"major"`
	Native     string   `json:"native"`    // language code with RFC-5646 format
	Country    string   `json:"country"`   // ISO-3166 format
	Learnings  []string `json:"learnings"` // languages code with RFC-5646 format
	Interests  []string `json:"interests"`
	Age        int      `json:"age"`
}
