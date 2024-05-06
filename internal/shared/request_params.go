package shared

import (
	"encoding/json"
	"fmt"
	"ps-cats-social/pkg/base/app"
)

func ExtractUserId(ctx *app.Context) (int64, error) {
	user := ctx.Context().Value("user_info")
	jsonData, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal context value. error: %v", err)
	}

	var userInfo map[string]interface{}
	err = json.Unmarshal(jsonData, &userInfo)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal json. error: %v", err)
	}

	userId, ok := userInfo["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_info")
	}

	return int64(userId), nil
}
