// This is a generated file - do not edit.
//
// Generated from x/custommint/types/query.proto.

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

class QueryTotalMintedRequest extends $pb.GeneratedMessage {
  factory QueryTotalMintedRequest() => create();

  QueryTotalMintedRequest._();

  factory QueryTotalMintedRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalMintedRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalMintedRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalMintedRequest clone() => QueryTotalMintedRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalMintedRequest copyWith(void Function(QueryTotalMintedRequest) updates) => super.copyWith((message) => updates(message as QueryTotalMintedRequest)) as QueryTotalMintedRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalMintedRequest create() => QueryTotalMintedRequest._();
  @$core.override
  QueryTotalMintedRequest createEmptyInstance() => create();
  static $pb.PbList<QueryTotalMintedRequest> createRepeated() => $pb.PbList<QueryTotalMintedRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalMintedRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalMintedRequest>(create);
  static QueryTotalMintedRequest? _defaultInstance;
}

class QueryTotalMintedResponse extends $pb.GeneratedMessage {
  factory QueryTotalMintedResponse({
    $0.Coin? amount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    return result;
  }

  QueryTotalMintedResponse._();

  factory QueryTotalMintedResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalMintedResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalMintedResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..aOM<$0.Coin>(1, _omitFieldNames ? '' : 'amount', subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalMintedResponse clone() => QueryTotalMintedResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalMintedResponse copyWith(void Function(QueryTotalMintedResponse) updates) => super.copyWith((message) => updates(message as QueryTotalMintedResponse)) as QueryTotalMintedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalMintedResponse create() => QueryTotalMintedResponse._();
  @$core.override
  QueryTotalMintedResponse createEmptyInstance() => create();
  static $pb.PbList<QueryTotalMintedResponse> createRepeated() => $pb.PbList<QueryTotalMintedResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalMintedResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalMintedResponse>(create);
  static QueryTotalMintedResponse? _defaultInstance;

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

class QueryTotalBurnedRequest extends $pb.GeneratedMessage {
  factory QueryTotalBurnedRequest() => create();

  QueryTotalBurnedRequest._();

  factory QueryTotalBurnedRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalBurnedRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalBurnedRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalBurnedRequest clone() => QueryTotalBurnedRequest()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalBurnedRequest copyWith(void Function(QueryTotalBurnedRequest) updates) => super.copyWith((message) => updates(message as QueryTotalBurnedRequest)) as QueryTotalBurnedRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalBurnedRequest create() => QueryTotalBurnedRequest._();
  @$core.override
  QueryTotalBurnedRequest createEmptyInstance() => create();
  static $pb.PbList<QueryTotalBurnedRequest> createRepeated() => $pb.PbList<QueryTotalBurnedRequest>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalBurnedRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalBurnedRequest>(create);
  static QueryTotalBurnedRequest? _defaultInstance;
}

class QueryTotalBurnedResponse extends $pb.GeneratedMessage {
  factory QueryTotalBurnedResponse({
    $0.Coin? amount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    return result;
  }

  QueryTotalBurnedResponse._();

  factory QueryTotalBurnedResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryTotalBurnedResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryTotalBurnedResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..aOM<$0.Coin>(1, _omitFieldNames ? '' : 'amount', subBuilder: $0.Coin.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalBurnedResponse clone() => QueryTotalBurnedResponse()..mergeFromMessage(this);
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  QueryTotalBurnedResponse copyWith(void Function(QueryTotalBurnedResponse) updates) => super.copyWith((message) => updates(message as QueryTotalBurnedResponse)) as QueryTotalBurnedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static QueryTotalBurnedResponse create() => QueryTotalBurnedResponse._();
  @$core.override
  QueryTotalBurnedResponse createEmptyInstance() => create();
  static $pb.PbList<QueryTotalBurnedResponse> createRepeated() => $pb.PbList<QueryTotalBurnedResponse>();
  @$core.pragma('dart2js:noInline')
  static QueryTotalBurnedResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<QueryTotalBurnedResponse>(create);
  static QueryTotalBurnedResponse? _defaultInstance;

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

class QueryParamsRequest extends $pb.GeneratedMessage {
  factory QueryParamsRequest() => create();

  QueryParamsRequest._();

  factory QueryParamsRequest.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsRequest.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
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

  QueryParamsResponse._();

  factory QueryParamsResponse.fromBuffer($core.List<$core.int> data, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(data, registry);
  factory QueryParamsResponse.fromJson($core.String json, [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'QueryParamsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'x.custommint.types'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'burnThreshold')
    ..aOS(2, _omitFieldNames ? '' : 'burnRatio')
    ..aInt64(3, _omitFieldNames ? '' : 'blocksPerYear')
    ..aOS(4, _omitFieldNames ? '' : 'firstYearInflation')
    ..aOS(5, _omitFieldNames ? '' : 'secondYearInflation')
    ..aOS(6, _omitFieldNames ? '' : 'otherYearInflation')
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

class QueryApi {
  final $pb.RpcClient _client;

  QueryApi(this._client);

  $async.Future<QueryTotalMintedResponse> totalMinted($pb.ClientContext? ctx, QueryTotalMintedRequest request) =>
    _client.invoke<QueryTotalMintedResponse>(ctx, 'Query', 'TotalMinted', request, QueryTotalMintedResponse())
  ;
  $async.Future<QueryTotalBurnedResponse> totalBurned($pb.ClientContext? ctx, QueryTotalBurnedRequest request) =>
    _client.invoke<QueryTotalBurnedResponse>(ctx, 'Query', 'TotalBurned', request, QueryTotalBurnedResponse())
  ;
  $async.Future<QueryParamsResponse> params($pb.ClientContext? ctx, QueryParamsRequest request) =>
    _client.invoke<QueryParamsResponse>(ctx, 'Query', 'Params', request, QueryParamsResponse())
  ;
}


const $core.bool _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
