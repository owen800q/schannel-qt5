package widgets

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
	"schannel-qt5/parser"
	"strconv"
)

// InvoiceDialog 显示全部的账单信息
type InvoiceDialog struct {
	widgets.QDialog

	table    *widgets.QTableWidget
	infobar  *widgets.QStatusBar
	selected *widgets.QLabel
	link     *widgets.QLabel

	invoices []*parser.Invoice
}

// NewInvoiceDialogWithData 生成dialog
func NewInvoiceDialogWithData(data []*parser.Invoice) *InvoiceDialog {
	dialog := NewInvoiceDialog(nil, 0)
	dialog.invoices = data

	// 设置infobar，选中内容时显示账单链接
	dialog.infobar = widgets.NewQStatusBar(nil)
	dialog.selected = widgets.NewQLabel2("未选中", nil, 0)
	dialog.link = widgets.NewQLabel(nil, 0)
	dialog.infobar.AddPermanentWidget(dialog.selected, 0)
	dialog.infobar.AddPermanentWidget(dialog.link, 0)

	// 初始化table，数据已经被排序
	dialog.table = widgets.NewQTableWidget(nil)
	// 设置行数，不设置将不显示任何数据
	dialog.table.SetRowCount(len(dialog.invoices))
	// 设置表头
	dialog.table.SetColumnCount(6)
	dialog.table.SetHorizontalHeaderLabels([]string{
		"账单编号",
		"链接",
		"开始日期",
		"结束日期",
		"金额（元）",
		"支付状态",
	})
	// 去除边框
	dialog.table.SetShowGrid(false)
	dialog.table.SetFrameShape(widgets.QFrame__NoFrame)
	// 去除行号
	dialog.table.VerticalHeader().SetVisible(false)
	// 设置table的数据项目
	dialog.setTable()
	dialog.table.ConnectItemClicked(dialog.setLink)

	//TODO: 点击colorlabel也可以切换链接显示
	/*dialog.table.ConnectCellClicked(func (row, col int) {
		invoice := dialog.invoices[row]
		dialog.selected.SetText(fmt.Sprintf("选中第%d行", row+1))
		dialog.link.SetText(invoice.Link)
	})*/

	// 设置不可编辑table
	dialog.table.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)

	vbox := widgets.NewQVBoxLayout()
	vbox.AddWidget(dialog.table, 0, 0)
	vbox.AddStretch(0)
	vbox.AddWidget(dialog.infobar, 0, 0)
	dialog.SetLayout(vbox)

	return dialog
}

// setTable 设置table
func (dialog *InvoiceDialog) setTable() {
	for row := 0; row < len(dialog.invoices); row++ {
		invoice := dialog.invoices[row]

		number := widgets.NewQTableWidgetItem2(invoice.Number, 0)
		dialog.table.SetItem(row, 0, number)
		link := widgets.NewQTableWidgetItem2(invoice.Link, 0)
		dialog.table.SetItem(row, 1, link)

		startTime := time2string(invoice.StartDate)
		start := widgets.NewQTableWidgetItem2(startTime, 0)
		dialog.table.SetItem(row, 2, start)
		expireTime := time2string(invoice.ExpireDate)
		expire := widgets.NewQTableWidgetItem2(expireTime, 0)
		dialog.table.SetItem(row, 3, expire)

		payment := strconv.FormatInt(invoice.Payment, 10)
		pay := widgets.NewQTableWidgetItem2(payment, 0)
		dialog.table.SetItem(row, 4, pay)

		text := ""
		color := ""
		if invoice.State == parser.NeedPay {
			text = "未付款"
			color = "red"
		} else if invoice.State == parser.FinishedPay {
			text = "已付款"
			color = "green"
		}
		label := NewColorLabelWithColor(text, color)
		dialog.table.SetCellWidget(row, 5, label)
	}
}

// setLink 当选中row中的单元格时将链接更新到infobar
func (dialog *InvoiceDialog) setLink(item *widgets.QTableWidgetItem) {
	index := item.Row()
	invoice := dialog.invoices[index]
	dialog.selected.SetText(fmt.Sprintf("选中第%d行", index+1))
	dialog.link.SetText(invoice.Link)
}