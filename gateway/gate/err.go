package gateway

import (
	"errors"
)

var (
	/* Config Errors */
	ErrMissingValueForConfigOption  = errors.New("Missing value for config option")
	ErrTooManyValuesForConfigOption = errors.New("Too many values for config option")
	ErrUnknownConfigOption          = errors.New("Unknown config option")
	ErrNoTransportSpecified         = errors.New("Missing transport")
	ErrInvalidModeSpecified         = errors.New("Invalid mode")
	ErrNotANumber                   = errors.New("Not a number")

	/* Protocol Errors */
	ErrZeroLengthClientID = errors.New("Zero-length clientID is invalid")
	ErrClientIDTooLong    = errors.New("ClientID too long")

	/* Topic Errors */
	ErrTopicFilterEmptyString     = errors.New("TopicFilter cannot be empty string")
	ErrTopicFilterInvalidWildcard = errors.New("TopicFilter contains invalid wildcard")
	ErrTopicNameEmptyString       = errors.New("TopicName cannot be empty string")
	ErrTopicNameContainsWildcard  = errors.New("TopicName cannot contain wildcard")

	/* Topic Tree Errors */
	ErrNoSuchSubscriptionExists = errors.New("Subscription does not exist")
	ErrNoSubscribers            = errors.New("No subscribers")
	ErrClientNotSubscribed      = errors.New("Client not subscribed")
)

func chkerr(e error) {
	if e != nil {
		panic(e)
	}
}
