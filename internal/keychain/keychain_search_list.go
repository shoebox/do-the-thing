package keychain

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (k keychain) addKeyChainToSearchList(ctx context.Context) error {
	list, err := k.getSearchList(ctx)
	if err != nil {
		return err
	}

	return k.setSearchList(ctx, append(list, k.API.PathService.KeyChain()))
}

func (k keychain) listCall(ctx context.Context, args []string) ([]byte, error) {
	return k.API.Exec.
		CommandContext(ctx, SecurityUtil, append([]string{ActionListKeyChains}, args...)...).
		Output()
}

func (k keychain) getSearchList(ctx context.Context) ([]string, error) {
	var res []string
	// Display the the keychain search list without any specified domain
	// TODO: Maybe necessary to define the domain later?
	b, err := k.listCall(ctx, []string{})

	if err != nil {
		return res, err
	}

	return parseSearchList(b), nil
}

func (k keychain) setSearchList(ctx context.Context, list []string) error {
	log.Info().Msg(fmt.Sprintf("Set search list to %v", list))
	b, err := k.listCall(ctx, append([]string{"-s"}, list...))
	log.Debug().Bytes("Result", b).Msg("Search list applied")
	return err
}

func parseSearchList(data []byte) []string {
	var res []string

	// Parse each line and try to parse the keychain p ath
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		txt := scanner.Text()
		l := isSearchListEntry(txt)
		if l != "" {
			res = append(res, l)
		}
	}

	return res
}

func isSearchListEntry(txt string) string {
	if searchListRegexp.MatchString(txt) {
		sm := searchListRegexp.FindStringSubmatch(txt)
		return sm[1]
	}

	return ""
}
