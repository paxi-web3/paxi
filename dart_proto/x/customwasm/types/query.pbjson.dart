// This is a generated file - do not edit.
//
// Generated from x/customwasm/types/query.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

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
    {'1': 'store_code_base_gas', '3': 1, '4': 1, '5': 4, '10': 'storeCodeBaseGas'},
    {'1': 'store_code_multiplier', '3': 2, '4': 1, '5': 4, '10': 'storeCodeMultiplier'},
    {'1': 'inst_base_gas', '3': 3, '4': 1, '5': 4, '10': 'instBaseGas'},
    {'1': 'inst_multiplier', '3': 4, '4': 1, '5': 4, '10': 'instMultiplier'},
  ],
};

/// Descriptor for `QueryParamsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queryParamsResponseDescriptor = $convert.base64Decode(
    'ChNRdWVyeVBhcmFtc1Jlc3BvbnNlEi0KE3N0b3JlX2NvZGVfYmFzZV9nYXMYASABKARSEHN0b3'
    'JlQ29kZUJhc2VHYXMSMgoVc3RvcmVfY29kZV9tdWx0aXBsaWVyGAIgASgEUhNzdG9yZUNvZGVN'
    'dWx0aXBsaWVyEiIKDWluc3RfYmFzZV9nYXMYAyABKARSC2luc3RCYXNlR2FzEicKD2luc3RfbX'
    'VsdGlwbGllchgEIAEoBFIOaW5zdE11bHRpcGxpZXI=');

const $core.Map<$core.String, $core.dynamic> QueryServiceBase$json = {
  '1': 'Query',
  '2': [
    {'1': 'Params', '2': '.x.customwasm.types.QueryParamsRequest', '3': '.x.customwasm.types.QueryParamsResponse', '4': {}},
  ],
};

@$core.Deprecated('Use queryServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> QueryServiceBase$messageJson = {
  '.x.customwasm.types.QueryParamsRequest': QueryParamsRequest$json,
  '.x.customwasm.types.QueryParamsResponse': QueryParamsResponse$json,
};

/// Descriptor for `Query`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List queryServiceDescriptor = $convert.base64Decode(
    'CgVRdWVyeRJ6CgZQYXJhbXMSJi54LmN1c3RvbXdhc20udHlwZXMuUXVlcnlQYXJhbXNSZXF1ZX'
    'N0GicueC5jdXN0b213YXNtLnR5cGVzLlF1ZXJ5UGFyYW1zUmVzcG9uc2UiH4LT5JMCGRIXL3Bh'
    'eGkvY3VzdG9td2FzbS9wYXJhbXM=');

