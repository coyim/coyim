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
			rooms: []*accountRoom{
				&accountRoom{
					id:     "#coyim:matrix:autonomia.digital",
					status: statusOnline,
				},
				&accountRoom{
					id:     "#wahay:matrix:autonomia.digital",
					status: statusOnline,
				},
			},
			groups: []*group{
				&group{
					rosterItem: &rosterItem{
						id:   "autonomia.digital",
						name: "CAD",
					},
					rooms: []*accountRoom{
						&accountRoom{
							id:     "#main:matrix:autonomia.digital",
							status: statusOnline,
						},
						&accountRoom{
							id:     "#admin:matrix:autonomia.digital",
							status: statusOnline,
						},
					},
					contacts: []*contact{
						&contact{
							rosterItem: &rosterItem{
								id:     "pedro@coy.im",
								name:   "Pedro CoyIM",
								status: statusOnline,
							},
						},
						&contact{
							rosterItem: &rosterItem{
								id:     "ola@coy.im",
								name:   "Ola Bini",
								status: statusOnline,
							},
						},
						&contact{
							rosterItem: &rosterItem{
								id:     "sandy@coy.im",
								status: statusOffline,
							},
						},
					},
				},
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
			rooms: []*accountRoom{
				&accountRoom{
					id:     "#main:matrix:autonomia.digital",
					status: statusOffline,
				},
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
		"#admin:matrix:autonomia.digital": &room{
			name: "Administration",
		},
	}

	return rooms
}
