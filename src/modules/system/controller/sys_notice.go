package controller

import (
	"fmt"
	"mask_api_gin/src/framework/utils/ctx"
	"mask_api_gin/src/framework/utils/parse"
	"mask_api_gin/src/framework/vo/result"
	"mask_api_gin/src/modules/system/model"
	"mask_api_gin/src/modules/system/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 实例化控制层 SysNoticeController 结构体
var NewSysNotice = &SysNoticeController{
	sysNoticeService: service.NewSysNoticeImpl,
}

// 通知公告信息
//
// PATH /system/notice
type SysNoticeController struct {
	// 公告服务
	sysNoticeService service.ISysNotice
}

// 通知公告列表
//
// GET /list
func (s *SysNoticeController) List(c *gin.Context) {
	querys := ctx.QueryMap(c)
	data := s.sysNoticeService.SelectNoticePage(querys)
	c.JSON(200, result.Ok(data))
}

// 通知公告信息
//
// GET /:noticeId
func (s *SysNoticeController) Info(c *gin.Context) {
	noticeId := c.Param("noticeId")
	if noticeId == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	data := s.sysNoticeService.SelectNoticeById(noticeId)
	if data.NoticeID == noticeId {
		c.JSON(200, result.OkData(data))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 通知公告新增
//
// POST /
func (s *SysNoticeController) Add(c *gin.Context) {
	var body model.SysNotice
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.NoticeID != "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	body.CreateBy = ctx.LoginUserToUserName(c)
	insertId := s.sysNoticeService.InsertNotice(body)
	if insertId != "" {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 通知公告修改
//
// PUT /
func (s *SysNoticeController) Edit(c *gin.Context) {
	var body model.SysNotice
	err := c.ShouldBindBodyWith(&body, binding.JSON)
	if err != nil || body.NoticeID == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}

	// 检查是否存在
	notice := s.sysNoticeService.SelectNoticeById(body.NoticeID)
	if notice.NoticeID != body.NoticeID {
		c.JSON(200, result.ErrMsg("没有权限访问公告信息数据！"))
		return
	}

	body.UpdateBy = ctx.LoginUserToUserName(c)
	rows := s.sysNoticeService.UpdateNotice(body)
	if rows > 0 {
		c.JSON(200, result.Ok(nil))
		return
	}
	c.JSON(200, result.Err(nil))
}

// 通知公告删除
//
// DELETE /:noticeIds
func (s *SysNoticeController) Remove(c *gin.Context) {
	noticeIds := c.Param("noticeIds")
	if noticeIds == "" {
		c.JSON(400, result.CodeMsg(400, "参数错误"))
		return
	}
	// 处理字符转id数组后去重
	ids := strings.Split(noticeIds, ",")
	uniqueIDs := parse.RemoveDuplicates(ids)
	if len(uniqueIDs) <= 0 {
		c.JSON(200, result.Err(nil))
		return
	}
	rows, err := s.sysNoticeService.DeleteNoticeByIds(uniqueIDs)
	if err != nil {
		c.JSON(200, result.ErrMsg(err.Error()))
		return
	}
	msg := fmt.Sprintf("删除成功：%d", rows)
	c.JSON(200, result.OkMsg(msg))
}
