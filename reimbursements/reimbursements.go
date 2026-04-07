package reimbursements

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"yukoreimburse/database"
	"yukoreimburse/lark"
	"yukoreimburse/tools"
)

type Expense_reports struct {
	Report_id      int    `gorm:"primary_key"`
	Reason         string `gorm:"size:255"`
	Report_time    time.Time
	User_id        int     `gorm:"primary_key"`
	Username       string  `gorm:""` // 非数据库字段，用于存储联表查询结果
	Remarks        string  `gorm:"size:255"`
	Amount         float64 `gorm:"type:decimal(10,2)"`
	Status         int     `gorm:"default:0"`
	Has_attachment int     `gorm:"default:0"`
}

type New_reports struct {
	Reason      string                  `form:"reason" binding:"required"` // 报销原因
	Amount      float64                 `form:"amount" binding:"required"` // 报销金额
	ExpenseTime string                  `form:"date" binding:"required"`   // 报销时间
	Remarks     string                  `form:"remarks"`                   // 备注
	Attachments []*multipart.FileHeader `form:"attachments[]"`             // 附件
}

// 获取报销单的函数（分页查询，每页30条）
func GetExpenseReports(user_id int, is_admin int, page int) ([]Expense_reports, error) {
	// 初始化空的报销单列表
	var expense_reports []Expense_reports

	// 每页记录数
	limit := 30
	// 计算偏移量
	offset := (page - 1) * limit

	// 构建基础 SQL 查询
	sql := `
		SELECT 
			e.report_id,
			e.reason,
			e.report_time,
			e.user_id,
			u.username,
			e.remarks,
			e.amount,
			e.status,
			e.has_attachment
		FROM 
			expense_reports e
		LEFT JOIN 
			users u 
		ON 
			e.user_id = u.user_id
	`

	// 添加权限条件
	if is_admin == 0 {
		// 如果是普通用户，查询该用户的报销单，并按 report_id 倒序排序
		sql += " WHERE e.user_id = ?"
	}
	// 无论是管理员还是普通用户，都需要按照 report_id 倒序排序
	sql += " ORDER BY e.report_id DESC"
	sql += " LIMIT ? OFFSET ?"

	// 执行查询
	var err error
	if is_admin == 0 {
		err = database.Gdb.Raw(sql, user_id, limit, offset).Scan(&expense_reports).Error
	} else {
		err = database.Gdb.Raw(sql, limit, offset).Scan(&expense_reports).Error
	}

	// 错误处理
	if err != nil {
		return nil, errors.New("获取报销单失败：" + err.Error())
	}

	// 返回查询结果
	return expense_reports, nil
}

// 获取所有报销单数量
func GetAllReportsCount(user_id int, is_admin int) (int, error) {
	// 初始化报销单计数器
	var count int

	// 构建基础 SQL 查询
	sql := `
		SELECT 
			COUNT(*) 
		FROM 
			expense_reports 
	`

	// 添加权限条件
	if is_admin == 0 {
		sql += " WHERE user_id = ?"
		if err := database.Gdb.Raw(sql, user_id).Scan(&count).Error; err != nil {
			return 0, errors.New("获取报销单数量失败：" + err.Error())
		}
	} else {
		if err := database.Gdb.Raw(sql).Scan(&count).Error; err != nil {
			return 0, errors.New("获取报销单数量失败：" + err.Error())
		}

	}
	return count, nil
}

// 获取已报销数量
func GetExpenseReportsCount(user_id int, is_admin int) (int, error) {
	// 初始化报销单计数器
	var count int

	// 构建基础 SQL 查询
	sql := `
		SELECT 
			COUNT(*) 
		FROM 
			expense_reports 
		WHERE 
			status = 1
	`

	// 添加权限条件
	if is_admin == 0 {
		sql += " AND user_id = ?"
		if err := database.Gdb.Raw(sql, user_id).Scan(&count).Error; err != nil {
			return 0, errors.New("获取报销单数量失败：" + err.Error())
		}
	} else {
		if err := database.Gdb.Raw(sql).Scan(&count).Error; err != nil {
			return 0, errors.New("获取报销单数量失败：" + err.Error())
		}
	}

	// 返回报销单数量
	return count, nil
}

func GetTotalReimbursedAmount(user_id int, is_admin int) (float64, error) {
	// 初始化已报销金额
	var totalAmount sql.NullFloat64

	// 构建基础 SQL 查询
	sqlQuery := `
		SELECT 
			COALESCE(SUM(amount), 0) 
		FROM 
			expense_reports 
		WHERE 
			status = 1
	`

	// 添加权限条件
	if is_admin == 0 {
		sqlQuery += " AND user_id = ?"
		if err := database.Gdb.Raw(sqlQuery, user_id).Scan(&totalAmount).Error; err != nil {
			return 0, errors.New("获取已报销金额失败：" + err.Error())
		}
	} else {
		if err := database.Gdb.Raw(sqlQuery).Scan(&totalAmount).Error; err != nil {
			return 0, errors.New("获取已报销金额失败：" + err.Error())
		}
	}

	// 如果 totalAmount 为 NULL，返回 0
	if !totalAmount.Valid {
		return 0, nil
	}

	// 返回已报销金额
	return totalAmount.Float64, nil
}

func GetPendingExpenseReportsCount(user_id int, is_admin int) (int, error) {
	// 初始化待报销单计数器
	var count int

	// 构建基础 SQL 查询
	sql := `
		SELECT 
			COUNT(*) 
		FROM 
			expense_reports 
		WHERE 
			status = 0
	`

	// 添加权限控制：如果是管理员，获取所有待报销单；如果是普通用户，获取该用户的待报销单
	if is_admin == 0 {
		sql += " AND user_id = ?"
		if err := database.Gdb.Raw(sql, user_id).Scan(&count).Error; err != nil {
			return 0, errors.New("获取待报销总数失败：" + err.Error())
		}
	} else {
		if err := database.Gdb.Raw(sql).Scan(&count).Error; err != nil {
			return 0, errors.New("获取待报销总数失败：" + err.Error())
		}
	}

	// 返回待报销的总数
	return count, nil
}

func GetPendingExpenseReportsAmount(user_id int, is_admin int) (float64, error) {
	// 初始化待报销金额
	var totalAmount float64

	// 构建基础 SQL 查询
	sql := `
		SELECT 
			COALESCE(SUM(amount), 0) AS total_amount
		FROM 
			expense_reports 
		WHERE 
			status = 0
	`

	// 添加权限控制：如果是管理员，获取所有待报销的金额；如果是普通用户，获取该用户的待报销金额
	if is_admin == 0 {
		sql += " AND user_id = ?"
		if err := database.Gdb.Raw(sql, user_id).Scan(&totalAmount).Error; err != nil {
			return 0, errors.New("获取待报销金额失败：" + err.Error())
		}
	} else {
		if err := database.Gdb.Raw(sql).Scan(&totalAmount).Error; err != nil {
			return 0, errors.New("获取待报销金额失败：" + err.Error())
		}
	}

	// 返回待报销的总金额（金额可能为 0，不作为错误处理）
	return totalAmount, nil
}

// 按用户 ID、状态及时间范围筛选所有报销单（不分页）
func GetFilteredExpenseReports(user_id *int, status *int, startTime *time.Time, endTime *time.Time) ([]Expense_reports, error) {
	// 初始化空的报销单列表
	var expense_reports []Expense_reports

	// 构建基础 SQL 查询
	sql := `
        SELECT 
            e.report_id,
            e.reason,
            e.report_time,
            e.user_id,
            u.username,
            e.remarks,
            e.amount,
            e.status,
            e.has_attachment
        FROM 
            expense_reports e
        LEFT JOIN 
            users u 
        ON 
            e.user_id = u.user_id
    `

	// 条件拼接
	var conditions []string
	var args []interface{}

	if user_id != nil {
		conditions = append(conditions, "e.user_id = ?")
		args = append(args, *user_id)
	}
	if status != nil {
		conditions = append(conditions, "e.status = ?")
		args = append(args, *status)
	}
	if startTime != nil {
		conditions = append(conditions, "e.report_time >= ?")
		args = append(args, *startTime)
	}
	if endTime != nil {
		conditions = append(conditions, "e.report_time <= ?")
		args = append(args, *endTime)
	}

	// 如果有条件，添加 WHERE 子句
	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	// 添加排序
	sql += " ORDER BY e.report_id DESC"

	// 执行查询
	err := database.Gdb.Raw(sql, args...).Scan(&expense_reports).Error

	// 错误处理
	if err != nil {
		return nil, errors.New("获取报销单失败：" + err.Error())
	}

	// 返回查询结果
	return expense_reports, nil
}

func UpdateReportsStatus(reportIDs []int, status int, is_admin int, user_id int, username string) error {
	var currentStatuses []int
	var userIDs []int

	// 创建日志文件或打开已有日志文件
	logFile, err := os.OpenFile("update_reports_status.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("无法打开日志文件：" + err.Error())
	}
	defer logFile.Close()

	// 设置日志输出
	logger := log.New(logFile, "", log.LstdFlags)

	// 查询选中报销单的状态和用户 ID
	err = database.Gdb.Model(&Expense_reports{}).
		Where("report_id IN (?)", reportIDs).
		Pluck("status", &currentStatuses).Error

	if err != nil {
		logger.Printf("操作人: %s (ID: %d), 查询报销单状态失败: %v\n", username, user_id, err)
		return errors.New("查询报销单状态失败：" + err.Error())
	}

	err = database.Gdb.Model(&Expense_reports{}).
		Where("report_id IN (?)", reportIDs).
		Pluck("user_id", &userIDs).Error

	if err != nil {
		logger.Printf("操作人: %s (ID: %d), 查询报销单用户失败: %v\n", username, user_id, err)
		return errors.New("查询报销单用户失败：" + err.Error())
	}

	// 检查是否有报销单
	if len(currentStatuses) == 0 || len(userIDs) == 0 {
		logger.Printf("操作人: %s (ID: %d), 未找到相关的报销单\n", username, user_id)
		return errors.New("未找到相关的报销单")
	}

	// 检查状态是否一致
	for _, currentStatus := range currentStatuses {
		if currentStatus != currentStatuses[0] {
			logger.Printf("操作人: %s (ID: %d), 选中的报销单状态不一致\n", username, user_id)
			return errors.New("选中的报销单状态不一致，只能批量修改相同状态的报销单")
		}
	}

	// 检查用户 ID 是否一致
	for _, reportUserID := range userIDs {
		if reportUserID != userIDs[0] {
			logger.Printf("操作人: %s (ID: %d), 选中的报销单不属于同一个用户\n", username, user_id)
			return errors.New("选中的报销单不属于同一个用户，只能批量处理同一个用户的报销单")
		}
	}

	// 如果是管理员，且原状态为“已撤回”，返回错误
	if is_admin == 1 && currentStatuses[0] == 3 {
		logger.Printf("管理员: %s (ID: %d), 尝试处理已撤回的报销单\n", username, user_id)
		return errors.New("已撤回的报销单不可处理")
	}

	// 如果是普通用户，且想修改状态为 1 或 2 的报销单，返回无权限错误
	if is_admin == 0 && (currentStatuses[0] == 1 || currentStatuses[0] == 2) {
		logger.Printf("普通用户: %s (ID: %d), 尝试修改已通过或被驳回的报销单\n", username, user_id)
		return errors.New("普通用户无权限修改已通过或被驳回的报销单")
	}

	// 记录修改前状态日志
	logger.Printf("操作人: %s (ID: %d), 准备更新报销单状态: IDs=%v, 当前状态=%d, 新状态=%d, 操作角色=%d\n",
		username, user_id, reportIDs, currentStatuses[0], status, is_admin)

	// 如果通过校验，执行状态更新
	err = database.Gdb.Model(&Expense_reports{}).
		Where("report_id IN (?)", reportIDs).
		Update("status", status).Error

	if err != nil {
		logger.Printf("操作人: %s (ID: %d), 更新报销单状态失败: %v\n", username, user_id, err)
		return errors.New("更新报销单状态失败：" + err.Error())
	}

	// 记录成功日志
	logger.Printf("操作人: %s (ID: %d), 成功更新报销单状态: IDs=%v, 原状态=%d, 新状态=%d, 操作角色=%d\n",
		username, user_id, reportIDs, currentStatuses[0], status, is_admin)

	var larkID string
	err = database.Gdb.
		Table("users").
		Select("lark_id").
		Where("user_id = ?", userIDs[0]).
		Scan(&larkID).Error

	if err != nil {
		return err
	}
	// 通过飞书通知对应用户
	uuid, err := tools.GetUuid()
	if err != nil {
		return err
	}

	fmt.Println(larkID)

	text := "你有报销单的状态发生变化，请前往网页端查看\nyukoexp.com"

	err = lark.SendText(larkID, uuid, text)
	if err != nil {

		return errors.New("飞书通知报销人失败！\n请手动通知！并向开发人员反馈失败原因！\n: " + err.Error())
	}

	return nil
}

// 新增报销单并上传附件
func AddExpenseReport(expenseReport *New_reports, userID int, username string, hasAttachment int, attachmentDir string) (int, error) {
	// 使用原生SQL插入报销单数据
	sql := `INSERT INTO expense_reports (reason, report_time, user_id, remarks, amount, status, has_attachment) 
	        VALUES (?, ?, ?, ?, ?, 0, ?)`
	// 执行SQL插入数据
	result := database.Gdb.Exec(sql, expenseReport.Reason, expenseReport.ExpenseTime, userID, expenseReport.Remarks, expenseReport.Amount, hasAttachment)
	if result.Error != nil {
		return 0, errors.New("插入报销单失败: " + result.Error.Error())
	}

	// 获取插入的报销单ID（通过 GORM 的 `Last()` 方法获取最新记录的 ID）
	var reportID int
	if err := database.Gdb.Raw("SELECT LAST_INSERT_ID()").Scan(&reportID).Error; err != nil {
		return 0, errors.New("获取报销单ID失败: " + err.Error())
	}

	// 如果有附件，执行文件上传
	if hasAttachment == 1 {
		err := tools.UploadAttachments(reportID, expenseReport.Attachments, attachmentDir)
		if err != nil {
			database.Gdb.Exec(`DELETE FROM expense_reports WHERE report_id = ?`, reportID)
			_ = os.RemoveAll(tools.AttachmentReportDir(reportID, attachmentDir))
			return 0, errors.New("附件上传失败，报销单已删除: " + err.Error())
		}

		updateSQL := `UPDATE expense_reports SET has_attachment = 1 WHERE report_id = ?`
		updateResult := database.Gdb.Exec(updateSQL, reportID)
		if updateResult.Error != nil {
			database.Gdb.Exec(`DELETE FROM expense_reports WHERE report_id = ?`, reportID)
			_ = os.RemoveAll(tools.AttachmentReportDir(reportID, attachmentDir))
			return 0, errors.New("更新附件标志失败，报销单已删除: " + updateResult.Error.Error())
		}
	}

	// // 通过飞书通知管理员
	// uuid, err := tools.GetUuid()
	// if err != nil {
	// 	return 0, err
	// }

	// text := fmt.Sprintf("%s 提交了新的报销单，请及时处理\nyukoexp.com", username)

	// err = lark.SendText("ou_361a84bb42e1210aafb90cc262e6ef77", uuid, text)
	// if err != nil {

	// 	return 0, errors.New("飞书通知管理员失败！\n请手动通知管理员！并向开发人员反馈失败原因！\n: " + err.Error())
	// }

	// 返回报销单ID
	return reportID, nil
}
