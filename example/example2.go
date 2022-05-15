package example

type MixedReceiver struct{}

func (r *MixedReceiver) PointerReceiver() {}

func (r MixedReceiver) ValueReceiver() {}
