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
					id:     "coyim:matrix:autonomia.digital",
					status: statusOnline,
				},
				&accountRoom{
					id:     "wahay:matrix:autonomia.digital",
					status: statusOnline,
				},
				&accountRoom{
					id:     "gtk-ui:matrix:autonomia.digital",
					status: statusOffline,
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
							id:     "main:matrix:autonomia.digital",
							status: statusOnline,
						},
						&accountRoom{
							id:     "admin:matrix:autonomia.digital",
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
					id:     "main:matrix:autonomia.digital",
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
		"coyim:matrix:autonomia.digital": &room{
			rosterItem: &rosterItem{
				name: "CoyIM",
			},
			members: membersList{
				"pedro": &member{
					rosterItem: &rosterItem{
						id:     "pedro@autonomia.digital",
						name:   "Pedro Enrique Palau",
						status: statusOnline,
					},
					role: roleAdministrator,
				},
				"sandy": &member{
					rosterItem: &rosterItem{
						id:     "sandy@autonomia.digital",
						name:   "Sandy Acurio",
						status: statusOnline,
					},
					role: roleModerator,
				},
				"reinaldo": &member{
					rosterItem: &rosterItem{
						id:     "reinaldo@autonomia.digital",
						status: statusOnline,
					},
					role: roleParticipant,
				},
				"cristian": &member{
					rosterItem: &rosterItem{
						id:     "cristian@autonomia.digital",
						status: statusOffline,
					},
					role: roleNone,
				},
				"sara": &member{
					rosterItem: &rosterItem{
						id:     "sara@autonomia.digital",
						name:   "Sara Zambrano",
						status: statusOffline,
					},
					role: roleVisitor,
				},
			},
		},
		"wahay:matrix:autonomia.digital": &room{
			rosterItem: &rosterItem{
				name: "Wahay",
			},
			members: membersList{
				"cristina": &member{
					rosterItem: &rosterItem{
						id:   "cristina@autonomia.digital",
						name: "Cristina",
					},
					role: roleAdministrator,
				},
				"alvaro": &member{
					rosterItem: &rosterItem{
						id:   "alvaro@autonomia.digital",
						name: "Alvaro",
					},
					role: roleModerator,
				},
				"mauro": &member{
					rosterItem: &rosterItem{
						id: "mauro@autonomia.digital",
					},
					role: roleParticipant,
				},
			},
		},
		"main:matrix:autonomia.digital": &room{
			rosterItem: &rosterItem{
				name: "Main",
			},
			description: "The Main room for CAD",
			members: membersList{
				"ola": &member{
					rosterItem: &rosterItem{
						id:   "ola@autonomia.digital",
						name: "Ola Bini",
					},
					role: roleAdministrator,
				},
				"sandy": &member{
					rosterItem: &rosterItem{
						id:   "sandy@autonomia.digital",
						name: "Sandy Acurio",
					},
					role: roleModerator,
				},
				"cristina": &member{
					rosterItem: &rosterItem{
						id:   "cristina@autonomia.digital",
						name: "Cristina",
					},
					role: roleParticipant,
				},
			},
		},
		"admin:matrix:autonomia.digital": &room{
			rosterItem: &rosterItem{
				name: "Administration",
			},
			members: membersList{
				"ola": &member{
					rosterItem: &rosterItem{
						id:   "ola@autonomia.digital",
						name: "Ola Bini",
					},
					role: roleAdministrator,
				},
				"sara": &member{
					rosterItem: &rosterItem{
						id:   "sara@autonomia.digital",
						name: "Sara",
					},
					role: roleAdministrator,
				},
				"alvaro": &member{
					rosterItem: &rosterItem{
						id:   "alvaro@autonomia.digital",
						name: "Alvaro Paredes",
					},
					role: roleModerator,
				},
			},
		},
		"gtk-ui:matrix:autonomia.digital": &room{
			members: membersList{
				"sandy": &member{
					rosterItem: &rosterItem{
						id:   "sandy@autonomia.digital",
						name: "Sandy Acurio",
					},
					role: roleAdministrator,
				},
				"pedro": &member{
					rosterItem: &rosterItem{
						id:   "pedro@autonomia.digital",
						name: "Pedro Palau",
					},
					role: roleAdministrator,
				},
			},
		},
	}

	return rooms
}
