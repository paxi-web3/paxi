// This is a generated file - do not edit.
//
// Generated from x/paxi/types/query.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import '../../../cosmos/base/v1beta1/coin.pbjson.dart' as $0;

@$core.Deprecated('Use queryLockedVestingRequestDescriptor instead')
const QueryLockedVestingRequest$json = {
  '1': 'QueryLockedVestingRequest',
};

/// Descriptor for `QueryLockedVestingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLockedVestingRequestDescriptor = $convert.base64Decode(
    'ChlRdWVyeUxvY2tlZFZlc3RpbmdSZXF1ZXN0');

@$core.Deprecated('Use queryLockedVestingResponseDescriptor instead')
const QueryLockedVestingResponse$json = {
  '1': 'QueryLockedVestingResponse',
  '2': [
    {'1': 'amount', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
};

/// Descriptor for `QueryLockedVestingResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLockedVestingResponseDescriptor = $convert.base64Decode(
    'ChpRdWVyeUxvY2tlZFZlc3RpbmdSZXNwb25zZRIxCgZhbW91bnQYASABKAsyGS5jb3Ntb3MuYm'
    'FzZS52MWJldGExLkNvaW5SBmFtb3VudA==');

@$core.Deprecated('Use queryCirculatingSupplyRequestDescriptor instead')
const QueryCirculatingSupplyRequest$json = {
  '1': 'QueryCirculatingSupplyRequest',
};

/// Descriptor for `QueryCirculatingSupplyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryCirculatingSupplyRequestDescriptor = $convert.base64Decode(
    'Ch1RdWVyeUNpcmN1bGF0aW5nU3VwcGx5UmVxdWVzdA==');

@$core.Deprecated('Use queryCirculatingSupplyResponseDescriptor instead')
const QueryCirculatingSupplyResponse$json = {
  '1': 'QueryCirculatingSupplyResponse',
  '2': [
    {'1': 'amount', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
};

/// Descriptor for `QueryCirculatingSupplyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryCirculatingSupplyResponseDescriptor = $convert.base64Decode(
    'Ch5RdWVyeUNpcmN1bGF0aW5nU3VwcGx5UmVzcG9uc2USMQoGYW1vdW50GAEgASgLMhkuY29zbW'
    '9zLmJhc2UudjFiZXRhMS5Db2luUgZhbW91bnQ=');

@$core.Deprecated('Use queryTotalSupplyRequestDescriptor instead')
const QueryTotalSupplyRequest$json = {
  '1': 'QueryTotalSupplyRequest',
};

/// Descriptor for `QueryTotalSupplyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalSupplyRequestDescriptor = $convert.base64Decode(
    'ChdRdWVyeVRvdGFsU3VwcGx5UmVxdWVzdA==');

@$core.Deprecated('Use queryTotalSupplyResponseDescriptor instead')
const QueryTotalSupplyResponse$json = {
  '1': 'QueryTotalSupplyResponse',
  '2': [
    {'1': 'amount', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
};

/// Descriptor for `QueryTotalSupplyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalSupplyResponseDescriptor = $convert.base64Decode(
    'ChhRdWVyeVRvdGFsU3VwcGx5UmVzcG9uc2USMQoGYW1vdW50GAEgASgLMhkuY29zbW9zLmJhc2'
    'UudjFiZXRhMS5Db2luUgZhbW91bnQ=');

@$core.Deprecated('Use queryEstimatedGasPriceRequestDescriptor instead')
const QueryEstimatedGasPriceRequest$json = {
  '1': 'QueryEstimatedGasPriceRequest',
};

/// Descriptor for `QueryEstimatedGasPriceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryEstimatedGasPriceRequestDescriptor = $convert.base64Decode(
    'Ch1RdWVyeUVzdGltYXRlZEdhc1ByaWNlUmVxdWVzdA==');

@$core.Deprecated('Use queryEstimatedGasPriceResponseDescriptor instead')
const QueryEstimatedGasPriceResponse$json = {
  '1': 'QueryEstimatedGasPriceResponse',
  '2': [
    {'1': 'gas_price', '3': 1, '4': 1, '5': 9, '10': 'gasPrice'},
  ],
};

/// Descriptor for `QueryEstimatedGasPriceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryEstimatedGasPriceResponseDescriptor = $convert.base64Decode(
    'Ch5RdWVyeUVzdGltYXRlZEdhc1ByaWNlUmVzcG9uc2USGwoJZ2FzX3ByaWNlGAEgASgJUghnYX'
    'NQcmljZQ==');

@$core.Deprecated('Use queryLastBlockGasUsedRequestDescriptor instead')
const QueryLastBlockGasUsedRequest$json = {
  '1': 'QueryLastBlockGasUsedRequest',
};

/// Descriptor for `QueryLastBlockGasUsedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLastBlockGasUsedRequestDescriptor = $convert.base64Decode(
    'ChxRdWVyeUxhc3RCbG9ja0dhc1VzZWRSZXF1ZXN0');

@$core.Deprecated('Use queryLastBlockGasUsedResponseDescriptor instead')
const QueryLastBlockGasUsedResponse$json = {
  '1': 'QueryLastBlockGasUsedResponse',
  '2': [
    {'1': 'gas_used', '3': 1, '4': 1, '5': 4, '10': 'gasUsed'},
  ],
};

/// Descriptor for `QueryLastBlockGasUsedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryLastBlockGasUsedResponseDescriptor = $convert.base64Decode(
    'Ch1RdWVyeUxhc3RCbG9ja0dhc1VzZWRSZXNwb25zZRIZCghnYXNfdXNlZBgBIAEoBFIHZ2FzVX'
    'NlZA==');

@$core.Deprecated('Use queryTotalTxsRequestDescriptor instead')
const QueryTotalTxsRequest$json = {
  '1': 'QueryTotalTxsRequest',
};

/// Descriptor for `QueryTotalTxsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalTxsRequestDescriptor = $convert.base64Decode(
    'ChRRdWVyeVRvdGFsVHhzUmVxdWVzdA==');

@$core.Deprecated('Use queryTotalTxsResponseDescriptor instead')
const QueryTotalTxsResponse$json = {
  '1': 'QueryTotalTxsResponse',
  '2': [
    {'1': 'total_txs', '3': 1, '4': 1, '5': 4, '10': 'totalTxs'},
  ],
};

/// Descriptor for `QueryTotalTxsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalTxsResponseDescriptor = $convert.base64Decode(
    'ChVRdWVyeVRvdGFsVHhzUmVzcG9uc2USGwoJdG90YWxfdHhzGAEgASgEUgh0b3RhbFR4cw==');

@$core.Deprecated('Use unlockScheduleDescriptor instead')
const UnlockSchedule$json = {
  '1': 'UnlockSchedule',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'time_str', '3': 2, '4': 1, '5': 9, '10': 'timeStr'},
    {'1': 'time_unix', '3': 3, '4': 1, '5': 3, '10': 'timeUnix'},
    {'1': 'amount', '3': 4, '4': 1, '5': 3, '10': 'amount'},
    {'1': 'denom', '3': 5, '4': 1, '5': 9, '10': 'denom'},
  ],
};

/// Descriptor for `UnlockSchedule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlockScheduleDescriptor = $convert.base64Decode(
    'Cg5VbmxvY2tTY2hlZHVsZRIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNzEhkKCHRpbWVfc3RyGA'
    'IgASgJUgd0aW1lU3RyEhsKCXRpbWVfdW5peBgDIAEoA1IIdGltZVVuaXgSFgoGYW1vdW50GAQg'
    'ASgDUgZhbW91bnQSFAoFZGVub20YBSABKAlSBWRlbm9t');

@$core.Deprecated('Use queryUnlockSchedulesRequestDescriptor instead')
const QueryUnlockSchedulesRequest$json = {
  '1': 'QueryUnlockSchedulesRequest',
};

/// Descriptor for `QueryUnlockSchedulesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryUnlockSchedulesRequestDescriptor = $convert.base64Decode(
    'ChtRdWVyeVVubG9ja1NjaGVkdWxlc1JlcXVlc3Q=');

@$core.Deprecated('Use queryUnlockSchedulesResponseDescriptor instead')
const QueryUnlockSchedulesResponse$json = {
  '1': 'QueryUnlockSchedulesResponse',
  '2': [
    {'1': 'unlock_schedules', '3': 1, '4': 3, '5': 11, '6': '.x.paxi.types.UnlockSchedule', '10': 'unlockSchedules'},
  ],
};

/// Descriptor for `QueryUnlockSchedulesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryUnlockSchedulesResponseDescriptor = $convert.base64Decode(
    'ChxRdWVyeVVubG9ja1NjaGVkdWxlc1Jlc3BvbnNlEkcKEHVubG9ja19zY2hlZHVsZXMYASADKA'
    'syHC54LnBheGkudHlwZXMuVW5sb2NrU2NoZWR1bGVSD3VubG9ja1NjaGVkdWxlcw==');

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
    {'1': 'extra_gas_per_new_account', '3': 1, '4': 1, '5': 4, '10': 'extraGasPerNewAccount'},
  ],
};

/// Descriptor for `QueryParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryParamsResponseDescriptor = $convert.base64Decode(
    'ChNRdWVyeVBhcmFtc1Jlc3BvbnNlEjgKGWV4dHJhX2dhc19wZXJfbmV3X2FjY291bnQYASABKA'
    'RSFWV4dHJhR2FzUGVyTmV3QWNjb3VudA==');

const $core.Map<$core.String, $core.dynamic> QueryServiceBase$json = {
  '1': 'Query',
  '2': [
    {'1': 'LockedVesting', '2': '.x.paxi.types.QueryLockedVestingRequest', '3': '.x.paxi.types.QueryLockedVestingResponse', '4': {}},
    {'1': 'CirculatingSupply', '2': '.x.paxi.types.QueryCirculatingSupplyRequest', '3': '.x.paxi.types.QueryCirculatingSupplyResponse', '4': {}},
    {'1': 'TotalSupply', '2': '.x.paxi.types.QueryTotalSupplyRequest', '3': '.x.paxi.types.QueryTotalSupplyResponse', '4': {}},
    {'1': 'EstimatedGasPrice', '2': '.x.paxi.types.QueryEstimatedGasPriceRequest', '3': '.x.paxi.types.QueryEstimatedGasPriceResponse', '4': {}},
    {'1': 'LastBlockGasUsed', '2': '.x.paxi.types.QueryLastBlockGasUsedRequest', '3': '.x.paxi.types.QueryLastBlockGasUsedResponse', '4': {}},
    {'1': 'TotalTxs', '2': '.x.paxi.types.QueryTotalTxsRequest', '3': '.x.paxi.types.QueryTotalTxsResponse', '4': {}},
    {'1': 'UnlockSchedules', '2': '.x.paxi.types.QueryUnlockSchedulesRequest', '3': '.x.paxi.types.QueryUnlockSchedulesResponse', '4': {}},
    {'1': 'Params', '2': '.x.paxi.types.QueryParamsRequest', '3': '.x.paxi.types.QueryParamsResponse', '4': {}},
  ],
};

@$core.Deprecated('Use queryServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> QueryServiceBase$messageJson = {
  '.x.paxi.types.QueryLockedVestingRequest': QueryLockedVestingRequest$json,
  '.x.paxi.types.QueryLockedVestingResponse': QueryLockedVestingResponse$json,
  '.cosmos.base.v1beta1.Coin': $0.Coin$json,
  '.x.paxi.types.QueryCirculatingSupplyRequest': QueryCirculatingSupplyRequest$json,
  '.x.paxi.types.QueryCirculatingSupplyResponse': QueryCirculatingSupplyResponse$json,
  '.x.paxi.types.QueryTotalSupplyRequest': QueryTotalSupplyRequest$json,
  '.x.paxi.types.QueryTotalSupplyResponse': QueryTotalSupplyResponse$json,
  '.x.paxi.types.QueryEstimatedGasPriceRequest': QueryEstimatedGasPriceRequest$json,
  '.x.paxi.types.QueryEstimatedGasPriceResponse': QueryEstimatedGasPriceResponse$json,
  '.x.paxi.types.QueryLastBlockGasUsedRequest': QueryLastBlockGasUsedRequest$json,
  '.x.paxi.types.QueryLastBlockGasUsedResponse': QueryLastBlockGasUsedResponse$json,
  '.x.paxi.types.QueryTotalTxsRequest': QueryTotalTxsRequest$json,
  '.x.paxi.types.QueryTotalTxsResponse': QueryTotalTxsResponse$json,
  '.x.paxi.types.QueryUnlockSchedulesRequest': QueryUnlockSchedulesRequest$json,
  '.x.paxi.types.QueryUnlockSchedulesResponse': QueryUnlockSchedulesResponse$json,
  '.x.paxi.types.UnlockSchedule': UnlockSchedule$json,
  '.x.paxi.types.QueryParamsRequest': QueryParamsRequest$json,
  '.x.paxi.types.QueryParamsResponse': QueryParamsResponse$json,
};

/// Descriptor for `Query`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List queryServiceDescriptor = $convert.base64Decode(
    'CgVRdWVyeRKFAQoNTG9ja2VkVmVzdGluZxInLngucGF4aS50eXBlcy5RdWVyeUxvY2tlZFZlc3'
    'RpbmdSZXF1ZXN0GigueC5wYXhpLnR5cGVzLlF1ZXJ5TG9ja2VkVmVzdGluZ1Jlc3BvbnNlIiGC'
    '0+STAhsSGS9wYXhpL3BheGkvbG9ja2VkX3Zlc3RpbmcSlQEKEUNpcmN1bGF0aW5nU3VwcGx5Ei'
    'sueC5wYXhpLnR5cGVzLlF1ZXJ5Q2lyY3VsYXRpbmdTdXBwbHlSZXF1ZXN0GiwueC5wYXhpLnR5'
    'cGVzLlF1ZXJ5Q2lyY3VsYXRpbmdTdXBwbHlSZXNwb25zZSIlgtPkkwIfEh0vcGF4aS9wYXhpL2'
    'NpcmN1bGF0aW5nX3N1cHBseRJ9CgtUb3RhbFN1cHBseRIlLngucGF4aS50eXBlcy5RdWVyeVRv'
    'dGFsU3VwcGx5UmVxdWVzdBomLngucGF4aS50eXBlcy5RdWVyeVRvdGFsU3VwcGx5UmVzcG9uc2'
    'UiH4LT5JMCGRIXL3BheGkvcGF4aS90b3RhbF9zdXBwbHkSlgEKEUVzdGltYXRlZEdhc1ByaWNl'
    'EisueC5wYXhpLnR5cGVzLlF1ZXJ5RXN0aW1hdGVkR2FzUHJpY2VSZXF1ZXN0GiwueC5wYXhpLn'
    'R5cGVzLlF1ZXJ5RXN0aW1hdGVkR2FzUHJpY2VSZXNwb25zZSImgtPkkwIgEh4vcGF4aS9wYXhp'
    'L2VzdGltYXRlZF9nYXNfcHJpY2USkwEKEExhc3RCbG9ja0dhc1VzZWQSKi54LnBheGkudHlwZX'
    'MuUXVlcnlMYXN0QmxvY2tHYXNVc2VkUmVxdWVzdBorLngucGF4aS50eXBlcy5RdWVyeUxhc3RC'
    'bG9ja0dhc1VzZWRSZXNwb25zZSImgtPkkwIgEh4vcGF4aS9wYXhpL2xhc3RfYmxvY2tfZ2FzX3'
    'VzZWQScQoIVG90YWxUeHMSIi54LnBheGkudHlwZXMuUXVlcnlUb3RhbFR4c1JlcXVlc3QaIy54'
    'LnBheGkudHlwZXMuUXVlcnlUb3RhbFR4c1Jlc3BvbnNlIhyC0+STAhYSFC9wYXhpL3BheGkvdG'
    '90YWxfdHhzEo0BCg9VbmxvY2tTY2hlZHVsZXMSKS54LnBheGkudHlwZXMuUXVlcnlVbmxvY2tT'
    'Y2hlZHVsZXNSZXF1ZXN0GioueC5wYXhpLnR5cGVzLlF1ZXJ5VW5sb2NrU2NoZWR1bGVzUmVzcG'
    '9uc2UiI4LT5JMCHRIbL3BheGkvcGF4aS91bmxvY2tfc2NoZWR1bGVzEmgKBlBhcmFtcxIgLngu'
    'cGF4aS50eXBlcy5RdWVyeVBhcmFtc1JlcXVlc3QaIS54LnBheGkudHlwZXMuUXVlcnlQYXJhbX'
    'NSZXNwb25zZSIZgtPkkwITEhEvcGF4aS9wYXhpL3BhcmFtcw==');

