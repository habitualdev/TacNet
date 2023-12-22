package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const addMarker = "/api/markers/add"
const getMarkers = "/api/markers"
const deleteMarker = "/api/markers/delete"

var sessionUuid = ""

func main() {
	a := app.New()
	sessionUuid = uuid.New().String()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("TacNet")
	w.Resize(fyne.NewSize(800, 600))
	habTacUser := widget.NewEntry()
	habTacUser.SetPlaceHolder("User")
	habTacSession := widget.NewEntry()
	habTacSession.SetPlaceHolder("Session")
	habTacUrl := widget.NewEntry()
	habTacUrl.SetPlaceHolder("Url")
	tabContainer := container.NewAppTabs()

	cffTargetName := widget.NewEntry()
	cffTargetName.SetPlaceHolder("Target Name")
	cffTargetGrid := widget.NewEntry()
	cffTargetGrid.SetPlaceHolder("Target Grid")
	cffTargetAltitude := widget.NewEntry()
	cffTargetAltitude.SetPlaceHolder("Target Altitude")
	cffComments := widget.NewMultiLineEntry()
	cffComments.SetPlaceHolder("Comments")
	cffSendButton := widget.NewButton("Send", func() {
		defer func() {
			if r := recover(); r != nil {
				a.SendNotification(fyne.NewNotification("Error", fmt.Sprintf("%v", r)))
			}
		}()
		northing := ""
		easting := ""
		switch len(cffTargetGrid.Text) {
		case 6:
			northing = cffTargetGrid.Text[:3] + "00"
			easting = cffTargetGrid.Text[3:] + "00"
		case 8:
			northing = cffTargetGrid.Text[:4] + "0"
			easting = cffTargetGrid.Text[4:] + "0"
		case 10:
			northing = cffTargetGrid.Text[:5]
			easting = cffTargetGrid.Text[5:]

		}
		reqForm := url.Values{
			"lat":      {easting},
			"lng":      {northing},
			"elv":      {cffTargetAltitude.Text},
			"title":    {cffTargetName.Text},
			"unit":     {"fm"},
			"comments": {strings.ReplaceAll(cffComments.Text, "\n", "<br>")},
			"session":  {habTacSession.Text},
			"id":       {uuid.New().String()},
		}

		_, err := http.PostForm(habTacUrl.Text+addMarker, reqForm)
		if err != nil {
			a.SendNotification(fyne.NewNotification("Error", err.Error()))
			return
		}

	})

	cffBox := container.NewVBox(cffTargetName, cffTargetGrid, cffTargetAltitude, cffComments, cffSendButton)

	tabContainer.Append(container.NewTabItem("CFF", cffBox))

	friendlyGrid := widget.NewEntry()
	friendlyGrid.SetPlaceHolder("Friendly Grid")
	friendlyComments := widget.NewMultiLineEntry()
	friendlyComments.SetPlaceHolder("Comments")
	friendlySendButton := widget.NewButton("Send", func() {
		if r := recover(); r != nil {
			a.SendNotification(fyne.NewNotification("Error", fmt.Sprintf("%v", r)))
		}
		northing := ""
		easting := ""
		switch len(friendlyGrid.Text) {
		case 6:
			northing = friendlyGrid.Text[:3] + "00"
			easting = friendlyGrid.Text[3:] + "00"
		case 8:
			northing = friendlyGrid.Text[:4] + "0"
			easting = friendlyGrid.Text[4:] + "0"
		case 10:
			northing = friendlyGrid.Text[:5]
			easting = friendlyGrid.Text[5:]

		}
		reqForm := url.Values{
			"lat":      {easting},
			"lng":      {northing},
			"title":    {habTacUser.Text},
			"unit":     {"1-8"},
			"comments": {strings.ReplaceAll(friendlyComments.Text, "\n", "<br>")},
			"session":  {habTacSession.Text},
			"id":       {uuid.New().String()},
		}
		_, err := http.Get(habTacUrl.Text + deleteMarker + "?session=" + habTacSession.Text + "&title=" + habTacUser.Text)
		if err != nil {
			a.SendNotification(fyne.NewNotification("Error", err.Error()))
			return
		}
		_, err = http.PostForm(habTacUrl.Text+addMarker, reqForm)
		if err != nil {
			a.SendNotification(fyne.NewNotification("Error", err.Error()))
			return
		}
	})

	friendlyBox := container.NewVBox(friendlyGrid, friendlyComments, friendlySendButton)

	tabContainer.Append(container.NewTabItem("Friendly", friendlyBox))

	habTacWindow := container.NewVBox(habTacUser, habTacSession, habTacUrl, tabContainer)

	observerPos := widget.NewEntry()
	observerPos.SetPlaceHolder("Observer Position")
	targetAz := widget.NewEntry()
	targetAz.SetPlaceHolder("Target Azimuth (Mils)")
	targetDist := widget.NewEntry()
	targetDist.SetPlaceHolder("Target Distance")
	targetPos := widget.NewEntry()
	targetPos.SetPlaceHolder("Target Position")

	calculatePolar := widget.NewButton("Calculate Polar From Observer", func() {
		targetPosCalc, err := PolarToGrid(targetDist.Text, targetAz.Text, observerPos.Text)
		if err != nil {
			a.SendNotification(fyne.NewNotification("Error", err.Error()))
			return
		}
		targetPos.SetText(targetPosCalc)
		targetPos.Refresh()
	})

	degreesToMilsCalcLabel := widget.NewLabel("")
	degreesToMilsCalcEntry := widget.NewEntry()
	degreesToMilsCalcEntry.SetPlaceHolder("Degrees")
	degreesToMilsCalcButton := widget.NewButton("Calculate", func() {
		degreesInt, err := strconv.Atoi(degreesToMilsCalcEntry.Text)
		if err != nil {
			return
		}
		mils := float64(degreesInt) * 17.777778

		degreesToMilsCalcLabel.SetText(fmt.Sprintf("%f", mils))
		degreesToMilsCalcLabel.Refresh()

	})

	utilityTab := container.NewTabItem("Utilities", container.NewVBox(container.NewVBox(widget.NewLabel("Polar from Observer Grid Calculator"),
		observerPos, targetAz, targetDist, calculatePolar, targetPos, widget.NewLabel("Degrees to Mils Conversions"), degreesToMilsCalcEntry,
		degreesToMilsCalcLabel, degreesToMilsCalcButton)))

	tabContainer.Append(utilityTab)

	w.SetContent(habTacWindow)
	w.ShowAndRun()
}
