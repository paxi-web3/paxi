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

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import '../../../cosmos/base/v1beta1/coin.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class MsgBurnToken extends $pb.GeneratedMessage {
  factory MsgBurnToken({
    $core.String? sender,
    $core.Iterable<$0.Coin>? amount,
  }) {
    final result = create();
    if (sender != null) result.sender = sender;
    if (amount != null) result.amount.addAll(amount);
    return result;
  }

  MsgBurnToken._();

  factory MsgBurnToken.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgBurnToken.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgBurnToken', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sender')
    ..pc<$0.Coin>(2, _omitFieldNames ? '' : 'amount', $pb.PbFieldType.PM, subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgBurnToken clone() => MsgBurnToken()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgBurnToken copyWith(void Function(MsgBurnToken) updates) => super.copyWith((message) => updates(message as MsgBurnToken)) as MsgBurnToken;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgBurnToken create() => MsgBurnToken._();
  @$core.override
  MsgBurnToken createEmptyInstance() => create();
  static $pb.PbList<MsgBurnToken> createRepeated() => $pb.PbList<MsgBurnToken>();
  @$core.pragma('dart2js:noInline')
  static MsgBurnToken getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgBurnToken>(create);
  static MsgBurnToken? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$0.Coin> get amount => $_getList(1);
}

class MsgBurnTokenResponse extends $pb.GeneratedMessage {
  factory MsgBurnTokenResponse() => create();

  MsgBurnTokenResponse._();

  factory MsgBurnTokenResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgBurnTokenResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgBurnTokenResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgBurnTokenResponse clone() => MsgBurnTokenResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgBurnTokenResponse copyWith(void Function(MsgBurnTokenResponse) updates) => super.copyWith((message) => updates(message as MsgBurnTokenResponse)) as MsgBurnTokenResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgBurnTokenResponse create() => MsgBurnTokenResponse._();
  @$core.override
  MsgBurnTokenResponse createEmptyInstance() => create();
  static $pb.PbList<MsgBurnTokenResponse> createRepeated() => $pb.PbList<MsgBurnTokenResponse>();
  @$core.pragma('dart2js:noInline')
  static MsgBurnTokenResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgBurnTokenResponse>(create);
  static MsgBurnTokenResponse? _defaultInstance;
}

class ParamsInput extends $pb.GeneratedMessage {
  factory ParamsInput({
    $fixnum.Int64? extraGasPerNewAccount,
  }) {
    final result = create();
    if (extraGasPerNewAccount != null) result.extraGasPerNewAccount = extraGasPerNewAccount;
    return result;
  }

  ParamsInput._();

  factory ParamsInput.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ParamsInput.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ParamsInput', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'extraGasPerNewAccount', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParamsInput clone() => ParamsInput()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ParamsInput copyWith(void Function(ParamsInput) updates) => super.copyWith((message) => updates(message as ParamsInput)) as ParamsInput;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ParamsInput create() => ParamsInput._();
  @$core.override
  ParamsInput createEmptyInstance() => create();
  static $pb.PbList<ParamsInput> createRepeated() => $pb.PbList<ParamsInput>();
  @$core.pragma('dart2js:noInline')
  static ParamsInput getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ParamsInput>(create);
  static ParamsInput? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get extraGasPerNewAccount => $_getI64(0);
  @$pb.TagNumber(1)
  set extraGasPerNewAccount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasExtraGasPerNewAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearExtraGasPerNewAccount() => $_clearField(1);
}

class MsgUpdateParams extends $pb.GeneratedMessage {
  factory MsgUpdateParams({
    $core.String? authority,
    ParamsInput? params,
  }) {
    final result = create();
    if (authority != null) result.authority = authority;
    if (params != null) result.params = params;
    return result;
  }

  MsgUpdateParams._();

  factory MsgUpdateParams.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgUpdateParams.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParams', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'authority')
    ..aOM<ParamsInput>(2, _omitFieldNames ? '' : 'params', subBuilder: ParamsInput.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgUpdateParams clone() => MsgUpdateParams()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgUpdateParams copyWith(void Function(MsgUpdateParams) updates) => super.copyWith((message) => updates(message as MsgUpdateParams)) as MsgUpdateParams;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgUpdateParams create() => MsgUpdateParams._();
  @$core.override
  MsgUpdateParams createEmptyInstance() => create();
  static $pb.PbList<MsgUpdateParams> createRepeated() => $pb.PbList<MsgUpdateParams>();
  @$core.pragma('dart2js:noInline')
  static MsgUpdateParams getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgUpdateParams>(create);
  static MsgUpdateParams? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get authority => $_getSZ(0);
  @$pb.TagNumber(1)
  set authority($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAuthority() => $_has(0);
  @$pb.TagNumber(1)
  void clearAuthority() => $_clearField(1);

  @$pb.TagNumber(2)
  ParamsInput get params => $_getN(1);
  @$pb.TagNumber(2)
  set params(ParamsInput value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasParams() => $_has(1);
  @$pb.TagNumber(2)
  void clearParams() => $_clearField(2);
  @$pb.TagNumber(2)
  ParamsInput ensureParams() => $_ensure(1);
}

class MsgUpdateParamsResponse extends $pb.GeneratedMessage {
  factory MsgUpdateParamsResponse() => create();

  MsgUpdateParamsResponse._();

  factory MsgUpdateParamsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgUpdateParamsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgUpdateParamsResponse clone() => MsgUpdateParamsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgUpdateParamsResponse copyWith(void Function(MsgUpdateParamsResponse) updates) => super.copyWith((message) => updates(message as MsgUpdateParamsResponse)) as MsgUpdateParamsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgUpdateParamsResponse create() => MsgUpdateParamsResponse._();
  @$core.override
  MsgUpdateParamsResponse createEmptyInstance() => create();
  static $pb.PbList<MsgUpdateParamsResponse> createRepeated() => $pb.PbList<MsgUpdateParamsResponse>();
  @$core.pragma('dart2js:noInline')
  static MsgUpdateParamsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgUpdateParamsResponse>(create);
  static MsgUpdateParamsResponse? _defaultInstance;
}

class MsgApi {
  final $pb.RpcClient _client;

  MsgApi(this._client);

  $async.Future<MsgBurnTokenResponse> burnToken($pb.ClientContext? ctx, MsgBurnToken request) =>
    _client.invoke<MsgBurnTokenResponse>(ctx, 'Msg', 'BurnToken', request, MsgBurnTokenResponse())
  ;
  $async.Future<MsgUpdateParamsResponse> updateParams($pb.ClientContext? ctx, MsgUpdateParams request) =>
    _client.invoke<MsgUpdateParamsResponse>(ctx, 'Msg', 'UpdateParams', request, MsgUpdateParamsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
