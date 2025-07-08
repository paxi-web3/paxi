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

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class ParamsInput extends $pb.GeneratedMessage {
  factory ParamsInput({
    $fixnum.Int64? codeId,
    $fixnum.Int64? swapFeeBps,
    $fixnum.Int64? minLiquidity,
  }) {
    final result = create();
    if (codeId != null) result.codeId = codeId;
    if (swapFeeBps != null) result.swapFeeBps = swapFeeBps;
    if (minLiquidity != null) result.minLiquidity = minLiquidity;
    return result;
  }

  ParamsInput._();

  factory ParamsInput.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ParamsInput.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ParamsInput', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'codeId', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'swapFeeBps', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'minLiquidity', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
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
  $fixnum.Int64 get codeId => $_getI64(0);
  @$pb.TagNumber(1)
  set codeId($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCodeId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCodeId() => $_clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get swapFeeBps => $_getI64(1);
  @$pb.TagNumber(2)
  set swapFeeBps($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSwapFeeBps() => $_has(1);
  @$pb.TagNumber(2)
  void clearSwapFeeBps() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get minLiquidity => $_getI64(2);
  @$pb.TagNumber(3)
  set minLiquidity($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasMinLiquidity() => $_has(2);
  @$pb.TagNumber(3)
  void clearMinLiquidity() => $_clearField(3);
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

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParams', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
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

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgUpdateParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
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

class MsgProvideLiquidity extends $pb.GeneratedMessage {
  factory MsgProvideLiquidity({
    $core.String? creator,
    $core.String? prc20,
    $core.String? paxiAmount,
    $core.String? prc20Amount,
  }) {
    final result = create();
    if (creator != null) result.creator = creator;
    if (prc20 != null) result.prc20 = prc20;
    if (paxiAmount != null) result.paxiAmount = paxiAmount;
    if (prc20Amount != null) result.prc20Amount = prc20Amount;
    return result;
  }

  MsgProvideLiquidity._();

  factory MsgProvideLiquidity.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgProvideLiquidity.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgProvideLiquidity', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'creator')
    ..aOS(2, _omitFieldNames ? '' : 'prc20')
    ..aOS(3, _omitFieldNames ? '' : 'paxiAmount')
    ..aOS(4, _omitFieldNames ? '' : 'prc20Amount')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgProvideLiquidity clone() => MsgProvideLiquidity()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgProvideLiquidity copyWith(void Function(MsgProvideLiquidity) updates) => super.copyWith((message) => updates(message as MsgProvideLiquidity)) as MsgProvideLiquidity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgProvideLiquidity create() => MsgProvideLiquidity._();
  @$core.override
  MsgProvideLiquidity createEmptyInstance() => create();
  static $pb.PbList<MsgProvideLiquidity> createRepeated() => $pb.PbList<MsgProvideLiquidity>();
  @$core.pragma('dart2js:noInline')
  static MsgProvideLiquidity getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgProvideLiquidity>(create);
  static MsgProvideLiquidity? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get creator => $_getSZ(0);
  @$pb.TagNumber(1)
  set creator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCreator() => $_has(0);
  @$pb.TagNumber(1)
  void clearCreator() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get prc20 => $_getSZ(1);
  @$pb.TagNumber(2)
  set prc20($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPrc20() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrc20() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get paxiAmount => $_getSZ(2);
  @$pb.TagNumber(3)
  set paxiAmount($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPaxiAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearPaxiAmount() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get prc20Amount => $_getSZ(3);
  @$pb.TagNumber(4)
  set prc20Amount($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPrc20Amount() => $_has(3);
  @$pb.TagNumber(4)
  void clearPrc20Amount() => $_clearField(4);
}

class MsgProvideLiquidityResponse extends $pb.GeneratedMessage {
  factory MsgProvideLiquidityResponse() => create();

  MsgProvideLiquidityResponse._();

  factory MsgProvideLiquidityResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgProvideLiquidityResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgProvideLiquidityResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgProvideLiquidityResponse clone() => MsgProvideLiquidityResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgProvideLiquidityResponse copyWith(void Function(MsgProvideLiquidityResponse) updates) => super.copyWith((message) => updates(message as MsgProvideLiquidityResponse)) as MsgProvideLiquidityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgProvideLiquidityResponse create() => MsgProvideLiquidityResponse._();
  @$core.override
  MsgProvideLiquidityResponse createEmptyInstance() => create();
  static $pb.PbList<MsgProvideLiquidityResponse> createRepeated() => $pb.PbList<MsgProvideLiquidityResponse>();
  @$core.pragma('dart2js:noInline')
  static MsgProvideLiquidityResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgProvideLiquidityResponse>(create);
  static MsgProvideLiquidityResponse? _defaultInstance;
}

class MsgWithdrawLiquidity extends $pb.GeneratedMessage {
  factory MsgWithdrawLiquidity({
    $core.String? creator,
    $core.String? prc20,
    $core.String? lpAmount,
  }) {
    final result = create();
    if (creator != null) result.creator = creator;
    if (prc20 != null) result.prc20 = prc20;
    if (lpAmount != null) result.lpAmount = lpAmount;
    return result;
  }

  MsgWithdrawLiquidity._();

  factory MsgWithdrawLiquidity.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgWithdrawLiquidity.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgWithdrawLiquidity', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'creator')
    ..aOS(2, _omitFieldNames ? '' : 'prc20')
    ..aOS(3, _omitFieldNames ? '' : 'lpAmount')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgWithdrawLiquidity clone() => MsgWithdrawLiquidity()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgWithdrawLiquidity copyWith(void Function(MsgWithdrawLiquidity) updates) => super.copyWith((message) => updates(message as MsgWithdrawLiquidity)) as MsgWithdrawLiquidity;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgWithdrawLiquidity create() => MsgWithdrawLiquidity._();
  @$core.override
  MsgWithdrawLiquidity createEmptyInstance() => create();
  static $pb.PbList<MsgWithdrawLiquidity> createRepeated() => $pb.PbList<MsgWithdrawLiquidity>();
  @$core.pragma('dart2js:noInline')
  static MsgWithdrawLiquidity getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgWithdrawLiquidity>(create);
  static MsgWithdrawLiquidity? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get creator => $_getSZ(0);
  @$pb.TagNumber(1)
  set creator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCreator() => $_has(0);
  @$pb.TagNumber(1)
  void clearCreator() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get prc20 => $_getSZ(1);
  @$pb.TagNumber(2)
  set prc20($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPrc20() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrc20() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get lpAmount => $_getSZ(2);
  @$pb.TagNumber(3)
  set lpAmount($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLpAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearLpAmount() => $_clearField(3);
}

class MsgWithdrawLiquidityResponse extends $pb.GeneratedMessage {
  factory MsgWithdrawLiquidityResponse() => create();

  MsgWithdrawLiquidityResponse._();

  factory MsgWithdrawLiquidityResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgWithdrawLiquidityResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgWithdrawLiquidityResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgWithdrawLiquidityResponse clone() => MsgWithdrawLiquidityResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgWithdrawLiquidityResponse copyWith(void Function(MsgWithdrawLiquidityResponse) updates) => super.copyWith((message) => updates(message as MsgWithdrawLiquidityResponse)) as MsgWithdrawLiquidityResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgWithdrawLiquidityResponse create() => MsgWithdrawLiquidityResponse._();
  @$core.override
  MsgWithdrawLiquidityResponse createEmptyInstance() => create();
  static $pb.PbList<MsgWithdrawLiquidityResponse> createRepeated() => $pb.PbList<MsgWithdrawLiquidityResponse>();
  @$core.pragma('dart2js:noInline')
  static MsgWithdrawLiquidityResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgWithdrawLiquidityResponse>(create);
  static MsgWithdrawLiquidityResponse? _defaultInstance;
}

class MsgSwap extends $pb.GeneratedMessage {
  factory MsgSwap({
    $core.String? creator,
    $core.String? prc20,
    $core.String? offerDenom,
    $core.String? offerAmount,
    $core.String? minReceive,
  }) {
    final result = create();
    if (creator != null) result.creator = creator;
    if (prc20 != null) result.prc20 = prc20;
    if (offerDenom != null) result.offerDenom = offerDenom;
    if (offerAmount != null) result.offerAmount = offerAmount;
    if (minReceive != null) result.minReceive = minReceive;
    return result;
  }

  MsgSwap._();

  factory MsgSwap.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgSwap.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgSwap', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'creator')
    ..aOS(2, _omitFieldNames ? '' : 'prc20')
    ..aOS(3, _omitFieldNames ? '' : 'offerDenom')
    ..aOS(4, _omitFieldNames ? '' : 'offerAmount')
    ..aOS(5, _omitFieldNames ? '' : 'minReceive')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgSwap clone() => MsgSwap()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgSwap copyWith(void Function(MsgSwap) updates) => super.copyWith((message) => updates(message as MsgSwap)) as MsgSwap;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgSwap create() => MsgSwap._();
  @$core.override
  MsgSwap createEmptyInstance() => create();
  static $pb.PbList<MsgSwap> createRepeated() => $pb.PbList<MsgSwap>();
  @$core.pragma('dart2js:noInline')
  static MsgSwap getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgSwap>(create);
  static MsgSwap? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get creator => $_getSZ(0);
  @$pb.TagNumber(1)
  set creator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCreator() => $_has(0);
  @$pb.TagNumber(1)
  void clearCreator() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get prc20 => $_getSZ(1);
  @$pb.TagNumber(2)
  set prc20($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPrc20() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrc20() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get offerDenom => $_getSZ(2);
  @$pb.TagNumber(3)
  set offerDenom($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasOfferDenom() => $_has(2);
  @$pb.TagNumber(3)
  void clearOfferDenom() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get offerAmount => $_getSZ(3);
  @$pb.TagNumber(4)
  set offerAmount($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasOfferAmount() => $_has(3);
  @$pb.TagNumber(4)
  void clearOfferAmount() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get minReceive => $_getSZ(4);
  @$pb.TagNumber(5)
  set minReceive($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasMinReceive() => $_has(4);
  @$pb.TagNumber(5)
  void clearMinReceive() => $_clearField(5);
}

class MsgSwapResponse extends $pb.GeneratedMessage {
  factory MsgSwapResponse() => create();

  MsgSwapResponse._();

  factory MsgSwapResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory MsgSwapResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MsgSwapResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgSwapResponse clone() => MsgSwapResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MsgSwapResponse copyWith(void Function(MsgSwapResponse) updates) => super.copyWith((message) => updates(message as MsgSwapResponse)) as MsgSwapResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MsgSwapResponse create() => MsgSwapResponse._();
  @$core.override
  MsgSwapResponse createEmptyInstance() => create();
  static $pb.PbList<MsgSwapResponse> createRepeated() => $pb.PbList<MsgSwapResponse>();
  @$core.pragma('dart2js:noInline')
  static MsgSwapResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MsgSwapResponse>(create);
  static MsgSwapResponse? _defaultInstance;
}

class MsgApi {
  final $pb.RpcClient _client;

  MsgApi(this._client);

  $async.Future<MsgUpdateParamsResponse> updateParams($pb.ClientContext? ctx, MsgUpdateParams request) =>
    _client.invoke<MsgUpdateParamsResponse>(ctx, 'Msg', 'UpdateParams', request, MsgUpdateParamsResponse())
  ;
  $async.Future<MsgProvideLiquidityResponse> provideLiquidity($pb.ClientContext? ctx, MsgProvideLiquidity request) =>
    _client.invoke<MsgProvideLiquidityResponse>(ctx, 'Msg', 'ProvideLiquidity', request, MsgProvideLiquidityResponse())
  ;
  $async.Future<MsgWithdrawLiquidityResponse> withdrawLiquidity($pb.ClientContext? ctx, MsgWithdrawLiquidity request) =>
    _client.invoke<MsgWithdrawLiquidityResponse>(ctx, 'Msg', 'WithdrawLiquidity', request, MsgWithdrawLiquidityResponse())
  ;
  $async.Future<MsgSwapResponse> swap($pb.ClientContext? ctx, MsgSwap request) =>
    _client.invoke<MsgSwapResponse>(ctx, 'Msg', 'Swap', request, MsgSwapResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
