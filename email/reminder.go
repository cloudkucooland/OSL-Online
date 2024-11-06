package email

import (
	"log/slog"
	"strings"

	"github.com/cloudkucooland/OSL-Online/model"
	"github.com/matcornic/hermes/v2"
)

const reminder = `Membership in the Order of Saint Luke is maintained by annually reaffirming the vows of the Order. This reaffirmation is done by communicating the reaffirmation of vows to the Chancellor General each year, typically in October near the Feast of Saint Luke. Members are also encouraged to publicly reaffirm their vows in their chapter or at a retreat. These public reaffirmations do not replace the communication to the Chancellor General.

Our records show your last reaffirmation was more than a year ago. We encourage you to take a moment and fill out the online reaffirmation form.

If you no longer wish to reaffirm the vows you can reply to this email and inform the Chancellor General. You can choose to remain a friend of the Order and still receive communication from us or be removed from our rolls.

Membership in the order does not require financial donations. If you choose to donate you may opt to receive the periodicals in printed form. Those who do not donate will receive the periodicals via email.`

func SendReminder(id model.MemberID) error {
	h, err := setup()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	member, err := id.Get(true)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if member.PrimaryEmail == "" {
		return nil
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name:   member.OSLShortName(),
			Intros: strings.Split(reminder, "\n"),
			Actions: []hermes.Action{
				{
					Instructions: "To reaffirm, please click here:",
					Button: hermes.Button{
						Text: "Online Reaffirmation",
						Link: "https://saint-luke.net/reaffirmation/",
					},
				},
				{
					Button: hermes.Button{
						Text: "Printable Reaffirmation Form",
						Link: "https://saint-luke.net/wp-content/uploads/2024/09/OSL-2024-Reaffirmation-GENERIC.pdf",
					},
				},
			},
			Outros: []string{
				"Living the sacramental life",
			},
		},
	}

	body, err := h.GenerateHTML(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	text, err := h.GeneratePlainText(e)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := send(member.PrimaryEmail, "OSL Reaffirmation Reminder", body, text); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}