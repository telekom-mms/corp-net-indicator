// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"de": &dictionary{index: deIndex, data: deData},
		"en": &dictionary{index: enIndex, data: enData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"About":                                 5,
	"Already connected to trusted network.": 7,
	"Cancel":                                24,
	"Certificate expires":                   32,
	"Connect":                               23,
	"Connect VPN":                           2,
	"Connect to VPN":                        3,
	"Connected":                             28,
	"Connected at":                          29,
	"Could not connect. Please Retry.":      8,
	"Could not disconnect. Please Retry.":   9,
	"Could not query certification expire date.":      12,
	"Could not query current Identity status.":        14,
	"Could not query current VPN status.":             10,
	"Could not query server list.":                    11,
	"Could not refresh identity login. Please Retry.": 13,
	"Device":                      31,
	"Disconnect VPN":              4,
	"Error: [%v]":                 15,
	"IP":                          30,
	"Identity Details":            17,
	"Kerberos ticket valid until": 20,
	"Last Refresh":                19,
	"Logged in":                   18,
	"Password":                    21,
	"Physical network":            27,
	"Program to show corporate network status.": 6,
	"ReLogin":     16,
	"Server":      22,
	"Show Status": 1,
	"Status":      0,
	"VPN Details": 26,
	"not trusted": 25,
	"trusted":     33,
}

var deIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000007, 0x00000017, 0x00000025,
	0x0000003b, 0x00000047, 0x0000004d, 0x00000086,
	0x000000c4, 0x000000f7, 0x00000128, 0x0000015c,
	0x00000188, 0x000001c3, 0x000001ee, 0x00000229,
	0x00000239, 0x00000246, 0x00000257, 0x00000262,
	0x00000273, 0x0000028f, 0x00000298, 0x0000029f,
	0x000002a9, 0x000002b3, 0x000002cb, 0x000002d7,
	0x000002e6, 0x000002f0, 0x000002fd, 0x00000300,
	// Entry 20 - 3F
	0x00000307, 0x0000031f, 0x00000331,
} // Size: 164 bytes

const deData string = "" + // Size: 817 bytes
	"\x02Status\x02Status anzeigen\x02VPN verbinden\x02Mit dem VPN verbinden" +
	"\x02VPN trennen\x02Über\x02Ein Programm zur Anzeige des Unternehmensnetz" +
	"werkstatus.\x02Verbindung zu einem vertrauenswürdigen Netzwerk hergestel" +
	"lt.\x02Verbindung fehlgeschlagen. Bitte erneut versuchen.\x02Trennung fe" +
	"hlgeschlagen. Bitte erneut versuchen.\x02Aktueller VPN-Status konnte nic" +
	"ht abgefragt werden.\x02Server-Liste konnte nicht abgefragt werden.\x02A" +
	"blaufdatum des Zertifikats konnte nicht abgefragt werden.\x02Identität k" +
	"onnte nicht angemeldet werden.\x02Aktueller Identitätsstatus konnte nich" +
	"t abgefragt werden.\x02Fehler: [%[1]v]\x02Neu anmelden\x02Identity Detai" +
	"ls\x02Angemeldet\x02Letzte Anmeldung\x02Kerberos Ticket gültig bis\x02Pa" +
	"sswort\x02Server\x02Verbinden\x02Abbrechen\x02nicht vertrauenswürdig\x02" +
	"VPN Details\x02Phys. Netzwerk\x02Verbunden\x02Verbunden am\x02IP\x02Gerä" +
	"t\x02Zertifikat läuft ab am\x02vertrauenswürdig"

var enIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000007, 0x00000013, 0x0000001f,
	0x0000002e, 0x0000003d, 0x00000043, 0x0000006d,
	0x00000093, 0x000000b4, 0x000000d8, 0x000000fc,
	0x00000119, 0x00000144, 0x00000174, 0x0000019d,
	0x000001ac, 0x000001b4, 0x000001c5, 0x000001cf,
	0x000001dc, 0x000001f8, 0x00000201, 0x00000208,
	0x00000210, 0x00000217, 0x00000223, 0x0000022f,
	0x00000240, 0x0000024a, 0x00000257, 0x0000025a,
	// Entry 20 - 3F
	0x00000261, 0x00000275, 0x0000027d,
} // Size: 164 bytes

const enData string = "" + // Size: 637 bytes
	"\x02Status\x02Show Status\x02Connect VPN\x02Connect to VPN\x02Disconnect" +
	" VPN\x02About\x02Program to show corporate network status.\x02Already co" +
	"nnected to trusted network.\x02Could not connect. Please Retry.\x02Could" +
	" not disconnect. Please Retry.\x02Could not query current VPN status." +
	"\x02Could not query server list.\x02Could not query certification expire" +
	" date.\x02Could not refresh identity login. Please Retry.\x02Could not q" +
	"uery current Identity status.\x02Error: [%[1]v]\x02ReLogin\x02Identity D" +
	"etails\x02Logged in\x02Last Refresh\x02Kerberos ticket valid until\x02Pa" +
	"ssword\x02Server\x02Connect\x02Cancel\x02not trusted\x02VPN Details\x02P" +
	"hysical network\x02Connected\x02Connected at\x02IP\x02Device\x02Certific" +
	"ate expires\x02trusted"

	// Total table size 1782 bytes (1KiB); checksum: 4908D755
