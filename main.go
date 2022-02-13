package main

import (
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var HandlerList []string

func main() {

	b, err := tb.NewBot(tb.Settings{
		Token:  TBToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	Config(DeleteDice, tb.OnDice)
	Config(DeleteUserJoined, tb.OnUserJoined)
	Config(DeleteUserLeft, tb.OnUserLeft)
	Config(DeleteDNewGroupTitle, tb.OnNewGroupTitle)
	Config(DeleteNewGroupPhoto, tb.OnNewGroupPhoto)
	Config(DeleteGroupPhotoDeleted, tb.OnGroupPhotoDeleted)
	Config(DeleteOnPinned, tb.OnPinned)

	log.Println(fmt.Sprintf("I will delete %s message.", HandlerList))

	dab := func(m *tb.Message, deleteChannel bool) {
		myRights, _ := b.ChatMemberOf(m.Chat, b.Me)
		if !myRights.Rights.CanDeleteMessages || !myRights.Rights.CanRestrictMembers {
			_, _ = b.Send(m.Chat, "爷权限不足，告辞！")
			_ = b.Leave(m.Chat)
			return
		}
		if deleteChannel && m.Sender.ID == 777000 {
			return
		}
		_ = b.Delete(m)
		if deleteChannel {
			_, _ = b.Raw("banChatSenderChat", map[string]int64{
				"chat_id":        m.Chat.ID,
				"sender_chat_id": m.SenderChat.ID,
			})
		}
	}

	for i := 0; i < len(HandlerList); i++ {
		OnEvent := HandlerList[i]
		b.Handle(OnEvent, func(m *tb.Message) {
			dab(m, DeleteChannel)
		})
	}

	if DeleteChannel {
		if b.Poller == nil {
			panic("telebot: can't start without a poller")
		}
		stop := make(chan struct{})
		go b.Poller.Poll(b, b.Updates, stop)

		for {
			upd := <-b.Updates
			if upd.Message != nil && upd.Message.SenderChat != nil {
				dab(upd.Message, DeleteChannel)
				continue
			}
			b.ProcessUpdate(upd)
		}
	} else {
		b.Start()
	}

}

func Config(check bool, handler string) {
	if check {
		HandlerList = append(HandlerList, handler)
	}
}
