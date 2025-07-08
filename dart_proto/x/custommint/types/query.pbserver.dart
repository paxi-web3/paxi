// This is a generated file - do not edit.
//
// Generated from x/custommint/types/query.proto.

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
  $async.Future<$1.QueryTotalMintedResponse> totalMinted($pb.ServerContext ctx, $1.QueryTotalMintedRequest request);
  $async.Future<$1.QueryTotalBurnedResponse> totalBurned($pb.ServerContext ctx, $1.QueryTotalBurnedRequest request);
  $async.Future<$1.QueryParamsResponse> params($pb.ServerContext ctx, $1.QueryParamsRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'TotalMinted': return $1.QueryTotalMintedRequest();
      case 'TotalBurned': return $1.QueryTotalBurnedRequest();
      case 'Params': return $1.QueryParamsRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'TotalMinted': return totalMinted(ctx, request as $1.QueryTotalMintedRequest);
      case 'TotalBurned': return totalBurned(ctx, request as $1.QueryTotalBurnedRequest);
      case 'Params': return params(ctx, request as $1.QueryParamsRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => QueryServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => QueryServiceBase$messageJson;
}

