# 移动魔百盒直播源转m3u8播放列表

## 使用方法

### 1. 获取频道列表接口地址和基础URL

参照 [@zhantong](https://github.com/zhantong) 的文章 [苏州移动IPTV抓包](https://www.polarxiong.com/archives/苏州移动IPTV抓包.html)

### 2. 修改 `main.go` 中的 `CHANNEL_API_URL` 和 `BASE_URL`

### 3. 编译运行

`go build`

`./ottcn2m3u8`

会在当前目录下生成 channel.m3u8
