package handlers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"gitlab.com/Dank-del/lastfm-tgbot/config"
	"gitlab.com/Dank-del/lastfm-tgbot/database"
	last_fm "gitlab.com/Dank-del/lastfm-tgbot/last.fm"
)

const (
	historyCB     = "h"
	historyPrefix = historyCB + "_"
	nextSuffix    = "_" + nextValue
	nextValue     = "n"
	backSuffix    = "_" + backValue
	backValue     = "b"

	limitTracks  float64 = 20
	sleepTimeout         = 30 * time.Minute
)

type historyData struct {
	tracks      [][]last_fm.Track
	currentPage int
	totalPages  int
	owner       *gotgbot.User
	t           time.Time
}

var historyMap map[string]*historyData
var historyMutex *sync.Mutex = &sync.Mutex{}

func checkHistoryMap() {
	if historyMap == nil {
		return
	}

	for {
		time.Sleep(sleepTimeout)
		if len(historyMap) == 0 {
			if historyMap == nil {
				return
			}

			continue
		}
		for key, value := range historyMap {
			if time.Since(value.t) > sleepTimeout {
				historyMutex.Lock()
				delete(historyMap, key)
				historyMutex.Unlock()
			}
		}
	}
}

func historyCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.Message
	user := msg.From
	_, err := b.SendChatAction(msg.Chat.Id, "typing")
	if err != nil {
		return err
	}
	dbuser, err := database.GetLastFMUserFromDB(user.Id)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}
	if dbuser.LastFmUsername == "" {
		_, err := msg.Reply(b, "<i>You haven't registered yourself on this bot yet</i>\n<b>Use /setusername</b>",
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	var theuser *last_fm.LastFMUser
	theuser, err = last_fm.GetLastFMUser(dbuser.LastFmUsername)
	if err != nil {
		return err
	}
	var limit int
	limit, err = strconv.Atoi(theuser.User.Playcount)
	if err != nil {
		_, err := msg.Reply(b,
			mdparser.GetItalic("failed to get your total play count").ToString(),
			config.GetDefaultMdOpt())
		if err != nil {
			return err
		}
		return err
	}

	grc, err := last_fm.GetRecentTracksByUsername(dbuser.LastFmUsername, limit)
	if err != nil {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", err.Error()),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	if grc.Error != 0 {
		_, err := msg.Reply(b, fmt.Sprintf("<i>Error: %s</i>", grc.Message),
			&gotgbot.SendMessageOpts{ParseMode: "html"})
		return err
	}

	tracks := grc.Recenttracks.Track

	key := fuckAbs(ctx.EffectiveChat.Id) +
		"_" + fuckAbs(ctx.EffectiveUser.Id)
	cb := historyPrefix + key

	kb := &gotgbot.InlineKeyboardMarkup{}
	kb.InlineKeyboard = append(kb.InlineKeyboard, []gotgbot.InlineKeyboardButton{})
	kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "⇽",
		CallbackData: cb + backSuffix,
	})
	kb.InlineKeyboard[0] = append(kb.InlineKeyboard[0], gotgbot.InlineKeyboardButton{
		Text:         "⇾",
		CallbackData: cb + nextSuffix,
	})

	if len(tracks) < int(limitTracks) {
		_, err = msg.Reply(b, getSimpleList(tracks, user, 0).ToString(),
			&gotgbot.SendMessageOpts{
				ParseMode:             "markdownv2",
				DisableWebPagePreview: true,
				//ReplyMarkup:           kb,
			})
		return err
	}

	thevalue := &historyData{
		//tracks:      tracks,
		currentPage: 1,
		totalPages:  int(math.Ceil(float64(len(tracks)) / limitTracks)),
		owner:       user,
		t:           time.Now(),
	}

	thevalue.GenerateWholeList(tracks)

	m := strings.Replace(thevalue.GetParsedText().ToString(), "|", `\|`, -1)
	// fmt.Println(m)
	_, err = msg.Reply(b, m,
		&gotgbot.SendMessageOpts{
			ParseMode:             "markdownv2",
			DisableWebPagePreview: true,
			ReplyMarkup:           kb,
		})

	// make sure that the err is nil; then allocate the memory.
	if err == nil || strings.Contains(err.Error(), "failed to execute") {
		historyMutex.Lock()
		if historyMap == nil {
			historyMap = make(map[string]*historyData)
			go checkHistoryMap()
		}
		historyMap[key] = thevalue
		historyMutex.Unlock()
	}
	return err
}

func getSimpleList(tracks []last_fm.Track,
	user *gotgbot.User, offset int) mdparser.WMarkDown {
	m := mdparser.GetBold("Recently played tracks by ").AppendMention(user.FirstName, user.Id).AppendNormal("\n\n")
	for a, e := range tracks {
		m = m.AppendNormal(fmt.Sprintf("%d", offset+a+1)).AppendNormal(": ")
		m = m.AppendHyperLink(fmt.Sprintf("%s - %s", e.Artist.Name, e.Name), e.URL)
		if e.Loved == "1" {
			m = m.AppendItalic(" (Loved ♥)")
		}

		m = m.AppendNormal("\n")
		if e.Album.Text != "" {
			m = m.AppendBold("From album: ")
			m = m.AppendItalic(fmt.Sprintf("%s\n", e.Album.Text))
		}

		if a > int(limitTracks) {
			break
		}
	}

	return m
}

func fuckAbs(i int64) string {
	return strconv.FormatInt(i, 10)
}

//  func(cq *gotgbot.CallbackQuery)
func historyCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, historyPrefix)
}

// type Response func(b *gotgbot.Bot, ctx *ext.Context) error
func historyCallBackResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	mystrs := strings.Split(ctx.CallbackQuery.Data, "_")
	if len(mystrs) < 4 || mystrs[0] != historyCB {
		// not for history? seems impossible... just added it in
		// case so we can prevent from panics...
		return nil
	}

	// h_grpID_userUD_back
	thevalue := historyMap[mystrs[1]+"_"+mystrs[2]]

	// check if the value is already deleted from the memory or not.
	// this will reduce the high memory usage and will prevent the memory
	// from being fucked up.
	if thevalue == nil {
		_, err := ctx.EffectiveMessage.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{})
		if err != nil {
			return err
		}
		return nil
	} else if thevalue.owner.Id != ctx.EffectiveUser.Id {
		_, err := ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "this list isn't for ya",
			ShowAlert: true,
		})
		if err != nil {
			return err
		}
		return nil
	}

	if mystrs[3] == nextValue {
		// go to the next page

		// check if the current page is equal to the last page or not.
		// if yes, we should take the user to the first ('1') page.
		if thevalue.currentPage == thevalue.totalPages {
			thevalue.currentPage = 1
		} else {
			thevalue.currentPage++
		}
	} else {
		// go to the previous page

		// check if the current page is '1' or not.
		// if yes, we should take the user to the last page.
		if thevalue.currentPage == 1 {
			thevalue.currentPage = thevalue.totalPages
		} else {
			thevalue.currentPage--
		}
	}

	msg := strings.Replace(thevalue.GetParsedText().ToString(), "|", `\|`, -1)
	_, err := ctx.EffectiveMessage.EditText(b, msg,
		&gotgbot.EditMessageTextOpts{
			ParseMode:             "markdownv2",
			ReplyMarkup:           *ctx.EffectiveMessage.ReplyMarkup,
			DisableWebPagePreview: true,
		})
	if err != nil {
		return err
	}

	return nil
}

func (h *historyData) GetCurrentTrackList() []last_fm.Track {
	return h.tracks[h.currentPage-1]
}

func (h *historyData) SetAsCurrentTrackList(tracks []last_fm.Track) {
	if len(h.tracks) < h.totalPages {
		h.tracks = make([][]last_fm.Track, h.totalPages)
	}

	h.tracks[h.currentPage-1] = tracks
}

func (h *historyData) GenerateWholeList(tracks []last_fm.Track) {
	var currentList []last_fm.Track
	h.tracks = make([][]last_fm.Track, h.totalPages)
	num := 0
	index := 0
	//whole := 0
	for _, t := range tracks {
		currentList = append(currentList, t)
		num++
		if num >= int(limitTracks) {
			h.tracks[index] = currentList
			currentList = nil
			num = 0
			index++
			continue
		}
	}

	if currentList != nil {
		h.tracks[index] = currentList
	}
}

func (h *historyData) GetParsedText() mdparser.WMarkDown {
	index := h.currentPage - 1
	offset := int(index) * int(limitTracks)
	m := getSimpleList(h.tracks[index], h.owner, offset)

	m = m.AppendNormal(" ").AppendItalic("\n(Page " + strconv.Itoa(h.currentPage) +
		" from " + strconv.Itoa(h.totalPages) + ")")
	return m
}
