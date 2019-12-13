package protocol

import (
	"golang.org/x/crypto/sha3"

	"github.com/status-im/status-go/whisper"
)

const (
	// MailServerPassword is required to make requests to MailServers.
	MailServerPassword = "status-offline-inbox"
)

// PublicChatTopic returns a Whisper topic for a public channel name.
func PublicChatTopic(name []byte) (whisper.TopicType, error) {
	hash := sha3.NewLegacyKeccak256()
	if _, err := hash.Write(name); err != nil {
		return whisper.TopicType{}, err
	}

	return whisper.BytesToTopic(hash.Sum(nil)), nil
}
