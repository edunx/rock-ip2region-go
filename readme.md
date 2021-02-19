---
   磐石IP2Region封装
---

# 配置
```lua
    rock.download{
        save = "resource/ip2region.db",
        url = "http://xxxx.com/ip2region.db"
    }

    local ip = rock.ip2region{
        db = "resource/ip2region.db",
    }

    local cityid , info , err = ip.memory_search("202.96.209.133")
    print(info)
```

# 调用
```golang
    import (
        ip "githbu.com/edunx/rock-ip2region-go"
    )

    var ud  *ip.Ip2Region

    ud.Search( "127.0.0.1" )
```