package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GeneratePDF(lk *LicnaKarta) {
	t := time.Now()
	timestamp := fmt.Sprintf("Datum štampe:     %02d.%02d.%d",
		t.Day(), t.Month(), t.Year())

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	// m.SetBorder(true)
	m.SetPageMargins(20, 20, 18)

	m.AddUTF8Font("CustomArial", consts.Normal, "arialUnicodeFont.ttf")
	m.AddUTF8Font("CustomArial", consts.Italic, "arialUnicodeFont.ttf")
	m.AddUTF8Font("CustomArial", consts.Bold, "arialUnicodeFont.ttf")
	m.AddUTF8Font("CustomArial", consts.BoldItalic, "arialUnicodeFont.ttf")
	m.SetDefaultFontFamily("CustomArial")

	mestoOpstinaDrzava := lk.mestoRodjenja + ", " + lk.opstinaRodjenja + ", " + lk.drzavaRodjenja
	// adresaFull := lk.adresaMesto + ", " + lk.adresaOpstina + ", " + lk.adresaUlica
	adresaFull := lk.prebivaliste

	fmt.Print(adresaFull)

	m.Line(0)
	m.Row(10, func() {
		m.Col(10, func() {
			m.Text("    \u010cITA\u010c ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA ", props.Text{
				Top:         2,
				Size:        13,
				Extrapolate: true,
				// Align: consts.Center,
			})
		})
		m.ColSpace(0)
	})
	m.Line(0)

	base64image := base64.StdEncoding.EncodeToString(lk.slika)

	m.Row(63, func() {
		m.Col(4, func() {
			_ = m.Base64Image(base64image, consts.Jpg, props.Rect{
				Top:     4,
				Center:  false,
				Percent: 93,
			})
		})
		m.ColSpace(4)
	})

	m.Line(7.0, props.Line{
		Width: 0.4,
	})
	m.Row(1, func() {
		m.Col(10, func() {
			m.Text("Podaci o građaninu", props.Text{
				Size: 11,
				// Top:  0,
			})
		})
	})
	m.Line(9.0, props.Line{
		Width: 0.5,
	})

	m.Row(73, func() {
		m.Col(3, func() {
			m.Text("Prezime: ", props.Text{
				Size: 11,
			})
			m.Text("Ime: ", props.Text{
				Top:  8,
				Size: 11,
			})
			m.Text("Ime jednog roditelja: ", props.Text{
				Top:  16,
				Size: 11,
			})
			m.Text("Datum rođenja: ", props.Text{
				Top:  24,
				Size: 11,
			})
			m.Text("Mesto rođenja,         opština i država:", props.Text{
				Top:  32,
				Size: 11,
			})
			m.Text("Prebivalište i adresa     stana:", props.Text{
				Top:  44,
				Size: 11,
			})
			m.Text("JMBG:", props.Text{
				Top:  56,
				Size: 11,
			})
			m.Text("Pol:", props.Text{
				Top:  66,
				Size: 11,
			})
		})
		m.Col(8, func() {
			m.Text(lk.prezime, props.Text{
				Size: 11,
			})
			m.Text(lk.ime, props.Text{
				Size: 11,
				Top:  8,
			})
			m.Text(lk.imeRoditelja, props.Text{
				Size: 11,
				Top:  16,
			})
			m.Text(lk.datumRodjenja, props.Text{
				Size: 11,
				Top:  24,
			})
			m.Text(mestoOpstinaDrzava, props.Text{
				Size: 11,
				Top:  32,
			})
			m.Text(adresaFull, props.Text{
				Size: 11,
				Top:  44,
			})
			m.Text(lk.JMBG, props.Text{
				Size: 11,
				Top:  56,
			})
			m.Text(lk.pol, props.Text{
				Size: 11,
				Top:  66,
			})
		})
		m.ColSpace(30)
	})

	m.Line(3.0, props.Line{
		Width: 0.4,
	})
	m.Row(2, func() {
		m.Col(10, func() {
			m.Text("Podaci o dokumentu", props.Text{
				Size: 11,
				// Top:  0,
			})
		})
	})
	m.Line(9.0, props.Line{
		Width: 0.4,
	})

	m.Row(30, func() {
		m.Col(3, func() {
			m.Text("Dokument izdaje: ", props.Text{
				Size: 11,
			})
			m.Text("Broj dokumenta: ", props.Text{
				Top:  8,
				Size: 11,
			})
			m.Text("Datum izdavanja: ", props.Text{
				Top:  16,
				Size: 11,
			})
			m.Text("Važi do: ", props.Text{
				Top:  24,
				Size: 11,
			})
		})
		m.Col(8, func() {
			m.Text(lk.dokumentIzdaje, props.Text{
				Size: 11,
			})
			m.Text(lk.brojDokumenta, props.Text{
				Size: 11,
				Top:  8,
			})
			m.Text(lk.datumIzdavanja, props.Text{
				Size: 11,
				Top:  16,
			})
			m.Text(lk.vaziDo, props.Text{
				Size: 11,
				Top:  24,
			})
		})
		m.ColSpace(30)
	})
	// m.SetBorder(true)
	m.Line(2.0, props.Line{
		Width: 0.4,
	})
	m.Line(0.5, props.Line{
		Width: 0.4,
	})

	m.Row(25, func() {
		m.Col(5, func() {
			m.Text(timestamp, props.Text{
				Top:  2,
				Size: 11,
				// Extrapolate: true,
			})
		})
	})

	m.Line(0.5, props.Line{
		Width: 0.3,
	})

	m.Row(17, func() {
		m.Col(12, func() {
			m.Text("1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.", props.Text{
				Top:  2,
				Size: 9,
				// Extrapolate: true,
			})
			m.Text("2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu skraćuje se na prva dva karaktera.", props.Text{
				Top:  9,
				Size: 9,
				// Extrapolate: true,
			})
		})
	})
	m.Row(1, func() {
		// m.Line(4.5)
		m.Line(4.5, props.Line{
			Width: 0.3,
		})
	})

	err := m.OutputFileAndClose("zpl.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}
}
