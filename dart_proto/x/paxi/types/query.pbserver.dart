// This is a generated file - do not edit.
//
// Generated from x/paxi/types/query.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'query.pb.dart' as $1;
import 'query.pbjson.dart';

export 'query.pb.dart';

abstract class QueryServiceBase extends $pb.GeneratedService {
  $async.Future<$1.QueryLockedVestingResponse> lockedVesting($pb.ServerContext ctx, $1.QueryLockedVestingRequest request);
  $async.Future<$1.QueryCirculatingSupplyResponse> circulatingSupply($pb.ServerContext ctx, $1.QueryCirculatingSupplyRequest request);
  $async.Future<$1.QueryTotalSupplyResponse> totalSupply($pb.ServerContext ctx, $1.QueryTotalSupplyRequest request);
  $async.Future<$1.QueryEstimatedGasPriceResponse> estimatedGasPrice($pb.ServerContext ctx, $1.QueryEstimatedGasPriceRequest request);
  $async.Future<$1.QueryLastBlockGasUsedResponse> lastBlockGasUsed($pb.ServerContext ctx, $1.QueryLastBlockGasUsedRequest request);
  $async.Future<$1.QueryTotalTxsResponse> totalTxs($pb.ServerContext ctx, $1.QueryTotalTxsRequest request);
  $async.Future<$1.QueryUnlockSchedulesResponse> unlockSchedules($pb.ServerContext ctx, $1.QueryUnlockSchedulesRequest request);
  $async.Future<$1.QueryParamsResponse> params($pb.ServerContext ctx, $1.QueryParamsRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'LockedVesting': return $1.QueryLockedVestingRequest();
      case 'CirculatingSupply': return $1.QueryCirculatingSupplyRequest();
      case 'TotalSupply': return $1.QueryTotalSupplyRequest();
      case 'EstimatedGasPrice': return $1.QueryEstimatedGasPriceRequest();
      case 'LastBlockGasUsed': return $1.QueryLastBlockGasUsedRequest();
      case 'TotalTxs': return $1.QueryTotalTxsRequest();
      case 'UnlockSchedules': return $1.QueryUnlockSchedulesRequest();
      case 'Params': return $1.QueryParamsRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'LockedVesting': return lockedVesting(ctx, request as $1.QueryLockedVestingRequest);
      case 'CirculatingSupply': return circulatingSupply(ctx, request as $1.QueryCirculatingSupplyRequest);
      case 'TotalSupply': return totalSupply(ctx, request as $1.QueryTotalSupplyRequest);
      case 'EstimatedGasPrice': return estimatedGasPrice(ctx, request as $1.QueryEstimatedGasPriceRequest);
      case 'LastBlockGasUsed': return lastBlockGasUsed(ctx, request as $1.QueryLastBlockGasUsedRequest);
      case 'TotalTxs': return totalTxs(ctx, request as $1.QueryTotalTxsRequest);
      case 'UnlockSchedules': return unlockSchedules(ctx, request as $1.QueryUnlockSchedulesRequest);
      case 'Params': return params(ctx, request as $1.QueryParamsRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => QueryServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => QueryServiceBase$messageJson;
}

