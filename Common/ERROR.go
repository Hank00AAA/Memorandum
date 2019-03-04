package Common

import "github.com/pkg/errors"

var(
	ERR_NO_LOCAL_ANY_IP_FOUND = errors.New("没有网卡IP")
	ERR_NODE_NUMBER_OVER_MAX = errors.New("节点id超出最大范围")
)
