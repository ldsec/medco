package handlers

import (
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_network"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
)

// MedCoNetworkGetMetadataHandler handles /medco/network API endpoint
func MedCoNetworkGetMetadataHandler(params medco_network.GetMetadataParams, principal *models.User) middleware.Responder {

	nodes := make([]*medco_network.GetMetadataOKBodyNodesItems0, 0)
	for idx, url := range utilserver.MedCoNodesURL {
		idxInt64 := int64(idx)
		nodes = append(nodes, &medco_network.GetMetadataOKBodyNodesItems0{
			Index: &idxInt64,
			Name:  "Node " + strconv.Itoa(idx), // todo: config to specify node name
			URL:   url,
		})
	}

	pubKey, err := unlynx.GetCothorityKey()
	if err != nil {
		return medco_network.NewGetMetadataDefault(500).WithPayload(&medco_network.GetMetadataDefaultBody{
			Message: "Could not get cothority key: " + err.Error(),
		})
	}

	medcoNodeIdxInt64 := int64(utilserver.MedCoNodeIdx)
	return medco_network.NewGetMetadataOK().WithPayload(&medco_network.GetMetadataOKBody{
		NodeIndex: &medcoNodeIdxInt64,
		Nodes:     nodes,
		PublicKey: pubKey,
	})
}
