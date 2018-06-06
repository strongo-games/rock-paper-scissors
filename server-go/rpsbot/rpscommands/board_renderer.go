package rpscommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/app"
	"bytes"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"fmt"
	"time"
	"strconv"
	"net/url"
	"github.com/strongo/bots-api-telegram"
	"github.com/prizarena/turn-based"
)

func renderGameMessage(whc bots.WebhookContext, t strongo.SingleLocaleTranslator, board turnbased.Board) (m bots.MessageFromBot, err error) {
	m.IsEdit = true
	m.Format = bots.MessageFormatHTML
	var s bytes.Buffer
	s.WriteString(whc.Translate(rpstrans.GameCardTitle))
	s.WriteString("\n")
	if board.BoardEntity == nil {
		board.BoardEntity = &turnbased.BoardEntity{Round: 1}
	}
	userMoves := board.UserMoves.Strings()
	movesCount := len(userMoves)
	switch {
	case movesCount == 0: // not started
		s.WriteString(whc.Translate(rpstrans.AskToMakeMove))
	case movesCount == 1 || userMoves[0] == "" || userMoves[1] == "": // one player made move, waiting for second player
		s.WriteString("\n\n")
		firstPlayerName := board.UserNames[0]
		if firstPlayerName == "" {
			firstPlayerName = "Player #1"
		}
		s.WriteString(t.Translate(rpstrans.FirstMoveDoneAwaitingSecond, firstPlayerName))
	case movesCount == 2: // game completed
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("Game completed, moves: %v v %v", userMoves[0], userMoves[1]))
		s.WriteString("\n\n")
		s.WriteString(whc.Translate(rpstrans.AskToMakeMove))
	default:
		panic(fmt.Sprintf("len(game.UserMoves) > 2: %v", movesCount))
	}
	s.WriteString(fmt.Sprintf("\n\nLast updated: %v", time.Now()))
	m.Text = s.String()
	callbackPrefix := fmt.Sprintf("bet?l=%v&r=%v&", board.Lang, strconv.Itoa(board.Round))
	if board.ID != "" {
		callbackPrefix += "b=" + url.QueryEscape(board.ID) + "&"
	}

	inlineQuery := ""
	inlineMove := func(text, id string) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{Text: whc.Translate(text), CallbackData: callbackPrefix + "m=" + t.Translate(id)}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			inlineMove(rpstrans.Option1emoji, rpstrans.Option1code),
			inlineMove(rpstrans.Option2emoji, rpstrans.Option2code),
			inlineMove(rpstrans.Option3emoji, rpstrans.Option3code),
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:              whc.Translate(rpstrans.ChallengeFriendCommandText),
				SwitchInlineQuery: &inlineQuery,
			},
		},
	)

	if board.ID != "" && (board.BoardEntity == nil || board.LeftTournament.IsZero()) {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{
				Text:         "🚫 Leave tournament",
				CallbackData: fmt.Sprintf("%v?board=%v", turnbased.LeaveTournamentCommandCode, board.ID),
			},
		})
	}

	m.Keyboard = keyboard
	return
}
