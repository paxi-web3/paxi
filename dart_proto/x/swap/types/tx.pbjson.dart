// This is a generated file - do not edit.
//
// Generated from x/swap/types/tx.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use paramsInputDescriptor instead')
const ParamsInput$json = {
  '1': 'ParamsInput',
  '2': [
    {'1': 'code_id', '3': 1, '4': 1, '5': 4, '10': 'codeId'},
    {'1': 'swap_fee_bps', '3': 2, '4': 1, '5': 4, '10': 'swapFeeBps'},
    {'1': 'min_liquidity', '3': 3, '4': 1, '5': 4, '10': 'minLiquidity'},
  ],
};

/// Descriptor for `ParamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List paramsInputDescriptor = $convert.base64Decode(
    'CgtQYXJhbXNJbnB1dBIXCgdjb2RlX2lkGAEgASgEUgZjb2RlSWQSIAoMc3dhcF9mZWVfYnBzGA'
    'IgASgEUgpzd2FwRmVlQnBzEiMKDW1pbl9saXF1aWRpdHkYAyABKARSDG1pbkxpcXVpZGl0eQ==');

@$core.Deprecated('Use msgUpdateParamsDescriptor instead')
const MsgUpdateParams$json = {
  '1': 'MsgUpdateParams',
  '2': [
    {'1': 'authority', '3': 1, '4': 1, '5': 9, '10': 'authority'},
    {'1': 'params', '3': 2, '4': 1, '5': 11, '6': '.x.swap.types.ParamsInput', '8': {}, '10': 'params'},
  ],
  '7': {},
};

/// Descriptor for `MsgUpdateParams`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgUpdateParamsDescriptor = $convert.base64Decode(
    'Cg9Nc2dVcGRhdGVQYXJhbXMSHAoJYXV0aG9yaXR5GAEgASgJUglhdXRob3JpdHkSNwoGcGFyYW'
    '1zGAIgASgLMhkueC5zd2FwLnR5cGVzLlBhcmFtc0lucHV0QgTI3h8AUgZwYXJhbXM6DoLnsCoJ'
    'YXV0aG9yaXR5');

@$core.Deprecated('Use msgUpdateParamsResponseDescriptor instead')
const MsgUpdateParamsResponse$json = {
  '1': 'MsgUpdateParamsResponse',
};

/// Descriptor for `MsgUpdateParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgUpdateParamsResponseDescriptor = $convert.base64Decode(
    'ChdNc2dVcGRhdGVQYXJhbXNSZXNwb25zZQ==');

@$core.Deprecated('Use msgProvideLiquidityDescriptor instead')
const MsgProvideLiquidity$json = {
  '1': 'MsgProvideLiquidity',
  '2': [
    {'1': 'creator', '3': 1, '4': 1, '5': 9, '10': 'creator'},
    {'1': 'prc20', '3': 2, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'paxi_amount', '3': 3, '4': 1, '5': 9, '10': 'paxiAmount'},
    {'1': 'prc20_amount', '3': 4, '4': 1, '5': 9, '10': 'prc20Amount'},
  ],
  '7': {},
};

/// Descriptor for `MsgProvideLiquidity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgProvideLiquidityDescriptor = $convert.base64Decode(
    'ChNNc2dQcm92aWRlTGlxdWlkaXR5EhgKB2NyZWF0b3IYASABKAlSB2NyZWF0b3ISFAoFcHJjMj'
    'AYAiABKAlSBXByYzIwEh8KC3BheGlfYW1vdW50GAMgASgJUgpwYXhpQW1vdW50EiEKDHByYzIw'
    'X2Ftb3VudBgEIAEoCVILcHJjMjBBbW91bnQ6DILnsCoHY3JlYXRvcg==');

@$core.Deprecated('Use msgProvideLiquidityResponseDescriptor instead')
const MsgProvideLiquidityResponse$json = {
  '1': 'MsgProvideLiquidityResponse',
};

/// Descriptor for `MsgProvideLiquidityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgProvideLiquidityResponseDescriptor = $convert.base64Decode(
    'ChtNc2dQcm92aWRlTGlxdWlkaXR5UmVzcG9uc2U=');

@$core.Deprecated('Use msgWithdrawLiquidityDescriptor instead')
const MsgWithdrawLiquidity$json = {
  '1': 'MsgWithdrawLiquidity',
  '2': [
    {'1': 'creator', '3': 1, '4': 1, '5': 9, '10': 'creator'},
    {'1': 'prc20', '3': 2, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'lp_amount', '3': 3, '4': 1, '5': 9, '10': 'lpAmount'},
  ],
  '7': {},
};

/// Descriptor for `MsgWithdrawLiquidity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgWithdrawLiquidityDescriptor = $convert.base64Decode(
    'ChRNc2dXaXRoZHJhd0xpcXVpZGl0eRIYCgdjcmVhdG9yGAEgASgJUgdjcmVhdG9yEhQKBXByYz'
    'IwGAIgASgJUgVwcmMyMBIbCglscF9hbW91bnQYAyABKAlSCGxwQW1vdW50OgyC57AqB2NyZWF0'
    'b3I=');

@$core.Deprecated('Use msgWithdrawLiquidityResponseDescriptor instead')
const MsgWithdrawLiquidityResponse$json = {
  '1': 'MsgWithdrawLiquidityResponse',
};

/// Descriptor for `MsgWithdrawLiquidityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgWithdrawLiquidityResponseDescriptor = $convert.base64Decode(
    'ChxNc2dXaXRoZHJhd0xpcXVpZGl0eVJlc3BvbnNl');

@$core.Deprecated('Use msgSwapDescriptor instead')
const MsgSwap$json = {
  '1': 'MsgSwap',
  '2': [
    {'1': 'creator', '3': 1, '4': 1, '5': 9, '10': 'creator'},
    {'1': 'prc20', '3': 2, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'offer_denom', '3': 3, '4': 1, '5': 9, '10': 'offerDenom'},
    {'1': 'offer_amount', '3': 4, '4': 1, '5': 9, '10': 'offerAmount'},
    {'1': 'min_receive', '3': 5, '4': 1, '5': 9, '10': 'minReceive'},
  ],
  '7': {},
};

/// Descriptor for `MsgSwap`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgSwapDescriptor = $convert.base64Decode(
    'CgdNc2dTd2FwEhgKB2NyZWF0b3IYASABKAlSB2NyZWF0b3ISFAoFcHJjMjAYAiABKAlSBXByYz'
    'IwEh8KC29mZmVyX2Rlbm9tGAMgASgJUgpvZmZlckRlbm9tEiEKDG9mZmVyX2Ftb3VudBgEIAEo'
    'CVILb2ZmZXJBbW91bnQSHwoLbWluX3JlY2VpdmUYBSABKAlSCm1pblJlY2VpdmU6DILnsCoHY3'
    'JlYXRvcg==');

@$core.Deprecated('Use msgSwapResponseDescriptor instead')
const MsgSwapResponse$json = {
  '1': 'MsgSwapResponse',
};

/// Descriptor for `MsgSwapResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgSwapResponseDescriptor = $convert.base64Decode(
    'Cg9Nc2dTd2FwUmVzcG9uc2U=');

const $core.Map<$core.String, $core.dynamic> MsgServiceBase$json = {
  '1': 'Msg',
  '2': [
    {'1': 'UpdateParams', '2': '.x.swap.types.MsgUpdateParams', '3': '.x.swap.types.MsgUpdateParamsResponse'},
    {'1': 'ProvideLiquidity', '2': '.x.swap.types.MsgProvideLiquidity', '3': '.x.swap.types.MsgProvideLiquidityResponse', '4': {}},
    {'1': 'WithdrawLiquidity', '2': '.x.swap.types.MsgWithdrawLiquidity', '3': '.x.swap.types.MsgWithdrawLiquidityResponse', '4': {}},
    {'1': 'Swap', '2': '.x.swap.types.MsgSwap', '3': '.x.swap.types.MsgSwapResponse', '4': {}},
  ],
  '3': {},
};

@$core.Deprecated('Use msgServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> MsgServiceBase$messageJson = {
  '.x.swap.types.MsgUpdateParams': MsgUpdateParams$json,
  '.x.swap.types.ParamsInput': ParamsInput$json,
  '.x.swap.types.MsgUpdateParamsResponse': MsgUpdateParamsResponse$json,
  '.x.swap.types.MsgProvideLiquidity': MsgProvideLiquidity$json,
  '.x.swap.types.MsgProvideLiquidityResponse': MsgProvideLiquidityResponse$json,
  '.x.swap.types.MsgWithdrawLiquidity': MsgWithdrawLiquidity$json,
  '.x.swap.types.MsgWithdrawLiquidityResponse': MsgWithdrawLiquidityResponse$json,
  '.x.swap.types.MsgSwap': MsgSwap$json,
  '.x.swap.types.MsgSwapResponse': MsgSwapResponse$json,
};

/// Descriptor for `Msg`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List msgServiceDescriptor = $convert.base64Decode(
    'CgNNc2cSVAoMVXBkYXRlUGFyYW1zEh0ueC5zd2FwLnR5cGVzLk1zZ1VwZGF0ZVBhcmFtcxolLn'
    'guc3dhcC50eXBlcy5Nc2dVcGRhdGVQYXJhbXNSZXNwb25zZRKHAQoQUHJvdmlkZUxpcXVpZGl0'
    'eRIhLnguc3dhcC50eXBlcy5Nc2dQcm92aWRlTGlxdWlkaXR5GikueC5zd2FwLnR5cGVzLk1zZ1'
    'Byb3ZpZGVMaXF1aWRpdHlSZXNwb25zZSIlgtPkkwIfIhovdHgvc3dhcC9wcm92aWRlX2xpcXVp'
    'ZGl0eToBKhKLAQoRV2l0aGRyYXdMaXF1aWRpdHkSIi54LnN3YXAudHlwZXMuTXNnV2l0aGRyYX'
    'dMaXF1aWRpdHkaKi54LnN3YXAudHlwZXMuTXNnV2l0aGRyYXdMaXF1aWRpdHlSZXNwb25zZSIm'
    'gtPkkwIgIhsvdHgvc3dhcC93aXRoZHJhd19saXF1aWRpdHk6ASoSVgoEU3dhcBIVLnguc3dhcC'
    '50eXBlcy5Nc2dTd2FwGh0ueC5zd2FwLnR5cGVzLk1zZ1N3YXBSZXNwb25zZSIYgtPkkwISIg0v'
    'dHgvc3dhcC9zd2FwOgEqGgWA57AqAQ==');

