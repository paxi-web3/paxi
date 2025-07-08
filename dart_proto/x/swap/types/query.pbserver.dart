// This is a generated file - do not edit.
//
// Generated from x/swap/types/query.proto.

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
  $async.Future<$1.QueryParamsResponse> params($pb.ServerContext ctx, $1.QueryParamsRequest request);
  $async.Future<$1.QueryPositionResponse> position($pb.ServerContext ctx, $1.QueryPositionRequest request);
  $async.Future<$1.QueryPoolResponse> pool($pb.ServerContext ctx, $1.QueryPoolRequest request);
  $async.Future<$1.QueryAllPoolsResponse> allPools($pb.ServerContext ctx, $1.QueryAllPoolsRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'Params': return $1.QueryParamsRequest();
      case 'Position': return $1.QueryPositionRequest();
      case 'Pool': return $1.QueryPoolRequest();
      case 'AllPools': return $1.QueryAllPoolsRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'Params': return params(ctx, request as $1.QueryParamsRequest);
      case 'Position': return position(ctx, request as $1.QueryPositionRequest);
      case 'Pool': return pool(ctx, request as $1.QueryPoolRequest);
      case 'AllPools': return allPools(ctx, request as $1.QueryAllPoolsRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => QueryServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => QueryServiceBase$messageJson;
}

