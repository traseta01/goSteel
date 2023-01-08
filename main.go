package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sf1/go-card/smartcard"
)

// type to hold e-card read data
type LicnaKarta struct {
	prezime         string
	ime             string
	imeRoditelja    string
	datumRodjenja   string
	mestoRodjenja   string
	opstinaRodjenja string
	drzavaRodjenja  string
	JMBG            string
	pol             string
	dokumentIzdaje  string
	brojDokumenta   string
	datumIzdavanja  string
	vaziDo          string
	prebivaliste    string
	adresaMesto     string
	adresaOpstina   string
	adresaUlica     string
	adresaBroj      string
	adresaStan      string
	slika           []byte
}

// variable to hold e-card read data
var lk *LicnaKarta

var reader *smartcard.Reader

func main() {

	// establish reader context
	context, err := smartcard.EstablishContext()
	if err != nil {
		fmt.Println("Error establishing context!")
		return
	}
	// get readers list
	readers_list, err := context.ListReadersWithCard()
	if err != nil {
		fmt.Println("Error getting card reader list!")
		return
	}

	// handle list of readers & get card
	if len(readers_list) == 0 {
		fmt.Println("\nPlease insert smart card\n")
		return
	}

	// work with one reader for now
	if len(readers_list) == 1 {
		reader = readers_list[0]
	}

	// connect to reader
	card, err := reader.Connect()
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%T\n------------------------------------\n", card)

	defer card.Disconnect()

	defer context.Release()

	// create empty LicnaKarta
	lk = new(LicnaKarta)

	// create new fyne app
	myApp := app.New()
	windw := myApp.NewWindow("GoSteel")

	text1 := canvas.NewText("***   ČITAČ ELEKTRONSKE LIČNE KARTE   ***", color.White)

	contentHello := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text1, layout.NewSpacer())

	labelRuler := widget.NewLabel("_________________________________________________________________________________________\nPodaci o građaninu\n _________________________________________________________________________________________")
	// labelRuler01 := widget.NewLabel("________________________________________________________________________________________")
	labelRuler02 := widget.NewLabel("________________________________________________________________________________________")
	labelPrezime := widget.NewLabel("Prezime:  ")
	labelIme := widget.NewLabel("Ime:    ")
	labelImeRoditelja := widget.NewLabel("Ime jednog roditelja:    ")
	labelPol := widget.NewLabel("Pol:    ")
	labeljmbg := widget.NewLabel("JMBG:    ")
	labelDatumRodjenja := widget.NewLabel("Datum rodjenja:    ")
	labelMestoRodjenjaOpstinaDrzava := widget.NewLabel("Mesto rodjenja, opstina i drzava:    ")

	labelPrebivalisteAdresa := widget.NewLabel("Prebivalište i adresa stana:    ")
	contentImePrezime := container.New(
		layout.NewVBoxLayout(),
		labelRuler,
		labelIme,
		labelPrezime,
		labelImeRoditelja,
		labelDatumRodjenja,
		labelMestoRodjenjaOpstinaDrzava,
		labelPrebivalisteAdresa,
		labeljmbg,
		labelPol,
		// labelRuler01,
		labelRuler02,
	)
	// container to display e-card image
	contentImage := container.New(layout.NewHBoxLayout())

	imagee := canvas.NewImageFromResource(theme.FyneLogo())
	imagee.FillMode = canvas.ImageFillOriginal

	buttonScan := widget.NewButton("Scan", func() {

		// read smartcard data
		lk.getImage(card)
		lk.readDataOne(card)
		lk.readDataTwo(card)
		lk.readDataThree(card)

		// prepare image
		img, _, err := image.Decode(bytes.NewReader(lk.slika))
		if err != nil {
			log.Fatalln(err)
		}

		imagee = canvas.NewImageFromImage(img)
		imagee.FillMode = canvas.ImageFillOriginal

		contentImage.AddObject(imagee)
		contentImage.Refresh()
	})

	buttonGenPdf := widget.NewButton("Create PDF", func() {

	})

	contetButtons := container.New(layout.NewHBoxLayout(), buttonScan, layout.NewSpacer(), buttonGenPdf)

	// set fyne window contents and attributes
	windw.SetContent(container.New(layout.NewVBoxLayout(), contetButtons, contentHello, contentImage, contentImePrezime))
	windw.Resize(fyne.NewSize(800, 900))
	windw.SetFixedSize(true)
	windw.CenterOnScreen()

	windw.ShowAndRun()
}

// Send APDU's to smart card
func sendCommand(apducommand []byte, card *smartcard.Card) []byte {

	res, err := card.TransmitAPDU(apducommand)
	if err != nil {
		panic(err)
	}
	// return apdu command response, minus last two bytes 90 00
	return res[:len(res)-2]
}

// get image from ID card
func (lkarta *LicnaKarta) getImage(card *smartcard.Card) {
	image := []byte{}
	apdu := []byte{}

	// initialize application for type2 cards, no impact on type1
	apdu = []byte{}
	apdu = append(apdu, 0x00, 0xA4, 0x04, 0x00, 0x0B, 0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x49, 0x44, 0x01, 0x00)
	sendCommand(smartcard.CommandAPDU(apdu), card)

	// select image file
	apdu = []byte{}
	apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x02, 0x0f, 0x06)
	// apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x00, 0x0f, 0x06)
	sendCommand(smartcard.CommandAPDU(apdu), card)

	// go smartcard API won't let us send 00 0B 00 00 00 so we'll use
	// 00 B0 00 00 FF and make sure to read an extra trailing byte
	// first read gives us number of sections to read and number of bytes
	// we need to read in the last section
	apdu = []byte{}
	apdu = append(apdu, 0x00, 0xB0, 0x00, 0x00, 0xff)
	pom := sendCommand(smartcard.CommandAPDU(apdu), card)

	fmt.Printf("%x\n\n\n\n", pom)

	imglen := pom[7]
	imglast := pom[6]

	// image = append(image, pom...)
	// image = append(image, 0xFF, 0xF8)
	image = append(image, pom...)
	apdu[3] = apdu[3] + 1
	pom = sendCommand(smartcard.CommandAPDU(apdu), card)
	image = append(image, pom[len(pom)-1])

	// fmt.Printf("duzina %d, \t broj krajnjih bajtova %d\n", imglen, imglast)

	// read the rest of the image
	var i byte
	for i = 1; i <= imglen; i++ {

		if i == imglen {
			apdu = []byte{}
			// apdu = append(apdu, 0x00, 0xB0, i, 0x00, imglast+6)
			apdu = append(apdu, 0x00, 0xB0, i, 0x00, imglast+8)
			pomocna01 := sendCommand(smartcard.CommandAPDU(apdu), card)
			image = append(image, pomocna01...)

			break
		}

		apdu = []byte{}
		apdu = append(apdu, 0x00, 0xB0, i, 0x00, 0xff)

		// read ith part of the image
		pomocna01 := sendCommand(smartcard.CommandAPDU(apdu), card)
		image = append(image, pomocna01...)
		// read an extra trailing byte
		apdu[3] = apdu[3] + 1
		pomocna01 = sendCommand(smartcard.CommandAPDU(apdu), card)
		image = append(image, pomocna01[len(pomocna01)-1])

	}
	// new card types
	// ATR: 3B FF 94 00 00 81 31 80 43 80 31 80 65 B0 85 02 01 F3 12 0F FF 82 90 00 79
	// ATR: 3B B9 18 00 81 31 FE 9E 80 73 FF 61 40 83 00 00 00 DF

	// old card types
	// file starts at 10 bytes offset
	// lk.slika = image[10:]

	// Type2 Card
	// file starts at 10 bytes offset
	lkarta.slika = image[8:]

	fmt.Printf("%x", lk.slika)
}

// read first set of data from ID card
func (lkarta *LicnaKarta) readDataOne(card *smartcard.Card) {

	apdu := []byte{}

	// apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x02, 0x0f, 0x02)
	// Type2
	apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x02, 0x0f, 0x02, 0x00)
	offset := sendCommand(smartcard.CommandAPDU(apdu), card)[3]
	fmt.Printf("\n\n%x\n\n", offset)

	apdu = []byte{}
	// apdu = append(apdu, 0x00, 0xB0, 0x00, 0x08, 0x00)
	apdu = append(apdu, 0x00, 0xB0, 0x00, 0x00, offset)
	pom := sendCommand(smartcard.CommandAPDU(apdu), card)

	fmt.Println("%x\n------------------------\n\n", pom)

	// drop first 4 bytes
	pom = pom[4:]
	fmt.Print("%x\n", pom)

	// read bytes of data
	lkarta.brojDokumenta = string(pom[4 : pom[2]+4])
	fmt.Printf("\n\nBR DOKUMENTA (Prvi podatak) : %s", lkarta.brojDokumenta)
	// pom = pom[2+pom[0]:]
	pom = pom[4+pom[2]:]
	fmt.Printf("\n\nBR (Drugi podatak) : %s", string(pom[4:pom[2]+4]))
	pom = pom[4+pom[2]:]
	fmt.Printf("\n\nBR (Treci podatak) : %s", string(pom[4:pom[2]+4]))
	pom = pom[4+pom[2]:]
	fmt.Printf("\n\nBR (Cetvrti podatak) : %s", string(pom[4:pom[2]+4]))
	pom = pom[4+pom[2]:]
	lkarta.datumIzdavanja = pretyDate(string(pom[4 : pom[2]+4]))
	fmt.Printf("\n\nBR (Peti podatak) Datum Izdavanja: %s", lkarta.datumIzdavanja)
	pom = pom[4+pom[2]:]
	lkarta.vaziDo = pretyDate(string(pom[4 : pom[2]+4]))
	fmt.Printf("\n\nBR (Peti podatak) Datum Izdavanja: %s", lkarta.vaziDo)
	pom = pom[4+pom[2]:]
	lkarta.dokumentIzdaje = pretyDate(string(pom[4 : pom[2]+4]))
	fmt.Printf("\n\nIZADATA OD : %s", lkarta.dokumentIzdaje)
	// fmt.Printf("\n\nZEMLJA : %s", string(pom[4:pom[2]+4]))
	// pom = pom[4+pom[2]:]
}

// read second set of data from ID card
func (lkarta *LicnaKarta) readDataTwo(card *smartcard.Card) {
	apdu := []byte{}

	apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x02, 0x0f, 0x03, 0x00)
	offset := sendCommand(smartcard.CommandAPDU(apdu), card)[3]
	fmt.Printf("\n\n%x\n\n", offset)

	apdu = []byte{}
	// apdu = append(apdu, 0x00, 0xB0, 0x00, 0x08, 0x00)
	apdu = append(apdu, 0x00, 0xB0, 0x00, 0x00, offset)
	pom := sendCommand(smartcard.CommandAPDU(apdu), card)
	pom = pom[4:]
	fmt.Printf("\n\n (Prvi podatak) : %s", string(pom[4:pom[2]+4]))

	lkarta.JMBG = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nPREZIME : %s", string(pom[4:pom[2]+4]))
	lkarta.prezime = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nIME : %s", string(pom[4:pom[2]+4]))
	lkarta.ime = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nIME OCA : %s", string(pom[4:pom[2]+4]))
	lkarta.imeRoditelja = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nPOL : %s", string(pom[4:pom[2]+4]))
	lkarta.pol = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nMESTO RODJENJAr:   %s", string(pom[4:pom[2]+4]))
	lkarta.mestoRodjenja = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nOpstina rodjenja:  %s", string(pom[4:pom[2]+4]))
	lkarta.opstinaRodjenja = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nDRZAVA : %s", string(pom[4:pom[2]+4]))
	lkarta.drzavaRodjenja = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nDATUM RODJENJA : %s", string(pom[4:pom[2]+4]))
	lkarta.datumRodjenja = pretyDate(string(pom[4 : pom[2]+4]))
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nDRZAVA SKRACENO : %s", string(pom[4:pom[2]+4]))
}

// read third set of information
func (lk *LicnaKarta) readDataThree(card *smartcard.Card) {
	apdu := []byte{}

	apdu = append(apdu, 0x00, 0xA4, 0x08, 0x00, 0x02, 0x0f, 0x04, 0x00)
	offset := sendCommand(smartcard.CommandAPDU(apdu), card)[3]
	fmt.Printf("\n\n%x\n\n", offset)

	apdu = []byte{}
	// apdu = append(apdu, 0x00, 0xB0, 0x00, 0x08, 0x00)
	apdu = append(apdu, 0x00, 0xB0, 0x00, 0x00, offset)
	pom := sendCommand(smartcard.CommandAPDU(apdu), card)
	pom = pom[4:]
	fmt.Printf("\n\n (Prvi podatak) : %s", string(pom[4:pom[2]+4]))

	lk.adresaMesto = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nPREZIME : %s", string(pom[4:pom[2]+4]))
	lk.adresaOpstina = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nIME : %s", string(pom[4:pom[2]+4]))
	lk.adresaUlica = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	fmt.Printf("\n\nIME OCA : %s", string(pom[4:pom[2]+4]))
	lk.imeRoditelja = string(pom[4 : pom[2]+4])
	pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nSPRAT: %s", string(pom[4:pom[2]+4]))
	// pom = pom[4+pom[2]:]
	// fmt.Printf("\n\nSPRAT: %s", string(pom[4:pom[2]+4]))
}
func pretyDate(date string) string {

	return date[:2] + "." + date[2:4] + "." + date[4:]

}
