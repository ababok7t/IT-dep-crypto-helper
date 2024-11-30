package bot

const (
	StateMain             string = "MAIN"
	StateCoinsList        string = "COINS_LIST"
	StateCoinInfo         string = "COIN_INFO"
	StateSetAlert         string = "SET_ALERT"
	StateAlert            string = "ALERT"
	StateSetCollection    string = "SET_COLLECTION"
	StateRemoveCollection string = "REMOVE_COLLECTION"
	StateCollection       string = "COLLECTION"
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
		} else {
			if string(message[0]) == "@" {
				sm.currentState = StateSetCollection
			} else if string(message[0]) == "$" {
				sm.currentState = StateRemoveCollection
			} else {
				sm.currentState = StateAlert
			}
		}

	case StateAlert:
		if message != "" {
			sm.currentState = StateSetAlert
		}

	case StateSetAlert:
		{
			if message == "далее" {
				sm.currentState = StateCoinsList
			}
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
