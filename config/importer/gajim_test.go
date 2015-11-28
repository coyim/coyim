package importer

import (
	"fmt"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/otr3"
	. "gopkg.in/check.v1"
)

type GajimSuite struct{}

var _ = Suite(&GajimSuite{})

func (s *GajimSuite) Test_GajimImporter_canImportFingerprintsFromFile(c *C) {
	importer := gajimImporter{}

	nm, res, ok := importer.importFingerprintsFrom(testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de.fpr"))

	c.Assert(ok, Equals, true)
	c.Assert(nm, Equals, "aba.baba@jabber.ccc.de")
	c.Assert(len(res), Equals, 6)

	c.Assert(res[0].UserID, Equals, "abcde@thoughtworks.com")
	c.Assert(res[0].FingerprintHex, Equals, "57d8ea36c76d5d800fe790c56dc33feb254e899b")
	c.Assert(res[0].Untrusted, Equals, true)

	c.Assert(res[1].UserID, Equals, "coyim@thoughtworks.com")
	c.Assert(res[1].FingerprintHex, Equals, "c8123327e389e3d036ba91cf92d722f515057b61")
	c.Assert(res[1].Untrusted, Equals, true)

	c.Assert(res[2].UserID, Equals, "someone@where.com")
	c.Assert(res[2].FingerprintHex, Equals, "a334e9d582da18f15028f7f7412bc8d15d0a1558")
	c.Assert(res[2].Untrusted, Equals, true)

	c.Assert(res[3].UserID, Equals, "bla@rose.com")
	c.Assert(res[3].FingerprintHex, Equals, "7c6c74ddb307c95fa30c3ecab25ee64a54124447")
	c.Assert(res[3].Untrusted, Equals, false)

	c.Assert(res[4].UserID, Equals, "not@coyim.com")
	c.Assert(res[4].FingerprintHex, Equals, "4157eea3bb3cf86cc0379e4c270e89b976bc34da")
	c.Assert(res[4].Untrusted, Equals, true)

	c.Assert(res[5].UserID, Equals, "not@coyim.com")
	c.Assert(res[5].FingerprintHex, Equals, "edd6274423cd2fb6993da928d923075be2d0d52a")
	c.Assert(res[5].Untrusted, Equals, true)
}

func (s *GajimSuite) Test_GajimImporter_canImportPrivateKey(c *C) {
	importer := gajimImporter{}

	nm, res, ok := importer.importKeyFrom(testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de.key3"))

	c.Assert(ok, Equals, true)
	c.Assert(nm, Equals, "aba.baba@jabber.ccc.de")

	pk := &otr3.DSAPrivateKey{}
	pk.Parse(res)

	c.Assert(fmt.Sprintf("%X", otr3.NewConversationWithVersion(3).DefaultFingerprintFor(pk)), Equals, "0AB95107E9457E494F7FA68E8AAD1B86EE96935E")
}

func (s *GajimSuite) Test_GajimImporter_canImportPluginSettings(c *C) {
	importer := gajimImporter{}
	res, res2, ok := importer.importOTRSettings(testResourceFilename("gajim_test_data/gotr"))

	c.Assert(ok, Equals, true)

	c.Assert(len(res), Equals, 27)
	c.Assert(len(res2), Equals, 59)

	c.Assert(res["AFinalLongNameMaybea"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["AnameAnameAnameAnm"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["AnotherName"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["Local"], Equals, gajimOTRSettings{allowV2: true, errorStartAke: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["aaaaa"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["aba.baba@jabber.ccc.de"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["abc@coy.im"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["abc@coy.org"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["b@coy.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["b@coy.im"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["coy.test@riseup.net"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["gmail.com"], Equals, gajimOTRSettings{allowV2: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["jabber.ccc.de"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["jabber.org"], Equals, gajimOTRSettings{allowV2: true, errorStartAke: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["jabber.org.by"], Equals, gajimOTRSettings{allowV2: true, errorStartAke: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["ola@coy.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["ola@olabini.se"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["one.two@jabber.ccc.de"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: false})
	c.Assert(res["q@coy.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["q@coy.im"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["riseup.net"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["thoughtworks.com"], Equals, gajimOTRSettings{allowV2: true, errorStartAke: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["ventwo22@jabber.org"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["xxxxyyyy@gmail.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["yyyyy@thoughtworks.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: false, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["z@coy.com"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res["z@coy.im"], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})

	c.Assert(res2[gajimAccountAndPeer{"AFinalLongNameMaybea", "ccaaccaac@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"AnameAnameAnameAnm", "anothernamesomewhe@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"AnameAnameAnameAnm", "onetwoth@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "a234@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "a23accccccccc@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "a56@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "aaaaaaqqqq@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "aapaapa@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "ab123@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "ab@cddefg.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "ababab@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "abcd7575@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "abcd@abcdabcd.org"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "abcdefa@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "abcdefb@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "abcdefg@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "blueblue@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "ccccc-ddddd@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "foofoofooa@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "ggggg@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "llll@jabber.xxxxxxx.yyy.se"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "lllll@nnnnnnnn.im"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "nlnln@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "pmpmm@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "rraarraa111@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "rrrrrrrr@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "spspspsps.marttt@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "sysysysysysysy@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "vsvsvs@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "xxxcc@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"aba.baba@jabber.ccc.de", "zzzzzzz.xxxxx@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "aabbccdde@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "abcabca@dukgo.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "abcdefgabcdefabce@dukgo.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "cthrr@coy.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "fourfive.sixe@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "onetwot@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"abc@coy.org", "onetwoth@jabber.calyxinstitute.org"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "ab@mabcab.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "abcd@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "abcdef@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "bbbbb@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "hhbbhhh@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "onefoursixseven@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "sbaasdwwww@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "whynotthisonea@jabber.ccc.de"}], Equals, gajimOTRSettings{allowV2: true, requireEncryption: true, sendTag: true, whitespaceStartAke: true})
	c.Assert(res2[gajimAccountAndPeer{"jabber.ccc.de", "xxaayyss@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"ola@olabini.se", "z12312312@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"ola@olabini.se", "z1231231@jabber.org"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"ola@olabini.se", "z123123@riseup.net"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"one.two@jabber.ccc.de", "one.two.thr@jabber.ccc.de"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"thoughtworks.com", "x123123@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"thoughtworks.com", "y1231231@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"ventwo22@jabber.org", "abcabc@jabber.org"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"yyyyy@thoughtworks.com", "eightni@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"yyyyy@thoughtworks.com", "neelev@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"yyyyy@thoughtworks.com", "ofourfiv@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"yyyyy@thoughtworks.com", "onetwoth@thoughtworks.com"}], Equals, gajimOTRSettings{})
	c.Assert(res2[gajimAccountAndPeer{"yyyyy@thoughtworks.com", "sixseve@thoughtworks.com"}], Equals, gajimOTRSettings{})
}

func (s *GajimSuite) Test_GajimImporter_canImportAccountInformation(c *C) {
	importer := gajimImporter{}
	res, ok := importer.importAccounts(testResourceFilename("gajim_test_data/config"))

	c.Assert(ok, Equals, true)
	c.Assert(len(res), Equals, 17)
	c.Assert(res["Local"], Equals, gajimAccountInfo{accountNickName: "Local", password: "zeroconf", hostname: "coy.com", name: "b", port: "5298"})
	c.Assert(res["aaaaa@riseup.net"], Equals, gajimAccountInfo{accountNickName: "aaaaa@riseup.net", password: "sdddddddddddddddddddddddddddddddddddddddddd", sslFingerprint: "1CFD0A83738A497B0399FB74E1E978A459F8546F", hostname: "riseup.net", name: "aaaaa", port: "5222"})
	c.Assert(res["aba.baba@jabber.ccc.de"], Equals, gajimAccountInfo{accountNickName: "aba.baba@jabber.ccc.de", password: "foo bar barium", sslFingerprint: "4E09F9D9F224174684768D467A84B139B86A021F", hostname: "jabber.ccc.de", server: "", name: "aba.baba", proxy: "", port: "5222"})
	c.Assert(res["abc@coy.org"], Equals, gajimAccountInfo{accountNickName: "abc@coy.org", password: "foo=bar", sslFingerprint: "DA1F3DD285CDE8B1B5BA254B6610F8A7F4BA0B0F", hostname: "coy.org", server: "", name: "abc", proxy: "", port: "5222"})
	c.Assert(res["b@coy.com"], Equals, gajimAccountInfo{accountNickName: "b@coy.com", password: "sdsafjc zxcxvxcvxcvxdcvfbf dfgdf", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.com", server: "", name: "b", proxy: "", port: "5222"})
	c.Assert(res["b@coy.im"], Equals, gajimAccountInfo{accountNickName: "b@coy.im", password: "7yfghfdghdfghfghfgfhdgh", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.im", server: "", name: "b", proxy: "", port: "5222"})
	c.Assert(res["coy.test@riseup.net"], Equals, gajimAccountInfo{accountNickName: "coy.test@riseup.net", password: "abcdef", sslFingerprint: "1CFD0A83738A497B0399FB74E1E978A459F8546F", hostname: "riseup.net", server: "", name: "coy.test", proxy: "", port: "5222"})
	c.Assert(res["ola@coy.com"], Equals, gajimAccountInfo{accountNickName: "ola@coy.com", password: "76ydfbgfdhdswewa", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.com", server: "", name: "ola", proxy: "", port: "5222"})
	c.Assert(res["ola@coy.im"], Equals, gajimAccountInfo{accountNickName: "ola@coy.im", password: "aabbbccdeedef", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.im", server: "", name: "ola", proxy: "", port: "5222"})
	c.Assert(res["ola@olabini.se"], Equals, gajimAccountInfo{accountNickName: "ola@olabini.se", password: "11aa11aa11aa", sslFingerprint: "DA1F3DD285CDE8B1B5BA254B6610F8A7F4BA0B0F", hostname: "olabini.se", server: "", name: "ola", proxy: "socks5://tor3:tor3@localhost:9050", port: "5222"})
	c.Assert(res["q@coy.com"], Equals, gajimAccountInfo{accountNickName: "q@coy.com", password: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaBB", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.com", server: "", name: "q", proxy: "", port: "5222"})
	c.Assert(res["q@coy.im"], Equals, gajimAccountInfo{accountNickName: "q@coy.im", password: "6575rgtgfsertgrty34t", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.im", server: "", name: "q", proxy: "", port: "5222"})
	c.Assert(res["ventwo22@jabber.org"], Equals, gajimAccountInfo{accountNickName: "ventwo22@jabber.org", password: "GDFGDFgsdfgsdfgdDFGSDFGSDFGS", sslFingerprint: "321AF3670EF2945ABEA04D42ECD2E8EC99CDEFCD", hostname: "jabber.org", server: "", name: "ventwo22", proxy: "", port: "5222"})
	c.Assert(res["xxxxyyyy@gmail.com"], Equals, gajimAccountInfo{accountNickName: "xxxxyyyy@gmail.com", password: "aaaaaaaaaaaaaaaaaaaaaaaa", sslFingerprint: "FB5E82C43B9D589CF85BE41EE459FBCE27B1A6E8", hostname: "gmail.com", server: "", name: "xxxxyyyy", proxy: "", port: "5222"})
	c.Assert(res["yyyyy@thoughtworks.com"], Equals, gajimAccountInfo{accountNickName: "yyyyy@thoughtworks.com", password: "aa11ytjtj(*(&*&*", sslFingerprint: "B9BEAFB5C198DEFC7D81AC873978315715475925", hostname: "thoughtworks.com", server: "talk.google.com", name: "yyyyy", proxy: "socks5://tor1:tor1@localhost:9050", port: "5222"})
	c.Assert(res["z@coy.com"], Equals, gajimAccountInfo{accountNickName: "z@coy.com", password: "zzzzzzzzzzzzzzzzzzzzz", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.com", server: "", name: "z", proxy: "", port: "5222"})
	c.Assert(res["z@coy.im"], Equals, gajimAccountInfo{accountNickName: "z@coy.im", password: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", sslFingerprint: "FB30D19A4C6D2ECA4530DECE9EA7A1106B824E60", hostname: "coy.im", server: "", name: "z", proxy: "", port: "5222"})
}

func (s *GajimSuite) Test_GajimImporter_canDoAFullImport(c *C) {
	importer := gajimImporter{}
	res, ok := importer.importAllFrom(
		testResourceFilename("gajim_test_data/config"),
		testResourceFilename("gajim_test_data/gotr"),
		[]string{testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de.key3")},
		[]string{testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de.fpr")},
	)

	c.Assert(ok, Equals, true)
	c.Assert(res, Not(IsNil))
	c.Assert(len(res.Accounts), Equals, 16)

	c.Assert(*res.Accounts[0], DeepEquals, config.Account{
		Account:           "aaaaa@riseup.net",
		Server:            "riseup.net",
		Password:          "sdddddddddddddddddddddddddddddddddddddddddd",
		Port:              5222,
		OTRAutoTearDown:   true,
		AlwaysEncryptWith: []string{},
		DontEncryptWith:   []string{},
	})

	c.Assert(*res.Accounts[1], DeepEquals, config.Account{
		Account:     "aba.baba@jabber.ccc.de",
		Server:      "jabber.ccc.de",
		Password:    "foo bar barium",
		Port:        5222,
		PrivateKeys: [][]byte{[]uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0xe6, 0x1d, 0xff, 0x9e, 0xae, 0xdd, 0x1c, 0x5e, 0xc8, 0xf5, 0x3f, 0x4c, 0x1, 0x32, 0x40, 0x2a, 0xe0, 0xfd, 0x60, 0x2e, 0x52, 0x6c, 0x31, 0xde, 0xfa, 0xc, 0x38, 0x6e, 0x91, 0x56, 0x3d, 0x81, 0x54, 0x26, 0xca, 0x6a, 0x76, 0x91, 0x9c, 0x64, 0xca, 0x4a, 0x26, 0x21, 0x70, 0xb3, 0xec, 0xc4, 0x59, 0x60, 0xd8, 0x3d, 0xf5, 0xe8, 0x61, 0x44, 0xf3, 0x1d, 0xa4, 0xa6, 0x2f, 0x9f, 0x51, 0x8d, 0x57, 0x8, 0x93, 0x3d, 0x5e, 0x25, 0x2a, 0x16, 0x15, 0x18, 0xb3, 0x2, 0x6, 0x31, 0x8e, 0x5d, 0x5f, 0x2d, 0x3f, 0xd8, 0x25, 0xe0, 0xd1, 0x96, 0xae, 0x3c, 0x11, 0xd1, 0x6, 0xc6, 0xd7, 0x9c, 0x63, 0xde, 0xc2, 0x9c, 0x2, 0x87, 0x8, 0x82, 0x67, 0xcf, 0x8f, 0x9d, 0xe7, 0x61, 0xc5, 0xed, 0x44, 0x78, 0xc8, 0xf5, 0xdc, 0xd6, 0xb3, 0x8e, 0xf8, 0x84, 0xf0, 0x87, 0x53, 0xcf, 0x72, 0x21, 0x0, 0x0, 0x0, 0x14, 0x95, 0x6f, 0xbc, 0xe8, 0x7a, 0xb9, 0x73, 0x42, 0x55, 0x17, 0x1d, 0x4b, 0xe1, 0x94, 0x9f, 0x98, 0x43, 0x86, 0x3c, 0xe5, 0x0, 0x0, 0x0, 0x80, 0xc, 0xca, 0x27, 0xf8, 0xf5, 0x99, 0xec, 0xc, 0xd9, 0xae, 0x6e, 0x89, 0x9f, 0x9, 0x3b, 0xfa, 0x55, 0x3, 0x25, 0x54, 0x5e, 0x54, 0xa1, 0xa2, 0x35, 0x84, 0x88, 0xfe, 0x7c, 0x14, 0x49, 0x25, 0xe2, 0x1f, 0xd5, 0xd7, 0x9b, 0x88, 0x98, 0x5c, 0x53, 0x57, 0x1c, 0xba, 0xd, 0xec, 0x9f, 0x5e, 0x5e, 0x5f, 0x3b, 0x20, 0xf5, 0x6c, 0xd8, 0x3b, 0xf6, 0x31, 0xdf, 0x8f, 0xe6, 0x92, 0x8e, 0x2e, 0xf8, 0xec, 0xf5, 0xd6, 0x1f, 0x42, 0xe3, 0x59, 0x84, 0x4d, 0x3c, 0xa6, 0xe7, 0x95, 0x58, 0x2e, 0x4a, 0xc3, 0xf9, 0xf7, 0x1a, 0x94, 0x15, 0x8, 0xde, 0x80, 0xdc, 0x95, 0x1c, 0xd3, 0xd7, 0xa6, 0x2e, 0x17, 0x17, 0x48, 0xdf, 0xd2, 0xfa, 0x91, 0x74, 0x8e, 0x81, 0xbb, 0x2b, 0x5b, 0xae, 0xb1, 0x91, 0x8e, 0x4c, 0x54, 0xb6, 0xf0, 0x40, 0x20, 0x29, 0xb3, 0x71, 0xde, 0x10, 0xed, 0x4e, 0xaa, 0x0, 0x0, 0x0, 0x80, 0xa3, 0x70, 0x98, 0x8e, 0x1b, 0x46, 0xf6, 0xda, 0xde, 0xf4, 0xbb, 0xa6, 0x7b, 0x3c, 0x7e, 0x1b, 0x2e, 0x4a, 0x5, 0x6d, 0x9d, 0x1a, 0x4, 0x5f, 0x56, 0x5, 0x12, 0xf2, 0xb5, 0xf8, 0xb5, 0xce, 0x8c, 0x9, 0x96, 0x1d, 0xb5, 0x54, 0x93, 0x6e, 0x40, 0x9, 0x6d, 0x1e, 0xa9, 0x2c, 0x55, 0xf2, 0xb5, 0x9d, 0x43, 0x9b, 0x65, 0xda, 0xd7, 0xe8, 0x95, 0x45, 0x79, 0x6, 0xba, 0xfc, 0x31, 0x91, 0xe4, 0x25, 0x37, 0x8c, 0x2e, 0xe8, 0xe9, 0xdf, 0x4, 0xfb, 0xdb, 0x19, 0x9b, 0xae, 0x9c, 0xb2, 0x8, 0xd8, 0x6d, 0xf9, 0xbe, 0x91, 0xce, 0xb2, 0x2b, 0x7e, 0x84, 0xd0, 0x4b, 0x32, 0xe2, 0xd0, 0xc7, 0xdb, 0xb2, 0x49, 0xf3, 0x6, 0xe9, 0x88, 0x3d, 0x1d, 0x32, 0x42, 0xd, 0x10, 0xfa, 0x38, 0x9f, 0x8a, 0xda, 0x3b, 0xf7, 0x7c, 0xc8, 0x1d, 0x77, 0xbe, 0x88, 0x25, 0xc1, 0xe6, 0x76, 0x78, 0x0, 0x0, 0x0, 0x14, 0x8a, 0x9a, 0x8f, 0xb3, 0x44, 0x17, 0xd1, 0x8, 0x3d, 0xa8, 0xf4, 0xc7, 0xe7, 0x6d, 0x51, 0x88, 0xe0, 0x96, 0x4b, 0xc2}},
		KnownFingerprints: []config.KnownFingerprint{
			config.KnownFingerprint{UserID: "abcde@thoughtworks.com", FingerprintHex: "57d8ea36c76d5d800fe790c56dc33feb254e899b", Untrusted: true},
			config.KnownFingerprint{UserID: "bla@rose.com", FingerprintHex: "7c6c74ddb307c95fa30c3ecab25ee64a54124447"},
			config.KnownFingerprint{UserID: "coyim@thoughtworks.com", FingerprintHex: "c8123327e389e3d036ba91cf92d722f515057b61", Untrusted: true},
			config.KnownFingerprint{UserID: "not@coyim.com", FingerprintHex: "4157eea3bb3cf86cc0379e4c270e89b976bc34da", Untrusted: true},
			config.KnownFingerprint{UserID: "not@coyim.com", FingerprintHex: "edd6274423cd2fb6993da928d923075be2d0d52a", Untrusted: true},
			config.KnownFingerprint{UserID: "someone@where.com", FingerprintHex: "a334e9d582da18f15028f7f7412bc8d15d0a1558", Untrusted: true},
		},
		RequireTor:          false,
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
	})

	c.Assert(*res.Accounts[2], DeepEquals, config.Account{
		Account:             "abc@coy.org",
		Server:              "coy.org",
		Password:            "foo=bar",
		Port:                5222,
		AlwaysEncrypt:       true,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
	})

	c.Assert(*res.Accounts[3], DeepEquals, config.Account{
		Account:             "b@coy.com",
		Server:              "coy.com",
		Password:            "sdsafjc zxcxvxcvxcvxdcvfbf dfgdf",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[4], DeepEquals, config.Account{
		Account:             "b@coy.im",
		Server:              "coy.im",
		Password:            "7yfghfdghdfghfghfgfhdgh",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[5], DeepEquals, config.Account{
		Account:             "coy.test@riseup.net",
		Server:              "riseup.net",
		Password:            "abcdef",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[6], DeepEquals, config.Account{
		Account:             "ola@coy.com",
		Server:              "coy.com",
		Password:            "76ydfbgfdhdswewa",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[7], DeepEquals, config.Account{
		Account:           "ola@coy.im",
		Server:            "coy.im",
		Password:          "aabbbccdeedef",
		Port:              5222,
		AlwaysEncryptWith: []string{},
		DontEncryptWith:   []string{},
		OTRAutoTearDown:   true,
	})

	c.Assert(*res.Accounts[8], DeepEquals, config.Account{
		Account:             "ola@olabini.se",
		Server:              "olabini.se",
		Password:            "11aa11aa11aa",
		Port:                5222,
		Proxies:             []string{"socks5://tor3:tor3@localhost:9050"},
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		RequireTor:          true,
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[9], DeepEquals, config.Account{
		Account:             "q@coy.com",
		Server:              "coy.com",
		Password:            "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaBB",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[10], DeepEquals, config.Account{
		Account:             "q@coy.im",
		Server:              "coy.im",
		Password:            "6575rgtgfsertgrty34t",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[11], DeepEquals, config.Account{
		Account:             "ventwo22@jabber.org",
		Server:              "jabber.org",
		Password:            "GDFGDFgsdfgsdfgdDFGSDFGSDFGS",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[12], DeepEquals, config.Account{
		Account:             "xxxxyyyy@gmail.com",
		Server:              "gmail.com",
		Password:            "aaaaaaaaaaaaaaaaaaaaaaaa",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
	})

	c.Assert(*res.Accounts[13], DeepEquals, config.Account{
		Account:             "yyyyy@thoughtworks.com",
		Server:              "talk.google.com",
		Password:            "aa11ytjtj(*(&*&*",
		Proxies:             []string{"socks5://tor1:tor1@localhost:9050"},
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{"eightni@thoughtworks.com", "neelev@thoughtworks.com", "ofourfiv@thoughtworks.com", "onetwoth@thoughtworks.com", "sixseve@thoughtworks.com"},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		RequireTor:          true,
	})

	c.Assert(*res.Accounts[14], DeepEquals, config.Account{
		Account:             "z@coy.com",
		Server:              "coy.com",
		Password:            "zzzzzzzzzzzzzzzzzzzzz",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})

	c.Assert(*res.Accounts[15], DeepEquals, config.Account{
		Account:             "z@coy.im",
		Server:              "coy.im",
		Password:            "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Port:                5222,
		AlwaysEncryptWith:   []string{},
		DontEncryptWith:     []string{},
		OTRAutoTearDown:     true,
		OTRAutoAppendTag:    true,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
	})
}

func (s *GajimSuite) Test_GajimImporter_canFailAFullImport(c *C) {
	importer := gajimImporter{}
	_, ok := importer.importAllFrom(
		testResourceFilename("gajim_test_data/config2"),
		testResourceFilename("gajim_test_data/gotr2"),
		[]string{testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de2.key3")},
		[]string{testResourceFilename("gajim_test_data/aba.baba@jabber.ccc.de2.fpr")},
	)

	c.Assert(ok, Equals, false)
}
