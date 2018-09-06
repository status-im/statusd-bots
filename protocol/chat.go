package protocol

import (
	"github.com/ethereum/go-ethereum/crypto/sha3"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

// PublicChatTopic returns a Whisper topic for a public channel name.
func PublicChatTopic(name []byte) (whisper.TopicType, error) {
	hash := sha3.NewKeccak256()
	if _, err := hash.Write(name); err != nil {
		return whisper.TopicType{}, err
	}

	return whisper.BytesToTopic(hash.Sum(nil)), nil
}
