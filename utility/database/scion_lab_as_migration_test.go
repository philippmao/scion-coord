package main

import (
	"testing"
	"strconv"

	"github.com/netsec-ethz/scion-coord/config"
	"github.com/netsec-ethz/scion-coord/models"
	"github.com/netsec-ethz/scion-coord/utility"
)

func Test(t *testing.T) {
	// Initialize DB
	server := &models.SCIONLabServer{
		IA:			config.SERVER_IA,
		IP: 		config.SERVER_IP,
		LastAssignedPort:  50003,
		VPNIP:		config.SERVER_VPN_IP,
		VPNLastAssignedIP:  config.SERVER_VPN_END_IP,
	}
	err := server.Insert()
	if err != nil {
		t.Fatal(err)
	}
	var i = 0
	for i < 5{
		iString := strconv.Itoa(i)
		u, err := models.RegisterUser("ac" + iString, "ETH", "mail" + iString, "pw", iString, iString)
		if err != nil {
			t.Fatal(err)
		}
		as := &models.As{
			Isd:	1,
			As: 	i,
			Account: u.Account,
		}
		err = as.Insert()
		if err != nil{
			t.Fatal(err)
		}
		vm := &models.SCIONLabVM{
			UserEmail:    "mail" + strconv.Itoa(i),
			IP:		utility.IPIncrement(config.SERVER_VPN_START_IP, int32(i)),
			IA:		as,
			IsVPN:  true,
			RemoteIA:	config.SERVER_IA,
			RemoteIAPort:	50000+i,
		}
		err = vm.Insert()
		if err != nil{
			t.Fatal(err)
		}
		i++
	}
	i = 5
	for i< 7{
		iString := strconv.Itoa(i)
		u, err := models.RegisterUser("ac" + iString, "ETH", "mail" + iString, "pw", iString, iString)
		if err != nil {
			t.Fatal(err)
		}
		as := &models.As{
			Isd:	1,
			As: 	i,
			Account: u.Account,
		}
		err = as.Insert()
		if err != nil{
			t.Fatal(err)
		}
		vm := &models.SCIONLabVM{
			UserEmail:    "mail" + strconv.Itoa(i),
			IP:		utility.IPIncrement("1.2.3.4", int32(i)),
			IA:		as,
			IsVPN:  false,
			RemoteIA:	config.SERVER_IA,
			RemoteIAPort:	50000+i,
		}
		err = vm.Insert()
		if err != nil{
			t.Fatal(err)
		}
		i++
	}
}
