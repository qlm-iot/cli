package main

import (
	"github.com/qlm-iot/qlm/df"
	"github.com/qlm-iot/qlm/mi"
)

func createEmptyReadRequest() []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read:    &mi.ReadRequest{},
	})
	return ret
}

func createQLMMessage(id, name string) string {
	objects := df.Objects{
		Objects: []df.Object{
			df.Object{
				Id: &df.QLMID{Text: id},
				InfoItems: []df.InfoItem{
					df.InfoItem{
						Name: name,
					},
				},
			},
		},
	}
	data, _ := df.Marshal(objects)
	return (string)(data)
}

func createQLMMessageWithValue(id, name, value string) string {
	objects := df.Objects{
		Objects: []df.Object{
			df.Object{
				Id: &df.QLMID{Text: id},
				InfoItems: []df.InfoItem{
					df.InfoItem{
						Name: name,
						Values: []df.Value{
							df.Value{
								Text: value,
							},
						},
					},
				},
			},
		},
	}
	data, _ := df.Marshal(objects)
	return (string)(data)
}

func createReadRequest(id, name string) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat: "odf",
			Message: &mi.Message{
				Data: createQLMMessage(id, name),
			},
		},
	})
	return ret
}

func createSubscriptionRequest(id, name string, interval float64) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat: "odf",
			Interval:  interval,
			Message: &mi.Message{
				Data: createQLMMessage(id, name),
			},
		},
	})
	return ret
}

func createReadSubscriptionRequest(requestId string) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Read: &mi.ReadRequest{
			MsgFormat: "odf",
			RequestIds: []mi.Id{
				mi.Id{Text: requestId},
			},
		},
	})
	return ret
}

func createCancelSubscriptionRequest(requestId string) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Cancel: &mi.CancelRequest{
			RequestIds: []mi.Id{
				mi.Id{Text: requestId},
			},
		},
	})
	return ret
}

func createWriteRequest(id, name, value string) []byte {
	ret, _ := mi.Marshal(mi.OmiEnvelope{
		Version: "1.0",
		Ttl:     -1,
		Write: &mi.WriteRequest{
			MsgFormat:  "odf",
			TargetType: "device",
			Message: &mi.Message{
				Data: createQLMMessageWithValue(id, name, value),
			},
		},
	})
	return ret
}
