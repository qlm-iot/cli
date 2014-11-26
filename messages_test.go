package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateReadRequest(t *testing.T) {
	data := createReadRequest("oven1", "temperature")
	expected := `<omiEnvelope version="1.0" ttl="-1">
    <read msgformat="odf">
        <msg><Objects xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="odf.xsd">
    <Object>
        <id>oven1</id>
        <InfoItem name="temperature"></InfoItem>
    </Object>
</Objects></msg>
    </read>
</omiEnvelope>`
	assert.Equal(t, expected, string(data))
}
