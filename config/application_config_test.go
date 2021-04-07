package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	"github.com/coyim/otr3"
	. "gopkg.in/check.v1"
)

type AccountsSuite struct{}

var _ = Suite(&AccountsSuite{})

func (s *AccountsSuite) Test_Accounts_RemoveAccount(c *C) {
	ac1 := &Account{Account: "account@one.com"}
	ac2 := &Account{Account: "account@two.com"}
	acs := ApplicationConfig{
		Accounts: []*Account{ac1, ac2},
	}

	acs.Remove(ac1)

	c.Check(len(acs.Accounts), Equals, 1)
	_, found := acs.GetAccount("account@two.com")
	c.Check(found, Equals, true)
}

func (s *AccountsSuite) Test_Accounts_DontRemoveWhenDoesntExist(c *C) {
	ac1 := &Account{Account: "account@one.com"}
	ac2 := &Account{Account: "account@two.com"}
	ac3 := &Account{Account: "nohay@anywhere.com"}
	acs := ApplicationConfig{
		Accounts: []*Account{ac1, ac2},
	}

	acs.Remove(ac3)

	c.Check(len(acs.Accounts), Equals, 2)
	_, found := acs.GetAccount("account@two.com")
	c.Check(found, Equals, true)
}

func (s *AccountsSuite) Test_Account_ByAccountNameAlphabetic(c *C) {
	ac1 := &Account{Account: "account@one.com"}
	ac2 := &Account{Account: "xccount@two.com"}
	ac3 := &Account{Account: "nohay@anywhere.com"}
	one := []*Account{ac1, ac2, ac3}

	sort.Sort(ByAccountNameAlphabetic(one))

	c.Assert(one[0], Equals, ac1)
	c.Assert(one[1], Equals, ac3)
	c.Assert(one[2], Equals, ac2)
}

func (s *AccountsSuite) Test_ApplicationConfig_serialize(c *C) {
	a := &ApplicationConfig{}

	res, e := a.serialize()
	c.Assert(e, IsNil)
	c.Assert(string(res), Equals, ""+
		"{\n"+
		"\t\"Accounts\": null,\n"+
		"\t\"Bell\": false,\n"+
		"\t\"ConnectAutomatically\": false,\n"+
		"\t\"Display\": {\n"+
		"\t\t\"MergeAccounts\": false,\n"+
		"\t\t\"ShowOnlyOnline\": false,\n"+
		"\t\t\"HideFeedbackBar\": false,\n"+
		"\t\t\"ShowOnlyConfirmed\": false,\n"+
		"\t\t\"SortByStatus\": false\n"+
		"\t},\n"+
		"\t\"AdvancedOptions\": false,\n"+
		"\t\"UniqueConfigurationID\": \"\"\n"+
		"}")
}

func (s *AccountsSuite) Test_ApplicationConfig_UpdateToLatestVersion_updatesAllAccounts(c *C) {
	a := &ApplicationConfig{
		Accounts: []*Account{
			&Account{},
			&Account{
				LegacyKnownFingerprints: []KnownFingerprint{
					KnownFingerprint{
						UserID:      "one@some.org",
						Fingerprint: []byte{0x01, 0x02, 0x03},
						Untrusted:   true,
					},
					KnownFingerprint{
						UserID:      "ignored@fingerprint.com",
						Fingerprint: []byte{},
						Untrusted:   true,
					},
					KnownFingerprint{
						UserID:      "one@some.org",
						Fingerprint: []byte{0x02, 0x02, 0x05},
						Untrusted:   false,
					},
				},
			},
		},
	}

	c.Assert(a.UpdateToLatestVersion(), Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_GetAccount_ReturnsNotOKForWrongAccount(c *C) {
	a := &ApplicationConfig{
		Accounts: []*Account{
			&Account{Account: "some@one.com"},
		},
	}

	res, ok := a.GetAccount("another@one.com")
	c.Assert(ok, Equals, false)
	c.Assert(res, IsNil)
}

func (s *AccountsSuite) Test_ApplicationConfig_AddNewAccount_returnsErrorIfEncountered(c *C) {
	originalFn := generateMissingKeysFunc
	generateMissingKeysFunc = func([][]byte) ([]otr3.PrivateKey, error) {
		return nil, errors.New("an unexpected IO error")
	}

	defer func() {
		generateMissingKeysFunc = originalFn
	}()

	a := &ApplicationConfig{}

	res, e := a.AddNewAccount()
	c.Assert(res, IsNil)
	c.Assert(e, ErrorMatches, "an unexpected IO error")
}

func (s *AccountsSuite) Test_ApplicationConfig_AddNewAccount_addsANewAccount(c *C) {
	a := &ApplicationConfig{}

	res, e := a.AddNewAccount()
	c.Assert(e, IsNil)
	c.Assert(a.Accounts, HasLen, 1)
	c.Assert(a.Accounts[0], Equals, res)
}

func (s *AccountsSuite) Test_ApplicationConfig_HasEncryptedStorage(c *C) {
	c.Assert((&ApplicationConfig{shouldEncrypt: false}).HasEncryptedStorage(), Equals, false)
	c.Assert((&ApplicationConfig{shouldEncrypt: true}).HasEncryptedStorage(), Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_onBeforeSave_generatesUniqueIDIfNeeded(c *C) {
	a := &ApplicationConfig{}
	a.onBeforeSave()
	c.Assert(a.UniqueConfigurationID, Not(Equals), "")
	old := a.UniqueConfigurationID
	a.onBeforeSave()
	c.Assert(a.UniqueConfigurationID, Equals, old)
}

func (s *AccountsSuite) Test_ApplicationConfig_GetUniqueID_generatesIfNecessary(c *C) {
	a := &ApplicationConfig{}
	v := a.GetUniqueID()
	c.Assert(a.UniqueConfigurationID, Not(Equals), "")
	c.Assert(a.UniqueConfigurationID, Equals, v)

	a.UniqueConfigurationID = "hello"

	c.Assert(a.GetUniqueID(), Equals, "hello")
	c.Assert(a.UniqueConfigurationID, Equals, "hello")
}

func (s *AccountsSuite) Test_ApplicationConfig_WhenLoaded_callsTheFunctionDirectlyIfCalledOnAnExistingConfig(c *C) {
	a := &ApplicationConfig{}
	called := false
	a.WhenLoaded(func(*ApplicationConfig) {
		called = true
	})

	c.Assert(called, Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_WhenLoaded_addsTheFunctionToLoadLaterIfCalledOnNilConfig(c *C) {
	previous := loadEntries
	loadEntries = []func(*ApplicationConfig){}
	defer func() {
		loadEntries = previous
	}()

	var a *ApplicationConfig = nil

	called := false
	a.WhenLoaded(func(*ApplicationConfig) {
		called = true
	})
	c.Assert(called, Equals, false)
	c.Assert(loadEntries, HasLen, 1)
	loadEntries[0](nil)
	c.Assert(called, Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_accountLoaded_callsAllFunctionsWhenCalled(c *C) {
	previous := loadEntries
	defer func() {
		loadEntries = previous
	}()

	called1 := false
	called2 := false

	wg := sync.WaitGroup{}
	wg.Add(2)

	loadEntries = []func(*ApplicationConfig){
		func(*ApplicationConfig) {
			called1 = true
			wg.Done()
		},
		func(*ApplicationConfig) {
			called2 = true
			wg.Done()
		},
	}

	a := &ApplicationConfig{}
	a.accountLoaded()
	wg.Wait()
	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, true)
	c.Assert(loadEntries, HasLen, 0)
}

func (s *AccountsSuite) Test_ApplicationConfig_doAfterSave(c *C) {
	a := &ApplicationConfig{}
	called := false
	a.doAfterSave(func() {
		called = true
	})

	c.Assert(a.afterSave, HasLen, 1)
	a.afterSave[0]()
	c.Assert(called, Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_onAfterSave(c *C) {
	a := &ApplicationConfig{}
	called1 := false
	called2 := false
	a.doAfterSave(func() {
		called1 = true
	})
	a.doAfterSave(func() {
		called2 = true
	})

	a.onAfterSave()
	c.Assert(a.afterSave, IsNil)
	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_savesWithoutEncryption(c *C) {
	tmpfileName := generateTempFileName()

	a := &ApplicationConfig{
		shouldEncrypt: false,
		filename:      tmpfileName,
	}

	e := a.Save(nil)
	defer os.Remove(tmpfileName)

	c.Assert(e, IsNil)

	content, _ := ioutil.ReadFile(a.filename)
	c.Assert(string(content), Equals, fmt.Sprintf(""+
		"{\n"+
		"\t\"Accounts\": null,\n"+
		"\t\"Bell\": false,\n"+
		"\t\"ConnectAutomatically\": false,\n"+
		"\t\"Display\": {\n"+
		"\t\t\"MergeAccounts\": false,\n"+
		"\t\t\"ShowOnlyOnline\": false,\n"+
		"\t\t\"HideFeedbackBar\": false,\n"+
		"\t\t\"ShowOnlyConfirmed\": false,\n"+
		"\t\t\"SortByStatus\": false\n"+
		"\t},\n"+
		"\t\"AdvancedOptions\": false,\n"+
		"\t\"UniqueConfigurationID\": \"%s\"\n"+
		"}", a.UniqueConfigurationID))
}

type mockKeySupplier struct {
	generateKey func(EncryptionParameters) ([]byte, []byte, bool)
}

func (m *mockKeySupplier) GenerateKey(params EncryptionParameters) ([]byte, []byte, bool) {
	return m.generateKey(params)
}

func (m *mockKeySupplier) Invalidate() {
}

func (m *mockKeySupplier) LastAttemptFailed() {
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_savesWithEncryption_withNewParameters(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	os.Remove(tmpfile.Name())

	a := &ApplicationConfig{
		shouldEncrypt:         true,
		filename:              tmpfile.Name(),
		UniqueConfigurationID: "123hello",
	}

	e := a.Save(&mockKeySupplier{
		generateKey: func(EncryptionParameters) ([]byte, []byte, bool) {
			return testKey, testMacKey, true
		},
	})
	defer os.Remove(a.filename)

	c.Assert(e, IsNil)
	c.Assert(a.filename, Equals, tmpfile.Name()+".enc")

	content, _ := ioutil.ReadFile(a.filename)

	ed, e2 := parseEncryptedData(content)
	c.Assert(e2, IsNil)
	c.Assert(ed.Data, Not(Equals), "")
	c.Assert(ed.Params.Nonce, Equals, a.params.Nonce)
	c.Assert(ed.Params.Salt, Equals, a.params.Salt)
	c.Assert(ed.Params.N, Equals, a.params.N)
	c.Assert(ed.Params.P, Equals, a.params.P)
	c.Assert(ed.Params.R, Equals, a.params.R)
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_savesWithEncryption_withExistingParameters(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	os.Remove(tmpfile.Name())

	a := &ApplicationConfig{
		shouldEncrypt:         true,
		filename:              tmpfile.Name(),
		UniqueConfigurationID: "123hello",
	}

	a.params = &EncryptionParameters{
		Nonce:         "dbd8f7642b05349123d59d1b",
		Salt:          "E18CB93A823465D2797539EBC5F3C0FD",
		nonceInternal: testNonce,
		saltInternal:  testSalt,
		N:             242144,
		R:             4,
		P:             2,
	}

	e := a.Save(&mockKeySupplier{
		generateKey: func(EncryptionParameters) ([]byte, []byte, bool) {
			return testKey, testMacKey, true
		},
	})
	defer os.Remove(a.filename)

	c.Assert(e, IsNil)
	c.Assert(a.filename, Equals, tmpfile.Name()+".enc")
	c.Assert(a.params.Salt, Equals, "e18cb93a823465d2797539ebc5f3c0fd")
	c.Assert(a.params.N, Equals, 242144)
	c.Assert(a.params.R, Equals, 4)
	c.Assert(a.params.P, Equals, 2)
	c.Assert(a.params.Nonce, Not(Equals), "dbd8f7642b05349123d59d1b")

	content, _ := ioutil.ReadFile(a.filename)

	ed, e2 := parseEncryptedData(content)
	c.Assert(e2, IsNil)
	c.Assert(ed.Data, Not(Equals), "")
	c.Assert(ed.Params.Nonce, Equals, a.params.Nonce)
	c.Assert(ed.Params.Salt, Equals, "e18cb93a823465d2797539ebc5f3c0fd")
	c.Assert(ed.Params.N, Equals, 242144)
	c.Assert(ed.Params.P, Equals, 2)
	c.Assert(ed.Params.R, Equals, 4)
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_savesWithEncryption_doesntAddExtensionIfNotNecessary(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	os.Remove(tmpfile.Name())

	a := &ApplicationConfig{
		shouldEncrypt:         true,
		filename:              tmpfile.Name() + ".enc",
		UniqueConfigurationID: "123hello",
	}

	e := a.Save(&mockKeySupplier{
		generateKey: func(EncryptionParameters) ([]byte, []byte, bool) {
			return testKey, testMacKey, true
		},
	})
	defer os.Remove(a.filename)

	c.Assert(e, IsNil)
	c.Assert(a.filename, Equals, tmpfile.Name()+".enc")

	content, _ := ioutil.ReadFile(a.filename)

	ed, e2 := parseEncryptedData(content)
	c.Assert(e2, IsNil)
	c.Assert(ed.Data, Not(Equals), "")
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_failsOnSerialization(c *C) {
	orgJSONMarshalIndentFn := jsonMarshalIndentFunc
	defer func() {
		jsonMarshalIndentFunc = orgJSONMarshalIndentFn
	}()

	jsonMarshalIndentFunc = func(v interface{}, v2, v3 string) ([]byte, error) {
		return nil, errors.New("ser went wrong")
	}

	a := &ApplicationConfig{}

	e := a.Save(nil)
	c.Assert(e, ErrorMatches, "ser went wrong")
}

func (s *AccountsSuite) Test_ApplicationConfig_Save_savesWithEncryption_failsIfEncryptionFails(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	os.Remove(tmpfile.Name())

	a := &ApplicationConfig{
		shouldEncrypt:         true,
		filename:              tmpfile.Name() + ".enc",
		UniqueConfigurationID: "123hello",
	}

	e := a.Save(&mockKeySupplier{
		generateKey: func(EncryptionParameters) ([]byte, []byte, bool) {
			return nil, nil, false
		},
	})
	defer os.Remove(a.filename)

	c.Assert(e, ErrorMatches, "no password supplied, aborting")
}

func (s *AccountsSuite) Test_ApplicationConfig_turnOnEncryption_doesntDoAnythingIfAlreadyEncrypted(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: true,
	}
	res := a.turnOnEncryption()
	c.Assert(res, Equals, false)
}

func (s *AccountsSuite) Test_ApplicationConfig_turnOnEncryption_turnsOnEncryptionButNothingElseWhenFilenameAlreadyHasSuffix(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: false,
		filename:      "test1.config.enc",
	}
	res := a.turnOnEncryption()
	c.Assert(res, Equals, true)
	c.Assert(a.shouldEncrypt, Equals, true)
	c.Assert(a.filename, Equals, "test1.config.enc")
}

func (s *AccountsSuite) Test_ApplicationConfig_turnOnEncryption_turnsOnEncryptionAndChangesFilename(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: false,
		filename:      "test1.config",
	}
	res := a.turnOnEncryption()
	c.Assert(res, Equals, true)
	c.Assert(a.afterSave, HasLen, 1)
	c.Assert(a.shouldEncrypt, Equals, true)
	c.Assert(a.filename, Equals, "test1.config.enc")
}

func (s *AccountsSuite) Test_ApplicationConfig_turnOffEncryption_doesntDoAnythingIfAlreadyNotEncrypted(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: false,
	}
	res := a.turnOffEncryption()
	c.Assert(res, Equals, false)
}

func (s *AccountsSuite) Test_ApplicationConfig_turnOffEncryption_turnsOffEncryption(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: true,
		filename:      "test1.config.enc",
	}
	res := a.turnOffEncryption()
	c.Assert(res, Equals, true)
	c.Assert(a.afterSave, HasLen, 1)
	c.Assert(a.shouldEncrypt, Equals, false)
	c.Assert(a.filename, Equals, "test1.config")
}

func (s *AccountsSuite) Test_ApplicationConfig_SetShouldSaveFileEncrypted_turningOn(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: true,
		filename:      "test1.config.enc",
	}
	res := a.SetShouldSaveFileEncrypted(true)
	c.Assert(res, Equals, false)
}

func (s *AccountsSuite) Test_ApplicationConfig_SetShouldSaveFileEncrypted_turningOff(c *C) {
	a := &ApplicationConfig{
		shouldEncrypt: false,
		filename:      "test1.config",
	}
	res := a.SetShouldSaveFileEncrypted(false)
	c.Assert(res, Equals, false)
}

func (s *AccountsSuite) Test_LoadOrCreate(c *C) {
	a, ok, e := LoadOrCreate("test111.conf", nil)
	c.Assert(a, Not(IsNil))
	c.Assert(a.filename, Equals, "test111.conf")
	c.Assert(ok, Equals, true)
	c.Assert(e, ErrorMatches, "Failed to parse config file")
}

func (s *AccountsSuite) Test_ApplicationConfig_removeOldFileOnNextSave_removesFileIfIsNotCurrentFilename(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	tmpFileName := tmpfile.Name()
	defer os.Remove(tmpFileName)

	a := &ApplicationConfig{filename: tmpFileName}

	a.removeOldFileOnNextSave()
	a.filename = "somethingelse"
	a.afterSave[0]()

	c.Assert(fileExists(tmpFileName), Equals, false)
}

func (s *AccountsSuite) Test_ApplicationConfig_removeOldFileOnNextSave_dontRemoveFileIfIsCurrentFilename(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	a := &ApplicationConfig{filename: tmpfile.Name()}

	a.removeOldFileOnNextSave()
	a.afterSave[0]()

	c.Assert(fileExists(tmpfile.Name()), Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_removeOldFileOnNextSave_doesntDoAnythingIfFileDoesntExist(c *C) {
	tmpFileName := generateTempFileName()

	a := &ApplicationConfig{filename: tmpFileName}

	a.removeOldFileOnNextSave()
	a.filename = "somethingelse"
	a.afterSave[0]()

	c.Assert(fileExists(tmpFileName), Equals, false)
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_failsWhenReadingNonExistingFile(c *C) {
	a := &ApplicationConfig{filename: "non-existing-file"}
	e := a.tryLoad(nil)
	c.Assert(e, Equals, errInvalidConfigFile)
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_loadsCorrectFile(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{
	"Accounts": [
		{
			"Account": "hello@foo.com",
			"Peers": null,
			"HideStatusUpdates": false,
			"OTRAutoTearDown": false,
			"OTRAutoAppendTag": false,
			"OTRAutoStartSession": false,
			"ConnectAutomatically": false
		}
	],
	"Bell": false,
	"ConnectAutomatically": false,
	"Display": {
		"MergeAccounts": false,
		"ShowOnlyOnline": false,
		"HideFeedbackBar": false,
		"ShowOnlyConfirmed": false,
		"SortByStatus": false
	},
	"AdvancedOptions": false,
	"UniqueConfigurationID": ""
}`))

	a := &ApplicationConfig{filename: tmpfile.Name()}

	previous := loadEntries
	defer func() {
		loadEntries = previous
	}()

	called1 := false
	called2 := false

	wg := sync.WaitGroup{}
	wg.Add(2)

	loadEntries = []func(*ApplicationConfig){
		func(*ApplicationConfig) {
			called1 = true
			wg.Done()
		},
		func(*ApplicationConfig) {
			called2 = true
			wg.Done()
		},
	}

	e := a.tryLoad(nil)
	c.Assert(e, IsNil)
	wg.Wait()
	c.Assert(called1, Equals, true)
	c.Assert(called2, Equals, true)
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_failsIfThereAreNoAccounts(c *C) {
	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(`{
	"Accounts": [
	],
	"Bell": false,
	"ConnectAutomatically": false,
	"Display": {
		"MergeAccounts": false,
		"ShowOnlyOnline": false,
		"HideFeedbackBar": false,
		"ShowOnlyConfirmed": false,
		"SortByStatus": false
	},
	"AdvancedOptions": false,
	"UniqueConfigurationID": ""
}`))

	a := &ApplicationConfig{filename: tmpfile.Name()}

	e := a.tryLoad(nil)
	c.Assert(e, Equals, errInvalidConfigFile)
}

const encryptedDataFileExample = `{
	"Params": {
		"Nonce": "dbd8f7642b05349123d59d1b",
		"Salt": "e18cb93a823465d2797539ebc5f3c0fd",
		"N": 262144,
		"R": 8,
		"P": 1
	},
	"Data": "24570e3858dbef3818ab86898a38a3855c1b89326e88b087697f09964173dc5e7a985aee75092852b6f6d56501f7e66e29cb19cc57138ab044853e2e4f1773d2a882dacf44e43800335e4d4973928cacd787c1c376db89b6cb46429f928048c2e3571a6b1184754ca40ef0dde9f745c18a640cee30fed2a1886f9377c9d3a60269f4b393fe14fbf43f76c742c55bdf6e7c3bc8dbcde71e66aed3564ff2ca8b5baa030350959f625654002f5bb13db318f2681d665f0bba2951dfd264e6e1493670a5b943931e841dd8424d5c179ffc3c63f5b2707a563842c225e7f18474d84a26a4414437061354a490dbe64ca4f2688fedf87631f3d3f2315ebfc7661e69c0b4f94f5c0ab7b900cfe582f974fa67264dd1c361266dc31a4a008d62a031748142c6e13758e013dacd7e21ed5da2c3316e72c79adbc490edab429a8e0d0bf4d3d2dd26d68b64bcde4c9229e402558537fdcb0c22bc9f4a444ebb3937de864c33187dd9edd8053542abeaf86259a5ceda401d588d355592192f4ed9ecdb56f96a214000bdb94fe914fe027ecd028ed77cc9fb74f6eb7555ca50af63a445c2ad4b8e33ef01f16d8242179dbd42c4e94e17e1ce5a3019a532c740c5263db0002feb85e4a649c43fadc360bb494d126a3981c65f567ad22603554b52ce567db8b9b54545bcbbc5429df267c034ddbf18c80d3f38b23fde5b65f2100d4fff188384d54a807d2e0250"
}`

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_encryptedFileWorks(c *C) {
	ks := FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	})

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(encryptedDataFileExample))

	a := &ApplicationConfig{
		filename: tmpfile.Name(),
	}

	e := a.tryLoad(ks)
	c.Assert(e, IsNil)
	c.Assert(a.shouldEncrypt, Equals, true)
	c.Assert(a.Accounts, HasLen, 1)
	c.Assert(a.Accounts[0].Account, Equals, "test1@example.com")
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_failsIfJSONDataIsInvalid(c *C) {
	data := `{`
	ks := FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, true
	})

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(data))

	a := &ApplicationConfig{
		filename: tmpfile.Name(),
	}

	e := a.tryLoad(ks)
	c.Assert(e, Equals, errInvalidConfigFile)
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_failsIfNoPasswordIsSupplied(c *C) {
	ks := FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKey, testMacKey, false
	})

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(encryptedDataFileExample))

	a := &ApplicationConfig{
		filename: tmpfile.Name(),
	}

	e := a.tryLoad(ks)
	c.Assert(e, Equals, errNoPasswordSupplied)
}

func (s *AccountsSuite) Test_ApplicationConfig_tryLoad_failsIfWrongPasswordIsSupplied(c *C) {
	ks := FunctionKeySupplier(func(params EncryptionParameters, _ bool) ([]byte, []byte, bool) {
		return testKeyWrong, testMacKey, true
	})

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(encryptedDataFileExample))

	a := &ApplicationConfig{
		filename: tmpfile.Name(),
	}

	e := a.tryLoad(ks)
	c.Assert(e, Equals, errDecryptionFailed)
}

func generateTempFileName() string {
	tmpfile, _ := ioutil.TempFile("", "")
	tmpfileName := tmpfile.Name()
	tmpfile.Close()
	os.Remove(tmpfileName)

	return tmpfileName
}
