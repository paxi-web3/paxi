// This is a generated file - do not edit.
//
// Generated from x/custommint/types/query.proto.

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

@$core.Deprecated('Use queryTotalMintedRequestDescriptor instead')
const QueryTotalMintedRequest$json = {
  '1': 'QueryTotalMintedRequest',
};

/// Descriptor for `QueryTotalMintedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalMintedRequestDescriptor = $convert.base64Decode(
    'ChdRdWVyeVRvdGFsTWludGVkUmVxdWVzdA==');

@$core.Deprecated('Use queryTotalMintedResponseDescriptor instead')
const QueryTotalMintedResponse$json = {
  '1': 'QueryTotalMintedResponse',
  '2': [
    {'1': 'amount', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
};

/// Descriptor for `QueryTotalMintedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalMintedResponseDescriptor = $convert.base64Decode(
    'ChhRdWVyeVRvdGFsTWludGVkUmVzcG9uc2USMQoGYW1vdW50GAEgASgLMhkuY29zbW9zLmJhc2'
    'UudjFiZXRhMS5Db2luUgZhbW91bnQ=');

@$core.Deprecated('Use queryTotalBurnedRequestDescriptor instead')
const QueryTotalBurnedRequest$json = {
  '1': 'QueryTotalBurnedRequest',
};

/// Descriptor for `QueryTotalBurnedRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalBurnedRequestDescriptor = $convert.base64Decode(
    'ChdRdWVyeVRvdGFsQnVybmVkUmVxdWVzdA==');

@$core.Deprecated('Use queryTotalBurnedResponseDescriptor instead')
const QueryTotalBurnedResponse$json = {
  '1': 'QueryTotalBurnedResponse',
  '2': [
    {'1': 'amount', '3': 1, '4': 1, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
};

/// Descriptor for `QueryTotalBurnedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryTotalBurnedResponseDescriptor = $convert.base64Decode(
    'ChhRdWVyeVRvdGFsQnVybmVkUmVzcG9uc2USMQoGYW1vdW50GAEgASgLMhkuY29zbW9zLmJhc2'
    'UudjFiZXRhMS5Db2luUgZhbW91bnQ=');

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
    {'1': 'burn_threshold', '3': 1, '4': 1, '5': 9, '10': 'burnThreshold'},
    {'1': 'burn_ratio', '3': 2, '4': 1, '5': 9, '10': 'burnRatio'},
    {'1': 'blocks_per_year', '3': 3, '4': 1, '5': 3, '10': 'blocksPerYear'},
    {'1': 'first_year_inflation', '3': 4, '4': 1, '5': 9, '10': 'firstYearInflation'},
    {'1': 'second_year_inflation', '3': 5, '4': 1, '5': 9, '10': 'secondYearInflation'},
    {'1': 'other_year_inflation', '3': 6, '4': 1, '5': 9, '10': 'otherYearInflation'},
  ],
};

/// Descriptor for `QueryParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryParamsResponseDescriptor = $convert.base64Decode(
    'ChNRdWVyeVBhcmFtc1Jlc3BvbnNlEiUKDmJ1cm5fdGhyZXNob2xkGAEgASgJUg1idXJuVGhyZX'
    'Nob2xkEh0KCmJ1cm5fcmF0aW8YAiABKAlSCWJ1cm5SYXRpbxImCg9ibG9ja3NfcGVyX3llYXIY'
    'AyABKANSDWJsb2Nrc1BlclllYXISMAoUZmlyc3RfeWVhcl9pbmZsYXRpb24YBCABKAlSEmZpcn'
    'N0WWVhckluZmxhdGlvbhIyChVzZWNvbmRfeWVhcl9pbmZsYXRpb24YBSABKAlSE3NlY29uZFll'
    'YXJJbmZsYXRpb24SMAoUb3RoZXJfeWVhcl9pbmZsYXRpb24YBiABKAlSEm90aGVyWWVhckluZm'
    'xhdGlvbg==');

const $core.Map<$core.String, $core.dynamic> QueryServiceBase$json = {
  '1': 'Query',
  '2': [
    {'1': 'TotalMinted', '2': '.x.custommint.types.QueryTotalMintedRequest', '3': '.x.custommint.types.QueryTotalMintedResponse', '4': {}},
    {'1': 'TotalBurned', '2': '.x.custommint.types.QueryTotalBurnedRequest', '3': '.x.custommint.types.QueryTotalBurnedResponse', '4': {}},
    {'1': 'Params', '2': '.x.custommint.types.QueryParamsRequest', '3': '.x.custommint.types.QueryParamsResponse', '4': {}},
  ],
};

@$core.Deprecated('Use queryServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> QueryServiceBase$messageJson = {
  '.x.custommint.types.QueryTotalMintedRequest': QueryTotalMintedRequest$json,
  '.x.custommint.types.QueryTotalMintedResponse': QueryTotalMintedResponse$json,
  '.cosmos.base.v1beta1.Coin': $0.Coin$json,
  '.x.custommint.types.QueryTotalBurnedRequest': QueryTotalBurnedRequest$json,
  '.x.custommint.types.QueryTotalBurnedResponse': QueryTotalBurnedResponse$json,
  '.x.custommint.types.QueryParamsRequest': QueryParamsRequest$json,
  '.x.custommint.types.QueryParamsResponse': QueryParamsResponse$json,
};

/// Descriptor for `Query`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List queryServiceDescriptor = $convert.base64Decode(
    'CgVRdWVyeRKPAQoLVG90YWxNaW50ZWQSKy54LmN1c3RvbW1pbnQudHlwZXMuUXVlcnlUb3RhbE'
    '1pbnRlZFJlcXVlc3QaLC54LmN1c3RvbW1pbnQudHlwZXMuUXVlcnlUb3RhbE1pbnRlZFJlc3Bv'
    'bnNlIiWC0+STAh8SHS9wYXhpL2N1c3RvbW1pbnQvdG90YWxfbWludGVkEo8BCgtUb3RhbEJ1cm'
    '5lZBIrLnguY3VzdG9tbWludC50eXBlcy5RdWVyeVRvdGFsQnVybmVkUmVxdWVzdBosLnguY3Vz'
    'dG9tbWludC50eXBlcy5RdWVyeVRvdGFsQnVybmVkUmVzcG9uc2UiJYLT5JMCHxIdL3BheGkvY3'
    'VzdG9tbWludC90b3RhbF9idXJuZWQSegoGUGFyYW1zEiYueC5jdXN0b21taW50LnR5cGVzLlF1'
    'ZXJ5UGFyYW1zUmVxdWVzdBonLnguY3VzdG9tbWludC50eXBlcy5RdWVyeVBhcmFtc1Jlc3Bvbn'
    'NlIh+C0+STAhkSFy9wYXhpL2N1c3RvbW1pbnQvcGFyYW1z');

