package services

type InMemoryTokenBlacklister struct {
	blacklistedFamilies map[string]bool
	blacklistedJTIs     map[string]bool
}

func NewInMemoryTokenBlacklister() *InMemoryTokenBlacklister {
	return &InMemoryTokenBlacklister{
		blacklistedFamilies: make(map[string]bool),
		blacklistedJTIs:     make(map[string]bool),
	}
}

func (blacklister *InMemoryTokenBlacklister) BlackListFamily(family string) {
	blacklister.blacklistedFamilies[family] = true
}

func (blacklister *InMemoryTokenBlacklister) InvalidateJTI(jti string) {
	blacklister.blacklistedJTIs[jti] = true
}

func (blacklister *InMemoryTokenBlacklister) ValidateJTI(family string, jti string) bool {
	familyBlacklisted := blacklister.blacklistedFamilies[family]
	jtiBlacklisted := blacklister.blacklistedJTIs[jti]

	return !(familyBlacklisted || jtiBlacklisted)
}
