package cache

type AuthToken struct {
	Key string
}

func (a AuthToken) String() string {
	return "AuthToken_" + a.Key
}

type SocialUrl struct {
	Key string
}

func (s SocialUrl) String() string {
	return "SocialUrl_" + s.Key
}

type HashUrl struct {
	Key string
}

func (h HashUrl) String() string {
	return "HashUrl_" + h.Key
}
