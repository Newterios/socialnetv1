package model

type Emoji struct {
	ID    string `json:"id"`
	Char  string `json:"char"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

var PredefinedEmojis = []Emoji{
	{ID: "smile", Char: "😊", Name: "Smile"},
	{ID: "cool", Char: "😎", Name: "Cool"},
	{ID: "fire", Char: "🔥", Name: "Fire"},
	{ID: "lightning", Char: "⚡", Name: "Lightning"},
	{ID: "gaming", Char: "🎮", Name: "Gaming"},
	{ID: "art", Char: "🎨", Name: "Art"},
	{ID: "music", Char: "🎵", Name: "Music"},
	{ID: "book", Char: "📚", Name: "Book"},
	{ID: "star", Char: "🌟", Name: "Star"},
	{ID: "strong", Char: "💪", Name: "Strong"},
}

func GetEmojiByID(id string) *Emoji {
	for _, e := range PredefinedEmojis {
		if e.ID == id {
			return &Emoji{ID: e.ID, Char: e.Char, Name: e.Name}
		}
	}
	return nil
}

func IsValidEmoji(id string) bool {
	return GetEmojiByID(id) != nil
}

type EmojiStat struct {
	EmojiID string `json:"emoji_id" bson:"emoji_id"`
	Count   int64  `json:"count" bson:"count"`
}

type EmojiUserInfo struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	EmojiAvatar string `json:"emoji_avatar"`
}
