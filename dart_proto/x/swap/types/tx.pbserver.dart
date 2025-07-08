// This is a generated file - do not edit.
//
// Generated from x/swap/types/tx.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'tx.pb.dart' as $0;
import 'tx.pbjson.dart';

export 'tx.pb.dart';

abstract class MsgServiceBase extends $pb.GeneratedService {
  $async.Future<$0.MsgUpdateParamsResponse> updateParams($pb.ServerContext ctx, $0.MsgUpdateParams request);
  $async.Future<$0.MsgProvideLiquidityResponse> provideLiquidity($pb.ServerContext ctx, $0.MsgProvideLiquidity request);
  $async.Future<$0.MsgWithdrawLiquidityResponse> withdrawLiquidity($pb.ServerContext ctx, $0.MsgWithdrawLiquidity request);
  $async.Future<$0.MsgSwapResponse> swap($pb.ServerContext ctx, $0.MsgSwap request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'UpdateParams': return $0.MsgUpdateParams();
      case 'ProvideLiquidity': return $0.MsgProvideLiquidity();
      case 'WithdrawLiquidity': return $0.MsgWithdrawLiquidity();
      case 'Swap': return $0.MsgSwap();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'UpdateParams': return updateParams(ctx, request as $0.MsgUpdateParams);
      case 'ProvideLiquidity': return provideLiquidity(ctx, request as $0.MsgProvideLiquidity);
      case 'WithdrawLiquidity': return withdrawLiquidity(ctx, request as $0.MsgWithdrawLiquidity);
      case 'Swap': return swap(ctx, request as $0.MsgSwap);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => MsgServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => MsgServiceBase$messageJson;
}

