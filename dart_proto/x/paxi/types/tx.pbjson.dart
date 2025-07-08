// This is a generated file - do not edit.
//
// Generated from x/paxi/types/tx.proto.

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

@$core.Deprecated('Use msgBurnTokenDescriptor instead')
const MsgBurnToken$json = {
  '1': 'MsgBurnToken',
  '2': [
    {'1': 'sender', '3': 1, '4': 1, '5': 9, '10': 'sender'},
    {'1': 'amount', '3': 2, '4': 3, '5': 11, '6': '.cosmos.base.v1beta1.Coin', '10': 'amount'},
  ],
  '7': {},
};

/// Descriptor for `MsgBurnToken`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgBurnTokenDescriptor = $convert.base64Decode(
    'CgxNc2dCdXJuVG9rZW4SFgoGc2VuZGVyGAEgASgJUgZzZW5kZXISMQoGYW1vdW50GAIgAygLMh'
    'kuY29zbW9zLmJhc2UudjFiZXRhMS5Db2luUgZhbW91bnQ6C4LnsCoGc2VuZGVy');

@$core.Deprecated('Use msgBurnTokenResponseDescriptor instead')
const MsgBurnTokenResponse$json = {
  '1': 'MsgBurnTokenResponse',
};

/// Descriptor for `MsgBurnTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgBurnTokenResponseDescriptor = $convert.base64Decode(
    'ChRNc2dCdXJuVG9rZW5SZXNwb25zZQ==');

@$core.Deprecated('Use paramsInputDescriptor instead')
const ParamsInput$json = {
  '1': 'ParamsInput',
  '2': [
    {'1': 'extra_gas_per_new_account', '3': 1, '4': 1, '5': 4, '10': 'extraGasPerNewAccount'},
  ],
};

/// Descriptor for `ParamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List paramsInputDescriptor = $convert.base64Decode(
    'CgtQYXJhbXNJbnB1dBI4ChlleHRyYV9nYXNfcGVyX25ld19hY2NvdW50GAEgASgEUhVleHRyYU'
    'dhc1Blck5ld0FjY291bnQ=');

@$core.Deprecated('Use msgUpdateParamsDescriptor instead')
const MsgUpdateParams$json = {
  '1': 'MsgUpdateParams',
  '2': [
    {'1': 'authority', '3': 1, '4': 1, '5': 9, '10': 'authority'},
    {'1': 'params', '3': 2, '4': 1, '5': 11, '6': '.x.paxi.types.ParamsInput', '8': {}, '10': 'params'},
  ],
  '7': {},
};

/// Descriptor for `MsgUpdateParams`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List msgUpdateParamsDescriptor = $convert.base64Decode(
    'Cg9Nc2dVcGRhdGVQYXJhbXMSHAoJYXV0aG9yaXR5GAEgASgJUglhdXRob3JpdHkSNwoGcGFyYW'
    '1zGAIgASgLMhkueC5wYXhpLnR5cGVzLlBhcmFtc0lucHV0QgTI3h8AUgZwYXJhbXM6DoLnsCoJ'
    'YXV0aG9yaXR5');

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
    {'1': 'BurnToken', '2': '.x.paxi.types.MsgBurnToken', '3': '.x.paxi.types.MsgBurnTokenResponse', '4': {}},
    {'1': 'UpdateParams', '2': '.x.paxi.types.MsgUpdateParams', '3': '.x.paxi.types.MsgUpdateParamsResponse'},
  ],
  '3': {},
};

@$core.Deprecated('Use msgServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> MsgServiceBase$messageJson = {
  '.x.paxi.types.MsgBurnToken': MsgBurnToken$json,
  '.cosmos.base.v1beta1.Coin': $0.Coin$json,
  '.x.paxi.types.MsgBurnTokenResponse': MsgBurnTokenResponse$json,
  '.x.paxi.types.MsgUpdateParams': MsgUpdateParams$json,
  '.x.paxi.types.ParamsInput': ParamsInput$json,
  '.x.paxi.types.MsgUpdateParamsResponse': MsgUpdateParamsResponse$json,
};

/// Descriptor for `Msg`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List msgServiceDescriptor = $convert.base64Decode(
    'CgNNc2cSawoJQnVyblRva2VuEhoueC5wYXhpLnR5cGVzLk1zZ0J1cm5Ub2tlbhoiLngucGF4aS'
    '50eXBlcy5Nc2dCdXJuVG9rZW5SZXNwb25zZSIegtPkkwIYIhMvdHgvcGF4aS9idXJuX3Rva2Vu'
    'OgEqElQKDFVwZGF0ZVBhcmFtcxIdLngucGF4aS50eXBlcy5Nc2dVcGRhdGVQYXJhbXMaJS54Ln'
    'BheGkudHlwZXMuTXNnVXBkYXRlUGFyYW1zUmVzcG9uc2UaBYDnsCoB');

