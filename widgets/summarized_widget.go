package widgets

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"

	"schannel-qt5/config"
	"schannel-qt5/parser"
)

// SummarizedWidget 综合服务信息显示，包括用户信息，服务信息
type SummarizedWidget struct {
	widgets.QWidget

	// 收到数据变动
	_ func() `signal:"dataRefresh,auto"`
	// 发出数据变动，让上层控件更新service
	// 上层控件完成service的更新后发送DataRefresh信号，int值为当前的index
	_ func(int) `signal:"serviceNeedUpdate"`

	// 用户数据接口
	dataBridge UserDataBridge

	// 服务信息面板
	servicePanel *ServicePanel
	// ssr开关面板
	switchPanel *SSRSwitchPanel
	// 使用量信息
	usedPanel *UsedPanel
	// 是否需要付款
	invoicePanel *InvoicePanel

	// 用户名-email
	user string
	// 用户配置
	conf *config.UserConfig
	// 服务信息
	service *parser.Service
	// 综合信息面板编号，因为可能不止一个服务，所以用来做身份区别
	// index与services数组中的索引相同
	index int
}

// NewSummarizedWidget2 创建综合信息控件
func NewSummarizedWidget2(index int,
	user string,
	service *parser.Service,
	conf *config.UserConfig,
	dataBridge UserDataBridge) *SummarizedWidget {
	if user == "" || dataBridge == nil {
		return nil
	}
	sw := NewSummarizedWidget(nil, 0)

	sw.user = user
	sw.dataBridge = dataBridge
	sw.service = service
	sw.conf = conf
	sw.index = index
	sw.InitUI()

	return sw
}

// InitUI 初始化UI
func (sw *SummarizedWidget) InitUI() {
	ssrInfo := sw.dataBridge.SSRInfos(sw.service)
	logger := sw.dataBridge.GetLogger()
	sw.servicePanel = NewServicePanel2(sw.user, ssrInfo)
	sw.invoicePanel = NewInvoicePanelWithData(sw.dataBridge.Invoices())
	sw.switchPanel = NewSSRSwitchPanel2(sw.conf, ssrInfo.Nodes, logger)
	sw.usedPanel = NewUsedPanelWithInfo(sw.user, ssrInfo, logger)

	updateButton := widgets.NewQPushButton2("刷新", nil)
	// 通知上层控件更新sw的service
	updateButton.ConnectClicked(func(_ bool) {
		sw.ServiceNeedUpdate(sw.index)
	})
	leftLayout := widgets.NewQVBoxLayout()
	leftLayout.AddWidget(sw.servicePanel, 0, 0)
	leftLayout.AddWidget(sw.invoicePanel, 0, 0)
	leftLayout.AddStretch(0)
	leftLayout.AddWidget(updateButton, 0, core.Qt__AlignRight)

	rightLayout := widgets.NewQVBoxLayout()
	rightLayout.AddWidget(sw.switchPanel, 0, 0)
	rightLayout.AddWidget(sw.usedPanel, 0, 0)

	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddLayout(leftLayout, 0)
	mainLayout.AddLayout(rightLayout, 0)
	sw.SetLayout(mainLayout)
}

// dataRefresh 处理数据更新
// 一般在SetService之后调用，直接调用将更新servicePanel以外的数据
func (sw *SummarizedWidget) dataRefresh() {
	// sw.service已经被外部更新
	ssrInfo := sw.dataBridge.SSRInfos(sw.service)
	sw.servicePanel.UpadteInfo(sw.user, ssrInfo)
	sw.invoicePanel.UpdateInvoices(sw.dataBridge.Invoices())
	sw.switchPanel.DataRefresh(sw.conf, ssrInfo.Nodes)
	sw.usedPanel.DataRefresh(ssrInfo)
}

// SetService 重新设置service，可用于更新数据
// 调用后一般需要出发DataRefresh信号
func (sw *SummarizedWidget) SetService(service *parser.Service) {
	sw.service = service
}

// UpdateConfig 当config更新时刷新switchPanel
// 一般用作ConfigWidget的信号处理函数
func (sw *SummarizedWidget) UpdateConfig(conf *config.UserConfig) {
	sw.conf = conf
	nodes := sw.dataBridge.SSRInfos(sw.service).Nodes
	sw.switchPanel.DataRefresh(sw.conf, nodes)
}
