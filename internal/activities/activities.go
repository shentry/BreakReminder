package activities

import (
	"math/rand"
)

type Category string

const (
	Stretch   Category = "拉伸"
	Eyes      Category = "护眼"
	Movement  Category = "运动"
	Breathing Category = "呼吸"
)

type Activity struct {
	Name        string
	Description string
	Category    Category
	DurationSec int
}

var activities = []Activity{
	// 拉伸 (8)
	{Name: "颈部环绕", Description: "缓慢转动头部画圈，顺时针、逆时针各 5 次。", Category: Stretch, DurationSec: 60},
	{Name: "耸肩放松", Description: "双肩向上耸至耳朵位置，保持 5 秒后放下，重复 10 次。", Category: Stretch, DurationSec: 60},
	{Name: "扩胸运动", Description: "双手在背后交握，伸直手臂微微上抬，保持 20 秒。", Category: Stretch, DurationSec: 30},
	{Name: "坐姿脊柱扭转", Description: "坐直身体，将上身扭向右侧保持 15 秒，再换左侧。", Category: Stretch, DurationSec: 45},
	{Name: "手腕画圈", Description: "伸出双臂，缓慢转动手腕，顺时针、逆时针各 10 圈。", Category: Stretch, DurationSec: 30},
	{Name: "上背拉伸", Description: "十指交叉，掌心向前推出，弓起上背部，保持 20 秒。", Category: Stretch, DurationSec: 30},
	{Name: "侧颈拉伸", Description: "轻轻将头倾向一侧，保持 15 秒后换边，每侧重复 2 次。", Category: Stretch, DurationSec: 60},
	{Name: "前臂拉伸", Description: "伸出一只手臂，用另一只手轻轻将手指向后拉，每臂 15 秒。", Category: Stretch, DurationSec: 45},

	// 护眼 (7)
	{Name: "20-20-20 法则", Description: "看向 6 米外的物体，持续 20 秒，缓解眼部疲劳。", Category: Eyes, DurationSec: 30},
	{Name: "眼球转动", Description: "缓慢转动眼球画大圈，顺时针 5 次、逆时针 5 次。", Category: Eyes, DurationSec: 30},
	{Name: "掌心热敷", Description: "双手搓热后轻轻覆盖在闭合的双眼上，保持 30 秒。", Category: Eyes, DurationSec: 45},
	{Name: "焦点切换", Description: "将笔举到一臂远处注视，然后注视远处物体，交替 10 次。", Category: Eyes, DurationSec: 30},
	{Name: "快速眨眼", Description: "快速眨眼 20 次，然后闭眼放松 20 秒。", Category: Eyes, DurationSec: 30},
	{Name: "8 字眼操", Description: "用眼睛追踪一个大大的 8 字形，每个方向 5 次。", Category: Eyes, DurationSec: 30},
	{Name: "远眺窗外", Description: "望向窗外，注视你能看到的最远处物体，保持 30 秒。", Category: Eyes, DurationSec: 30},

	// 运动 (8)
	{Name: "站立拉伸", Description: "站起来，双手向天花板伸展，保持 10 秒，重复 3 次。", Category: Movement, DurationSec: 45},
	{Name: "起身走走", Description: "在房间里走一走，顺便去倒杯水吧！", Category: Movement, DurationSec: 120},
	{Name: "踮脚运动", Description: "站立踮起脚尖，短暂保持后放下，重复 15 次。", Category: Movement, DurationSec: 45},
	{Name: "桌面俯卧撑", Description: "双手撑在桌边，做 10 个斜面俯卧撑。", Category: Movement, DurationSec: 45},
	{Name: "腿部摆动", Description: "扶住稳定物体，一条腿前后摆动 10 次，换腿重复。", Category: Movement, DurationSec: 60},
	{Name: "开合跳", Description: "做 20 个开合跳，让血液流动起来！", Category: Movement, DurationSec: 45},
	{Name: "靠墙静蹲", Description: "背靠墙壁，膝盖弯曲 90 度，保持 30 秒。", Category: Movement, DurationSec: 45},
	{Name: "原地踏步", Description: "原地高抬腿踏步 60 秒。", Category: Movement, DurationSec: 60},

	// 呼吸 (7)
	{Name: "方块呼吸", Description: "吸气 4 秒、屏息 4 秒、呼气 4 秒、屏息 4 秒，重复 4 轮。", Category: Breathing, DurationSec: 75},
	{Name: "4-7-8 呼吸法", Description: "吸气 4 秒、屏息 7 秒、呼气 8 秒，重复 3 次。", Category: Breathing, DurationSec: 60},
	{Name: "腹式深呼吸", Description: "手放腹部，深吸气使腹部鼓起，吸 5 秒呼 5 秒，重复 6 次。", Category: Breathing, DurationSec: 60},
	{Name: "交替鼻孔呼吸", Description: "按住右鼻孔从左侧吸气，再按住左侧从右侧呼气，交替 5 轮。", Category: Breathing, DurationSec: 60},
	{Name: "狮子呼吸", Description: "深吸一口气，然后张大嘴、伸出舌头用力呼出，重复 5 次。", Category: Breathing, DurationSec: 30},
	{Name: "数息练习", Description: "自然呼吸并数每次呼气，数到 10 重新开始，持续 2 分钟。", Category: Breathing, DurationSec: 120},
	{Name: "叹气放松", Description: "用鼻子深吸一口气，然后发出长长的叹气声呼出，重复 5 次。", Category: Breathing, DurationSec: 30},
}

func All() []Activity {
	return activities
}

func Random() Activity {
	return activities[rand.Intn(len(activities))]
}

func RandomFromCategory(cat Category) Activity {
	var filtered []Activity
	for _, a := range activities {
		if a.Category == cat {
			filtered = append(filtered, a)
		}
	}
	if len(filtered) == 0 {
		return Random()
	}
	return filtered[rand.Intn(len(filtered))]
}
