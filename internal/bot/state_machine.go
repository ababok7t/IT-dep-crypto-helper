package bot

const (
	StateMain          string = "MAIN"
	StateCoinsList     string = "COINS_LIST"
	StateCoinInfo      string = "COIN_INFO"
	StateSetAlert      string = "SET_ALERT"
	StateSetCollection string = "SET_COLLECTION"
	StateCollection    string = "COLLECTION"
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

		if message == "избранное" {
			sm.currentState = StateCollection
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
		} else if message == "установить alert" {
			sm.currentState = StateSetAlert
		} else {
			sm.currentState = StateSetCollection
		}

	case StateSetAlert:
		if message == "назад" {
			sm.currentState = StateCoinInfo
		}

	case StateSetCollection:
		if message == "далее" {
			sm.currentState = StateCoinsList
		}

	case StateCollection:
		if message == "назад" {
			sm.currentState = StateMain
		} else {
			sm.currentState = StateCoinInfo
		}
	}

	return sm.currentState
}
