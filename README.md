# 移动魔百盒直播源生成m3u8播放列表

## 使用方法

参照 [@zhantong](https://github.com/zhantong) 的文章 [苏州移动IPTV抓包](https://www.polarxiong.com/archives/苏州移动IPTV抓包.html) 获取频道列表接口地址和基础URL

    Usage of ottcn2m3u8:
      -api string
        	API URL to fetch channel list. (default "http://looktvepg.jsa.bcs.ottcn.com:8080/ysten-lvoms-epg/epg/getChannelIndexs.shtml?deviceGroupId=1697")
      -base string
        	Base URL for stream. (default "http://183.207.248.71:80/cntv/live1")
      -h	This help.
      -o string
        	Output file path. (default "channel.m3u8")
      -v	Verbose mode.
