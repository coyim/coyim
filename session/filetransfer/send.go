package filetransfer

import (
	"errors"
	"strings"

	"github.com/coyim/coyim/session/access"
)

func discoverSupport(s access.Session, peer string) (profiles []string, err error) {
	if res, ok := s.Conn().DiscoveryFeatures(peer); ok {
		foundSI := false
		for _, feature := range res {
			if feature == "http://jabber.org/protocolx/si" {
				foundSI = true
			} else if strings.HasPrefix(feature, "http://jabber.org/protocolx/si/profile/") {
				profiles = append(profiles, feature)
			}
		}

		if !foundSI {
			return nil, errors.New("Peer doesn't support stream initiation")
		}

		if len(profiles) == 0 {
			return nil, errors.New("Peer doesn't support any stream initiation profiles")
		}

		return profiles, nil
	}
	return nil, errors.New("Problem discovering the features of the peer")
}

// InitSend starts the process of sending a file to a peer
func InitSend(s access.Session, peer string, file string) error {
	_, err := discoverSupport(s, peer)
	if err != nil {
		return err
	}
	return nil
}
