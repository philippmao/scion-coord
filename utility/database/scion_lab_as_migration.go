package main

import (
	"fmt"

	"github.com/netsec-ethz/scion-coord/config"
	"github.com/netsec-ethz/scion-coord/models"
	"github.com/netsec-ethz/scion/go/lib/addr"
)

func migrateDB() error {
	// Migrate ScionLabserver 1-7 to the Database
	sLabServer, err := models.FindSCIONLabServer(config.SERVER_IA)
	if err != nil {
		return err
	}
	ia, err := addr.IAFromString(config.SERVER_IA)
	if err != nil {
		return err
	}
	apServer := &models.AttachmentPoint{
		VPNIP:      config.SERVER_VPN_IP,
		StartVPNIP: config.SERVER_VPN_START_IP,
		EndVPNIP:   config.SERVER_VPN_END_IP,
	}
	err = apServer.Insert()
	if err != nil {
		return err
	}
	slasServer := &models.SCIONLabAS{
		UserMail:  "", // TODO Admin UserMail ?
		PublicIP:  sLabServer.IP,
		StartPort: config.SERVER_START_PORT,
		ISD:       ia.I,
		AS:        ia.A,
		Status:    models.CREATE,
		Type:      models.SERVER,
		AP:        apServer,
	}
	err = slasServer.Insert()
	if err != nil {
		return err
	}
	sLabVMs, err := models.FindSCIONLabVMsByRemoteIA(config.SERVER_IA)
	if err != nil {
		return fmt.Errorf("Error looking for VMs IA: %v", config.SERVER_IA)
	}
	var ip string
	for _, sLabVM := range sLabVMs {
		// add SCIONLabAS entry to database
		if sLabVM.IsVPN {
			// VM is connected via VPN -> no public IP
			ip = "0.0.0.0"
		} else {
			ip = sLabVM.IP
		}
		slas := &models.SCIONLabAS{
			UserMail:  sLabVM.UserEmail,
			PublicIP:  ip,
			StartPort: config.SERVER_START_PORT,
			ISD:       sLabVM.IA.Isd,
			AS:        sLabVM.IA.As,
			Status:    models.CREATE,
			Type:      models.VM,
		}
		err = slas.Insert()
		if err != nil {
			return fmt.Errorf("Error inserting slas: %v", slas)
		}
		var respondIP string
		if sLabVM.IsVPN {
			respondIP = config.SERVER_VPN_IP
		} else {
			respondIP = config.SERVER_IP
		}
		// TODO what are the BrIds of the SlabServer ?
		var ServerBRIDoffset = 0
		// add Connection for that VM
		connection := models.Connection{
			JoinIP:        sLabVM.IP,
			RespondIP:     respondIP,
			JoinAS:        slas,
			RespondAP:     apServer,
			JoinBRID:      1,
			RespondBRID:   sLabVM.RemoteIAPort - 50001 + ServerBRIDoffset,
			Linktype:      models.PARENT,
			IsVPN:         sLabVM.IsVPN,
			JoinStatus:    models.ACTIVE,
			RespondStatus: models.ACTIVE,
		}
		err = connection.Insert()
		if err != nil {
			return fmt.Errorf("Error inserting Connection: %v", connection)
		}
	}
	return nil
}

func main() {
	err := migrateDB()
	if err != nil {
		fmt.Println("Error migrating DB: %v", err)
	} else {
		fmt.Println("Succesfully migrated DB")
	}
}
