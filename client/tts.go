package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/garlicgarrison/elevenlabs2/client/types"
)

func (c Client) TTSWriter(ctx context.Context, w io.Writer, modelID, text, voiceID string, options types.SynthesisOptions) error {
	options.Clamp()
	url := fmt.Sprintf(c.endpoint+"/v1/text-to-speech/%s", voiceID)
	opts := types.TTS{
		ModelID:       modelID,
		Text:          text,
		VoiceSettings: options,
	}
	b, _ := json.Marshal(opts)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("User-Agent", "github.com/garlicgarrison/elevenlabs")
	req.Header.Set("accept", "audio/mpeg")
	res, err := client.Do(req)

	switch res.StatusCode {
	case 401:
		return ErrUnauthorized
	case 200:
		if err != nil {
			return err
		}
		defer res.Body.Close()
		io.Copy(w, res.Body)
		return nil
	case 422:
		fallthrough
	default:
		ve := types.ValidationError{}
		defer res.Body.Close()
		jerr := json.NewDecoder(res.Body).Decode(&ve)
		if jerr != nil {
			err = errors.Join(err, jerr)
		} else {
			err = errors.Join(err, ve)
		}
		return err
	}
}

func (c Client) TTS(ctx context.Context, modelID, text, voiceID string, options types.SynthesisOptions) ([]byte, error) {
	options.Clamp()
	url := fmt.Sprintf(c.endpoint+"/v1/text-to-speech/%s", voiceID)
	client := &http.Client{}
	opts := types.TTS{
		ModelID:       modelID,
		Text:          text,
		VoiceSettings: options,
	}
	b, _ := json.Marshal(opts)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("User-Agent", "github.com/garlicgarrison/elevenlabs")
	req.Header.Set("accept", "audio/mpeg")
	res, err := client.Do(req)

	switch res.StatusCode {
	case 401:
		return []byte{}, ErrUnauthorized
	case 200:
		if err != nil {
			return []byte{}, err
		}
		b := bytes.Buffer{}
		w := bufio.NewWriter(&b)

		defer res.Body.Close()
		io.Copy(w, res.Body)
		return b.Bytes(), nil
	case 422:
		fallthrough
	default:
		ve := types.ValidationError{}
		defer res.Body.Close()
		jerr := json.NewDecoder(res.Body).Decode(&ve)
		if jerr != nil {
			err = errors.Join(err, jerr)
		} else {
			err = errors.Join(err, ve)
		}
		return []byte{}, err
	}
}

func (c Client) TTSStream(ctx context.Context, w io.Writer, text, voiceID string, options types.SynthesisOptions) error {
	options.Clamp()
	url := fmt.Sprintf(c.endpoint+"/v1/text-to-speech/%s/stream", voiceID)
	opts := types.TTS{
		Text:          text,
		VoiceSettings: options,
	}
	b, _ := json.Marshal(opts)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("User-Agent", "github.com/garlicgarrison/elevenlabs")
	req.Header.Set("accept", "audio/mpeg")
	res, err := client.Do(req)

	switch res.StatusCode {
	case 401:
		return ErrUnauthorized
	case 200:
		if err != nil {
			return err
		}
		defer res.Body.Close()
		io.Copy(w, res.Body)
		return nil
	case 422:
		fallthrough
	default:
		ve := types.ValidationError{}
		defer res.Body.Close()
		jerr := json.NewDecoder(res.Body).Decode(&ve)
		if jerr != nil {
			err = errors.Join(err, jerr)
		} else {
			err = errors.Join(err, ve)
		}
		return err
	}
}
