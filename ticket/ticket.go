package ticket

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ap8322/brant/config"
)

var tickets Tickets

type Tickets struct {
	Tickets []Ticket `toml:"Tickets"`
}

type Ticket struct {
	ID    string `toml:"id"`
	Title string `toml:"title"`
}

// Load reads toml file.
func (tickets *Tickets) Load() error {
	cache := config.Conf.Core.TicketCache
	if _, err := toml.DecodeFile(cache, tickets); err != nil {
		return fmt.Errorf("Failed to load ticket cache file. %v", err)
	}
	return nil
}

func (tickets *Tickets) Save() error {
	cache := config.Conf.Core.TicketCache
	f, err := os.Create(cache)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save ticket cache file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(tickets)
}
