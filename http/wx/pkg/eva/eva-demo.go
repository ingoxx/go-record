package eva

import (
	"crypto/rand"
	"encoding/json"
	"github.com/ingoxx/go-record/http/wx/pkg/form"
	"math/big"
)

type SportType struct {
	Name string
}

func NewSportType(name string) SportType {
	return SportType{
		Name: name,
	}
}

func (s SportType) DefaultEvaBoard() ([]*form.MsgBoard, error) {
	var data string
	var evaData []*form.MsgBoard
	if s.Name == "bks" {
		data = `[
		  {"evaluate_id": "aaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "打波啊", "evaluate": "我去，打球怎么还带个扳手？", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/wx_4.JPG", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "bbb", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "广州刘得华", "evaluate": "这里打球得掉层皮才能走", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/wx_1.JPG", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "ccc", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "马路傻手", "evaluate": "打球5分钟，吵架10分钟", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/wx_3.JPG", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "ddd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "你笑个锤子", "evaluate": "打架为啥带个球？", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/wx_4.JPG", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "eee", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "挖得hell", "evaluate": "热身运动一定要做足，篮底内线肉搏才能赢", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/wx-fbs/wx_2.JPG", "time": "2024-04-05 16:32:15"}
		]`
	}

	if s.Name == "bms" {
		data = `[
		{"evaluate_id": "ggg", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "jkasd", "evaluate": "来个人吧", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/profile2/2147.png", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "sss", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "阿毛", "evaluate": "来个高手0.0", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/profile2/2147.png", "time": "2024-04-05 16:32:15"},
		{"evaluate_id": "asd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "16888", "evaluate": "环境还行，就是夏天太热", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/profile2/2146.png", "time": "2024-04-05 16:32:15"}
		]`
	}

	if s.Name == "fbs" {
		data = `[
		  {"evaluate_id": "ada", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "上海吴燕组", "evaluate": "差点被踢爆l_l", "img": "https://mp-578c2584-f82c-45e7-9d53-51332c711501.cdn.bspapp.com/profile2/2145.png", "time": "2024-04-05 16:32:15"},
			{ "evaluate_id": "drj", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "门将界的漏勺", "evaluate": "一场比赛丢了八个球，对手都不好意思庆祝了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "odg", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "爱踩球的阿斌", "evaluate": "想带球过人结果自己被球绊倒，全场掌声送给我", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15" },
		  { "evaluate_id": "werw", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "边路小旋风", "evaluate": "速度很快，就是球留在原地没跟上", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "sbt6", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "只会大脚解围", "evaluate": "全场最远的射门是我解围踢出的", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "ajsdfaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "假装C罗", "evaluate": "学C罗庆祝倒是像，就是射门全飞看台", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
		]`
	}

	if s.Name == "sws" {
		data = `[
		  {"evaluate_id": "ewgw", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "深圳郭富城", "evaluate": "蛙泳太难了l_l", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
			{"evaluate_id": "qt7un", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "游泳别放屁", "evaluate": "恳求各位不要在泳池里边放屁拉屎!!!", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
			{"evaluate_id": "ikho", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "666", "evaluate": "环境还行，水质也干净", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
		]`
	}

	if s.Name == "tns" {
		data = `[
		  { "evaluate_id": "rtgthb", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "反手如刮风", "evaluate": "打出去的球像流星，别人还没反应过来就出界了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "hefwf", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "球拍杀手小王", "evaluate": "打了三局摔坏两把拍子，商家都笑开花了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "wet67un", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "没接到过正手", "evaluate": "对手发球速度太快，我全程负责捡球", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "wethwsg", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "爱网球的阿三", "evaluate": "打球像拍蚊子，姿势全凭感觉", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  { "evaluate_id": "ytwre", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "刘小慢", "evaluate": "别人打单打，我打单人羽毛球模式，全程自己发自己接", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
		]
`
	}
	if s.Name == "gos" {
		data = `[
			{"evaluate_id": "rwt56urt", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "猴哥", "evaluate": "抽了一下午的空气", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
		  {"evaluate_id": "w7utrwt5", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "高尔夫穷人", "evaluate": "玩不起l_l", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
			{"evaluate_id": "aaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "222", "evaluate": "还行吧l_l", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
		]`
	}
	if s.Name == "hxc" {
		data = `[
  { "evaluate_id": "7u84yk", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "雪地翻滚王", "evaluate": "本来是来滑雪的，结果一路滚下山比滑还快", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "i9rgeh", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "刹车全靠脸", "evaluate": "下坡不会刹车，最后是靠撞雪人停下来的", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "78hrstg", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "单板冲浪侠", "evaluate": "滑着滑着冲进了咖啡厅，老板还问我喝不喝热可可", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "mhg91", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "雪场迷路王", "evaluate": "迷路了半小时，结果滑到了儿童初级道", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "8uf", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "防寒全靠抖", "evaluate": "装备没穿好，冷到像在跳广场舞取暖", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
]
`
	}
	if s.Name == "yjg" {
		data = `[
  { "evaluate_id": "wrt25", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "下犬式翻车王", "evaluate": "做着做着下犬式，直接变成趴地式睡觉", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "5r34r", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "体式全靠蒙", "evaluate": "老师说进入战士二式，我摆了个招财猫式", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "234u6m", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "呼吸全乱套", "evaluate": "吸气呼气配错节奏，差点原地起飞", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "929yhub", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "柔韧度负数", "evaluate": "别人能劈叉，我只能坐那像抱膝哭泣", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "0jidfnjn1", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "平板撑30秒", "evaluate": "教练说撑两分钟，我撑了三十秒就去喝水了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
]
`
	}
	if s.Name == "tqd" {
		data = `[
  { "evaluate_id": "9a0hfih5y90ha", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "黑带只会劈叉", "evaluate": "教练让我踢高一点，我直接踢飞了鞋", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "259ajsd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "腿短踢不到", "evaluate": "别人一脚踢到头，我一脚只能踢到膝盖", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "aaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "出拳像摸头", "evaluate": "打沙包力度太轻，教练说像帮它按摩", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "aaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "防守全靠脸", "evaluate": "实战时忘了抬手，结果脸吃了三脚", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "aaa", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "道服系不紧", "evaluate": "踢腿太猛裤子松了，全馆都看见了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
]`
	}
	if s.Name == "gym" {
		data = `[
  { "evaluate_id": "2149ihasd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "杠铃界郭德纲", "evaluate": "深蹲做到怀疑人生，站起来发现裤子裂了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "4ui9nhn", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "跑步机上的蜗牛", "evaluate": "本来想慢跑，结果被隔壁大爷超了三圈", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "bmk485y", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "仰卧起坐小王子", "evaluate": "做到第十个就开始想中午吃啥了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "1294u9ahd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "杠铃举不起来", "evaluate": "器械区健身帅哥太多，我只好去喝水了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "fa9u13da", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "汗水中的咸鱼", "evaluate": "出了健身房就去撸串，卡路里原地复活", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "2349yaihsd", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "哑铃界的段子手", "evaluate": "今天练二头肌，明天就要买止痛药", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "asd9u2934", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "瑜伽垫翻车王", "evaluate": "做了个下犬式，结果抽筋变成趴地式", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "vmk90167", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "跑步机电量不足", "evaluate": "跑了一公里就气喘吁吁，怀疑自己上辈子是树懒", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "z9u5612", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "杠铃举到飞起", "evaluate": "硬拉太猛，回家发现沙发都抬得轻松了", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"},
  { "evaluate_id": "324a9dus", "is_like": false, "like": 0, "group_id": "aaa-bbb-ccc-ddd-eee", "user": "游泳池边的咸鸭蛋", "evaluate": "蛙泳还没学会，先被水呛了三升", "img": "https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0", "time": "2024-04-05 16:32:15"}
]`
	}

	if err := json.Unmarshal([]byte(data), &evaData); err != nil {
		return evaData, err
	}

	return evaData, nil
}

func (s SportType) RandomNickname() string {
	nicknames := []string{
		"三分铁匠铺",
		"篮筐守护神",
		"球场捡球大队长",
		"扶墙上篮王",
		"空位不进侠",
		"篮板漏风王",
		"全场犯规制造机",
		"运球到界外",
		"传球失误艺术家",
		"球场背景板",
		"篮筐终结者（自己）",
		"球鞋打滑侠",
		"上篮撞墙王",
		"三步走成五步",
		"球场空气掌控者",
		"裁判的好朋友",
		"扣篮靠意念",
		"篮球场摄影师",
		"球衣永远是干净的",
		"上场五分钟气喘两小时",
		"罚球不进研究员",
		"球场假动作大师",
		"一秒钟掉球侠",
		"篮板永远抢不到",
		"篮筐的守门员",
		"投篮靠蒙王",
		"球场指定背锅侠",
		"关键球必失先生",
		"一条龙跑偏侠",
		"上篮打成传球",
		"空位三不沾专家",
		"跑位永远跑错",
		"防守空气专业户",
		"抢断失败大王",
		"盖帽打自己脸",
		"球场传球黑洞",
		"投篮手抖症患者",
		"上篮要扶梯",
		"全场最响的喘气声",
		"篮板跳不高",
		"球场逃跑王",
		"街球晃倒自己",
		"半场都在喊要球",
		"失误制造机",
		"篮筐铁打的兄弟",
		"运球看地板",
		"打铁声交响乐",
		"上篮弹框王",
		"球场划水大师",
	}

	// 随机取一个昵称
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(nicknames))))
	if err != nil {
		// 退化处理：如果随机失败，返回第一个
		return nicknames[0]
	}
	return nicknames[n.Int64()]
}
