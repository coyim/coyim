package muc

func fakeAccounts() []*mucAccount {
	accounts := []*mucAccount{
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "sandy@autonomia.digital",
				status: mucStatusOnline,
			},
			contacts: []*mucAccountContact{
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "pedro@autonomia.digital",
						name:   "Pedro Enrique",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "rafael@autonomia.digital",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: mucStatusOffline,
					},
				},
			},
			rooms: []string{
				"#coyim:matrix:autonomia.digital",
				"#wahay:matrix:autonomia.digital",
			},
		},
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "pedro@autonomia.digital",
				status: mucStatusOnline,
			},
			contacts: []*mucAccountContact{
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "sandy@autonomia.digital",
						name:   "Sandy Acurio",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "rafael@autonomia.digital",
						status: mucStatusOnline,
					},
				},
				&mucAccountContact{
					mucRosterItem: &mucRosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: mucStatusOffline,
					},
				},
			},
			rooms: []string{
				"#main:matrix:autonomia.digital",
			},
		},
		&mucAccount{
			mucRosterItem: &mucRosterItem{
				id:     "pedro@coy.im",
				name:   "Pedro CoyIM",
				status: mucStatusOffline,
			},
		},
	}

	return accounts
}

func fakeRooms() map[string]*mucRoom {
	rooms := map[string]*mucRoom{
		"#coyim:matrix:autonomia.digital": &mucRoom{
			name: "CoyIM",
		},
		"#wahay:matrix:autonomia.digital": &mucRoom{
			name: "Wahay",
		},
		"#main:matrix:autonomia.digital": &mucRoom{
			name: "Main",
		},
	}

	return rooms
}
