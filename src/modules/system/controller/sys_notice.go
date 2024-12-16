package controller

import (
	"mask_api_gin/src/framework/context"
	"mask_api_gin/src/framework/response"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"

	"fmt"

	"github.com/gin-gonic/gin"
)

// NewSysNotice 实例化控制层
var NewSysNotice = &SysNoticeController{
	sysNoticeService: service.NewSysNotice,
}

// SysNoticeController 通知公告信息
//
// PATH /system/notice
type SysNoticeController struct {
	sysNoticeService *service.SysNotice // 公告服务
}

// List 通知公告列表
//
// GET /list
func (s SysNoticeController) List(c *gin.Context) {
	query := context.QueryMap(c)
	rows, total := s.sysNoticeService.FindByPage(query)
	c.JSON(200, response.OkData(map[string]any{"rows": rows, "total": total}))
}

// Info 通知公告信息
//
// GET /:noticeId
func (s SysNoticeController) Info(c *gin.Context) {
	noticeId := c.Param("noticeId")
	if noticeId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: noticeId is empty"))
		return
	}

	data := s.sysNoticeService.FindById(noticeId)
	if data.NoticeId == noticeId {
		c.JSON(200, response.OkData(data))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Add 通知公告新增
//
// POST /
func (s SysNoticeController) Add(c *gin.Context) {
	var body model.SysNotice
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.NoticeId != "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: noticeId not is empty"))
		return
	}

	body.CreateBy = context.LoginUserToUserName(c)
	insertId := s.sysNoticeService.Insert(body)
	if insertId != "" {
		c.JSON(200, response.OkData(insertId))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Edit 通知公告修改
//
// PUT /
func (s SysNoticeController) Edit(c *gin.Context) {
	var body model.SysNotice
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", response.FormatBindError(err))
		c.JSON(400, response.CodeMsg(40010, errMsgs))
		return
	}
	if body.NoticeId == "" {
		c.JSON(400, response.CodeMsg(40010, "bind err: noticeId is empty"))
		return
	}

	// 检查是否存在
	noticeInfo := s.sysNoticeService.FindById(body.NoticeId)
	if noticeInfo.NoticeId != body.NoticeId {
		c.JSON(200, response.ErrMsg("没有权限访问公告信息数据！"))
		return
	}

	noticeInfo.NoticeTitle = body.NoticeTitle
	noticeInfo.NoticeType = body.NoticeType
	noticeInfo.NoticeContent = body.NoticeContent
	noticeInfo.StatusFlag = body.StatusFlag
	noticeInfo.Remark = body.Remark
	noticeInfo.UpdateBy = context.LoginUserToUserName(c)
	rows := s.sysNoticeService.Update(noticeInfo)
	if rows > 0 {
		c.JSON(200, response.Ok(nil))
		return
	}
	c.JSON(200, response.Err(nil))
}

// Remove 通知公告删除
//
// DELETE /:noticeId
func (s SysNoticeController) Remove(c *gin.Context) {
	noticeIdsStr := c.Param("noticeId")
	noticeIds := parse.RemoveDuplicatesToArray(noticeIdsStr, ",")
	if noticeIdsStr == "" || len(noticeIds) <= 0 {
		c.JSON(400, response.CodeMsg(40010, "bind err: noticeId is empty"))
		return
	}

	rows, err := s.sysNoticeService.DeleteByIds(noticeIds)
	if err != nil {
		c.JSON(200, response.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, response.OkMsg(msg))
}
