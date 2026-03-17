# BreakReminder

macOS 菜单栏休息提醒工具。定时提醒你站起来活动身体，保护健康。

## 功能

- **定时提醒** — 自定义间隔（10~120 分钟），到时弹窗提醒休息
- **30 种休息活动** — 拉伸、护眼、运动、呼吸四大类，随机推荐
- **息屏自动暂停** — 锁屏 / 合盖时自动暂停计时，唤醒后自动恢复
- **菜单栏常驻** — 实时显示倒计时，支持暂停、跳过、立即休息
- **开机自启** — 可选开机自动启动

## 截图

<p align="center">
  <img src="assets/icon.png" width="128" alt="App Icon">
</p>

## 安装

### 从源码构建

需要 Go 1.21+ 和 macOS 11.0+。

```bash
git clone https://github.com/shentry/BreakReminder.git
cd BreakReminder
make package
```

构建完成后将 `BreakReminder.app` 拖入「应用程序」文件夹即可。

### 直接运行

```bash
make run
```

## 使用

启动后应用会常驻菜单栏，显示距下次休息的倒计时。

| 菜单项 | 说明 |
|--------|------|
| 暂停 / 继续 | 手动暂停或恢复计时 |
| 立即休息 | 跳过等待，立刻开始休息 |
| 设置 | 调整间隔、时长、通知方式等 |

**息屏暂停**：锁屏（⌃⌘Q）或屏幕休眠时，计时器自动暂停；解锁后自动恢复，不会浪费你的休息间隔。

## 配置

配置文件位于 `~/.breakreminder/config.json`：

| 选项 | 默认值 | 说明 |
|------|--------|------|
| interval_minutes | 30 | 提醒间隔（分钟） |
| break_duration_sec | 300 | 休息时长（秒） |
| notification_style | both | 通知方式：system / popup / both |
| launch_at_login | false | 开机自启 |
| sound_enabled | true | 提示音 |

## 技术栈

- **Go** + **Fyne v2.7.2** — 跨平台 GUI 框架
- **CGo** — 调用 macOS 原生 API 检测屏幕锁定 / 唤醒
- **systray** — 菜单栏集成

## License

MIT
