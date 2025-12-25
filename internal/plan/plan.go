package plan

type Plan struct {
	Actions []Action
}

func (p *Plan) IsEmpty() bool {
	return len(p.Actions) == 0
}

func (p *Plan) Validate() error {
	for _, action := range p.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	return nil
}
