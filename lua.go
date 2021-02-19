package ip2region

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

const (
	MT string = "ROCK_IP2REGION_MT"
)


func CheckIp2RegionUserdata( L *lua.LState , idx int ) *Ip2Region {
	ud := L.CheckUserData( idx )

	switch v := ud.Value.(type) {
	case *Ip2Region:
		return ud.Value.(*Ip2Region)
	default:
		L.RaiseError("expect invalid type , must be Ip2geion , got %T" , v )
		return nil
	}
}

func CreateIp2RegionUserData(L *lua.LState) int {
	opt := L.CheckTable(1)

	v := &Ip2Region{
		dbFile:  opt.CheckString("db" , "resource/ip2region.db"),
	}

	if err := v.Start(); err != nil {
		L.RaiseError("start ip2region fail , e: %v" , err)
		return 0
	}
	pub.Out.Debug("start ip2regin successful , info: %s" , v.dbFile)

	ud := L.NewUserDataByInterface( v , MT)
	L.Push(ud)
	return 1
}

func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	mt := L.NewTypeMetatable( MT )

	L.SetField(mt , "__index" , L.NewFunction(Get))
	L.SetField(mt , "__newindex" , L.NewFunction(Set))

	L.SetField(parent , "ip2region" , L.NewFunction(CreateIp2RegionUserData))
}

func Get(L *lua.LState) int {
	self := CheckIp2RegionUserdata(L , 1)
	name := L.CheckString(2)
	switch name {
	case "memory_search":
		L.Push(L.NewFunction( func (L *lua.LState) int {
			ip := L.CheckString(1)
			city , info , err := self.Search( ip )
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 3
			}
			L.Push(lua.LNumber(city))
			L.Push(lua.LString(info))
			L.Push(lua.LNil)
			return 3
		}))
		return 1
	default:
		return 0
	}

	return 0
}

func Set(L *lua.LState) int {
	return 0
}

func (this *Ip2Region) ToUserData(L *lua.LState) *lua.LUserData {
	return L.NewUserDataByInterface( this , MT )
}
