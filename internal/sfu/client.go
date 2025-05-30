package sfu

import "github.com/pion/webrtc/v4"

type Client struct {
	pc *webrtc.PeerConnection
}

func (p *Client) GetOffer() (string, error) {
	gatherComplete := webrtc.GatheringCompletePromise(p.pc)

	offer, err := p.pc.CreateOffer(&webrtc.OfferOptions{
		OfferAnswerOptions: webrtc.OfferAnswerOptions{
			VoiceActivityDetection: true,
		},
	})
	if err != nil {
		return "", err
	}

	err = p.pc.SetLocalDescription(offer)
	if err != nil {
		return "", err
	}

	<-gatherComplete

	return p.pc.LocalDescription().SDP, nil
}

func (p *Client) GetAnswer() (string, error) {
	gatherComplete := webrtc.GatheringCompletePromise(p.pc)

	answer, err := p.pc.CreateAnswer(&webrtc.AnswerOptions{
		OfferAnswerOptions: webrtc.OfferAnswerOptions{
			VoiceActivityDetection: true,
		},
	})
	if err != nil {
		return "", err
	}

	err = p.pc.SetLocalDescription(answer)
	if err != nil {
		return "", err
	}

	<-gatherComplete

	return p.pc.LocalDescription().SDP, nil
}

func (p *Client) SetOffer(offer string) error {
	err := p.pc.SetRemoteDescription(webrtc.SessionDescription{
		SDP:  offer,
		Type: webrtc.SDPTypeOffer,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) SetAnswer(answer string) error {
	err := p.pc.SetRemoteDescription(webrtc.SessionDescription{
		SDP:  answer,
		Type: webrtc.SDPTypeAnswer,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Client) Close() error {
	return p.pc.Close()
}
