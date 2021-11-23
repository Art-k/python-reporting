package include

import (
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type DBBuilding struct {
	gorm.Model
	ExternalID                   string
	Name                         string
	Address                      string
	Propertynumber               string
	Companyid                    int
	Companyname                  string
	Numberofunits                int
	Numberofresidents            int
	Numberofusers                int
	Numberofemails               int
	Numberofinvitations          int
	CorporationnumberName        string
	CorporationnumberCompanyid   int
	CorporationnumberCompanyname string
	Disabledauthorizationat      bool
	Syncluxerone                 bool
	Snailintegration             bool
}

type InBuilding struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	//Country             interface{}   `json:"country"`
	//City                interface{}   `json:"city"`
	//Province            interface{}   `json:"province"`
	//Postalcode          interface{}   `json:"postalCode"`
	//Latitude            interface{}   `json:"latitude"`
	//Longitude           interface{}   `json:"longitude"`
	Propertynumber string `json:"propertyNumber"`
	//Propertytype        interface{}   `json:"propertyType"`
	//Backgroundimage     interface{}   `json:"backgroundImage"`
	Companyid           int    `json:"companyId"`
	Companyname         string `json:"companyName"`
	Numberofunits       int    `json:"numberOfUnits"`
	Numberofresidents   int    `json:"numberOfResidents"`
	Numberofusers       int    `json:"numberOfUsers"`
	Numberofemails      int    `json:"numberOfEmails"`
	Numberofinvitations int    `json:"numberOfInvitations"`
	//Createdat           time.Time     `json:"createdAt"`
	//Updatedat           time.Time     `json:"updatedAt"`
	//Suites              []string      `json:"suites"`
	//Contacts            []interface{} `json:"contacts"`
	//Createdby           interface{}   `json:"createdBy"`
	//Welcomemessage      interface{}   `json:"welcomeMessage"`
	//Contactus           string        `json:"contactUs"`
	//Timezone            interface{}   `json:"timeZone"`
	Corporationnumber struct {
		//Buildings                     []string      `json:"buildings"`
		//Reminders                     []interface{} `json:"reminders"`
		//Settings                      string        `json:"settings"`
		//Logouturl                     interface{}   `json:"logoutUrl"`
		//ID                            string        `json:"id"`
		Name string `json:"name"`
		//Securitycompanylogo           interface{}   `json:"securityCompanyLogo"`
		//Invitetemplate                string        `json:"inviteTemplate"`
		//Invitesubject                 interface{}   `json:"inviteSubject"`
		//Signaturetype                 interface{}   `json:"signatureType"`
		//Signatureimagepath            interface{}   `json:"signatureImagePath"`
		//Signaturetext                 interface{}   `json:"signatureText"`
		//Signaturefont                 interface{}   `json:"signatureFont"`
		//Suitecustomfields             []interface{} `json:"suiteCustomFields"`
		//Suiteusercustomfields         []interface{} `json:"suiteUserCustomFields"`
		//Residentsnotificationsettings []interface{} `json:"residentsNotificationSettings"`
		//Staffnotificationsettings     []interface{} `json:"staffNotificationSettings"`
		Companyid   int    `json:"companyId"`
		Companyname string `json:"companyName"`
	} `json:"corporationNumber"`
	//Purpose                 string      `json:"purpose"`
	Disabledauthorizationat interface{} `json:"disabledAuthorizationAt"`
	Syncluxerone            bool        `json:"syncLuxerone"`
	Snailintegration        bool        `json:"snailIntegration"`
	//Securitycompanylogo     interface{} `json:"securityCompanyLogo"`
	//Cnsettings              struct {
	//	Petfriendly                    bool `json:"petFriendly"`
	//	Useresidentinstruction         bool `json:"useResidentInstruction"`
	//	Useemailwaiver                 bool `json:"useEmailWaiver"`
	//	Useparcelswaiver               bool `json:"useParcelsWaiver"`
	//	Usekeywaiver                   bool `json:"useKeyWaiver"`
	//	Autoinviteresidents            bool `json:"autoInviteResidents"`
	//	Autoapproveresidentdatachanges bool `json:"autoApproveResidentDataChanges"`
	//	Usedateofbirth                 bool `json:"useDateOfBirth"`
	//} `json:"cNSettings"`
}

type InBuildingResponse struct {
	Total    int          `json:"total"`
	Entities []InBuilding `json:"entities"`
}

func SyncBuildings() {

	Message := "<p>Hi</p><br>Here is a list of MaxCondoClub changes: <br>"
	isAdded := false
	isChanged := false

	allBuildings := GetAllMCCBuildings()
	Log.Info("Buildings received")
	for _, rec := range allBuildings {

		var exBuilding DBBuilding
		db.Where("external_id = ?", rec.ID).Find(&exBuilding)

		if exBuilding.ID == 0 {
			newBld := DBBuilding{
				ExternalID:                   rec.ID,
				Name:                         rec.Name,
				Address:                      rec.Address,
				Propertynumber:               rec.Propertynumber,
				Companyid:                    rec.Companyid,
				Companyname:                  rec.Companyname,
				Numberofunits:                rec.Numberofunits,
				Numberofresidents:            rec.Numberofresidents,
				Numberofusers:                rec.Numberofusers,
				Numberofemails:               rec.Numberofemails,
				Numberofinvitations:          rec.Numberofinvitations,
				CorporationnumberName:        rec.Corporationnumber.Name,
				CorporationnumberCompanyid:   rec.Corporationnumber.Companyid,
				CorporationnumberCompanyname: rec.Corporationnumber.Companyname,
				Disabledauthorizationat:      false,
				Syncluxerone:                 rec.Syncluxerone,
				Snailintegration:             rec.Snailintegration,
			}
			if rec.Disabledauthorizationat != nil {
				newBld.Disabledauthorizationat = true
			}

			db.Create(&newBld)

			Message += "<br><br><i>Building added</i>"
			Message += "<br>Building Name :<b>" + rec.Name + "</b><br>"
			Message += "Building Address :<b>" + rec.Address + "</b><br>"
			Message += "Company Name :<b>" + rec.Companyname + "</b><br>"
			Message += "Corporation :<b>" + rec.Corporationnumber.Name + "</b><br><br>"

			isAdded = true

			Log.Infof("Building ADDED '%d'", newBld.ExternalID)

		} else {

			changed := false
			tmpMessage := "<br><br><i>Building changed '" + exBuilding.Address + "'</i><br>"

			if exBuilding.Name != rec.Name {
				tmpMessage += "Building Name :<b>" + exBuilding.Name + "</b> -> <b>" + rec.Name + "</b><br>"
				exBuilding.Name = rec.Name
				changed = true
			}

			if exBuilding.Address != rec.Address {
				tmpMessage += "Building Address :<b>" + exBuilding.Address + "</b> -> <b>" + rec.Address + "</b><br>"
				exBuilding.Address = rec.Address
				changed = true
			}

			if exBuilding.Companyname != rec.Companyname {
				tmpMessage += "Company Name :<b>" + exBuilding.Companyname + "</b> -> <b>" + rec.Companyname + "</b><br>"
				exBuilding.Companyname = rec.Companyname
				changed = true
			}

			if exBuilding.CorporationnumberName != rec.Corporationnumber.Name {
				tmpMessage += "Company Name :<b>" + exBuilding.CorporationnumberName + "</b> -> <b>" + rec.Corporationnumber.Name + "</b><br><br>"
				exBuilding.CorporationnumberName = rec.Corporationnumber.Name
				changed = true
			}

			if !exBuilding.Disabledauthorizationat && rec.Disabledauthorizationat != nil {
				tmpMessage += "Authorisation is <b>disabled</b>"
				exBuilding.Disabledauthorizationat = true
				changed = true
			}

			if exBuilding.Disabledauthorizationat && rec.Disabledauthorizationat == nil {
				tmpMessage += "Authorisation is <b>enabled</b>"
				exBuilding.Disabledauthorizationat = false
				changed = true
			}

			db.Save(&exBuilding)

			if changed {
				Log.Infof("Building changed '%d'", exBuilding.ExternalID)
				Message += tmpMessage
				isChanged = true
			}

		}
	}

	if isAdded || isChanged {
		SendEmailOAUTH2(os.Getenv("BUILDINGS_CHANGED1"), "MaxCondoClub Building list changed ", Message)
		SendEmailOAUTH2(os.Getenv("BUILDINGS_CHANGED2"), "MaxCondoClub Building list changed ", Message)
		SendEmailOAUTH2(os.Getenv("BUILDINGS_CHANGED3"), "MaxCondoClub Building list changed ", Message)
		Log.Infof("Messages sent!")
	}

	Log.Infof("Compare DONE!")

}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func GetMCCToken() Token {

	var token Token

	var client *http.Client
	client = &http.Client{}

	URL := "https://auth.maxcondoclub.com/oauth/v2/token?username=" + os.Getenv("MCC_LOGIN") + "&password=" + os.Getenv("MCC_PASSWORD") + "&grant_type=password&client_id=" + os.Getenv("MCC_CLIENT_ID") + "&client_secret=" + os.Getenv("MCC_CLIENT_SECRET")
	req, _ := http.NewRequest("GET", URL, nil)
	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		//Log.Error("[UPDATE HIVESTACK INVENTORY]", err, "ERROR HIVESTACK GET UNITS")
		return token
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &token)
	if err != nil {
		//Log.Error("UpdateInternalHivestackInventory", err)
		return token
	}

	return token
}

func GetAllMCCBuildings() []InBuilding {

	var client *http.Client
	client = &http.Client{}
	var allBuilding []InBuilding

	token := GetMCCToken()
	page := 1
	for {
		pageStr := strconv.Itoa(page)
		page += 1
		URL := "https://building.maxcondoclub.com/api/v1/buildings.json?expand=corporationNumber&page=" + pageStr
		req, _ := http.NewRequest("GET", URL, nil)

		req.Header.Add("Authorization", "Bearer "+token.AccessToken)

		resp, err := client.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			//Log.Error("[UPDATE HIVESTACK INVENTORY]", err, "ERROR HIVESTACK GET UNITS")
			break
		}

		body, err := ioutil.ReadAll(resp.Body)

		var buildingsResponse InBuildingResponse
		err = json.Unmarshal(body, &buildingsResponse)
		if err != nil {
			//Log.Error("UpdateInternalHivestackInventory", err)
			break
		}

		for _, rec := range buildingsResponse.Entities {
			allBuilding = append(allBuilding, rec)
		}

		if len(allBuilding) >= buildingsResponse.Total {
			break
		}

	}

	return allBuilding
}
