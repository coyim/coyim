package importer

import "github.com/twstrike/coyim/config"

// TryImportAll will try to import from all known importers
func TryImportAll() map[string][]*config.Accounts {
	res := make(map[string][]*config.Accounts)

	res["Adium"] = (&adiumImporter{}).TryImport()
	res["Gajim"] = (&gajimImporter{}).TryImport()
	res["Pidgin"] = (&pidginImporter{}).TryImport()
	res["xmpp-client"] = (&xmppClientImporter{}).TryImport()

	return res
}
