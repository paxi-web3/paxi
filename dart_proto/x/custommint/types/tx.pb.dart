// This is a generated file - do not edit.
//
// Generated from x/custommint/types/tx.proto.

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

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class ParamsInput extends $pb.GeneratedMessage {
  factory ParamsInput({
    $core.String? burnThreshold,
    $core.String? burnRatio,
    $fixnum.Int64? blocksPerYear,
    $core.String? firstYearInflation,
    $core.String? secondYearInflation,
    $core.String? otherYearInflation,
  }) {
    final result = create();
    if (burnThreshold != null) result.burnThreshold = burnThreshold;
    if (burnRatio != null) result.burnRatio = burnRatio;
    if (blocksPerYear != null) result.blocksPerYear = blocksPerYear;
    if (firstYearInflation != null) result.firstYearInflation = firstYearInflation;
    if (secondYearInflation != null) result.secondYearInflation = secondYearInflation;
    if (otherYearInflation != null) result.otherYearInflation = otherYearInflation;
    return result;
  }

  ParamsInput._();

  factory ParamsInput.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ParamsInput.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ParamsInput', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'burnThreshold')
    ..aOS(2, _omitFieldNames ? '' : 'burnRatio')
    ..aInt64(3, _omitFieldNames ? '' : 'blocksPerYear')
    ..aOS(4, _omitFieldNames ? '' : 'firstYearInflation')
    ..aOS(5, _omitFieldNames ? '' : 'secondYearInflation')
    ..aOS(6, _omitFieldNames ? '' : 'otherYearInflation')
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
  $core.String get burnThreshold => $_getSZ(0);
  @$pb.TagNumber(1)
  set burnThreshold($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasBurnThreshold() => $_has(0);
  @$pb.TagNumber(1)
  void clearBurnThreshold() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get burnRatio => $_getSZ(1);
  @$pb.TagNumber(2)
  set burnRatio($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBurnRatio() => $_has(1);
  @$pb.TagNumber(2)
  void clearBurnRatio() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get blocksPerYear => $_getI64(2);
  @$pb.TagNumber(3)
  set blocksPerYear($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasBlocksPerYear() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlocksPerYear() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get firstYearInflation => $_getSZ(3);
  @$pb.TagNumber(4)
  set firstYearInflation($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasFirstYearInflation() => $_has(3);
  @$pb.TagNumber(4)
  void clearFirstYearInflation() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get secondYearInflation => $_getSZ(4);
  @$pb.TagNumber(5)
  set secondYearInflation($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSecondYearInflation() => $_has(4);
  @$pb.TagNumber(5)
  void clearSecondYearInflation() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get otherYearInflation => $_getSZ(5);
  @$pb.TagNumber(6)
  set otherYearInflation($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasOtherYearInflation() => $_has(5);
  @$pb.TagNumber(6)
  void clearOtherYearInflation() => $_clearField(6);
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

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParams', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
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

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
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

  $async.Future<MsgUpdateParamsResponse> updateParams($pb.ClientContext? ctx, MsgUpdateParams request) =>
    _client.invoke<MsgUpdateParamsResponse>(ctx, 'Msg', 'UpdateParams', request, MsgUpdateParamsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
