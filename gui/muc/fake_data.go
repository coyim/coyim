package muc

func fakeAccounts() []*account {
	accounts := []*account{
		&account{
			rosterItem: &rosterItem{
				id:     "sandy@autonomia.digital",
				status: statusOnline,
			},
			contacts: []*contact{
				&contact{
					rosterItem: &rosterItem{
						id:     "pedro@autonomia.digital",
						name:   "Pedro Enrique",
						status: statusOnline,
					},
				},
				&contact{
					rosterItem: &rosterItem{
						id:     "rafael@autonomia.digital",
						status: statusOnline,
					},
				},
				&contact{
					rosterItem: &rosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: statusOffline,
					},
				},
			},
			rooms: []string{
				"#coyim:matrix:autonomia.digital",
				"#wahay:matrix:autonomia.digital",
			},
		},
		&account{
			rosterItem: &rosterItem{
				id:     "pedro@autonomia.digital",
				status: statusOnline,
			},
			contacts: []*contact{
				&contact{
					rosterItem: &rosterItem{
						id:     "sandy@autonomia.digital",
						name:   "Sandy Acurio",
						status: statusOnline,
					},
				},
				&contact{
					rosterItem: &rosterItem{
						id:     "rafael@autonomia.digital",
						status: statusOnline,
					},
				},
				&contact{
					rosterItem: &rosterItem{
						id:     "cristina@autonomia.digital",
						name:   "Cristina Salcedo",
						status: statusOffline,
					},
				},
			},
			rooms: []string{
				"#main:matrix:autonomia.digital",
			},
		},
		&account{
			rosterItem: &rosterItem{
				id:     "pedro@coy.im",
				name:   "Pedro CoyIM",
				status: statusOffline,
			},
		},
	}

	return accounts
}

func fakeRooms() map[string]*room {
	rooms := map[string]*room{
		"#coyim:matrix:autonomia.digital": &room{
			name: "CoyIM",
		},
		"#wahay:matrix:autonomia.digital": &room{
			name: "Wahay",
		},
		"#main:matrix:autonomia.digital": &room{
			name: "Main",
		},
	}

	return rooms
}
