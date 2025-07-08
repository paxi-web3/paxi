// This is a generated file - do not edit.
//
// Generated from x/paxi/types/query.proto.

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

class QueryLockedVestingRequest extends $pb.GeneratedMessage {
  factory QueryLockedVestingRequest() => create();

  QueryLockedVestingRequest._();

  factory QueryLockedVestingRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryLockedVestingRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryLockedVestingRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLockedVestingRequest clone() => QueryLockedVestingRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLockedVestingRequest copyWith(void Function(QueryLockedVestingRequest) updates) => super.copyWith((message) => updates(message as QueryLockedVestingRequest)) as QueryLockedVestingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLockedVestingRequest create() => QueryLockedVestingRequest._();
  @$core.override
  QueryLockedVestingRequest createEmptyInstance() => create();
  static $pb.PbList<QueryLockedVestingRequest> createRepeated() => $pb.PbList<QueryLockedVestingRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryLockedVestingRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryLockedVestingRequest>(create);
  static QueryLockedVestingRequest? _defaultInstance;
}

class QueryLockedVestingResponse extends $pb.GeneratedMessage {
  factory QueryLockedVestingResponse({
    $0.Coin? amount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    return result;
  }

  QueryLockedVestingResponse._();

  factory QueryLockedVestingResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryLockedVestingResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryLockedVestingResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOM<$0.Coin>(1, _omitFieldNames ? '' : 'amount', subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLockedVestingResponse clone() => QueryLockedVestingResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLockedVestingResponse copyWith(void Function(QueryLockedVestingResponse) updates) => super.copyWith((message) => updates(message as QueryLockedVestingResponse)) as QueryLockedVestingResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLockedVestingResponse create() => QueryLockedVestingResponse._();
  @$core.override
  QueryLockedVestingResponse createEmptyInstance() => create();
  static $pb.PbList<QueryLockedVestingResponse> createRepeated() => $pb.PbList<QueryLockedVestingResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryLockedVestingResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryLockedVestingResponse>(create);
  static QueryLockedVestingResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Coin get amount => $_getN(0);
  @$pb.TagNumber(1)
  set amount($0.Coin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Coin ensureAmount() => $_ensure(0);
}

class QueryCirculatingSupplyRequest extends $pb.GeneratedMessage {
  factory QueryCirculatingSupplyRequest() => create();

  QueryCirculatingSupplyRequest._();

  factory QueryCirculatingSupplyRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryCirculatingSupplyRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryCirculatingSupplyRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryCirculatingSupplyRequest clone() => QueryCirculatingSupplyRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryCirculatingSupplyRequest copyWith(void Function(QueryCirculatingSupplyRequest) updates) => super.copyWith((message) => updates(message as QueryCirculatingSupplyRequest)) as QueryCirculatingSupplyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryCirculatingSupplyRequest create() => QueryCirculatingSupplyRequest._();
  @$core.override
  QueryCirculatingSupplyRequest createEmptyInstance() => create();
  static $pb.PbList<QueryCirculatingSupplyRequest> createRepeated() => $pb.PbList<QueryCirculatingSupplyRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryCirculatingSupplyRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryCirculatingSupplyRequest>(create);
  static QueryCirculatingSupplyRequest? _defaultInstance;
}

class QueryCirculatingSupplyResponse extends $pb.GeneratedMessage {
  factory QueryCirculatingSupplyResponse({
    $0.Coin? amount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    return result;
  }

  QueryCirculatingSupplyResponse._();

  factory QueryCirculatingSupplyResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryCirculatingSupplyResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryCirculatingSupplyResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOM<$0.Coin>(1, _omitFieldNames ? '' : 'amount', subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryCirculatingSupplyResponse clone() => QueryCirculatingSupplyResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryCirculatingSupplyResponse copyWith(void Function(QueryCirculatingSupplyResponse) updates) => super.copyWith((message) => updates(message as QueryCirculatingSupplyResponse)) as QueryCirculatingSupplyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryCirculatingSupplyResponse create() => QueryCirculatingSupplyResponse._();
  @$core.override
  QueryCirculatingSupplyResponse createEmptyInstance() => create();
  static $pb.PbList<QueryCirculatingSupplyResponse> createRepeated() => $pb.PbList<QueryCirculatingSupplyResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryCirculatingSupplyResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryCirculatingSupplyResponse>(create);
  static QueryCirculatingSupplyResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Coin get amount => $_getN(0);
  @$pb.TagNumber(1)
  set amount($0.Coin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Coin ensureAmount() => $_ensure(0);
}

class QueryTotalSupplyRequest extends $pb.GeneratedMessage {
  factory QueryTotalSupplyRequest() => create();

  QueryTotalSupplyRequest._();

  factory QueryTotalSupplyRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalSupplyRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalSupplyRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalSupplyRequest clone() => QueryTotalSupplyRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalSupplyRequest copyWith(void Function(QueryTotalSupplyRequest) updates) => super.copyWith((message) => updates(message as QueryTotalSupplyRequest)) as QueryTotalSupplyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalSupplyRequest create() => QueryTotalSupplyRequest._();
  @$core.override
  QueryTotalSupplyRequest createEmptyInstance() => create();
  static $pb.PbList<QueryTotalSupplyRequest> createRepeated() => $pb.PbList<QueryTotalSupplyRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalSupplyRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalSupplyRequest>(create);
  static QueryTotalSupplyRequest? _defaultInstance;
}

class QueryTotalSupplyResponse extends $pb.GeneratedMessage {
  factory QueryTotalSupplyResponse({
    $0.Coin? amount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    return result;
  }

  QueryTotalSupplyResponse._();

  factory QueryTotalSupplyResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalSupplyResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalSupplyResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOM<$0.Coin>(1, _omitFieldNames ? '' : 'amount', subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalSupplyResponse clone() => QueryTotalSupplyResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalSupplyResponse copyWith(void Function(QueryTotalSupplyResponse) updates) => super.copyWith((message) => updates(message as QueryTotalSupplyResponse)) as QueryTotalSupplyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalSupplyResponse create() => QueryTotalSupplyResponse._();
  @$core.override
  QueryTotalSupplyResponse createEmptyInstance() => create();
  static $pb.PbList<QueryTotalSupplyResponse> createRepeated() => $pb.PbList<QueryTotalSupplyResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalSupplyResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalSupplyResponse>(create);
  static QueryTotalSupplyResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $0.Coin get amount => $_getN(0);
  @$pb.TagNumber(1)
  set amount($0.Coin value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => $_clearField(1);
  @$pb.TagNumber(1)
  $0.Coin ensureAmount() => $_ensure(0);
}

class QueryEstimatedGasPriceRequest extends $pb.GeneratedMessage {
  factory QueryEstimatedGasPriceRequest() => create();

  QueryEstimatedGasPriceRequest._();

  factory QueryEstimatedGasPriceRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryEstimatedGasPriceRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryEstimatedGasPriceRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryEstimatedGasPriceRequest clone() => QueryEstimatedGasPriceRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryEstimatedGasPriceRequest copyWith(void Function(QueryEstimatedGasPriceRequest) updates) => super.copyWith((message) => updates(message as QueryEstimatedGasPriceRequest)) as QueryEstimatedGasPriceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryEstimatedGasPriceRequest create() => QueryEstimatedGasPriceRequest._();
  @$core.override
  QueryEstimatedGasPriceRequest createEmptyInstance() => create();
  static $pb.PbList<QueryEstimatedGasPriceRequest> createRepeated() => $pb.PbList<QueryEstimatedGasPriceRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryEstimatedGasPriceRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryEstimatedGasPriceRequest>(create);
  static QueryEstimatedGasPriceRequest? _defaultInstance;
}

class QueryEstimatedGasPriceResponse extends $pb.GeneratedMessage {
  factory QueryEstimatedGasPriceResponse({
    $core.String? gasPrice,
  }) {
    final result = create();
    if (gasPrice != null) result.gasPrice = gasPrice;
    return result;
  }

  QueryEstimatedGasPriceResponse._();

  factory QueryEstimatedGasPriceResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryEstimatedGasPriceResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryEstimatedGasPriceResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'gasPrice')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryEstimatedGasPriceResponse clone() => QueryEstimatedGasPriceResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryEstimatedGasPriceResponse copyWith(void Function(QueryEstimatedGasPriceResponse) updates) => super.copyWith((message) => updates(message as QueryEstimatedGasPriceResponse)) as QueryEstimatedGasPriceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryEstimatedGasPriceResponse create() => QueryEstimatedGasPriceResponse._();
  @$core.override
  QueryEstimatedGasPriceResponse createEmptyInstance() => create();
  static $pb.PbList<QueryEstimatedGasPriceResponse> createRepeated() => $pb.PbList<QueryEstimatedGasPriceResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryEstimatedGasPriceResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryEstimatedGasPriceResponse>(create);
  static QueryEstimatedGasPriceResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get gasPrice => $_getSZ(0);
  @$pb.TagNumber(1)
  set gasPrice($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasGasPrice() => $_has(0);
  @$pb.TagNumber(1)
  void clearGasPrice() => $_clearField(1);
}

class QueryLastBlockGasUsedRequest extends $pb.GeneratedMessage {
  factory QueryLastBlockGasUsedRequest() => create();

  QueryLastBlockGasUsedRequest._();

  factory QueryLastBlockGasUsedRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryLastBlockGasUsedRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryLastBlockGasUsedRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLastBlockGasUsedRequest clone() => QueryLastBlockGasUsedRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLastBlockGasUsedRequest copyWith(void Function(QueryLastBlockGasUsedRequest) updates) => super.copyWith((message) => updates(message as QueryLastBlockGasUsedRequest)) as QueryLastBlockGasUsedRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLastBlockGasUsedRequest create() => QueryLastBlockGasUsedRequest._();
  @$core.override
  QueryLastBlockGasUsedRequest createEmptyInstance() => create();
  static $pb.PbList<QueryLastBlockGasUsedRequest> createRepeated() => $pb.PbList<QueryLastBlockGasUsedRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryLastBlockGasUsedRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryLastBlockGasUsedRequest>(create);
  static QueryLastBlockGasUsedRequest? _defaultInstance;
}

class QueryLastBlockGasUsedResponse extends $pb.GeneratedMessage {
  factory QueryLastBlockGasUsedResponse({
    $fixnum.Int64? gasUsed,
  }) {
    final result = create();
    if (gasUsed != null) result.gasUsed = gasUsed;
    return result;
  }

  QueryLastBlockGasUsedResponse._();

  factory QueryLastBlockGasUsedResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryLastBlockGasUsedResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryLastBlockGasUsedResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'gasUsed', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLastBlockGasUsedResponse clone() => QueryLastBlockGasUsedResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryLastBlockGasUsedResponse copyWith(void Function(QueryLastBlockGasUsedResponse) updates) => super.copyWith((message) => updates(message as QueryLastBlockGasUsedResponse)) as QueryLastBlockGasUsedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryLastBlockGasUsedResponse create() => QueryLastBlockGasUsedResponse._();
  @$core.override
  QueryLastBlockGasUsedResponse createEmptyInstance() => create();
  static $pb.PbList<QueryLastBlockGasUsedResponse> createRepeated() => $pb.PbList<QueryLastBlockGasUsedResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryLastBlockGasUsedResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryLastBlockGasUsedResponse>(create);
  static QueryLastBlockGasUsedResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get gasUsed => $_getI64(0);
  @$pb.TagNumber(1)
  set gasUsed($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasGasUsed() => $_has(0);
  @$pb.TagNumber(1)
  void clearGasUsed() => $_clearField(1);
}

class QueryTotalTxsRequest extends $pb.GeneratedMessage {
  factory QueryTotalTxsRequest() => create();

  QueryTotalTxsRequest._();

  factory QueryTotalTxsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalTxsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalTxsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalTxsRequest clone() => QueryTotalTxsRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalTxsRequest copyWith(void Function(QueryTotalTxsRequest) updates) => super.copyWith((message) => updates(message as QueryTotalTxsRequest)) as QueryTotalTxsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalTxsRequest create() => QueryTotalTxsRequest._();
  @$core.override
  QueryTotalTxsRequest createEmptyInstance() => create();
  static $pb.PbList<QueryTotalTxsRequest> createRepeated() => $pb.PbList<QueryTotalTxsRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalTxsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalTxsRequest>(create);
  static QueryTotalTxsRequest? _defaultInstance;
}

class QueryTotalTxsResponse extends $pb.GeneratedMessage {
  factory QueryTotalTxsResponse({
    $fixnum.Int64? totalTxs,
  }) {
    final result = create();
    if (totalTxs != null) result.totalTxs = totalTxs;
    return result;
  }

  QueryTotalTxsResponse._();

  factory QueryTotalTxsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalTxsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalTxsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'totalTxs', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalTxsResponse clone() => QueryTotalTxsResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalTxsResponse copyWith(void Function(QueryTotalTxsResponse) updates) => super.copyWith((message) => updates(message as QueryTotalTxsResponse)) as QueryTotalTxsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalTxsResponse create() => QueryTotalTxsResponse._();
  @$core.override
  QueryTotalTxsResponse createEmptyInstance() => create();
  static $pb.PbList<QueryTotalTxsResponse> createRepeated() => $pb.PbList<QueryTotalTxsResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalTxsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalTxsResponse>(create);
  static QueryTotalTxsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get totalTxs => $_getI64(0);
  @$pb.TagNumber(1)
  set totalTxs($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTotalTxs() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalTxs() => $_clearField(1);
}

class UnlockSchedule extends $pb.GeneratedMessage {
  factory UnlockSchedule({
    $core.String? address,
    $core.String? timeStr,
    $fixnum.Int64? timeUnix,
    $fixnum.Int64? amount,
    $core.String? denom,
  }) {
    final result = create();
    if (address != null) result.address = address;
    if (timeStr != null) result.timeStr = timeStr;
    if (timeUnix != null) result.timeUnix = timeUnix;
    if (amount != null) result.amount = amount;
    if (denom != null) result.denom = denom;
    return result;
  }

  UnlockSchedule._();

  factory UnlockSchedule.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory UnlockSchedule.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UnlockSchedule', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'timeStr')
    ..aInt64(3, _omitFieldNames ? '' : 'timeUnix')
    ..aInt64(4, _omitFieldNames ? '' : 'amount')
    ..aOS(5, _omitFieldNames ? '' : 'denom')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlockSchedule clone() => UnlockSchedule()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnlockSchedule copyWith(void Function(UnlockSchedule) updates) => super.copyWith((message) => updates(message as UnlockSchedule)) as UnlockSchedule;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnlockSchedule create() => UnlockSchedule._();
  @$core.override
  UnlockSchedule createEmptyInstance() => create();
  static $pb.PbList<UnlockSchedule> createRepeated() => $pb.PbList<UnlockSchedule>();
  @$core.pragma('dart2js:noInline')
  static UnlockSchedule getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnlockSchedule>(create);
  static UnlockSchedule? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get timeStr => $_getSZ(1);
  @$pb.TagNumber(2)
  set timeStr($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTimeStr() => $_has(1);
  @$pb.TagNumber(2)
  void clearTimeStr() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get timeUnix => $_getI64(2);
  @$pb.TagNumber(3)
  set timeUnix($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTimeUnix() => $_has(2);
  @$pb.TagNumber(3)
  void clearTimeUnix() => $_clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get amount => $_getI64(3);
  @$pb.TagNumber(4)
  set amount($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAmount() => $_has(3);
  @$pb.TagNumber(4)
  void clearAmount() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.String get denom => $_getSZ(4);
  @$pb.TagNumber(5)
  set denom($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasDenom() => $_has(4);
  @$pb.TagNumber(5)
  void clearDenom() => $_clearField(5);
}

class QueryUnlockSchedulesRequest extends $pb.GeneratedMessage {
  factory QueryUnlockSchedulesRequest() => create();

  QueryUnlockSchedulesRequest._();

  factory QueryUnlockSchedulesRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryUnlockSchedulesRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryUnlockSchedulesRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryUnlockSchedulesRequest clone() => QueryUnlockSchedulesRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryUnlockSchedulesRequest copyWith(void Function(QueryUnlockSchedulesRequest) updates) => super.copyWith((message) => updates(message as QueryUnlockSchedulesRequest)) as QueryUnlockSchedulesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryUnlockSchedulesRequest create() => QueryUnlockSchedulesRequest._();
  @$core.override
  QueryUnlockSchedulesRequest createEmptyInstance() => create();
  static $pb.PbList<QueryUnlockSchedulesRequest> createRepeated() => $pb.PbList<QueryUnlockSchedulesRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryUnlockSchedulesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryUnlockSchedulesRequest>(create);
  static QueryUnlockSchedulesRequest? _defaultInstance;
}

class QueryUnlockSchedulesResponse extends $pb.GeneratedMessage {
  factory QueryUnlockSchedulesResponse({
    $core.Iterable<UnlockSchedule>? unlockSchedules,
  }) {
    final result = create();
    if (unlockSchedules != null) result.unlockSchedules.addAll(unlockSchedules);
    return result;
  }

  QueryUnlockSchedulesResponse._();

  factory QueryUnlockSchedulesResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryUnlockSchedulesResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryUnlockSchedulesResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..pc<UnlockSchedule>(1, _omitFieldNames ? '' : 'unlockSchedules', $pb.PbFieldType.PM, subBuilder: UnlockSchedule.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryUnlockSchedulesResponse clone() => QueryUnlockSchedulesResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryUnlockSchedulesResponse copyWith(void Function(QueryUnlockSchedulesResponse) updates) => super.copyWith((message) => updates(message as QueryUnlockSchedulesResponse)) as QueryUnlockSchedulesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryUnlockSchedulesResponse create() => QueryUnlockSchedulesResponse._();
  @$core.override
  QueryUnlockSchedulesResponse createEmptyInstance() => create();
  static $pb.PbList<QueryUnlockSchedulesResponse> createRepeated() => $pb.PbList<QueryUnlockSchedulesResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryUnlockSchedulesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryUnlockSchedulesResponse>(create);
  static QueryUnlockSchedulesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<UnlockSchedule> get unlockSchedules => $_getList(0);
}

class QueryParamsRequest extends $pb.GeneratedMessage {
  factory QueryParamsRequest() => create();

  QueryParamsRequest._();

  factory QueryParamsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
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
    $fixnum.Int64? extraGasPerNewAccount,
  }) {
    final result = create();
    if (extraGasPerNewAccount != null) result.extraGasPerNewAccount = extraGasPerNewAccount;
    return result;
  }

  QueryParamsResponse._();

  factory QueryParamsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.paxi.types'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'extraGasPerNewAccount', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
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
  $fixnum.Int64 get extraGasPerNewAccount => $_getI64(0);
  @$pb.TagNumber(1)
  set extraGasPerNewAccount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasExtraGasPerNewAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearExtraGasPerNewAccount() => $_clearField(1);
}

class QueryApi {
  final $pb.RpcClient _client;

  QueryApi(this._client);

  $async.Future<QueryLockedVestingResponse> lockedVesting($pb.ClientContext? ctx, QueryLockedVestingRequest request) =>
    _client.invoke<QueryLockedVestingResponse>(ctx, 'Query', 'LockedVesting', request, QueryLockedVestingResponse())
  ;
  $async.Future<QueryCirculatingSupplyResponse> circulatingSupply($pb.ClientContext? ctx, QueryCirculatingSupplyRequest request) =>
    _client.invoke<QueryCirculatingSupplyResponse>(ctx, 'Query', 'CirculatingSupply', request, QueryCirculatingSupplyResponse())
  ;
  $async.Future<QueryTotalSupplyResponse> totalSupply($pb.ClientContext? ctx, QueryTotalSupplyRequest request) =>
    _client.invoke<QueryTotalSupplyResponse>(ctx, 'Query', 'TotalSupply', request, QueryTotalSupplyResponse())
  ;
  $async.Future<QueryEstimatedGasPriceResponse> estimatedGasPrice($pb.ClientContext? ctx, QueryEstimatedGasPriceRequest request) =>
    _client.invoke<QueryEstimatedGasPriceResponse>(ctx, 'Query', 'EstimatedGasPrice', request, QueryEstimatedGasPriceResponse())
  ;
  $async.Future<QueryLastBlockGasUsedResponse> lastBlockGasUsed($pb.ClientContext? ctx, QueryLastBlockGasUsedRequest request) =>
    _client.invoke<QueryLastBlockGasUsedResponse>(ctx, 'Query', 'LastBlockGasUsed', request, QueryLastBlockGasUsedResponse())
  ;
  $async.Future<QueryTotalTxsResponse> totalTxs($pb.ClientContext? ctx, QueryTotalTxsRequest request) =>
    _client.invoke<QueryTotalTxsResponse>(ctx, 'Query', 'TotalTxs', request, QueryTotalTxsResponse())
  ;
  $async.Future<QueryUnlockSchedulesResponse> unlockSchedules($pb.ClientContext? ctx, QueryUnlockSchedulesRequest request) =>
    _client.invoke<QueryUnlockSchedulesResponse>(ctx, 'Query', 'UnlockSchedules', request, QueryUnlockSchedulesResponse())
  ;
  $async.Future<QueryParamsResponse> params($pb.ClientContext? ctx, QueryParamsRequest request) =>
    _client.invoke<QueryParamsResponse>(ctx, 'Query', 'Params', request, QueryParamsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
