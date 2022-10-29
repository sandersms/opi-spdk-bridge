// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/v1/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ////////////////////////////////////////////////////////
var subsystems = map[string]*pb.NVMeSubsystem{}

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("NVMeSubsystemCreate: Received from client: %v", in)
	params := NvmfCreateSubsystemParams{
		Nqn:          in.GetSubsystem().GetNqn(),
		SerialNumber: "SPDK0",
		AllowAnyHost: true,
	}
	var result NvmfCreateSubsystemResult
	err := call("nvmf_create_subsystem", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	subsystems[in.Subsystem.Id.Value] = in.Subsystem
	log.Printf("Received from SPDK: %v", result)
	return &pb.NVMeSubsystemCreateResponse{}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("NVMeSubsystemDelete: Received from client: %v", in)
	subsys, ok := subsystems[in.SubsystemId.Value]
	if !ok {
		err := fmt.Errorf("unable to find key %s", in.SubsystemId)
		log.Printf("error: %v", err)
		return nil, err
	}
	params := NvmfDeleteSubsystemParams{
		Nqn: subsys.Nqn,
	}
	var result NvmfDeleteSubsystemResult
	err := call("nvmf_delete_subsystem", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not delete: %v", in)
	}
	delete(subsystems, subsys.Id.Value)
	return &pb.NVMeSubsystemDeleteResponse{}, nil
}

func (s *server) NVMeSubsystemUpdate(ctx context.Context, in *pb.NVMeSubsystemUpdateRequest) (*pb.NVMeSubsystemUpdateResponse, error) {
	log.Printf("NVMeSubsystemUpdate: Received from client: %v", in)
	return &pb.NVMeSubsystemUpdateResponse{}, nil
}

func (s *server) NVMeSubsystemList(ctx context.Context, in *pb.NVMeSubsystemListRequest) (*pb.NVMeSubsystemListResponse, error) {
	log.Printf("NVMeSubsystemList: Received from client: %v", in)
	var result []NvmfGetSubsystemsResult
	err := call("nvmf_get_subsystems", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.NVMeSubsystem, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.NVMeSubsystem{Nqn: r.Nqn}
	}
	return &pb.NVMeSubsystemListResponse{Subsystem: Blobarray}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("NVMeSubsystemGet: Received from client: %v", in)
	subsys, ok := subsystems[in.SubsystemId.Value]
	if !ok {
		err := fmt.Errorf("unable to find key %s", in.SubsystemId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}

	var result []NvmfGetSubsystemsResult
	err := call("nvmf_get_subsystems", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)

	for i := range result {
		r := &result[i]
		if r.Nqn == subsys.Nqn {
			return &pb.NVMeSubsystemGetResponse{Subsystem: &pb.NVMeSubsystem{Nqn: r.Nqn}}, nil
		}
	}
	msg := fmt.Sprintf("Could not find NQN: %s", subsys.Nqn)
	log.Print(msg)
	return nil, status.Errorf(codes.InvalidArgument, msg)
}

func (s *server) NVMeSubsystemStats(ctx context.Context, in *pb.NVMeSubsystemStatsRequest) (*pb.NVMeSubsystemStatsResponse, error) {
	log.Printf("NVMeSubsystemStats: Received from client: %v", in)
	var result NvmfGetSubsystemStatsResult
	err := call("nvmf_get_stats", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	return &pb.NVMeSubsystemStatsResponse{Stats: fmt.Sprint(result.TickRate)}, nil
}

// ////////////////////////////////////////////////////////
var controllers = map[string]*pb.NVMeController{}

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received from client: %v", in.Controller)
	controllers[in.Controller.Id.Value] = in.Controller
	return &pb.NVMeControllerCreateResponse{}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received from client: %v", in.ControllerId)
	controller, ok := controllers[in.ControllerId.Value]
	if !ok {
		return nil, fmt.Errorf("error finding controller %s", in.ControllerId.Value)
	}
	delete(controllers, controller.Id.Value)
	return &pb.NVMeControllerDeleteResponse{}, nil
}

func (s *server) NVMeControllerUpdate(ctx context.Context, in *pb.NVMeControllerUpdateRequest) (*pb.NVMeControllerUpdateResponse, error) {
	log.Printf("Received from client: %v", in.Controller)
	controllers[in.Controller.Id.Value] = in.Controller
	return &pb.NVMeControllerUpdateResponse{}, nil
}

func (s *server) NVMeControllerList(ctx context.Context, in *pb.NVMeControllerListRequest) (*pb.NVMeControllerListResponse, error) {
	log.Printf("Received from client: %v", in.SubsystemId)
	Blobarray := []*pb.NVMeController{}
	for _, controller := range controllers {
		Blobarray = append(Blobarray, controller)
	}
	return &pb.NVMeControllerListResponse{Controller: Blobarray}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received from client: %v", in.ControllerId)
	controller, ok := controllers[in.ControllerId.Value]
	if !ok {
		return nil, fmt.Errorf("error finding controller %s", in.ControllerId.Value)
	}
	return &pb.NVMeControllerGetResponse{Controller: controller}, nil
}

func (s *server) NVMeControllerStats(ctx context.Context, in *pb.NVMeControllerStatsRequest) (*pb.NVMeControllerStatsResponse, error) {
	log.Printf("Received from client: %v", in.Id)
	return &pb.NVMeControllerStatsResponse{}, nil
}

// ////////////////////////////////////////////////////////
var namespaces = map[string]*pb.NVMeNamespace{}

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("NVMeNamespaceCreate: Received from client: %v", in)
	subsys, ok := subsystems[in.Namespace.SubsystemId.Value]
	if !ok {
		err := fmt.Errorf("unable to find subsystem %s", in.Namespace.SubsystemId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}

	params := NvmfSubsystemAddNsParams{
		Nqn: subsys.Nqn,
	}

	// TODO: using bdev for crypto key id as a middle end handle for now
	params.Namespace.BdevName = in.Namespace.CryptoKeyId.Value

	var result NvmfSubsystemAddNsResult
	err := call("nvmf_subsystem_add_ns", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	namespaces[in.Namespace.Id.Value] = in.Namespace
	return &pb.NVMeNamespaceCreateResponse{}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("NVMeNamespaceDelete: Received from client: %v", in)
	namespace, ok := namespaces[in.NamespaceId.Value]
	if !ok {
		err := fmt.Errorf("unable to find key %s", in.NamespaceId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}
	subsys, ok := subsystems[namespace.SubsystemId.Value]
	if !ok {
		err := fmt.Errorf("unable to find subsystem %s", namespace.SubsystemId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}

	params := NvmfSubsystemRemoveNsParams{
		Nqn:  subsys.Nqn,
		Nsid: int(namespace.HostNsid),
	}
	var result NvmfSubsystemRemoveNsResult
	err := call("nvmf_subsystem_remove_ns", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	delete(namespaces, namespace.Id.Value)
	return &pb.NVMeNamespaceDeleteResponse{}, nil
}

func (s *server) NVMeNamespaceUpdate(ctx context.Context, in *pb.NVMeNamespaceUpdateRequest) (*pb.NVMeNamespaceUpdateResponse, error) {
	log.Printf("Received from client: %v", in.Namespace)
	namespaces[in.Namespace.Id.Value] = in.Namespace
	return &pb.NVMeNamespaceUpdateResponse{}, nil
}

func (s *server) NVMeNamespaceList(ctx context.Context, in *pb.NVMeNamespaceListRequest) (*pb.NVMeNamespaceListResponse, error) {
	log.Printf("NVMeNamespaceList: Received from client: %v", in)

	nqn := ""
	if in.SubsystemId != nil {
		subsys, ok := subsystems[in.SubsystemId.Value]
		if !ok {
			err := fmt.Errorf("unable to find subsystem %s", in.SubsystemId.Value)
			log.Printf("error: %v", err)
			return nil, err
		}
		nqn = subsys.Nqn
	}
	var result []NvmfGetSubsystemsResult
	err := call("nvmf_get_subsystems", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)

	Blobarray := []*pb.NVMeNamespace{}
	for i := range result {
		rr := &result[i]
		if rr.Nqn == nqn || nqn == "" {
			for j := range rr.Namespaces {
				r := &rr.Namespaces[j]
				Blobarray = append(Blobarray, &pb.NVMeNamespace{HostNsid: uint32(r.Nsid)})
			}
		}
	}
	if len(Blobarray) > 0 {
		return &pb.NVMeNamespaceListResponse{Namespace: Blobarray}, nil
	}

	msg := fmt.Sprintf("Could not find any namespaces for NQN: %s", nqn)
	log.Print(msg)
	return nil, status.Errorf(codes.InvalidArgument, msg)
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("NVMeNamespaceGet: Received from client: %v", in)
	namespace, ok := namespaces[in.NamespaceId.Value]
	if !ok {
		err := fmt.Errorf("unable to find key %s", in.NamespaceId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}
	// TODO: do we even query SPDK to confirm if namespace is present?
	// return &pb.NVMeNamespaceGetResponse{Namespace: &namespace}, nil

	// fetch subsystems -> namespaces from server, match the nsid to find the corresponding namespace
	subsys, ok := subsystems[namespace.SubsystemId.Value]
	if !ok {
		err := fmt.Errorf("unable to find subsystem %s", namespace.SubsystemId.Value)
		log.Printf("error: %v", err)
		return nil, err
	}

	var result []NvmfGetSubsystemsResult
	err := call("nvmf_get_subsystems", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	for i := range result {
		rr := &result[i]
		if rr.Nqn == subsys.Nqn {
			for j := range rr.Namespaces {
				r := &rr.Namespaces[j]
				if uint32(r.Nsid) == namespace.HostNsid {
					return &pb.NVMeNamespaceGetResponse{Namespace: namespace}, nil
				}
			}
			msg := fmt.Sprintf("Could not find NSID: %d", namespace.HostNsid)
			log.Print(msg)
			return nil, status.Errorf(codes.InvalidArgument, msg)
		}
	}
	msg := fmt.Sprintf("Could not find NQN: %s", subsys.Nqn)
	log.Print(msg)
	return nil, status.Errorf(codes.InvalidArgument, msg)
}

func (s *server) NVMeNamespaceStats(ctx context.Context, in *pb.NVMeNamespaceStatsRequest) (*pb.NVMeNamespaceStatsResponse, error) {
	log.Printf("Received from client: %v", in.NamespaceId)
	return &pb.NVMeNamespaceStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) VirtioBlkCreate(ctx context.Context, in *pb.VirtioBlkCreateRequest) (*pb.VirtioBlkCreateResponse, error) {
	log.Printf("VirtioBlkCreate: Received from client: %v", in)
	params := VhostCreateBlkControllerParams{
		Ctrlr:   in.GetController().GetName(),
		DevName: in.GetController().GetBdev(),
	}
	var result VhostCreateBlkControllerResult
	err := call("vhost_create_blk_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not create: %v", in)
	}
	return &pb.VirtioBlkCreateResponse{}, nil
}

func (s *server) VirtioBlkDelete(ctx context.Context, in *pb.VirtioBlkDeleteRequest) (*pb.VirtioBlkDeleteResponse, error) {
	log.Printf("VirtioBlkDelete: Received from client: %v", in)
	params := VhostDeleteControllerParams{
		Ctrlr: fmt.Sprint("VirtioBlk", in.GetControllerId()),
	}
	var result VhostDeleteControllerResult
	err := call("vhost_delete_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not delete: %v", in)
	}
	return &pb.VirtioBlkDeleteResponse{}, nil
}

func (s *server) VirtioBlkUpdate(ctx context.Context, in *pb.VirtioBlkUpdateRequest) (*pb.VirtioBlkUpdateResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkUpdateResponse{}, nil
}

func (s *server) VirtioBlkList(ctx context.Context, in *pb.VirtioBlkListRequest) (*pb.VirtioBlkListResponse, error) {
	log.Printf("VirtioBlkList: Received from client: %v", in)
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.VirtioBlk, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.VirtioBlk{Name: r.Ctrlr}
	}
	return &pb.VirtioBlkListResponse{Controller: Blobarray}, nil
}

func (s *server) VirtioBlkGet(ctx context.Context, in *pb.VirtioBlkGetRequest) (*pb.VirtioBlkGetResponse, error) {
	log.Printf("VirtioBlkGet: Received from client: %v", in)
	params := VhostGetControllersParams{
		Name: fmt.Sprint("VirtioBlk", in.GetControllerId()),
	}
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if len(result) != 1 {
		msg := fmt.Sprintf("expecting exactly 1 result, got %d", len(result))
		log.Print(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	return &pb.VirtioBlkGetResponse{Controller: &pb.VirtioBlk{Name: result[0].Ctrlr}}, nil
}

func (s *server) VirtioBlkStats(ctx context.Context, in *pb.VirtioBlkStatsRequest) (*pb.VirtioBlkStatsResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) VirtioScsiControllerCreate(ctx context.Context, in *pb.VirtioScsiControllerCreateRequest) (*pb.VirtioScsiControllerCreateResponse, error) {
	log.Printf("VirtioScsiControllerCreate: Received from client: %v", in)
	params := VhostCreateScsiControllerParams{
		Ctrlr: in.GetController().GetName(),
	}
	var result VhostCreateScsiControllerResult
	err := call("vhost_create_scsi_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not create: %v", in)
	}
	return &pb.VirtioScsiControllerCreateResponse{}, nil
}

func (s *server) VirtioScsiControllerDelete(ctx context.Context, in *pb.VirtioScsiControllerDeleteRequest) (*pb.VirtioScsiControllerDeleteResponse, error) {
	log.Printf("VirtioScsiControllerDelete: Received from client: %v", in)
	params := VhostDeleteControllerParams{
		Ctrlr: fmt.Sprint("OPI-VirtioScsi", in.GetControllerId()),
	}
	var result VhostDeleteControllerResult
	err := call("vhost_delete_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not delete: %v", in)
	}
	return &pb.VirtioScsiControllerDeleteResponse{}, nil
}

func (s *server) VirtioScsiControllerUpdate(ctx context.Context, in *pb.VirtioScsiControllerUpdateRequest) (*pb.VirtioScsiControllerUpdateResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioScsiControllerUpdateResponse{}, nil
}

func (s *server) VirtioScsiControllerList(ctx context.Context, in *pb.VirtioScsiControllerListRequest) (*pb.VirtioScsiControllerListResponse, error) {
	log.Printf("VirtioScsiControllerList: Received from client: %v", in)
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.VirtioScsiController, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.VirtioScsiController{Name: r.Ctrlr}
	}
	return &pb.VirtioScsiControllerListResponse{Controller: Blobarray}, nil
}

func (s *server) VirtioScsiControllerGet(ctx context.Context, in *pb.VirtioScsiControllerGetRequest) (*pb.VirtioScsiControllerGetResponse, error) {
	log.Printf("VirtioScsiControllerGet: Received from client: %v", in)
	params := VhostGetControllersParams{
		Name: fmt.Sprint("OPI-VirtioScsi", in.GetControllerId()),
	}
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if len(result) != 1 {
		msg := fmt.Sprintf("expecting exactly 1 result, got %d", len(result))
		log.Print(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	return &pb.VirtioScsiControllerGetResponse{Controller: &pb.VirtioScsiController{Name: result[0].Ctrlr}}, nil
}

func (s *server) VirtioScsiControllerStats(ctx context.Context, in *pb.VirtioScsiControllerStatsRequest) (*pb.VirtioScsiControllerStatsResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioScsiControllerStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) VirtioScsiLunCreate(ctx context.Context, in *pb.VirtioScsiLunCreateRequest) (*pb.VirtioScsiLunCreateResponse, error) {
	log.Printf("VirtioScsiLunCreate: Received from client: %v", in)
	params := struct {
		Name string `json:"ctrlr"`
		Num  int    `json:"scsi_target_num"`
		Bdev string `json:"bdev_name"`
	}{
		Name: fmt.Sprint("OPI-VirtioScsi", in.GetLun().GetControllerId()),
		Num:  5,
		Bdev: in.GetLun().GetBdev(),
	}
	var result int
	err := call("vhost_scsi_controller_add_target", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	return &pb.VirtioScsiLunCreateResponse{}, nil
}

func (s *server) VirtioScsiLunDelete(ctx context.Context, in *pb.VirtioScsiLunDeleteRequest) (*pb.VirtioScsiLunDeleteResponse, error) {
	log.Printf("VirtioScsiLunDelete: Received from client: %v", in)
	params := struct {
		Name string `json:"ctrlr"`
		Num  int    `json:"scsi_target_num"`
	}{
		Name: fmt.Sprint("OPI-VirtioScsi", in.GetControllerId()),
		Num:  5,
	}
	var result bool
	err := call("vhost_scsi_controller_remove_target", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if !result {
		log.Printf("Could not delete: %v", in)
	}
	return &pb.VirtioScsiLunDeleteResponse{}, nil
}

func (s *server) VirtioScsiLunUpdate(ctx context.Context, in *pb.VirtioScsiLunUpdateRequest) (*pb.VirtioScsiLunUpdateResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioScsiLunUpdateResponse{}, nil
}

func (s *server) VirtioScsiLunList(ctx context.Context, in *pb.VirtioScsiLunListRequest) (*pb.VirtioScsiLunListResponse, error) {
	log.Printf("VirtioScsiLunList: Received from client: %v", in)
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.VirtioScsiLun, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.VirtioScsiLun{Bdev: r.Ctrlr}
	}
	return &pb.VirtioScsiLunListResponse{Lun: Blobarray}, nil
}

func (s *server) VirtioScsiLunGet(ctx context.Context, in *pb.VirtioScsiLunGetRequest) (*pb.VirtioScsiLunGetResponse, error) {
	log.Printf("VirtioScsiLunGet: Received from client: %v", in)
	params := VhostGetControllersParams{
		Name: fmt.Sprint("OPI-VirtioScsi", in.GetControllerId()),
	}
	var result []VhostGetControllersResult
	err := call("vhost_get_controllers", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if len(result) != 1 {
		msg := fmt.Sprintf("expecting exactly 1 result, got %d", len(result))
		log.Print(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	return &pb.VirtioScsiLunGetResponse{Lun: &pb.VirtioScsiLun{Bdev: result[0].Ctrlr}}, nil
}

func (s *server) VirtioScsiLunStats(ctx context.Context, in *pb.VirtioScsiLunStatsRequest) (*pb.VirtioScsiLunStatsResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioScsiLunStatsResponse{}, nil
}
