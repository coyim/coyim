package gui

import "github.com/coyim/gotk3adapter/gtki"

type infoBarColor struct {
	background string
	titleColor string
}

type infoBarColorStyles map[gtki.MessageType]infoBarColor

func newInfoBarColorStyles(c mucColorSet) infoBarColorStyles {
	return infoBarColorStyles{
		gtki.MESSAGE_INFO: infoBarColor{
			background: "linear-gradient(0deg, rgba(14,116,144,1) 0%, rgba(8,145,178,1) 100%)",
			titleColor: "#ECFEFF",
		},
		gtki.MESSAGE_WARNING: infoBarColor{
			background: "linear-gradient(0deg, rgba(234,88,12,1) 0%, rgba(249,115,22,1) 100%)",
			titleColor: "#FFF7ED",
		},
		gtki.MESSAGE_QUESTION: infoBarColor{
			background: "linear-gradient(0deg, rgba(153,27,27,1) 0%, rgba(185,28,28,1) 100%)",
			titleColor: "#FEFCE8",
		},
		gtki.MESSAGE_ERROR: infoBarColor{
			background: "linear-gradient(0deg, rgba(136,19,55,1) 0%, rgba(159,18,57,1) 100%)",
			titleColor: "#FFF1F2",
		},
		gtki.MESSAGE_OTHER: infoBarColor{
			background: "linear-gradient(0deg, rgba(6,95,70,1) 0%, rgba(4,120,87,1) 100%)",
			titleColor: "#F0FDFA",
		},
	}
}

func (s *mucStylesProvider) setInfoBarStyle(ib gtki.InfoBar) {
	if st, ok := s.infoBarColorStyles[ib.GetMessageType()]; ok {
		s.setWidgetStyles(ib, styles{
			".infobar": style{
				"background":  st.background,
				"text-shadow": "none",
				"font-weight": "500",
				"padding":     "8px 10px",
			},
			".infobar .content": style{
				"text-shadow": "none",
			},
			".infobar .title": style{
				"color":       st.titleColor,
				"text-shadow": "none",
			},
		})
	}
}
