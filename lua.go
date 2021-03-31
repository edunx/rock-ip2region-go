package ip2region

import (
	"fmt"
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

func (self *Ip2Region) debug(L *lua.LState , args *lua.Args) lua.LValue {
	n := args.Len()
	if n <= 0 {
		return lua.LNil
	}

	for i := 1 ; i<=n ; i++ {
		ip := args.CheckString(L , i)
		city , info , err := self.Search( ip )
		if err != nil {
			return lua.LString( err.Error() )
		}
		fmt.Printf("ip: %s city: %d , info: %s\n" ,ip , city , info)
	}
	return lua.LNil
}

func (self *Ip2Region) Index(L *lua.LState , key string ) lua.LValue {
	if key == "debug" { return lua.NewGFunction(self.debug) }
	return lua.LNil
}

func (self *Ip2Region) ToLightUserData(L *lua.LState) *lua.LightUserData {
	return L.NewLightUserData( self )
}

func createRegionLightUserData(L *lua.LState , args *lua.Args ) lua.LValue {
	path := args.CheckString(L , 1)
	v := &Ip2Region{ dbFile: path }

	if err := v.Start(); err != nil {
		L.RaiseError("start ip2region fail , e: %v" , err)
		return lua.LNil
	}
	pub.Out.Debug("start ip2regin successful , info: %s" , v.dbFile)
	return v.ToLightUserData(L)
}

func LuaInjectApi(L *lua.LState , parent *lua.UserKV) {
	parent.Set("ip2region" , lua.NewGFunction( createRegionLightUserData ) )
}