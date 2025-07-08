// This is a generated file - do not edit.
//
// Generated from x/swap/types/query.proto.

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

import '../../../cosmos/base/query/v1beta1/pagination.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class QueryParamsRequest extends $pb.GeneratedMessage {
  factory QueryParamsRequest() => create();

  QueryParamsRequest._();

  factory QueryParamsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsRequest clone() => QueryParamsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsRequest copyWith(void Function(QueryParamsRequest) updates) => super.copyWith((message) => updates(message as QueryParamsRequest)) as QueryParamsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryParamsRequest create() => QueryParamsRequest._();
  @$core.override
  QueryParamsRequest createEmptyInstance() => create();
  static $pb.PbList<QueryParamsRequest> createRepeated() => $pb.PbList<QueryParamsRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryParamsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryParamsRequest>(create);
  static QueryParamsRequest? _defaultInstance;
}

class QueryParamsResponse extends $pb.GeneratedMessage {
  factory QueryParamsResponse({
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

  QueryParamsResponse._();

  factory QueryParamsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'codeId', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'swapFeeBps', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'minLiquidity', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsResponse clone() => QueryParamsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryParamsResponse copyWith(void Function(QueryParamsResponse) updates) => super.copyWith((message) => updates(message as QueryParamsResponse)) as QueryParamsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryParamsResponse create() => QueryParamsResponse._();
  @$core.override
  QueryParamsResponse createEmptyInstance() => create();
  static $pb.PbList<QueryParamsResponse> createRepeated() => $pb.PbList<QueryParamsResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryParamsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryParamsResponse>(create);
  static QueryParamsResponse? _defaultInstance;

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

class ProviderPosition extends $pb.GeneratedMessage {
  factory ProviderPosition({
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

  ProviderPosition._();

  factory ProviderPosition.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory ProviderPosition.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProviderPosition', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'creator')
    ..aOS(2, _omitFieldNames ? '' : 'prc20')
    ..aOS(3, _omitFieldNames ? '' : 'lpAmount')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProviderPosition clone() => ProviderPosition()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProviderPosition copyWith(void Function(ProviderPosition) updates) => super.copyWith((message) => updates(message as ProviderPosition)) as ProviderPosition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProviderPosition create() => ProviderPosition._();
  @$core.override
  ProviderPosition createEmptyInstance() => create();
  static $pb.PbList<ProviderPosition> createRepeated() => $pb.PbList<ProviderPosition>();
  @$core.pragma('dart2js:noInline')
  static ProviderPosition getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProviderPosition>(create);
  static ProviderPosition? _defaultInstance;

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

class QueryPositionRequest extends $pb.GeneratedMessage {
  factory QueryPositionRequest({
    $core.String? creator,
    $core.String? prc20,
  }) {
    final result = create();
    if (creator != null) result.creator = creator;
    if (prc20 != null) result.prc20 = prc20;
    return result;
  }

  QueryPositionRequest._();

  factory QueryPositionRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryPositionRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryPositionRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'creator')
    ..aOS(2, _omitFieldNames ? '' : 'prc20')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPositionRequest clone() => QueryPositionRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPositionRequest copyWith(void Function(QueryPositionRequest) updates) => super.copyWith((message) => updates(message as QueryPositionRequest)) as QueryPositionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryPositionRequest create() => QueryPositionRequest._();
  @$core.override
  QueryPositionRequest createEmptyInstance() => create();
  static $pb.PbList<QueryPositionRequest> createRepeated() => $pb.PbList<QueryPositionRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryPositionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryPositionRequest>(create);
  static QueryPositionRequest? _defaultInstance;

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
}

class QueryPositionResponse extends $pb.GeneratedMessage {
  factory QueryPositionResponse({
    ProviderPosition? position,
    $core.String? expectedPaxi,
    $core.String? expectedPrc20,
  }) {
    final result = create();
    if (position != null) result.position = position;
    if (expectedPaxi != null) result.expectedPaxi = expectedPaxi;
    if (expectedPrc20 != null) result.expectedPrc20 = expectedPrc20;
    return result;
  }

  QueryPositionResponse._();

  factory QueryPositionResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryPositionResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryPositionResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOM<ProviderPosition>(1, _omitFieldNames ? '' : 'position', subBuilder: ProviderPosition.create)
    ..aOS(2, _omitFieldNames ? '' : 'expectedPaxi')
    ..aOS(3, _omitFieldNames ? '' : 'expectedPrc20')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPositionResponse clone() => QueryPositionResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPositionResponse copyWith(void Function(QueryPositionResponse) updates) => super.copyWith((message) => updates(message as QueryPositionResponse)) as QueryPositionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryPositionResponse create() => QueryPositionResponse._();
  @$core.override
  QueryPositionResponse createEmptyInstance() => create();
  static $pb.PbList<QueryPositionResponse> createRepeated() => $pb.PbList<QueryPositionResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryPositionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryPositionResponse>(create);
  static QueryPositionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ProviderPosition get position => $_getN(0);
  @$pb.TagNumber(1)
  set position(ProviderPosition value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPosition() => $_has(0);
  @$pb.TagNumber(1)
  void clearPosition() => $_clearField(1);
  @$pb.TagNumber(1)
  ProviderPosition ensurePosition() => $_ensure(0);

  @$pb.TagNumber(2)
  $core.String get expectedPaxi => $_getSZ(1);
  @$pb.TagNumber(2)
  set expectedPaxi($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasExpectedPaxi() => $_has(1);
  @$pb.TagNumber(2)
  void clearExpectedPaxi() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get expectedPrc20 => $_getSZ(2);
  @$pb.TagNumber(3)
  set expectedPrc20($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasExpectedPrc20() => $_has(2);
  @$pb.TagNumber(3)
  void clearExpectedPrc20() => $_clearField(3);
}

class QueryPoolRequest extends $pb.GeneratedMessage {
  factory QueryPoolRequest({
    $core.String? prc20,
  }) {
    final result = create();
    if (prc20 != null) result.prc20 = prc20;
    return result;
  }

  QueryPoolRequest._();

  factory QueryPoolRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryPoolRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryPoolRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'prc20')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPoolRequest clone() => QueryPoolRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPoolRequest copyWith(void Function(QueryPoolRequest) updates) => super.copyWith((message) => updates(message as QueryPoolRequest)) as QueryPoolRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryPoolRequest create() => QueryPoolRequest._();
  @$core.override
  QueryPoolRequest createEmptyInstance() => create();
  static $pb.PbList<QueryPoolRequest> createRepeated() => $pb.PbList<QueryPoolRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryPoolRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryPoolRequest>(create);
  static QueryPoolRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get prc20 => $_getSZ(0);
  @$pb.TagNumber(1)
  set prc20($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPrc20() => $_has(0);
  @$pb.TagNumber(1)
  void clearPrc20() => $_clearField(1);
}

class QueryPoolResponse extends $pb.GeneratedMessage {
  factory QueryPoolResponse({
    $core.String? prc20,
    $core.String? reservePaxi,
    $core.String? reservePrc20,
    $core.String? pricePaxiPerPrc20,
    $core.String? pricePrc20PerPaxi,
    $core.String? totalShares,
  }) {
    final result = create();
    if (prc20 != null) result.prc20 = prc20;
    if (reservePaxi != null) result.reservePaxi = reservePaxi;
    if (reservePrc20 != null) result.reservePrc20 = reservePrc20;
    if (pricePaxiPerPrc20 != null) result.pricePaxiPerPrc20 = pricePaxiPerPrc20;
    if (pricePrc20PerPaxi != null) result.pricePrc20PerPaxi = pricePrc20PerPaxi;
    if (totalShares != null) result.totalShares = totalShares;
    return result;
  }

  QueryPoolResponse._();

  factory QueryPoolResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryPoolResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryPoolResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'prc20')
    ..aOS(2, _omitFieldNames ? '' : 'reservePaxi')
    ..aOS(3, _omitFieldNames ? '' : 'reservePrc20')
    ..aOS(4, _omitFieldNames ? '' : 'pricePaxiPerPrc20')
    ..aOS(5, _omitFieldNames ? '' : 'pricePrc20PerPaxi')
    ..aOS(6, _omitFieldNames ? '' : 'totalShares')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPoolResponse clone() => QueryPoolResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryPoolResponse copyWith(void Function(QueryPoolResponse) updates) => super.copyWith((message) => updates(message as QueryPoolResponse)) as QueryPoolResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryPoolResponse create() => QueryPoolResponse._();
  @$core.override
  QueryPoolResponse createEmptyInstance() => create();
  static $pb.PbList<QueryPoolResponse> createRepeated() => $pb.PbList<QueryPoolResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryPoolResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryPoolResponse>(create);
  static QueryPoolResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get prc20 => $_getSZ(0);
  @$pb.TagNumber(1)
  set prc20($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPrc20() => $_has(0);
  @$pb.TagNumber(1)
  void clearPrc20() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get reservePaxi => $_getSZ(1);
  @$pb.TagNumber(2)
  set reservePaxi($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReservePaxi() => $_has(1);
  @$pb.TagNumber(2)
  void clearReservePaxi() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get reservePrc20 => $_getSZ(2);
  @$pb.TagNumber(3)
  set reservePrc20($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasReservePrc20() => $_has(2);
  @$pb.TagNumber(3)
  void clearReservePrc20() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get pricePaxiPerPrc20 => $_getSZ(3);
  @$pb.TagNumber(4)
  set pricePaxiPerPrc20($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPricePaxiPerPrc20() => $_has(3);
  @$pb.TagNumber(4)
  void clearPricePaxiPerPrc20() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get pricePrc20PerPaxi => $_getSZ(4);
  @$pb.TagNumber(5)
  set pricePrc20PerPaxi($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPricePrc20PerPaxi() => $_has(4);
  @$pb.TagNumber(5)
  void clearPricePrc20PerPaxi() => $_clearField(5);

  @$pb.TagNumber(6)
  $core.String get totalShares => $_getSZ(5);
  @$pb.TagNumber(6)
  set totalShares($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasTotalShares() => $_has(5);
  @$pb.TagNumber(6)
  void clearTotalShares() => $_clearField(6);
}

class QueryAllPoolsRequest extends $pb.GeneratedMessage {
  factory QueryAllPoolsRequest({
    $0.PageRequest? pagination,
  }) {
    final result = create();
    if (pagination != null) result.pagination = pagination;
    return result;
  }

  QueryAllPoolsRequest._();

  factory QueryAllPoolsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryAllPoolsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryAllPoolsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOM<$0.PageRequest>(1, _omitFieldNames ? '' : 'pagination', subBuilder: $0.PageRequest.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryAllPoolsRequest clone() => QueryAllPoolsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryAllPoolsRequest copyWith(void Function(QueryAllPoolsRequest) updates) => super.copyWith((message) => updates(message as QueryAllPoolsRequest)) as QueryAllPoolsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryAllPoolsRequest create() => QueryAllPoolsRequest._();
  @$core.override
  QueryAllPoolsRequest createEmptyInstance() => create();
  static $pb.PbList<QueryAllPoolsRequest> createRepeated() => $pb.PbList<QueryAllPoolsRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryAllPoolsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryAllPoolsRequest>(create);
  static QueryAllPoolsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $0.PageRequest get pagination => $_getN(0);
  @$pb.TagNumber(1)
  set pagination($0.PageRequest value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPagination() => $_has(0);
  @$pb.TagNumber(1)
  void clearPagination() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.PageRequest ensurePagination() => $_ensure(0);
}

class QueryAllPoolsResponse extends $pb.GeneratedMessage {
  factory QueryAllPoolsResponse({
    $core.Iterable<QueryPoolResponse>? pools,
    $0.PageResponse? pagination,
  }) {
    final result = create();
    if (pools != null) result.pools.addAll(pools);
    if (pagination != null) result.pagination = pagination;
    return result;
  }

  QueryAllPoolsResponse._();

  factory QueryAllPoolsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryAllPoolsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryAllPoolsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..pc<QueryPoolResponse>(1, _omitFieldNames ? '' : 'pools', $pb.PbFieldType.PM, subBuilder: QueryPoolResponse.create)
    ..aOM<$0.PageResponse>(2, _omitFieldNames ? '' : 'pagination', subBuilder: $0.PageResponse.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryAllPoolsResponse clone() => QueryAllPoolsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryAllPoolsResponse copyWith(void Function(QueryAllPoolsResponse) updates) => super.copyWith((message) => updates(message as QueryAllPoolsResponse)) as QueryAllPoolsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryAllPoolsResponse create() => QueryAllPoolsResponse._();
  @$core.override
  QueryAllPoolsResponse createEmptyInstance() => create();
  static $pb.PbList<QueryAllPoolsResponse> createRepeated() => $pb.PbList<QueryAllPoolsResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryAllPoolsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryAllPoolsResponse>(create);
  static QueryAllPoolsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<QueryPoolResponse> get pools => $_getList(0);

  @$pb.TagNumber(2)
  $0.PageResponse get pagination => $_getN(1);
  @$pb.TagNumber(2)
  set pagination($0.PageResponse value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasPagination() => $_has(1);
  @$pb.TagNumber(2)
  void clearPagination() => $_clearField(2);
  @$pb.TagNumber(2)
  $0.PageResponse ensurePagination() => $_ensure(1);
}

class PoolProto extends $pb.GeneratedMessage {
  factory PoolProto({
    $core.String? prc20,
    $core.String? reservePaxi,
    $core.String? reservePrc20,
    $core.String? totalShares,
  }) {
    final result = create();
    if (prc20 != null) result.prc20 = prc20;
    if (reservePaxi != null) result.reservePaxi = reservePaxi;
    if (reservePrc20 != null) result.reservePrc20 = reservePrc20;
    if (totalShares != null) result.totalShares = totalShares;
    return result;
  }

  PoolProto._();

  factory PoolProto.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory PoolProto.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'PoolProto', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.swap.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'prc20')
    ..aOS(2, _omitFieldNames ? '' : 'reservePaxi')
    ..aOS(3, _omitFieldNames ? '' : 'reservePrc20')
    ..aOS(4, _omitFieldNames ? '' : 'totalShares')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PoolProto clone() => PoolProto()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PoolProto copyWith(void Function(PoolProto) updates) => super.copyWith((message) => updates(message as PoolProto)) as PoolProto;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PoolProto create() => PoolProto._();
  @$core.override
  PoolProto createEmptyInstance() => create();
  static $pb.PbList<PoolProto> createRepeated() => $pb.PbList<PoolProto>();
  @$core.pragma('dart2js:noInline')
  static PoolProto getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PoolProto>(create);
  static PoolProto? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get prc20 => $_getSZ(0);
  @$pb.TagNumber(1)
  set prc20($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPrc20() => $_has(0);
  @$pb.TagNumber(1)
  void clearPrc20() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get reservePaxi => $_getSZ(1);
  @$pb.TagNumber(2)
  set reservePaxi($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReservePaxi() => $_has(1);
  @$pb.TagNumber(2)
  void clearReservePaxi() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get reservePrc20 => $_getSZ(2);
  @$pb.TagNumber(3)
  set reservePrc20($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasReservePrc20() => $_has(2);
  @$pb.TagNumber(3)
  void clearReservePrc20() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get totalShares => $_getSZ(3);
  @$pb.TagNumber(4)
  set totalShares($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasTotalShares() => $_has(3);
  @$pb.TagNumber(4)
  void clearTotalShares() => $_clearField(4);
}

class QueryApi {
  final $pb.RpcClient _client;

  QueryApi(this._client);

  $async.Future<QueryParamsResponse> params($pb.ClientContext? ctx, QueryParamsRequest request) =>
    _client.invoke<QueryParamsResponse>(ctx, 'Query', 'Params', request, QueryParamsResponse())
  ;
  $async.Future<QueryPositionResponse> position($pb.ClientContext? ctx, QueryPositionRequest request) =>
    _client.invoke<QueryPositionResponse>(ctx, 'Query', 'Position', request, QueryPositionResponse())
  ;
  $async.Future<QueryPoolResponse> pool($pb.ClientContext? ctx, QueryPoolRequest request) =>
    _client.invoke<QueryPoolResponse>(ctx, 'Query', 'Pool', request, QueryPoolResponse())
  ;
  $async.Future<QueryAllPoolsResponse> allPools($pb.ClientContext? ctx, QueryAllPoolsRequest request) =>
    _client.invoke<QueryAllPoolsResponse>(ctx, 'Query', 'AllPools', request, QueryAllPoolsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
