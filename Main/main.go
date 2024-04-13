package main

import (
	"GoTarantool/Server"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
)

var MyUser string
var myPass string
var MyRole string
var SelectedGroupName string
var SelectedGroupId string
var userId string
var selectedRowMsg *gtk.ListBoxRow
var selectedRowGroup *gtk.ListBoxRow
var selectedRowUsers *gtk.ListBoxRow

type TimedMsg struct {
	msg  string
	time uint64
}

type MessageStruct struct {
	message string
	user    string
	msgTime uint64
	msgId   string
}

var messagesArr = make([]TimedMsg, 0)

func main() {
	t := time.NewTimer(1 * time.Second)
	//isLogin := false
	conn := Server.Server()
	defer conn.Close()

	gtk.Init(nil)

	// Создаём билдер
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = b.AddFromFile("Chat_gui.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID

	//Объекты

	objMain, err := b.GetObject("main_window")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	objReg, err := b.GetObject("regWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objChangeGroup, err := b.GetObject("changeWinGroup")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objChangeMsg, err := b.GetObject("changeWinMsg")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objNewGroupWin, err := b.GetObject("newGroupWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objJoinGroupWin, err := b.GetObject("joinGroupWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objWarningWin, err := b.GetObject("warningWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objUsersWin, err := b.GetObject("usersWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	winMain := objMain.(*gtk.Window)
	//winMain.Move(0, 0)

	winMain.Connect("destroy", func() {
		gtk.MainQuit()
	})

	winReg := objReg.(*gtk.Window)
	winReg.Connect("delete-event", func() {
		gtk.MainQuit()

	})

	winChangeGroup := objChangeGroup.(*gtk.Dialog)
	winChangeGroup.Connect("delete-event", func() {
		winChangeGroup.Hide()

	})

	winChangeMsg := objChangeMsg.(*gtk.Dialog)
	winChangeMsg.Connect("delete-event", func() {
		winChangeMsg.Hide()

	})

	winNewGroup := objNewGroupWin.(*gtk.Dialog)

	winJoinGroup := objJoinGroupWin.(*gtk.Dialog)

	warningWin := objWarningWin.(*gtk.Dialog)
	warningWin.Connect("delete-event", func() {
		warningWin.Hide()

	})

	usersWin := objUsersWin.(*gtk.Dialog)

	objMain, _ = b.GetObject("msg_entry")
	msgEntry := objMain.(*gtk.Entry)

	// Получаем кнопку

	objMain, _ = b.GetObject("msg_btn")
	msgBtn := objMain.(*gtk.Button)

	// Получаем метку

	objMain, _ = b.GetObject("guild_label")
	guildLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("tag_label")
	tagLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("group_text_label")
	guildTextLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("tag_text_label")
	tagTextLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("user_label")
	userLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("registration_success_label")
	registrationSuccessLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("labelCheckId")
	labelCheckId := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("msg_scroll")
	scrolledWindow := objMain.(*gtk.ScrolledWindow)

	objReg, _ = b.GetObject("login_reg_btn")
	loginRegBtn := objReg.(*gtk.Button)

	objChangeMsg, _ = b.GetObject("changeBtnMsg")
	changeBtnMsgConfirm := objChangeMsg.(*gtk.Button)

	objChangeGroup, _ = b.GetObject("changeBtnGroup")
	changeBtnGroup := objChangeGroup.(*gtk.Button)

	objChangeGroup, _ = b.GetObject("changeGroupEntry")
	changeGroupEntry := objChangeGroup.(*gtk.Entry)

	objMain, _ = b.GetObject("btn_change_msg")
	btnChangeMsg := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("newuser_reg_btn")
	newUserRegBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("login_reg_entry")
	loginRegEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("pass_reg_entry")
	passRegEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("joinGroupEntry")
	joinGroupEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("listbox_groups")
	groupsListbox := objMain.(*gtk.ListBox)

	objNewGroupWin, _ = b.GetObject("entryNewGroupName")
	entryNewGroupName := objNewGroupWin.(*gtk.Entry)

	objNewGroupWin, _ = b.GetObject("entryNewGroupId")
	entryNewGroupId := objNewGroupWin.(*gtk.Entry)

	objChangeMsg, _ = b.GetObject("changeMsgEntry")
	changeMsgEntry := objChangeMsg.(*gtk.Entry)

	objMain, _ = b.GetObject("add_group_btn")
	addGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("btn_del_group")
	delGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("btn_change_group")
	changeGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("exit_btn")
	exitBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("listbox_msg")
	msgListbox := objMain.(*gtk.ListBox)

	objUsersWin, _ = b.GetObject("listbox_users")
	usersListbox := objUsersWin.(*gtk.ListBox)

	objNewGroupWin, _ = b.GetObject("newGroupBtn")
	newGroupBtn := objNewGroupWin.(*gtk.Button)

	objNewGroupWin, _ = b.GetObject("btnCheckId")
	btnCheckId := objNewGroupWin.(*gtk.Button)

	objMain, _ = b.GetObject("join_group_btn")
	joinGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("joinGroupBtnConfirm")
	joinGroupBtnConfirm := objMain.(*gtk.Button)

	objWarningWin, _ = b.GetObject("warningWinBtn")
	warningWinBtn := objWarningWin.(*gtk.Button)

	objMain, _ = b.GetObject("btn_del_msg")
	btnDelMsg := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("users_btn")
	usersBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("statistics_btn")
	statisticsBtn := objMain.(*gtk.Button)

	objUsersWin, _ = b.GetObject("kick_btn")
	kickBtn := objUsersWin.(*gtk.Button)

	objUsersWin, _ = b.GetObject("promoteBtn")
	promoteBtn := objUsersWin.(*gtk.Button)

	objUsersWin, _ = b.GetObject("downgradeBtn")
	downgradeBtn := objUsersWin.(*gtk.Button)

	//objUsersWin, _ = b.GetObject("test_entry")
	//testEntry := objUsersWin.(*gtk.Entry)

	usersWin.Connect("delete-event", func() {
		selectedRowUsers, _ = gtk.ListBoxRowNew()
		selectedRowUsers = nil

		fmt.Println(selectedRowGroup.GetName())
		fmt.Println(MyUser)
		fmt.Println(MyRole)

		usersWin.Hide()

		selectedRowMsg, _ = gtk.ListBoxRowNew() ///////////
		selectedRowMsg = nil
		selectedRowGroup = groupsListbox.GetSelectedRow()
		selectedRowGroupId, _ := selectedRowGroup.GetName()
		fmt.Println("group id ", selectedRowGroupId)

		labelRow, _ := selectedRowGroup.GetChild()
		groupLabel := labelRow.(*gtk.Label)
		textLabel, _ := groupLabel.GetText()
		SelectedGroupName = textLabel
		SelectedGroupId = selectedRowGroupId
		guildLabel.SetText(SelectedGroupName)
		tagLabel.SetText("@" + selectedRowGroupId)
		messagesArr = messagesArr[:0]
		GetMsg(conn, msgListbox)
		AutoScroll(scrolledWindow)

		info, _ := conn.Call("fn.group_users_cnt", []interface{}{SelectedGroupId})
		usersCntTuples := info.Tuples()
		usersCnt := usersCntTuples[0][0].(string)
		usersBtn.SetLabel("Участники: " + usersCnt)

		info, _ = conn.Call("fn.get_user_role", []interface{}{MyUser, SelectedGroupId})
		userRoleTuples := info.Tuples()
		MyRole = userRoleTuples[0][0].(string)

		//Показываем элементы группы
		btnDelMsg.Show()
		btnChangeMsg.Show()
		usersBtn.Show()
		statisticsBtn.Show()
		changeGroupBtn.Show()
		delGroupBtn.Show()
		msgBtn.Show()
		msgEntry.Show()
		guildLabel.Show()
		guildTextLabel.Show()
		tagLabel.Show()
		tagTextLabel.Show()

		winMain.ShowAll()

		//Таймер потом включить
		// go func() {
		// 	for {
		// 		t.Reset(2 * time.Second)
		// 		fmt.Println("Показываем в таймере")
		// 		info, _ = conn.Call("fn.group_exists", []interface{}{MyUser, selectedRowGroupId})
		// 		fmt.Println("Проверка")
		// 		groupExistsTuples := info.Tuples()
		// 		fmt.Println("туплес")
		// 		groupExists := groupExistsTuples[0][0].(bool)
		// 		if groupExists == false {
		// 			fmt.Println("НЕ СУЩЕСТВУЕТ")
		// 			t.Stop()
		// 			// groupsListbox.UnselectAll()
		// 			// msgListbox.UnselectAll()
		// 			// clearListbox(groupsListbox)
		// 			// clearListbox(msgListbox)
		// 			// selectedRowGroup, _ = gtk.ListBoxRowNew()
		// 			// selectedRowGroup = nil
		// 			// selectedRowMsg, _ = gtk.ListBoxRowNew()
		// 			// selectedRowMsg = nil
		// 			// userLabel.SetText("")
		// 			// guildLabel.SetText("")
		// 			// tagLabel.SetText("")
		// 			// MyUser = ""
		// 			// myPass = ""
		// 			// MyRole = ""
		// 			// messagesArr = messagesArr[:0]
		// 			// fmt.Println("Отключен")
		// 			break
		// 		}
		// 		if groupExists == true {
		// 			GetMsg(conn, msgListbox)
		// 		}
		// 		<-t.C
		// 	}

		// }()

	})

	winNewGroup.Connect("delete-event", func() {
		entryNewGroupId.SetText("")
		entryNewGroupName.SetText("")
		labelCheckId.SetText("")
		winNewGroup.Hide()

	})

	winJoinGroup.Connect("delete-event", func() {
		joinGroupEntry.SetText("")
		winJoinGroup.Hide()
	})

	groupsListbox.Connect("button-press-event", func(box *gtk.ListBox, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS || buttonEvent.Type() == gdk.EVENT_3BUTTON_PRESS {
			fmt.Println("DOUBLE")
			groupsListbox.UnselectAll()
			btnDelMsg.Hide()
			btnChangeMsg.Hide()
			usersBtn.Hide()
			statisticsBtn.Hide()
			changeGroupBtn.Hide()
			delGroupBtn.Hide()
			groupsListbox.Hide()
			groupsListbox.Show()
			selectedRowGroup = groupsListbox.GetSelectedRow()
			clearListbox(groupsListbox)
			clearListbox(msgListbox)
			GetGroups(conn, groupsListbox)
			guildLabel.SetText("")
			tagLabel.SetText("")
			usersBtn.SetLabel("Участники: 0")
			selectedRowGroup, _ = gtk.ListBoxRowNew()
			selectedRowGroup = nil
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
			fmt.Println("двойное", selectedRowGroup.GetIndex())
			return // Ignore double-click events
		}

		// Handle single-click events
		row := box.GetSelectedRow()
		if row != nil {
			index := row.GetIndex()
			fmt.Printf("Clicked on row %d\n", index)
		}
	})

	msgListbox.Connect("row-activated", func() {

		tempRow := msgListbox.GetSelectedRow()
		tempName, _ := tempRow.GetName()
		if tempName == "system" {
			msgListbox.UnselectAll()
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
			return
		}
		fmt.Println("tempRow.GetIndex() msg:", tempRow.GetIndex())
		fmt.Println("selectedRow.GetIndex()) msg:", selectedRowMsg.GetIndex())

		fmt.Println(tempRow.GetName())

		if tempRow.GetIndex() == selectedRowMsg.GetIndex() {
			msgListbox.UnselectAll()
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
		} else {
			selectedRowMsg = msgListbox.GetSelectedRow()
		}

		//widgetRow, _ := (selectedRowMsg.GetChild())
		//labelRow := widgetRow.(*gtk.Label)
		//fmt.Println(labelRow.GetText())
	})

	btnChangeMsg.Connect("clicked", func() {
		if selectedRowMsg != nil {
			msgId, _ := selectedRowMsg.GetName()
			info, _ := conn.Call("fn.get_selected_msg", []interface{}{msgId})
			msgTuples := info.Tuples()
			msgUser := msgTuples[0][0].(string)
			newMsgText := msgTuples[2][0].(string)
			if msgUser == MyUser {
				changeMsgEntry.SetText(newMsgText)
				winChangeMsg.Run()
			} else {
				warningWin.Run()
				fmt.Println("Вы можете изменять только свои сообщения!")
			}
		}
	})

	changeBtnMsgConfirm.Connect("clicked", func() {
		fmt.Println(selectedRowMsg.GetName())
		fmt.Println(changeMsgEntry.GetText())
		msgId, _ := selectedRowMsg.GetName()
		msgText, _ := changeMsgEntry.GetText()
		info, _ := conn.Call("fn.edit_msg", []interface{}{msgId, msgText})
		msgTuples := info.Tuples()
		fmt.Println(msgTuples)
		newMsgUser := msgTuples[0][0].(string)
		newMsgGroup := msgTuples[1][0].(string)
		newMsgText := msgTuples[2][0].(string)
		row, _ := selectedRowMsg.GetChild()
		rowLabel := row.(*gtk.Label)
		rowLabel.SetText(newMsgUser + "(" + newMsgGroup + "): " + newMsgText)

		//TODO автозаполнение поля изменения сообщения,изменение только своего сообщения
		winChangeMsg.Hide()
	})

	changeBtnGroup.Connect("clicked", func() {
		newGroup, _ := changeGroupEntry.GetText()
		_, _ = conn.Call("fn.edit_group", []interface{}{MyUser, SelectedGroupId, newGroup})
		GetGroups(conn, groupsListbox)
		selectedRowGroup, _ = gtk.ListBoxRowNew()
		selectedRowGroup = nil
		SelectedGroupName = newGroup
		clearListbox(msgListbox)
		messagesArr = messagesArr[:0]
		msgEntry.SetText("")
		winChangeGroup.Hide()
	})

	usersListbox.Connect("row-activated", func() {
		fmt.Println(selectedRowUsers.GetIndex())
		tempRow := usersListbox.GetSelectedRow()
		if tempRow.GetIndex() == selectedRowUsers.GetIndex() {
			usersListbox.UnselectAll()
			selectedRowUsers, _ = gtk.ListBoxRowNew()
			selectedRowUsers = nil
		} else {
			selectedRowUsers = usersListbox.GetSelectedRow()
			labelRow, _ := selectedRowUsers.GetChild()
			userLabel := labelRow.(*gtk.Label)
			userName, _ := userLabel.GetText()
			userRole, _ := selectedRowUsers.GetName()
			fmt.Println(userName)
			fmt.Println(userRole)

		}
	})

	groupsListbox.Connect("row-activated", func() {
		fmt.Println(selectedRowGroup)
		fmt.Println("Выбрана группа")
		AutoScroll(scrolledWindow)
		tempRow := groupsListbox.GetSelectedRow()
		fmt.Println("tempRow.GetIndex() group:", tempRow.GetIndex())
		fmt.Println("selectedRow.GetIndex()) group:", selectedRowGroup.GetIndex())
		if tempRow.GetIndex() == selectedRowGroup.GetIndex() {
			t.Stop()
			groupsListbox.UnselectAll()
			guildLabel.SetText("")
			tagLabel.SetText("")
			clearListbox(msgListbox)
			selectedRowGroup, _ = gtk.ListBoxRowNew()
			selectedRowGroup = nil
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
			MyRole = ""
			fmt.Println("Индекс сообщения", selectedRowMsg.GetIndex())
			fmt.Println("sel gr после обнуления", selectedRowGroup)
			usersBtn.SetLabel("Участники: 0")

			//Скрываем элементы группы
			btnDelMsg.Hide()
			btnChangeMsg.Hide()
			usersBtn.Hide()
			statisticsBtn.Hide()
			changeGroupBtn.Hide()
			delGroupBtn.Hide()
			msgBtn.Hide()
			msgEntry.Hide()
			guildLabel.Hide()
			guildTextLabel.Hide()
			tagLabel.Hide()
			tagTextLabel.Hide()

		} else {
			selectedRowMsg, _ = gtk.ListBoxRowNew() ///////////
			selectedRowMsg = nil
			selectedRowGroup = groupsListbox.GetSelectedRow()
			selectedRowGroupId, _ := selectedRowGroup.GetName()
			fmt.Println("group id ", selectedRowGroupId)

			labelRow, _ := selectedRowGroup.GetChild()
			groupLabel := labelRow.(*gtk.Label)
			textLabel, _ := groupLabel.GetText()
			SelectedGroupName = textLabel
			SelectedGroupId = selectedRowGroupId
			guildLabel.SetText(SelectedGroupName)
			tagLabel.SetText("@" + selectedRowGroupId)
			messagesArr = messagesArr[:0]
			GetMsg(conn, msgListbox)
			AutoScroll(scrolledWindow)

			info, _ := conn.Call("fn.group_users_cnt", []interface{}{SelectedGroupId})
			usersCntTuples := info.Tuples()
			usersCnt := usersCntTuples[0][0].(string)
			usersBtn.SetLabel("Участники: " + usersCnt)

			info, _ = conn.Call("fn.get_user_role", []interface{}{MyUser, SelectedGroupId})
			userRoleTuples := info.Tuples()
			MyRole = userRoleTuples[0][0].(string)

			//Показываем элементы группы
			btnDelMsg.Show()
			btnChangeMsg.Show()
			usersBtn.Show()
			statisticsBtn.Show()
			changeGroupBtn.Show()
			delGroupBtn.Show()
			msgBtn.Show()
			msgEntry.Show()
			guildLabel.Show()
			guildTextLabel.Show()
			tagLabel.Show()
			tagTextLabel.Show()

			winMain.ShowAll()

			//Таймер потом включить
			// go func() {
			// 	for {
			// 		t.Reset(2 * time.Second)
			// 		GetMsg(conn, msgListbox)
			// 		fmt.Println("Показываем в таймере")
			// 		info, _ = conn.Call("fn.group_exists", []interface{}{MyUser, selectedRowGroupId})
			// 		fmt.Println("Проверка")
			// 		groupExistsTuples := info.Tuples()
			// 		fmt.Println("туплес")
			// 		groupExists := groupExistsTuples[0][0].(bool)
			// 		if groupExists == false {
			// 			t.Stop()
			// 			// // fmt.Println("Не существует")
			// 			// // groupsListbox.UnselectAll()
			// 			// // msgListbox.UnselectAll()
			// 			// clearListbox(groupsListbox)
			// 			// GetGroups(conn, groupsListbox)
			// 			// clearListbox(msgListbox)
			// 			// selectedRowGroup, _ = gtk.ListBoxRowNew()
			// 			// selectedRowGroup = nil
			// 			// selectedRowMsg, _ = gtk.ListBoxRowNew()
			// 			// selectedRowMsg = nil
			// 			// userLabel.SetText("")
			// 			// guildLabel.SetText("")
			// 			// tagLabel.SetText("")
			// 			// // MyUser = ""
			// 			// // myPass = ""
			// 			// // MyRole = ""
			// 			// messagesArr = messagesArr[:0]
			// 			// break

			// 			groupsListbox.UnselectAll()
			// 			msgListbox.UnselectAll()
			// 			clearListbox(groupsListbox)
			// 			clearListbox(msgListbox)
			// 			selectedRowGroup, _ = gtk.ListBoxRowNew()
			// 			selectedRowGroup = nil
			// 			selectedRowMsg, _ = gtk.ListBoxRowNew()
			// 			selectedRowMsg = nil
			// 			userLabel.SetText("")
			// 			guildLabel.SetText("")
			// 			tagLabel.SetText("")
			// 			MyUser = ""
			// 			myPass = ""
			// 			MyRole = ""
			// 			messagesArr = messagesArr[:0]
			// 			fmt.Println("Выход")
			// 			winReg.ShowAll()
			// 			winMain.Hide()
			// 		}
			// 		<-t.C
			// 	}

			// }()



			go func() {
				for {
					t.Reset(2 * time.Second)
					GetMsg(conn, msgListbox)
					fmt.Println("Показываем в таймере")
					info, _ = conn.Call("fn.group_exists", []interface{}{MyUser, selectedRowGroupId})
					fmt.Println("Проверка")
					groupExistsTuples := info.Tuples()
					fmt.Println("туплес")
					groupExists := groupExistsTuples[0][0].(bool)
					if groupExists == false {
						t.Stop()
						// // fmt.Println("Не существует")
						// // groupsListbox.UnselectAll()
						// // msgListbox.UnselectAll()
						// clearListbox(groupsListbox)
						// GetGroups(conn, groupsListbox)
						// clearListbox(msgListbox)
						// selectedRowGroup, _ = gtk.ListBoxRowNew()
						// selectedRowGroup = nil
						// selectedRowMsg, _ = gtk.ListBoxRowNew()
						// selectedRowMsg = nil
						// userLabel.SetText("")
						// guildLabel.SetText("")
						// tagLabel.SetText("")
						// // MyUser = ""
						// // myPass = ""
						// // MyRole = ""
						// messagesArr = messagesArr[:0]
						// break

						groupsListbox.UnselectAll()
						msgListbox.UnselectAll()
						clearListbox(groupsListbox)
						clearListbox(msgListbox)
						selectedRowGroup, _ = gtk.ListBoxRowNew()
						selectedRowGroup = nil
						selectedRowMsg, _ = gtk.ListBoxRowNew()
						selectedRowMsg = nil
						userLabel.SetText("")
						guildLabel.SetText("")
						tagLabel.SetText("")
						MyUser = ""
						myPass = ""
						MyRole = ""
						messagesArr = messagesArr[:0]
						fmt.Println("Выход")
						winReg.ShowAll()
						winMain.Hide()
					}
					<-t.C
				}

			}()
		}

	})

	// Сигнал по нажатию на кнопку

	exitBtn.Connect("clicked", func() {
		groupsListbox.UnselectAll()
		msgListbox.UnselectAll()
		clearListbox(groupsListbox)
		clearListbox(msgListbox)
		selectedRowGroup, _ = gtk.ListBoxRowNew()
		selectedRowGroup = nil
		selectedRowMsg, _ = gtk.ListBoxRowNew()
		selectedRowMsg = nil
		userLabel.SetText("")
		guildLabel.SetText("")
		tagLabel.SetText("")
		MyUser = ""
		myPass = ""
		MyRole = ""
		messagesArr = messagesArr[:0]
		fmt.Println("Выход")
		winReg.ShowAll()
		winMain.Hide()
	})

	newGroupBtn.Connect("clicked", func() {
		groupName, _ := entryNewGroupName.GetText()
		groupId, _ := entryNewGroupId.GetText()

		info, _ := conn.Call("fn.check_group_id", []interface{}{groupId})
		tuples := info.Tuples()
		checked := tuples[0][0].(string)
		if checked == "true" {
			if MyUser != "" && !isOnlyWhitespace(groupName) && !isOnlyWhitespace(groupId) {
				fmt.Println(MyUser, groupName, groupId)
				_, _ = conn.Call("fn.new_group", []interface{}{MyUser, groupId, groupName})

				rowGroup, _ := gtk.ListBoxRowNew()
				rowGroup.SetName(groupId)
				labelGroup, _ := gtk.LabelNew(groupName)
				labelGroup.SetSizeRequest(-1, 50)
				markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", groupName)
				labelGroup.SetMarkup(markup)
				rowGroup.Add(labelGroup)
				groupsListbox.Insert(rowGroup, 0)

				entryNewGroupId.SetText("")
				entryNewGroupName.SetText("")
				labelCheckId.SetText("")

				winNewGroup.Close()
				winMain.ShowAll()
			}
		}
		if checked == "false" {
			checkedText := "Тег занят"
			markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}

	})

	btnCheckId.Connect("clicked", func() {
		groupId, _ := entryNewGroupId.GetText()
		info, _ := conn.Call("fn.check_group_id", []interface{}{groupId})
		tuples := info.Tuples()
		checked := tuples[0][0].(string)
		if checked == "true" {
			checkedText := "Тег свободен"
			markup := fmt.Sprintf("<span size='15000' foreground='green'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}
		if checked == "false" {
			checkedText := "Тег занят"
			markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}
	})

	delGroupBtn.Connect("clicked", func() {
		if selectedRowGroup != nil {
			message := fmt.Sprintf("Пользователь %s вышел из группы", MyUser)
			_, _ = conn.Call("fn.new_msg", []interface{}{message, SelectedGroupId, "system"})
			_, _ = conn.Call("fn.del_group", []interface{}{MyUser, SelectedGroupId})
			groupsListbox.Remove(selectedRowGroup)
			selectedRowGroup, _ = gtk.ListBoxRowNew()
			selectedRowGroup = nil
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
			clearListbox(msgListbox)

			guildLabel.SetText("")
			tagLabel.SetText("")
			MyRole = ""
			messagesArr = messagesArr[:0]
			msgEntry.SetText("")
		}
	})

	changeGroupBtn.Connect("clicked", func() {
		fmt.Println(selectedRowGroup)
		fmt.Println("......")
		if selectedRowGroup != nil {
			groupRow, _ := selectedRowGroup.GetChild()
			groupLabel := groupRow.(*gtk.Label)
			groupText, _ := groupLabel.GetText()
			changeGroupEntry.SetText(groupText)
			winChangeGroup.Run()
		}

	})

	loginRegBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			newUser, _ := loginRegEntry.GetText()
			newPass, _ := passRegEntry.GetText()

			clearListbox(groupsListbox)

			info, _ := conn.Call("fn.login", []interface{}{newUser, newPass})
			userTuples := info.Tuples()
			if userTuples[0][0].(bool) == true {
				MyUser = newUser
				myPass = newPass

				userLabel.SetText(MyUser)
				GetGroups(conn, groupsListbox)
				loginRegEntry.SetText("")
				passRegEntry.SetText("")
				registrationSuccessLabel.SetText("\n")
				winMain.ShowAll()
				winReg.Hide()

				//Скрываем элементы группы
				btnDelMsg.Hide()
				btnChangeMsg.Hide()
				usersBtn.Hide()
				statisticsBtn.Hide()
				changeGroupBtn.Hide()
				delGroupBtn.Hide()
				msgBtn.Hide()
				msgEntry.Hide()
				guildLabel.Hide()
				guildTextLabel.Hide()
				tagLabel.Hide()
				tagTextLabel.Hide()

			} else {
				errText := "Ошибка! \n Неверный логин или пароль!"
				markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", errText)
				registrationSuccessLabel.SetText(errText)
				registrationSuccessLabel.SetMarkup(markup)
			}

		}

	})

	promoteBtn.Connect("clicked", func() {
		if selectedRowUsers != nil {
			selectedRowUsers = usersListbox.GetSelectedRow()
			labelRow, _ := selectedRowUsers.GetChild()
			userLabel := labelRow.(*gtk.Label)
			userName, _ := userLabel.GetText()
			userRole, _ := selectedRowUsers.GetName()
			fmt.Println(userName)
			fmt.Println(userRole)
			if MyRole == "admin" {
				message := fmt.Sprintf("Пользователь %s повышен до администратора", userName)
				_, _ = conn.Call("fn.new_msg", []interface{}{message, SelectedGroupId, "system"})
				_, _ = conn.Call("fn.promote_user", []interface{}{userName, SelectedGroupId, "admin"})
				usersListbox.Remove(selectedRowUsers)
				selectedRowUsers, _ = gtk.ListBoxRowNew()
				selectedRowUsers = nil
				GetMsg(conn, msgListbox)
				clearListbox(usersListbox)
				GetUsers(conn, usersListbox)

			} else if MyRole != "admin" {
				fmt.Println("Недостаточно прав")
			} else if MyUser == userName {
				fmt.Println("Вы уже являетесь администратором")
			}

			// clearListbox(msgListbox)

			// guildLabel.SetText("")
			// tagLabel.SetText("")
			// messagesArr = messagesArr[:0]
			// msgEntry.SetText("")
		}

	})

	downgradeBtn.Connect("clicked", func() {
		if selectedRowUsers != nil {
			selectedRowUsers = usersListbox.GetSelectedRow()
			labelRow, _ := selectedRowUsers.GetChild()
			userLabel := labelRow.(*gtk.Label)
			userName, _ := userLabel.GetText()
			userRole, _ := selectedRowUsers.GetName()
			userName = getTextBeforeSubstring(userName, " (")
			fmt.Println(userName)
			fmt.Println(userRole)
			if MyUser == userName {
				fmt.Println("Вы не можете понизить свой ранг")
			} else if MyRole != "admin" {
				fmt.Println("Недостаточно прав")
			} else if MyRole == "admin" {
				message := fmt.Sprintf("%s понижен до пользователя", userName)
				_, _ = conn.Call("fn.new_msg", []interface{}{message, SelectedGroupId, "system"})
				_, _ = conn.Call("fn.promote_user", []interface{}{userName, SelectedGroupId, "user"})
				usersListbox.Remove(selectedRowUsers)
				selectedRowUsers, _ = gtk.ListBoxRowNew()
				selectedRowUsers = nil
				GetMsg(conn, msgListbox)
				clearListbox(usersListbox)
				GetUsers(conn, usersListbox)

			}

			// clearListbox(msgListbox)

			// guildLabel.SetText("")
			// tagLabel.SetText("")
			// messagesArr = messagesArr[:0]
			// msgEntry.SetText("")
		}

	})

	addGroupBtn.Connect("clicked", func() {
		winNewGroup.Run()

	})

	joinGroupBtn.Connect("clicked", func() {
		winJoinGroup.Run()
	})

	warningWinBtn.Connect("clicked", func() {
		warningWin.Hide()
	})

	usersBtn.Connect("clicked", func() {
		clearListbox(usersListbox)
		res := GetUsers(conn, usersListbox)
		if res {
			if MyRole == "admin" {
				kickBtn.Show()
				promoteBtn.Show()
				downgradeBtn.Show()
			} else {
				kickBtn.Hide()
				promoteBtn.Hide()
				downgradeBtn.Hide()
			}
			t.Stop()
			usersWin.Run()
		}

	})

	kickBtn.Connect("clicked", func() {
		if selectedRowUsers != nil {
			selectedRowUsers = usersListbox.GetSelectedRow()
			labelRow, _ := selectedRowUsers.GetChild()
			userLabel := labelRow.(*gtk.Label)
			userName, _ := userLabel.GetText()
			userName = getTextBeforeSubstring(userName, " (")
			userRole, _ := selectedRowUsers.GetName()
			fmt.Println(userName)
			fmt.Println(userRole)
			if userName == MyUser {
				fmt.Println("Вы не можете удалить сами себя")
			} else if MyRole != "admin" {
				fmt.Println("Недостаточно прав")
			} else if MyRole == "admin" {
				message := fmt.Sprintf("Пользователь %s исключён из группы", userName)
				_, _ = conn.Call("fn.new_msg", []interface{}{message, SelectedGroupId, "system"})
				_, _ = conn.Call("fn.del_group", []interface{}{userName, SelectedGroupId})
				usersListbox.Remove(selectedRowUsers)
				selectedRowUsers, _ = gtk.ListBoxRowNew()
				selectedRowUsers = nil
				GetMsg(conn, msgListbox)

				//убрать исключение самого себя

				//У исключённого не забыть TODO
				// selectedRowGroup, _ = gtk.ListBoxRowNew()
				// selectedRowGroup = nil
				// selectedRowMsg, _ = gtk.ListBoxRowNew()
				// selectedRowMsg = nil
			}
			// clearListbox(msgListbox)

			// guildLabel.SetText("")
			// tagLabel.SetText("")
			// messagesArr = messagesArr[:0]
			// msgEntry.SetText("")
		}
	})

	btnDelMsg.Connect("clicked", func() {
		if selectedRowMsg != nil {
			msgId, _ := selectedRowMsg.GetName()
			_, _ = conn.Call("fn.del_msg", []interface{}{msgId})
			msgListbox.Remove(selectedRowMsg)
			selectedRowMsg, _ = gtk.ListBoxRowNew()
			selectedRowMsg = nil
		}
	})

	joinGroupBtnConfirm.Connect("clicked", func() {
		groupId, _ := joinGroupEntry.GetText()
		info, _ := conn.Call("fn.join_group", []interface{}{MyUser, groupId})
		fmt.Println(info)
		tuples := info.Tuples()
		fmt.Println(tuples[0][0])
		if tuples[0][0].(string) != "false" {
			joinGroupName := tuples[0][0].(string)

			rowGroup, _ := gtk.ListBoxRowNew()
			rowGroup.SetName(groupId)
			labelGroup, _ := gtk.LabelNew(joinGroupName)
			labelGroup.SetSizeRequest(-1, 50)
			markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", joinGroupName)
			labelGroup.SetMarkup(markup)
			rowGroup.Add(labelGroup)
			groupsListbox.Insert(rowGroup, 0)

			// rowJoinUser, _ := gtk.ListBoxRowNew()
			// rowJoinUser.SetName("notMsg")
			// labelJoinUser, _ := gtk.LabelNew(message)
			// rowJoinUser.Add(labelJoinUser)
			// msgListbox.Insert(rowJoinUser, 0)
			message := fmt.Sprintf("Пользователь %s вступил в группу", MyUser)
			_, _ = conn.Call("fn.new_msg", []interface{}{message, groupId, "system"})

			joinGroupEntry.SetText("")
			winJoinGroup.Hide()
			winMain.ShowAll()
		}
	})

	msgBtn.Connect("button-press-event", func(btn *gtk.Button, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS || buttonEvent.Type() == gdk.EVENT_3BUTTON_PRESS {
			newMsg, _ := msgEntry.GetText()
			_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupId, MyUser})
			GetMsg(conn, msgListbox)
			AutoScroll(scrolledWindow)
			//msgListbox.ShowAll()
			//winMain.ShowAll()
			fmt.Println("Double click")
			return // Ignore double-click events
		}

		newMsg, _ := msgEntry.GetText()
		_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupId, MyUser})
		GetMsg(conn, msgListbox)
		AutoScroll(scrolledWindow)
		msgEntry.SetText("")

	})

	msgEntry.Connect("key-press-event", func(entry *gtk.Entry, event *gdk.Event) {
		keyEvent := &gdk.EventKey{Event: event}
		keyVal := keyEvent.KeyVal()
		if keyVal == gdk.KEY_Return {
			text, _ := entry.GetText()
			log.Println("Enter key pressed in entry. Text entered:", text)
			newMsg, _ := msgEntry.GetText()
			_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupId, MyUser})
			GetMsg(conn, msgListbox)
			AutoScroll(scrolledWindow)
			msgEntry.SetText("")
		}
	})

	newUserRegBtn.Connect("clicked", func() {
		newLogin, _ := loginRegEntry.GetText()
		newPass, _ := passRegEntry.GetText()
		if newLogin != "" && newPass != "" {
			info, _ := conn.Call("fn.new_user", []interface{}{newLogin, newPass})
			successTuples := info.Tuples()
			success := successTuples[0][0].(bool)
			successText := ""
			markup := ""
			if success == true {
				fmt.Println("Успешная регистрация!")
				successText = "Успешная регистрация! \n"
				markup = fmt.Sprintf("<span size='15000' foreground='green'>%s</span>", successText)

			} else {
				successText = "Ошибка! \n Пользователь уже существует"
				markup = fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", successText)
			} // TODO неверный пароль
			registrationSuccessLabel.SetText(successText)
			registrationSuccessLabel.SetMarkup(markup)
		}
	})

	//})

	// Отображаем все виджеты в окне

	// provider, _ := gtk.CssProviderNew()
	// provider.LoadFromData(`
	// 	window {
	// 		background-color: rgba(255, 255, 255, 0.5); /* Белый цвет фона с прозрачностью 50% */
	// 	}
	// `)

	// context, _ := winReg.GetStyleContext()
	// context.AddProvider(provider, uint(gtk.STYLE_PROVIDER_PRIORITY_APPLICATION))

	winReg.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

}

func AutoScroll(scrolledWindow *gtk.ScrolledWindow) {
	//scrolledWindow.ShowAll()
	adjustment := scrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetLower())
	// scrolledWindow.ShowAll()
	// adjustment := scrolledWindow.GetVAdjustment()
	// tempUp := adjustment.GetUpper()
	// adjustment.SetUpper(adjustment.GetUpper() + 22)
	// adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
	// fmt.Println("Up до", adjustment.GetUpper())
	// fmt.Println("Page до", adjustment.GetPageSize())
	// fmt.Println("value до", adjustment.GetUpper() - adjustment.GetPageSize())

	// scrolledWindow.ShowAll()
	// adjustment.SetUpper(tempUp)
	// adjustment.SetValue(adjustment.GetValue() - 22)
	// fmt.Println("Up после", adjustment.GetUpper())
	// fmt.Println("Page после", adjustment.GetPageSize())
}

// func delGroup(groupId int, groupsListbox *gtk.ListBox) {
// 	groupsListbox
// }

func clearListbox(ListBox *gtk.ListBox) {
	messagesArr = messagesArr[:0] //TODO ПЕРЕНЕСТИ ИЗ ФУНКЦИИ

	children := ListBox.GetChildren()

	for children.Length() > 0 {
		child := children.NthData(0)
		if widget, ok := child.(*gtk.Widget); ok {
			ListBox.Remove(widget)
		}
		children = ListBox.GetChildren()
	}
	ListBox.ShowAll()
}

func isOnlyWhitespace(input string) bool {
	for _, char := range input {
		if !unicode.IsSpace(char) { // Проверяем, что символ не является пробелом
			return false
		}
	}
	return true
}

func GetGroups(conn *tarantool.Connection, groupsListbox *gtk.ListBox) {
	clearListbox(groupsListbox)
	infoUserGroups, _ := conn.Call("fn.get_user_groups", []interface{}{MyUser})
	userGroupsTuples := infoUserGroups.Tuples()
	fmt.Println(userGroupsTuples)
	fmt.Println(len(userGroupsTuples))
	fmt.Println(len(userGroupsTuples[0]))
	if len(userGroupsTuples) == 1 && len(userGroupsTuples[0]) == 0 {
		return
	}
	for i := 0; i < len(userGroupsTuples); i++ {
		rowGroup, _ := gtk.ListBoxRowNew()
		rowGroup.SetName(userGroupsTuples[i][0].(string))
		labelGroup, _ := gtk.LabelNew(userGroupsTuples[i][1].(string))
		labelGroup.SetSizeRequest(-1, 50)
		tempText, _ := labelGroup.GetText()
		markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", tempText)
		labelGroup.SetMarkup(markup)
		rowGroup.Add(labelGroup)
		groupsListbox.Insert(rowGroup, 0)

	}
	groupsListbox.ShowAll()
}

func GetUsers(conn *tarantool.Connection, usersListbox *gtk.ListBox) bool {
	usersListbox.UnselectAll()
	info, _ := conn.Call("fn.group_users", []interface{}{SelectedGroupId})
	fmt.Println(info)
	groupUsersTuples := info.Tuples()
	fmt.Println(groupUsersTuples)
	if len(groupUsersTuples) == 1 && len(groupUsersTuples[0]) == 0 {
		fmt.Println("no")
		return false
	} else {
		for i := 0; i < len(groupUsersTuples); i++ {
			userName := groupUsersTuples[i][0].(string)
			userRole := groupUsersTuples[i][1].(string)
			fmt.Println(userName, userRole)
			rowUser, _ := gtk.ListBoxRowNew()
			if userRole == "admin" {
				userName = userName + " (admin)"
			}
			labelUser, _ := gtk.LabelNew(userName)
			markup := fmt.Sprintf("<span font_desc='Serif Bold 15'>%s</span>", userName)
			labelUser.SetMarkup(markup)
			labelUser.SetHAlign(gtk.ALIGN_START)
			labelUser.SetJustify(gtk.JUSTIFY_CENTER)
			rowUser.Add(labelUser)
			rowUser.SetName(userRole)
			usersListbox.Prepend(rowUser)

		}
		fmt.Println(SelectedGroupId)
		fmt.Println(info)
		usersListbox.ShowAll()
		return true
	}
}

func getTextBeforeSubstring(input string, substring string) string {
	index := strings.Index(input, substring)
	if index != -1 {
		return input[:index]
	} else {
		return input
	}
}

func GetMsg(conn *tarantool.Connection, msgListBox *gtk.ListBox) {
	var lastTimedMsg uint64

	fmt.Println("Начало")

	fmt.Println("len(messagesArr) ", len(messagesArr))
	if len(messagesArr) == 0 {
		clearListbox(msgListBox)
		lastTimedMsg = 0 // В самом начале загружаем все сообщения
	} else {
		lastTimedMsg = messagesArr[len(messagesArr)-1].time
	}

	infoTimedMsg, err := conn.Call("fn.time_group_msg", []interface{}{SelectedGroupId, lastTimedMsg})
	if err != nil {
		panic(err)
	}
	fmt.Println("Получаем смс")

	newMessagesCntTuples := infoTimedMsg.Tuples()
	cntMsg := int(newMessagesCntTuples[0][0].(uint64))
	fmt.Println("lastTimedMsg", lastTimedMsg)
	//fmt.Println(newMessagesCntTuples)

	var newMessages []MessageStruct
	for i := 0; i < cntMsg; i++ {
		newMessagesTuples := newMessagesCntTuples[1][i].([]interface{})
		newMessages = append(newMessages, MessageStruct{newMessagesTuples[0].(string), newMessagesTuples[1].(string), newMessagesTuples[2].(uint64), newMessagesTuples[3].(string)})
	}

	fmt.Println("Заполняем массив")

	if cntMsg != 0 {
		for i := 0; i < cntMsg; i++ {

			msgText := newMessages[i].message
			msgUser := newMessages[i].user
			msgTime := newMessages[i].msgTime
			msgId := newMessages[i].msgId

			newMsg := ""

			if msgUser == "system" {
				newMsg = msgText
			} else {
				newMsg = msgUser + "(" + SelectedGroupName + "): " + msgText
			}
			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})
			//listbox

			fmt.Println(newMsg)
			rowMsg, err := gtk.ListBoxRowNew()
			if err != nil {
				panic(err)
			}
			labelMsg, err := gtk.LabelNew(newMsg)
			if err != nil {
				panic(err)
			}
			labelMsg.SetHAlign(gtk.ALIGN_START)
			labelMsg.SetJustify(gtk.JUSTIFY_CENTER)
			rowMsg.Add(labelMsg)
			rowMsg.SetName(msgId)
			if msgUser == "system" {
				markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 10' color='#323ea8'>%s</span>", msgText)
				labelMsg.SetMarkup(markup)
				rowMsg.SetName("system")
			}
			//msgListBox.Insert(rowMsg, 1)
			msgListBox.Prepend(rowMsg)
			fmt.Println("Добавляем")
			//msgListBox.ShowAll()
			//rowMsg.ShowAll()
			fmt.Println("Показываем всё")

		}

		msgListBox.ShowAll()

	}
}
