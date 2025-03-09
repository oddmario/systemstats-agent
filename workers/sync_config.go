package workers

import (
	"time"

	"github.com/oddmario/systemstats-agent/utils"
)

func syncConfig() {
	for {
		utils.LoadConfig(false)

		time.Sleep(3 * time.Second)
	}
}
