package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_network"
	"github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

// MedCoNetworkGetMetadataHandler handles /medco/network API endpoint
func MedCoNetworkGetMetadataHandler(params medco_network.GetMetadataParams, principal *models.User) middleware.Responder {

	nodes := make([]*medco_network.NodesItems0, 0)
	for idx, url := range utilserver.MedCoNodesURL {
		nodes = append(nodes, &medco_network.NodesItems0{
			Index: int64(idx),
			Name: "Node " + string(idx), // todo: config to specify node name
			URL: url,
		})
	}

	pubKey, err := unlynx.GetCothorityKey()
	if err != nil {
		return medco_network.NewGetMetadataDefault(500).WithPayload(&medco_network.GetMetadataDefaultBody{
			Message: "Could not get cothority key: " + err.Error(),
		})
	}

	return medco_network.NewGetMetadataOK().WithPayload(&medco_network.GetMetadataOKBody{
		NodeIndex: int64(utilserver.MedCoNodeIdx),
		Nodes: nodes,
		PublicKey: pubKey,
	})
}
