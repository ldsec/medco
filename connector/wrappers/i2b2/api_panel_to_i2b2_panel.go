package i2b2

import (
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"
)

func apiPanelToI2b2Panel(queryPanel *models.Panel) Panel {
	invert := "0"
	if *queryPanel.Not {
		invert = "1"
	}

	i2b2Panel := Panel{
		PanelAccuracyScale:   "100",
		Invert:               invert,
		PanelTiming:          strings.ToUpper(string(queryPanel.PanelTiming)),
		TotalItemOccurrences: "1",
	}

	for _, queryItem := range queryPanel.ConceptItems {

		i2b2Item := Item{
			ItemKey: convertPathToI2b2Format(*queryItem.QueryTerm),
		}
		if queryItem.Operator != "" && queryItem.Modifier == nil {
			i2b2Item.ConstrainByValue = &ConstrainByValue{
				ValueType:       queryItem.Type,
				ValueOperator:   queryItem.Operator,
				ValueConstraint: queryItem.Value,
			}
		}
		if queryItem.Modifier != nil {
			i2b2Item.ConstrainByModifier = &ConstrainByModifier{
				AppliedPath: strings.ReplaceAll(*queryItem.Modifier.AppliedPath, "/", `\`),
				ModifierKey: convertPathToI2b2Format(*queryItem.Modifier.ModifierKey),
			}
			if queryItem.Operator != "" {
				i2b2Item.ConstrainByModifier.ConstrainByValue = &ConstrainByValue{
					ValueType:       queryItem.Type,
					ValueOperator:   queryItem.Operator,
					ValueConstraint: queryItem.Value,
				}
			}
		}
		i2b2Panel.Items = append(i2b2Panel.Items, i2b2Item)

	}
	for _, cohort := range queryPanel.CohortItems {

		i2b2Item := Item{
			ItemKey: cohort,
		}
		i2b2Panel.Items = append(i2b2Panel.Items, i2b2Item)
	}
	return i2b2Panel
}
