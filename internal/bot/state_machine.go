package bot

const (
	StateMain          string = "MAIN"
	StateCoinsList     string = "COINS_LIST"
	StateCoinInfo      string = "COIN_INFO"
	StateSetAlert      string = "SET_ALERT"
	StateSetCollection string = "SET_COLLECTION"
)

type StateMachine struct {
	currentState string
}

func NewStateMachine() *StateMachine {
	return &StateMachine{currentState: StateMain}
}

func (sm *StateMachine) SetState(message string) string {

	switch sm.currentState {

	case StateMain:
		if message == "перейти к списку монет" {
			sm.currentState = StateCoinsList
		}

	case StateCoinsList:
		if message == "назад" {
			sm.currentState = StateMain
		} else {
			sm.currentState = StateCoinInfo
		}

	case StateCoinInfo:
		if message == "назад" {
			sm.currentState = StateCoinsList
		}
	}
	return sm.currentState
}
