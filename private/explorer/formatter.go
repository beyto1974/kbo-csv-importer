package main

import (
	"fmt"
	"strings"
)

// FormatEnterpriseLLM produces a comprehensive, structured plain-text
// description of an enterprise and all its related entities, suitable
// for consumption by an LLM or human reader.
func FormatEnterpriseLLM(b *EntityBundle) string {
	var sb strings.Builder

	// ── Enterprise ──────────────────────────────────────────────
	e := b.Enterprise
	sb.WriteString("# Enterprise\n\n")

	writeField(&sb, "Enterprise Number", e.EnterpriseNumber)

	if e.StatusCode.Description != "" {
		writeField(&sb, "Status", fmt.Sprintf("%s (%s)", e.StatusCode.Description, e.Status))
	} else {
		writeField(&sb, "Status", e.Status)
	}

	if e.JuridicalSituationCode.Description != "" {
		writeField(&sb, "Juridical Situation", fmt.Sprintf("%s (%s)", e.JuridicalSituationCode.Description, e.JuridicalSituation))
	} else {
		writeField(&sb, "Juridical Situation", e.JuridicalSituation)
	}

	if e.TypeOfEnterpriseCode.Description != "" {
		writeField(&sb, "Type of Enterprise", fmt.Sprintf("%s (%s)", e.TypeOfEnterpriseCode.Description, e.TypeOfEnterprise))
	} else {
		writeField(&sb, "Type of Enterprise", e.TypeOfEnterprise)
	}

	if e.JuridicalFormCode.Description != "" {
		writeField(&sb, "Juridical Form", fmt.Sprintf("%s (%s)", e.JuridicalFormCode.Description, e.JuridicalForm))
	} else {
		writeField(&sb, "Juridical Form", e.JuridicalForm)
	}

	if e.JuridicalFormCacCode.Description != "" {
		writeField(&sb, "Juridical Form (CAC)", fmt.Sprintf("%s (%s)", e.JuridicalFormCacCode.Description, e.JuridicalFormCac))
	} else if e.JuridicalFormCac != "" {
		writeField(&sb, "Juridical Form (CAC)", e.JuridicalFormCac)
	}

	if e.StartDate.Time.IsZero() {
		writeField(&sb, "Start Date", "")
	} else {
		writeField(&sb, "Start Date", e.StartDate.Time.Format("2006-01-02"))
	}

	// ── Denominations (names) ───────────────────────────────────
	if len(b.Denominations) > 0 {
		sb.WriteString("\n## Denominations / Names\n\n")
		for _, d := range b.Denominations {
			name := ""
			if d.Denomination != nil {
				name = *d.Denomination
			}
			typeLabel := d.TypeOfDenominationCode.Description
			if typeLabel == "" {
				typeLabel = d.TypeOfDenomination
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s (type: %s, entity: %s)\n",
				d.Language, name, typeLabel, d.EntityNumber))
		}
	}

	// ── Establishments ──────────────────────────────────────────
	if len(b.Establishments) > 0 {
		sb.WriteString("\n## Establishments\n\n")
		for i, est := range b.Establishments {
			startDate := ""
			if !est.StartDate.Time.IsZero() {
				startDate = est.StartDate.Time.Format("2006-01-02")
			}
			sb.WriteString(fmt.Sprintf("%d. Establishment %s (start: %s, enterprise: %s)\n",
				i+1, est.EstablishmentNumber, startDate, est.EnterpriseNumber))
		}
	}

	// ── Addresses ───────────────────────────────────────────────
	if len(b.Addresses) > 0 {
		sb.WriteString("\n## Addresses\n\n")
		for i, a := range b.Addresses {
			typeLabel := a.TypeOfAddressCode.Description
			if typeLabel == "" {
				typeLabel = a.TypeOfAddress
			}
			sb.WriteString(fmt.Sprintf("### Address %d (entity: %s, type: %s)\n", i+1, a.EntityNumber, typeLabel))
			street := coalesce(a.StreetFr, a.StreetNl)
			municipality := coalesce(a.MunicipalityFr, a.MunicipalityNl)
			country := coalesce(a.CountryFr, a.CountryNl)

			addr := street
			if a.HouseNumber != "" {
				addr += " " + a.HouseNumber
			}
			if a.Box != "" {
				addr += " box " + a.Box
			}
			if a.Zipcode != "" || municipality != "" {
				addr += ", " + a.Zipcode + " " + municipality
			}
			if country != "" {
				addr += ", " + country
			}
			sb.WriteString(fmt.Sprintf("  %s\n", strings.TrimSpace(addr)))

			if a.ExtraAddressInfo != "" {
				sb.WriteString(fmt.Sprintf("  Extra: %s\n", a.ExtraAddressInfo))
			}
			if !a.DateStrikingOff.Time.IsZero() {
				sb.WriteString(fmt.Sprintf("  Struck off: %s\n", a.DateStrikingOff.Time.Format("2006-01-02")))
			}
		}
	}

	// ── Contacts ────────────────────────────────────────────────
	if len(b.Contacts) > 0 {
		sb.WriteString("\n## Contacts\n\n")
		for _, c := range b.Contacts {
			typeLabel := c.ContactTypeCode.Description
			if typeLabel == "" {
				typeLabel = c.ContactType
			}
			sb.WriteString(fmt.Sprintf("- %s: %s (contact-of: %s, entity: %s)\n",
				typeLabel, c.Value, c.EntityContact, c.EntityNumber))
		}
	}

	// ── Activities ──────────────────────────────────────────────
	if len(b.Activities) > 0 {
		sb.WriteString("\n## Activities / NACE Codes\n\n")
		for _, act := range b.Activities {
			parts := []string{}

			groupLabel := act.ActivityGroupCode.Description
			if groupLabel == "" {
				groupLabel = act.ActivityGroup
			}
			parts = append(parts, fmt.Sprintf("Group: %s", groupLabel))

			naceLabel := act.NaceCodeCode.Description
			if naceLabel == "" {
				naceLabel = act.NaceCode
			}
			parts = append(parts, fmt.Sprintf("NACE %s: %s", act.NaceCode, naceLabel))

			if act.NaceVersionCode.Description != "" {
				parts = append(parts, fmt.Sprintf("Version: %s", act.NaceVersionCode.Description))
			}

			if act.ClassificationCode.Description != "" {
				parts = append(parts, fmt.Sprintf("Classification: %s", act.ClassificationCode.Description))
			}

			sb.WriteString(fmt.Sprintf("- [%s] %s\n", act.EntityNumber, strings.Join(parts, " | ")))
		}
	}

	// ── Branches ────────────────────────────────────────────────
	if len(b.Branches) > 0 {
		sb.WriteString("\n## Branches\n\n")
		for _, br := range b.Branches {
			startDate := ""
			if !br.StartDate.Time.IsZero() {
				startDate = br.StartDate.Time.Format("2006-01-02")
			}
			sb.WriteString(fmt.Sprintf("- Branch %s (start: %s)\n", br.ID, startDate))
		}
	}

	return sb.String()
}

func writeField(sb *strings.Builder, label, value string) {
	if value == "" {
		value = "—"
	}
	sb.WriteString(fmt.Sprintf("- **%s**: %s\n", label, value))
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
