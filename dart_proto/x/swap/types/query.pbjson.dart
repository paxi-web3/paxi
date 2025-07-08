// This is a generated file - do not edit.
//
// Generated from x/swap/types/query.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import '../../../cosmos/base/query/v1beta1/pagination.pbjson.dart' as $0;

@$core.Deprecated('Use queryParamsRequestDescriptor instead')
const QueryParamsRequest$json = {
  '1': 'QueryParamsRequest',
};

/// Descriptor for `QueryParamsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryParamsRequestDescriptor = $convert.base64Decode(
    'ChJRdWVyeVBhcmFtc1JlcXVlc3Q=');

@$core.Deprecated('Use queryParamsResponseDescriptor instead')
const QueryParamsResponse$json = {
  '1': 'QueryParamsResponse',
  '2': [
    {'1': 'code_id', '3': 1, '4': 1, '5': 4, '10': 'codeId'},
    {'1': 'swap_fee_bps', '3': 2, '4': 1, '5': 4, '10': 'swapFeeBps'},
    {'1': 'min_liquidity', '3': 3, '4': 1, '5': 4, '10': 'minLiquidity'},
  ],
};

/// Descriptor for `QueryParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryParamsResponseDescriptor = $convert.base64Decode(
    'ChNRdWVyeVBhcmFtc1Jlc3BvbnNlEhcKB2NvZGVfaWQYASABKARSBmNvZGVJZBIgCgxzd2FwX2'
    'ZlZV9icHMYAiABKARSCnN3YXBGZWVCcHMSIwoNbWluX2xpcXVpZGl0eRgDIAEoBFIMbWluTGlx'
    'dWlkaXR5');

@$core.Deprecated('Use providerPositionDescriptor instead')
const ProviderPosition$json = {
  '1': 'ProviderPosition',
  '2': [
    {'1': 'creator', '3': 1, '4': 1, '5': 9, '10': 'creator'},
    {'1': 'prc20', '3': 2, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'lp_amount', '3': 3, '4': 1, '5': 9, '10': 'lpAmount'},
  ],
};

/// Descriptor for `ProviderPosition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List providerPositionDescriptor = $convert.base64Decode(
    'ChBQcm92aWRlclBvc2l0aW9uEhgKB2NyZWF0b3IYASABKAlSB2NyZWF0b3ISFAoFcHJjMjAYAi'
    'ABKAlSBXByYzIwEhsKCWxwX2Ftb3VudBgDIAEoCVIIbHBBbW91bnQ=');

@$core.Deprecated('Use queryPositionRequestDescriptor instead')
const QueryPositionRequest$json = {
  '1': 'QueryPositionRequest',
  '2': [
    {'1': 'creator', '3': 1, '4': 1, '5': 9, '10': 'creator'},
    {'1': 'prc20', '3': 2, '4': 1, '5': 9, '10': 'prc20'},
  ],
};

/// Descriptor for `QueryPositionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryPositionRequestDescriptor = $convert.base64Decode(
    'ChRRdWVyeVBvc2l0aW9uUmVxdWVzdBIYCgdjcmVhdG9yGAEgASgJUgdjcmVhdG9yEhQKBXByYz'
    'IwGAIgASgJUgVwcmMyMA==');

@$core.Deprecated('Use queryPositionResponseDescriptor instead')
const QueryPositionResponse$json = {
  '1': 'QueryPositionResponse',
  '2': [
    {'1': 'position', '3': 1, '4': 1, '5': 11, '6': '.x.swap.types.ProviderPosition', '10': 'position'},
    {'1': 'expected_paxi', '3': 2, '4': 1, '5': 9, '10': 'expectedPaxi'},
    {'1': 'expected_prc20', '3': 3, '4': 1, '5': 9, '10': 'expectedPrc20'},
  ],
};

/// Descriptor for `QueryPositionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryPositionResponseDescriptor = $convert.base64Decode(
    'ChVRdWVyeVBvc2l0aW9uUmVzcG9uc2USOgoIcG9zaXRpb24YASABKAsyHi54LnN3YXAudHlwZX'
    'MuUHJvdmlkZXJQb3NpdGlvblIIcG9zaXRpb24SIwoNZXhwZWN0ZWRfcGF4aRgCIAEoCVIMZXhw'
    'ZWN0ZWRQYXhpEiUKDmV4cGVjdGVkX3ByYzIwGAMgASgJUg1leHBlY3RlZFByYzIw');

@$core.Deprecated('Use queryPoolRequestDescriptor instead')
const QueryPoolRequest$json = {
  '1': 'QueryPoolRequest',
  '2': [
    {'1': 'prc20', '3': 1, '4': 1, '5': 9, '10': 'prc20'},
  ],
};

/// Descriptor for `QueryPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryPoolRequestDescriptor = $convert.base64Decode(
    'ChBRdWVyeVBvb2xSZXF1ZXN0EhQKBXByYzIwGAEgASgJUgVwcmMyMA==');

@$core.Deprecated('Use queryPoolResponseDescriptor instead')
const QueryPoolResponse$json = {
  '1': 'QueryPoolResponse',
  '2': [
    {'1': 'prc20', '3': 1, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'reserve_paxi', '3': 2, '4': 1, '5': 9, '10': 'reservePaxi'},
    {'1': 'reserve_prc20', '3': 3, '4': 1, '5': 9, '10': 'reservePrc20'},
    {'1': 'price_paxi_per_prc20', '3': 4, '4': 1, '5': 9, '10': 'pricePaxiPerPrc20'},
    {'1': 'price_prc20_per_paxi', '3': 5, '4': 1, '5': 9, '10': 'pricePrc20PerPaxi'},
    {'1': 'total_shares', '3': 6, '4': 1, '5': 9, '10': 'totalShares'},
  ],
};

/// Descriptor for `QueryPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryPoolResponseDescriptor = $convert.base64Decode(
    'ChFRdWVyeVBvb2xSZXNwb25zZRIUCgVwcmMyMBgBIAEoCVIFcHJjMjASIQoMcmVzZXJ2ZV9wYX'
    'hpGAIgASgJUgtyZXNlcnZlUGF4aRIjCg1yZXNlcnZlX3ByYzIwGAMgASgJUgxyZXNlcnZlUHJj'
    'MjASLwoUcHJpY2VfcGF4aV9wZXJfcHJjMjAYBCABKAlSEXByaWNlUGF4aVBlclByYzIwEi8KFH'
    'ByaWNlX3ByYzIwX3Blcl9wYXhpGAUgASgJUhFwcmljZVByYzIwUGVyUGF4aRIhCgx0b3RhbF9z'
    'aGFyZXMYBiABKAlSC3RvdGFsU2hhcmVz');

@$core.Deprecated('Use queryAllPoolsRequestDescriptor instead')
const QueryAllPoolsRequest$json = {
  '1': 'QueryAllPoolsRequest',
  '2': [
    {'1': 'pagination', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.query.v1beta1.PageRequest', '10': 'pagination'},
  ],
};

/// Descriptor for `QueryAllPoolsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryAllPoolsRequestDescriptor = $convert.base64Decode(
    'ChRRdWVyeUFsbFBvb2xzUmVxdWVzdBJGCgpwYWdpbmF0aW9uGAEgASgLMiYuY29zbW9zLmJhc2'
    'UucXVlcnkudjFiZXRhMS5QYWdlUmVxdWVzdFIKcGFnaW5hdGlvbg==');

@$core.Deprecated('Use queryAllPoolsResponseDescriptor instead')
const QueryAllPoolsResponse$json = {
  '1': 'QueryAllPoolsResponse',
  '2': [
    {'1': 'pools', '3': 1, '4': 3, '5': 11, '6': '.x.swap.types.QueryPoolResponse', '10': 'pools'},
    {'1': 'pagination', '3': 2, '4': 1, '5': 11, '6': '.cosmos.base.query.v1beta1.PageResponse', '10': 'pagination'},
  ],
};

/// Descriptor for `QueryAllPoolsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryAllPoolsResponseDescriptor = $convert.base64Decode(
    'ChVRdWVyeUFsbFBvb2xzUmVzcG9uc2USNQoFcG9vbHMYASADKAsyHy54LnN3YXAudHlwZXMuUX'
    'VlcnlQb29sUmVzcG9uc2VSBXBvb2xzEkcKCnBhZ2luYXRpb24YAiABKAsyJy5jb3Ntb3MuYmFz'
    'ZS5xdWVyeS52MWJldGExLlBhZ2VSZXNwb25zZVIKcGFnaW5hdGlvbg==');

@$core.Deprecated('Use poolProtoDescriptor instead')
const PoolProto$json = {
  '1': 'PoolProto',
  '2': [
    {'1': 'prc20', '3': 1, '4': 1, '5': 9, '10': 'prc20'},
    {'1': 'reserve_paxi', '3': 2, '4': 1, '5': 9, '10': 'reservePaxi'},
    {'1': 'reserve_prc20', '3': 3, '4': 1, '5': 9, '10': 'reservePrc20'},
    {'1': 'total_shares', '3': 4, '4': 1, '5': 9, '10': 'totalShares'},
  ],
};

/// Descriptor for `PoolProto`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List poolProtoDescriptor = $convert.base64Decode(
    'CglQb29sUHJvdG8SFAoFcHJjMjAYASABKAlSBXByYzIwEiEKDHJlc2VydmVfcGF4aRgCIAEoCV'
    'ILcmVzZXJ2ZVBheGkSIwoNcmVzZXJ2ZV9wcmMyMBgDIAEoCVIMcmVzZXJ2ZVByYzIwEiEKDHRv'
    'dGFsX3NoYXJlcxgEIAEoCVILdG90YWxTaGFyZXM=');

const $core.Map<$core.String, $core.dynamic> QueryServiceBase$json = {
  '1': 'Query',
  '2': [
    {'1': 'Params', '2': '.x.swap.types.QueryParamsRequest', '3': '.x.swap.types.QueryParamsResponse', '4': {}},
    {'1': 'Position', '2': '.x.swap.types.QueryPositionRequest', '3': '.x.swap.types.QueryPositionResponse', '4': {}},
    {'1': 'Pool', '2': '.x.swap.types.QueryPoolRequest', '3': '.x.swap.types.QueryPoolResponse', '4': {}},
    {'1': 'AllPools', '2': '.x.swap.types.QueryAllPoolsRequest', '3': '.x.swap.types.QueryAllPoolsResponse', '4': {}},
  ],
};

@$core.Deprecated('Use queryServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> QueryServiceBase$messageJson = {
  '.x.swap.types.QueryParamsRequest': QueryParamsRequest$json,
  '.x.swap.types.QueryParamsResponse': QueryParamsResponse$json,
  '.x.swap.types.QueryPositionRequest': QueryPositionRequest$json,
  '.x.swap.types.QueryPositionResponse': QueryPositionResponse$json,
  '.x.swap.types.ProviderPosition': ProviderPosition$json,
  '.x.swap.types.QueryPoolRequest': QueryPoolRequest$json,
  '.x.swap.types.QueryPoolResponse': QueryPoolResponse$json,
  '.x.swap.types.QueryAllPoolsRequest': QueryAllPoolsRequest$json,
  '.cosmos.base.query.v1beta1.PageRequest': $0.PageRequest$json,
  '.x.swap.types.QueryAllPoolsResponse': QueryAllPoolsResponse$json,
  '.cosmos.base.query.v1beta1.PageResponse': $0.PageResponse$json,
};

/// Descriptor for `Query`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List queryServiceDescriptor = $convert.base64Decode(
    'CgVRdWVyeRJoCgZQYXJhbXMSIC54LnN3YXAudHlwZXMuUXVlcnlQYXJhbXNSZXF1ZXN0GiEueC'
    '5zd2FwLnR5cGVzLlF1ZXJ5UGFyYW1zUmVzcG9uc2UiGYLT5JMCExIRL3BheGkvc3dhcC9wYXJh'
    'bXMSggEKCFBvc2l0aW9uEiIueC5zd2FwLnR5cGVzLlF1ZXJ5UG9zaXRpb25SZXF1ZXN0GiMueC'
    '5zd2FwLnR5cGVzLlF1ZXJ5UG9zaXRpb25SZXNwb25zZSItgtPkkwInEiUvcGF4aS9zd2FwL3Bv'
    'c2l0aW9uL3tjcmVhdG9yfS97cHJjMjB9EmgKBFBvb2wSHi54LnN3YXAudHlwZXMuUXVlcnlQb2'
    '9sUmVxdWVzdBofLnguc3dhcC50eXBlcy5RdWVyeVBvb2xSZXNwb25zZSIfgtPkkwIZEhcvcGF4'
    'aS9zd2FwL3Bvb2wve3ByYzIwfRJxCghBbGxQb29scxIiLnguc3dhcC50eXBlcy5RdWVyeUFsbF'
    'Bvb2xzUmVxdWVzdBojLnguc3dhcC50eXBlcy5RdWVyeUFsbFBvb2xzUmVzcG9uc2UiHILT5JMC'
    'FhIUL3BheGkvc3dhcC9hbGxfcG9vbHM=');

