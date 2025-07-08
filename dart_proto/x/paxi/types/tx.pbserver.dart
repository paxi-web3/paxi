// This is a generated file - do not edit.
//
// Generated from x/paxi/types/tx.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'tx.pb.dart' as $1;
import 'tx.pbjson.dart';

export 'tx.pb.dart';

abstract class MsgServiceBase extends $pb.GeneratedService {
  $async.Future<$1.MsgBurnTokenResponse> burnToken($pb.ServerContext ctx, $1.MsgBurnToken request);
  $async.Future<$1.MsgUpdateParamsResponse> updateParams($pb.ServerContext ctx, $1.MsgUpdateParams request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'BurnToken': return $1.MsgBurnToken();
      case 'UpdateParams': return $1.MsgUpdateParams();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'BurnToken': return burnToken(ctx, request as $1.MsgBurnToken);
      case 'UpdateParams': return updateParams(ctx, request as $1.MsgUpdateParams);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => MsgServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => MsgServiceBase$messageJson;
}

