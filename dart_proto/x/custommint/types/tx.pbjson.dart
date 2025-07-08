// This is a generated file - do not edit.
//
// Generated from x/custommint/types/tx.proto.

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
    {'1': 'burn_threshold', '3': 1, '4': 1, '5': 9, '10': 'burnThreshold'},
    {'1': 'burn_ratio', '3': 2, '4': 1, '5': 9, '10': 'burnRatio'},
    {'1': 'blocks_per_year', '3': 3, '4': 1, '5': 3, '10': 'blocksPerYear'},
    {'1': 'first_year_inflation', '3': 4, '4': 1, '5': 9, '10': 'firstYearInflation'},
    {'1': 'second_year_inflation', '3': 5, '4': 1, '5': 9, '10': 'secondYearInflation'},
    {'1': 'other_year_inflation', '3': 6, '4': 1, '5': 9, '10': 'otherYearInflation'},
  ],
};

/// Descriptor for `ParamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List paramsInputDescriptor = $convert.base64Decode(
    'CgtQYXJhbXNJbnB1dBIlCg5idXJuX3RocmVzaG9sZBgBIAEoCVINYnVyblRocmVzaG9sZBIdCg'
    'pidXJuX3JhdGlvGAIgASgJUglidXJuUmF0aW8SJgoPYmxvY2tzX3Blcl95ZWFyGAMgASgDUg1i'
    'bG9ja3NQZXJZZWFyEjAKFGZpcnN0X3llYXJfaW5mbGF0aW9uGAQgASgJUhJmaXJzdFllYXJJbm'
    'ZsYXRpb24SMgoVc2Vjb25kX3llYXJfaW5mbGF0aW9uGAUgASgJUhNzZWNvbmRZZWFySW5mbGF0'
    'aW9uEjAKFG90aGVyX3llYXJfaW5mbGF0aW9uGAYgASgJUhJvdGhlclllYXJJbmZsYXRpb24=');

@$core.Deprecated('Use msgUpdateParamsDescriptor instead')
const MsgUpdateParams$json = {
  '1': 'MsgUpdateParams',
  '2': [
    {'1': 'authority', '3': 1, '4': 1, '5': 9, '10': 'authority'},
    {'1': 'params', '3': 2, '4': 1, '5': 11, '6': '.x.custommint.types.ParamsInput', '8': {}, '10': 'params'},
  ],
  '7': {},
};

/// Descriptor for `MsgUpdateParams`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgUpdateParamsDescriptor = $convert.base64Decode(
    'Cg9Nc2dVcGRhdGVQYXJhbXMSHAoJYXV0aG9yaXR5GAEgASgJUglhdXRob3JpdHkSPQoGcGFyYW'
    '1zGAIgASgLMh8ueC5jdXN0b21taW50LnR5cGVzLlBhcmFtc0lucHV0QgTI3h8AUgZwYXJhbXM6'
    'DoLnsCoJYXV0aG9yaXR5');

@$core.Deprecated('Use msgUpdateParamsResponseDescriptor instead')
const MsgUpdateParamsResponse$json = {
  '1': 'MsgUpdateParamsResponse',
};

/// Descriptor for `MsgUpdateParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgUpdateParamsResponseDescriptor = $convert.base64Decode(
    'ChdNc2dVcGRhdGVQYXJhbXNSZXNwb25zZQ==');

const $core.Map<$core.String, $core.dynamic> MsgServiceBase$json = {
  '1': 'Msg',
  '2': [
    {'1': 'UpdateParams', '2': '.x.custommint.types.MsgUpdateParams', '3': '.x.custommint.types.MsgUpdateParamsResponse'},
  ],
  '3': {},
};

@$core.Deprecated('Use msgServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> MsgServiceBase$messageJson = {
  '.x.custommint.types.MsgUpdateParams': MsgUpdateParams$json,
  '.x.custommint.types.ParamsInput': ParamsInput$json,
  '.x.custommint.types.MsgUpdateParamsResponse': MsgUpdateParamsResponse$json,
};

/// Descriptor for `Msg`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List msgServiceDescriptor = $convert.base64Decode(
    'CgNNc2cSYAoMVXBkYXRlUGFyYW1zEiMueC5jdXN0b21taW50LnR5cGVzLk1zZ1VwZGF0ZVBhcm'
    'FtcxorLnguY3VzdG9tbWludC50eXBlcy5Nc2dVcGRhdGVQYXJhbXNSZXNwb25zZRoFgOewKgE=');

