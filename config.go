package lilraft

// Config ...
type Config struct {
	Addresses       []string
	Clusters        []string
	ElectionTimeout uint
}
