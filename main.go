package main

import (
	"fmt"
	"net/http"
	// "path/filepath"
	"strconv"
	"yukoreimburse/files"
	"yukoreimburse/reimbursements"
	"yukoreimburse/tools"
	"yukoreimburse/users"

	"github.com/gin-gonic/gin"
)

func main() {

	// 搭建web服务器
	r := gin.Default()

	// 加载静态资源
	r.StaticFS("/static", http.Dir("./static"))

	attachmentDir := "/attachments"

	// 加载附件目录
	r.StaticFS(attachmentDir, http.Dir("."+attachmentDir))

	// 加载模板
	r.LoadHTMLGlob("./tpl/*")

	// // 用户全局中间件
	var user users.Users
	r.Use(func(ctx *gin.Context) {
		// 判断是否登录
		val := tools.GetSession(ctx, "user")

		user, _ = val.(users.Users)

		ctx.Next()
	})

	// 首页
	r.GET("/", func(ctx *gin.Context) {
		// 判断是否登录
		if user.Username == "" {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    "请先登录",
				"target": "login",
			})
			return
		}

		// 获取最新30条记录（第一页数据）
		var expense_reports []reimbursements.Expense_reports
		expense_reports, err := reimbursements.GetExpenseReports(user.User_id, user.Is_admin, 1)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "/",
			})
			return
		}

		// 获取总记录数，用于计算总页数
		count, err := reimbursements.GetAllReportsCount(user.User_id, user.Is_admin)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{"msg": err.Error()})
			return
		}

		totalPages := (count + 29) / 30 // 计算总页数，向上取整

		// 获取已报销数量
		expense_report_count, err := reimbursements.GetExpenseReportsCount(user.User_id, user.Is_admin)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "/",
			})
			return
		}

		// 获取已报销金额
		total_reimbursed_amount, err := reimbursements.GetTotalReimbursedAmount(user.User_id, user.Is_admin)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "/",
			})
			return
		}

		// 获取待报销数量
		pending_expense_report_count, err := reimbursements.GetPendingExpenseReportsCount(user.User_id, user.Is_admin)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "/",
			})
			return
		}

		// 获取待报销金额
		pending_expense_report_amount, err := reimbursements.GetPendingExpenseReportsAmount(user.User_id, user.Is_admin)
		if err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "/",
			})
			return
		}

		// 计算上一页和下一页页码
		prevPage := 1
		var nextPage int
		if totalPages > 1 {
			nextPage = 2
		} else {
			nextPage = 1
		}
		// fmt.Println(nextPage)

		// 将首页数据传递给模板
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"user":                          user,
			"expense_reports":               expense_reports,
			"expense_report_count":          expense_report_count,
			"total_reimbursed_amount":       total_reimbursed_amount,
			"pending_expense_report_count":  pending_expense_report_count,
			"pending_expense_report_amount": pending_expense_report_amount,
			"totalPages":                    totalPages,
			"prevPage":                      prevPage,
			"nextPage":                      nextPage,
		})
	})

	// 首页处理分页查询的 POST 请求
	r.POST("/expense_reports", func(c *gin.Context) {
		pageStr := c.PostForm("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		expense_reports, err := reimbursements.GetExpenseReports(user.User_id, user.Is_admin, page)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		// 获取总记录数和总页数
		count, err := reimbursements.GetAllReportsCount(user.User_id, user.Is_admin)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		totalPages := (count + 29) / 30

		// 创建报销单 HTML
		expenseReportsHTML := ""
		for _, report := range expense_reports {
			// 判断状态并设置相应的文本
			var statusText string
			switch report.Status {
			case 0:
				statusText = "待审批"
			case 1:
				statusText = "已通过"
			case 2:
				statusText = "驳回"
			default:
				statusText = "报销人撤回"
			}

			// 生成报销单的 HTML，根据是否有附件来决定是否显示查看附件按钮
			var attachmentButton string
			if report.Has_attachment == 1 {
				attachmentButton = fmt.Sprintf(`<a href="javascript:void(0);" class="bg-transparent hover:bg-red-500 text-red-700 font-semibold hover:text-white py-2 px-4 border border-red-500hover:border-transparent rounded" onclick="fetchFiles(%d)">查看附件</a>`, report.Report_id)
			} else {
				attachmentButton = "无附件"
			}

			// 生成报销单的 HTML
			expenseReportsHTML += fmt.Sprintf(`
       		<tr>
                <td>%d</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>%s</td>
                <td>¥%.2f</td>
                <td>%s</td>
            </tr>`,
				report.Report_id, report.Reason, report.Report_time, report.Username, report.Remarks, attachmentButton, report.Amount, statusText)
		}

		// 生成分页 HTML
		paginationHTML := fmt.Sprintf(`
        <button onclick="loadExpenseReports(%d)" class="bg-gray-200 hover:bg-gray-500 text-gray-900 font-bold py-2 px-4 rounded-l">上一页</button>
        <span>第 %d 页 / 共 %d 页</span>
        <button onclick="loadExpenseReports(%d)" class="bg-gray-200 hover:bg-gray-500 text-gray-900 font-bold py-2 px-4 rounded-r">下一页</button>`,
			page-1, page, totalPages, page+1)

		// 返回数据
		c.JSON(http.StatusOK, gin.H{
			"expense_reports_html": expenseReportsHTML,
			"pagination_html":      paginationHTML,
		})
	})

	// 路由处理：获取文件列表
	r.GET("/get_files", func(ctx *gin.Context) {
		// 从请求中获取报销单 ID
		reportID := ctx.DefaultQuery("report_id", "") // 获取报销单 ID
		if reportID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "report_id is required"})
			return
		}

		// 将 report_id 转换为整数
		reportIDInt, err := strconv.Atoi(reportID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report_id"})
			return
		}

		// 获取文件夹中的所有文件
		files, err := files.GetFiles(reportIDInt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回文件列表
		ctx.JSON(http.StatusOK, gin.H{"files": files})
	})

	// 登录页
	r.GET("/login", func(ctx *gin.Context) {
		captcha := tools.GetCaptcha(ctx)
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"user":    user,
			"captcha": captcha,
		})
	})

	// 验证登录
	r.POST("/checkLogin", func(ctx *gin.Context) {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		captcha := ctx.PostForm("captcha")

		// 验证验证码
		if !tools.CheckCaptcha(ctx, captcha) {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    "验证码错误",
				"target": "login",
			})
			return
		}

		// 验证数据
		user, err := users.CheckUser(username, password)

		if err != nil {
			// 验证失败
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "login",
			})
			return
		}

		// fmt.Println(user)

		// 验证成功，设置session
		if err := tools.SetSession(ctx, "user", user); err != nil {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    err.Error(),
				"target": "login",
			})
			return
		}

		ctx.Redirect(http.StatusFound, "/")
	})

	// 退出
	r.GET("/logout", func(ctx *gin.Context) {
		tools.DelSession(ctx, "user")
		ctx.Redirect(http.StatusFound, "login")
	})

	// 获取验证码
	r.GET("/getCaptcha", func(ctx *gin.Context) {
		captcha := tools.GetCaptcha(ctx)
		ctx.String(200, captcha)
	})

	// 编辑页面
	r.GET("/edit", func(ctx *gin.Context) {
		// 判断是否登录
		if user.Username == "" {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    "请先登录",
				"target": "login",
			})
			return
		}
		var usernameList []users.UsernameList
		var err error
		if user.Is_admin == 1 {
			usernameList, err = users.GetUsernameList()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}
		ctx.HTML(http.StatusOK, "edit.html", gin.H{
			"user":         user,
			"usernameList": usernameList,
		})
	})

	// 获取经筛选的报销单
	r.POST("/get_filtered_reports", func(ctx *gin.Context) {
		var user_id *int
		var status *int

		// 从表单获取字符串值
		if user.Is_admin == 1 {
			user_id_string := ctx.PostForm("user_id")
			if user_id_string != "" {
				id, err := strconv.Atoi(user_id_string)
				if err != nil {
					ctx.JSON(400, gin.H{"error": "用户ID格式错误"})
					return
				}
				user_id = &id
			}
		} else {
			user_id = &user.User_id
		}

		status_string := ctx.PostForm("status")

		// 将字符串转为数字，空字符串转为 nil

		if status_string != "" {
			s, err := strconv.Atoi(status_string)
			if err != nil {
				ctx.JSON(400, gin.H{"error": "状态格式错误"})
				return
			}
			status = &s
		}

		// 获取筛选结果
		reports, err := reimbursements.GetFilteredExpenseReports(user_id, status)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "获取报销单失败: " + err.Error()})
			return
		}

		// fmt.Println(reports)

		// 返回筛选结果
		ctx.JSON(200, gin.H{"reports": reports})
	})

	r.POST("/update_reports", func(ctx *gin.Context) {
		var payload struct {
			Report_ids []int `json:"report_ids"` // 确保字段名和前端一致
			Status     int   `json:"status"`
		}

		if user.Is_admin == 0 && payload.Status == 3 {
			ctx.JSON(400, gin.H{"error": "管理员不可撤回报销单，请直接驳回"})
		}

		if user.Is_admin == 0 && (payload.Status == 1 || payload.Status == 2) {
			ctx.JSON(400, gin.H{"error": "无权限"})
		}

		// 绑定 JSON 数据
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(400, gin.H{"error": "请求参数错误", "details": err.Error()})
			return
		}

		err := reimbursements.UpdateReportsStatus(payload.Report_ids, payload.Status, user.Is_admin, user.User_id, user.Username)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "更新报销单状态失败!\n" + err.Error()})
			return
		}

		ctx.JSON(200, gin.H{"message": "更新成功"})

	})

	r.GET("/add", func(ctx *gin.Context) {
		// 判断是否登录
		if user.Username == "" {
			ctx.HTML(http.StatusOK, "error.html", gin.H{
				"msg":    "请先登录",
				"target": "login",
			})
			return
		}
		ctx.HTML(http.StatusOK, "add.html", gin.H{
			"user": user,
		})
	})

	// 路由：接收报销单数据
	r.POST("/add_expense_report", func(ctx *gin.Context) {
		// 创建一个 ExpenseReport 结构体实例
		var expenseReport reimbursements.New_reports

		// 绑定表单数据，包括附件
		if err := ctx.ShouldBind(&expenseReport); err != nil {
			// 如果绑定失败，返回错误信息
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		has_attachment := 0

		// 检查是否有有效附件上传
		if expenseReport.Attachments != nil {
			for _, file := range expenseReport.Attachments {
				if file.Size > 0 { // 检查文件大小是否大于 0
					has_attachment = 1
					break
				}
			}
		}
		fmt.Println(has_attachment)

		_, err := reimbursements.AddExpenseReport(&expenseReport, user.User_id, user.Username, has_attachment, attachmentDir)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 返回成功响应
		ctx.JSON(http.StatusOK, gin.H{
			"message":        "报销单提交成功",
			"has_attachment": has_attachment,
		})
	})

	r.Run(":80")
}
