package entity

type Hero struct {
	Name         string
	Health       int
	maxHealth    int
	recoveryRate int
}

func (h *Hero) DealDamage(i int) {
	h.Health -= 1
	if h.Health < 0 {
		h.Health = 0
	}
}

func (h *Hero) IsDead() bool {
	return h.Health == 0
}
