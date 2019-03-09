package Common

import "github.com/pkg/errors"

var(
	ERR_NO_LOCAL_ANY_IP_FOUND = errors.New("没有网卡IP")
	ERR_NODE_NUMBER_OVER_MAX = errors.New("节点id超出最大范围")
	ERR_NO_FOUND_ACCOUNT = errors.New("帐号密码不对应")
	ERR_ACCOUNT_IS_NIL = errors.New("帐号为空")
	ERR_ACCOUNT_DONT_EXIST = errors.New("不存在用户")
	ERR_MULTI_EMAIL_EXIST  = errors.New("存在多个相同email的用户")
)
