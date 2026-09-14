package logic

import (
	"context"
	"encoding/json"

	"backend/internal/errorx"
	"backend/internal/model"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateWorkflowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateWorkflowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateWorkflowLogic {
	return &ValidateWorkflowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateWorkflowLogic) ValidateWorkflow(req *types.ValidateWorkflowReq) (resp *types.ValidateWorkflowResp, err error) {
	workflow, err := l.svcCtx.WorkflowModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NewNotFoundError("workflow not found")
		}
		l.Logger.Errorf("find workflow failed: %v", err)
		return nil, errorx.NewInternalError("failed to find workflow")
	}

	var nodes []types.Node
	if workflow.Nodes.Valid && workflow.Nodes.String != "" {
		if err := json.Unmarshal([]byte(workflow.Nodes.String), &nodes); err != nil {
			return &types.ValidateWorkflowResp{
				Valid:  false,
				Errors: []types.ValidationError{{Code: "INVALID_NODES", Message: "failed to parse nodes JSON"}},
			}, nil
		}
	}

	var edges []types.Edge
	if workflow.Edges.Valid && workflow.Edges.String != "" {
		if err := json.Unmarshal([]byte(workflow.Edges.String), &edges); err != nil {
			return &types.ValidateWorkflowResp{
				Valid:  false,
				Errors: []types.ValidationError{{Code: "INVALID_EDGES", Message: "failed to parse edges JSON"}},
			}, nil
		}
	}

	errors := validateGraph(nodes, edges, workflow.EntryNodeId.String)

	return &types.ValidateWorkflowResp{
		Valid:  len(errors) == 0,
		Errors: errors,
	}, nil
}

func validateGraph(nodes []types.Node, edges []types.Edge, entryNodeId string) []types.ValidationError {
	var errors []types.ValidationError

	if len(nodes) == 0 {
		errors = append(errors, types.ValidationError{
			Code:    "NO_NODES",
			Message: "workflow must have at least one node",
		})
		return errors
	}

	nodeSet := make(map[string]bool)
	nodeTypes := make(map[string]string)
	for _, node := range nodes {
		if node.Id == "" {
			errors = append(errors, types.ValidationError{
				Code:    "EMPTY_NODE_ID",
				Message: "node has empty id",
			})
			continue
		}
		if nodeSet[node.Id] {
			errors = append(errors, types.ValidationError{
				Code:    "DUPLICATE_NODE_ID",
				Message: "duplicate node id: " + node.Id,
				NodeId:  node.Id,
			})
		}
		nodeSet[node.Id] = true
		nodeTypes[node.Id] = node.Type

		validTypes := map[string]bool{"script": true, "http": true, "human": true}
		if !validTypes[node.Type] {
			errors = append(errors, types.ValidationError{
				Code:    "INVALID_NODE_TYPE",
				Message: "invalid node type: " + node.Type + ", expected script|http|human",
				NodeId:  node.Id,
			})
		}
	}

	edgeSet := make(map[string]bool)
	outgoing := make(map[string][]string)
	incoming := make(map[string][]string)

	for _, edge := range edges {
		if edge.Id == "" {
			errors = append(errors, types.ValidationError{
				Code:    "EMPTY_EDGE_ID",
				Message: "edge has empty id",
			})
			continue
		}
		if edgeSet[edge.Id] {
			errors = append(errors, types.ValidationError{
				Code:    "DUPLICATE_EDGE_ID",
				Message: "duplicate edge id: " + edge.Id,
				EdgeId:  edge.Id,
			})
		}
		edgeSet[edge.Id] = true

		if !nodeSet[edge.Source] {
			errors = append(errors, types.ValidationError{
				Code:    "DANGLING_EDGE_SOURCE",
				Message: "edge source node not found: " + edge.Source,
				EdgeId:  edge.Id,
			})
		}
		if !nodeSet[edge.Target] {
			errors = append(errors, types.ValidationError{
				Code:    "DANGLING_EDGE_TARGET",
				Message: "edge target node not found: " + edge.Target,
				EdgeId:  edge.Id,
			})
		}

		validOutlets := map[string]bool{"success": true, "failure": true}
		if !validOutlets[edge.Outlet] {
			errors = append(errors, types.ValidationError{
				Code:    "INVALID_OUTLET",
				Message: "invalid edge outlet: " + edge.Outlet + ", expected success|failure",
				EdgeId:  edge.Id,
			})
		}

		outgoing[edge.Source] = append(outgoing[edge.Source], edge.Target)
		incoming[edge.Target] = append(incoming[edge.Target], edge.Source)
	}

	if entryNodeId == "" {
		if len(nodes) > 0 {
			foundEntry := false
			for _, node := range nodes {
				if len(incoming[node.Id]) == 0 {
					if !foundEntry {
						foundEntry = true
					} else {
						errors = append(errors, types.ValidationError{
							Code:    "MISSING_ENTRY_NODE",
							Message: "multiple root nodes found, entry_node_id must be specified",
						})
						break
					}
				}
			}
			if !foundEntry && len(nodes) > 0 {
				errors = append(errors, types.ValidationError{
					Code:    "NO_ROOT_NODE",
					Message: "no root node found (all nodes have incoming edges)",
				})
			}
		}
	} else {
		if !nodeSet[entryNodeId] {
			errors = append(errors, types.ValidationError{
				Code:    "INVALID_ENTRY_NODE",
				Message: "entry_node_id refers to non-existent node: " + entryNodeId,
			})
		}
	}

	if hasCycle(nodes, outgoing) {
		errors = append(errors, types.ValidationError{
			Code:    "CYCLE_DETECTED",
			Message: "workflow graph contains a cycle",
		})
	}

	if entryNodeId != "" && nodeSet[entryNodeId] {
		reachable := getReachableNodes(entryNodeId, outgoing)
		for _, node := range nodes {
			if !reachable[node.Id] && node.Id != entryNodeId {
				errors = append(errors, types.ValidationError{
					Code:    "UNREACHABLE_NODE",
					Message: "node is not reachable from entry node: " + node.Id,
					NodeId:  node.Id,
				})
			}
		}
	}

	return errors
}

func hasCycle(nodes []types.Node, outgoing map[string][]string) bool {
	white := make(map[string]bool)
	gray := make(map[string]bool)
	black := make(map[string]bool)

	for _, node := range nodes {
		white[node.Id] = true
	}

	var dfs func(nodeId string) bool
	dfs = func(nodeId string) bool {
		white[nodeId] = false
		gray[nodeId] = true

		for _, neighbor := range outgoing[nodeId] {
			if gray[neighbor] {
				return true
			}
			if white[neighbor] {
				if dfs(neighbor) {
					return true
				}
			}
		}

		gray[nodeId] = false
		black[nodeId] = true
		return false
	}

	for _, node := range nodes {
		if white[node.Id] {
			if dfs(node.Id) {
				return true
			}
		}
	}

	return false
}

func getReachableNodes(startNode string, outgoing map[string][]string) map[string]bool {
	reachable := make(map[string]bool)
	var visit func(nodeId string)
	visit = func(nodeId string) {
		if reachable[nodeId] {
			return
		}
		reachable[nodeId] = true
		for _, neighbor := range outgoing[nodeId] {
			visit(neighbor)
		}
	}
	visit(startNode)
	return reachable
}
