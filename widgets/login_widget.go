package widgets

import (
	"log"
	"net/http"

	"github.com/go-xorm/xorm"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/crawler"
	"schannel-qt5/models"
)

// LoginWidget 登录界面
type LoginWidget struct {
	widgets.QWidget

	// loginUser 将登录所用的用户名传递给父控件
	_ func(string)         `signal:"loginUser"`
	_ func()               `signal:"loginFailed,auto"`
	_ func([]*http.Cookie) `signal:"logined"`

	username    *widgets.QComboBox
	password    *widgets.QLineEdit
	loginStatus *ColorLabel
	remember    *widgets.QCheckBox
	conf        *config.UserConfig
	logger      *log.Logger
	db          *xorm.Engine
}

// NewLoginWidget2 根据config，logger，db生成登录控件
func NewLoginWidget2(conf *config.UserConfig, logger *log.Logger, db *xorm.Engine) *LoginWidget {
	if conf == nil || logger == nil {
		return nil
	}

	widget := NewLoginWidget(nil, 0)
	widget.conf = conf
	widget.logger = logger
	widget.db = db
	widget.InitUI()

	return widget
}

func (l *LoginWidget) InitUI() {
	userLabel := widgets.NewQLabel2("用户名: ", nil, 0)
	l.username = widgets.NewQComboBox(nil)
	l.username.SetEditable(true)
	users, err := models.GetAllUsers(l.db)
	if err != nil {
		l.logger.Fatalln(err)
	}
	// 第一项为空
	names := make([]string, 1, len(users)+1)
	for _, v := range users {
		names = append(names, v.Name)
	}
	l.username.AddItems(names)
	// 实现记住用户密码
	l.username.ConnectCurrentTextChanged(l.setPassword)

	userLabel.SetBuddy(l.username)
	userInputLayout := widgets.NewQHBoxLayout()
	userInputLayout.AddWidget(userLabel, 0, 0)
	userInputLayout.AddWidget(l.username, 0, 0)

	passwdLabel := widgets.NewQLabel2("密码:", nil, 0)
	l.password = widgets.NewQLineEdit(nil)
	l.password.SetPlaceholderText("密码")
	l.password.SetEchoMode(widgets.QLineEdit__Password)

	passwdLabel.SetBuddy(l.password)
	passwdInputLayout := widgets.NewQHBoxLayout()
	passwdInputLayout.AddWidget(passwdLabel, 0, 0)
	passwdInputLayout.AddWidget(l.password, 0, 0)

	l.loginStatus = NewColorLabelWithColor("用户名或密码错误，请重试", "red")
	l.loginStatus.Hide()

	l.remember = widgets.NewQCheckBox2("记住用户名和密码", nil)
	loginButton := widgets.NewQPushButton2("登录", nil)
	loginButton.ConnectClicked(l.checkLogin)

	loginLayout := widgets.NewQHBoxLayout()
	loginLayout.AddStretch(0)
	loginLayout.AddWidget(loginButton, 0, 0)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(l.loginStatus, 0, 0)
	mainLayout.AddLayout(userInputLayout, 0)
	mainLayout.AddLayout(passwdInputLayout, 0)
	mainLayout.AddWidget(l.remember, 0, 0)
	mainLayout.AddLayout(loginLayout, 0)
	l.SetLayout(mainLayout)
}

func (l *LoginWidget) checkLogin(_ bool) {
	// 防止多次点击登录按钮或在登录时改变lineedit内容
	l.SetEnabled(false)
	defer l.SetEnabled(true)

	passwd := l.password.Text()
	user := l.username.CurrentText()
	if user == "" || passwd == "" {
		l.LoginFailed()
		return
	}

	// 登录
	cookies, err := crawler.GetAuth(user, passwd, l.conf.Proxy.String())
	if err != nil {
		l.logger.Printf("crawler failed: %v\n", err)
		l.LoginFailed()
		return
	}

	// 登陆成功，记住密码
	if l.remember.IsChecked() {
		if err := models.SetUserPassword(l.db, user, []byte(passwd)); err != nil {
			l.logger.Println(err)
		}
	}

	// 传递登录信息
	l.logger.Printf("logined as [%s] success\n", user)
	l.LoginUser(user)
	l.Logined(cookies)
}

func (l *LoginWidget) loginFailed() {
	if l.loginStatus.IsHidden() {
		l.loginStatus.Show()
	}
}

// setPassword 将密码不为null的用户显示
func (l *LoginWidget) setPassword(user string) {
	info, err := models.GetUserPassword(l.db, user)
	if err != nil {
		l.logger.Println(err)
		l.username.SetCurrentText(user)
		l.password.SetText("")
		return
	} else if info.Passwd != nil {
		l.password.SetText(string(info.Passwd))
		return
	}
}
